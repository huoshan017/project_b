package base

import "math"

// 坐标位置
type Pos struct {
	X, Y int32 // 注意：x轴向右，y轴向上 为正方向
}

func NewPos(x, y int32) Pos {
	return Pos{X: x, Y: y}
}

func (p *Pos) Clear() {
	p.X = math.MinInt32
	p.Y = math.MinInt32
}

// 矩形
type Rect struct {
	LeftBottom Pos // 左上
	RightTop   Pos // 右下
}

// 距離
func Distance(p1, p2 *Pos) uint32 {
	return Sqrt(uint32((p1.X-p2.X)*(p1.X-p2.X) + (p1.Y-p2.Y)*(p1.Y-p2.Y)))
}

// 方向位置
func DirPos(x, y int32, distance int32, dir Angle) (int32, int32) {
	sn, sd := Sine(dir)
	cn, cd := Cosine(dir)
	dx, dy := distance*cn/cd, distance*sn/sd
	return x + dx, y + dy
}

// 計算移動位置
func MovePos(x, y int32, moveDir Angle, speed int32, tickMs uint32) (int32, int32) {
	distance := int32(int64(speed) * int64(tickMs) / 1000)
	return DirPos(x, y, distance, moveDir)
}

// 計算點乘
func Dot(pos1, pos2 *Pos) int32 {
	return pos1.X*pos2.X + pos1.Y*pos2.Y
}

// 計算叉乘
func Cross(pos1, pos2 *Pos) int32 {
	return pos1.X*pos2.Y - pos1.Y*pos2.X
}
