package common

import (
	"project_b/common/base"
	"project_b/common/ds"
	"project_b/common/effect"
	"project_b/common/math"
	"project_b/common/object"
	"project_b/common/weapon"
	"project_b/common_data"
	"project_b/game_map"
	"project_b/log"

	"github.com/huoshan017/ponu/heap"
)

type PlayerTankKV struct {
	PlayerId uint64
	Tank     *object.Tank
}

type TankBornPosInfo struct {
	x, y int32
	flag int16
}

// 场景圖，沒有玩家(Player)概念的游戲邏輯
// 必须在单个goroutine中执行
type World struct {
	mapConfig                            *game_map.Config                               // 地圖配置
	mapWidth, mapHeight                  int32                                          // 地圖寬高
	gmap                                 *GridMap                                       // 網格地圖
	tankBornPosList                      []TankBornPosInfo                              // 坦克出生位置信息列表
	eventMgr                             base.IEventManager                             // 事件管理器
	staticObjList                        *ds.MapListUnion[uint32, *object.StaticObject] // 靜態對象列表，用map和list的聯合體是爲了遍歷時的有序性
	tankList                             *ds.MapListUnion[uint32, *object.Tank]         // 坦克列表，不區分玩家和BOT
	shellList                            *ds.MapListUnion[uint32, *object.Shell]        // 炮彈列表，不區分坦克的炮彈
	surroundObjList                      *ds.MapListUnion[uint32, *object.SurroundObj]  // 環繞物體列表
	objFactory                           *object.ObjectFactory                          // 對象池
	effectList                           *ds.MapListUnion[uint32, *effect.Effect]       // 效果列表
	effectPool                           *effect.EffectPool                             // 效果池
	objAddedEvent, objRemovedEvent       base.Event                                     // obj添加和刪除事件
	effectAddedEvent, effectRemovedEvent base.Event                                     // 效果添加刪除事件
	staticObjRecycleList                 []*object.StaticObject                         // 靜態對象回收列表
	tankRecycleList                      []*object.Tank                                 // 坦克對象回收列表
	shellRecycleList                     []*object.Shell                                // 炮彈對象回收列表
	surroundObjRecycleList               []*object.SurroundObj                          // 環繞物體對象回收列表
	effectRecycleList                    []*effect.Effect                               // 效果回收列表
	effectSearchedList                   []uint32                                       // 效果搜索結果列表
}

func NewWorld(eventMgr base.IEventManager) *World {
	return &World{
		eventMgr:        eventMgr,
		staticObjList:   ds.NewMapListUnion[uint32, *object.StaticObject](),
		tankList:        ds.NewMapListUnion[uint32, *object.Tank](),
		shellList:       ds.NewMapListUnion[uint32, *object.Shell](),
		surroundObjList: ds.NewMapListUnion[uint32, *object.SurroundObj](),
		objFactory:      object.NewObjectFactory(true),
		effectList:      ds.NewMapListUnion[uint32, *effect.Effect](),
		effectPool:      effect.NewEffectPool(),
		gmap:            NewGridMap(1),
	}
}

func (s *World) GetMapId() int32 {
	return s.mapConfig.Id
}

func (s *World) LoadMap(m *game_map.Config) bool {
	// 载入地图
	s.gmap.Load(m)
	s.loadMap(m, false)
	s.mapConfig = m
	s.mapWidth = int32(len(m.Layers[0])) * m.TileWidth
	s.mapHeight = int32(len(m.Layers)) * m.TileHeight
	log.Info("Load map %v done, map width %v, map height %v", m.Id, s.mapWidth, s.mapHeight)
	return true
}

