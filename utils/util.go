package utils

import (
	"math"
	"math/rand"
	"project_b/common/object"
)

func RandomPosInRect(rect object.Rect) object.Pos {
	x := rand.Int31n(int32(rect.RightBottom.X - rect.LeftTop.X))
	y := rand.Int31n(int32(rect.RightBottom.Y - rect.LeftTop.Y))
	return object.Pos{X: float64(x), Y: float64(y)}
}

func RandomPosInRectWithRand(r rand.Rand, rect object.Rect) object.Pos {
	x := r.Int31n(int32(rect.RightBottom.X - rect.LeftTop.X))
	y := r.Int31n(int32(rect.RightBottom.Y - rect.LeftTop.Y))
	return object.Pos{X: float64(x), Y: float64(y)}
}

func Float32IsEqual(a, b float32) bool {
	return math.Abs(float64(a - b)) < 0.0001
}

func Float64IsEqual(a, b float64) bool {
	 return math.Abs(a - b) < 0.000001
}