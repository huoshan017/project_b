package math

import (
	"log"
	"math"
	"testing"
)

func TestAngle(t *testing.T) {
	for a := int32(0); a <= oneValueUnits; a++ {
		var angle = NewAngle(a)
		/* 计算cos */
		v1 := float32(angle.Cos()) / float32(oneValueUnits)
		// math.Cos参数是弧度
		v2 := float32(math.Cos(float64(a) * math.Pi / float64(half360DegUnits)))
		var delta = v1 - v2
		if delta > 0.001 {
			log.Printf("value a=%v  Angle.cos(a/oneValueUnits):v1=%v  math.cos(a*PI/half360DegUnits)v2:=%v  delta=%v", a, v1, v2, delta)
		}

		/* 计算sin */
		v1 = float32(angle.Sin()) / float32(oneValueUnits)
		v2 = float32(math.Sin(float64(a) * math.Pi / float64(half360DegUnits)))
		delta = v1 - v2
		if delta > 0.001 {
			log.Printf("value a=%v  Angle.sin(a/oneValueUnits):v1=%v  math.sin(a*PI/half360DegUnits)v2:=%v  delta=%v", a, v1, v2, delta)
		}

		/* 计算tan */
		v1 = float32(angle.Tan()) / float32(oneValueUnits)
		v2 = float32(math.Tan(float64(a) * math.Pi / float64(half360DegUnits)))
		delta = v1 - v2
		if delta > 0.001 {
			log.Printf("value a=%v  Angle.tan(a/oneValueUnits):v1=%v  math.sin(a*PI/half360DegUnit)v2:=%v  delta=%v", a, v1, v2, delta)
		}
	}
}
