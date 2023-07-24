package base

import "testing"

func TestSineCosine(t *testing.T) {
	for degree := int16(0); degree <= 720; degree++ {
		for minute := int16(0); minute < 60; minute += 10 {
			Sine(Angle{degree, minute})
			Cosine(Angle{degree, minute})
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
