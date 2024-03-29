package object

import (
	"project_b/common/base"
)

// 基礎物体静态信息
type ObjStaticInfo struct {
	id        int32             // 配置id
	typ       base.ObjectType   //  類型
	subType   base.ObjSubtype   // 子類型
	ownerType base.ObjOwnerType // 所有者类型
	camp      base.CampType     // 陣營
	x0, y0    int32             // 统一：矩形中心相对于位于局部坐标系的坐标
	w, l      int32             // 宽度w指x軸方向 長度l指y軸方向,
	rotation  int32             // 在本地(局部)坐標系中的旋轉朝向與"零角度朝向"的角度差值，度數[0, 360)，放到世界坐標中則要逆變換到零角度
	speed     int32             // 速度
	layer     int32             // 0-5
	collision bool              // 是否碰撞
}

// 创建物体静态信息
func NewObjStaticInfo(id int32, typ base.ObjectType, subType base.ObjSubtype, camp base.CampType, x0, y0, w, l int32, rotation int32, speed int32, layer int32, collision bool) *ObjStaticInfo {
	return &ObjStaticInfo{
		id: id, typ: typ, subType: subType, camp: camp, x0: x0, y0: y0, w: w, l: l, rotation: rotation, speed: speed, layer: layer, collision: collision,
	}
}

func (info ObjStaticInfo) Id() int32 {
	return info.id
}

func (info ObjStaticInfo) Pos0() base.Pos {
	return base.Pos{X: int32(info.x0), Y: int32(info.y0)}
}

func (info ObjStaticInfo) Width() int32 {
	return info.w
}

func (info ObjStaticInfo) Length() int32 {
	return info.l
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

// 移動物體靜態配置
type MovableObjStaticInfo struct {
	ObjStaticInfo
	MoveFunc func(IMovableObject, uint32) (int32, int32) // 移動函數
}

// 環繞物體靜態配置
type SurroundObjStaticInfo struct {
	MovableObjStaticInfo
	AroundRadius    int32 // 環繞半徑
	AngularVelocity int32 // 環繞角速度
	Clockwise       bool  // 是否順時針
}

// 炮彈靜態配置
type ShellStaticInfo struct {
	MovableObjStaticInfo
	Range                   int32 // 射程
	Damage                  int32 // 傷害
	BlastRadius             int32 // 爆炸半徑
	TrackTarget             bool  // 是否追蹤目標
	SearchTargetRadius      int32 // 搜索目標半徑
	SteeringAngularVelocity int32 // 轉向角速度(分(1/60度)/秒)
}

// 坦克炮彈配置
type TankShellConfig struct {
	ShellInfo         *ShellStaticInfo // 子彈配置ID
	AmountFireOneTime int8             // 一次發射炮彈量
	IntervalInFireMs  int32            // 一次發射的幾發炮彈之間的間隔時間(毫秒)
	CooldownMs        int32            // 每次發射冷卻時間(毫秒)
}

// 坦克靜態配置
type TankStaticInfo struct {
	MovableObjStaticInfo
	Level                   int32                // 等級
	SteeringAngularVelocity int32                // 轉向角速度(分(1/60度)/秒)
	ShellLaunchPos          base.Vec2            // 炮彈發射點 坦克局部坐標系中相對於坦克中心位置的坐標
	ShellConfig             TankShellConfig      // 炮彈配置
	ShieldConfig            TankShieldStaticInfo // 護盾靜態配置
	LaserId                 int32                // 激光id
}

// 坦克護盾配置
type TankShieldStaticInfo struct {
	Width, Length int32  // 長寬
	DurationMs    uint32 // 持續時間(毫秒)，0表示無限
}
