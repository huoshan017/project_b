package base

import "reflect"

type EventId int32

type EventHandle func(args ...any)

type Event struct {
	handles []EventHandle
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

func (e *Event) Call(args ...any) {
	for _, h := range e.handles {
		if h != nil {
			h(args...)
		}
	}
}

func (e *Event) Size() int32 {
	return int32(len(e.handles))
}

func (e *Event) Clear() {
	clear(e.handles)
	e.handles = e.handles[:0]
}

type eventData struct {
	eid  EventId
	args []any
}

type EventManager struct {
	id2EventHandles map[EventId]*Event
	edList          []eventData
	eventPool       *ObjectPool[Event]
}

func NewEventManager() *EventManager {
	return &EventManager{
		id2EventHandles: make(map[EventId]*Event),
		edList:          make([]eventData, 0),
		eventPool:       NewObjectPool[Event](),
	}
}

func (e *EventManager) Clear() {
	clear(e.id2EventHandles)
	if len(e.edList) > 0 {
		clear(e.edList)
		e.edList = nil
	}
}

func (e *EventManager) RegisterEvent(id EventId, handle EventHandle) {
	handles, o := e.id2EventHandles[id]
	if !o {
		handles = e.eventPool.Get()
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
	if handles.Size() == 0 {
		delete(e.id2EventHandles, id)
		e.eventPool.Put(handles)
	}
}

func (e *EventManager) InvokeEvent(id EventId, args ...any) {
	handles, o := e.id2EventHandles[id]
	if o {
		handles.Call(args...)
	}
}

func (e *EventManager) DispatchEvent(id EventId, args ...any) {
	e.edList = append(e.edList, eventData{eid: id, args: args})
}

func (e *EventManager) Update() {
	for _, ed := range e.edList {
		handles, o := e.id2EventHandles[ed.eid]
		if !o {
			continue
		}
		handles.Call(ed.args...)
	}
}
