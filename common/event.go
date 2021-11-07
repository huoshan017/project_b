package common

import "reflect"

type EventId int32

type EventHandle func(args ...interface{})

type Event struct {
	handles []EventHandle
}

func NewEvent() *Event {
	return &Event{}
}

func (e *Event) Register(eh EventHandle) {
	for _, h := range e.handles {
		if reflect.DeepEqual(h, eh) {
			return
		}
	}
	e.handles = append(e.handles, eh)
}

func (e *Event) Unregister(eh EventHandle) {
	for i, h := range e.handles {
		if reflect.DeepEqual(h, eh) {
			e.handles = append(e.handles[:i], e.handles[i+1:]...)
			return
		}
	}
}

func (e *Event) Call(args ...interface{}) {
	for _, h := range e.handles {
		if h != nil {
			h(args...)
		}
	}
}

type EventManager struct {
	id2EventHandles map[EventId]*Event
}

func NewEventManager() *EventManager {
	return &(EventManager{id2EventHandles: make(map[EventId]*Event)})
}

func (e *EventManager) RegisterEvent(id EventId, handle EventHandle) {
	handles, o := e.id2EventHandles[id]
	if !o {
		handles = NewEvent()
		e.id2EventHandles[id] = handles
	}
	handles.Register(handle)
}

func (e *EventManager) UnregisterEvent(id EventId, handle EventHandle) {
	handles, o := e.id2EventHandles[id]
	if !o {
		return
	}
	handles.Unregister(handle)
}

func (e *EventManager) InvokeEvent(id EventId, args ...interface{}) {
	handles, o := e.id2EventHandles[id]
	if o {
		handles.Call(args...)
	}
}
