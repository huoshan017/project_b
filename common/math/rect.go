package math

// 矩形区域，(x, y)表示左下角，w宽度，h高度
type Rect struct {
	x, y, w, h int32
}

func NewRect(x, y, w, h int32) *Rect {
	return &Rect{
		x: x, y: y, w: w, h: h,
	}
}

func (r Rect) X() int32 {
	return r.x
}

func (r Rect) Y() int32 {
	return r.y
}

func (r Rect) W() int32 {
	return r.w
}

func (r Rect) H() int32 {
	return r.h
}
