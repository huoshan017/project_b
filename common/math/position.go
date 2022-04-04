package math

import "math"

var (
	ZeroPosition = NewPosition(0, 0, 0)
)

type Position struct {
	x, y, z int32
}

func NewPosition(x, y, z int32) Position {
	return Position{x: x, y: y, z: z}
}

func NewPositionWithDist(x, y, z Dist) Position {
	return Position{x: x.length, y: y.length, z: z.length}
}

func (p *Position) Add(pos Position) {
	p.x += pos.x
	p.y += pos.y
	p.z += pos.z
}

func (p *Position) Sub(pos Position) {
	p.x -= pos.x
	p.y -= pos.y
	p.z -= pos.z
}

func (p *Position) Mul(mul int32) {
	p.x *= mul
	p.y *= mul
	p.z *= mul
}

func (p *Position) Div(divisor int32) {
	p.x /= divisor
	p.y /= divisor
	p.z /= divisor
}

func (p Position) EqualTo(pos Position) bool {
	return p.x == pos.x && p.y == pos.y && p.z == pos.z
}

func (p Position) NotEqualTo(pos Position) bool {
	return p.x != pos.x || p.y != pos.y || p.z != pos.z
}

func SubWithPosition2Vector(a, b *Position) Vector {
	return NewVector(a.x-b.x, a.y-b.y, a.z-b.z)
}

func PositionLerp(a, b Position, mul, div int32) Position {
	// (a + (b-a)*mul/div)
	b.Sub(a)
	b.Mul(mul)
	b.Div(div)
	a.Add(b)
	return a
}

func PositionLerpQuadratic(a, b Position, pitch Angle, mul, div int32) Position {
	// Start with a linear lerp between the points
	var ret = PositionLerp(a, b, mul, div)

	if pitch.Get() == 0 {
		return ret
	}

	// Add an additional quadratic variation to height
	// uses decimal to avoid integer overflow
	var v = SubWithPosition2Vector(&b, &a)
	var l = int64(int32(v.Length()) * pitch.Tan() * mul * (div - mul))
	var offset = float64(l) / float64(oneValueUnits*div*div)
	offset += float64(ret.z)
	clampedOffset := offset
	if clampedOffset < math.MinInt32 {
		clampedOffset = math.MinInt32
	} else if clampedOffset > math.MaxInt32 {
		clampedOffset = math.MaxInt32
	}
	return NewPosition(ret.x, ret.y, int32(clampedOffset))
}
