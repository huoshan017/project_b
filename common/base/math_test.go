package base

import "testing"

func TestSine(t *testing.T) {
	for degree := int16(0); degree < 90; degree++ {
		for minute := int16(0); minute < 60; minute += 10 {
			angle := NewAngle(degree, minute)
			sn, sd := Sine(angle)
			// 比較sine(a) sine(pi-a)
			angle2 := AngleSub(PiAngle(), angle)
			sn2, sd2 := Sine(angle2)
			if sn != sn2 || sd != sd2 {
				t.Errorf("error: sine(%v)=(%v,%v) not equal to sine(%v)=(%v,%v)", angle, sn, sd, angle2, sn2, sd2)
			}
			// 比較sine(a) sine(pi+a)
			angle3 := AngleAdd(PiAngle(), angle)
			sn3, sd3 := Sine(angle3)
			if sn != -sn3 || sd != sd3 {
				t.Errorf("error: sine(%v)=(%v,%v)  sine(%v)=(%v,%v)", angle, sn, sd, angle3, sn3, sd3)
			}
			// 比較sin(a) sine(2pi-a)
			angle4 := AngleSub(TwoPiAngle(), angle)
			sn4, sd4 := Sine(angle4)
			if sn != -sn4 || sd != sd4 {
				t.Errorf("error: sine(%v)=(%v,%v)  sine(%v)=(%v,%v)", angle, sn, sd, angle4, sn4, sd4)
			}
		}
	}
}

func TestCosine(t *testing.T) {
	for degree := int16(0); degree < 90; degree++ {
		for minute := int16(0); minute < 60; minute += 10 {
			angle := NewAngle(degree, minute)
			cn, cd := Cosine(angle)
			// 比較cosine(a) cosine(pi-a)
			angle2 := AngleSub(PiAngle(), angle)
			cn2, cd2 := Cosine(angle2)
			if cn != -cn2 || cd != cd2 {
				t.Errorf("error: cosine(%v)=(%v,%v)  cosine(%v)=(%v,%v)", angle, cn, cd, angle2, cn2, cd2)
			}
			// 比較cosine(a) cosine(pi+a)
			angle3 := AngleAdd(PiAngle(), angle)
			cn3, cd3 := Cosine(angle3)
			if cn != -cn3 || cd != cd3 {
				t.Errorf("error: cosine(%v)=(%v,%v)  cosine(%v)=(%v,%v)", angle, cn, cd, angle3, cn3, cd3)
			}
			// 比較cosin(a) cosine(2pi-a)
			angle4 := AngleSub(TwoPiAngle(), angle)
			cn4, cd4 := Cosine(angle4)
			if cn != cn4 || cd != cd4 {
				t.Errorf("error: cosine(%v)=(%v,%v)  cosine(%v)=(%v,%v)", angle, cn, cd, angle4, cn4, cd4)
			}
		}
	}
}

func TestTangent(t *testing.T) {
	for degree := int16(0); degree < 90; degree++ {
		for minute := int16(0); minute < 60; minute += 10 {
			angle := NewAngle(degree, minute)
			tn, td := Tangent(angle)
			// 比較tangent(a) tangent(pi-a)
			angle2 := AngleSub(PiAngle(), angle)
			tn2, td2 := Tangent(angle2)
			if tn != -tn2 || td != td2 {
				t.Errorf("error: (tangent(a) vs tangent(pi-a)) tangent(%v)=(%v,%v)  tangent(%v)=(%v,%v)", angle, tn, td, angle2, tn2, td2)
			}
			// 比較tangent(a) tangent(pi+a)
			angle3 := AngleAdd(PiAngle(), angle)
			tn3, td3 := Tangent(angle3)
			if tn != tn3 || td != td3 {
				t.Errorf("error: (tangent(a) vs tangent(pi+a)) tangent(%v)=(%v,%v)  tangent(%v)=(%v,%v)", angle, tn, td, angle3, tn3, td3)
			}
			// 比較tangent(a) tangent(2pi-a)
			angle4 := AngleSub(TwoPiAngle(), angle)
			tn4, td4 := Tangent(angle4)
			if tn != -tn4 || td != td4 {
				t.Errorf("error: (tangent(a) vs tangent(2pi-a)) tangent(%v)=(%v,%v)  tangent(%v)=(%v,%v)", angle, tn, td, angle4, tn4, td4)
			}
		}
	}
}

