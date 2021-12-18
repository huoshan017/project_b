package common

import (
	"project_b/common/base"
	"project_b/common/log"
	"project_b/common/object"
	"project_b/common/time"
	"project_b/game_map"
)

const (
	InitTankListLength = 10
)

// 场景结构必须在单个goroutine中执行
type Scene struct {
	gmap           *game_map.Config
	eventMgr       base.IEventManager
	playerTanks    map[uint64]*object.Tank
	playerTankList []*object.Tank
	enemyTanks     map[int32]*object.Tank
	enemyTankList  []*object.Tank
}

func NewScene(eventMgr base.IEventManager) *Scene {
	return &Scene{
		eventMgr:       eventMgr,
		playerTanks:    make(map[uint64]*object.Tank),
		playerTankList: make([]*object.Tank, InitTankListLength),
		enemyTanks:     make(map[int32]*object.Tank),
		enemyTankList:  make([]*object.Tank, InitTankListLength),
	}
}

func (s *Scene) LoadMap(m *game_map.Config) {
	s.gmap = m
}

func (s *Scene) GetPlayerTank(id uint64) *object.Tank {
	return s.playerTanks[id]
}

func (s *Scene) AddPlayerTank(id uint64, tank *object.Tank) {
	s.playerTanks[id] = tank
}

func (s *Scene) RemovePlayerTank(id uint64) {
	delete(s.playerTanks, id)
}

func (s *Scene) GetPlayerTankList() []*object.Tank {
	s.playerTankList = s.playerTankList[:0]
	for _, t := range s.playerTanks {
		s.playerTankList = append(s.playerTankList, t)
	}
	return s.playerTankList
}

func (s *Scene) GetPlayerTanks() map[uint64]*object.Tank {
	return s.playerTanks
}

func (s *Scene) PlayerTankMove(uid uint64, dir object.Direction) {
	tank := s.playerTanks[uid]
	if tank == nil {
		log.Error("player %v tank not found", uid)
		return
	}
	tank.Move(dir)
}

func (s *Scene) PlayerTankStopMove(uid uint64) {
	tank := s.playerTanks[uid]
	if tank == nil {
		log.Error("player %v tank not found", uid)
		return
	}
	tank.Stop()
}

func (s *Scene) PlayerTankChange(uid uint64, staticInfo *object.ObjStaticInfo) bool {
	tank := s.playerTanks[uid]
	if tank == nil {
		return false
	}
	tank.Change(staticInfo)
	return true
}

func (s *Scene) PlayerTankRestore(uid uint64) int32 {
	tank := s.playerTanks[uid]
	if tank == nil {
		return 0
	}
	tank.Restore()
	return tank.Id()
}

func (s *Scene) GetEnemyTank(id int32) *object.Tank {
	return s.enemyTanks[id]
}

func (s *Scene) GetEnemyTankList() []*object.Tank {
	return s.enemyTankList[:0]
}

func (s *Scene) GetEnemyTanks() map[int32]*object.Tank {
	return s.enemyTanks
}

func (s *Scene) Update(tick time.Duration) {
	for _, tank := range s.playerTanks {
		tank.Update(tick)
	}
	for _, tank := range s.enemyTanks {
		tank.Update(tick)
	}
}

func (s *Scene) RegisterPlayerEvent(uid uint64, eid base.EventId, handle func(args ...interface{})) {
	tank, o := s.playerTanks[uid]
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

func (s *Scene) UnregisterPlayerEvent(uid uint64, eid base.EventId, handle func(args ...interface{})) {
	tank, o := s.playerTanks[uid]
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

func (s *Scene) RegisterEnemyEvent(id int32, eid base.EventId, handle func(args ...interface{})) {
	tank, o := s.enemyTanks[id]
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

func (s *Scene) UnregisterEnemyEvent(id int32, eid base.EventId, handle func(args ...interface{})) {
	tank, o := s.enemyTanks[id]
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

func (s *Scene) RegisterAllPlayersEvent(eid base.EventId, handle func(args ...interface{})) {
	for id := range s.playerTanks {
		s.RegisterPlayerEvent(id, eid, handle)
	}
}

func (s *Scene) UnregisterAllPlayersEvent(eid base.EventId, handle func(args ...interface{})) {
	for id := range s.playerTanks {
		s.UnregisterPlayerEvent(id, eid, handle)
	}
}

func (s *Scene) RegisterAllEnemiesEvent(eid base.EventId, handle func(args ...interface{})) {
	for id := range s.enemyTanks {
		s.RegisterEnemyEvent(id, eid, handle)
	}
}

func (s *Scene) UnregisterAllEnemiesEvent(eid base.EventId, handle func(args ...interface{})) {
	for id := range s.enemyTanks {
		s.UnregisterEnemyEvent(id, eid, handle)
	}
}
