package object

import (
	"fmt"
	"math"
)

type CollisionResult int32

const (
	CollisionNone     CollisionResult = iota // 無碰撞
	CollisionOnly                            // 僅僅碰撞
	CollisionAndBlock                        // 碰撞且阻擋
)

type CollisionInfo struct {
	Result           CollisionResult // 碰撞結果
	ObjList          []IObject       // 碰撞物體的列表
	MovingObj        IMovableObject  // 移動物體
	MovingObjPos     Pos             // 相切時移動物體的位置
	TangentPoints    [2]Pos          // 相切時碰撞點列表，最多兩個
	TangentPointsNum int8            // 相切時碰撞點數量[0,1,2]，兩個點表示物體相切的綫段兩端點
}

// CheckMovingObjCollisionObj 檢查可移動物體和物體是否碰撞
func CheckMovingObjCollisionObj(mobj IMovableObject, dx, dy int32, obj IObject) CollisionResult {
	var cr CollisionResult
	if !(mobj.StaticInfo().Layer() == obj.StaticInfo().Layer() ||
		(obj.Type() == ObjTypeStatic && (obj.Subtype() == ObjSubtypeWater || obj.Subtype() == ObjSubtypeIce))) {
		return cr
	}
	collisionComp1 := mobj.GetColliderComp()
	if collisionComp1 == nil {
		return cr
	}
	collisionComp2 := obj.GetColliderComp()
	if collisionComp2 == nil {
		return cr
	}
	aabb1 := collisionComp1.GetAABB()
	aabb2 := collisionComp2.GetAABB()
	if aabb1.MoveIntersect(dx, dy, &aabb2) {
		cr = movingObjCollision2ObjResult(mobj, obj)
	}
	return cr
}

