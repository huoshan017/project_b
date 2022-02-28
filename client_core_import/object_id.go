package main

import "sync"

type ObjectId int32

var refs struct {
	sync.Mutex
	objs map[ObjectId]interface{}
	next ObjectId
}

func init() {
	refs.Lock()
	defer refs.Unlock()

	refs.objs = make(map[ObjectId]interface{})
	refs.next = 1000
}

func NewObjectId(obj interface{}) ObjectId {
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

func ObjectGet(id ObjectId) interface{} {
	refs.Lock()
	defer refs.Unlock()

	return refs.objs[id]
}

func ObjectFree(id ObjectId) interface{} {
	refs.Lock()
	defer refs.Unlock()

	obj := refs.objs[id]
	delete(refs.objs, id)
	return obj
}
