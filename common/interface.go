package common

import (
	"project_b/common/base"
)

// 游戏逻辑接口
type IGameLogic interface {
	EventMgr() base.IEventManager
	Update()
}

// 玩家接口
type IPlayer interface {
	Id() uint64
	Account() string
	Token() string
	SetId(id uint64)
	SetAccount(account string)
	SetToken(token string)
	Entered()
	Left(force bool)
	IsEntered() bool
	IsLeft() bool
}

// 服务器玩家接口
type ISPlayer interface {
	IPlayer
}
