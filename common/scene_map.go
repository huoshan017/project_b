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
	gmap                *game_map.Config
	mapWidth, mapHeight int32
	eventMgr            base.IEventManager
	playerTankList      *ds.MapListUnion[uint64, *object.Tank]
	enemyTankList       *ds.MapListUnion[uint32, *object.Tank]
	playerTankListCache []PlayerTankKV
	objFactory          *object.ObjectFactory
	mapInstance         *game_map.MapInstance
}

func NewSceneMap(eventMgr base.IEventManager) *SceneMap {
	return &SceneMap{
		eventMgr:       eventMgr,
		playerTankList: ds.NewMapListUnion[uint64, *object.Tank](),
		enemyTankList:  ds.NewMapListUnion[uint32, *object.Tank](),
		objFactory:     object.NewObjectFactory(true),
		mapInstance:    game_map.NewMapInstance(0),
	}
}

func (s *SceneMap) GetMapId() int32 {
	return s.gmap.Id
}

func (s *SceneMap) LoadMap(m *game_map.Config) bool {
	// 载入地图
	s.mapInstance.Load(m)
	// 地图载入前事件
	s.eventMgr.InvokeEvent(EventIdBeforeMapLoad)
	for line := 0; line < len(m.Layers); line++ {
		for col := 0; col < len(m.Layers[line]); col++ {
			st := object.StaticObjType(m.Layers[line][col])
			if common_data.StaticObjectConfigData[st] == nil {
				continue
			}
			tileObj := s.objFactory.NewStaticObject(common_data.StaticObjectConfigData[st])
			// 二維數組Y軸是自上而下的，而世界坐標Y軸是自下而上的，所以設置Y坐標要倒過來
			tileObj.SetPos(m.TileWidth*int32(col), m.TileHeight*int32(len(m.Layers)-1-line))
			s.mapInstance.AddTile(int16(len(m.Layers)-1-line), int16(col), tileObj)
		}
	}
	s.gmap = m
	s.mapWidth = int32(len(m.Layers[0])) * m.TileWidth
	s.mapHeight = int32(len(m.Layers)) * m.TileHeight
	// 地图载入完成事件
	s.eventMgr.InvokeEvent(EventIdMapLoaded, s)
	log.Info("Load map %v done, map width %v, map height %v", m.Id, s.mapWidth, s.mapHeight)
	return true
}

func (s *SceneMap) UnloadMap() {
	/*if s.tileObjectArray == nil {
		return
	}
	// 地图卸载前事件
	s.eventMgr.InvokeEvent(EventIdBeforeMapUnload)
	for i := 0; i < len(s.tileObjectArray); i++ {
		if s.tileObjectArray[i] == nil {
			continue
		}
		for j := 0; j <= len(s.tileObjectArray[i]); j++ {
			s.objFactory.RecycleStaticObject(s.tileObjectArray[i][j])
			s.tileObjectArray[i][j] = nil
		}
	}
	// 地图卸载后事件
	s.eventMgr.InvokeEvent(EventIdMapUnloaded)
	s.tileObjectArray = nil*/
}

func (s *SceneMap) GetMapConfig() *game_map.Config {
	return s.gmap
}

func (s *SceneMap) GetLayerObjsWithRange(rect *math.Rect) [game_map.MapMaxLayer]*heap.BinaryHeapKV[uint32, int32] {
	return s.mapInstance.GetLayerObjsWithRange(rect)
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
	tank := s.objFactory.NewTank(&s.gmap.PlayerTankInitData)
	tank.RegisterCheckPosEventHandle(s.checkObjPosEventHandle)
	// 随机并设置坦克位置
	pos := utils.RandomPosInRect(s.gmap.PlayerTankInitRect)
	tank.SetPos(pos.X, pos.Y)
	// 加入到玩家坦克列表
	s.playerTankList.Add(pid, tank)
	s.mapInstance.AddObj(tank)
	return tank
}

func (s *SceneMap) AddPlayerTank(pid uint64, tank *object.Tank) {
	tank.RegisterCheckPosEventHandle(s.checkObjPosEventHandle)
	s.playerTankList.Add(pid, tank)
	s.mapInstance.AddObj(tank)
}

func (s *SceneMap) AddPlayerTankWithInfo(pid uint64, id int32, level int32, x, y int32, dir object.Direction, currSpeed int32) {
	tank := s.objFactory.NewTank(common_data.TankConfigData[id])
	tank.SetPos(x, y)
	tank.SetLevel(level)
	tank.SetDir(dir)
	tank.SetCurrentSpeed(currSpeed)
	tank.RegisterCheckPosEventHandle(s.checkObjPosEventHandle)
	s.playerTankList.Add(pid, tank)
	s.mapInstance.AddObj(tank)
}

func (s *SceneMap) RemovePlayerTank(pid uint64) {
	tank := s.playerTankList.Remove(pid)
	tank.UnregisterCheckPosEventHandle(s.checkObjPosEventHandle)
	s.mapInstance.RemoveObj(tank.InstId())
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
	if s.checkObjCanMove(tank, nil, nil) {
		tank.Move(dir)
	}
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
	tank.RegisterCheckPosEventHandle(s.checkObjPosEventHandle)
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
		v.UnregisterCheckPosEventHandle(s.checkObjPosEventHandle)
	}
	s.enemyTankList.Remove(id)
}

func (s *SceneMap) Update(tick time.Duration) {
	count := s.playerTankList.Count()
	for i := int32(0); i < count; i++ {
		_, tank := s.playerTankList.GetByIndex(i)
		tank.Update(tick)
		s.mapInstance.UpdateObj(tank)
	}
	count = s.enemyTankList.Count()
	for i := int32(0); i < count; i++ {
		_, tank := s.enemyTankList.GetByIndex(i)
		tank.Update(tick)
		s.mapInstance.UpdateObj(tank)
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

func (s *SceneMap) checkObjPosEventHandle(args ...any) {
	obj := args[0].(object.IMovableObject)
	var x, y int32
	if !s.checkObjCanMove(obj, &x, &y) {
		obj.SetPos(x, y)
		obj.Stop()
	}
}

func (s *SceneMap) checkObjCanMove(obj object.IMovableObject, rx, ry *int32) bool {
	x, y := obj.Pos()
	var move bool = true
	if x <= s.gmap.X && obj.Dir() == object.DirLeft {
		move = false
		x = s.gmap.X
	} else if x >= s.gmap.X+s.mapWidth-obj.Width() && obj.Dir() == object.DirRight {
		move = false
		x = s.gmap.X + s.mapWidth - obj.Width()
	}
	if y <= s.gmap.Y && obj.Dir() == object.DirDown {
		move = false
		y = s.gmap.Y
	} else if y >= s.gmap.Y+s.mapHeight-obj.Height() && obj.Dir() == object.DirUp {
		move = false
		y = s.gmap.Y + s.mapHeight - obj.Height()
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
