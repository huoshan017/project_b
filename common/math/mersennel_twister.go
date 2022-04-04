package math

import "math"

// Quick & dirty Mersenne Twister [MT19937] implementation
type MersennelTwister struct {
	mt         [624]uint32
	index      int32
	last       int32
	totalCount int32
}

func NewMersennelTwister(seed int32) *MersennelTwister {
	m := &MersennelTwister{}
	m.mt[0] = uint32(seed)
	l := uint32(len(m.mt))
	for i := uint32(1); i < l; i++ {
		m.mt[i] = 1812433253*(m.mt[i-1]^(m.mt[i-1]>>30)) + i
	}
	return m
}

func (m *MersennelTwister) Next() int32 {
	if m.index == 0 {
		m.Generate()
	}
	var y = m.mt[m.index]
	y ^= y >> 11
	y ^= (y << 7) & 2636928640
	y ^= (y << 15) & 4022730752
	y ^= y >> 18

	m.index = (m.index + 1) % 624
	m.totalCount += 1
	m.last = int32(y % math.MaxInt32)
	return m.last
}

func (m *MersennelTwister) Next2(low, high int32) int32 {
	if high < low {
		panic("Maximum value is less than minimum value")
	}

	var diff = high - low
	if diff <= 1 {
		return low
	}

	return low + m.Next()%diff
}

func (m *MersennelTwister) NextFloat() float32 {
	return float32(math.Abs(float64(m.Next()) / float64(0x7fffffff)))
}

func (m *MersennelTwister) Generate() {
	l := len(m.mt)
	for i := 0; i < l; i++ {
		var y = (m.mt[i] & 0x80000000) | (m.mt[(i+1)%624] & 0x7fffffff)
		m.mt[i] = m.mt[(i+397)%624] ^ (y >> 1)
		if y&1 == 1 {
			m.mt[i] = m.mt[i] ^ 2567483615
		}
	}
}
