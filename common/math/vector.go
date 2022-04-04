package math

type Vector struct {
	x, y, z int32
}

func NewVector(x, y, z int32) Vector {
	return Vector{x: x, y: y, z: z}
}

func NewVectorWithDist(x, y, z Dist) Vector {
	return Vector{
		x: x.length,
		y: y.length,
		z: z.length,
	}
}

func NewVectorWithPosition(pos *Position) Vector {
	return Vector{
		x: pos.x,
		y: pos.y,
		z: pos.z,
	}
}

func NewZeroVector() Vector {
	return Vector{x: 0, y: 0, z: 0}
}

func (v Vector) X() int32 {
	return v.x
}

func (v Vector) Y() int32 {
	return v.y
}

func (v Vector) Z() int32 {
	return v.z
}

func (v *Vector) Add(vec Vector) {
	v.x += vec.x
	v.y += vec.y
	v.z += vec.z
}

func VectorAdd(a, b Vector) Vector {
	v := Vector{}
	v.x = a.x + b.x
	v.y = a.y + b.y
	v.z = a.z + b.z
	return v
}

func (v *Vector) Sub(vec Vector) {
	v.x -= vec.x
	v.y -= vec.y
	v.z -= vec.z
}

func VectorSub(a, b Vector) Vector {
	return Vector{x: a.x - b.x, y: a.y - b.y, z: a.z - b.z}
}

func VectorNegative(vec Vector) Vector {
	return Vector{x: -vec.x, y: -vec.y, z: -vec.z}
}

func (v *Vector) Div(divisor int32) {
	v.x /= divisor
	v.y /= divisor
	v.z /= divisor
}

func VectorDiv(vec Vector, divisor int32) Vector {
	return Vector{x: vec.x / divisor, y: vec.y / divisor, z: vec.z / divisor}
}

func (v *Vector) Mul(a int32) {
	v.x *= a
	v.y *= a
	v.z *= a
}

func VectorMul(vec Vector, a int32) Vector {
	return Vector{x: vec.x * a, y: vec.y * a, z: vec.z * a}
}

func (v Vector) EqualTo(vec Vector) bool {
	return v.x == vec.x && v.y == vec.y && v.z == vec.z
}

func VectorEqual(a, b Vector) bool {
	return a.EqualTo(b)
}

func (v Vector) NotEqualTo(vec Vector) bool {
	return v.x != vec.x || v.y != vec.y || v.z != vec.z
}

func VectorNotEqual(a, b Vector) bool {
	return a.NotEqualTo(b)
}

func (v Vector) Dot(vec Vector) int32 {
	return v.x*vec.x + v.y*vec.y + v.z*vec.z
}

func VectorDot(a, b Vector) int32 {
	return a.Dot(b)
}

func (v Vector) LengthSquared() uint64 {
	return uint64(v.x*v.x) + uint64(v.y*v.y) + uint64(v.z*v.z)
}

func (v Vector) Length() uint32 {
	return uint32(Sqrt64Floor(v.LengthSquared()))
}

func (v Vector) HorizontalLengthSquared() uint64 {
	return uint64(v.x*v.x) + uint64(v.y*v.y)
}

func (v Vector) HorizontalLength() uint32 {
	return uint32(Sqrt64Floor(v.HorizontalLengthSquared()))
}

func (v Vector) VerticalLengthSquared() uint64 {
	return uint64(v.z * v.z)
}

func (v Vector) VerticalLength() uint32 {
	return uint32(Sqrt64Floor(v.VerticalLengthSquared()))
}

func (v *Vector) Rotate(rot *Rotation) {
	var mat Int32Matrix4x4
	rot.AsMatrix(&mat)
	v.RotateWithMatrix(&mat)
}

func (v *Vector) RotateWithMatrix(mat *Int32Matrix4x4) {
	v.x = (v.x*mat.M11 + v.y*mat.M21 + v.z*mat.M31) / mat.M44
	v.y = (v.x*mat.M12 + v.y*mat.M22 + v.z*mat.M23) / mat.M44
	v.z = (v.x*mat.M13 + v.y*mat.M23 + v.z*mat.M33) / mat.M44
}

func (v Vector) Yaw() Angle {
	if v.LengthSquared() == 0 {
		return ZeroAngle
	}

	a := AngleArcTan(-v.Y(), v.X())
	return a.Sub(NewAngle(quarter360DegUnits))
}

func VectorLerp(a, b Vector, mul, div int32) Vector {
	// a + (b-a)*mul/div
	b.Sub(a)
	b.Mul(mul)
	b.Div(div)
	a.Add(b)
	return a
}

func VectorLerpQuadratic(a, b Vector, pitch Angle, mul, div int32) Vector {
	// Start with a linear lerp between the points
	var ret = VectorLerp(a, b, mul, div)
	if pitch.Get() == 0 {
		return ret
	}

	// Add an additional quadratic variation to height
	// Uses decimal to avoid integer overflow
	b.Sub(a)
	var offset = (int32(b.Length()) * pitch.Tan() * mul * (div - mul) / (oneValueUnits * div * div))
	return NewVector(ret.X(), ret.Y(), ret.Z()+offset)
}

// Sampled a N-sample probability density function in the range [-1024..1024, -1024..1024]
// 1 sample produces a rectangular probability
// 2 samples produces a triangular probability
// ...
// N samples approximates a true Gaussian
func VectorFromPDF(m *MersennelTwister, samples int32) Vector {
	return NewVectorWithDist(DistFromPDF(m, samples), DistFromPDF(m, samples), DistZero)
}
