package math

import (
	"math"
	"project_b/common/utils"
)

var (
	DistZero     = NewDist(0)
	DistMaxValue = NewDist(math.MaxInt32)
)

/**
 * 1d world distance - 1024 units = 1 cell
 */
type Dist struct {
	length int32
}

func NewDist(r int32) Dist {
	return Dist{length: r}
}

func NewDistFromCells(cells int32) Dist {
	return NewDist(cells * oneValueUnits)
}

func (d Dist) LengthSquared() int64 {
	return int64(d.length * d.length)
}

func (d *Dist) Add(dist Dist) {
	d.length += dist.length
}

func (d *Dist) Sub(dist Dist) {
	d.length -= dist.length
}

func (d *Dist) Negative() {
	d.length = -d.length
}

func (d *Dist) Div(divisor int32) {
	d.length /= divisor
}

func (d *Dist) Mul(a int32) {
	d.length *= a
}

func (d Dist) LessThan(dist Dist) bool {
	return d.length < dist.length
}

func (d Dist) GreaterThan(dist Dist) bool {
	return d.length > dist.length
}

func (d Dist) EqualTo(dist Dist) bool {
	return d.length == dist.length
}

func (d Dist) LessAndEqual(dist Dist) bool {
	return d.length <= dist.length
}

func (d Dist) GreaterAndEqual(dist Dist) bool {
	return d.length >= dist.length
}

func (d Dist) NotEqualTo(dist Dist) bool {
	return d.length != dist.length
}

func DistLessThan(a, b Dist) bool {
	return a.LessThan(b)
}

func DistGreaterThan(a, b Dist) bool {
	return a.GreaterThan(b)
}

func DistLessAndEqual(a, b Dist) bool {
	return a.LessAndEqual(b)
}

func DistGreaterAndEqual(a, b Dist) bool {
	return a.GreaterAndEqual(b)
}

func DistEqual(a, b Dist) bool {
	return a.EqualTo(b)
}

func DistNotEqual(a, b Dist) bool {
	return a.NotEqualTo(b)
}

/**
 * Sampled a N-sample probability density function in the range [-1024 1024]
 * 1 sample produces a rectangular probability
 * 2 samples produces a triangular probability
 * ...
 * N samples approximates a true Gaussian
 */
func DistFromPDF(m *MersennelTwister, samples int32) Dist {
	arr := utils.MakeArray(samples, func(int32) interface{} {
		return m.Next2(-oneValueUnits, oneValueUnits)
	})
	var sum int32
	for _, a := range arr {
		sum += a.(int32)
	}
	return NewDist(sum / samples)
}
