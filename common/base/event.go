package base

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

type eventData struct {
	eid  EventId
	args []interface{}
}

type EventManager struct {
	id2EventHandles map[EventId]*Event
	edList          []*eventData
}

func NewEventManager() *EventManager {
	return &EventManager{
		id2EventHandles: make(map[EventId]*Event),
		edList:          make([]*eventData, 0),
	}
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

func (e *EventManager) DispatchEvent(id EventId, args ...interface{}) {
	e.edList = append(e.edList, &eventData{eid: id, args: args})
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
