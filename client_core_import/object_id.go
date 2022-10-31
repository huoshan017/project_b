package main

import "sync"

type ObjectId int32

var refs struct {
	sync.RWMutex
	objs map[ObjectId]any
	next ObjectId
}

func init() {
	refs.Lock()
	defer refs.Unlock()

	refs.objs = make(map[ObjectId]any)
	refs.next = 1000
}

func NewObjectId(obj any) ObjectId {
	refs.Lock()
	defer refs.Unlock()

	id := refs.next
	refs.next++

	refs.objs[id] = obj
	return id
}

func ObjectIsNil(id ObjectId) bool {
	return id == 0
}

func ObjectGet(id ObjectId) any {
	refs.RLock()
	defer refs.RUnlock()

	return refs.objs[id]
}

func ObjectFree(id ObjectId) any {
	refs.Lock()
	defer refs.Unlock()

	obj := refs.objs[id]
	delete(refs.objs, id)
	return obj
}
