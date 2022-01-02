package object

import (
	"math"
	"project_b/common/log"
	"unsafe"
)

type ObjectFactory struct {
	objIdIncrementer uint32             // 增量器
	objMap           map[uint32]IObject // 对象map
	isRecycleObjId   bool               // 是否回收对象id
	freeObjIds       []uint32           // 空闲id列表
	freeStaticObjs   []*StaticObject    // 空闲静态对象池
	freeMovableObjs  [][]*MovableObject // 空闲可运动物体对象池
}

func NewObjectFactory(isRecycleObjId bool) *ObjectFactory {
	return &ObjectFactory{
		objMap:          make(map[uint32]IObject),
		freeMovableObjs: make([][]*MovableObject, MovableObjEnumMax),
	}
}

func (f *ObjectFactory) incrObjId() uint32 {
	if f.objIdIncrementer == math.MaxUint32 {
		panic("Object Id overflow")
	}
	f.objIdIncrementer += 1
	return f.objIdIncrementer
}

func (f *ObjectFactory) getNewObjId() uint32 {
	l := len(f.freeObjIds)
	if !f.isRecycleObjId || l == 0 {
		return f.incrObjId()
	}
	id := f.freeObjIds[l-1]
	f.freeObjIds = f.freeObjIds[:l-1]
	return id
}

func (f *ObjectFactory) NewStaticObject(info *ObjStaticInfo) *StaticObject {
	if info.typ != ObjTypeStatic {
		log.Error("object type is invalid, must ObjTypeStatic")
		return nil
	}
	id := f.getNewObjId()
	l := len(f.freeStaticObjs)
	var obj *StaticObject
	if l == 0 {
		obj = NewStaticObject(id, info)
	} else {
		obj = f.freeStaticObjs[l-1]
		f.freeStaticObjs = f.freeStaticObjs[:l-1]
	}
	f.objMap[id] = obj
	return obj
}

func (f *ObjectFactory) RecycleStaticObject(obj *StaticObject) bool {
	if obj.Type() != ObjTypeStatic {
		log.Error("object type is invalid, must ObjTypeStatic")
		return false
	}
	if _, o := f.objMap[obj.InstId()]; !o {
		return false
	}
	if f.isRecycleObjId {
		f.freeObjIds = append(f.freeObjIds, obj.InstId())
	}
	f.freeStaticObjs = append(f.freeStaticObjs, obj)
	delete(f.objMap, obj.InstId())
	return true
}

func (f *ObjectFactory) NewTank(info *ObjStaticInfo) *Tank {
	if info.typ != ObjTypeMovable && info.subType != ObjSubTypeTank {
		log.Error("object type and subtype is invalid, must ObjTypeMovable and ObjSubTypeTank")
		return nil
	}
	id := f.getNewObjId()
	l := len(f.freeMovableObjs[MovableObjTank])
	var obj *Tank
	if l == 0 {
		obj = NewTank(id, info)
	} else {
		obj = (*Tank)(unsafe.Pointer((f.freeMovableObjs[MovableObjTank][l-1])))
		obj.Init(id, info)
		f.freeMovableObjs[MovableObjTank] = f.freeMovableObjs[MovableObjTank][:l-1]
	}
	f.objMap[id] = obj
	return obj
}

func (f *ObjectFactory) RecycleTank(tank *Tank) bool {
	tank.Uninit()
	mobj := (*MovableObject)(unsafe.Pointer(tank))
	return f.recycleMovableObject(mobj)
}

func (f *ObjectFactory) NewBullet(info *ObjStaticInfo) *Bullet {
	if info.typ != ObjTypeMovable && info.subType != ObjSubTypeBullet {
		log.Error("object type and subtype is invalid, must ObjTypeMovable and ObjSubTypeBullet")
		return nil
	}
	id := f.getNewObjId()
	l := len(f.freeMovableObjs[MovableObjBullet])
	var obj *Bullet
	if l == 0 {
		obj = NewBullet(id, info)
	} else {
		obj = (*Bullet)(unsafe.Pointer((f.freeMovableObjs[MovableObjBullet][l-1])))
		obj.Init(id, info)
		f.freeMovableObjs[MovableObjBullet] = f.freeMovableObjs[MovableObjBullet][:l-1]
	}
	f.objMap[id] = obj
	return obj
}

func (f *ObjectFactory) RecycleBullet(bullet *Bullet) bool {
	bullet.Uninit()
	mobj := (*MovableObject)(unsafe.Pointer(bullet))
	return f.recycleMovableObject(mobj)
}

func (f *ObjectFactory) recycleMovableObject(mobj *MovableObject) bool {
	if _, o := f.objMap[mobj.InstId()]; !o {
		return false
	}
	if f.isRecycleObjId {
		f.freeObjIds = append(f.freeObjIds, mobj.InstId())
	}
	f.freeMovableObjs[mobj.Subtype()] = append(f.freeMovableObjs[mobj.Subtype()], mobj)
	delete(f.objMap, mobj.InstId())
	return true
}
