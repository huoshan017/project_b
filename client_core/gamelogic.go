package client_core

import (
	"project_b/common"
	"project_b/common/base"
	"project_b/common/object"
	"project_b/game_proto"
)

type GameLogic struct {
	common.GameLogic
	myId uint64
}

func CreateGameLogic(eventMgr base.IEventManager) *GameLogic {
	return &GameLogic{
		GameLogic: *common.NewGameLogic(eventMgr),
	}
}

func (l *GameLogic) SetMyId(id uint64) {
	l.myId = id
}

func (l *GameLogic) GetMyId() uint64 {
	return l.myId
}

// 玩家坦克进入
func (l *GameLogic) PlayerEnterWithTankInfo(cplayer *CPlayer, tankProtoInfo *game_proto.TankInfo) {
	if l.CurrentSceneMap().GetPlayerTank(cplayer.Id()) == nil {
		l.CurrentSceneMap().AddPlayerTankWithInfo(cplayer.Id(), tankProtoInfo.Id, tankProtoInfo.Level, tankProtoInfo.CurrPos.X, tankProtoInfo.CurrPos.Y, object.Direction(tankProtoInfo.Direction), tankProtoInfo.CurrSpeed)
	}
}

func (l *GameLogic) MyPlayerTankMove(moveDir object.Direction) {
	l.PlayerTankMove(l.myId, moveDir)
}

func (l *GameLogic) MyPlayerTankStopMove() {
	l.PlayerTankStopMove(l.myId)
}
