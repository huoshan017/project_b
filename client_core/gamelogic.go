package client_core

import (
	"project_b/common"
	"project_b/common/object"
)

type GameLogic struct {
	common.GameLogic
	myId uint64
}

func CreateGameLogic() *GameLogic {
	return &GameLogic{
		GameLogic: *common.NewGameLogic(),
	}
}

func (l *GameLogic) SetMyId(id uint64) {
	l.myId = id
}

func (l *GameLogic) GetMyId() uint64 {
	return l.myId
}

func (l *GameLogic) MyPlayerTankMove(moveDir object.Direction) {
	l.GameLogic.PlayerTankMove(l.myId, moveDir)
}

func (l *GameLogic) MyPlayerTankStopMove() {
	l.GameLogic.PlayerTankStopMove(l.myId)
}
