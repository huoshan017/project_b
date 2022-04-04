package math

type Rotation struct {
	roll, pitch, yaw Angle // Euler angle
	x, y, z, w       int32 // quaternion
}

func NewRotationFromEuler(roll, pitch, yaw Angle) Rotation {
	rot := Rotation{roll: roll, pitch: pitch, yaw: yaw}

	// Angles increase clockwise
	qr := NewAngle(-roll.value / 2)
	qp := NewAngle(-pitch.value / 2)
	qy := NewAngle(-yaw.value / 2)
	cr := qr.Cos()
	sr := qr.Sin()
	cp := qp.Cos()
	sp := qp.Sin()
	cy := qy.Cos()
	sy := qy.Sin()

	// Normalise to 1024 = 1.0
	sqrtOneValueUnits := oneValueUnits * oneValueUnits
	rot.x = (sr*cp*cy - cr*sp*sy) / (sqrtOneValueUnits)
	rot.y = (cr*sp*cy + sr*cp*sy) / (sqrtOneValueUnits)
	rot.z = (cr*cp*sy - sr*sp*cy) / (sqrtOneValueUnits)
	rot.w = (cr*cp*cy + sr*sp*sy) / (sqrtOneValueUnits)
	return rot
}

/**
 * Construct a rotation from an axis and angle.
 * the axis is expected to be normalized to length 1024
 */
func NewRotationWithAxisAndAngle(axis Vector, angle Angle) Rotation {
	rot := Rotation{
		x: axis.X() * NewAngle(-angle.Get()/2).Sin() / full360DegUnits,
		y: axis.Y() * NewAngle(-angle.Get()/2).Sin() / full360DegUnits,
		z: axis.Z() * NewAngle(-angle.Get()/2).Sin() / full360DegUnits,
		w: NewAngle(-angle.Get() / 2).Cos(),
	}
	rot.roll, rot.pitch, rot.yaw = Quaternion2Euler(rot.x, rot.y, rot.z, rot.w)
	return rot
}

func newRotationFromQuaternion(x, y, z, w int32) Rotation {
	rot := Rotation{
		x: x,
		y: y,
		z: z,
		w: w,
	}
	rot.roll, rot.pitch, rot.roll = Quaternion2Euler(rot.x, rot.y, rot.z, rot.w)
	return rot
}

func newRotation(x, y, z, w int32, roll, pitch, yaw Angle) Rotation {
	return Rotation{
		roll:  roll,
		pitch: pitch,
		yaw:   yaw,
		x:     x,
		y:     y,
		z:     z,
		w:     w,
	}
}

func (r *Rotation) Rotate(rot *Rotation) Rotation {
	if *r == NoneRotation {
		return *rot
	}

	if *rot == NoneRotation {
		return *r
	}

	rx := (rot.w*r.x + rot.x*r.w + rot.y*r.z + rot.z*r.y) / oneValueUnits
	ry := (rot.w*r.y + rot.x*r.z + rot.y*r.w + rot.z*r.x) / oneValueUnits
	rz := (rot.w*r.z + rot.x*r.y + rot.y*r.x + rot.z*r.w) / oneValueUnits
	rw := (rot.w*r.w + rot.x*r.x + rot.y*r.y + rot.z*r.z) / oneValueUnits

	return newRotationFromQuaternion(rx, ry, rz, rw)
}

func (r Rotation) WithRoll(roll Angle) Rotation {
	return NewRotationFromEuler(roll, r.pitch, r.yaw)
}

func (r Rotation) WithPitch(pitch Angle) Rotation {
	return NewRotationFromEuler(r.roll, pitch, r.yaw)
}

func (r Rotation) WithYaw(yaw Angle) Rotation {
	return NewRotationFromEuler(r.roll, r.pitch, yaw)
}

func (r Rotation) AsMatrix(mat *Int32Matrix4x4) {
	// Theoretically 1024 squared, but may differ slightly due to rounding
	lsq := r.x*r.x + r.y*r.y + r.z*r.z + r.w*r.w

	// Quaternion components use 10 bits, so there's no risk of overflow
	mat.M11 = lsq - 2*(r.y*r.y+r.z*r.z)
	mat.M12 = 2 * (r.x*r.y + r.z*r.w)
	mat.M13 = 2 * (r.x*r.z - r.y*r.w)
	mat.M14 = 0

	mat.M21 = 2 * (r.x*r.y - r.z*r.w)
	mat.M22 = lsq - 2*(r.x*r.x+r.z*r.z)
	mat.M23 = 2 * (r.y*r.z + r.x*r.w)
	mat.M24 = 0

	mat.M31 = 2 * (r.x*r.z + r.y*r.w)
	mat.M32 = 2 * (r.y*r.z - r.x*r.w)
	mat.M33 = lsq - 2*(r.x*r.x+r.y*r.y)
	mat.M34 = 0

	mat.M41 = 0
	mat.M42 = 0
	mat.M43 = 0
	mat.M44 = lsq
}

