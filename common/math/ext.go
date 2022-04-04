package math

type SqrtRoundMode int

const (
	RoundFloor   SqrtRoundMode = iota
	RoundNearest SqrtRoundMode = 1
	RoundCeiling SqrtRoundMode = 2
)

func Sqrt32(number uint32, rmode SqrtRoundMode) uint32 {
	divisor := uint32(1 << 30)
	root := uint32(0)
	remainder := number

	for divisor > number {
		divisor >>= 2
	}

	for divisor != 0 {
		if root+divisor <= remainder {
			remainder -= root + divisor
			root += 2 * divisor
		}
		root >>= 1
		divisor >>= 2
	}

	if rmode == RoundNearest && remainder > root {
		root += 1
	} else if rmode == RoundCeiling && root*root < number {
		root += 1
	}

	return root
}

func Sqrt32Floor(number uint32) uint32 {
	return Sqrt32(number, RoundFloor)
}

func Sqrt32Nearest(number uint32) uint32 {
	return Sqrt32(number, RoundNearest)
}

func Sqrt32Ceiling(number uint32) uint32 {
	return Sqrt32(number, RoundCeiling)
}

func Sqrt64(number uint64, rmode SqrtRoundMode) uint64 {
	divisor := uint64(1 << 62)
	root := uint64(0)
	remainder := number

	for divisor > number {
		divisor >>= 2
	}

	for divisor != 0 {
		if root+divisor <= remainder {
			remainder -= root + divisor
			root += 2 * divisor
		}
		root >>= 1
		divisor >>= 2
	}

	if rmode == RoundNearest && remainder > root {
		root += 1
	} else if rmode == RoundCeiling && root*root < number {
		root += 1
	}

	return root
}

func Sqrt64Floor(number uint64) uint64 {
	return Sqrt64(number, RoundFloor)
}

func Sqrt64Nearest(number uint64) uint64 {
	return Sqrt64(number, RoundNearest)
}

func Sqrt64Ceiling(number uint64) uint64 {
	return Sqrt64(number, RoundCeiling)
}
