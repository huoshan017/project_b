package main

import "project_b/common/object"

/****************** 游戏主逻辑消息 ******************/

const (
	// 坦克进入同步
	MsgIdPlayerTankEnterSync = "MsgIdPlayerTankEnterSync"
	// 坦克进入回应
	MsgIdPlayerTankEnterAck = "MsgIdPlayerTankEnterAck"
	// 坦克离开同步
	MsgIdPlayerTankLeaveSync = "MsgIdPlayerTankLeaveSync"
	// 坦克离开回应
	MsgIdPlayerTankLeaveAck = "MsgIdPlayerTankLeaveAck"
	// 坦克移动同步
	MsgIdPlayerTankMoveSync = "MsgIdPlayerTankMoveSync"
	// 坦克移动回应
	MsgIdPlayerTankMoveAck = "MsgIdPlayerTankMoveAck"
	// 敌人移动坦克更新
	MsgIdEnemyTankMoveUpdateNotify = "MsgIdEnemyTankMoveNotify"
	// 玩家坦克改变请求
	MsgIdPlayerTankChangeReq = "MsgIdPlayerTankChangeReq"
	// 玩家坦克改变回应
	MsgIdPlayerTankChangeAck = "MsgIdPlayerTankChangeAck"
)

// 玩家坦克进入同步
type MsgPlayerTankEnterSync struct {
	playerId uint64
}

// 玩家坦克进入回应
type MsgPlayerTankEnterAck struct {
	result int32
}

// 玩家坦克离开同步
type MsgPlayerTankLeaveSync struct {
	playerId uint64
}

// 玩家坦克离开回应
type MsgPlayerTankLeaveAck struct {
	result int32
}

// 玩家坦克移动同步
type MsgPlayerTankMoveSync struct {
	playerId  uint64
	direction object.Direction
}

// 玩家坦克移动回应
type MsgPlayerTankMoveAck struct {
	playerId  uint64
	direction object.Direction
}

// 敌人坦克移动更新通知
type MsgEnemyTankMoveUpdateNotify struct {
}

// 玩家坦克改变
type MsgPlayerTankChange struct {
	playerId uint64
}

/********************** 好友消息 ********************/

/********************** 聊天消息 ********************/

/********************** 公会消息 ********************/
