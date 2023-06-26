package object

// 物体静态信息
type ObjStaticInfo struct {
	id        int32 // 配置id
	typ       ObjectType
	subType   ObjSubType
	ownerType ObjOwnerType // 所有者类型
	x0, y0    int32        // 统一：矩形左下角相对于位于局部坐标系的坐标
	w, h      int32        // 宽度高度
	dir       Direction
	speed     int32
	layer     int32 // 0-5
	collision bool  // 是否碰撞
}

// 创建物体静态信息
func NewObjStaticInfo(id int32, typ ObjectType, subType ObjSubType, x0, y0, w, h int32, speed int32, dir Direction, layer int32, collision bool) *ObjStaticInfo {
	return &ObjStaticInfo{
		id: id, typ: typ, subType: subType, x0: x0, y0: y0, w: w, h: h, dir: dir, speed: speed, layer: layer, collision: collision,
	}
}

func (info ObjStaticInfo) Id() int32 {
	return info.id
}

func (info ObjStaticInfo) Pos0() Pos {
	return Pos{int32(info.x0), int32(info.y0)}
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

func (info ObjStaticInfo) Speed() int32 {
	return info.speed
}

func (info ObjStaticInfo) Layer() int32 {
	return info.layer
}

func (info ObjStaticInfo) Collision() bool {
	return info.collision
}

type BulletStaticInfo struct {
	ObjStaticInfo
	Range       int32 // 射程
	Damage      int32 // 傷害
	BlastRadius int32 // 爆炸半徑
}

type TankBulletConfig struct {
	BulletId          int32 // 子彈配置ID
	AmountFireOneTime int8  // 一次發射炮彈量
	IntervalInFire    int32 // 一次發射的幾發炮彈之間的間隔時間(毫秒)
	Cooldown          int32 // 每次發射冷卻時間(毫秒)
}

type TankStaticInfo struct {
	ObjStaticInfo
	TankBulletConfig
}
