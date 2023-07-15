package object

import "project_b/common/time"

// 物体静态信息
type ObjStaticInfo struct {
	id        int32 // 配置id
	typ       ObjectType
	subType   ObjSubType
	ownerType ObjOwnerType // 所有者类型
	camp      CampType     // 陣營
	x0, y0    int32        // 统一：矩形中心相对于位于局部坐标系的坐标
	w, l      int32        // 宽度長度
	//dir       Direction
	speed     int32
	layer     int32 // 0-5
	collision bool  // 是否碰撞
}

// 创建物体静态信息
func NewObjStaticInfo(id int32, typ ObjectType, subType ObjSubType, camp CampType, x0, y0, w, l int32, speed int32 /*dir Direction, */, layer int32, collision bool) *ObjStaticInfo {
	return &ObjStaticInfo{
		id: id, typ: typ, subType: subType, camp: camp, x0: x0, y0: y0, w: w, l: l /*dir: dir,*/, speed: speed, layer: layer, collision: collision,
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

func (info ObjStaticInfo) Length() int32 {
	return info.l
}

/*func (info ObjStaticInfo) Dir() Direction {
	return info.dir
}*/

func (info ObjStaticInfo) Speed() int32 {
	return info.speed
}

func (info ObjStaticInfo) Layer() int32 {
	return info.layer
}

func (info ObjStaticInfo) Collision() bool {
	return info.collision
}

// 移動物體靜態配置
type MovableObjStaticInfo struct {
	ObjStaticInfo
	MoveFunc func(IMovableObject, time.Duration) (int32, int32) // 移動函數
}

// 環繞物體靜態配置
type SurroundObjStaticInfo struct {
	MovableObjStaticInfo
	AroundRadius    int32 // 環繞半徑
	AngularVelocity int32 // 環繞角速度
	Clockwise       bool  // 是否順時針
}

// 子彈靜態配置
type BulletStaticInfo struct {
	MovableObjStaticInfo
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

// 坦克靜態配置
type TankStaticInfo struct {
	MovableObjStaticInfo
	Orientation  int32 // 朝向
	Level        int32 // 等級
	BulletConfig TankBulletConfig
}

// 效果靜態信息
type EffectStaticInfo struct {
	Id            int32      // 配置id
	Et            EffectType // 效果類型
	Param         int32      // 參數
	Width, Height int32      // 寬高
}
