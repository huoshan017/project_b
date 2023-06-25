package common

import (
	"project_b/common/base"
	"project_b/common/ds"
	"project_b/common/log"
	"project_b/common/math"
	"project_b/common/object"
	"project_b/common/time"
	"project_b/common_data"
	"project_b/game_map"
	"project_b/utils"

	"github.com/huoshan017/ponu/heap"
)

// 场景圖结构必须在单个goroutine中执行
type SceneMap struct {
	mapConfig           *game_map.Config
	mapWidth, mapHeight int32
	eventMgr            base.IEventManager
	playerTankList      *ds.MapListUnion[uint64, *object.Tank]
	enemyTankList       *ds.MapListUnion[uint32, *object.Tank]
	playerTankListCache []PlayerTankKV
	objFactory          *object.ObjectFactory
	pmap                *PartitionMap
}

func NewSceneMap(eventMgr base.IEventManager) *SceneMap {
	return &SceneMap{
		eventMgr:       eventMgr,
		playerTankList: ds.NewMapListUnion[uint64, *object.Tank](),
		enemyTankList:  ds.NewMapListUnion[uint32, *object.Tank](),
		objFactory:     object.NewObjectFactory(true),
		pmap:           NewPartitionMap(0),
	}
}

func (s *SceneMap) GetMapId() int32 {
	return s.mapConfig.Id
}

func (s *SceneMap) LoadMap(m *game_map.Config) bool {
	// 载入地图
	s.pmap.Load(m)
	// 地图载入前事件
	s.eventMgr.InvokeEvent(EventIdBeforeMapLoad)
	for line := 0; line < len(m.Layers); line++ {
		for col := 0; col < len(m.Layers[line]); col++ {
			st := object.StaticObjType(m.Layers[line][col])
			if common_data.StaticObjectConfigData[st] == nil {
				continue
			}
			tileObj := s.objFactory.NewStaticObject(common_data.StaticObjectConfigData[st])
			// 碰撞組件
			if common_data.StaticObjectConfigData[st].Collision() {
				tileObj.AddComp(&object.CollisionComp{})
			}
			// 二維數組Y軸是自上而下的，而世界坐標Y軸是自下而上的，所以設置Y坐標要倒過來
			tileObj.SetPos(m.TileWidth*int32(col), m.TileHeight*int32(len(m.Layers)-1-line))
			// 加入網格分區地圖
			s.pmap.AddTile(int16(len(m.Layers)-1-line), int16(col), tileObj)
		}
	}
	s.mapConfig = m
	s.mapWidth = int32(len(m.Layers[0])) * m.TileWidth
	s.mapHeight = int32(len(m.Layers)) * m.TileHeight
	// 地图载入完成事件
	s.eventMgr.InvokeEvent(EventIdMapLoaded, s)
	log.Info("Load map %v done, map width %v, map height %v", m.Id, s.mapWidth, s.mapHeight)
	return true
}

func (s *SceneMap) UnloadMap() {
	// 地图卸载前事件
	s.eventMgr.InvokeEvent(EventIdBeforeMapUnload)
	s.mapWidth = 0
	s.mapHeight = 0
	s.playerTankList.Clear()
	s.enemyTankList.Clear()
	s.playerTankListCache = s.playerTankListCache[:0]
	s.pmap.Unload()
	s.objFactory.Clear()
	// 地图卸载后事件
	s.eventMgr.InvokeEvent(EventIdMapUnloaded)
}

func (s *SceneMap) GetMapConfig() *game_map.Config {
	return s.mapConfig
}

func (s *SceneMap) GetLayerObjsWithRange(rect *math.Rect) [MapMaxLayer]*heap.BinaryHeapKV[uint32, int32] {
	return s.pmap.GetLayerObjsWithRange(rect)
}

func (s *SceneMap) GetObj(instId uint32) object.IObject {
	return s.objFactory.GetObj(instId)
}

func (s *SceneMap) GetPlayerTank(pid uint64) *object.Tank {
	v, o := s.playerTankList.Get(pid)
	if !o {
		return nil
	}
	return v
}

func (s *SceneMap) NewPlayerTank(pid uint64) *object.Tank {
	tank := s.objFactory.NewTank(&s.mapConfig.PlayerTankInitData)
	// 注冊檢測移動事件處理
	tank.RegisterCheckMoveEventHandle(s.checkObjMoveEventHandle)
	// 設置碰撞組件
	if common_data.TankConfigData[s.mapConfig.PlayerTankInitData.Id()].Collision() {
		tank.AddComp(&object.CollisionComp{})
	}
	// 随机并设置坦克位置
	pos := utils.RandomPosInRect(s.mapConfig.PlayerTankInitRect)
	tank.SetPos(pos.X, pos.Y)
	// 加入到玩家坦克列表
	s.playerTankList.Add(pid, tank)
	// 加入網格分區地圖
	s.pmap.AddObj(tank)
	return tank
}

