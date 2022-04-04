package math

type Int32Matrix4x4 struct {
	M11, M12, M13, M14 int32
	M21, M22, M23, M24 int32
	M31, M32, M33, M34 int32
	M41, M42, M43, M44 int32
}

func NewInt32Matrix4x4(m11, m12, m13, m14, m21, m22, m23, m24, m31, m32, m33, m34, m41, m42, m43, m44 int32) Int32Matrix4x4 {
	m := Int32Matrix4x4{
		M11: m11, M12: m12, M13: m13, M14: m14,
		M21: m21, M22: m22, M23: m23, M24: m24,
		M31: m31, M32: m32, M33: m33, M34: m34,
		M41: m41, M42: m42, M43: m43, M44: m44,
	}
	return m
}

func (m Int32Matrix4x4) EqualTo(ma *Int32Matrix4x4) bool {
	return m.M11 == ma.M11 && m.M12 == ma.M12 && m.M13 == ma.M13 && m.M14 == ma.M14 &&
		m.M21 == ma.M21 && m.M22 == ma.M22 && m.M23 == ma.M23 && m.M24 == ma.M24 &&
		m.M31 == ma.M31 && m.M32 == ma.M32 && m.M33 == ma.M33 && m.M34 == ma.M34 &&
		m.M41 == ma.M41 && m.M42 == ma.M42 && m.M43 == ma.M43 && m.M44 == ma.M44
}

func (m Int32Matrix4x4) NotEqualTo(ma *Int32Matrix4x4) bool {
	return !m.EqualTo(ma)
}

func Int32Matrix4x4IsEqual(a, b *Int32Matrix4x4) bool {
	return a.EqualTo(b)
}
