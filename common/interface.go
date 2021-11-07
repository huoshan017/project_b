package common

// 事件注册器接口
type IEventRegistrar interface {
	RegisterEvent(id EventId, handle EventHandle)
	UnregisterEvent(id EventId, handle EventHandle)
}

// 事件调用器
type IEventInvoker interface {
	InvokeEvent(id EventId, args ...interface{})
}

// 事件管理器
type IEventManager interface {
	IEventRegistrar
	IEventInvoker
}

// 游戏逻辑接口
type IGameLogic interface {
	EventMgr() IEventManager
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
