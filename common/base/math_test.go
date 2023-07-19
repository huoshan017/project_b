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

func TestArcTangent(t *testing.T) {
	for i := 0; i < len(tanval); i++ {
		for j := 0; j < len(tanval[i]); j++ {
			ArcTangent(tanval[i][j], denominator)
		}
	}
}
