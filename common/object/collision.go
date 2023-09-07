package object

import (
	"fmt"
	"math"

	"project_b/common/base"
	"project_b/log"
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
	MovingObjPos     base.Pos        // 相切時移動物體的位置
	TangentPoints    [2]base.Pos     // 相切時碰撞點列表，最多兩個
	TangentPointsNum int8            // 相切時碰撞點數量[0,1,2]，兩個點表示物體相切的綫段兩端點
}

func (ci *CollisionInfo) Clear() {
	ci.Result = CollisionNone
	if len(ci.ObjList) > 0 {
		clear(ci.ObjList)
		ci.ObjList = ci.ObjList[:0]
	}
	ci.MovingObj = nil
	ci.MovingObjPos.X = 0
	ci.MovingObjPos.Y = 0
	ci.TangentPointsNum = 0
}

// 點是否在物體内
func IsPointInObject(pos base.Pos, obj IObject) bool {
	posA := base.NewPos(obj.LeftBottom())
	posB := base.NewPos(obj.RightBottom())
	posC := base.NewPos(obj.RightTop())
	posD := base.NewPos(obj.LeftTop())

	// AB X AP
	ab := base.NewVec2(posB.X-posA.X, posB.Y-posA.Y)
	ap := base.NewVec2(pos.X-posA.X, pos.Y-posA.Y)
	c1 := ab.Cross(ap)
	// CD X CP
	cd := base.NewVec2(posD.X-posC.X, posD.Y-posC.Y)
	cp := base.NewVec2(pos.X-posC.X, pos.Y-posC.Y)
	c2 := cd.Cross(cp)

	log.Debug("(AB X AP) * (CD X CP) = (%v*%v) = %v", c1, c2, c1*c2)
	if int64(c1)*int64(c2) < 0 {
		return false
	}

	// BC X BP
	bc := base.NewVec2(posC.X-posB.X, posC.Y-posB.Y)
	bp := base.NewVec2(pos.X-posB.X, pos.Y-posB.Y)
	c3 := bc.Cross(bp)
	// DA X DP
	da := base.NewVec2(posA.X-posD.X, posA.Y-posD.Y)
	dp := base.NewVec2(pos.X-posD.X, pos.Y-posD.Y)
	c4 := da.Cross(dp)

	log.Debug("(BC X BP) * (DA X DP) = (%v*%v) = %v", c3, c4, c3*c4)
	return int64(c3)*int64(c4) >= 0
}

// 檢測兩條綫段是否相交
func CheckTwoLineSegmentIntersect(start1, end1, start2, end2 *base.Pos) bool {
	// 綫段(start1, end1)
	se := base.NewVec2(end1.X-start1.X, end1.Y-start1.Y)
	// 綫段(start1, start2)
	sa := base.NewVec2(start2.X-start1.X, start2.Y-start1.Y)
	// 綫段(start1, end2)
	sb := base.NewVec2(end2.X-start1.X, end2.Y-start1.Y)

	c1 := se.Cross(sa)
	c2 := se.Cross(sb)

	// 綫段(start2, end2)
	ab := base.NewVec2(end2.X-start2.X, end2.Y-start2.Y)
	// 綫段(start2, start1)
	as := base.NewVec2(start1.X-start2.X, start1.Y-start2.Y)
	// 綫段(start2, end1)
	bs := base.NewVec2(end1.X-start2.X, end1.Y-start2.Y)

	c3 := ab.Cross(as)
	c4 := ab.Cross(bs)

	// c1*c2 和 c3*c4 都小於0才能判斷相交
	// TODO c1*c2=0 或者 c3*c4=0 屬於相切，不是嚴格意義上的相交
	if int64(c1)*int64(c2) < 0 && int64(c3)*int64(c4) < 0 {
		return true
	}

	return false
}