// 進入到narrow phase階段的碰撞物體列表都是會有碰撞的，所以這裏要做的是找出最先碰撞的，并且計算碰撞點和此時mobj的位置
func NarrowPhaseCheckMovingObjCollision2ObjList(mobj IMovableObject, dx, dy int32, collisionObjList []IObject, collisionInfo *CollisionInfo) CollisionResult {
	var (
		obj IObject
		ml  = int32(len(collisionObjList))
	)

	aabb := mobj.GetColliderComp().GetAABB()
	collisionInfo.Result = CollisionAndBlock
	collisionInfo.MovingObj = mobj
	if dx == 0 || dy == 0 {
		var (
			md      = int32(math.MaxInt32)
			idxList []int32
		)
		fn := func(d, md, i int32) int32 {
			if d < md {
				md = d
				if len(idxList) > 0 {
					idxList = idxList[:1]
					idxList[0] = i
				} else {
					idxList = append(idxList, i)
				}
			} else if d == md {
				idxList = append(idxList, i)
			}
			return md
		}
		px, py := mobj.Pos()
		if dy > 0 { // 正上方移動
			for i := int32(0); i < ml; i++ {
				obj = collisionObjList[i]
				aabb2 := obj.GetColliderComp().GetAABB()
				d := aabb2.Bottom - aabb.Top
				md = fn(d, md, i)
			}
			collisionInfo.MovingObjPos.X = px
			collisionInfo.MovingObjPos.Y = py + md
		} else if dy < 0 { // 正下方移動
			for i := int32(0); i < ml; i++ {
				obj = collisionObjList[i]
				aabb2 := obj.GetColliderComp().GetAABB()
				d := aabb.Bottom - aabb2.Top
				md = fn(d, md, i)
			}
			collisionInfo.MovingObjPos.X = px
			collisionInfo.MovingObjPos.Y = py - md
		}
		if dx > 0 { // 正右方移動
			for i := int32(0); i < ml; i++ {
				obj = collisionObjList[i]
				aabb2 := obj.GetColliderComp().GetAABB()
				d := aabb2.Left - aabb.Right
				md = fn(d, md, i)
			}
			collisionInfo.MovingObjPos.X = px + md
			collisionInfo.MovingObjPos.Y = py
		} else if dx < 0 { // 正左方移動
			for i := int32(0); i < ml; i++ {
				obj = collisionObjList[i]
				aabb2 := obj.GetColliderComp().GetAABB()
				d := aabb.Left - aabb2.Right
				md = fn(d, md, i)
			}
			collisionInfo.MovingObjPos.X = px - md
			collisionInfo.MovingObjPos.Y = py
		}
		for i := 0; i < len(idxList); i++ {
			collisionInfo.ObjList = append(collisionInfo.ObjList, collisionObjList[idxList[i]])
		}
		/*if md < 0 {
			if mobj.Type() == ObjTypeMovable && mobj.Subtype() == ObjSubtypeTank {
				var str string
				for i := 0; i < len(collisionObjList); i++ {
					str += fmt.Sprintf("%v ", collisionObjList[i].InstId())
				}
				panic(fmt.Sprintf("dx(%v) dy(%v)  md(%v)  collisionObjList(%s), mobj(%+v)", dx, dy, md, str, *mobj.(*Tank)))
			}
		}*/
	} else {
		var (
			mx, my       = int32(math.MaxInt32), int32(math.MaxInt32)
			xList, yList []int32
		)
		fn := func(x, y, i int32) (int32, int32) {
			if x < mx {
				mx = x
				if len(xList) > 0 {
					xList = xList[:0]
				}
				xList = append(xList, i)
			} else if x == mx {
				xList = append(xList, i)
			}
			if y < my {
				my = y
				if len(yList) > 0 {
					yList = yList[:0]
				}
				yList = append(yList, i)
			} else if y == my {
				yList = append(yList, i)
			}
			return mx, my
		}
		if dx > 0 && dy > 0 { // 右上方移動
			for i := int32(0); i < ml; i++ {
				obj = collisionObjList[i]
				aabb2 := obj.GetColliderComp().GetAABB()
				//                o o o o o o o o
				//                o             o
				//                o             o
				//                o             o
				//                o             o
				//                o             o
				// -------------- o o o o o o o o ------
				//                |
				//    o o o o o   |
				//    o       o   |
				//    o       o   |
				//    o o o o o   |
				//                |
				var x, y int32
				if aabb.Right <= aabb2.Left && aabb.Top > aabb2.Bottom {
					x = aabb2.Left - aabb.Right
				} else if aabb.Right <= aabb2.Left && aabb.Top <= aabb2.Bottom {
					x = aabb2.Left - aabb.Right
					y = aabb2.Bottom - aabb.Top
				} else if aabb.Right > aabb2.Left && aabb.Top <= aabb2.Bottom {
					y = aabb2.Bottom - aabb.Top
				} else {
					panic(fmt.Sprintf("dx(%v)>0 dy(%v)>0  aabb.Right(%v) > aabb2.Left(%v) && aabb.Top(%v) > aabb2.Bottom(%v)", dx, dy, aabb.Right, aabb2.Left, aabb.Top, aabb2.Bottom))
				}
				mx, my = fn(x, y, i)
				if movingObjCollision2ObjResult(mobj, obj) == CollisionAndBlock {
					collisionInfo.Result = CollisionAndBlock
				}
			}
		} else if dx > 0 && dy < 0 { // 右下方移動
			for i := int32(0); i < ml; i++ {
				obj = collisionObjList[i]
				aabb2 := obj.GetColliderComp().GetAABB()
				//                 |
				//     o o o o o   |
				//     o       o   |
				//     o       o   |
				//     o o o o o   |
				//                 |
				// --------------- o o o o o o o o ------
				//                 o             o
				//                 o             o
				//                 o             o
				//                 o             o
				//                 o o o o o o o o
				//                 |
				var x, y int32
				if aabb.Right <= aabb2.Left && aabb.Bottom < aabb2.Top {
					x = aabb2.Left - aabb.Right
				} else if aabb.Right <= aabb2.Left && aabb.Bottom >= aabb2.Top {
					x = aabb2.Left - aabb.Right
					y = aabb.Bottom - aabb2.Top
				} else if aabb.Right > aabb2.Left && aabb.Bottom >= aabb2.Top {
					y = aabb.Bottom - aabb2.Top
				} else {
					panic(fmt.Sprintf("dx(%v)>0 dy(%v)<0  aabb.Right(%v) > aabb2.Left(%v) && aabb.Bottom(%v) < aabb2.Top(%v)", dx, dy, aabb.Right, aabb2.Left, aabb.Bottom, aabb2.Top))
				}
				mx, my = fn(x, y, i)
				if movingObjCollision2ObjResult(mobj, obj) == CollisionAndBlock {
					collisionInfo.Result = CollisionAndBlock
				}
			}
		} else if dx < 0 && dy > 0 { // 左上方移動
			for i := int32(0); i < ml; i++ {
				obj = collisionObjList[i]
				aabb2 := obj.GetColliderComp().GetAABB()
				//   + + + + + + + +
				//   +             +
				//   +             +
				//   +             +
				//   +             +
				//   +             +
				// - + + + + + + + + --------------
				//                 |   + + + + +
				//                 |   +       +
				//                 |   +       +
				//                 |   + + + + +
				//                 |
				var x, y int32
				if aabb.Left >= aabb2.Right && aabb.Top > aabb2.Bottom {
					x = aabb.Left - aabb2.Right
				} else if aabb.Left >= aabb2.Right && aabb.Top <= aabb2.Bottom {
					x = aabb.Left - aabb2.Right
					y = aabb2.Bottom - aabb.Top
				} else if aabb.Left < aabb2.Right && aabb.Top <= aabb2.Bottom {
					y = aabb2.Bottom - aabb.Top
				} else {
					panic(fmt.Sprintf("dx(%v)<0 dy(%v)>0  aabb.Left(%v) < aabb2.Right(%v) && aabb.Top(%v) > aabb2.Bottom(%v)", dx, dy, aabb.Left, aabb2.Right, aabb.Top, aabb2.Bottom))
				}
				mx, my = fn(x, y, i)
				if movingObjCollision2ObjResult(mobj, obj) == CollisionAndBlock {
					collisionInfo.Result = CollisionAndBlock
				}
			}
		} else { // 左下方移動
			for i := int32(0); i < ml; i++ {
				obj = collisionObjList[i]
				aabb2 := obj.GetColliderComp().GetAABB()
				//                 |
				//                 |   + + + + +
				//                 |   +       +
				//                 |   +       +
				//                 |   + + + + +
				// - + + + + + + + + ----------------
				//   +             +
				//   +             +
				//   +             +
				//   +             +
				//   +             +
				//   + + + + + + + +
				var x, y int32
				if aabb.Left >= aabb2.Right && aabb.Top < aabb2.Bottom {
					x = aabb.Left - aabb2.Right
				} else if aabb.Left >= aabb2.Right && aabb.Top >= aabb2.Bottom {
					x = aabb.Left - aabb2.Right
					y = aabb.Top - aabb2.Bottom
				} else if aabb.Left < aabb2.Right && aabb.Top >= aabb2.Bottom {
					y = aabb.Top - aabb2.Bottom
				} else {
					panic(fmt.Sprintf("dx(%v)<0 dy(%v)<0  aabb.Left(%v) < aabb2.Right(%v) && aabb.Top(%v) >= aabb2.Bottom(%v)", dx, dy, aabb.Left, aabb2.Right, aabb.Top, aabb2.Bottom))
				}
				mx, my = fn(x, y, i)
				if movingObjCollision2ObjResult(mobj, obj) == CollisionAndBlock {
					collisionInfo.Result = CollisionAndBlock
				}
			}
		}
		px, py := mobj.Pos()
		if mx == math.MaxInt32 && my != math.MaxInt32 {
			for i := 0; i < len(yList); i++ {
				collisionInfo.ObjList = append(collisionInfo.ObjList, collisionObjList[yList[i]])
			}
		} else if mx != math.MaxInt32 && my == math.MaxInt32 {
			for i := 0; i < len(xList); i++ {
				collisionInfo.ObjList = append(collisionInfo.ObjList, collisionObjList[xList[i]])
			}
		} else {
			ddx, ddy := dx, dy
			if ddx < 0 {
				ddx = -ddx
			}
			if ddy < 0 {
				ddy = -ddy
			}
			if ddx*my > mx*ddy { // ddx / ddy > mx / my
				// 説明y方向移動過大，調整my適應mx
				my = ddy * mx / ddx
				for i := 0; i < len(xList); i++ {
					collisionInfo.ObjList = append(collisionInfo.ObjList, collisionObjList[xList[i]])
				}
			} else if ddx*my < mx*ddy { // ddx / ddy < mx / my
				// 説明x方向移動過大，調整mx適應my
				mx = ddx * my / ddy
				for i := 0; i < len(yList); i++ {
					collisionInfo.ObjList = append(collisionInfo.ObjList, collisionObjList[yList[i]])
				}
			} else {
				// 不用調整
				for i := 0; i < len(xList); i++ {
					collisionInfo.ObjList = append(collisionInfo.ObjList, collisionObjList[xList[i]])
				}
				for i := 0; i < len(yList); i++ {
					collisionInfo.ObjList = append(collisionInfo.ObjList, collisionObjList[yList[i]])
				}
			}
		}
		if dx < 0 {
			mx = -mx
		}
		if dy < 0 {
			my = -my
		}
		collisionInfo.MovingObjPos.X = px + mx
		collisionInfo.MovingObjPos.Y = py + my
	}
	// 測試用
	/*{
		var _aabb AABB
		if mobj.Subtype() == ObjSubtypeTank {
			var tank Tank = *(mobj.(*Tank))
			tank.SetPos(collisionInfo.MovingObjPos.X, collisionInfo.MovingObjPos.Y)
			_aabb = tank.colliderComp.GetAABB()
		} else if mobj.Subtype() == ObjSubtypeShell {
			var shell Shell = *(mobj.(*Shell))
			shell.SetPos(collisionInfo.MovingObjPos.X, collisionInfo.MovingObjPos.Y)
			_aabb = shell.colliderComp.GetAABB()
		} else {
			panic(fmt.Sprintf("mobj 子類型(%v)錯誤", mobj.Subtype()))
		}
		for i := 0; i < len(collisionObjList); i++ {
			obj := collisionObjList[i]
			_aabb2 := obj.GetColliderComp().GetAABB()
			if _aabb.Intersect(&_aabb2) {
				panic(fmt.Sprintf("!!! 測試失敗 aabb %v intersect aabb %v, when dx(%v) dy(%v)", _aabb, _aabb2, dx, dy))
			}
		}
	}*/
	return collisionInfo.Result
}

