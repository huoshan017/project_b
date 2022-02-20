package common

import (
	"project_b/common/base"
	"project_b/common/ds"
	"project_b/common/log"
	"project_b/common/object"
	"project_b/common/time"
	"project_b/common_data"
	"project_b/game_map"
	"project_b/utils"
)

type PlayerTankKV struct {
	PlayerId uint64
	Tank     *object.Tank
}

// 场景结构必须在单个goroutine中执行
type Scene struct {
	gmap                *game_map.Config
	tileObjectArray     [][]*object.StaticObject
	eventMgr            base.IEventManager
	playerTankList      *ds.MapListUnion
	enemyTankList       *ds.MapListUnion
	playerTankListCache []PlayerTankKV
	objFactory          *object.ObjectFactory
}

func NewScene(eventMgr base.IEventManager) *Scene {
	return &Scene{
		eventMgr:       eventMgr,
		playerTankList: ds.NewMapListUnion(),
		enemyTankList:  ds.NewMapListUnion(),
		objFactory:     object.NewObjectFactory(true),
	}
}

func (s *Scene) LoadMap(m *game_map.Config) bool {
	// 载入地图
	if s.tileObjectArray == nil {
		s.tileObjectArray = make([][]*object.StaticObject, len(m.Layers))
	}
	// 地图载入前事件
	s.eventMgr.InvokeEvent(EventIdBeforeMapLoad)
	for i := 0; i < len(s.tileObjectArray); i++ {
		s.tileObjectArray[i] = make([]*object.StaticObject, len(m.Layers[i]))
		for j := 0; j < len(s.tileObjectArray[i]); j++ {
			st := object.StaticObjType(m.Layers[i][j])
			s.tileObjectArray[i][j] = s.objFactory.NewStaticObject(common_data.StaticObjectConfigData[st])
			if s.tileObjectArray[i][j] == nil {
				log.Error("Can't create static object %v by map layer (%v, %v) in loading map", st, i, j)
				continue
			}
		}
	}
	s.gmap = m
	// 地图载入完成事件
	s.eventMgr.InvokeEvent(EventIdMapLoaded, m.Id, s.tileObjectArray)
	log.Info("Load map %v done", m.Id)
	return true
}

