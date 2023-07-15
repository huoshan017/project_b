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

func (a Angle) ToMinutes() int16 {
	return a.degree*60 + a.minute
}

func (a Angle) Negative() Angle {
	return Angle{degree: -a.degree, minute: -a.minute}
}
