package object

import (
	"math"
	"project_b/common/base"
	"project_b/log"
)

type ObjectFactory struct {
	objIdIncrementer uint32             // 增量器
	objMap           map[uint32]IObject // 对象map
	isRecycleObjId   bool               // 是否回收对象id
	freeObjIds       []uint32           // 空闲id列表
	freeStaticObjs   []IStaticObject    // 空闲静态对象池
	freeMovableObjs  [][]IMovableObject // 空闲可运动物体对象池
}

func NewObjectFactory(isRecycleObjId bool) *ObjectFactory {
	return &ObjectFactory{
		objMap:          make(map[uint32]IObject),
		freeMovableObjs: make([][]IMovableObject, base.MovableObjEnumMax),
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

func (f *ObjectFactory) GetObj(instId uint32) IObject {
	return f.objMap[instId]
}

func (f *ObjectFactory) NewStaticObject(info *ObjStaticInfo) IStaticObject {
	if info.typ != base.ObjTypeStatic {
		log.Error("object type is invalid, must ObjTypeStatic")
		return nil
	}
	id := f.getNewObjId()
	l := len(f.freeStaticObjs)
	var obj IStaticObject
	if l == 0 {
		obj = NewStaticObject( /*id, info*/ )
	} else {
		obj = f.freeStaticObjs[l-1]
		f.freeStaticObjs = f.freeStaticObjs[:l-1]
	}
	obj.Init(id, info)
	f.objMap[id] = obj
	return obj
}

func (f *ObjectFactory) RecycleStaticObject(obj IStaticObject) bool {
	if obj.Type() != base.ObjTypeStatic {
		log.Error("object type is invalid, must ObjTypeStatic")
		return false
	}
	if _, o := f.objMap[obj.InstId()]; !o {
		return false
	}
	id := obj.InstId()
	obj.Uninit()
	if f.isRecycleObjId {
		f.freeObjIds = append(f.freeObjIds, id)
	}
	f.freeStaticObjs = append(f.freeStaticObjs, obj)
	delete(f.objMap, id)
	return true
}

func (f *ObjectFactory) NewTank(info *TankStaticInfo) *Tank {
	if info.typ != base.ObjTypeMovable && info.subType != base.ObjSubtypeTank {
		log.Error("object type and subtype is invalid, must ObjTypeMovable and ObjSubTypeTank")
		return nil
	}
	id := f.getNewObjId()
	l := len(f.freeMovableObjs[base.MovableObjTank])
	var obj *Tank
	if l == 0 {
		obj = NewTank()
	} else {
		obj = f.freeMovableObjs[base.MovableObjTank][l-1].(*Tank)
		f.freeMovableObjs[base.MovableObjTank] = f.freeMovableObjs[base.MovableObjTank][:l-1]
	}
	obj.Init(id, &info.ObjStaticInfo)
	f.objMap[id] = obj
	return obj
}

func (f *ObjectFactory) RecycleTank(tank *Tank) bool {
	res := f.recycleMovableObject(tank)
	if res {
		tank.Uninit()
	}
	return res
}

func (f *ObjectFactory) NewShell(info *ShellStaticInfo) *Shell {
	if info.typ != base.ObjTypeMovable && info.subType != base.ObjSubtypeShell {
		log.Error("object type and subtype is invalid, must ObjTypeMovable and ObjSubTypeShell")
		return nil
	}
	id := f.getNewObjId()
	l := len(f.freeMovableObjs[base.MovableObjShell])
	var obj *Shell
	if l == 0 {
		obj = NewShell()
	} else {
		obj = f.freeMovableObjs[base.MovableObjShell][l-1].(*Shell)
		f.freeMovableObjs[base.MovableObjShell] = f.freeMovableObjs[base.MovableObjShell][:l-1]
	}
	obj.Init(id, &info.ObjStaticInfo)
	f.objMap[id] = obj
	return obj
}

func (f *ObjectFactory) RecycleShell(shell *Shell) bool {
	res := f.recycleMovableObject(shell)
	if res {
		shell.Uninit()
	}
	return res
}

func (f *ObjectFactory) NewSurroundObj(info *SurroundObjStaticInfo) *SurroundObj {
	if info.typ != base.ObjTypeMovable && info.subType != base.ObjSubtypeSurroundObj {
		log.Error("object type and subtype is invalid, must ObjTypeMovable and ObjSubTypeSurroundObj")
		return nil
	}
	id := f.getNewObjId()
	l := len(f.freeMovableObjs[base.MovableObjSurroundObj])
	var obj *SurroundObj
	if l == 0 {
		obj = NewSurroundObj()
	} else {
		obj = f.freeMovableObjs[base.MovableObjSurroundObj][l-1].(*SurroundObj)
		f.freeMovableObjs[base.MovableObjSurroundObj] = f.freeMovableObjs[base.MovableObjSurroundObj][:l-1]
	}
	obj.Init(id, &info.ObjStaticInfo)
	f.objMap[id] = obj
	return obj
}

func (f *ObjectFactory) RecycleSurroundObj(sobj *SurroundObj) bool {
	res := f.recycleMovableObject(sobj)
	if res {
		sobj.Uninit()
	}
	return res
}

func (f *ObjectFactory) recycleMovableObject(mobj IMovableObject) bool {
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

func (f *ObjectFactory) Clear() {
	for _, obj := range f.objMap {
		if obj.Type() == base.ObjTypeStatic {
			f.RecycleStaticObject(obj.(*StaticObject))
			obj.Uninit()
		} else if obj.Type() == base.ObjTypeMovable {
			f.recycleMovableObject(obj.(*MovableObject))
			obj.Uninit()
		}
	}
	clear(f.objMap)
	f.freeObjIds = f.freeObjIds[:0]
	f.freeStaticObjs = f.freeStaticObjs[:0]
	for i := 0; i < len(f.freeMovableObjs); i++ {
		f.freeMovableObjs[i] = f.freeMovableObjs[i][:0]
	}
}