func (s *SceneMap) AddPlayerTank(pid uint64, tank *object.Tank) {
	tank.RegisterCheckMoveEventHandle(s.checkObjMoveEventHandle)
	s.playerTankList.Add(pid, tank)
	s.pmap.AddObj(tank)
}

func (s *SceneMap) AddPlayerTankWithInfo(pid uint64, id int32, level int32, x, y int32, dir object.Direction, currSpeed int32) {
	tank := s.objFactory.NewTank(common_data.TankConfigData[id])
	tank.SetPos(x, y)
	tank.SetLevel(level)
	tank.SetDir(dir)
	tank.SetCurrentSpeed(currSpeed)
	// 注冊檢測移動事件處理
	tank.RegisterCheckMoveEventHandle(s.checkObjMoveEventHandle)
	// 設置碰撞組件
	if common_data.TankConfigData[id].Collision() {
		tank.AddComp(&object.CollisionComp{})
	}
	// 加入玩家坦克列表
	s.playerTankList.Add(pid, tank)
	// 加入網格分區地圖
	s.pmap.AddObj(tank)
}

func (s *SceneMap) RemovePlayerTank(pid uint64) {
	tank := s.playerTankList.Remove(pid)
	tank.UnregisterCheckMoveEventHandle(s.checkObjMoveEventHandle)
	s.pmap.RemoveObj(tank.InstId())
	s.objFactory.RecycleTank(tank)
}

func (s *SceneMap) GetPlayerTankList() []PlayerTankKV {
	if s.playerTankListCache != nil && len(s.playerTankListCache) > 0 {
		s.playerTankListCache = s.playerTankListCache[:0]
	}
	count := s.playerTankList.Count()
	for i := int32(0); i < count; i++ {
		k, v := s.playerTankList.GetByIndex(i)
		s.playerTankListCache = append(s.playerTankListCache, PlayerTankKV{
			PlayerId: k,
			Tank:     v,
		})
	}
	return s.playerTankListCache
}

func (s *SceneMap) PlayerTankMove(uid uint64, dir object.Direction) {
	tank := s.GetPlayerTank(uid)
	if tank == nil {
		log.Error("player %v tank not found", uid)
		return
	}
	tank.SetDir(dir)
	tank.Move(dir)
}

func (s *SceneMap) PlayerTankStopMove(uid uint64) {
	tank := s.GetPlayerTank(uid)
	if tank == nil {
		log.Error("player %v tank not found", uid)
		return
	}
	tank.Stop()
}

func (s *SceneMap) PlayerTankChange(uid uint64, staticInfo *object.ObjStaticInfo) bool {
	tank := s.GetPlayerTank(uid)
	if tank == nil {
		return false
	}
	tank.Change(staticInfo)
	return true
}

func (s *SceneMap) PlayerTankRestore(uid uint64) int32 {
	tank := s.GetPlayerTank(uid)
	if tank == nil {
		return 0
	}
	tank.Restore()
	return tank.Id()
}

func (s *SceneMap) AddEnemyTank(instId uint32, tank *object.Tank) {
	s.enemyTankList.Add(instId, tank)
	tank.RegisterCheckMoveEventHandle(s.checkObjMoveEventHandle)
}

func (s *SceneMap) GetEnemyTank(instId uint32) *object.Tank {
	v, o := s.enemyTankList.Get(instId)
	if !o {
		return nil
	}
	return v
}

func (s *SceneMap) RemoveEnemyTank(id uint32) {
	v, o := s.enemyTankList.Get(id)
	if o {
		v.UnregisterCheckMoveEventHandle(s.checkObjMoveEventHandle)
	}
	s.enemyTankList.Remove(id)
}

func (s *SceneMap) EnemyTankMove(id uint32, dir object.Direction) {
	tank := s.GetEnemyTank(id)
	if tank == nil {
		log.Error("enemy tank %v not found", id)
		return
	}
	tank.SetDir(dir)
	tank.Move(dir)
}

func (s *SceneMap) EnemyTankStopMove(id uint32) {
	tank := s.GetEnemyTank(id)
	if tank == nil {
		log.Error("enemy tank %v not found", id)
		return
	}
	tank.Stop()
}

