package base

type Angle struct {
	degree int16
	minute int16
}

func NewAngleObj(degree, minute int16) Angle {
	return Angle{degree: degree, minute: minute}
}

func (a *Angle) Reset(degree, minute int16) {
	a.degree = degree
	a.minute = minute
}

func (a Angle) Get() (degree, minute int16) {
	return a.degree, a.minute
}

func (a *Angle) Set(minutes int16) {
	a.degree, a.minute = minutes/60, minutes%60
}

func (a *Angle) Add(angle Angle) {
	a.degree += angle.degree
	a.minute += angle.minute
	if a.degree >= 360 {
		a.degree -= 360
	}
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

func (a Angle) ToDir() Vec2 {
	n, d := Tangent(a)
	return Vec2{d, n}
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