func (s *World) loadMap(m *game_map.Config, reload bool) {
	for line := 0; line < len(m.Layers); line++ {
		for col := 0; col < len(m.Layers[line]); col++ {
			st := base.StaticObjType(m.Layers[line][col])
			if common_data.StaticObjectConfigData[st] == nil {
				if reload {
					continue
				}
				s.tankBornPosList = append(s.tankBornPosList, TankBornPosInfo{
					x:    m.TileWidth*int32(col) + m.TileWidth/2,
					y:    m.TileHeight*int32(len(m.Layers)-1-line) + m.TileHeight/2,
					flag: m.Layers[line][col],
				})
				continue
			}
			tileObj := s.objFactory.NewStaticObject(common_data.StaticObjectConfigData[st])
			// 二維數組Y軸是自上而下的，而世界坐標Y軸是自下而上的，所以設置Y坐標要倒過來
			tileObj.SetPos(m.TileWidth*int32(col)+m.TileWidth/2, m.TileHeight*int32(len(m.Layers)-1-line)+m.TileHeight/2)
			s.objAddedEvent.Call(tileObj)
			// 加入網格分區地圖
			s.gmap.AddTile(int16(len(m.Layers)-1-line), int16(col), tileObj)
		}
	}
}

func (s *World) UnloadMap() {
	s.mapWidth = 0
	s.mapHeight = 0
	s.mapConfig = nil
	if len(s.tankBornPosList) > 0 {
		s.tankBornPosList = s.tankBornPosList[:0]
	}
	s.clearEvents()
	s.gmap.Unload()
	s.clearObjsData()
}

func (s *World) clearEvents() {
	s.objAddedEvent.Clear()
	s.objRemovedEvent.Clear()
	s.effectAddedEvent.Clear()
	s.effectRemovedEvent.Clear()
}

func (s *World) clearObjsData() {
	for i := int32(0); i < s.staticObjList.Count(); i++ {
		_, v := s.staticObjList.GetByIndex(i)
		if v != nil {
			s.objFactory.RecycleStaticObject(v)
		}
	}
	for i := int32(0); i < s.tankList.Count(); i++ {
		_, v := s.tankList.GetByIndex(i)
		if v != nil {
			s.objFactory.RecycleTank(v)
		}
	}
	for i := int32(0); i < s.shellList.Count(); i++ {
		_, v := s.shellList.GetByIndex(i)
		if v != nil {
			s.objFactory.RecycleShell(v)
		}
	}
	for i := int32(0); i < s.surroundObjList.Count(); i++ {
		_, v := s.surroundObjList.GetByIndex(i)
		if v != nil {
			s.objFactory.RecycleSurroundObj(v)
		}
	}
	for i := int32(0); i < s.effectList.Count(); i++ {
		_, v := s.effectList.GetByIndex(i)
		if v != nil {
			s.effectPool.Put(v)
		}
	}
	s.staticObjList.Clear()
	s.tankList.Clear()
	s.shellList.Clear()
	s.surroundObjList.Clear()
	s.objFactory.Clear()
	s.effectList.Clear()
	s.effectPool.Clear()
	s.surroundObjRecycleList = s.surroundObjRecycleList[:0]
	s.tankRecycleList = s.tankRecycleList[:0]
	s.shellRecycleList = s.shellRecycleList[:0]
	s.staticObjRecycleList = s.staticObjRecycleList[:0]
	s.effectRecycleList = s.effectRecycleList[:0]
	s.effectSearchedList = s.effectSearchedList[:0]
}

func (s *World) ReloadMap() {
	s.gmap.ClearObjsData()
	s.clearObjsData()
	s.loadMap(s.mapConfig, true)
}

func (s *World) GetGridMap() *GridMap {
	return s.gmap
}

func (s *World) GetMapLeftBottom() (int32, int32) {
	return s.mapConfig.X, s.mapConfig.Y
}

func (s *World) GetMapWidthHeight() (int32, int32) {
	return s.mapWidth, s.mapHeight
}

func (s *World) Center() (int32, int32) {
	return s.mapConfig.X + s.mapWidth/2, s.mapConfig.Y + s.mapHeight/2
}

func (s *World) GetTankBornPosList() []TankBornPosInfo {
	return s.tankBornPosList
}

func (s *World) RegisterEventHandle(id base.EventId, handle func(...any)) {
	s.eventMgr.RegisterEvent(id, handle)
}

func (s *World) UnregisterEventHandle(id base.EventId, handle func(...any)) {
	s.eventMgr.UnregisterEvent(id, handle)
}

func (s *World) RegisterEffectAddedHandle(handle func(...any)) {
	s.effectAddedEvent.Register(handle)
}

func (s *World) UnregisterEffectAddedHandle(handle func(...any)) {
	s.effectAddedEvent.Unregister(handle)
}

