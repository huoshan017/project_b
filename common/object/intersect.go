package object

import (
	"math"
	"project_b/common/base"
)

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

	//log.Debug("(AB X AP) * (CD X CP) = (%v*%v) = %v", c1, c2, c1*c2)
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

	//log.Debug("(BC X BP) * (DA X DP) = (%v*%v) = %v", c3, c4, c3*c4)
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
		if (s1.Y-y)*(e1.Y-y) > 0 {
			return false
		}
	} else if e1.Y == s1.Y {
		if e2.Y == s2.Y {
			return false
		}
		// (e2.X-s2.X)/(e2.Y-s2.Y) == (x-s2.X)(y-s2.Y)
		y = e1.Y
		x = (e2.X-s2.X)*(e1.Y-s2.Y)/(e2.Y-s2.Y) + s2.X
		if (s1.X-x)*(e1.X-x) > 0 {
			return false
		}
	} else if e2.X == s2.X {
		if e1.X == s1.X {
			return false
		}
		// (e1.X-s1.X)/(e1.Y-s1.Y) == (x-s1.X)/(y-s1.Y)
		x = e2.X
		y = (e1.Y-s1.Y)*(e2.X-s1.X)/(e1.X-s1.X) + s1.Y
		if (s2.Y-y)*(e2.Y-y) > 0 {
			return false
		}
	} else if e2.Y == s2.Y {
		if e1.Y == s1.Y {
			return false
		}
		// (e1.X-s1.X)/(e1.Y-s1.Y) == (x-s1.X)/(y-s1.Y)
		y = e2.Y
		x = (e1.X-s1.X)*(e2.Y-s1.Y)/(e1.Y-s1.Y) + s1.X
		if (s2.X-x)*(e2.X-x) > 0 {
			return false
		}
	} else {
		// 利用斜率求出交點的坐標
		// (e1.X-s1.X)/(e1.Y-s1.Y) = (x-s1.X)/(y-s1.Y)
		// (e2.X-s2.X)/(e2.Y-s2.Y) = (x-s2.X)/(y-s2.Y)
		var fm int64 = int64(s1.X)*int64(e1.Y-s1.Y)*int64(e2.X-s2.X) - int64(s2.X)*int64(e2.Y-s2.Y)*int64(e1.X-s1.X) + int64(s2.Y-s1.Y)*int64(e1.X-s1.X)*int64(e2.X-s2.X)
		var fz int64 = int64(e1.Y-s1.Y)*int64(e2.X-s2.X) - int64(e1.X-s1.X)*int64(e2.Y-s2.Y)
		x = int32(fm / fz)
		y = s1.Y + int32(int64(e1.Y-s1.Y)*int64(x-s1.X)/int64(e1.X-s1.X))
		//log.Debug("get intersection (%v, %v)", x, y)
	}

	// 再判斷交點是否在兩條綫段上（因爲上面求出的交點有可能只是在綫段延長綫上）
	// 交點與綫段兩端的差值符號相反
	var a, b int64
	if s1.X != e1.X {
		a = int64(s1.X-x) * int64(e1.X-x)
	} else {
		a = int64(s1.Y-y) * int64(e1.Y-y)
	}
	if s2.X != e2.X {
		b = int64(s2.X-x) * int64(e2.X-x)
	} else {
		b = int64(s2.Y-y) * int64(e2.Y-y)
	}
	if a > 0 || b > 0 {
		//log.Debug("point(%v, %v) not within line segment (%+v, %+v) and (%+v, %+v)", x, y, s1, e1, s2, e2)
		return false
	}

	if in != nil {
		in.X, in.Y = x, y
	}

	//log.Debug("two line segment (%+v, %+v) and (%+v, %+v) with intersection point (%v, %v)", start1, end1, start2, end2, x, y)

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

// 獲得綫段和物體的交點(是距離綫段起始點start最近的交點，遠些的交點沒有太大用處)
// 綫段分別與四條邊判斷是否有交點，如果有交點在綫段上，這個點與綫段兩個端點構成的兩條綫段的斜率相等
func GetLineSegmentAndObjIntersection(start, end *base.Pos, obj IObject, intersection *base.Pos) bool {
	posA := base.NewPos(obj.LeftBottom())
	posB := base.NewPos(obj.RightBottom())
	posC := base.NewPos(obj.RightTop())
	posD := base.NewPos(obj.LeftTop())

	var (
		posList = []*base.Pos{&posA, &posB, &posC, &posD}
		n       int8
		in      [2]base.Pos
		inn     *base.Pos
	)

	//log.Debug("obj %v A(%v) B(%v) C(%v) D(%v)", obj.InstId(), posA, posB, posC, posD)

	for i := 0; i < len(posList) && n < 2; i++ {
		if intersection != nil {
			inn = &in[n]
		} else {
			inn = nil
		}
		if GetTwoLineSegmentIntersection(start, end, posList[i], posList[(i+1)%len(posList)], inn) {
			n += 1
		}
	}

	var result = (n >= 1)
	if result {
		var d uint32 = math.MaxUint32
		for i := int8(0); i < n; i++ {
			dd := base.Distance(start, &in[i])
			if d > dd {
				*intersection = in[i]
				d = dd
			}
		}
	}
	return result
}