func (r Rotation) AsMatrixTo() *Int32Matrix4x4 {
	mat := &Int32Matrix4x4{}
	r.AsMatrix(mat)
	return mat
}

var (
	NoneRotation = NewRotationFromEuler(NewZeroAngle(), NewZeroAngle(), NewZeroAngle())
)

func RotationNone() Rotation {
	return NewRotationFromEuler(NewZeroAngle(), NewZeroAngle(), NewZeroAngle())
}

func RotationFromFacing(facing int32) Rotation {
	return NewRotationFromEuler(NewZeroAngle(), NewZeroAngle(), NewAngleFromFacing(facing))
}

func RotationFromYaw(yaw Angle) Rotation {
	return NewRotationFromEuler(NewZeroAngle(), NewZeroAngle(), yaw)
}

func RotationAdd(a, b *Rotation) Rotation {
	return NewRotationFromEuler(a.roll.Add(b.roll), a.pitch.Add(b.pitch), a.yaw.Add(b.yaw))
}

func RotationSub(a, b *Rotation) Rotation {
	return NewRotationFromEuler(a.roll.Sub(b.roll), a.pitch.Sub(b.pitch), a.yaw.Sub(b.yaw))
}

func RotationNegative(rot *Rotation) Rotation {
	return newRotation(-rot.x, -rot.y, -rot.z, -rot.w, rot.roll.GetNegative(), rot.pitch.GetNegative(), rot.yaw.GetNegative())
}

func RotationEqualTo(a, b *Rotation) bool {
	return a.roll.EqualTo(b.roll) && a.pitch.EqualTo(b.pitch) && a.yaw.EqualTo(b.yaw)
}

func RotationSLerp(a, b *Rotation, mul, div int32) Rotation {
	// This implements the standard spherical linear interpolation
	// between two quaternions, according for OpenRA's interger math
	// conventions and WRot always using (nearly) normalized quaternions
	dot := a.x*b.x + a.y*b.y + a.z*b.z + a.w*b.w
	var flip int32
	if dot >= 0 {
		flip = 1
	} else {
		flip = -1
	}

	// a and b describe the same rotations
	if flip*dot >= oneValueUnits*oneValueUnits {
		return *a
	}

	var theta = AngleArcCos(dot / oneValueUnits)
	var s1 = NewAngle((div - mul) * theta.value / div).Sin()
	var s2 = NewAngle(mul * theta.value / div).Sin()
	var s3 = theta.Sin()

	var x = (a.x*s1 + flip*b.x*s2) / s3
	var y = (a.y*s1 + flip*b.y*s2) / s3
	var z = (a.z*s1 + flip*b.z*s2) / s3
	var w = (a.w*s1 + flip*b.w*s2) / s3

	// Normalize to 1024 == 1.0
	var l = Sqrt64(uint64(x*x)+uint64(y*y)+uint64(z*z)+uint64(w*w), RoundFloor)
	return newRotationFromQuaternion(int32(uint64(oneValueUnits*x)/l), int32(uint64(oneValueUnits*y)/l), int32(uint64(oneValueUnits*z)/l), int32(uint64(oneValueUnits*w)/l))
}

func Quaternion2Euler(x, y, z, w int32) (roll, pitch, yaw Angle) {
	// Theoretically 1024 squared, but may differ slightly due to rounding
	lsq := x*x + y*y + z*z

	srcp := 2 * (w*x + y*z)
	crcp := lsq - 2*(x*x+y*y)
	sp := (w*y - z*x) / half360DegUnits
	sycp := 2 * (w*z + x*y)
	cycp := lsq - 2*(y*y+z*x)

	roll = AngleArcTan(srcp, crcp)
	roll.Negative()

	asp := sp
	if asp < 0 {
		asp = -asp
	}
	if asp >= full360DegUnits {
		if sp < 0 {
			pitch = NewAngle(-quarter360DegUnits)
		} else if sp > 0 {
			pitch = NewAngle(quarter360DegUnits)
		} else {
			pitch = NewAngle(0)
		}
	} else {
		pitch = AngleArcSin(sp)
	}
	pitch.Negative()

	yaw = AngleArcTan(sycp, cycp)
	yaw.Negative()

	return roll, pitch, yaw
}
