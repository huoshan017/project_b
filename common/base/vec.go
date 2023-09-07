package base

type Vec2 struct {
	x, y int32
}

func NewVec2(x, y int32) Vec2 {
	return Vec2{x: x, y: y}
}

func (v Vec2) X() int32 {
	return v.x
}

func (v Vec2) Y() int32 {
	return v.y
}

func (v Vec2) Add(v2 Vec2) Vec2 {
	return Vec2{v.x + v2.x, v.y + v2.y}
}

func (v Vec2) Sub(v2 Vec2) Vec2 {
	return Vec2{v.x - v2.x, v.y - v2.y}
}

func (v Vec2) Mul(a int32) Vec2 {
	return Vec2{v.x * a, v.y * a}
}

func (v Vec2) Div(a int32) Vec2 {
	return Vec2{v.x / a, v.y / a}
}

func (v Vec2) Length() int32 {
	return int32(Sqrt(uint32(v.x*v.x + v.y*v.y)))
}

func (v Vec2) Dot(v2 Vec2) int32 {
	return v.x*v2.x + v.y*v2.y
}

func (v Vec2) Cross(v2 Vec2) int32 {
	return v.x*v2.y - v.y*v2.x
}

func (v Vec2) Translate(x, y int32) Vec2 {
	return Vec2{x: v.x + x, y: v.y + y}
}

func (v Vec2) Scale(sx, sy int32) Vec2 {
	return Vec2{v.x * sx, v.y * sy}
}

// TODO 返回的角度範圍是[-90, 90]
func (v Vec2) ToAngle() Angle {
	return ArcTangent(v.y, v.x)
}

// 返回的角度是[0, 360)
func (v Vec2) ToAngle360() Angle {
	if v.x == 0 && v.y < 0 {
		return Angle{degree: 270}
	}
	if v.y == 0 && v.x < 0 {
		return Angle{degree: 180}
	}
	angle := v.ToAngle()
	if v.x > 0 && v.y < 0 {
		angle.Add(TwoPiAngle())
	} else if v.x < 0 && v.y > 0 {
		angle.Add(PiAngle())
	} else if v.x < 0 && v.y < 0 {
		angle.Add(PiAngle())
	}
	return angle
}