func (s *Scene) UnloadMap() {
	if s.tileObjectArray == nil {
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
	s.tileObjectArray = nil
}

func (s *Scene) GetPlayerTank(pid uint64) *object.Tank {
	v, o := s.playerTankList.Get(pid)
	if !o {
		return nil
	}
	return v.(*object.Tank)
}

func (s *Scene) NewPlayerTank(pid uint64) *object.Tank {
	tank := s.objFactory.NewTank(&s.gmap.PlayerTankInitData)
	// 随机并设置坦克位置
	pos := utils.RandomPosInRect(s.gmap.PlayerTankInitRect)
	tank.SetPos(pos.X, pos.Y)
	// 加入到玩家坦克列表
	s.playerTankList.Add(pid, tank)
	return tank
}

func (s *Scene) AddPlayerTank(pid uint64, tank *object.Tank) {
	s.playerTankList.Add(pid, tank)
}

func (s *Scene) RemovePlayerTank(pid uint64) {
	tank := s.playerTankList.Remove(pid)
	s.objFactory.RecycleTank(tank.(*object.Tank))
}

func (s *Scene) GetPlayerTankList() []PlayerTankKV {
	if s.playerTankListCache != nil && len(s.playerTankListCache) > 0 {
		s.playerTankListCache = s.playerTankListCache[:0]
	}
	count := s.playerTankList.Count()
	for i := int32(0); i < count; i++ {
		k, v := s.playerTankList.GetByIndex(i)
		s.playerTankListCache = append(s.playerTankListCache, PlayerTankKV{
			PlayerId: k.(uint64),
			Tank:     v.(*object.Tank),
		})
	}
	return s.playerTankListCache
}

func (s *Scene) PlayerTankMove(uid uint64, dir object.Direction) {
	tank := s.GetPlayerTank(uid)
	if tank == nil {
		log.Error("player %v tank not found", uid)
		return
	}
	tank.Move(dir)
}

func (s *Scene) PlayerTankStopMove(uid uint64) {
	tank := s.GetPlayerTank(uid)
	if tank == nil {
		log.Error("player %v tank not found", uid)
		return
	}
	tank.Stop()
}

func (s *Scene) PlayerTankChange(uid uint64, staticInfo *object.ObjStaticInfo) bool {
	tank := s.GetPlayerTank(uid)
	if tank == nil {
		return false
	}
	tank.Change(staticInfo)
	return true
}

func (s *Scene) PlayerTankRestore(uid uint64) int32 {
	tank := s.GetPlayerTank(uid)
	if tank == nil {
		return 0
	}
	tank.Restore()
	return tank.Id()
}

func (s *Scene) AddEnemyTank(instId uint32, tank *object.Tank) {
	s.enemyTankList.Add(instId, tank)
}

func (s *Scene) GetEnemyTank(instId uint32) *object.Tank {
	v, o := s.enemyTankList.Get(instId)
	if !o {
		return nil
	}
	return v.(*object.Tank)
}

func (s *Scene) RemoveEnemyTank(id uint64) {
	s.enemyTankList.Remove(id)
}

func (s *Scene) Update(tick time.Duration) {
	count := s.playerTankList.Count()
	for i := int32(0); i < count; i++ {
		_, tank := s.playerTankList.GetByIndex(i)
		tank.(*object.Tank).Update(tick)
	}
	count = s.enemyTankList.Count()
	for i := int32(0); i < count; i++ {
		_, tank := s.enemyTankList.GetByIndex(i)
		tank.(*object.Tank).Update(tick)
	}
}

// 注册事件
func (s *Scene) RegisterEvent(eid base.EventId, handle func(args ...interface{})) {
	s.eventMgr.RegisterEvent(eid, handle)
}

// 注销事件
func (s *Scene) UnregisterEvent(eid base.EventId, handle func(args ...interface{})) {
	s.eventMgr.UnregisterEvent(eid, handle)
}

// 注册玩家事件
func (s *Scene) RegisterPlayerEvent(uid uint64, eid base.EventId, handle func(args ...interface{})) {
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
func (s *Scene) UnregisterPlayerEvent(uid uint64, eid base.EventId, handle func(args ...interface{})) {
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

func (s *Scene) RegisterEnemyEvent(instId uint32, eid base.EventId, handle func(args ...interface{})) {
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

func (s *Scene) UnregisterEnemyEvent(instId uint32, eid base.EventId, handle func(args ...interface{})) {
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

func (s *Scene) RegisterAllPlayersEvent(eid base.EventId, handle func(args ...interface{})) {
	count := s.playerTankList.Count()
	for i := int32(0); i < count; i++ {
		pid, _ := s.playerTankList.GetByIndex(i)
		s.RegisterPlayerEvent(pid.(uint64), eid, handle)
	}
}

func (s *Scene) UnregisterAllPlayersEvent(eid base.EventId, handle func(args ...interface{})) {
	count := s.playerTankList.Count()
	for i := int32(0); i < count; i++ {
		pid, _ := s.playerTankList.GetByIndex(i)
		s.UnregisterPlayerEvent(pid.(uint64), eid, handle)
	}
}

func (s *Scene) RegisterAllEnemiesEvent(eid base.EventId, handle func(args ...interface{})) {
	count := s.enemyTankList.Count()
	for i := int32(0); i < count; i++ {
		instId, _ := s.enemyTankList.GetByIndex(i)
		s.RegisterEnemyEvent(instId.(uint32), eid, handle)
	}
}

func (s *Scene) UnregisterAllEnemiesEvent(eid base.EventId, handle func(args ...interface{})) {
	count := s.enemyTankList.Count()
	for i := int32(0); i < count; i++ {
		instId, _ := s.enemyTankList.GetByIndex(i)
		s.UnregisterEnemyEvent(instId.(uint32), eid, handle)
	}
}