func TestCotangent(t *testing.T) {
	for degree := int16(0); degree < 90; degree++ {
		for minute := int16(0); minute < 60; minute += 10 {
			angle := NewAngle(degree, minute)
			cn, cd := Cotangent(angle)
			// 比較cotangent(a) cotangent(pi-a)
			angle2 := AngleSub(PiAngle(), angle)
			cn2, cd2 := Cotangent(angle2)
			if cn != -cn2 || cd != cd2 {
				t.Errorf("error: (cotangent(a) vs cotangent(pi-a)) cotangent(%v)=(%v,%v)  cotangent(%v)=(%v,%v)", angle, cn, cd, angle2, cn2, cd2)
			}
			// 比較cotangent(a) cotangent(pi+a)
			angle3 := AngleAdd(PiAngle(), angle)
			cn3, cd3 := Cotangent(angle3)
			if cn != cn3 || cd != cd3 {
				t.Errorf("error: (cotangent(a) vs cotangent(pi+a)) cotangent(%v)=(%v,%v)  cotangent(%v)=(%v,%v)", angle, cn, cd, angle3, cn3, cd3)
			}
			// 比較cotangent(a) cotangent(2pi-a)
			angle4 := AngleSub(TwoPiAngle(), angle)
			cn4, cd4 := Cotangent(angle4)
			if cn != -cn4 || cd != cd4 {
				t.Errorf("error: (cotangent(a) vs cotangent(2pi-a)) cotangent(%v)=(%v,%v)  cotangent(%v)=(%v,%v)", angle, cn, cd, angle4, cn4, cd4)
			}
		}
	}
}

func TestTangentCotangent(t *testing.T) {
	for degree := int16(0); degree <= 720; degree++ {
		for minute := int16(0); minute < 60; minute += 10 {
			Tangent(Angle{degree, minute})
			Cotangent(Angle{degree, minute})
		}
	}
}

func TestArcSine(t *testing.T) {
	for degree := int16(0); degree >= -90; degree-- {
		for minute := int16(0); minute >= -50; minute -= 10 {
			sn, sd := Sine(Angle{degree, minute})
			angle := ArcSine(sn, sd)
			if angle.degree != degree || angle.minute != minute {
				tsn, tsd := Sine(angle)
				if tsn != sn || tsd != sd {
					t.Errorf("arcsine get invalid angle %v not equaivalent original degree(%v) minute(%v)", angle, degree, minute)
				}
			}
			if degree == -90 {
				break
			}
		}
	}
	for degree := int16(0); degree <= 90; degree++ {
		for minute := int16(0); minute <= 50; minute += 10 {
			sn, sd := Sine(Angle{degree, minute})
			angle := ArcSine(sn, sd)
			if angle.degree != degree || angle.minute != minute {
				tsn, tsd := Sine(angle)
				if tsn != sn || tsd != sd {
					t.Errorf("arcsine get invalid angle %v not equaivalent original degree(%v) minute(%v)", angle, degree, minute)
				}
			}
			if degree == 90 {
				break
			}
		}
	}
}

func TestArcCosine(t *testing.T) {
	for degree := int16(0); degree <= 180; degree++ {
		for minute := int16(0); minute <= 50; minute += 10 {
			cn, cd := Cosine(Angle{degree, minute})
			angle := ArcCosine(cn, cd)
			if angle.degree != degree || angle.minute != minute {
				tcn, tcd := Cosine(angle)
				if tcn != cn || tcd != cd {
					t.Errorf("arccosine get invalid angle %v not equaivalent original degree(%v) minute(%v)", angle, degree, minute)
				}
			}
			if degree == 180 {
				break
			}
		}
	}
}

func TestArcTangent(t *testing.T) {
	for degree := int16(0); degree >= -90; degree-- {
		for minute := int16(0); minute >= -50; minute -= 10 {
			tn, td := Tangent(Angle{degree, minute})
			angle := ArcTangent(tn, td)
			if angle.degree != degree || angle.minute != minute {
				ttn, ttd := Tangent(angle)
				if ttn != tn || ttd != td {
					t.Errorf("arctangent get invalid angle %v not equaivalent original degree(%v) minute(%v)  tn_td(%v %v) ttn_ttd(%v %v)", angle, degree, minute, tn, td, ttn, ttd)
				}
			}
			if degree == -90 {
				break
			}
		}
	}
	for degree := int16(0); degree <= 90; degree++ {
		for minute := int16(0); minute <= 50; minute += 10 {
			tn, td := Tangent(Angle{degree, minute})
			angle := ArcTangent(tn, td)
			if angle.degree != degree || angle.minute != minute {
				ttn, ttd := Tangent(angle)
				if ttn != tn || ttd != td {
					t.Errorf("arctangent get invalid angle %v not equaivalent original degree(%v) minute(%v)  tn_td(%v %v) ttn_ttd(%v %v)", angle, degree, minute, tn, td, ttn, ttd)
				}
			}
			if degree == 90 {
				break
			}
		}
	}
}

func TestArcCotangent(t *testing.T) {
	for degree := int16(0); degree <= 180; degree++ {
		for minute := int16(0); minute <= 50; minute += 10 {
			cn, cd := Cotangent(Angle{degree, minute})
			angle := ArcCotangent(cn, cd)
			if angle.degree != degree || angle.minute != minute {
				ccn, ccd := Cotangent(angle)
				if ccn != cn || ccd != cd {
					t.Errorf("arccotangent get invalid angle %v not equaivalent original degree(%v) minute(%v)  cn_cd(%v %v) ccn_ccd(%v %v)", angle, degree, minute, cn, cd, ccn, ccd)
				}
			}
			if degree == 180 {
				break
			}
		}
	}
}