// 兩條綫段交點
func GetTwoLineSegmentIntersection(start1, end1, start2, end2 *base.Pos, intersection *base.Pos) bool {
	s1, e1, s2, e2, in := start1, end1, start2, end2, intersection
	// 平行沒有交點
	// (e1.X-s1.X)/(e1.Y-s1.Y) == (e2.X-s2.X)/(e2.Y-s2.Y)
	if (e1.X-s1.X)*(e2.Y-s2.Y) == (e1.Y-s1.Y)*(e2.X-s2.X) {
		return false
	}

	var x, y int32
	// 其中一條水平或者垂直
	if e1.X == s1.X {
		if e2.X == s2.X {
			return false
		}
		// (e2.X-s2.X)/(e2.Y-s2.Y) == (x-s2.X)/(y-s2.Y)
		x = e1.X
		y = (e1.X-s2.X)*(e2.Y-s2.Y)/(e2.X-s2.X) + s2.Y
	} else if e1.Y == s1.Y {
		if e2.Y == s2.Y {
			return false
		}
		// (e2.X-s2.X)/(e2.Y-s2.Y) == (x-s2.X)(y-s2.Y)
		y = e1.Y
		x = (e2.X-s2.X)*(e1.Y-s2.Y)/(e2.Y-s2.Y) + s2.X
	} else if e2.X == s2.X {
		if e1.X == s1.X {
			return false
		}
		// (e1.X-s1.X)/(e1.Y-s1.Y) == (x-s1.X)/(y-s1.Y)
		x = e2.X
		y = (e1.Y-s1.Y)*(e2.X-s1.X)/(e1.X-s1.X) + s1.Y
	} else if e2.Y == s2.Y {
		if e1.Y == s1.Y {
			return false
		}
		// (e1.X-s1.X)/(e1.Y-s1.Y) == (x-s1.X)/(y-s1.Y)
		y = e2.Y
		x = (e1.X-s1.X)*(e2.Y-s1.Y)/(e1.Y-s1.Y) + s1.X
	} else {
		// 利用斜率求出交點的坐標
		// (e1.X-s1.X)/(e1.Y-s1.Y) = (in.X-s1.X)/(in.Y-s1.Y)
		// (e2.X-s2.X)/(e2.Y-s2.Y) = (in.X-s2.X)/(in.Y-s2.Y)
		// in.X = (e1.X-s1.X)*(in.Y-s1.Y)/(e1.Y-s1.Y) + s1.X
		// in.Y = (e2.Y-s2.Y)*(in.X-s2.X)/(e2.X-s2.X) + s2.Y
		// in.X = (e1.X-s1.X)*((e2.Y-s2.Y)*(in.X-s2.X)/(e2.X-s2.X) + s2.Y-s1.Y)/(e1.Y-s1.Y) + s1.X
		// in.X*(e1.Y-s1.Y) = (e1.X-s1.X)*((e2.Y-s2.Y)*(in.X-s2.X)/(e2.X-s2.X)+s2.Y-s1.Y) +s1.X*(e1.Y-s1.Y)
		// in.X*[(e1.Y-s1.Y)-(e1.X-s1.X)*(e2.Y-s2.Y)/(e2.X-s2.X)] =
		//     (e1.X-s1.X)*((e2.Y-s2.Y)*(-s2.X)/(e2.X-s2.X)+s2.Y-s1.Y) + s1.X*(e1.Y-s1.Y)
		//in.X = ((e1.X-s1.X)*((e2.Y-s2.Y)*(-s2.X)/(e2.X-s2.X)+s2.Y-s1.Y) + s1.X*(e1.Y-s1.Y)) / ((e1.Y - s1.Y) - (e1.X-s1.X)*(e2.Y-s2.Y)/(e2.X-s2.X))
		x = ((e1.X-s1.X)*((e2.Y-s2.Y)*(-s2.X)+(s2.Y-s1.Y)*(e2.X-s2.X)) + s1.X*(e1.Y-s1.Y)*(e2.X-s2.X)) / ((e1.Y-s1.Y)*(e2.X-s2.X) - (e1.X-s1.X)*(e2.Y-s2.Y))
		y = (e2.Y-s2.Y)*(x-s2.X)/(e2.X-s2.X) + s2.Y
		log.Debug("get intersection (%v, %v)", x, y)
	}

	// 再判斷交點是否在兩條綫段上（因爲上面求出的交點有可能只是在綫段延長綫上）
	// 交點與綫段兩端的差值符號相反
	if start1.X != end1.X && (start1.X-x)*(end1.X-x) > 0 {
		log.Debug("point(%v, %v) not within line segment (%+v, %+v) and (%+v, %+v)", x, y, start1, end1, start2, end2)
		return false
	}
	if start2.X != end2.X && (start2.X-x)*(end2.X-x) > 0 {
		log.Debug("point(%v, %v) not within line segment (%+v, %+v) and (%+v, %+v)", x, y, start1, end1, start2, end2)
		return false
	}

	in.X, in.Y = x, y

	log.Debug("two line segment (%+v, %+v) and (%+v, %+v) with intersection point (%v, %v)", start1, end1, start2, end2, x, y)

	return true
}

// 綫段跟物體是否相交
func CheckLineSegmentIntersectObj(start, end *base.Pos, obj IObject) bool {
	posA := base.NewPos(obj.LeftBottom())
	posB := base.NewPos(obj.RightBottom())
	posC := base.NewPos(obj.RightTop())
	posD := base.NewPos(obj.LeftTop())

	var (
		posList = []*base.Pos{&posA, &posB, &posC, &posD}
		result  bool
	)

	for i := 0; i < len(posList); i++ {
		if CheckTwoLineSegmentIntersect(start, end, posList[i], posList[(i+1)%len(posList)]) {
			result = true
			break
		}
	}

	return result
}

