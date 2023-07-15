package base

import "testing"

func TestSineCosine(t *testing.T) {
	for degree := int16(0); degree <= 720; degree++ {
		for minute := int16(0); minute < 60; minute += 10 {
			Sine(Angle{degree, minute})
		}
	}
}
