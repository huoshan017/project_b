package base

import "math"

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

func (v Vec2) Rotate(angle int32) Vec2 {
	s, c := math.Sincos(math.Pi * float64(angle) / 180)
	return Vec2{x: int32(float64(v.x)*c - float64(v.y)*s), y: int32(float64(v.x)*s + float64(v.y)*c)}
}

func (v Vec2) Scale(sx, sy int32) Vec2 {
	return Vec2{v.x * sx, v.y * sy}
}

func (v Vec2) ToAngle() Angle {
	return ArcTangent(v.y, v.x)
}
