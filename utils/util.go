package utils

import (
	"math"
	"math/rand"
	"project_b/common/base"
	"time"
)

var (
	defaultRand = rand.New(rand.NewSource(time.Now().Unix()))
)

func RandomPosInRect(rect base.Rect) base.Pos {
	x := defaultRand.Int31n(int32(rect.RightTop.X - rect.LeftBottom.X))
	y := defaultRand.Int31n(int32(rect.RightTop.Y - rect.LeftBottom.Y))
	return base.Pos{X: x, Y: y}
}

func RandomPosInRectWithRand(r rand.Rand, rect base.Rect) base.Pos {
	x := r.Int31n(int32(rect.RightTop.X - rect.LeftBottom.X))
	y := r.Int31n(int32(rect.RightTop.Y - rect.LeftBottom.Y))
	return base.Pos{X: x, Y: y}
}

func Float32IsEqual(a, b float32) bool {
	return math.Abs(float64(a-b)) < 0.0001
}

func Float64IsEqual(a, b float64) bool {
	return math.Abs(a-b) < 0.000001
}