func (s *SceneMap) Update(tick time.Duration) {
	count := s.playerTankList.Count()
	for i := int32(0); i < count; i++ {
		_, tank := s.playerTankList.GetByIndex(i)
		tank.Update(tick)
		s.pmap.UpdateMovable(tank)
	}
	count = s.enemyTankList.Count()
	for i := int32(0); i < count; i++ {
		_, tank := s.enemyTankList.GetByIndex(i)
		tank.Update(tick)
		s.pmap.UpdateMovable(tank)
	}
}

// 注册事件
func (s *SceneMap) RegisterEvent(eid base.EventId, handle func(args ...any)) {
	s.eventMgr.RegisterEvent(eid, handle)
}

// 注销事件
func (s *SceneMap) UnregisterEvent(eid base.EventId, handle func(args ...any)) {
	s.eventMgr.UnregisterEvent(eid, handle)
}

// 注册玩家事件
func (s *SceneMap) RegisterPlayerEvent(uid uint64, eid base.EventId, handle func(args ...any)) {
	tank := s.GetPlayerTank(uid)
	if tank == nil {
		return
	}
	switch eid {
	case EventIdTankMove:
		tank.RegisterMoveEventHandle(handle)
	case EventIdTankStopMove:
		tank.RegisterStopMoveEventHandle(handle)
	case EventIdTankSetPos:
		tank.RegisterUpdateEventHandle(handle)
	}
}

// 注销玩家事件
func (s *SceneMap) UnregisterPlayerEvent(uid uint64, eid base.EventId, handle func(args ...any)) {
	tank := s.GetPlayerTank(uid)
	if tank == nil {
		return
	}
	switch eid {
	case EventIdTankMove:
		tank.UnregisterMoveEventHandle(handle)
	case EventIdTankStopMove:
		tank.UnregisterStopMoveEventHandle(handle)
	case EventIdTankSetPos:
		tank.UnregisterUpdateEventHandle(handle)
	}
}

func (s *SceneMap) RegisterEnemyEvent(instId uint32, eid base.EventId, handle func(args ...any)) {
	tank := s.GetEnemyTank(instId)
	if tank == nil {
		return
	}
	switch eid {
	case EventIdTankMove:
		tank.RegisterMoveEventHandle(handle)
	case EventIdTankStopMove:
		tank.RegisterStopMoveEventHandle(handle)
	case EventIdTankSetPos:
		tank.RegisterUpdateEventHandle(handle)
	}
}

func (s *SceneMap) UnregisterEnemyEvent(instId uint32, eid base.EventId, handle func(args ...any)) {
	tank := s.GetEnemyTank(instId)
	if tank == nil {
		return
	}
	switch eid {
	case EventIdTankMove:
		tank.UnregisterMoveEventHandle(handle)
	case EventIdTankStopMove:
		tank.UnregisterStopMoveEventHandle(handle)
	case EventIdTankSetPos:
		tank.UnregisterUpdateEventHandle(handle)
	}
}

func (s *SceneMap) RegisterAllPlayersEvent(eid base.EventId, handle func(args ...any)) {
	count := s.playerTankList.Count()
	for i := int32(0); i < count; i++ {
		pid, _ := s.playerTankList.GetByIndex(i)
		s.RegisterPlayerEvent(pid, eid, handle)
	}
}

func (s *SceneMap) UnregisterAllPlayersEvent(eid base.EventId, handle func(args ...any)) {
	count := s.playerTankList.Count()
	for i := int32(0); i < count; i++ {
		pid, _ := s.playerTankList.GetByIndex(i)
		s.UnregisterPlayerEvent(pid, eid, handle)
	}
}

func (s *SceneMap) RegisterAllEnemiesEvent(eid base.EventId, handle func(args ...any)) {
	count := s.enemyTankList.Count()
	for i := int32(0); i < count; i++ {
		instId, _ := s.enemyTankList.GetByIndex(i)
		s.RegisterEnemyEvent(instId, eid, handle)
	}
}

func (s *SceneMap) UnregisterAllEnemiesEvent(eid base.EventId, handle func(args ...any)) {
	count := s.enemyTankList.Count()
	for i := int32(0); i < count; i++ {
		instId, _ := s.enemyTankList.GetByIndex(i)
		s.UnregisterEnemyEvent(instId, eid, handle)
	}
}

