package base

type Angle struct {
	degree int16
	minute int16
}

func NewAngle(degree, minute int16) Angle {
	minutes := int32(degree*60 + minute)
	degree, minute = minutes2DegreeAndMinute(minutes)
	return Angle{degree: degree, minute: minute}
}

func (a *Angle) Reset(degree, minute int16) {
	minutes := int32(degree*60 + minute)
	a.degree, a.minute = minutes2DegreeAndMinute(minutes)
}

func (a Angle) Get() (degree, minute int16) {
	return a.degree, a.minute
}

func (a *Angle) Set(minutes int16) {
	a.degree, a.minute = minutes2DegreeAndMinute(int32(minutes))
}

func (a *Angle) Clear() {
	a.degree = 0
	a.minute = 0
}

func (a *Angle) Add(angle Angle) {
	minutes := int32(60*(a.degree+angle.degree) + a.minute + angle.minute)
	a.degree, a.minute = minutes2DegreeAndMinute(minutes)
}

func (a *Angle) AddMinutes(minutes int16) {
	newMinutes := int32(60*a.degree + a.minute + minutes)
	a.degree, a.minute = minutes2DegreeAndMinute(newMinutes)
}

func (a *Angle) Sub(angle Angle) {
	angle.degree = -angle.degree
	angle.minute = -angle.minute
	a.Add(angle)
}

func (a *Angle) Normalize() {
	minutes := a.degree*60 + a.minute
	if minutes >= 360*60 {
		minutes %= (360 * 60)
	} else if minutes < 0 {
		minutes = -minutes
		minutes %= (360 * 60)
		if minutes > 0 {
			minutes = (-minutes + 360*60)
		}
	}
	if minutes == 0 {
		a.degree, a.minute = 0, 0
	} else {
		a.degree = minutes / 60
		a.minute = minutes % 60
	}
}

func minutes2DegreeAndMinute(minutes int32) (int16, int16) {
	if minutes >= 0 {
		return int16(minutes / 60), int16(minutes % 60)
	}
	minutes = -minutes
	return -int16(minutes / 60), -int16(minutes % 60)
}

func (a Angle) DistanceToVec2(distance int32) Vec2 {
	sn, sd := Sine(a)
	cn, cd := Cosine(a)
	return Vec2{distance * cn / cd, distance * sn / sd}
}

func (a Angle) ToMinutes() int16 {
	return a.degree*60 + a.minute
}

func (a Angle) Negative() Angle {
	return Angle{degree: -a.degree, minute: -a.minute}
}

func (a Angle) IsNegative() bool {
	return (a.degree*a.minute >= 0 && (a.degree < 0 || a.minute < 0))
}

func (a *Angle) ToLeft() {
	a.degree = 180
	a.minute = 0
}

func (a *Angle) ToRight() {
	a.degree = 0
	a.minute = 0
}

func (a *Angle) ToUp() {
	a.degree = 90
	a.minute = 0
}

func (a *Angle) ToDown() {
	a.degree = 270
	a.minute = 0
}

func (a Angle) IsLeft() bool {
	return a.degree == 180 && a.minute == 0
}

func (a Angle) IsRight() bool {
	return a.degree == 0 && a.minute == 0
}

func (a Angle) IsUp() bool {
	return a.degree == 90 && a.minute == 0
}

func (a Angle) IsDown() bool {
	return a.degree == 270 && a.minute == 0
}

func (a Angle) Greater(b Angle) bool {
	return AngleGreater(a, b)
}

func (a Angle) GreaterEqual(b Angle) bool {
	return AngleGreater(a, b) || a == b
}

func (a Angle) Less(b Angle) bool {
	return AngleLess(a, b)
}

func (a Angle) LessEqual(b Angle) bool {
	return AngleLess(a, b) || a == b
}

func AngleAdd(a, b Angle) Angle {
	minutes := (a.degree+b.degree)*60 + a.minute + b.minute
	if minutes >= 0 {
		return Angle{degree: minutes / 60, minute: minutes % 60}
	}
	minutes = -minutes
	return Angle{degree: -(minutes / 60), minute: -(minutes % 60)}
}

func AngleSub(a, b Angle) Angle {
	minutes := a.degree*60 + a.minute - (b.degree*60 + b.minute)
	if minutes >= 0 {
		return Angle{degree: minutes / 60, minute: minutes % 60}
	}
	minutes = -minutes
	return Angle{degree: -(minutes / 60), minute: -(minutes % 60)}
}

func AngleGreater(a, b Angle) bool {
	return a.ToMinutes() > b.ToMinutes()
}

func AngleLess(a, b Angle) bool {
	return a.ToMinutes() < b.ToMinutes()
}

func HalfPiAngle() Angle {
	return Angle{degree: 90}
}

func PiAngle() Angle {
	return Angle{degree: 180}
}

func TwoPiAngle() Angle {
	return Angle{degree: 360}
}

func Rotate(x, y int32, x0, y0 int32, angle Angle) (int32, int32) {
	sn, sd := Sine(angle)
	cn, cd := Cosine(angle)
	return (x-x0)*cn/cd - (y-y0)*sn/sd + x0, (x-x0)*sn/sd + (y-y0)*cn/cd + y0
}
