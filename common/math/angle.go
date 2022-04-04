package math

import "math"

var (
	ZeroAngle = NewZeroAngle()
)

/**
 * 1D angle - 1024 units = 360 degrees.
 */
type Angle struct {
	value int32
}

/**
 * NewAngle
 */
func NewAngle(a int32) Angle {
	a = normalAngle(a)
	return Angle{value: a}
}

/**
 * NewZeroAngle
 */
func NewZeroAngle() Angle {
	return NewAngle(0)
}

/**
 * NewAngleFromFacing
 */
func NewAngleFromFacing(facing int32) Angle {
	return NewAngle(facing * 4)
}

/**
 * NewAngleFromDegrees
 */
func NewAngleFromDegrees(degrees int32) Angle {
	return NewAngle(degrees * full360DegUnits / 360)
}

/**
 * calcAngle
 */
func normalAngle(a int32) int32 {
	a %= full360DegUnits
	if a < 0 {
		a += full360DegUnits
	}
	return a
}

/**
 * Angle.Get
 */
func (a Angle) Get() int32 {
	return a.value
}

func (a Angle) GetUnits() int32 {
	return full360DegUnits
}

/**
 * Angle.Add
 */
func (a *Angle) Add(angle Angle) Angle {
	a.value = normalAngle(a.value + angle.value)
	return *a
}

/**
 * Angle.Sub
 */
func (a *Angle) Sub(angle Angle) Angle {
	a.value = normalAngle(a.value - angle.value)
	return *a
}

/**
 * Angle.Negative
 */
func (a *Angle) Negative() Angle {
	a.value = normalAngle(-a.value)
	return *a
}

/**
 * Angle.GetNegative
 */
func (a Angle) GetNegative() Angle {
	return NewAngle(-a.value)
}

/**
 * Angle.EqualTo
 */
func (a Angle) EqualTo(angle Angle) bool {
	return a.value == angle.value
}

/**
 * Angle.NotEqualTo
 */
func (a Angle) NotEqualTo(angle Angle) bool {
	return a.value != angle.value
}

/**
 * Angle.Cos
 */
func (a Angle) Cos() int32 {
	var c int32 = a.value
	if c > half360DegUnits {
		c -= half360DegUnits
	}
	if c <= quarter360DegUnits {
		c = cosineTable[c]
	} else if c <= half360DegUnits {
		c = -cosineTable[half360DegUnits-c]
	}
	if a.value > half360DegUnits {
		c = -c
	}
	return c
}

/**
 * Angle.Sin
 */
func (a Angle) Sin() int32 {
	return NewAngle(a.value - quarter360DegUnits).Cos()
}

/**
 * Angle.Tan
 */
func (a Angle) Tan() int32 {
	c := a.value
	if c > half360DegUnits {
		c -= half360DegUnits
	}
	if c <= quarter360DegUnits {
		c = tanTable[c]
	} else if c <= half360DegUnits {
		c = -tanTable[half360DegUnits-c]
	}
	return c
}

/**
 * Angle.FloatRadians
 */
func (a Angle) FloatRadians() float32 {
	return float32(a.value) * math.Pi / float32(half360DegUnits)
}

/**
 * Angle.FloatDegrees
 */
func (a Angle) FloatDegrees() float32 {
	return float32(a.value) / math.Pi
}

/**
 * AngleLerp
 */
func AngleLerp(a, b Angle, mul, div int32) Angle {
	aa := a.value
	bb := b.value
	if aa > bb && aa-bb > half360DegUnits {
		aa -= full360DegUnits
	}
	if bb > aa && bb-aa > half360DegUnits {
		bb -= full360DegUnits
	}
	return NewAngle(aa + (bb-aa)*mul/div)
}

/**
 * ArcSin
 */
func AngleArcSin(d int32) Angle {
	if d < -oneValueUnits || d > oneValueUnits {
		panic("common.math: ArcSin is only valid for values between -1024 and 1024.")
	}
	if d < 0 {
		return NewAngle(half360DegUnits + quarter360DegUnits + closestCosineIndex(-d))
	}
	return NewAngle(quarter360DegUnits - closestCosineIndex(d))
}

/**
 * ArcCos
 */
func AngleArcCos(d int32) Angle {
	if d < -oneValueUnits || d > oneValueUnits {
		panic("common.math: ArcCos is only valid for values between -1024 and 1024.")
	}
	if d < 0 {
		return NewAngle(half360DegUnits - closestCosineIndex(-d))
	}
	return NewAngle(closestCosineIndex(d))
}

/**
 * ArcTan
 */
func AngleArcTan(y, x int32) Angle {
	return angleArcTan(y, x, 1)
}

// arcTan
func angleArcTan(y, x int32, stride int32) Angle {
	if y == 0 {
		if x >= 0 {
			return NewAngle(0)
		}
		return NewAngle(half360DegUnits)
	}
	if x == 0 {
		if y > 0 {
			return NewAngle(quarter360DegUnits)
		} else if y < 0 {
			return NewAngle(-quarter360DegUnits)
		}
		return NewAngle(0)
	}

	yy := y
	if y < 0 {
		yy = -y
	}
	xx := x
	if x < 0 {
		xx = -x
	}
	bestVal := int32(math.MaxInt32)
	bestAngle := int32(0)
	for i := int32(0); i < quarter360DegUnits; i = i + stride {
		v := oneValueUnits*yy - xx*tanTable[i]
		if v < 0 {
			v = -v
		}
		if v < int32(bestVal) {
			bestVal = v
			bestAngle = i
		}
	}

	if x < 0 && y > 0 {
		bestAngle = half360DegUnits - bestAngle
	} else if x < 0 && y < 0 {
		bestAngle = half360DegUnits + bestAngle
	} else if x > 0 && y < 0 {
		bestAngle = full360DegUnits - bestAngle
	}
	return NewAngle(bestAngle)
}

