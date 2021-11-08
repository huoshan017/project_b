package common

import (
	"log"
	"project_b/common/object"
	"project_b/game_map"
	"time"
)

const (
	InitTankListLength = 10
)

// 场景结构必须在单个goroutine中执行
type Scene struct {
	gmap           *game_map.Config
	eventMgr       IEventManager
	playerTanks    map[uint64]*object.Tank
	playerTankList []*object.Tank
	enemyTanks     map[int32]*object.Tank
	enemyTankList  []*object.Tank
}

func NewScene(eventMgr IEventManager) *Scene {
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
		log.Printf("player %v tank not found", uid)
		return
	}
	tank.Move(dir)
}

func (s *Scene) PlayerTankStopMove(uid uint64) {
	tank := s.playerTanks[uid]
	if tank == nil {
		log.Printf("player %v tank not found", uid)
		return
	}
	tank.Stop()
}

func (s *Scene) PlayerTankChange(uid uint64, staticInfo *object.ObjStaticInfo) bool {
	tank := s.playerTanks[uid]
	if tank == nil {
		return false
	}
	tank.ChangeStaticInfo(staticInfo)
	tank.SetCurrentSpeed(staticInfo.Speed())
	return true
}

func (s *Scene) PlayerTankRestore(uid uint64) int32 {
	tank := s.playerTanks[uid]
	if tank == nil {
		return 0
	}
	tank.RestoreStaticInfo()
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