func (s *SceneMap) checkObjMoveEventHandle(args ...any) {
	obj := args[0].(object.IMovableObject)
	dir := args[1].(object.Direction)
	distance := args[2].(float64)
	isMove := args[3].(*bool)
	isCollision := args[4].(*bool)
	resObj := args[5].(*object.IObject)
	var (
		x, y int32
	)
	if !s.checkObjMoveRange(obj, dir, distance, &x, &y) {
		obj.SetPos(x, y)
		obj.Stop()
		*isMove = false
		*isCollision = false
	} else if s.checkMovableObjCollision(obj, dir, distance, resObj) {
		obj.Stop()
		*isCollision = true
		*isMove = false
	} else {
		*isMove = true
		*isCollision = false
	}
}

func (s *SceneMap) checkObjMoveRange(obj object.IMovableObject, dir object.Direction, distance float64, rx, ry *int32) bool {
	x, y := obj.Pos()
	var move bool = true
	switch dir {
	case object.DirLeft:
		if float64(x)-distance <= float64(s.mapConfig.X) {
			move = false
			x = s.mapConfig.X
		}
	case object.DirRight:
		if float64(x)+distance >= float64(s.mapConfig.X+s.mapWidth-obj.Width()) {
			move = false
			x = s.mapConfig.X + s.mapWidth - obj.Width()
		}
	case object.DirUp:
		if float64(y)+distance >= float64(s.mapConfig.Y+s.mapHeight-obj.Height()) {
			move = false
			y = s.mapConfig.Y + s.mapHeight - obj.Height()
		}
	case object.DirDown:
		if float64(y)-distance <= float64(s.mapConfig.Y) {
			move = false
			y = s.mapConfig.Y
		}
	}
	if !move {
		if rx != nil {
			*rx = x
		}
		if ry != nil {
			*ry = y
		}
	}
	return move
}

// 移動之前調用
func (s *SceneMap) checkMovableObjCollision(obj object.IMovableObject, dir object.Direction, distance float64, collisionObj *object.IObject) bool {
	// 是否擁有碰撞組件
	comp := obj.GetComp("Collision")
	if comp == nil {
		return false
	}

	// 獲取檢測碰撞範圍
	lx, by, rx, ty := s.pmap.objGridBounds(obj)
	if rx < lx || ty < by {
		return false
	}

	for y := by; y <= ty; y++ {
		for x := lx; x <= rx; x++ {
			gidx := s.pmap.gridLineCol2Index(y, x)
			lis := s.pmap.grids[gidx].getMObjs().GetList()
			for i := 0; i < len(lis); i++ {
				item := lis[i]
				obj2, o := s.pmap.mobjs.Get(item.Key)
				if !o {
					log.Warn("not found movable object %v", item.Key)
					continue
				}
				if obj2.InstId() != obj.InstId() {
					if s.checkMovableObjCollisionObj(obj, comp, dir, distance, obj2) {
						if collisionObj != nil {
							*collisionObj = obj2
						}
						return true
					}
				}
			}

			lis = s.pmap.grids[gidx].getSObjs().GetList()
			for i := 0; i < len(lis); i++ {
				item := lis[i]
				obj2, o := s.pmap.sobjs.Get(item.Key)
				if !o {
					log.Warn("not found static object %v", item.Key)
					continue
				}
				if s.checkMovableObjCollisionObj(obj, comp, dir, distance, obj2) {
					if collisionObj != nil {
						*collisionObj = obj2
					}
					return true
				}
			}
		}
	}
	return false
}

func (s *SceneMap) checkMovableObjCollisionObj(obj object.IMovableObject, comp object.IComponent, dir object.Direction, distance float64, obj2 object.IObject) bool {
	// 遍歷範圍内的所有碰撞物體
	var (
		collisionComp *object.CollisionComp
		aabb1         object.AABB
	)
	collisionComp = comp.(*object.CollisionComp)
	aabb1 = collisionComp.GetAABB(obj)
	aabb1.Move(dir, distance)
	comp2 := obj2.GetComp("Collision")
	if comp2 == nil {
		return false
	}
	collisionComp2 := comp2.(*object.CollisionComp)
	if collisionComp2 == nil {
		return false
	}
	aabb2 := collisionComp2.GetAABB(obj2)
	if aabb1.Intersect(&aabb2) {
		switch dir {
		case object.DirLeft:
			obj.SetPos(obj2.Right(), obj.Bottom())
		case object.DirRight:
			obj.SetPos(obj2.Left()-obj.Width(), obj.Bottom())
		case object.DirUp:
			obj.SetPos(obj.Left(), obj2.Bottom()-obj.Height())
		case object.DirDown:
			obj.SetPos(obj.Left(), obj2.Top())
		}
		return true
	}
	return false
}