func (s *World) RegisterEffectRemovedHandle(handle func(...any)) {
	s.effectRemovedEvent.Register(handle)
}

func (s *World) UnregisterEffectRemovedHandle(handle func(...any)) {
	s.effectRemovedEvent.Unregister(handle)
}

func (s *World) RegisterObjectAddedHandle(handle func(...any)) {
	s.objAddedEvent.Register(handle)
}

func (s *World) UnregisterObjectAddedHandle(handle func(...any)) {
	s.objAddedEvent.Unregister(handle)
}

func (s *World) RegisterObjectRemovedHandle(handle func(...any)) {
	s.objRemovedEvent.Register(handle)
}

func (s *World) UnregisterObjectRemovedHandle(handle func(...any)) {
	s.objRemovedEvent.Unregister(handle)
}

// 注册坦克事件
func (s *World) RegisterTankEvent(instId uint32, eid base.EventId, handle func(args ...any)) {
	tank, o := s.tankList.Get(instId)
	if !o {
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

// 注銷坦克事件
func (s *World) UnregisterTankEvent(instId uint32, eid base.EventId, handle func(args ...any)) {
	tank, o := s.tankList.Get(instId)
	if !o {
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

func (s *World) GetMapConfig() *game_map.Config {
	return s.mapConfig
}

func (s *World) GetLayerObjsWithRange(rect *math.Rect) [MapMaxLayer]*heap.BinaryHeapKV[uint32, int32] {
	return s.gmap.GetLayerObjsWithRange(rect)
}

func (s *World) GetObj(instId uint32) object.IObject {
	return s.objFactory.GetObj(instId)
}

func (s *World) GetTankListWithRange(rect *math.Rect) []uint32 {
	return s.gmap.GetMovableObjListWithRangeAndSubtype(rect, base.ObjSubtypeTank)
}

func (s *World) GetEffectListWithRange(rect *math.Rect) []uint32 {
	s.effectSearchedList = s.effectSearchedList[:0]
	count := s.effectList.Count()
	for i := int32(0); i < count; i++ {
		instId, effect := s.effectList.GetByIndex(i)
		ex, ey := effect.Pos()
		el := ex - effect.Width()/2
		er := ex + effect.Width()/2
		et := ey + effect.Height()/2
		eb := ey - effect.Height()/2
		if !(er <= rect.X() || el >= rect.X()+rect.W() || et <= rect.Y() || eb >= rect.Y()+rect.H()) {
			s.effectSearchedList = append(s.effectSearchedList, instId)
		}
	}
	return s.effectSearchedList
}

func (s *World) GetEffect(instId uint32) effect.IEffect {
	effect, o := s.effectList.Get(instId)
	if !o {
		return nil
	}
	return effect
}

func (s *World) GetTank(instId uint32) *object.Tank {
	tank, o := s.tankList.Get(instId)
	if !o {
		return nil
	}
	return tank
}

func (s *World) NewTankWithPos(x, y int32) *object.Tank {
	tank := s.objFactory.NewTank(&s.mapConfig.PlayerTankInitData)
	// 注冊檢測移動事件處理
	tank.RegisterCheckMoveEventHandle(s.checkObjMoveEventHandle)
	tank.SetPos(x, y)
	// 設置碰撞處理
	tank.SetCollisionHandle(s.onTankCollision)
	// 加入到玩家坦克列表
	s.tankList.Add(tank.InstId(), tank)
	// 加入網格分區地圖
	s.gmap.AddObj(tank)
	// 物體創建事件
	s.objAddedEvent.Call(tank)
	return tank
}

func (s *World) AddTank(tank *object.Tank) {
	tank.RegisterCheckMoveEventHandle(s.checkObjMoveEventHandle)
	s.tankList.Add(tank.InstId(), tank)
	s.gmap.AddObj(tank)
}

func (s *World) NewTankWithStaticInfo(id int32, level int32, x, y int32 /*, currSpeed int32*/) *object.Tank {
	tank := s.objFactory.NewTank(common_data.TankConfigData[id])
	tank.SetPos(x, y)
	tank.SetLevel(level)
	//tank.SetCurrentSpeed(currSpeed)
	// 注冊檢測移動事件處理
	tank.RegisterCheckMoveEventHandle(s.checkObjMoveEventHandle)
	// 設置碰撞處理
	tank.SetCollisionHandle(s.onTankCollision)
	// 加入玩家坦克列表
	s.tankList.Add(tank.InstId(), tank)
	// 加入網格分區地圖
	s.gmap.AddObj(tank)
	s.objAddedEvent.Call(tank)
	return tank
}

func (s *World) RemoveTank(instId uint32) {
	tank := s.tankList.Remove(instId)
	tank.UnregisterCheckMoveEventHandle(s.checkObjMoveEventHandle)
	s.gmap.RemoveObj(tank.InstId())
	s.objRemovedEvent.Call(tank)
	s.objFactory.RecycleTank(tank)
}

func (s *World) TankMove(instId uint32, orientation int32) {
	tank, o := s.tankList.Get(instId)
	if !o {
		log.Error("tank %v not found", instId)
		return
	}
	angle := base.NewAngle(int16(orientation), 0)
	tank.Move(angle)
}

func (s *World) TankStopMove(instId uint32) {
	tank, o := s.tankList.Get(instId)
	if !o {
		log.Error("tank %v not found", instId)
		return
	}
	tank.Stop()
}

func (s *World) TankFire(instId uint32) {
	tank, o := s.tankList.Get(instId)
	if !o {
		log.Error("tank %v not found", instId)
		return
	}
	shell := tank.CheckAndFire(s.objFactory.NewShell)
	if shell != nil {
		s.shellList.Add(shell.InstId(), shell)
		shell.RegisterCheckMoveEventHandle(s.checkObjMoveEventHandle)
		shell.SetCollisionHandle(s.onShellCollision)
		if shell.ShellStaticInfo().TrackTarget {
			shell.SetSearchTargetFunc(s.searchShellTarget)
			shell.SetFetchTargetFunc(s.GetObj)
		}
		s.gmap.AddObj(shell)
		shell.RotateTo(tank.Rotation())
		if tank.IsMoving() {
			shell.MoveNow(tank.Rotation())
		} else {
			shell.Move(tank.Rotation())
		}
	}
}

func (w *World) TankEmitLaser(tankInstId uint32) {
	tank, o := w.tankList.Get(tankInstId)
	if !o {
		log.Error("World: tank %v not found", tankInstId)
		return
	}
	laserId := tank.TankStaticInfo().LaserId
	laserStaticInfo := common_data.LaserConfigData[laserId]
	laser := tank.LaunchLaser(laserStaticInfo)
	laser.SetEffectFunc(w.LaserEffect)
	laser.Emit()
}

func (w *World) TankCancelLaser(tankInstId uint32) {
	tank, o := w.tankList.Get(tankInstId)
	if !o {
		log.Error("World: tank %v not found", tankInstId)
		return
	}
	tank.CancelLaser()
}

func (s *World) TankAddNewShell(instId uint32, shellConfigId int32) bool {
	tank, o := s.tankList.Get(instId)
	if !o {
		log.Error("World: tank %v not found", instId)
		return false
	}
	return tank.AppendShell(common_data.ShellConfigData[shellConfigId])
}

func (s *World) TankSwitchShell(instId uint32) {
	tank, o := s.tankList.Get(instId)
	if !o {
		log.Error("World: tank %v not found while switch shell", instId)
		return
	}
	tank.SwitchShell()
}

func (s *World) TankReleaseSurroundObj(instId uint32) {
	tank, o := s.tankList.Get(instId)
	if !o {
		log.Error("tank %v not found", instId)
		return
	}
	ball := s.objFactory.NewSurroundObj(common_data.SurroundObjConfigData[1])
	ball.SetAroundCenterObject(tank.InstId(), s.objFactory.GetObj)
	s.surroundObjList.Add(ball.InstId(), ball)
	s.objAddedEvent.Call(ball)
	s.gmap.AddObj(ball)
	ball.Move(base.NewAngle(0, 0))
}

func (s *World) TankRotate(instId uint32, degree int32) {
	tank, o := s.tankList.Get(instId)
	if !o {
		log.Error("tank %v not found", instId)
		return
	}
	angle := base.Angle{}
	angle.Reset(int16(degree), 0)
	tank.Rotate(angle)
}

func (s *World) TankChange(instId uint32, staticInfo *object.TankStaticInfo) bool {
	tank, o := s.tankList.Get(instId)
	if !o {
		return false
	}
	tank.Change(staticInfo)
	return true
}

func (s *World) TankRestore(instId uint32) int32 {
	tank, o := s.tankList.Get(instId)
	if !o {
		return 0
	}
	tank.Restore()
	return tank.Id()
}

func (s *World) TankShield(instId uint32, shieldId int32) {
	tank, o := s.tankList.Get(instId)
	if !o {
		return
	}
	if tank.HasShield() {
		tank.CancelShield()
	} else {
		tank.AddShield(common_data.TankShieldConfigData[shieldId])
	}
}

func (s *World) TankUnlimitedShield(instId uint32) {
	s.TankShield(instId, 1)
}

func (s *World) Update(tickMs uint32) {
	count := s.tankList.Count()
	for i := int32(0); i < count; i++ {
		_, tank := s.tankList.GetByIndex(i)
		tank.Update(tickMs)        // 相當於MonoBehevior.Update
		s.gmap.UpdateMovable(tank) // 相當於MonoBehevior.FixedUpdate
		if tank.IsRecycle() {
			s.tankRecycleList = append(s.tankRecycleList, tank)
		}
	}

	count = s.shellList.Count()
	for i := int32(0); i < count; i++ {
		_, shell := s.shellList.GetByIndex(i)
		shell.Update(tickMs)
		s.gmap.UpdateMovable(shell)
		if shell.IsRecycle() {
			s.shellRecycleList = append(s.shellRecycleList, shell)
		}
	}

	// todo 環繞物體一定要在坦克後面更新
	count = s.surroundObjList.Count()
	for i := int32(0); i < count; i++ {
		_, ball := s.surroundObjList.GetByIndex(i)
		ball.Update(tickMs)
		s.gmap.UpdateMovable(ball)
		if ball.IsRecycle() {
			s.surroundObjRecycleList = append(s.surroundObjRecycleList, ball)
		}
	}

	count = s.effectList.Count()
	for i := int32(0); i < count; i++ {
		_, effect := s.effectList.GetByIndex(i)
		effect.Update(tickMs) // todo 效果暫時不可移動，所以不需要在parition_map中更新位置
		if effect.IsOver() {
			s.effectRecycleList = append(s.effectRecycleList, effect)
		}
	}

	if len(s.staticObjRecycleList) > 0 {
		for _, obj := range s.staticObjRecycleList {
			s.staticObjList.Remove(obj.InstId())
			s.gmap.RemoveObj(obj.InstId())
			s.objRemovedEvent.Call(obj)
			s.objFactory.RecycleStaticObject(obj)
		}
		s.staticObjRecycleList = s.staticObjRecycleList[:0]
	}

	if len(s.tankRecycleList) > 0 {
		for _, tank := range s.tankRecycleList {
			s.tankList.Remove(tank.InstId())
			s.gmap.RemoveObj(tank.InstId())
			s.objRemovedEvent.Call(tank)
			s.objFactory.RecycleTank(tank)
		}
		s.tankRecycleList = s.tankRecycleList[:0]
	}

	if len(s.shellRecycleList) > 0 {
		for _, shell := range s.shellRecycleList {
			s.shellList.Remove(shell.InstId())
			s.gmap.RemoveObj(shell.InstId())
			s.objRemovedEvent.Call(shell)
			s.objFactory.RecycleShell(shell)
		}
		s.shellRecycleList = s.shellRecycleList[:0]
	}

	if len(s.surroundObjRecycleList) > 0 {
		for _, ball := range s.surroundObjRecycleList {
			s.surroundObjList.Remove(ball.InstId())
			s.gmap.RemoveObj(ball.InstId())
			s.objRemovedEvent.Call(ball)
			s.objFactory.RecycleSurroundObj(ball)
		}
		s.surroundObjRecycleList = s.surroundObjRecycleList[:0]
	}

	if len(s.effectRecycleList) > 0 {
		for _, effect := range s.effectRecycleList {
			s.effectList.Remove(effect.InstId())
			s.effectRemovedEvent.Call(effect)
			s.effectPool.Put(effect)
		}
		s.effectRecycleList = s.effectRecycleList[:0]
	}
}

func (s *World) Pause() {
	for i := int32(0); i < s.tankList.Count(); i++ {
		_, tank := s.tankList.GetByIndex(i)
		tank.Pause()
	}
	for i := int32(0); i < s.shellList.Count(); i++ {
		_, shell := s.shellList.GetByIndex(i)
		shell.Pause()
	}
	for i := int32(0); i < s.surroundObjList.Count(); i++ {
		_, sobj := s.surroundObjList.GetByIndex(i)
		sobj.Pause()
	}
}

func (s *World) Resume() {
	for i := int32(0); i < s.tankList.Count(); i++ {
		_, tank := s.tankList.GetByIndex(i)
		tank.Resume()
	}
	for i := int32(0); i < s.shellList.Count(); i++ {
		_, shell := s.shellList.GetByIndex(i)
		shell.Resume()
	}
	for i := int32(0); i < s.surroundObjList.Count(); i++ {
		_, sobj := s.surroundObjList.GetByIndex(i)
		sobj.Resume()
	}
}

func (s *World) checkObjMoveEventHandle(args ...any) {
	instId := args[0].(uint32)
	dx := args[1].(int32)
	dy := args[2].(int32)
	ci := args[3].(*object.CollisionInfo)

	obj := s.objFactory.GetObj(instId)
	if obj.Type() != base.ObjTypeMovable {
		log.Error("World.checkObjMoveEventHandle object %v must be movable", instId)
		return
	}

	var (
		mobj = obj.(object.IMovableObject)
	)
	if !s.checkObjMoveRange(mobj, dx, dy, ci) {
		s.onMovableObjReachMapBorder(mobj)
	} else {
		s.gmap.CheckMovingObjCollision(mobj, dx, dy, ci)
	}
}

func (s *World) checkObjMoveRange(obj object.IMovableObject, dx, dy int32, ci *object.CollisionInfo) bool {
	x, y := obj.Pos()
	colliderComp := obj.GetColliderComp()
	if dx != 0 {
		if colliderComp != nil {
			aabb := colliderComp.GetAABB()
			if aabb.Left+dx <= s.mapConfig.X {
				ci.Result = object.CollisionAndBlock
				x = s.mapConfig.X + (aabb.Right-aabb.Left)/2
			}
			if aabb.Right+dx >= s.mapConfig.X+s.mapWidth {
				ci.Result = object.CollisionAndBlock
				x = s.mapConfig.X + s.mapWidth - (aabb.Right-aabb.Left)/2
			}
		} else {
			if x-obj.Width()/2 <= s.mapConfig.X {
				ci.Result = object.CollisionAndBlock
				x = s.mapConfig.X + obj.Width()/2
			}
			if x+obj.Width()/2 >= s.mapConfig.X+s.mapWidth {
				ci.Result = object.CollisionAndBlock
				x = s.mapConfig.X + s.mapWidth - obj.Width()/2
			}
		}
	}
	if dy != 0 {
		if colliderComp != nil {
			aabb := colliderComp.GetAABB()
			if aabb.Top+dy >= s.mapConfig.Y+s.mapHeight {
				ci.Result = object.CollisionAndBlock
				y = s.mapConfig.Y + s.mapHeight - (aabb.Top-aabb.Bottom)/2
			}
			if aabb.Bottom+dy <= s.mapConfig.Y {
				ci.Result = object.CollisionAndBlock
				y = s.mapConfig.Y + (aabb.Top-aabb.Bottom)/2
			}
		} else {
			if y-obj.Length()/2 <= s.mapConfig.Y {
				ci.Result = object.CollisionAndBlock
				y = s.mapConfig.Y + obj.Length()/2
			}
			if y+obj.Length()/2 >= s.mapConfig.Y+s.mapHeight {
				ci.Result = object.CollisionAndBlock
				y = s.mapConfig.Y + s.mapHeight - obj.Length()/2
			}
		}
	}
	if ci.Result == object.CollisionAndBlock {
		ci.MovingObjPos.X = x
		ci.MovingObjPos.Y = y
	}
	return ci.Result != object.CollisionAndBlock
}

func (s *World) onMovableObjReachMapBorder(mobj object.IMovableObject) {
	if mobj.Subtype() == base.ObjSubtypeShell {
		mobj.ToRecycle()
	}
}

func (s *World) onTankCollision(mobj object.IMovableObject, ci *object.CollisionInfo) {
	tank, o := mobj.(*object.Tank)
	if !o {
		return
	}
	for i := 0; i < len(ci.ObjList); i++ {
		obj := ci.ObjList[i]
		objType := obj.Type()
		if objType == base.ObjTypeMovable {
			if obj.Subtype() == base.ObjSubtypeShell {
				s.shellEffect(obj.(*object.Shell), tank)
			}
		} else if objType == base.ObjTypeItem {
			objSubtype := obj.Subtype()
			switch objSubtype {
			case base.ObjSubtypeRewardLife:
			case base.ObjSubtypeReinforcement:
			case base.ObjSubtypeFrozen:
			case base.ObjSubtypeSelfUpgrade:
			case base.ObjSubtypeShield:
				tank.AddShield(common_data.TankShieldConfigData[2])
			case base.ObjSubtypeBomb:
			}
		}
	}
}

func (s *World) onShellCollision(mobj object.IMovableObject, ci *object.CollisionInfo) {
	shell := mobj.(*object.Shell)
	for i := 0; i < len(ci.ObjList); i++ {
		obj := ci.ObjList[i]
		s.shellEffect(shell, obj)
	}
}

func (s *World) shellEffect(shell *object.Shell, obj object.IObject) {
	var (
		objType      = obj.Type()
		effectParams = [2]struct {
			effectFunc func(...any)
			effectId   int32
			cx, cy     int32
		}{}
	)
	shell.ToRecycle()
	effectParams[0].effectId = 1
	effectParams[0].effectFunc = bulletExplodeEffect
	effectParams[0].cx, effectParams[0].cy = shell.Pos()
	if objType == base.ObjTypeMovable {
		if obj.Subtype() == base.ObjSubtypeShell {
			obj.ToRecycle()
			effectParams[1].effectId = 1
			effectParams[1].effectFunc = bulletExplodeEffect
			effectParams[1].cx, effectParams[1].cy = obj.Pos()
		} else if obj.Subtype() == base.ObjSubtypeTank {
			tank := obj.(*object.Tank)
			if !tank.HasShield() {
				obj.ToRecycle()
				effectParams[1].effectId = 2
				effectParams[1].effectFunc = bigBulletExplodeEffect
				effectParams[1].cx, effectParams[1].cy = obj.Pos()
			}
		}
	}
	// 生成爆炸效果
	for i := 0; i < len(effectParams); i++ {
		if effectParams[i].effectId <= 0 {
			continue
		}
		effect := s.effectPool.Get(common_data.EffectConfigData[effectParams[i].effectId], effectParams[i].effectFunc, s.gmap, shell)
		effect.SetPos(effectParams[i].cx, effectParams[i].cy)
		s.effectList.Add(effect.InstId(), effect)
	}
}

func (s *World) searchShellTarget(shell *object.Shell) object.IObject {
	staticInfo := shell.ShellStaticInfo()
	cx, cy := shell.Pos()
	rect := math.NewRect(cx-staticInfo.SearchTargetRadius, cy-staticInfo.SearchTargetRadius, cx+staticInfo.SearchTargetRadius, cy+staticInfo.SearchTargetRadius)
	objList := s.gmap.GetMovableObjListWithRangeAndSubtype(rect, base.ObjSubtypeTank)
	var (
		tank *object.Tank
		o    bool
		sd   int64 = -1
		n    int   = -1
	)
	for i := 0; i < len(objList); i++ {
		tank, o = s.tankList.Get(objList[i])
		if o && tank.Camp() != shell.Camp() {
			t := object.SquareOfDistance(tank, shell)
			if sd < 0 || sd > t {
				sd = t
				n = i
			}
		}
	}
	if n >= 0 {
		tank, o = s.tankList.Get(objList[n])
		if !o {
			return nil
		}
		return tank
	}
	return nil
}

// 激光效果
func (s *World) LaserEffect(laser *weapon.Laser, start, end base.Pos) (pos base.Pos, result bool) {
	log.Debug("LaserEffect: start %v, end %v", start, end)
	if !s.posIsValid(&start) {
		return
	}
	if start == end {
		return
	}
	if !s.posIsValid(&end) {
		if start.X == end.X {
			if end.Y < s.mapConfig.Y {
				end.Y = s.mapConfig.Y
			}
			if end.Y >= s.mapConfig.Y+s.mapHeight {
				end.Y = s.mapConfig.Y + s.mapHeight - 1
			}
		} else if start.Y == end.Y {
			if end.X < s.mapConfig.X {
				end.X = s.mapConfig.X
			}
			if end.X >= s.mapConfig.X+s.mapWidth {
				end.X = s.mapConfig.X + s.mapWidth - 1
			}
		} else {
			var x, y int32
			// (end.X-start.X)/(end.Y-start.Y) == (x-start.X)/(y-start.Y)
			// 判斷end點在start點的方位
			if end.X < start.X && end.Y < start.Y {
				// 左下方
				// 令x=s.mapConfig.X，求出y，如果y在(s.mapConfig.Y, s.mapConfig.Y+s.mapHeight)
				// 之間，無需再判斷，否則令y=s.mapConfig.Y，求出x
				x = s.mapConfig.X
				y = (end.Y-start.Y)*(x-start.X)/(end.X-start.X) + start.Y
				if y < s.mapConfig.Y {
					y = s.mapConfig.Y
					x = (end.X-start.X)*(y-start.Y)/(end.Y-start.Y) + start.X
				}
			} else if end.X > start.X && end.Y < start.Y {
				// 右下方
				x = s.mapConfig.X + s.mapWidth - 1
				y = (end.Y-start.Y)*(x-start.X)/(end.X-start.X) + start.Y
				if y < s.mapConfig.Y {
					y = s.mapConfig.Y
					x = (end.X-start.X)*(y-start.Y)/(end.Y-start.Y) + start.X
				}
			} else if end.X > start.X && end.Y > start.Y {
				// 右上方
				x = s.mapConfig.X + s.mapWidth - 1
				y = (end.Y-start.Y)*(x-start.X)/(end.X-start.X) + start.Y
				if y >= s.mapConfig.Y+s.mapHeight {
					y = s.mapConfig.Y + s.mapHeight - 1
					x = (end.X-start.X)*(y-start.Y)/(end.Y-start.Y) + start.X
				}
			} else {
				// 左上方
				x = s.mapConfig.X
				y = (end.Y-start.Y)*(s.mapConfig.X-start.X)/(end.X-start.X) + start.Y
				if y >= s.mapConfig.Y+s.mapHeight {
					y = s.mapConfig.Y + s.mapHeight - 1
					x = (end.X-start.X)*(y-start.Y)/(end.Y-start.Y) + start.X
				}
			}
			end.X, end.Y = x, y
		}
	}

	// 判斷起始點是否在某個物體内
	objList := s.gmap.PointInObjList(start)
	for _, id := range objList {
		obj := s.GetObj(id)
		if obj == nil {
			continue
		}
		typ := obj.Type()
		subtype := obj.Subtype()
		switch typ {
		case base.ObjTypeStatic:
			switch subtype {
			case base.ObjSubtypeBrick, base.ObjSubtypeIron:
				return
			}
		}
	}

	pos = end

	// 獲得綫段相交的物體列表
	var intersectInfo IntersectInfo
	if s.gmap.GetLineSegmentFirstIntersectInfo(&start, &end, &intersectInfo) {
		obj := intersectInfo.obj
		typ := obj.Type()
		subtype := obj.Subtype()
		switch typ {
		case base.ObjTypeStatic:
			switch subtype {
			case base.ObjSubtypeBrick, base.ObjSubtypeIron:
				pos = intersectInfo.pos
			}
		case base.ObjTypeMovable:
			switch subtype {
			case base.ObjSubtypeTank:
				if obj.Camp() != laser.Camp() {
					obj.ToRecycle()
				}
			}
		}
	} else {
		log.Debug("not found intersect point, pos is %v", pos)
	}
	return pos, true
}

func (w *World) posIsValid(pos *base.Pos) bool {
	if pos.X < w.mapConfig.X || pos.X >= w.mapConfig.X+w.mapWidth {
		return false
	}
	if pos.Y < w.mapConfig.Y || pos.Y >= w.mapConfig.Y+w.mapHeight {
		return false
	}
	return true
}
