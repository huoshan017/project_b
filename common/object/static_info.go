package object

// 物体静态信息
type ObjStaticInfo struct {
	id      int32 // 配置id
	typ     ObjectType
	subType ObjSubType
	x0, y0  int32 // 统一：矩形左下角相对于位于局部坐标系的坐标
	w, h    int32 // 宽度高度
	dir     Direction
	speed   float32
}

// 创建物体静态信息
func NewObjStaticInfo(id int32, typ ObjectType, subType ObjSubType, x0, y0, w, h int32, speed float32, dir Direction) *ObjStaticInfo {
	return &ObjStaticInfo{
		id: id, typ: typ, subType: subType, x0: x0, y0: y0, w: w, h: h, dir: dir, speed: speed,
	}
}

func (info ObjStaticInfo) Id() int32 {
	return info.id
}

func (info ObjStaticInfo) Pos0() Pos {
	return Pos{float64(info.x0), float64(info.y0)}
}

func (info ObjStaticInfo) Width() int32 {
	return info.w
}

func (info ObjStaticInfo) Height() int32 {
	return info.h
}

func (info ObjStaticInfo) Dir() Direction {
	return info.dir
}

func (info ObjStaticInfo) Speed() float32 {
	return info.speed
}
