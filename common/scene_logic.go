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
	mapConfig                            *game_map.Config                               // 地圖配置
	mapWidth, mapHeight                  int32                                          // 地圖寬高
	eventMgr                             base.IEventManager                             // 事件管理器
	staticObjList                        *ds.MapListUnion[uint32, *object.StaticObject] // 靜態對象列表，用map和list的聯合體是爲了遍歷時的有序性
	tankList                             *ds.MapListUnion[uint32, *object.Tank]         // 坦克列表，不區分玩家和BOT
	bulletList                           *ds.MapListUnion[uint32, *object.Bullet]       // 炮彈列表，不區分坦克的炮彈
	objFactory                           *object.ObjectFactory                          // 對象池
	effectList                           *ds.MapListUnion[uint32, *object.Effect]       // 效果列表
	effectPool                           *object.EffectPool                             // 效果池
	pmap                                 *PartitionMap                                  // 分區地圖
	objCreatedEvent, objRemovedEvent     base.Event                                     // 對象創建刪除事件
	effectAddedEvent, effectRemovedEvent base.Event                                     // 效果添加刪除事件
	staticObjRecycleList                 []*object.StaticObject                         // 靜態對象回收列表
	tankRecycleList                      []*object.Tank                                 // 坦克對象回收列表
	bulletRecycleList                    []*object.Bullet                               // 炮彈對象回收列表
	effectRecycleList                    []*object.Effect                               // 效果回收列表
}

func NewSceneLogic(eventMgr base.IEventManager) *SceneLogic {
	return &SceneLogic{
		eventMgr:      eventMgr,
		staticObjList: ds.NewMapListUnion[uint32, *object.StaticObject](),
		tankList:      ds.NewMapListUnion[uint32, *object.Tank](),
		bulletList:    ds.NewMapListUnion[uint32, *object.Bullet](),
		objFactory:    object.NewObjectFactory(true),
		effectList:    ds.NewMapListUnion[uint32, *object.Effect](),
		effectPool:    object.NewEffectPool(),
		pmap:          NewPartitionMap(0),
	}
}

func (s *SceneLogic) GetMapId() int32 {
	return s.mapConfig.Id
}

