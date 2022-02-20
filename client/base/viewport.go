package base

type Viewport struct {
	x, y int32 // 左上角坐标
	w, h int32 // 宽和高
}

func CreateViewport(x, y, w, h int32) *Viewport {
	return &Viewport{
		x: x, y: y, w: w, h: h,
	}
}
