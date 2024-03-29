package object

import (
	"project_b/common/base"
	"testing"
)

func TestObjectFactory(t *testing.T) {
	var of = NewObjectFactory(true)
	for i := 0; i < 1000; i++ {
		so := of.NewStaticObject(&ObjStaticInfo{typ: base.ObjTypeStatic})
		of.RecycleStaticObject(so)
	}
	for i := 0; i < 1000; i++ {
		tank := of.NewTank(&TankStaticInfo{
			MovableObjStaticInfo: MovableObjStaticInfo{
				ObjStaticInfo: ObjStaticInfo{
					typ:     base.ObjTypeMovable,
					subType: base.ObjSubtypeTank,
					id:      1,
					w:       100,
					l:       100,
					//dir:     DirUp,
					speed: 10,
				},
			},
		})
		of.RecycleTank(tank)
	}
	for i := 0; i < 1000; i++ {
		bullet := of.NewShell(&ShellStaticInfo{
			MovableObjStaticInfo: MovableObjStaticInfo{
				ObjStaticInfo: ObjStaticInfo{typ: base.ObjTypeMovable, subType: base.ObjSubtypeShell},
			},
		})
		of.RecycleShell(bullet)
	}
}
