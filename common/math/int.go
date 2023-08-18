package math

import (
	"math"
	"project_b/log"
)

const (
	defaultPower10N = int(2)
)

var (
	customPower10N int = defaultPower10N
)

func SetN4Power10(n int) {
	if n < 1 || n > 6 {
		log.Error("math: n for power 10 is invalid, the range is [1, 6]")
		return
	}
	customPower10N = n
}

type Int32P int32

func NewInt32P(n int32) Int32P {
	pow := math.Pow10(customPower10N)
	v := (n * int32(pow) * (n * int32(pow)))
	if v > int32(math.MaxInt32) {
		panic("n is invalid")
	}
	// 检查n的范围
	return Int32P(n * int32(math.Pow10(customPower10N)))
}

func (i Int32P) Int32() int32 {
	return int32(i) / int32(math.Pow10(customPower10N))
}

func (i Int32P) Add(i2 Int32P) Int32P {
	return Int32P(int32(i) + int32(i2))
}

func (i Int32P) Sub(i2 Int32P) Int32P {
	return Int32P(int32(i) - int32(i2))
}

func (i Int32P) Mul(i2 Int32P) Int32P {
	return Int32P(int32(i) * int32(i2))
}

func (i Int32P) Div(i2 Int32P) Int32P {
	return Int32P(int32(i) / int32(i2))
}

type Int64P int64

func NewInt64P(n int64) Int64P {
	pow := math.Pow10(customPower10N)
	v := (n * int64(pow) * (n * int64(pow)))
	if v > int64(math.MaxInt64) {
		panic("n is invalid")
	}
	// 检查n的范围
	return Int64P(n * int64(math.Pow10(customPower10N)))
}