func (s *SceneLogic) LoadMap(m *game_map.Config) bool {
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
			// 二維數組Y軸是自上而下的，而世界坐標Y軸是自下而上的，所以設置Y坐標要倒過來
			tileObj.SetPos(m.TileWidth*int32(col), m.TileHeight*int32(len(m.Layers)-1-line))
			s.objCreatedEvent.Call(tileObj)
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

func (s *SceneLogic) UnloadMap() {
	// 地图卸载前事件
	s.eventMgr.InvokeEvent(EventIdBeforeMapUnload)
	s.mapWidth = 0
	s.mapHeight = 0
	for i := int32(0); i < s.staticObjList.Count(); i++ {
		_, v := s.staticObjList.GetByIndex(i)
		if v != nil {
			s.objRemovedEvent.Call(v)
			s.objFactory.RecycleStaticObject(v)
		}
	}
	for i := int32(0); i < s.tankList.Count(); i++ {
		_, v := s.tankList.GetByIndex(i)
		if v != nil {
			s.objRemovedEvent.Call(v)
			s.objFactory.RecycleTank(v)
		}
	}
	for i := int32(0); i < s.bulletList.Count(); i++ {
		_, v := s.bulletList.GetByIndex(i)
		if v != nil {
			s.objRemovedEvent.Call(v)
			s.objFactory.RecycleBullet(v)
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
	s.bulletList.Clear()
	s.pmap.Unload()
	s.objFactory.Clear()
	s.effectList.Clear()
	s.effectPool.Clear()
	// 地图卸载后事件
	s.eventMgr.InvokeEvent(EventIdMapUnloaded)
}

func (s *SceneLogic) RegisterNewObjCreatedHandle(handle func(...any)) {
	s.objCreatedEvent.Register(handle)
}

func (s *SceneLogic) UnregisterNewObjCreatedHandle(handle func(...any)) {
	s.objCreatedEvent.Unregister(handle)
}

func (s *SceneLogic) RegisterObjRemovedHandle(handle func(...any)) {
	s.objRemovedEvent.Register(handle)
}

func (s *SceneLogic) UnregisterObjRemovedHandle(handle func(...any)) {
	s.objRemovedEvent.Unregister(handle)
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
	return s.pmap.GetLayerObjsWithRange(rect)
}

func (s *SceneLogic) GetObj(instId uint32) object.IObject {
	return s.objFactory.GetObj(instId)
}

func (s *SceneLogic) GetEffectListWithRange(rect *math.Rect) []uint32 {
	var effectIdList []uint32
	count := s.effectList.Count()
	for i := int32(0); i < count; i++ {
		instId, effect := s.effectList.GetByIndex(i)
		ex, ey := effect.Center()
		el := ex - effect.Width()/2
		er := ex + effect.Width()/2
		et := ey + effect.Height()/2
		eb := ey - effect.Height()/2
		if !(er <= rect.X() || el >= rect.X()+rect.W() || et <= rect.Y() || eb >= rect.Y()+rect.H()) {
			effectIdList = append(effectIdList, instId)
		}
	}
	return effectIdList
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
	s.pmap.AddObj(tank)
	// 物體創建事件
	s.objCreatedEvent.Call(tank)
	return tank
}

func (s *SceneLogic) AddTank(tank *object.Tank) {
	tank.RegisterCheckMoveEventHandle(s.checkObjMoveEventHandle)
	s.tankList.Add(tank.InstId(), tank)
	s.pmap.AddObj(tank)
}

func (s *SceneLogic) NewTankWithStaticInfo(id int32, level int32, x, y int32, dir object.Direction, currSpeed int32) uint32 {
	tank := s.objFactory.NewTank(common_data.TankConfigData[id])
	tank.SetPos(x, y)
	tank.SetLevel(level)
	tank.SetDir(dir)
	tank.SetCurrentSpeed(currSpeed)
	// 注冊檢測移動事件處理
	tank.RegisterCheckMoveEventHandle(s.checkObjMoveEventHandle)
	// 加入玩家坦克列表
	s.tankList.Add(tank.InstId(), tank)
	// 加入網格分區地圖
	s.pmap.AddObj(tank)
	s.objCreatedEvent.Call(tank)
	return tank.InstId()
}

func (s *SceneLogic) RemoveTank(instId uint32) {
	tank := s.tankList.Remove(instId)
	tank.UnregisterCheckMoveEventHandle(s.checkObjMoveEventHandle)
	s.pmap.RemoveObj(tank.InstId())
	s.objRemovedEvent.Call(tank)
	s.objFactory.RecycleTank(tank)
}

func (s *SceneLogic) TankMove(instId uint32, dir object.Direction) {
	tank, o := s.tankList.Get(instId)
	if !o {
		log.Error("tank %v not found", instId)
		return
	}
	tank.SetDir(dir)
	tank.Move(dir)
}

func (s *SceneLogic) TankStopMove(instId uint32) {
	tank, o := s.tankList.Get(instId)
	if !o {
		log.Error("player %v tank not found", instId)
		return
	}
	tank.Stop()
}

func (s *SceneLogic) TankFire(instId uint32) {
	tank, o := s.tankList.Get(instId)
	if !o {
		log.Error("player %v tank not found", instId)
		return
	}
	bulletConfig := common_data.BulletConfigData[tank.GetBulletConfig().BulletId]
	bullet := tank.CheckAndFire(s.objFactory.NewBullet, bulletConfig)
	if bullet != nil {
		s.bulletList.Add(bullet.InstId(), bullet)
		bullet.RegisterCheckMoveEventHandle(s.checkObjMoveEventHandle)
		collider := bullet.GetComp("Collider")
		if collider != nil {
			c := collider.(*object.ColliderComp)
			c.SetCollisionHandle(s.onBulletCollision)
		}
		// 子彈移動
		bullet.Move(tank.Dir())
		s.pmap.AddObj(bullet)
	}
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
		s.pmap.UpdateMovable(tank) // 相當於MonoBehevior.FixedUpdate
		if tank.IsRecycle() {
			s.tankRecycleList = append(s.tankRecycleList, tank)
		}
	}

	count = s.bulletList.Count()
	for i := int32(0); i < count; i++ {
		_, bullet := s.bulletList.GetByIndex(i)
		bullet.Update(tick)
		s.pmap.UpdateMovable(bullet)
		if bullet.IsRecycle() {
			s.bulletRecycleList = append(s.bulletRecycleList, bullet)
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
			s.pmap.RemoveObj(obj.InstId())
			s.objRemovedEvent.Call(obj)
			s.objFactory.RecycleStaticObject(obj)
		}
		s.staticObjRecycleList = s.staticObjRecycleList[:0]
	}

	if len(s.tankRecycleList) > 0 {
		for _, tank := range s.tankRecycleList {
			s.tankList.Remove(tank.InstId())
			s.pmap.RemoveObj(tank.InstId())
			s.objRemovedEvent.Call(tank)
			s.objFactory.RecycleTank(tank)
		}
		s.tankRecycleList = s.tankRecycleList[:0]
	}

	if len(s.bulletRecycleList) > 0 {
		for _, bullet := range s.bulletRecycleList {
			s.bulletList.Remove(bullet.InstId())
			s.pmap.RemoveObj(bullet.InstId())
			s.objRemovedEvent.Call(bullet)
			s.objFactory.RecycleBullet(bullet)
		}
		s.bulletRecycleList = s.bulletRecycleList[:0]
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
	dir := args[1].(object.Direction)
	distance := args[2].(float64)
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
	if !s.checkObjMoveRange(mobj, dir, distance, &x, &y) {
		mobj.SetPos(x, y)
		mobj.Stop()
		*isMove = false
		*isCollision = false
		s.onMovableObjReachMapBorder(mobj)
	} else if checkMovableObjCollision(s.pmap, mobj, dir, distance, resObj) {
		mobj.Stop()
		*isCollision = true
		*isMove = false
	} else {
		*isMove = true
		*isCollision = false
	}
}

func (s *SceneLogic) checkObjMoveRange(obj object.IMovableObject, dir object.Direction, distance float64, rx, ry *int32) bool {
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

func (s *SceneLogic) onMovableObjReachMapBorder(mobj object.IMovableObject) {
	if mobj.Subtype() == object.ObjSubTypeBullet {
		mobj.ToRecycle()
	}
}

func (s *SceneLogic) onBulletCollision(args ...any) {
	// todo 處理子彈碰撞
	bullet := args[0].(*object.Bullet)
	obj := args[1].(object.IObject)
	if obj.Type() == object.ObjTypeStatic {
		bullet.ToRecycle()
	} else if obj.Type() == object.ObjTypeMovable {
		bullet.ToRecycle()
		obj.ToRecycle()
	}
	// 生成爆炸效果
	effect := s.effectPool.Get(common_data.EffectConfigData[1], bulletExplodeEffect, s.pmap, bullet)
	effect.SetCenter(bullet.Center())
	s.effectList.Add(effect.InstId(), effect)
}
