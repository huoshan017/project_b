package client_base

type Viewport struct {
	x, y int32 // 左上角坐标
	w, h int32 // 宽和高
}

func (v Viewport) X() int32 {
	return v.x
}

func (v Viewport) Y() int32 {
	return v.y
}

func (v Viewport) W() int32 {
	return v.w
}

func (v Viewport) H() int32 {
	return v.h
}

func CreateViewport(x, y, w, h int32) *Viewport {
	return &Viewport{
		x: x, y: y, w: w, h: h,
	}
}