// inner function
func closestCosineIndex(value int32) int32 {
	a := int32(0)
	b := quarter360DegUnits
	for a != b-1 {
		m := (a + b) / 2
		v := cosineTable[m]

		if v == value {
			return int32(m)
		}

		if v < value {
			b = m
		} else {
			a = m
		}
	}
	av := cosineTable[a] - value
	vb := value - cosineTable[b]
	if av > vb {
		return b
	}
	return a
}

var (
	/**
	 * table of cosine
	 */
	cosineTable = []int32{
		1024, 1023, 1023, 1023, 1023, 1023, 1023, 1023, 1022, 1022, 1022, 1021,
		1021, 1020, 1020, 1019, 1019, 1018, 1017, 1017, 1016, 1015, 1014, 1013,
		1012, 1011, 1010, 1009, 1008, 1007, 1006, 1005, 1004, 1003, 1001, 1000,
		999, 997, 996, 994, 993, 991, 990, 988, 986, 985, 983, 981, 979, 978,
		976, 974, 972, 970, 968, 966, 964, 962, 959, 957, 955, 953, 950, 948,
		946, 943, 941, 938, 936, 933, 930, 928, 925, 922, 920, 917, 914, 911,
		908, 906, 903, 900, 897, 894, 890, 887, 884, 881, 878, 875, 871, 868,
		865, 861, 858, 854, 851, 847, 844, 840, 837, 833, 829, 826, 822, 818,
		814, 811, 807, 803, 799, 795, 791, 787, 783, 779, 775, 771, 767, 762,
		758, 754, 750, 745, 741, 737, 732, 728, 724, 719, 715, 710, 706, 701,
		696, 692, 687, 683, 678, 673, 668, 664, 659, 654, 649, 644, 639, 634,
		629, 625, 620, 615, 609, 604, 599, 594, 589, 584, 579, 574, 568, 563,
		558, 553, 547, 542, 537, 531, 526, 521, 515, 510, 504, 499, 493, 488,
		482, 477, 471, 466, 460, 454, 449, 443, 437, 432, 426, 420, 414, 409,
		403, 397, 391, 386, 380, 374, 368, 362, 356, 350, 344, 339, 333, 327,
		321, 315, 309, 303, 297, 291, 285, 279, 273, 267, 260, 254, 248, 242,
		236, 230, 224, 218, 212, 205, 199, 193, 187, 181, 175, 168, 162, 156,
		150, 144, 137, 131, 125, 119, 112, 106, 100, 94, 87, 81, 75, 69, 62,
		56, 50, 43, 37, 31, 25, 18, 12, 6, 0,
	}

	/**
	 * table of tan
	 */
	tanTable = []int32{
		0, 6, 12, 18, 25, 31, 37, 44, 50, 56, 62, 69, 75, 81, 88, 94, 100, 107,
		113, 119, 126, 132, 139, 145, 151, 158, 164, 171, 177, 184, 190, 197,
		203, 210, 216, 223, 229, 236, 243, 249, 256, 263, 269, 276, 283, 290,
		296, 303, 310, 317, 324, 331, 338, 345, 352, 359, 366, 373, 380, 387,
		395, 402, 409, 416, 424, 431, 438, 446, 453, 461, 469, 476, 484, 492,
		499, 507, 515, 523, 531, 539, 547, 555, 563, 571, 580, 588, 596, 605,
		613, 622, 630, 639, 648, 657, 666, 675, 684, 693, 702, 711, 721, 730,
		740, 749, 759, 769, 779, 789, 799, 809, 819, 829, 840, 850, 861, 872,
		883, 894, 905, 916, 928, 939, 951, 963, 974, 986, 999, 1011, 1023, 1036,
		1049, 1062, 1075, 1088, 1102, 1115, 1129, 1143, 1158, 1172, 1187, 1201,
		1216, 1232, 1247, 1263, 1279, 1295, 1312, 1328, 1345, 1363, 1380, 1398,
		1416, 1435, 1453, 1473, 1492, 1512, 1532, 1553, 1574, 1595, 1617, 1639,
		1661, 1684, 1708, 1732, 1756, 1782, 1807, 1833, 1860, 1887, 1915, 1944,
		1973, 2003, 2034, 2065, 2098, 2131, 2165, 2199, 2235, 2272, 2310, 2348,
		2388, 2429, 2472, 2515, 2560, 2606, 2654, 2703, 2754, 2807, 2861, 2918,
		2976, 3036, 3099, 3164, 3232, 3302, 3375, 3451, 3531, 3613, 3700, 3790,
		3885, 3984, 4088, 4197, 4311, 4432, 4560, 4694, 4836, 4987, 5147, 5318,
		5499, 5693, 5901, 6124, 6364, 6622, 6903, 7207, 7539, 7902, 8302, 8743,
		9233, 9781, 10396, 11094, 11891, 12810, 13882, 15148, 16667, 18524, 20843,
		23826, 27801, 33366, 41713, 55622, 83438, 166883, math.MaxInt32,
	}
)
