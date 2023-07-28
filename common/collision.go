package common

import (
	"project_b/common/object"
)

// checkMovableObjCollisionObj 檢查可移動物體和物體是否碰撞
func checkMovableObjCollisionObj(mobj object.IMovableObject, comp object.IComponent, dx, dy int32, obj object.IObject) bool {
	if !(mobj.StaticInfo().Layer() == obj.StaticInfo().Layer() ||
		(obj.Type() == object.ObjTypeStatic && (obj.Subtype() == object.ObjSubtypeWater || obj.Subtype() == object.ObjSubtypeIce))) {
		return false
	}

	var (
		collisionComp *object.ColliderComp
		aabb1         object.AABB
	)
	collisionComp = comp.(*object.ColliderComp)
	aabb1 = collisionComp.GetAABB(mobj)
	aabb1.Move(int32(dx), int32(dy))
	comp2 := obj.GetComp("Collider")
	if comp2 == nil {
		return false
	}
	collisionComp2 := comp2.(*object.ColliderComp)
	if collisionComp2 == nil {
		return false
	}
	aabb2 := collisionComp2.GetAABB(obj)
	if aabb1.MoveIntersect(dx, dy, &aabb2) {
		if onMovableObjCollisionObj(mobj, obj) {
			if isCollisionCorrectPosition(mobj, obj) {
				mx, my := mobj.Pos()
				if dx < 0 {
					mobj.SetPos(obj.OriginalRight()+mobj.Length()/2, my)
				} else if dx > 0 {
					mobj.SetPos(obj.OriginalLeft()-mobj.Length()/2, my)
				}
				if dy < 0 {
					mobj.SetPos(mx, obj.OriginalTop()+mobj.Length()/2)
				} else if dy > 0 {
					mobj.SetPos(mx, obj.OriginalBottom()-mobj.Length()/2)
				}
			}
			return true
		}
	}
	return false
}

// onMovalbeObjCollisionObj 碰撞檢測之後確認碰撞結果
func onMovableObjCollisionObj(mobj object.IMovableObject, obj object.IObject) bool {
	var (
		mobjSubtype = mobj.Subtype()
		collision   bool
	)
	if mobjSubtype == object.ObjSubtypeShell {
		switch obj.Type() {
		case object.ObjTypeStatic:
			switch obj.Subtype() {
			case object.ObjSubtypeBrick, object.ObjSubtypeIron, object.ObjSubtypeHome:
				collision = true
			default:
			}
		case object.ObjTypeMovable:
			switch obj.Subtype() {
			case object.ObjSubtypeTank:
				if mobj.Camp() != obj.Camp() {
					collision = true
				}
			case object.ObjSubtypeShell:
				if mobj.Camp() != obj.Camp() {
					collision = true
				}
			default:
			}
		}
	} else if mobjSubtype == object.ObjSubtypeTank {
		switch obj.Type() {
		case object.ObjTypeStatic:
			switch obj.Subtype() {
			case object.ObjSubtypeBrick, object.ObjSubtypeIron, object.ObjSubtypeWater, object.ObjSubtypeIce, object.ObjSubtypeHome:
				collision = true
			default:
			}
		case object.ObjTypeMovable:
			switch obj.Subtype() {
			case object.ObjSubtypeShell:
				if mobj.Camp() != obj.Camp() {
					collision = true
				}
			case object.ObjSubtypeTank:
				collision = true
			default:
			}
		}
	}
	return collision
}

func isCollisionCorrectPosition(mobj object.IMovableObject, obj object.IObject) bool {
	if mobj.Subtype() == object.ObjSubtypeTank {
		if obj.Type() == object.ObjTypeStatic {
			switch obj.Subtype() {
			case object.ObjSubtypeBrick, object.ObjSubtypeIron, object.ObjSubtypeWater, object.ObjSubtypeIce, object.ObjSubtypeHome:
				return true
			}
		} else if obj.Type() == object.ObjTypeMovable {
			switch obj.Subtype() {
			case object.ObjSubtypeTank:
				return true
			}
		}
	}
	return false
}
