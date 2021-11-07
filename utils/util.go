package utils

import (
	"math/rand"
	"project_b/common/object"
)

func RandomPosInRect(rect object.Rect) object.Pos {
	x := rand.Int31n(rect.RightBottom.X - rect.LeftTop.X)
	y := rand.Int31n(rect.RightBottom.Y - rect.LeftTop.Y)
	return object.Pos{X: x, Y: y}
}

func RandomPosInRectWithRand(r rand.Rand, rect object.Rect) object.Pos {
	x := r.Int31n(rect.RightBottom.X - rect.LeftTop.X)
	y := r.Int31n(rect.RightBottom.Y - rect.LeftTop.Y)
	return object.Pos{X: x, Y: y}
}