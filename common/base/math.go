package base

import "fmt"

// 三角函數值
type TrigonometricFunctionValue struct {
	Numerator   int16
	Denominator int16
}

func (sc TrigonometricFunctionValue) Negative() TrigonometricFunctionValue {
	return TrigonometricFunctionValue{-sc.Numerator, sc.Denominator}
}

var (
	sin0_00  = TrigonometricFunctionValue{Numerator: 0, Denominator: 10000}
	sin0_10  = TrigonometricFunctionValue{Numerator: 29, Denominator: 10000}
	sin0_20  = TrigonometricFunctionValue{Numerator: 58, Denominator: 10000}
	sin0_30  = TrigonometricFunctionValue{Numerator: 87, Denominator: 10000}
	sin0_40  = TrigonometricFunctionValue{Numerator: 116, Denominator: 10000}
	sin0_50  = TrigonometricFunctionValue{Numerator: 145, Denominator: 10000}
	sin1_00  = TrigonometricFunctionValue{Numerator: 175, Denominator: 10000}
	sin1_10  = TrigonometricFunctionValue{Numerator: 204, Denominator: 10000}
	sin1_20  = TrigonometricFunctionValue{Numerator: 233, Denominator: 10000}
	sin1_30  = TrigonometricFunctionValue{Numerator: 262, Denominator: 10000}
	sin1_40  = TrigonometricFunctionValue{Numerator: 291, Denominator: 10000}
	sin1_50  = TrigonometricFunctionValue{Numerator: 320, Denominator: 10000}
	sin2_00  = TrigonometricFunctionValue{Numerator: 349, Denominator: 10000}
	sin2_10  = TrigonometricFunctionValue{Numerator: 378, Denominator: 10000}
	sin2_20  = TrigonometricFunctionValue{Numerator: 407, Denominator: 10000}
	sin2_30  = TrigonometricFunctionValue{Numerator: 436, Denominator: 10000}
	sin2_40  = TrigonometricFunctionValue{Numerator: 465, Denominator: 10000}
	sin2_50  = TrigonometricFunctionValue{Numerator: 494, Denominator: 10000}
	sin3_00  = TrigonometricFunctionValue{Numerator: 523, Denominator: 10000}
	sin3_10  = TrigonometricFunctionValue{Numerator: 552, Denominator: 10000}
	sin3_20  = TrigonometricFunctionValue{Numerator: 581, Denominator: 10000}
	sin3_30  = TrigonometricFunctionValue{Numerator: 610, Denominator: 10000}
	sin3_40  = TrigonometricFunctionValue{Numerator: 640, Denominator: 10000}
	sin3_50  = TrigonometricFunctionValue{Numerator: 669, Denominator: 10000}
	sin4_00  = TrigonometricFunctionValue{Numerator: 698, Denominator: 10000}
	sin4_10  = TrigonometricFunctionValue{Numerator: 727, Denominator: 10000}
	sin4_20  = TrigonometricFunctionValue{Numerator: 756, Denominator: 10000}
	sin4_30  = TrigonometricFunctionValue{Numerator: 785, Denominator: 10000}
	sin4_40  = TrigonometricFunctionValue{Numerator: 814, Denominator: 10000}
	sin4_50  = TrigonometricFunctionValue{Numerator: 843, Denominator: 10000}
	sin5_00  = TrigonometricFunctionValue{Numerator: 872, Denominator: 10000}
	sin5_10  = TrigonometricFunctionValue{Numerator: 901, Denominator: 10000}
	sin5_20  = TrigonometricFunctionValue{Numerator: 929, Denominator: 10000}
	sin5_30  = TrigonometricFunctionValue{Numerator: 958, Denominator: 10000}
	sin5_40  = TrigonometricFunctionValue{Numerator: 987, Denominator: 10000}
	sin5_50  = TrigonometricFunctionValue{Numerator: 1016, Denominator: 10000}
	sin6_00  = TrigonometricFunctionValue{Numerator: 1045, Denominator: 10000}
	sin6_10  = TrigonometricFunctionValue{Numerator: 1074, Denominator: 10000}
	sin6_20  = TrigonometricFunctionValue{Numerator: 1103, Denominator: 10000}
	sin6_30  = TrigonometricFunctionValue{Numerator: 1132, Denominator: 10000}
	sin6_40  = TrigonometricFunctionValue{Numerator: 1161, Denominator: 10000}
	sin6_50  = TrigonometricFunctionValue{Numerator: 1190, Denominator: 10000}
	sin7_00  = TrigonometricFunctionValue{Numerator: 1219, Denominator: 10000}
	sin7_10  = TrigonometricFunctionValue{Numerator: 1248, Denominator: 10000}
	sin7_20  = TrigonometricFunctionValue{Numerator: 1276, Denominator: 10000}
	sin7_30  = TrigonometricFunctionValue{Numerator: 1305, Denominator: 10000}
	sin7_40  = TrigonometricFunctionValue{Numerator: 1334, Denominator: 10000}
	sin7_50  = TrigonometricFunctionValue{Numerator: 1363, Denominator: 10000}
	sin8_00  = TrigonometricFunctionValue{Numerator: 1392, Denominator: 10000}
	sin8_10  = TrigonometricFunctionValue{Numerator: 1421, Denominator: 10000}
	sin8_20  = TrigonometricFunctionValue{Numerator: 1449, Denominator: 10000}
	sin8_30  = TrigonometricFunctionValue{Numerator: 1478, Denominator: 10000}
	sin8_40  = TrigonometricFunctionValue{Numerator: 1507, Denominator: 10000}
	sin8_50  = TrigonometricFunctionValue{Numerator: 1536, Denominator: 10000}
	sin9_00  = TrigonometricFunctionValue{Numerator: 1564, Denominator: 10000}
	sin9_10  = TrigonometricFunctionValue{Numerator: 1593, Denominator: 10000}
	sin9_20  = TrigonometricFunctionValue{Numerator: 1622, Denominator: 10000}
	sin9_30  = TrigonometricFunctionValue{Numerator: 1650, Denominator: 10000}
	sin9_40  = TrigonometricFunctionValue{Numerator: 1679, Denominator: 10000}
	sin9_50  = TrigonometricFunctionValue{Numerator: 1708, Denominator: 10000}
	sin10_00 = TrigonometricFunctionValue{Numerator: 1736, Denominator: 10000}
	sin10_10 = TrigonometricFunctionValue{Numerator: 1765, Denominator: 10000}
	sin10_20 = TrigonometricFunctionValue{Numerator: 1794, Denominator: 10000}
	sin10_30 = TrigonometricFunctionValue{Numerator: 1822, Denominator: 10000}
	sin10_40 = TrigonometricFunctionValue{Numerator: 1851, Denominator: 10000}
	sin10_50 = TrigonometricFunctionValue{Numerator: 1880, Denominator: 10000}
	sin11_00 = TrigonometricFunctionValue{Numerator: 1908, Denominator: 10000}
	sin11_10 = TrigonometricFunctionValue{Numerator: 1937, Denominator: 10000}
	sin11_20 = TrigonometricFunctionValue{Numerator: 1965, Denominator: 10000}
	sin11_30 = TrigonometricFunctionValue{Numerator: 1994, Denominator: 10000}
	sin11_40 = TrigonometricFunctionValue{Numerator: 2022, Denominator: 10000}
	sin11_50 = TrigonometricFunctionValue{Numerator: 2051, Denominator: 10000}
	sin12_00 = TrigonometricFunctionValue{Numerator: 2079, Denominator: 10000}
	sin12_10 = TrigonometricFunctionValue{Numerator: 2108, Denominator: 10000}
	sin12_20 = TrigonometricFunctionValue{Numerator: 2136, Denominator: 10000}
	sin12_30 = TrigonometricFunctionValue{Numerator: 2164, Denominator: 10000}
	sin12_40 = TrigonometricFunctionValue{Numerator: 2193, Denominator: 10000}
	sin12_50 = TrigonometricFunctionValue{Numerator: 2221, Denominator: 10000}
	sin13_00 = TrigonometricFunctionValue{Numerator: 2250, Denominator: 10000}
	sin13_10 = TrigonometricFunctionValue{Numerator: 2278, Denominator: 10000}
	sin13_20 = TrigonometricFunctionValue{Numerator: 2306, Denominator: 10000}
	sin13_30 = TrigonometricFunctionValue{Numerator: 2334, Denominator: 10000}
	sin13_40 = TrigonometricFunctionValue{Numerator: 2363, Denominator: 10000}
	sin13_50 = TrigonometricFunctionValue{Numerator: 2391, Denominator: 10000}
	sin14_00 = TrigonometricFunctionValue{Numerator: 2419, Denominator: 10000}
	sin14_10 = TrigonometricFunctionValue{Numerator: 2447, Denominator: 10000}
	sin14_20 = TrigonometricFunctionValue{Numerator: 2476, Denominator: 10000}
	sin14_30 = TrigonometricFunctionValue{Numerator: 2504, Denominator: 10000}
	sin14_40 = TrigonometricFunctionValue{Numerator: 2532, Denominator: 10000}
	sin14_50 = TrigonometricFunctionValue{Numerator: 2560, Denominator: 10000}
	sin15_00 = TrigonometricFunctionValue{Numerator: 2588, Denominator: 10000}
	sin15_10 = TrigonometricFunctionValue{Numerator: 2616, Denominator: 10000}
	sin15_20 = TrigonometricFunctionValue{Numerator: 2644, Denominator: 10000}
	sin15_30 = TrigonometricFunctionValue{Numerator: 2672, Denominator: 10000}
	sin15_40 = TrigonometricFunctionValue{Numerator: 2700, Denominator: 10000}
	sin15_50 = TrigonometricFunctionValue{Numerator: 2728, Denominator: 10000}
	sin16_00 = TrigonometricFunctionValue{Numerator: 2756, Denominator: 10000}
	sin16_10 = TrigonometricFunctionValue{Numerator: 2784, Denominator: 10000}
	sin16_20 = TrigonometricFunctionValue{Numerator: 2812, Denominator: 10000}
	sin16_30 = TrigonometricFunctionValue{Numerator: 2840, Denominator: 10000}
	sin16_40 = TrigonometricFunctionValue{Numerator: 2868, Denominator: 10000}
	sin16_50 = TrigonometricFunctionValue{Numerator: 2896, Denominator: 10000}
	sin17_00 = TrigonometricFunctionValue{Numerator: 2924, Denominator: 10000}
	sin17_10 = TrigonometricFunctionValue{Numerator: 2952, Denominator: 10000}
	sin17_20 = TrigonometricFunctionValue{Numerator: 2979, Denominator: 10000}
	sin17_30 = TrigonometricFunctionValue{Numerator: 3007, Denominator: 10000}
	sin17_40 = TrigonometricFunctionValue{Numerator: 3035, Denominator: 10000}
	sin17_50 = TrigonometricFunctionValue{Numerator: 3062, Denominator: 10000}
	sin18_00 = TrigonometricFunctionValue{Numerator: 3090, Denominator: 10000}
	sin18_10 = TrigonometricFunctionValue{Numerator: 3118, Denominator: 10000}
	sin18_20 = TrigonometricFunctionValue{Numerator: 3145, Denominator: 10000}
	sin18_30 = TrigonometricFunctionValue{Numerator: 3173, Denominator: 10000}
	sin18_40 = TrigonometricFunctionValue{Numerator: 3201, Denominator: 10000}
	sin18_50 = TrigonometricFunctionValue{Numerator: 3228, Denominator: 10000}
	sin19_00 = TrigonometricFunctionValue{Numerator: 3256, Denominator: 10000}
	sin19_10 = TrigonometricFunctionValue{Numerator: 3283, Denominator: 10000}
	sin19_20 = TrigonometricFunctionValue{Numerator: 3311, Denominator: 10000}
	sin19_30 = TrigonometricFunctionValue{Numerator: 3338, Denominator: 10000}
	sin19_40 = TrigonometricFunctionValue{Numerator: 3365, Denominator: 10000}
	sin19_50 = TrigonometricFunctionValue{Numerator: 3393, Denominator: 10000}
	sin20_00 = TrigonometricFunctionValue{Numerator: 3420, Denominator: 10000}
	sin20_10 = TrigonometricFunctionValue{Numerator: 3448, Denominator: 10000}
	sin20_20 = TrigonometricFunctionValue{Numerator: 3475, Denominator: 10000}
	sin20_30 = TrigonometricFunctionValue{Numerator: 3502, Denominator: 10000}
	sin20_40 = TrigonometricFunctionValue{Numerator: 3529, Denominator: 10000}
	sin20_50 = TrigonometricFunctionValue{Numerator: 3557, Denominator: 10000}
	sin21_00 = TrigonometricFunctionValue{Numerator: 3584, Denominator: 10000}
	sin21_10 = TrigonometricFunctionValue{Numerator: 3611, Denominator: 10000}
	sin21_20 = TrigonometricFunctionValue{Numerator: 3638, Denominator: 10000}
	sin21_30 = TrigonometricFunctionValue{Numerator: 3665, Denominator: 10000}
	sin21_40 = TrigonometricFunctionValue{Numerator: 3692, Denominator: 10000}
	sin21_50 = TrigonometricFunctionValue{Numerator: 3719, Denominator: 10000}
	sin22_00 = TrigonometricFunctionValue{Numerator: 3746, Denominator: 10000}
	sin22_10 = TrigonometricFunctionValue{Numerator: 3773, Denominator: 10000}
	sin22_20 = TrigonometricFunctionValue{Numerator: 3800, Denominator: 10000}
	sin22_30 = TrigonometricFunctionValue{Numerator: 3827, Denominator: 10000}
	sin22_40 = TrigonometricFunctionValue{Numerator: 3854, Denominator: 10000}
	sin22_50 = TrigonometricFunctionValue{Numerator: 3881, Denominator: 10000}
	sin23_00 = TrigonometricFunctionValue{Numerator: 3907, Denominator: 10000}
	sin23_10 = TrigonometricFunctionValue{Numerator: 3934, Denominator: 10000}
	sin23_20 = TrigonometricFunctionValue{Numerator: 3961, Denominator: 10000}
	sin23_30 = TrigonometricFunctionValue{Numerator: 3987, Denominator: 10000}
	sin23_40 = TrigonometricFunctionValue{Numerator: 4014, Denominator: 10000}
	sin23_50 = TrigonometricFunctionValue{Numerator: 4041, Denominator: 10000}
	sin24_00 = TrigonometricFunctionValue{Numerator: 4067, Denominator: 10000}
	sin24_10 = TrigonometricFunctionValue{Numerator: 4094, Denominator: 10000}
	sin24_20 = TrigonometricFunctionValue{Numerator: 4120, Denominator: 10000}
	sin24_30 = TrigonometricFunctionValue{Numerator: 4147, Denominator: 10000}
	sin24_40 = TrigonometricFunctionValue{Numerator: 4173, Denominator: 10000}
	sin24_50 = TrigonometricFunctionValue{Numerator: 4200, Denominator: 10000}
	sin25_00 = TrigonometricFunctionValue{Numerator: 4226, Denominator: 10000}
	sin25_10 = TrigonometricFunctionValue{Numerator: 4253, Denominator: 10000}
	sin25_20 = TrigonometricFunctionValue{Numerator: 4279, Denominator: 10000}
	sin25_30 = TrigonometricFunctionValue{Numerator: 4305, Denominator: 10000}
	sin25_40 = TrigonometricFunctionValue{Numerator: 4331, Denominator: 10000}
	sin25_50 = TrigonometricFunctionValue{Numerator: 4358, Denominator: 10000}
	sin26_00 = TrigonometricFunctionValue{Numerator: 4384, Denominator: 10000}
	sin26_10 = TrigonometricFunctionValue{Numerator: 4410, Denominator: 10000}
	sin26_20 = TrigonometricFunctionValue{Numerator: 4436, Denominator: 10000}
	sin26_30 = TrigonometricFunctionValue{Numerator: 4462, Denominator: 10000}
	sin26_40 = TrigonometricFunctionValue{Numerator: 4488, Denominator: 10000}
	sin26_50 = TrigonometricFunctionValue{Numerator: 4514, Denominator: 10000}
	sin27_00 = TrigonometricFunctionValue{Numerator: 4540, Denominator: 10000}
	sin27_10 = TrigonometricFunctionValue{Numerator: 4566, Denominator: 10000}
	sin27_20 = TrigonometricFunctionValue{Numerator: 4592, Denominator: 10000}
	sin27_30 = TrigonometricFunctionValue{Numerator: 4617, Denominator: 10000}
	sin27_40 = TrigonometricFunctionValue{Numerator: 4643, Denominator: 10000}
	sin27_50 = TrigonometricFunctionValue{Numerator: 4669, Denominator: 10000}
	sin28_00 = TrigonometricFunctionValue{Numerator: 4695, Denominator: 10000}
	sin28_10 = TrigonometricFunctionValue{Numerator: 4720, Denominator: 10000}
	sin28_20 = TrigonometricFunctionValue{Numerator: 4746, Denominator: 10000}
	sin28_30 = TrigonometricFunctionValue{Numerator: 4772, Denominator: 10000}
	sin28_40 = TrigonometricFunctionValue{Numerator: 4797, Denominator: 10000}
	sin28_50 = TrigonometricFunctionValue{Numerator: 4823, Denominator: 10000}
	sin29_00 = TrigonometricFunctionValue{Numerator: 4848, Denominator: 10000}
	sin29_10 = TrigonometricFunctionValue{Numerator: 4874, Denominator: 10000}
	sin29_20 = TrigonometricFunctionValue{Numerator: 4899, Denominator: 10000}
	sin29_30 = TrigonometricFunctionValue{Numerator: 4924, Denominator: 10000}
	sin29_40 = TrigonometricFunctionValue{Numerator: 4950, Denominator: 10000}
	sin29_50 = TrigonometricFunctionValue{Numerator: 4975, Denominator: 10000}
	sin30_00 = TrigonometricFunctionValue{Numerator: 5000, Denominator: 10000}
	sin30_10 = TrigonometricFunctionValue{Numerator: 5025, Denominator: 10000}
	sin30_20 = TrigonometricFunctionValue{Numerator: 5050, Denominator: 10000}
	sin30_30 = TrigonometricFunctionValue{Numerator: 5075, Denominator: 10000}
	sin30_40 = TrigonometricFunctionValue{Numerator: 5100, Denominator: 10000}
	sin30_50 = TrigonometricFunctionValue{Numerator: 5125, Denominator: 10000}
	sin31_00 = TrigonometricFunctionValue{Numerator: 5150, Denominator: 10000}
	sin31_10 = TrigonometricFunctionValue{Numerator: 5175, Denominator: 10000}
	sin31_20 = TrigonometricFunctionValue{Numerator: 5200, Denominator: 10000}
	sin31_30 = TrigonometricFunctionValue{Numerator: 5225, Denominator: 10000}
	sin31_40 = TrigonometricFunctionValue{Numerator: 5250, Denominator: 10000}
	sin31_50 = TrigonometricFunctionValue{Numerator: 5275, Denominator: 10000}
	sin32_00 = TrigonometricFunctionValue{Numerator: 5299, Denominator: 10000}
	sin32_10 = TrigonometricFunctionValue{Numerator: 5324, Denominator: 10000}
	sin32_20 = TrigonometricFunctionValue{Numerator: 5348, Denominator: 10000}
	sin32_30 = TrigonometricFunctionValue{Numerator: 5373, Denominator: 10000}
	sin32_40 = TrigonometricFunctionValue{Numerator: 5398, Denominator: 10000}
	sin32_50 = TrigonometricFunctionValue{Numerator: 5422, Denominator: 10000}
	sin33_00 = TrigonometricFunctionValue{Numerator: 5446, Denominator: 10000}
	sin33_10 = TrigonometricFunctionValue{Numerator: 5471, Denominator: 10000}
	sin33_20 = TrigonometricFunctionValue{Numerator: 5495, Denominator: 10000}
	sin33_30 = TrigonometricFunctionValue{Numerator: 5519, Denominator: 10000}
	sin33_40 = TrigonometricFunctionValue{Numerator: 5544, Denominator: 10000}
	sin33_50 = TrigonometricFunctionValue{Numerator: 5568, Denominator: 10000}
	sin34_00 = TrigonometricFunctionValue{Numerator: 5592, Denominator: 10000}
	sin34_10 = TrigonometricFunctionValue{Numerator: 5616, Denominator: 10000}
	sin34_20 = TrigonometricFunctionValue{Numerator: 5640, Denominator: 10000}
	sin34_30 = TrigonometricFunctionValue{Numerator: 5664, Denominator: 10000}
	sin34_40 = TrigonometricFunctionValue{Numerator: 5688, Denominator: 10000}
	sin34_50 = TrigonometricFunctionValue{Numerator: 5712, Denominator: 10000}
	sin35_00 = TrigonometricFunctionValue{Numerator: 5736, Denominator: 10000}
	sin35_10 = TrigonometricFunctionValue{Numerator: 5760, Denominator: 10000}
	sin35_20 = TrigonometricFunctionValue{Numerator: 5783, Denominator: 10000}
	sin35_30 = TrigonometricFunctionValue{Numerator: 5807, Denominator: 10000}
	sin35_40 = TrigonometricFunctionValue{Numerator: 5831, Denominator: 10000}
	sin35_50 = TrigonometricFunctionValue{Numerator: 5854, Denominator: 10000}
	sin36_00 = TrigonometricFunctionValue{Numerator: 5878, Denominator: 10000}
	sin36_10 = TrigonometricFunctionValue{Numerator: 5901, Denominator: 10000}
	sin36_20 = TrigonometricFunctionValue{Numerator: 5925, Denominator: 10000}
	sin36_30 = TrigonometricFunctionValue{Numerator: 5948, Denominator: 10000}
	sin36_40 = TrigonometricFunctionValue{Numerator: 5972, Denominator: 10000}
	sin36_50 = TrigonometricFunctionValue{Numerator: 5995, Denominator: 10000}
	sin37_00 = TrigonometricFunctionValue{Numerator: 6018, Denominator: 10000}
	sin37_10 = TrigonometricFunctionValue{Numerator: 6041, Denominator: 10000}
	sin37_20 = TrigonometricFunctionValue{Numerator: 6065, Denominator: 10000}
	sin37_30 = TrigonometricFunctionValue{Numerator: 6088, Denominator: 10000}
	sin37_40 = TrigonometricFunctionValue{Numerator: 6111, Denominator: 10000}
	sin37_50 = TrigonometricFunctionValue{Numerator: 6134, Denominator: 10000}
	sin38_00 = TrigonometricFunctionValue{Numerator: 6157, Denominator: 10000}
	sin38_10 = TrigonometricFunctionValue{Numerator: 6180, Denominator: 10000}
	sin38_20 = TrigonometricFunctionValue{Numerator: 6202, Denominator: 10000}
	sin38_30 = TrigonometricFunctionValue{Numerator: 6225, Denominator: 10000}
	sin38_40 = TrigonometricFunctionValue{Numerator: 6248, Denominator: 10000}
	sin38_50 = TrigonometricFunctionValue{Numerator: 6271, Denominator: 10000}
	sin39_00 = TrigonometricFunctionValue{Numerator: 6293, Denominator: 10000}
	sin39_10 = TrigonometricFunctionValue{Numerator: 6316, Denominator: 10000}
	sin39_20 = TrigonometricFunctionValue{Numerator: 6338, Denominator: 10000}
	sin39_30 = TrigonometricFunctionValue{Numerator: 6361, Denominator: 10000}
	sin39_40 = TrigonometricFunctionValue{Numerator: 6383, Denominator: 10000}
	sin39_50 = TrigonometricFunctionValue{Numerator: 6406, Denominator: 10000}
	sin40_00 = TrigonometricFunctionValue{Numerator: 6428, Denominator: 10000}
	sin40_10 = TrigonometricFunctionValue{Numerator: 6450, Denominator: 10000}
	sin40_20 = TrigonometricFunctionValue{Numerator: 6472, Denominator: 10000}
	sin40_30 = TrigonometricFunctionValue{Numerator: 6494, Denominator: 10000}
	sin40_40 = TrigonometricFunctionValue{Numerator: 6517, Denominator: 10000}
	sin40_50 = TrigonometricFunctionValue{Numerator: 6539, Denominator: 10000}
	sin41_00 = TrigonometricFunctionValue{Numerator: 6561, Denominator: 10000}
	sin41_10 = TrigonometricFunctionValue{Numerator: 6583, Denominator: 10000}
	sin41_20 = TrigonometricFunctionValue{Numerator: 6604, Denominator: 10000}
	sin41_30 = TrigonometricFunctionValue{Numerator: 6626, Denominator: 10000}
	sin41_40 = TrigonometricFunctionValue{Numerator: 6648, Denominator: 10000}
	sin41_50 = TrigonometricFunctionValue{Numerator: 6670, Denominator: 10000}
	sin42_00 = TrigonometricFunctionValue{Numerator: 6691, Denominator: 10000}
	sin42_10 = TrigonometricFunctionValue{Numerator: 6713, Denominator: 10000}
	sin42_20 = TrigonometricFunctionValue{Numerator: 6734, Denominator: 10000}
	sin42_30 = TrigonometricFunctionValue{Numerator: 6756, Denominator: 10000}
	sin42_40 = TrigonometricFunctionValue{Numerator: 6777, Denominator: 10000}
	sin42_50 = TrigonometricFunctionValue{Numerator: 6799, Denominator: 10000}
	sin43_00 = TrigonometricFunctionValue{Numerator: 6820, Denominator: 10000}
	sin43_10 = TrigonometricFunctionValue{Numerator: 6841, Denominator: 10000}
	sin43_20 = TrigonometricFunctionValue{Numerator: 6862, Denominator: 10000}
	sin43_30 = TrigonometricFunctionValue{Numerator: 6884, Denominator: 10000}
	sin43_40 = TrigonometricFunctionValue{Numerator: 6905, Denominator: 10000}
	sin43_50 = TrigonometricFunctionValue{Numerator: 6926, Denominator: 10000}
	sin44_00 = TrigonometricFunctionValue{Numerator: 6947, Denominator: 10000}
	sin44_10 = TrigonometricFunctionValue{Numerator: 6967, Denominator: 10000}
	sin44_20 = TrigonometricFunctionValue{Numerator: 6988, Denominator: 10000}
	sin44_30 = TrigonometricFunctionValue{Numerator: 7009, Denominator: 10000}
	sin44_40 = TrigonometricFunctionValue{Numerator: 7030, Denominator: 10000}
	sin44_50 = TrigonometricFunctionValue{Numerator: 7050, Denominator: 10000}
	sin45_00 = TrigonometricFunctionValue{Numerator: 7071, Denominator: 10000}
	sin45_10 = TrigonometricFunctionValue{Numerator: 7092, Denominator: 10000}
	sin45_20 = TrigonometricFunctionValue{Numerator: 7112, Denominator: 10000}
	sin45_30 = TrigonometricFunctionValue{Numerator: 7133, Denominator: 10000}
	sin45_40 = TrigonometricFunctionValue{Numerator: 7153, Denominator: 10000}
	sin45_50 = TrigonometricFunctionValue{Numerator: 7173, Denominator: 10000}
	sin46_00 = TrigonometricFunctionValue{Numerator: 7193, Denominator: 10000}
	sin46_10 = TrigonometricFunctionValue{Numerator: 7214, Denominator: 10000}
	sin46_20 = TrigonometricFunctionValue{Numerator: 7234, Denominator: 10000}
	sin46_30 = TrigonometricFunctionValue{Numerator: 7254, Denominator: 10000}
	sin46_40 = TrigonometricFunctionValue{Numerator: 7274, Denominator: 10000}
	sin46_50 = TrigonometricFunctionValue{Numerator: 7294, Denominator: 10000}
	sin47_00 = TrigonometricFunctionValue{Numerator: 7314, Denominator: 10000}
	sin47_10 = TrigonometricFunctionValue{Numerator: 7333, Denominator: 10000}
	sin47_20 = TrigonometricFunctionValue{Numerator: 7353, Denominator: 10000}
	sin47_30 = TrigonometricFunctionValue{Numerator: 7373, Denominator: 10000}
	sin47_40 = TrigonometricFunctionValue{Numerator: 7392, Denominator: 10000}
	sin47_50 = TrigonometricFunctionValue{Numerator: 7412, Denominator: 10000}
	sin48_00 = TrigonometricFunctionValue{Numerator: 7431, Denominator: 10000}
	sin48_10 = TrigonometricFunctionValue{Numerator: 7451, Denominator: 10000}
	sin48_20 = TrigonometricFunctionValue{Numerator: 7470, Denominator: 10000}
	sin48_30 = TrigonometricFunctionValue{Numerator: 7490, Denominator: 10000}
	sin48_40 = TrigonometricFunctionValue{Numerator: 7509, Denominator: 10000}
	sin48_50 = TrigonometricFunctionValue{Numerator: 7528, Denominator: 10000}
	sin49_00 = TrigonometricFunctionValue{Numerator: 7547, Denominator: 10000}
	sin49_10 = TrigonometricFunctionValue{Numerator: 7566, Denominator: 10000}
	sin49_20 = TrigonometricFunctionValue{Numerator: 7585, Denominator: 10000}
	sin49_30 = TrigonometricFunctionValue{Numerator: 7604, Denominator: 10000}
	sin49_40 = TrigonometricFunctionValue{Numerator: 7623, Denominator: 10000}
	sin49_50 = TrigonometricFunctionValue{Numerator: 7642, Denominator: 10000}
	sin50_00 = TrigonometricFunctionValue{Numerator: 7660, Denominator: 10000}
	sin50_10 = TrigonometricFunctionValue{Numerator: 7679, Denominator: 10000}
	sin50_20 = TrigonometricFunctionValue{Numerator: 7698, Denominator: 10000}
	sin50_30 = TrigonometricFunctionValue{Numerator: 7716, Denominator: 10000}
	sin50_40 = TrigonometricFunctionValue{Numerator: 7735, Denominator: 10000}
	sin50_50 = TrigonometricFunctionValue{Numerator: 7753, Denominator: 10000}
	sin51_00 = TrigonometricFunctionValue{Numerator: 7771, Denominator: 10000}
	sin51_10 = TrigonometricFunctionValue{Numerator: 7790, Denominator: 10000}
	sin51_20 = TrigonometricFunctionValue{Numerator: 7808, Denominator: 10000}
	sin51_30 = TrigonometricFunctionValue{Numerator: 7826, Denominator: 10000}
	sin51_40 = TrigonometricFunctionValue{Numerator: 7844, Denominator: 10000}
	sin51_50 = TrigonometricFunctionValue{Numerator: 7862, Denominator: 10000}
	sin52_00 = TrigonometricFunctionValue{Numerator: 7880, Denominator: 10000}
	sin52_10 = TrigonometricFunctionValue{Numerator: 7898, Denominator: 10000}
	sin52_20 = TrigonometricFunctionValue{Numerator: 7916, Denominator: 10000}
	sin52_30 = TrigonometricFunctionValue{Numerator: 7934, Denominator: 10000}
	sin52_40 = TrigonometricFunctionValue{Numerator: 7951, Denominator: 10000}
	sin52_50 = TrigonometricFunctionValue{Numerator: 7969, Denominator: 10000}
	sin53_00 = TrigonometricFunctionValue{Numerator: 7986, Denominator: 10000}
	sin53_10 = TrigonometricFunctionValue{Numerator: 8004, Denominator: 10000}
	sin53_20 = TrigonometricFunctionValue{Numerator: 8021, Denominator: 10000}
	sin53_30 = TrigonometricFunctionValue{Numerator: 8039, Denominator: 10000}
	sin53_40 = TrigonometricFunctionValue{Numerator: 8056, Denominator: 10000}
	sin53_50 = TrigonometricFunctionValue{Numerator: 8073, Denominator: 10000}
	sin54_00 = TrigonometricFunctionValue{Numerator: 8090, Denominator: 10000}
	sin54_10 = TrigonometricFunctionValue{Numerator: 8107, Denominator: 10000}
	sin54_20 = TrigonometricFunctionValue{Numerator: 8124, Denominator: 10000}
	sin54_30 = TrigonometricFunctionValue{Numerator: 8141, Denominator: 10000}
	sin54_40 = TrigonometricFunctionValue{Numerator: 8158, Denominator: 10000}
	sin54_50 = TrigonometricFunctionValue{Numerator: 8175, Denominator: 10000}
	sin55_00 = TrigonometricFunctionValue{Numerator: 8192, Denominator: 10000}
	sin55_10 = TrigonometricFunctionValue{Numerator: 8208, Denominator: 10000}
	sin55_20 = TrigonometricFunctionValue{Numerator: 8225, Denominator: 10000}
	sin55_30 = TrigonometricFunctionValue{Numerator: 8241, Denominator: 10000}
	sin55_40 = TrigonometricFunctionValue{Numerator: 8258, Denominator: 10000}
	sin55_50 = TrigonometricFunctionValue{Numerator: 8274, Denominator: 10000}
	sin56_00 = TrigonometricFunctionValue{Numerator: 8290, Denominator: 10000}
	sin56_10 = TrigonometricFunctionValue{Numerator: 8307, Denominator: 10000}
	sin56_20 = TrigonometricFunctionValue{Numerator: 8323, Denominator: 10000}
	sin56_30 = TrigonometricFunctionValue{Numerator: 8339, Denominator: 10000}
	sin56_40 = TrigonometricFunctionValue{Numerator: 8355, Denominator: 10000}
	sin56_50 = TrigonometricFunctionValue{Numerator: 8371, Denominator: 10000}
	sin57_00 = TrigonometricFunctionValue{Numerator: 8387, Denominator: 10000}
	sin57_10 = TrigonometricFunctionValue{Numerator: 8403, Denominator: 10000}
	sin57_20 = TrigonometricFunctionValue{Numerator: 8418, Denominator: 10000}
	sin57_30 = TrigonometricFunctionValue{Numerator: 8434, Denominator: 10000}
	sin57_40 = TrigonometricFunctionValue{Numerator: 8450, Denominator: 10000}
	sin57_50 = TrigonometricFunctionValue{Numerator: 8465, Denominator: 10000}
	sin58_00 = TrigonometricFunctionValue{Numerator: 8480, Denominator: 10000}
	sin58_10 = TrigonometricFunctionValue{Numerator: 8496, Denominator: 10000}
	sin58_20 = TrigonometricFunctionValue{Numerator: 8511, Denominator: 10000}
	sin58_30 = TrigonometricFunctionValue{Numerator: 8526, Denominator: 10000}
	sin58_40 = TrigonometricFunctionValue{Numerator: 8542, Denominator: 10000}
	sin58_50 = TrigonometricFunctionValue{Numerator: 8557, Denominator: 10000}
	sin59_00 = TrigonometricFunctionValue{Numerator: 8572, Denominator: 10000}
	sin59_10 = TrigonometricFunctionValue{Numerator: 8587, Denominator: 10000}
	sin59_20 = TrigonometricFunctionValue{Numerator: 8601, Denominator: 10000}
	sin59_30 = TrigonometricFunctionValue{Numerator: 8616, Denominator: 10000}
	sin59_40 = TrigonometricFunctionValue{Numerator: 8631, Denominator: 10000}
	sin59_50 = TrigonometricFunctionValue{Numerator: 8646, Denominator: 10000}
	sin60_00 = TrigonometricFunctionValue{Numerator: 8660, Denominator: 10000}
	sin60_10 = TrigonometricFunctionValue{Numerator: 8675, Denominator: 10000}
	sin60_20 = TrigonometricFunctionValue{Numerator: 8689, Denominator: 10000}
	sin60_30 = TrigonometricFunctionValue{Numerator: 8704, Denominator: 10000}
	sin60_40 = TrigonometricFunctionValue{Numerator: 8718, Denominator: 10000}
	sin60_50 = TrigonometricFunctionValue{Numerator: 8732, Denominator: 10000}
	sin61_00 = TrigonometricFunctionValue{Numerator: 8746, Denominator: 10000}
	sin61_10 = TrigonometricFunctionValue{Numerator: 8760, Denominator: 10000}
	sin61_20 = TrigonometricFunctionValue{Numerator: 8774, Denominator: 10000}
	sin61_30 = TrigonometricFunctionValue{Numerator: 8788, Denominator: 10000}
	sin61_40 = TrigonometricFunctionValue{Numerator: 8802, Denominator: 10000}
	sin61_50 = TrigonometricFunctionValue{Numerator: 8816, Denominator: 10000}
	sin62_00 = TrigonometricFunctionValue{Numerator: 8829, Denominator: 10000}
	sin62_10 = TrigonometricFunctionValue{Numerator: 8843, Denominator: 10000}
	sin62_20 = TrigonometricFunctionValue{Numerator: 8857, Denominator: 10000}
	sin62_30 = TrigonometricFunctionValue{Numerator: 8870, Denominator: 10000}
	sin62_40 = TrigonometricFunctionValue{Numerator: 8884, Denominator: 10000}
	sin62_50 = TrigonometricFunctionValue{Numerator: 8897, Denominator: 10000}
	sin63_00 = TrigonometricFunctionValue{Numerator: 8910, Denominator: 10000}
	sin63_10 = TrigonometricFunctionValue{Numerator: 8923, Denominator: 10000}
	sin63_20 = TrigonometricFunctionValue{Numerator: 8936, Denominator: 10000}
	sin63_30 = TrigonometricFunctionValue{Numerator: 8949, Denominator: 10000}
	sin63_40 = TrigonometricFunctionValue{Numerator: 8962, Denominator: 10000}
	sin63_50 = TrigonometricFunctionValue{Numerator: 8975, Denominator: 10000}
	sin64_00 = TrigonometricFunctionValue{Numerator: 8988, Denominator: 10000}
	sin64_10 = TrigonometricFunctionValue{Numerator: 9001, Denominator: 10000}
	sin64_20 = TrigonometricFunctionValue{Numerator: 9013, Denominator: 10000}
	sin64_30 = TrigonometricFunctionValue{Numerator: 9026, Denominator: 10000}
	sin64_40 = TrigonometricFunctionValue{Numerator: 9038, Denominator: 10000}
	sin64_50 = TrigonometricFunctionValue{Numerator: 9051, Denominator: 10000}
	sin65_00 = TrigonometricFunctionValue{Numerator: 9063, Denominator: 10000}
	sin65_10 = TrigonometricFunctionValue{Numerator: 9075, Denominator: 10000}
	sin65_20 = TrigonometricFunctionValue{Numerator: 9088, Denominator: 10000}
	sin65_30 = TrigonometricFunctionValue{Numerator: 9100, Denominator: 10000}
	sin65_40 = TrigonometricFunctionValue{Numerator: 9112, Denominator: 10000}
	sin65_50 = TrigonometricFunctionValue{Numerator: 9124, Denominator: 10000}
	sin66_00 = TrigonometricFunctionValue{Numerator: 9135, Denominator: 10000}
	sin66_10 = TrigonometricFunctionValue{Numerator: 9147, Denominator: 10000}
	sin66_20 = TrigonometricFunctionValue{Numerator: 9159, Denominator: 10000}
	sin66_30 = TrigonometricFunctionValue{Numerator: 9171, Denominator: 10000}
	sin66_40 = TrigonometricFunctionValue{Numerator: 9182, Denominator: 10000}
	sin66_50 = TrigonometricFunctionValue{Numerator: 9194, Denominator: 10000}
	sin67_00 = TrigonometricFunctionValue{Numerator: 9205, Denominator: 10000}
	sin67_10 = TrigonometricFunctionValue{Numerator: 9216, Denominator: 10000}
	sin67_20 = TrigonometricFunctionValue{Numerator: 9228, Denominator: 10000}
	sin67_30 = TrigonometricFunctionValue{Numerator: 9239, Denominator: 10000}
	sin67_40 = TrigonometricFunctionValue{Numerator: 9250, Denominator: 10000}
	sin67_50 = TrigonometricFunctionValue{Numerator: 9261, Denominator: 10000}
	sin68_00 = TrigonometricFunctionValue{Numerator: 9272, Denominator: 10000}
	sin68_10 = TrigonometricFunctionValue{Numerator: 9283, Denominator: 10000}
	sin68_20 = TrigonometricFunctionValue{Numerator: 9293, Denominator: 10000}
	sin68_30 = TrigonometricFunctionValue{Numerator: 9304, Denominator: 10000}
	sin68_40 = TrigonometricFunctionValue{Numerator: 9315, Denominator: 10000}
	sin68_50 = TrigonometricFunctionValue{Numerator: 9325, Denominator: 10000}
	sin69_00 = TrigonometricFunctionValue{Numerator: 9336, Denominator: 10000}
	sin69_10 = TrigonometricFunctionValue{Numerator: 9346, Denominator: 10000}
	sin69_20 = TrigonometricFunctionValue{Numerator: 9356, Denominator: 10000}
	sin69_30 = TrigonometricFunctionValue{Numerator: 9367, Denominator: 10000}
	sin69_40 = TrigonometricFunctionValue{Numerator: 9377, Denominator: 10000}
	sin69_50 = TrigonometricFunctionValue{Numerator: 9387, Denominator: 10000}
	sin70_00 = TrigonometricFunctionValue{Numerator: 9397, Denominator: 10000}
	sin70_10 = TrigonometricFunctionValue{Numerator: 9407, Denominator: 10000}
	sin70_20 = TrigonometricFunctionValue{Numerator: 9417, Denominator: 10000}
	sin70_30 = TrigonometricFunctionValue{Numerator: 9426, Denominator: 10000}
	sin70_40 = TrigonometricFunctionValue{Numerator: 9436, Denominator: 10000}
	sin70_50 = TrigonometricFunctionValue{Numerator: 9446, Denominator: 10000}
	sin71_00 = TrigonometricFunctionValue{Numerator: 9455, Denominator: 10000}
	sin71_10 = TrigonometricFunctionValue{Numerator: 9465, Denominator: 10000}
	sin71_20 = TrigonometricFunctionValue{Numerator: 9474, Denominator: 10000}
	sin71_30 = TrigonometricFunctionValue{Numerator: 9483, Denominator: 10000}
	sin71_40 = TrigonometricFunctionValue{Numerator: 9492, Denominator: 10000}
	sin71_50 = TrigonometricFunctionValue{Numerator: 9502, Denominator: 10000}
	sin72_00 = TrigonometricFunctionValue{Numerator: 9511, Denominator: 10000}
	sin72_10 = TrigonometricFunctionValue{Numerator: 9520, Denominator: 10000}
	sin72_20 = TrigonometricFunctionValue{Numerator: 9528, Denominator: 10000}
	sin72_30 = TrigonometricFunctionValue{Numerator: 9537, Denominator: 10000}
	sin72_40 = TrigonometricFunctionValue{Numerator: 9546, Denominator: 10000}
	sin72_50 = TrigonometricFunctionValue{Numerator: 9555, Denominator: 10000}
	sin73_00 = TrigonometricFunctionValue{Numerator: 9563, Denominator: 10000}
	sin73_10 = TrigonometricFunctionValue{Numerator: 9572, Denominator: 10000}
	sin73_20 = TrigonometricFunctionValue{Numerator: 9580, Denominator: 10000}
	sin73_30 = TrigonometricFunctionValue{Numerator: 9588, Denominator: 10000}
	sin73_40 = TrigonometricFunctionValue{Numerator: 9596, Denominator: 10000}
	sin73_50 = TrigonometricFunctionValue{Numerator: 9605, Denominator: 10000}
	sin74_00 = TrigonometricFunctionValue{Numerator: 9613, Denominator: 10000}
	sin74_10 = TrigonometricFunctionValue{Numerator: 9621, Denominator: 10000}
	sin74_20 = TrigonometricFunctionValue{Numerator: 9628, Denominator: 10000}
	sin74_30 = TrigonometricFunctionValue{Numerator: 9636, Denominator: 10000}
	sin74_40 = TrigonometricFunctionValue{Numerator: 9644, Denominator: 10000}
	sin74_50 = TrigonometricFunctionValue{Numerator: 9652, Denominator: 10000}
	sin75_00 = TrigonometricFunctionValue{Numerator: 9659, Denominator: 10000}
	sin75_10 = TrigonometricFunctionValue{Numerator: 9667, Denominator: 10000}
	sin75_20 = TrigonometricFunctionValue{Numerator: 9674, Denominator: 10000}
	sin75_30 = TrigonometricFunctionValue{Numerator: 9681, Denominator: 10000}
	sin75_40 = TrigonometricFunctionValue{Numerator: 9689, Denominator: 10000}
	sin75_50 = TrigonometricFunctionValue{Numerator: 9696, Denominator: 10000}
	sin76_00 = TrigonometricFunctionValue{Numerator: 9703, Denominator: 10000}
	sin76_10 = TrigonometricFunctionValue{Numerator: 9710, Denominator: 10000}
	sin76_20 = TrigonometricFunctionValue{Numerator: 9717, Denominator: 10000}
	sin76_30 = TrigonometricFunctionValue{Numerator: 9724, Denominator: 10000}
	sin76_40 = TrigonometricFunctionValue{Numerator: 9730, Denominator: 10000}
	sin76_50 = TrigonometricFunctionValue{Numerator: 9737, Denominator: 10000}
	sin77_00 = TrigonometricFunctionValue{Numerator: 9744, Denominator: 10000}
	sin77_10 = TrigonometricFunctionValue{Numerator: 9750, Denominator: 10000}
	sin77_20 = TrigonometricFunctionValue{Numerator: 9757, Denominator: 10000}
	sin77_30 = TrigonometricFunctionValue{Numerator: 9763, Denominator: 10000}
	sin77_40 = TrigonometricFunctionValue{Numerator: 9769, Denominator: 10000}
	sin77_50 = TrigonometricFunctionValue{Numerator: 9775, Denominator: 10000}
	sin78_00 = TrigonometricFunctionValue{Numerator: 9781, Denominator: 10000}
	sin78_10 = TrigonometricFunctionValue{Numerator: 9787, Denominator: 10000}
	sin78_20 = TrigonometricFunctionValue{Numerator: 9793, Denominator: 10000}
	sin78_30 = TrigonometricFunctionValue{Numerator: 9799, Denominator: 10000}
	sin78_40 = TrigonometricFunctionValue{Numerator: 9805, Denominator: 10000}
	sin78_50 = TrigonometricFunctionValue{Numerator: 9811, Denominator: 10000}
	sin79_00 = TrigonometricFunctionValue{Numerator: 9816, Denominator: 10000}
	sin79_10 = TrigonometricFunctionValue{Numerator: 9822, Denominator: 10000}
	sin79_20 = TrigonometricFunctionValue{Numerator: 9827, Denominator: 10000}
	sin79_30 = TrigonometricFunctionValue{Numerator: 9833, Denominator: 10000}
	sin79_40 = TrigonometricFunctionValue{Numerator: 9838, Denominator: 10000}
	sin79_50 = TrigonometricFunctionValue{Numerator: 9843, Denominator: 10000}
	sin80_00 = TrigonometricFunctionValue{Numerator: 9848, Denominator: 10000}
	sin80_10 = TrigonometricFunctionValue{Numerator: 9853, Denominator: 10000}
	sin80_20 = TrigonometricFunctionValue{Numerator: 9858, Denominator: 10000}
	sin80_30 = TrigonometricFunctionValue{Numerator: 9863, Denominator: 10000}
	sin80_40 = TrigonometricFunctionValue{Numerator: 9868, Denominator: 10000}
	sin80_50 = TrigonometricFunctionValue{Numerator: 9872, Denominator: 10000}
	sin81_00 = TrigonometricFunctionValue{Numerator: 9877, Denominator: 10000}
	sin81_10 = TrigonometricFunctionValue{Numerator: 9881, Denominator: 10000}
	sin81_20 = TrigonometricFunctionValue{Numerator: 9886, Denominator: 10000}
	sin81_30 = TrigonometricFunctionValue{Numerator: 9890, Denominator: 10000}
	sin81_40 = TrigonometricFunctionValue{Numerator: 9894, Denominator: 10000}
	sin81_50 = TrigonometricFunctionValue{Numerator: 9899, Denominator: 10000}
	sin82_00 = TrigonometricFunctionValue{Numerator: 9903, Denominator: 10000}
	sin82_10 = TrigonometricFunctionValue{Numerator: 9907, Denominator: 10000}
	sin82_20 = TrigonometricFunctionValue{Numerator: 9911, Denominator: 10000}
	sin82_30 = TrigonometricFunctionValue{Numerator: 9914, Denominator: 10000}
	sin82_40 = TrigonometricFunctionValue{Numerator: 9918, Denominator: 10000}
	sin82_50 = TrigonometricFunctionValue{Numerator: 9922, Denominator: 10000}
	sin83_00 = TrigonometricFunctionValue{Numerator: 9925, Denominator: 10000}
	sin83_10 = TrigonometricFunctionValue{Numerator: 9929, Denominator: 10000}
	sin83_20 = TrigonometricFunctionValue{Numerator: 9932, Denominator: 10000}
	sin83_30 = TrigonometricFunctionValue{Numerator: 9936, Denominator: 10000}
	sin83_40 = TrigonometricFunctionValue{Numerator: 9939, Denominator: 10000}
	sin83_50 = TrigonometricFunctionValue{Numerator: 9942, Denominator: 10000}
	sin84_00 = TrigonometricFunctionValue{Numerator: 9945, Denominator: 10000}
	sin84_10 = TrigonometricFunctionValue{Numerator: 9948, Denominator: 10000}
	sin84_20 = TrigonometricFunctionValue{Numerator: 9951, Denominator: 10000}
	sin84_30 = TrigonometricFunctionValue{Numerator: 9954, Denominator: 10000}
	sin84_40 = TrigonometricFunctionValue{Numerator: 9957, Denominator: 10000}
	sin84_50 = TrigonometricFunctionValue{Numerator: 9959, Denominator: 10000}
	sin85_00 = TrigonometricFunctionValue{Numerator: 9962, Denominator: 10000}
	sin85_10 = TrigonometricFunctionValue{Numerator: 9964, Denominator: 10000}
	sin85_20 = TrigonometricFunctionValue{Numerator: 9967, Denominator: 10000}
	sin85_30 = TrigonometricFunctionValue{Numerator: 9969, Denominator: 10000}
	sin85_40 = TrigonometricFunctionValue{Numerator: 9971, Denominator: 10000}
	sin85_50 = TrigonometricFunctionValue{Numerator: 9974, Denominator: 10000}
	sin86_00 = TrigonometricFunctionValue{Numerator: 9976, Denominator: 10000}
	sin86_10 = TrigonometricFunctionValue{Numerator: 9978, Denominator: 10000}
	sin86_20 = TrigonometricFunctionValue{Numerator: 9980, Denominator: 10000}
	sin86_30 = TrigonometricFunctionValue{Numerator: 9981, Denominator: 10000}
	sin86_40 = TrigonometricFunctionValue{Numerator: 9983, Denominator: 10000}
	sin86_50 = TrigonometricFunctionValue{Numerator: 9985, Denominator: 10000}
	sin87_00 = TrigonometricFunctionValue{Numerator: 9986, Denominator: 10000}
	sin87_10 = TrigonometricFunctionValue{Numerator: 9988, Denominator: 10000}
	sin87_20 = TrigonometricFunctionValue{Numerator: 9989, Denominator: 10000}
	sin87_30 = TrigonometricFunctionValue{Numerator: 9990, Denominator: 10000}
	sin87_40 = TrigonometricFunctionValue{Numerator: 9992, Denominator: 10000}
	sin87_50 = TrigonometricFunctionValue{Numerator: 9993, Denominator: 10000}
	sin88_00 = TrigonometricFunctionValue{Numerator: 9994, Denominator: 10000}
	sin88_10 = TrigonometricFunctionValue{Numerator: 9995, Denominator: 10000}
	sin88_20 = TrigonometricFunctionValue{Numerator: 9996, Denominator: 10000}
	sin88_30 = TrigonometricFunctionValue{Numerator: 9997, Denominator: 10000}
	sin88_40 = TrigonometricFunctionValue{Numerator: 9997, Denominator: 10000}
	sin88_50 = TrigonometricFunctionValue{Numerator: 9998, Denominator: 10000}
	sin89_00 = TrigonometricFunctionValue{Numerator: 9998, Denominator: 10000}
	sin89_10 = TrigonometricFunctionValue{Numerator: 9999, Denominator: 10000}
	sin89_20 = TrigonometricFunctionValue{Numerator: 9999, Denominator: 10000}
	sin89_30 = TrigonometricFunctionValue{Numerator: 10000, Denominator: 10000}
	sin89_40 = TrigonometricFunctionValue{Numerator: 10000, Denominator: 10000}
	sin89_50 = TrigonometricFunctionValue{Numerator: 10000, Denominator: 10000}
	sin90_00 = TrigonometricFunctionValue{Numerator: 10000, Denominator: 10000}
)

