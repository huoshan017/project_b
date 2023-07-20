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

	"github.com/huoshan017/ponu/heap"
)

type PlayerTankKV struct {
	PlayerId uint64
	Tank     *object.Tank
}

// 场景圖，沒有玩家(Player)概念的游戲邏輯
// 必须在单个goroutine中执行
type SceneLogic struct {
	mapConfig                                      *game_map.Config                               // 地圖配置
	mapWidth, mapHeight                            int32                                          // 地圖寬高
	gmap                                           *GridMap                                       // 網格地圖
	eventMgr                                       base.IEventManager                             // 事件管理器
	staticObjList                                  *ds.MapListUnion[uint32, *object.StaticObject] // 靜態對象列表，用map和list的聯合體是爲了遍歷時的有序性
	tankList                                       *ds.MapListUnion[uint32, *object.Tank]         // 坦克列表，不區分玩家和BOT
	shellList                                      *ds.MapListUnion[uint32, *object.Shell]        // 炮彈列表，不區分坦克的炮彈
	surroundObjList                                *ds.MapListUnion[uint32, *object.SurroundObj]  // 環繞物體列表
	objFactory                                     *object.ObjectFactory                          // 對象池
	effectList                                     *ds.MapListUnion[uint32, *object.Effect]       // 效果列表
	effectPool                                     *object.EffectPool                             // 效果池
	staticObjAddedEvent, staticObjRemovedEvent     base.Event                                     // 靜態對象添加刪除事件
	tankAddedEvent, tankRemovedEvent               base.Event                                     // 坦克添加刪除事件
	shellAddedEvent, shellRemovedEvent             base.Event                                     // 子彈添加刪除事件
	surroundObjAddedEvent, surroundObjRemovedEvent base.Event                                     // 環繞物體添加刪除事件
	effectAddedEvent, effectRemovedEvent           base.Event                                     // 效果添加刪除事件
	staticObjRecycleList                           []*object.StaticObject                         // 靜態對象回收列表
	tankRecycleList                                []*object.Tank                                 // 坦克對象回收列表
	shellRecycleList                               []*object.Shell                                // 炮彈對象回收列表
	surroundObjRecycleList                         []*object.SurroundObj                          // 環繞物體對象回收列表
	effectRecycleList                              []*object.Effect                               // 效果回收列表
	effectSearchedList                             []uint32                                       // 效果搜索結果列表
}

func NewSceneLogic(eventMgr base.IEventManager) *SceneLogic {
	return &SceneLogic{
		eventMgr:        eventMgr,
		staticObjList:   ds.NewMapListUnion[uint32, *object.StaticObject](),
		tankList:        ds.NewMapListUnion[uint32, *object.Tank](),
		shellList:       ds.NewMapListUnion[uint32, *object.Shell](),
		surroundObjList: ds.NewMapListUnion[uint32, *object.SurroundObj](),
		objFactory:      object.NewObjectFactory(true),
		effectList:      ds.NewMapListUnion[uint32, *object.Effect](),
		effectPool:      object.NewEffectPool(),
		gmap:            NewGridMap(1),
	}
}

func (s *SceneLogic) GetMapId() int32 {
	return s.mapConfig.Id
}

func (s *SceneLogic) LoadMap(m *game_map.Config) bool {
	// 载入地图
	s.gmap.Load(m)
	for line := 0; line < len(m.Layers); line++ {
		for col := 0; col < len(m.Layers[line]); col++ {
			st := object.StaticObjType(m.Layers[line][col])
			if common_data.StaticObjectConfigData[st] == nil {
				continue
			}
			tileObj := s.objFactory.NewStaticObject(common_data.StaticObjectConfigData[st])
			// 二維數組Y軸是自上而下的，而世界坐標Y軸是自下而上的，所以設置Y坐標要倒過來
			tileObj.SetPos(m.TileWidth*int32(col), m.TileHeight*int32(len(m.Layers)-1-line))
			s.staticObjAddedEvent.Call(tileObj)
			// 加入網格分區地圖
			s.gmap.AddTile(int16(len(m.Layers)-1-line), int16(col), tileObj)
		}
	}
	s.mapConfig = m
	s.mapWidth = int32(len(m.Layers[0])) * m.TileWidth
	s.mapHeight = int32(len(m.Layers)) * m.TileHeight
	log.Info("Load map %v done, map width %v, map height %v", m.Id, s.mapWidth, s.mapHeight)
	return true
}

