package common

import (
	"project_b/common/log"
	"project_b/common/object"
)

// checkMovableObjCollision 遍歷碰撞範圍内的網格檢查碰撞結果 移動之前調用
func checkMovableObjCollision(pmap *PartitionMap, obj object.IMovableObject, dir object.Direction, distance float64, collisionObj *object.IObject) bool {
	// 是否擁有碰撞組件
	comp := obj.GetComp("Collider")
	if comp == nil {
		return false
	}

	// 獲取檢測碰撞範圍
	lx, by, rx, ty := pmap.objGridBounds(obj)
	if rx < lx || ty < by {
		return false
	}

	for y := by; y <= ty; y++ {
		for x := lx; x <= rx; x++ {
			gidx := pmap.gridLineCol2Index(y, x)
			lis := pmap.grids[gidx].getMObjs().GetList()
			for i := 0; i < len(lis); i++ {
				item := lis[i]
				obj2, o := pmap.mobjs.Get(item.Key)
				if !o {
					log.Warn("Collision: grid(x:%v y:%v) not found movable object %v", x, y, item.Key)
					continue
				}
				if obj2.InstId() != obj.InstId() && obj2.StaticInfo().Layer() == obj.StaticInfo().Layer() {
					if checkMovableObjCollisionObj(obj, comp, dir, distance, obj2) {
						if collisionObj != nil {
							*collisionObj = obj2
						}
						return true
					}
				}
			}

			lis = pmap.grids[gidx].getSObjs().GetList()
			for i := 0; i < len(lis); i++ {
				item := lis[i]
				obj2, o := pmap.sobjs.Get(item.Key)
				if !o {
					log.Warn("Collision: grid(x:%v y:%v) not found static object %v", x, y, item.Key)
					continue
				}
				if checkMovableObjCollisionObj(obj, comp, dir, distance, obj2) {
					if collisionObj != nil {
						*collisionObj = obj2
					}
					return true
				}
			}
		}
	}
	return false
}

// checkMovableObjCollisionObj 檢查可移動物體和物體是否碰撞
func checkMovableObjCollisionObj(mobj object.IMovableObject, comp object.IComponent, dir object.Direction, distance float64, obj object.IObject) bool {
	if !(mobj.StaticInfo().Layer() == obj.StaticInfo().Layer() ||
		(obj.Type() == object.ObjTypeStatic && (obj.Subtype() == object.ObjSubTypeWater || obj.Subtype() == object.ObjSubTypeIce))) {
		return false
	}

	var (
		collisionComp *object.ColliderComp
		aabb1         object.AABB
	)
	collisionComp = comp.(*object.ColliderComp)
	aabb1 = collisionComp.GetAABB(mobj)
	aabb1.Move(dir, distance)
	comp2 := obj.GetComp("Collider")
	if comp2 == nil {
		return false
	}
	collisionComp2 := comp2.(*object.ColliderComp)
	if collisionComp2 == nil {
		return false
	}
	aabb2 := collisionComp2.GetAABB(obj)
	if aabb1.MoveIntersect(dir, &aabb2) {
		if onMovableObjCollisionObj(mobj, obj) {
			switch dir {
			case object.DirLeft:
				mobj.SetPos(obj.Right(), mobj.Bottom())
			case object.DirRight:
				mobj.SetPos(obj.Left()-mobj.Width(), mobj.Bottom())
			case object.DirUp:
				mobj.SetPos(mobj.Left(), obj.Bottom()-mobj.Height())
			case object.DirDown:
				mobj.SetPos(mobj.Left(), obj.Top())
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
	if mobjSubtype == object.ObjSubTypeBullet {
		switch obj.Type() {
		case object.ObjTypeStatic:
			switch obj.Subtype() {
			case object.ObjSubTypeBrick, object.ObjSubTypeIron, object.ObjSubTypeHome:
				collision = true
			default:
			}
		case object.ObjTypeMovable:
			switch obj.Subtype() {
			case object.ObjSubTypeTank:
				if mobj.Camp() != obj.Camp() {
					collision = true
				}
			case object.ObjSubTypeBullet:
				if mobj.Camp() != obj.Camp() {
					collision = true
				}
			default:
			}
		}
	} else if mobjSubtype == object.ObjSubTypeTank {
		switch obj.Type() {
		case object.ObjTypeStatic:
			switch obj.Subtype() {
			case object.ObjSubTypeBrick, object.ObjSubTypeIron, object.ObjSubTypeWater, object.ObjSubTypeIce, object.ObjSubTypeHome:
				collision = true
			default:
			}
		case object.ObjTypeMovable:
			switch obj.Subtype() {
			case object.ObjSubTypeBullet:
				if mobj.Camp() != obj.Camp() {
					collision = true
				}
			case object.ObjSubTypeTank:
				collision = true
			default:
			}
		}
	}
	return collision
}