var sinval = [][6]TrigonometricFunctionValue{
	{sin0_00, sin0_10, sin0_20, sin0_30, sin0_40, sin0_50},
	{sin1_00, sin1_10, sin1_20, sin1_30, sin1_40, sin1_50},
	{sin2_00, sin2_10, sin2_20, sin2_30, sin2_40, sin2_50},
	{sin3_00, sin3_10, sin3_20, sin3_30, sin3_40, sin3_50},
	{sin4_00, sin4_10, sin4_20, sin4_30, sin4_40, sin4_50},
	{sin5_00, sin5_10, sin5_20, sin5_30, sin5_40, sin5_50},
	{sin6_00, sin6_10, sin6_20, sin6_30, sin6_40, sin6_50},
	{sin7_00, sin7_10, sin7_20, sin7_30, sin7_40, sin7_50},
	{sin8_00, sin8_10, sin8_20, sin8_30, sin8_40, sin8_50},
	{sin9_00, sin9_10, sin9_20, sin9_30, sin9_40, sin9_50},
	{sin10_00, sin10_10, sin10_20, sin10_30, sin10_40, sin10_50}, // 10
	{sin11_00, sin11_10, sin11_20, sin11_30, sin11_40, sin11_50},
	{sin12_00, sin12_10, sin12_20, sin12_30, sin12_40, sin12_50},
	{sin13_00, sin13_10, sin13_20, sin13_30, sin13_40, sin13_50},
	{sin14_00, sin14_10, sin14_20, sin14_30, sin14_40, sin14_50},
	{sin15_00, sin15_10, sin15_20, sin15_30, sin15_40, sin15_50},
	{sin16_00, sin16_10, sin16_20, sin16_30, sin16_40, sin16_50},
	{sin17_00, sin17_10, sin17_20, sin17_30, sin17_40, sin17_50},
	{sin18_00, sin18_10, sin18_20, sin18_30, sin18_40, sin18_50},
	{sin19_00, sin19_10, sin19_20, sin19_30, sin19_40, sin19_50},
	{sin20_00, sin20_10, sin20_20, sin20_30, sin20_40, sin20_50}, // 20
	{sin21_00, sin21_10, sin21_20, sin21_30, sin21_40, sin21_50},
	{sin22_00, sin22_10, sin22_20, sin22_30, sin22_40, sin22_50},
	{sin23_00, sin23_10, sin23_20, sin23_30, sin23_40, sin23_50},
	{sin24_00, sin24_10, sin24_20, sin24_30, sin24_40, sin24_50},
	{sin25_00, sin25_10, sin25_20, sin25_30, sin25_40, sin25_50},
	{sin26_00, sin26_10, sin26_20, sin26_30, sin26_40, sin26_50},
	{sin27_00, sin27_10, sin27_20, sin27_30, sin27_40, sin27_50},
	{sin28_00, sin28_10, sin28_20, sin28_30, sin28_40, sin28_50},
	{sin29_00, sin29_10, sin29_20, sin29_30, sin29_40, sin29_50},
	{sin30_00, sin30_10, sin30_20, sin30_30, sin30_40, sin30_50}, // 30
	{sin31_00, sin31_10, sin31_20, sin31_30, sin31_40, sin31_50},
	{sin32_00, sin32_10, sin32_20, sin32_30, sin32_40, sin32_50},
	{sin33_00, sin33_10, sin33_20, sin33_30, sin33_40, sin33_50},
	{sin34_00, sin34_10, sin34_20, sin34_30, sin34_40, sin34_50},
	{sin35_00, sin35_10, sin35_20, sin35_30, sin35_40, sin35_50},
	{sin36_00, sin36_10, sin36_20, sin36_30, sin36_40, sin36_50},
	{sin37_00, sin37_10, sin37_20, sin37_30, sin37_40, sin37_50},
	{sin38_00, sin38_10, sin38_20, sin38_30, sin38_40, sin38_50},
	{sin39_00, sin39_10, sin39_20, sin39_30, sin39_40, sin39_50},
	{sin40_00, sin40_10, sin40_20, sin40_30, sin40_40, sin40_50}, // 40
	{sin41_00, sin41_10, sin41_20, sin41_30, sin41_40, sin41_50},
	{sin42_00, sin42_10, sin42_20, sin42_30, sin42_40, sin42_50},
	{sin43_00, sin43_10, sin43_20, sin43_30, sin43_40, sin43_50},
	{sin44_00, sin44_10, sin44_20, sin44_30, sin44_40, sin44_50},
	{sin45_00, sin45_10, sin45_20, sin45_30, sin45_40, sin45_50},
	{sin46_00, sin46_10, sin46_20, sin46_30, sin46_40, sin46_50},
	{sin47_00, sin47_10, sin47_20, sin47_30, sin47_40, sin47_50},
	{sin48_00, sin48_10, sin48_20, sin48_30, sin48_40, sin48_50},
	{sin49_00, sin49_10, sin49_20, sin49_30, sin49_40, sin49_50},
	{sin50_00, sin50_10, sin50_20, sin50_30, sin50_40, sin50_50}, // 50
	{sin51_00, sin51_10, sin51_20, sin51_30, sin51_40, sin51_50},
	{sin52_00, sin52_10, sin52_20, sin52_30, sin52_40, sin52_50},
	{sin53_00, sin53_10, sin53_20, sin53_30, sin53_40, sin53_50},
	{sin54_00, sin54_10, sin54_20, sin54_30, sin54_40, sin54_50},
	{sin55_00, sin55_10, sin55_20, sin55_30, sin55_40, sin55_50},
	{sin56_00, sin56_10, sin56_20, sin56_30, sin56_40, sin56_50},
	{sin57_00, sin57_10, sin57_20, sin57_30, sin57_40, sin57_50},
	{sin58_00, sin58_10, sin58_20, sin58_30, sin58_40, sin58_50},
	{sin59_00, sin59_10, sin59_20, sin59_30, sin59_40, sin59_50},
	{sin60_00, sin60_10, sin60_20, sin60_30, sin60_40, sin60_50}, // 60
	{sin61_00, sin61_10, sin61_20, sin61_30, sin61_40, sin61_50},
	{sin62_00, sin62_10, sin62_20, sin62_30, sin62_40, sin62_50},
	{sin63_00, sin63_10, sin63_20, sin63_30, sin63_40, sin63_50},
	{sin64_00, sin64_10, sin64_20, sin64_30, sin64_40, sin64_50},
	{sin65_00, sin65_10, sin65_20, sin65_30, sin65_40, sin65_50},
	{sin66_00, sin66_10, sin66_20, sin66_30, sin66_40, sin66_50},
	{sin67_00, sin67_10, sin67_20, sin67_30, sin67_40, sin67_50},
	{sin68_00, sin68_10, sin68_20, sin68_30, sin68_40, sin68_50},
	{sin69_00, sin69_10, sin69_20, sin69_30, sin69_40, sin69_50},
	{sin70_00, sin70_10, sin70_20, sin70_30, sin70_40, sin70_50}, // 70
	{sin71_00, sin71_10, sin71_20, sin71_30, sin71_40, sin71_50},
	{sin72_00, sin72_10, sin72_20, sin72_30, sin72_40, sin72_50},
	{sin73_00, sin73_10, sin73_20, sin73_30, sin73_40, sin73_50},
	{sin74_00, sin74_10, sin74_20, sin74_30, sin74_40, sin74_50},
	{sin75_00, sin75_10, sin75_20, sin75_30, sin75_40, sin75_50},
	{sin76_00, sin76_10, sin76_20, sin76_30, sin76_40, sin76_50},
	{sin77_00, sin77_10, sin77_20, sin77_30, sin77_40, sin77_50},
	{sin78_00, sin78_10, sin78_20, sin78_30, sin78_40, sin78_50},
	{sin79_00, sin79_10, sin79_20, sin79_30, sin79_40, sin79_50},
	{sin80_00, sin80_10, sin80_20, sin80_30, sin80_40, sin80_50}, // 80
	{sin81_00, sin81_10, sin81_20, sin81_30, sin81_40, sin81_50},
	{sin82_00, sin82_10, sin82_20, sin82_30, sin82_40, sin82_50},
	{sin83_00, sin83_10, sin83_20, sin83_30, sin83_40, sin83_50},
	{sin84_00, sin84_10, sin84_20, sin84_30, sin84_40, sin84_50},
	{sin85_00, sin85_10, sin85_20, sin85_30, sin85_40, sin85_50},
	{sin86_00, sin86_10, sin86_20, sin86_30, sin86_40, sin86_50},
	{sin87_00, sin87_10, sin87_20, sin87_30, sin87_40, sin87_50},
	{sin88_00, sin88_10, sin88_20, sin88_30, sin88_40, sin88_50},
	{sin89_00, sin89_10, sin89_20, sin89_30, sin89_40, sin89_50},
	{sin90_00, TrigonometricFunctionValue{}, TrigonometricFunctionValue{}, TrigonometricFunctionValue{}, TrigonometricFunctionValue{}},
}

