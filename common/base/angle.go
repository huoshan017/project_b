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

func (a *Angle) Clear() {
	a.degree = 0
	a.minute = 0
}

func (a *Angle) Add(angle Angle) {
	minutes := 60*(a.degree+angle.degree) + a.minute + angle.minute
	if minutes < 0 {
		minutes = -minutes
		a.degree, a.minute = minutes/60, minutes%60
		if a.degree >= 360 {
			a.degree %= 360
		}
		a.degree = -a.degree
		a.minute = -a.minute
	} else {
		a.degree, a.minute = minutes/60, minutes%60
		if a.degree >= 360 {
			a.degree %= 360
		}
	}
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

func (a Angle) ToVec2() Vec2 {
	a.Normalize()
	if a.minute == 0 {
		if a.degree == 0 {
			return Vec2{1, 0}
		} else if a.degree == 90 {
			return Vec2{0, 1}
		} else if a.degree == 180 {
			return Vec2{-1, 0}
		} else if a.degree == 270 {
			return Vec2{0, -1}
		}
	}
	n, d := Tangent(a)
	if (n >= 32768 && n < 65535) || (n > -65535 && n <= -32768) {
		n /= 100
		d /= 100
	} else if (n >= 65535 && n < 65535*2) || (n > -65535*2 && n <= -65535) {
		n /= 200
		d /= 200
	} else if (n >= 65535*2 && n < 65535*4) || (n > -65535*4 && n <= -65535*2) {
		n /= 400
		d /= 400
	} else if (n >= 65535*4 && n < 65535*8) || (n > -65535*8 && n <= -65535*4) {
		n /= 1000
		d /= 1000
	} else if (n >= 65535*8 && n < 65535*16) || (n > -65535*16 && n <= -65535*8) {
		n /= 1000
		d /= 1000
	} else if (n >= 65535*16 && n < 65535*32) || (n > -65535*32 && n <= -65535*16) {
		n /= 2000
		d /= 2000
	} else if (n >= 65535*32 && n < 65535*64) || (n > -65535*64 && n <= -65535*32) {
		n /= 5000
		d /= 5000
	} else if n >= 65535*64 || n < -65535*64 {
		n /= 10000
		d /= 10000
	}
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

func (a Angle) GreaterEqual(b Angle) bool {
	return AngleGreater(a, b) || a == b
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
	a.Normalize()
	b.Normalize()
	return a.ToMinutes() > b.ToMinutes()
}

func AngleLess(a, b Angle) bool {
	a.Normalize()
	b.Normalize()
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
