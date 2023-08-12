package client_core

import (
	"project_b/common"
	"project_b/common/base"
	"project_b/common/object"
	"project_b/game_proto"
)

type GameLogic struct {
	common.GameLogic
	myId               uint64
	tankId             int32
	tankLevel          int32
	tankPosX, tankPosY int32
	tankOrientation    int32
	tankCurrSpeed      int32
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
	orientation := object.Dir2Orientation(object.Direction(tankProtoInfo.Direction))
	l.GameLogic.PlayerEnterWithStaticInfo(cplayer.Id(), tankProtoInfo.Id, tankProtoInfo.Level, tankProtoInfo.CurrPos.X, tankProtoInfo.CurrPos.Y, orientation, tankProtoInfo.CurrSpeed)
	l.tankId = tankProtoInfo.Id
	l.tankLevel = tankProtoInfo.Level
	l.tankPosX = tankProtoInfo.CurrPos.X
	l.tankPosY = tankProtoInfo.CurrPos.Y
	l.tankOrientation = orientation
	l.tankCurrSpeed = tankProtoInfo.CurrSpeed
}

func (l *GameLogic) MyPlayerTankMove(orientation int32) {
	l.PlayerTankMove(l.myId, orientation)
}

func (l *GameLogic) MyPlayerTankStopMove() {
	l.PlayerTankStopMove(l.myId)
}

func (l *GameLogic) MyPlayerTankFire(shellId int32) {
	l.PlayerTankFire(l.myId, shellId)
}

func (l *GameLogic) MyPlayerTankReleaseSurroundObj() {
	l.PlayerTankReleaseSurroundObj(l.myId)
}

func (l *GameLogic) MyPlayerTankRotate(angle int32) {
	l.PlayerTankRotate(l.myId, angle)
}

func (l *GameLogic) MyPlayerTankRevive() {
	l.PlayerTankRevive(l.myId, l.tankId, l.tankLevel, l.tankPosX, l.tankPosY, l.tankOrientation, l.tankCurrSpeed)
}

func (l *GameLogic) MyPlayerTankShield() {
	l.PlayerTankShield(l.myId)
}