func (s *SceneLogic) UnloadMap() {
	s.mapWidth = 0
	s.mapHeight = 0
	for i := int32(0); i < s.staticObjList.Count(); i++ {
		_, v := s.staticObjList.GetByIndex(i)
		if v != nil {
			s.staticObjRemovedEvent.Call(v)
			s.objFactory.RecycleStaticObject(v)
		}
	}
	for i := int32(0); i < s.tankList.Count(); i++ {
		_, v := s.tankList.GetByIndex(i)
		if v != nil {
			s.tankRemovedEvent.Call(v)
			s.objFactory.RecycleTank(v)
		}
	}
	for i := int32(0); i < s.shellList.Count(); i++ {
		_, v := s.shellList.GetByIndex(i)
		if v != nil {
			s.shellRemovedEvent.Call(v)
			s.objFactory.RecycleShell(v)
		}
	}
	for i := int32(0); i < s.surroundObjList.Count(); i++ {
		_, v := s.surroundObjList.GetByIndex(i)
		if v != nil {
			s.surroundObjRemovedEvent.Call(v)
			s.objFactory.RecycleSurroundObj(v)
		}
	}
	for i := int32(0); i < s.effectList.Count(); i++ {
		_, v := s.effectList.GetByIndex(i)
		if v != nil {
			s.effectRemovedEvent.Call(v)
			s.effectPool.Put(v)
		}
	}
	s.staticObjList.Clear()
	s.tankList.Clear()
	s.shellList.Clear()
	s.surroundObjList.Clear()
	s.gmap.Unload()
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

func (s *SceneLogic) RegisterStaticObjAddedHandle(handle func(...any)) {
	s.staticObjAddedEvent.Register(handle)
}

func (s *SceneLogic) UnregisterStaticObjAddedHandle(handle func(...any)) {
	s.staticObjAddedEvent.Unregister(handle)
}

func (s *SceneLogic) RegisterStaticObjRemovedHandle(handle func(...any)) {
	s.staticObjRemovedEvent.Register(handle)
}

func (s *SceneLogic) UnregisterStaticObjRemovedHandle(handle func(...any)) {
	s.staticObjRemovedEvent.Unregister(handle)
}

func (s *SceneLogic) RegisterTankAddedHandle(handle func(...any)) {
	s.tankAddedEvent.Register(handle)
}

func (s *SceneLogic) UnregisterTankAddedHandle(handle func(...any)) {
	s.tankAddedEvent.Unregister(handle)
}

func (s *SceneLogic) RegisterTankRemovedHandle(handle func(...any)) {
	s.tankRemovedEvent.Register(handle)
}

func (s *SceneLogic) UnregisterTankRemovedHandle(handle func(...any)) {
	s.tankRemovedEvent.Unregister(handle)
}

func (s *SceneLogic) RegisterShellAddedHandle(handle func(...any)) {
	s.shellAddedEvent.Register(handle)
}

func (s *SceneLogic) UnregisterShellAddedHandle(handle func(...any)) {
	s.shellAddedEvent.Unregister(handle)
}

func (s *SceneLogic) RegisterShellRemovedHandle(handle func(...any)) {
	s.shellRemovedEvent.Register(handle)
}

func (s *SceneLogic) UnregisterShellRemovedHandle(handle func(...any)) {
	s.shellRemovedEvent.Unregister(handle)
}

func (s *SceneLogic) RegisterSurroundObjAddedHandle(handle func(...any)) {
	s.surroundObjAddedEvent.Register(handle)
}

func (s *SceneLogic) UnregisterSurroundObjAddedHandle(handle func(...any)) {
	s.surroundObjAddedEvent.Unregister(handle)
}

func (s *SceneLogic) RegisterSurroundObjRemovedHandle(handle func(...any)) {
	s.surroundObjRemovedEvent.Register(handle)
}

func (s *SceneLogic) UnregisterSurroundObjRemovedHandle(handle func(...any)) {
	s.surroundObjRemovedEvent.Unregister(handle)
}

func (s *SceneLogic) RegisterEffectAddedHandle(handle func(...any)) {
	s.effectAddedEvent.Register(handle)
}

func (s *SceneLogic) UnregisterEffectAddedHandle(handle func(...any)) {
	s.effectAddedEvent.Unregister(handle)
}

func (s *SceneLogic) RegisterEffectRemovedHandle(handle func(...any)) {
	s.effectRemovedEvent.Register(handle)
}

func (s *SceneLogic) UnregisterEffectRemovedHandle(handle func(...any)) {
	s.effectRemovedEvent.Unregister(handle)
}

func (s *SceneLogic) GetMapConfig() *game_map.Config {
	return s.mapConfig
}

func (s *SceneLogic) GetLayerObjsWithRange(rect *math.Rect) [MapMaxLayer]*heap.BinaryHeapKV[uint32, int32] {
	return s.gmap.GetLayerObjsWithRange(rect)
}

func (s *SceneLogic) GetObj(instId uint32) object.IObject {
	return s.objFactory.GetObj(instId)
}

func (s *SceneLogic) GetTankListWithRange(rect *math.Rect) []uint32 {
	return s.gmap.GetMovableObjListWithRangeAndSubtype(rect, object.ObjSubtypeTank)
}

func (s *SceneLogic) GetEffectListWithRange(rect *math.Rect) []uint32 {
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

func (s *SceneLogic) GetEffect(instId uint32) object.IEffect {
	effect, o := s.effectList.Get(instId)
	if !o {
		return nil
	}
	return effect
}

func (s *SceneLogic) GetTank(instId uint32) *object.Tank {
	tank, o := s.tankList.Get(instId)
	if !o {
		return nil
	}
	return tank
}

func (s *SceneLogic) NewTankWithPos(x, y int32) *object.Tank {
	tank := s.objFactory.NewTank(&s.mapConfig.PlayerTankInitData)
	// 注冊檢測移動事件處理
	tank.RegisterCheckMoveEventHandle(s.checkObjMoveEventHandle)
	tank.SetPos(x, y)
	// 加入到玩家坦克列表
	s.tankList.Add(tank.InstId(), tank)
	// 加入網格分區地圖
	//s.pmap.AddObj(tank)
	s.gmap.AddObj(tank)
	// 物體創建事件
	s.tankAddedEvent.Call(tank)
	return tank
}

func (s *SceneLogic) AddTank(tank *object.Tank) {
	tank.RegisterCheckMoveEventHandle(s.checkObjMoveEventHandle)
	s.tankList.Add(tank.InstId(), tank)
	s.gmap.AddObj(tank)
}

func (s *SceneLogic) NewTankWithStaticInfo(id int32, level int32, x, y int32 /*dir object.Direction, */, currSpeed int32) *object.Tank {
	tank := s.objFactory.NewTank(common_data.TankConfigData[id])
	tank.SetPos(x, y)
	tank.SetLevel(level)
	tank.SetCurrentSpeed(currSpeed)
	// 注冊檢測移動事件處理
	tank.RegisterCheckMoveEventHandle(s.checkObjMoveEventHandle)
	// 加入玩家坦克列表
	s.tankList.Add(tank.InstId(), tank)
	// 加入網格分區地圖
	//s.pmap.AddObj(tank)
	s.gmap.AddObj(tank)
	s.tankAddedEvent.Call(tank)
	return tank
}

func (s *SceneLogic) RemoveTank(instId uint32) {
	tank := s.tankList.Remove(instId)
	tank.UnregisterCheckMoveEventHandle(s.checkObjMoveEventHandle)
	s.gmap.RemoveObj(tank.InstId())
	s.tankRemovedEvent.Call(tank)
	s.objFactory.RecycleTank(tank)
}

func (s *SceneLogic) TankMove(instId uint32, orientation int32) {
	tank, o := s.tankList.Get(instId)
	if !o {
		log.Error("tank %v not found", instId)
		return
	}
	angle := base.NewAngleObj(int16(orientation), 0)
	tank.Move( /*dir*/ angle)
}

func (s *SceneLogic) TankStopMove(instId uint32) {
	tank, o := s.tankList.Get(instId)
	if !o {
		log.Error("tank %v not found", instId)
		return
	}
	tank.Stop()
}

func (s *SceneLogic) TankFire(instId uint32, shellId int32) {
	tank, o := s.tankList.Get(instId)
	if !o {
		log.Error("tank %v not found", instId)
		return
	}
	shellConfig := common_data.ShellConfigData[shellId]
	shell := tank.CheckAndFire(s.objFactory.NewShell, shellConfig)
	if shell != nil {
		s.shellList.Add(shell.InstId(), shell)
		shell.RegisterCheckMoveEventHandle(s.checkObjMoveEventHandle)
		collider := shell.GetComp("Collider")
		if collider != nil {
			c := collider.(*object.ColliderComp)
			c.SetCollisionHandle(s.onShellCollision)
		}
		if shellConfig.TrackTarget {
			//shell.SetSearchTargetFunc(s.searchShellTarget)
			shell.SetSearchTargetFunc(s.testSearchShellTarget)
			shell.SetFetchTargetFunc(s.GetObj)
		}
		shell.Move(tank.MoveDir())
		s.gmap.AddObj(shell)
	}
}

func (s *SceneLogic) TankReleaseSurroundObj(instId uint32) {
	tank, o := s.tankList.Get(instId)
	if !o {
		log.Error("tank %v not found", instId)
		return
	}
	ball := s.objFactory.NewSurroundObj(common_data.SurroundObjConfigData[1])
	ball.SetAroundCenterObject(tank.InstId(), s.objFactory.GetObj)
	s.surroundObjList.Add(ball.InstId(), ball)
	s.surroundObjAddedEvent.Call(ball)
	s.gmap.AddObj(ball)
	ball.Move(base.NewAngleObj(0, 0))
}

func (s *SceneLogic) TankRotate(instId uint32, degree int32) {
	tank, o := s.tankList.Get(instId)
	if !o {
		log.Error("tank %v not found", instId)
		return
	}
	angle := base.Angle{}
	angle.Reset(int16(degree), 0)
	tank.Rotate(angle)
}

func (s *SceneLogic) TankChange(instId uint32, staticInfo *object.TankStaticInfo) bool {
	tank, o := s.tankList.Get(instId)
	if !o {
		return false
	}
	tank.Change(staticInfo)
	return true
}

func (s *SceneLogic) TankRestore(instId uint32) int32 {
	tank, o := s.tankList.Get(instId)
	if !o {
		return 0
	}
	tank.Restore()
	return tank.Id()
}

func (s *SceneLogic) Update(tick time.Duration) {
	count := s.tankList.Count()
	for i := int32(0); i < count; i++ {
		_, tank := s.tankList.GetByIndex(i)
		tank.Update(tick)          // 相當於MonoBehevior.Update
		s.gmap.UpdateMovable(tank) // 相當於MonoBehevior.FixedUpdate
		if tank.IsRecycle() {
			s.tankRecycleList = append(s.tankRecycleList, tank)
		}
	}

	count = s.shellList.Count()
	for i := int32(0); i < count; i++ {
		_, shell := s.shellList.GetByIndex(i)
		shell.Update(tick)
		s.gmap.UpdateMovable(shell)
		if shell.IsRecycle() {
			s.shellRecycleList = append(s.shellRecycleList, shell)
		}
	}

	// todo 環繞物體一定要在坦克後面更新
	count = s.surroundObjList.Count()
	for i := int32(0); i < count; i++ {
		_, ball := s.surroundObjList.GetByIndex(i)
		ball.Update(tick)
		s.gmap.UpdateMovable(ball)
		if ball.IsRecycle() {
			s.surroundObjRecycleList = append(s.surroundObjRecycleList, ball)
		}
	}

	count = s.effectList.Count()
	for i := int32(0); i < count; i++ {
		_, effect := s.effectList.GetByIndex(i)
		effect.Update() // todo 效果暫時不可移動，所以不需要在parition_map中更新位置
		if effect.IsOver() {
			s.effectRecycleList = append(s.effectRecycleList, effect)
		}
	}

	if len(s.staticObjRecycleList) > 0 {
		for _, obj := range s.staticObjRecycleList {
			s.staticObjList.Remove(obj.InstId())
			s.gmap.RemoveObj(obj.InstId())
			s.staticObjRemovedEvent.Call(obj)
			s.objFactory.RecycleStaticObject(obj)
		}
		s.staticObjRecycleList = s.staticObjRecycleList[:0]
	}

	if len(s.tankRecycleList) > 0 {
		for _, tank := range s.tankRecycleList {
			s.tankList.Remove(tank.InstId())
			s.gmap.RemoveObj(tank.InstId())
			s.tankRemovedEvent.Call(tank)
			s.objFactory.RecycleTank(tank)
		}
		s.tankRecycleList = s.tankRecycleList[:0]
	}

	if len(s.shellRecycleList) > 0 {
		for _, shell := range s.shellRecycleList {
			s.shellList.Remove(shell.InstId())
			s.gmap.RemoveObj(shell.InstId())
			s.shellRemovedEvent.Call(shell)
			s.objFactory.RecycleShell(shell)
		}
		s.shellRecycleList = s.shellRecycleList[:0]
	}

	if len(s.surroundObjRecycleList) > 0 {
		for _, ball := range s.surroundObjRecycleList {
			s.surroundObjList.Remove(ball.InstId())
			s.gmap.RemoveObj(ball.InstId())
			s.surroundObjRemovedEvent.Call(ball)
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

// 注册事件
func (s *SceneLogic) RegisterEvent(eid base.EventId, handle func(args ...any)) {
	s.eventMgr.RegisterEvent(eid, handle)
}

// 注销事件
func (s *SceneLogic) UnregisterEvent(eid base.EventId, handle func(args ...any)) {
	s.eventMgr.UnregisterEvent(eid, handle)
}

// 注册坦克事件
func (s *SceneLogic) RegisterTankEvent(instId uint32, eid base.EventId, handle func(args ...any)) {
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
func (s *SceneLogic) UnregisterTankEvent(instId uint32, eid base.EventId, handle func(args ...any)) {
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

func (s *SceneLogic) checkObjMoveEventHandle(args ...any) {
	instId := args[0].(uint32)
	dx := args[1].(int32)
	dy := args[2].(int32)
	isMove := args[3].(*bool)
	isCollision := args[4].(*bool)
	resObj := args[5].(*object.IObject)

	obj := s.objFactory.GetObj(instId)
	if obj.Type() != object.ObjTypeMovable {
		log.Error("SceneLogic.checkObjMoveEventHandle object %v must be movable", instId)
		return
	}

	var (
		x, y int32
		mobj = obj.(object.IMovableObject)
	)
	if !s.checkObjMoveRange(mobj /*dir,*/, dx, dy, &x, &y) {
		mobj.SetPos(x, y)
		mobj.Stop()
		*isMove = false
		*isCollision = false
		s.onMovableObjReachMapBorder(mobj)
	} else if s.gmap.CheckMovableObjCollision(mobj /*dir,*/, dx, dy, resObj) {
		mobj.Stop()
		*isCollision = true
		*isMove = false
	} else {
		*isMove = true
		*isCollision = false
	}
}

func (s *SceneLogic) checkObjMoveRange(obj object.IMovableObject /*dir object.Direction, */, dx, dy int32, rx, ry *int32) bool {
	x, y := obj.Pos()
	var move bool = true
	if x+dx <= s.mapConfig.X {
		move = false
		x = s.mapConfig.X
	}
	if x+dx >= s.mapConfig.X+s.mapWidth-obj.Width() {
		move = false
		x = s.mapConfig.X + s.mapWidth - obj.Width()
	}
	if y+dy >= s.mapConfig.Y+s.mapHeight-obj.Length() {
		move = false
		y = s.mapConfig.Y + s.mapHeight - obj.Length()
	}
	if y+dy <= s.mapConfig.Y {
		move = false
		y = s.mapConfig.Y
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

func (s *SceneLogic) onMovableObjReachMapBorder(mobj object.IMovableObject) {
	if mobj.Subtype() == object.ObjSubtypeShell {
		mobj.ToRecycle()
	}
}

func (s *SceneLogic) onShellCollision(args ...any) {
	// todo 處理炮彈碰撞
	bullet := args[0].(*object.Shell)
	obj := args[1].(object.IObject)
	var (
		effectParams = [2]struct {
			effectFunc func(...any)
			effectId   int32
			cx, cy     int32
		}{}
	)
	if obj.Type() == object.ObjTypeStatic {
		bullet.ToRecycle()
		effectParams[0].effectId = 1
		effectParams[0].effectFunc = bulletExplodeEffect
		effectParams[0].cx, effectParams[0].cy = bullet.Pos()
	} else if obj.Type() == object.ObjTypeMovable {
		bullet.ToRecycle()
		obj.ToRecycle()
		if obj.Subtype() == object.ObjSubtypeShell {
			effectParams[0].effectId = 1
			effectParams[0].effectFunc = bulletExplodeEffect
			effectParams[0].cx, effectParams[0].cy = bullet.Pos()
			effectParams[1].effectId = 1
			effectParams[1].effectFunc = bulletExplodeEffect
			effectParams[1].cx, effectParams[1].cy = obj.Pos()
		} else if obj.Subtype() == object.ObjSubtypeTank {
			effectParams[0].effectFunc = bigBulletExplodeEffect
			effectParams[0].effectId = 2
			effectParams[0].cx, effectParams[0].cy = obj.Pos()
		}
	}
	// 生成爆炸效果
	for i := 0; i < len(effectParams); i++ {
		if effectParams[i].effectId <= 0 {
			continue
		}
		effect := s.effectPool.Get(common_data.EffectConfigData[effectParams[i].effectId], effectParams[i].effectFunc /*s.pmap*/, s.gmap, bullet)
		effect.SetPos(effectParams[i].cx, effectParams[i].cy)
		s.effectList.Add(effect.InstId(), effect)
	}
}

func (s *SceneLogic) searchShellTarget(shell *object.Shell) object.IObject {
	staticInfo := shell.ShellStaticInfo()
	cx, cy := shell.Pos()
	rect := math.NewRect(cx-staticInfo.SearchTargetRadius, cy-staticInfo.SearchTargetRadius, cx+staticInfo.SearchTargetRadius, cy+staticInfo.SearchTargetRadius)
	objList := s.gmap.GetMovableObjListWithRangeAndSubtype(rect, object.ObjSubtypeTank)
	if len(objList) > 0 {
		log.Debug("searched shell target list: %v", objList)
	}
	var tank *object.Tank
	var o bool
	for i := 0; i < len(objList); i++ {
		tank, o = s.tankList.Get(objList[i])
		if o && tank.Camp() != shell.Camp() {
			return tank
		}
	}
	return nil
}

func (s *SceneLogic) testSearchShellTarget(shell *object.Shell) object.IObject {
	for i := int32(0); i < s.tankList.Count(); i++ {
		_, tank := s.tankList.GetByIndex(i)
		if tank.Camp() != shell.Camp() {
			return tank
		}
	}
	return nil
}