// 獲得綫段和物體的交點
func GetLineSegmentAndObjIntersection(start, end *base.Pos, obj IObject, intersection *base.Pos, intersection2 *base.Pos) bool {
	// 綫段分別與四條邊判斷是否有交點，如果有交點在綫段上，這個點與綫段兩個端點構成的兩條綫段的斜率相等
	posA := base.NewPos(obj.LeftBottom())
	posB := base.NewPos(obj.RightBottom())
	posC := base.NewPos(obj.RightTop())
	posD := base.NewPos(obj.LeftTop())

	var (
		posList = []*base.Pos{&posA, &posB, &posC, &posD}
		n       int8
		in      *base.Pos
		result  bool
	)
	for i := 0; i < len(posList); i++ {
		if n == 0 {
			in = intersection
		} else if n == 1 {
			in = intersection2
		} else {
			result = true
			break
		}
		if GetTwoLineSegmentIntersection(start, end, posList[i], posList[(i+1)%len(posList)], in) {
			n += 1
		}
	}
	return result
}

// CheckMovingObjCollisionObj 檢查可移動物體和物體是否碰撞
func CheckMovingObjCollisionObj(mobj IMovableObject, dx, dy int32, obj IObject) CollisionResult {
	var cr CollisionResult
	if !(mobj.StaticInfo().Layer() == obj.StaticInfo().Layer() ||
		(obj.Type() == base.ObjTypeStatic && (obj.Subtype() == base.ObjSubtypeWater || obj.Subtype() == base.ObjSubtypeIce))) {
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
					// y方向上的移動距離與x方向成正比
					y = x * dy / dx
				} else if aabb.Right <= aabb2.Left && aabb.Top <= aabb2.Bottom {
					x = aabb2.Left - aabb.Right
					y = aabb2.Bottom - aabb.Top
				} else if aabb.Right > aabb2.Left && aabb.Top <= aabb2.Bottom {
					y = aabb2.Bottom - aabb.Top
					x = y * dx / dy
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
					y = x * -dy / dx
				} else if aabb.Right <= aabb2.Left && aabb.Bottom >= aabb2.Top {
					x = aabb2.Left - aabb.Right
					y = aabb.Bottom - aabb2.Top
				} else if aabb.Right > aabb2.Left && aabb.Bottom >= aabb2.Top {
					y = aabb.Bottom - aabb2.Top
					x = y * dx / -dy
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
					y = x * dy / -dx
				} else if aabb.Left >= aabb2.Right && aabb.Top <= aabb2.Bottom {
					x = aabb.Left - aabb2.Right
					y = aabb2.Bottom - aabb.Top
				} else if aabb.Left < aabb2.Right && aabb.Top <= aabb2.Bottom {
					y = aabb2.Bottom - aabb.Top
					x = y * -dx / dy
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
				if aabb.Left >= aabb2.Right && aabb.Bottom < aabb2.Top {
					x = aabb.Left - aabb2.Right
					y = x * -dy / -dx
				} else if aabb.Left >= aabb2.Right && aabb.Bottom >= aabb2.Top {
					x = aabb.Left - aabb2.Right
					y = aabb.Bottom - aabb2.Top
				} else if aabb.Left < aabb2.Right && aabb.Bottom >= aabb2.Top {
					y = aabb.Bottom - aabb2.Top
					x = y * -dx / -dy
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
	return collisionInfo.Result
}

// movingObjCollision2Obj 移動物體是否被阻擋
func movingObjCollision2ObjResult(mobj IMovableObject, obj IObject) CollisionResult {
	var (
		mobjSubtype = mobj.Subtype()
		result      CollisionResult
	)
	if mobjSubtype == base.ObjSubtypeShell {
		switch obj.Type() {
		case base.ObjTypeStatic:
			switch obj.Subtype() {
			case base.ObjSubtypeBrick, base.ObjSubtypeIron, base.ObjSubtypeHome:
				result = CollisionOnly
			default:
			}
		case base.ObjTypeMovable:
			switch obj.Subtype() {
			case base.ObjSubtypeTank:
				if mobj.Camp() != obj.Camp() {
					result = CollisionOnly
				}
			case base.ObjSubtypeShell:
				if mobj.Camp() != obj.Camp() {
					result = CollisionOnly
				}
			default:
			}
		}
	} else if mobjSubtype == base.ObjSubtypeTank {
		switch obj.Type() {
		case base.ObjTypeStatic:
			switch obj.Subtype() {
			case base.ObjSubtypeBrick, base.ObjSubtypeIron, base.ObjSubtypeWater, base.ObjSubtypeIce, base.ObjSubtypeHome:
				result = CollisionAndBlock
			default:
			}
		case base.ObjTypeMovable:
			switch obj.Subtype() {
			case base.ObjSubtypeShell:
				if mobj.Camp() != obj.Camp() {
					result = CollisionOnly
				}
			case base.ObjSubtypeTank:
				result = CollisionAndBlock
			default:
			}
		case base.ObjTypeItem:
			result = CollisionOnly
		}
	}
	return result
}
