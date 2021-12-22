package common

import (
	"project_b/common/base"
	"project_b/common/ds"
	"project_b/common/log"
	"project_b/common/object"
	"project_b/common/time"
	"project_b/game_map"
)

type TankKV struct {
	PlayerId uint64
	Tank     *object.Tank
}

// 场景结构必须在单个goroutine中执行
type Scene struct {
	gmap           *game_map.Config
	eventMgr       base.IEventManager
	playerTankList *ds.MapListUnion
	enemyTankList  *ds.MapListUnion
}

func NewScene(eventMgr base.IEventManager) *Scene {
	return &Scene{
		eventMgr:       eventMgr,
		playerTankList: ds.NewMapListUnion(),
		enemyTankList:  ds.NewMapListUnion(),
	}
}

func (s *Scene) LoadMap(m *game_map.Config) {
	s.gmap = m
}

func (s *Scene) GetPlayerTank(id uint64) *object.Tank {
	v, o := s.playerTankList.Get(id)
	if !o {
		return nil
	}
	return v.(*object.Tank)
}

func (s *Scene) AddPlayerTank(id uint64, tank *object.Tank) {
	s.playerTankList.Add(id, tank)
}

func (s *Scene) RemovePlayerTank(id uint64) {
	s.playerTankList.Remove(id)
}

func (s *Scene) GetPlayerTankList() []TankKV {
	var l []TankKV
	count := s.playerTankList.Count()
	for i := int32(0); i < count; i++ {
		k, v := s.playerTankList.GetByIndex(i)
		l = append(l, TankKV{
			PlayerId: k.(uint64),
			Tank:     v.(*object.Tank),
		})
	}
	return l
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

func (s *Scene) AddEnemyTank(id uint64, tank *object.Tank) {
	s.enemyTankList.Add(id, tank)
}

func (s *Scene) GetEnemyTank(id uint64) *object.Tank {
	v, o := s.enemyTankList.Get(id)
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

func (s *Scene) RegisterEnemyEvent(id uint64, eid base.EventId, handle func(args ...interface{})) {
	tank := s.GetPlayerTank(id)
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

func (s *Scene) UnregisterEnemyEvent(id uint64, eid base.EventId, handle func(args ...interface{})) {
	tank := s.GetPlayerTank(id)
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
		tank, _ := s.playerTankList.GetByIndex(i)
		s.RegisterPlayerEvent(tank.(*object.Tank).InstId(), eid, handle)
	}
}

func (s *Scene) UnregisterAllPlayersEvent(eid base.EventId, handle func(args ...interface{})) {
	count := s.playerTankList.Count()
	for i := int32(0); i < count; i++ {
		tank, _ := s.playerTankList.GetByIndex(i)
		s.UnregisterPlayerEvent(tank.(*object.Tank).InstId(), eid, handle)
	}
}

func (s *Scene) RegisterAllEnemiesEvent(eid base.EventId, handle func(args ...interface{})) {
	count := s.enemyTankList.Count()
	for i := int32(0); i < count; i++ {
		tank, _ := s.enemyTankList.GetByIndex(i)
		s.RegisterEnemyEvent(tank.(*object.Tank).InstId(), eid, handle)
	}
}

func (s *Scene) UnregisterAllEnemiesEvent(eid base.EventId, handle func(args ...interface{})) {
	count := s.enemyTankList.Count()
	for i := int32(0); i < count; i++ {
		tank, _ := s.enemyTankList.GetByIndex(i)
		s.UnregisterEnemyEvent(tank.(*object.Tank).InstId(), eid, handle)
	}
}