// movingObjCollision2Obj 移動物體是否被阻擋
func movingObjCollision2ObjResult(mobj IMovableObject, obj IObject) CollisionResult {
	var (
		mobjSubtype = mobj.Subtype()
		result      CollisionResult
	)
	if mobjSubtype == ObjSubtypeShell {
		switch obj.Type() {
		case ObjTypeStatic:
			switch obj.Subtype() {
			case ObjSubtypeBrick, ObjSubtypeIron, ObjSubtypeHome:
				result = CollisionOnly
			default:
			}
		case ObjTypeMovable:
			switch obj.Subtype() {
			case ObjSubtypeTank:
				if mobj.Camp() != obj.Camp() {
					result = CollisionOnly
				}
			case ObjSubtypeShell:
				if mobj.Camp() != obj.Camp() {
					result = CollisionOnly
				}
			default:
			}
		}
	} else if mobjSubtype == ObjSubtypeTank {
		switch obj.Type() {
		case ObjTypeStatic:
			switch obj.Subtype() {
			case ObjSubtypeBrick, ObjSubtypeIron, ObjSubtypeWater, ObjSubtypeIce, ObjSubtypeHome:
				result = CollisionAndBlock
			default:
			}
		case ObjTypeMovable:
			switch obj.Subtype() {
			case ObjSubtypeShell:
				if mobj.Camp() != obj.Camp() {
					result = CollisionOnly
				}
			case ObjSubtypeTank:
				result = CollisionAndBlock
			default:
			}
		case ObjTypeItem:
			result = CollisionOnly
		}
	}
	return result
}
