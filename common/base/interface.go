package base

// 事件接口
type IEvent interface {
	Register(handle EventHandle)
	Unregister(handle EventHandle)
	Call(args ...any)
}

// 事件注册器接口
type IEventRegistrar interface {
	RegisterEvent(id EventId, handle EventHandle)
	UnregisterEvent(id EventId, handle EventHandle)
}

// 事件调用器
type IEventInvoker interface {
	InvokeEvent(id EventId, args ...any)
}

// 事件派发器
type IEventDispatcher interface {
	DispatchEvent(id EventId, args ...any)
	Update()
}

// 事件管理器
type IEventManager interface {
	IEventRegistrar
	IEventInvoker
	IEventDispatcher
}