func Sine(angle Angle) TrigonometricFunctionValue {
	if angle.minute >= 60 {
		panic(fmt.Sprintf("base: invalid minute param %v for sin value", angle.minute))
	}
	if int32(angle.degree)*int32(angle.minute) < 0 {
		panic(fmt.Sprintf("base: invalid degree param %v and minute param %v", angle.degree, angle.minute))
	}
	if angle.degree < 0 {
		return Sine(angle.Negative()).Negative()
	}
	if angle.degree >= 360 {
		angle.degree %= 360
	}
	if angle.degree >= 90 && angle.degree < 180 {
		angle.degree -= 90
		return Cosine(angle)
	} else if angle.degree >= 180 && angle.degree < 270 {
		angle.degree -= 180
		return Sine(angle).Negative()
	} else if angle.degree >= 270 && angle.degree < 360 {
		angle.degree -= 180
		return Sine(angle).Negative()
	}
	return sinval[angle.degree][angle.minute/10]
}

func Cosine(angle Angle) TrigonometricFunctionValue {
	if angle.minute >= 60 {
		panic(fmt.Sprintf("base: invalid minute param %v for cos value", angle.minute))
	}
	if int32(angle.degree)*int32(angle.minute) < 0 {
		panic(fmt.Sprintf("base: invalid degree param %v and minute param %v", angle.degree, angle.minute))
	}
	if angle.degree < 0 {
		angle.degree = -angle.degree
		angle.minute = -angle.minute
	}
	if angle.degree >= 360 {
		angle.degree %= 360
	}
	if angle.degree >= 90 && angle.degree < 180 {
		angle.degree -= 90
		return Sine(angle).Negative()
	} else if angle.degree >= 180 && angle.degree < 270 {
		angle.degree -= 180
		return Cosine(angle).Negative()
	} else if angle.degree >= 270 && angle.degree < 360 {
		angle.degree -= 180
		return Cosine(angle).Negative()
	}

	dm := 90*60 - (angle.degree*60 + angle.minute)
	angle.degree, angle.minute = dm/60, dm%60
	return sinval[angle.degree][angle.minute/10]
}
