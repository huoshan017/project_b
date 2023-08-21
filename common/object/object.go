package object

import (
	"project_b/common/base"
)

/*******************************
				^ y
				|
				|
				|              x
------------------------------->
				|
				|
				|
				|
*******************************/

// 坐标位置
type Pos struct {
	X, Y int32 // 注意：x轴向右，y轴向上 为正方向
}

// 矩形
type Rect struct {
	LeftBottom Pos // 左上
	RightTop   Pos // 右下
}

// 物体结构
type object struct {
	instId            uint32         // 实例id
	ownerType         ObjOwnerType   // 所有制类型，可被动态临时改变，所以需要在对象中另外缓存
	currentCamp       CampType       // 當前陣營
	staticInfo        *ObjStaticInfo // 静态常量数据
	x, y              int32          // 指本地坐标系在父坐标系的坐标，如果父坐标系是世界坐标系，x、y就是世界坐标
	rotation          base.Angle     // 旋轉角度，以X軸正方向為0度，逆時針方向為旋轉正方向
	components        []IComponent   // 組件
	changedStaticInfo *ObjStaticInfo // 改变的静态常量数据
	toRecycle         bool           // 去回收
	colliderComp      *ColliderComp  // 碰撞器組件
}

// 回收
func (o *object) ToRecycle() {
	o.toRecycle = true
}

// 是否回收
func (o object) IsRecycle() bool {
	return o.toRecycle
}

// 初始化
func (o *object) Init(instId uint32, staticInfo *ObjStaticInfo) {
	o.instId = instId
	o.ownerType = staticInfo.ownerType
	o.currentCamp = staticInfo.camp
	o.staticInfo = staticInfo
	if staticInfo.collision {
		o.AddComp(&ColliderComp{obj: o})
	}
}

// 反初始化
func (o *object) Uninit() {
	o.instId = 0
	o.ownerType = OwnerNone
	o.currentCamp = CampTypeNone
	o.staticInfo = nil
	o.x, o.y = 0, 0
	o.rotation.Clear()
	o.components = o.components[:0]
	o.changedStaticInfo = nil
	o.toRecycle = false
}

// 靜態信息
func (o *object) ObjStaticInfo() *ObjStaticInfo {
	return o.staticInfo
}

// 设置静态信息
func (o *object) SetStaticInfo(staticInfo *ObjStaticInfo) {
	o.staticInfo = staticInfo
}

// 改变静态信息
func (o *object) ChangeStaticInfo(staticInfo *ObjStaticInfo) {
	o.changedStaticInfo = staticInfo
}

// 还原静态信息
func (o *object) RestoreStaticInfo() {
	if o.changedStaticInfo != nil {
		o.changedStaticInfo = nil
	}
}

// 实例id
func (o object) InstId() uint32 {
	return o.instId
}

// 配置id
func (o object) Id() int32 {
	if o.changedStaticInfo != nil {
		return o.changedStaticInfo.id
	}
	return o.staticInfo.id
}

// 原始配置id
func (o object) OriginId() int32 {
	return o.staticInfo.id
}

// 类型
func (o object) Type() ObjectType {
	return o.staticInfo.typ
}

// 子类型
func (o object) Subtype() ObjSubtype {
	return o.staticInfo.subType
}

// 所有者类型
func (o object) OwnerType() ObjOwnerType {
	return o.ownerType
}

// 陣營
func (o object) Camp() CampType {
	return o.currentCamp
}

// 靜態信息
func (o object) StaticInfo() *ObjStaticInfo {
	if o.changedStaticInfo != nil {
		return o.changedStaticInfo
	}
	return o.staticInfo
}

// 位置
func (o object) Pos() (int32, int32) {
	if o.changedStaticInfo != nil {
		return o.x + o.changedStaticInfo.x0, o.y + o.changedStaticInfo.y0
	}
	return o.x + o.staticInfo.x0, o.y + o.staticInfo.y0
}

// 坐标位置，相对于父坐标系
func (o *object) SetPos(x, y int32) {
	if o.changedStaticInfo != nil {
		o.x = x - o.changedStaticInfo.x0
		o.y = y - o.changedStaticInfo.y0
	} else {
		o.x = x - o.staticInfo.x0
		o.y = y - o.staticInfo.y0
	}
}

// 中心點坐標
func (o object) Center() (x, y int32) {
	if o.changedStaticInfo != nil {
		return o.changedStaticInfo.x0, o.changedStaticInfo.y0
	}
	return o.staticInfo.x0, o.staticInfo.y0
}

// 宽度
func (o object) Width() int32 {
	if o.changedStaticInfo != nil {
		return o.changedStaticInfo.w
	}
	return o.staticInfo.w
}

// 长度
func (o object) Length() int32 {
	if o.changedStaticInfo != nil {
		return o.changedStaticInfo.l
	}
	return o.staticInfo.l
}

func (o *object) x_y(x1, y1 int32) (int32, int32) {
	var (
		x0, y0 int32
	)
	if o.changedStaticInfo != nil {
		x0, y0 = o.x+o.changedStaticInfo.x0, o.y+o.changedStaticInfo.y0
	} else {
		x0, y0 = o.x+o.staticInfo.x0, o.y+o.staticInfo.y0
	}

	rotation := o.Rotation()
	// 公式
	// x = (x1-x0)*cos(a) - (y1-y0)*sin(a) + x0
	// y = (x1-x0)*sin(a) + (y1-y0)*cos(a) + y0
	x, y := base.Rotate(x1, y1, x0, y0, rotation)
	return x, y
}

// 左上坐標(相對於本地坐標系)
func (o object) LeftTop() (int32, int32) {
	var (
		x1, y1 int32
	)
	if o.changedStaticInfo != nil {
		x1, y1 = o.x+o.changedStaticInfo.x0+o.changedStaticInfo.l/2, o.y+o.changedStaticInfo.y0+o.changedStaticInfo.w/2
	} else {
		x1, y1 = o.x+o.staticInfo.x0+o.staticInfo.l/2, o.y+o.staticInfo.y0+o.staticInfo.w/2
	}
	return o.x_y(x1, y1)
}

// 左下坐标（相对于本地坐标系）
func (o object) LeftBottom() (int32, int32) {
	var (
		x1, y1 int32
	)
	if o.changedStaticInfo != nil {
		x1, y1 = o.x+o.changedStaticInfo.x0-o.changedStaticInfo.l/2, o.y+o.changedStaticInfo.y0+o.changedStaticInfo.w/2
	} else {
		x1, y1 = o.x+o.staticInfo.x0-o.staticInfo.l/2, o.y+o.staticInfo.y0+o.staticInfo.w/2
	}
	return o.x_y(x1, y1)
}

// 右下坐標 (相對於本地坐標系)
func (o object) RightBottom() (int32, int32) {
	var (
		x1, y1 int32
	)
	if o.changedStaticInfo != nil {
		x1, y1 = o.x+o.changedStaticInfo.x0-o.changedStaticInfo.l/2, o.y+o.changedStaticInfo.y0-o.changedStaticInfo.w/2
	} else {
		x1, y1 = o.x+o.staticInfo.x0-o.staticInfo.l/2, o.y+o.staticInfo.y0-o.staticInfo.w/2
	}
	return o.x_y(x1, y1)
}

// 右上坐标（相对于本地坐标系）
func (o object) RightTop() (int32, int32) {
	var (
		x1, y1 int32
	)
	if o.changedStaticInfo != nil {
		x1, y1 = o.x+o.changedStaticInfo.x0+o.changedStaticInfo.l/2, o.y+o.changedStaticInfo.y0-o.changedStaticInfo.w/2
	} else {
		x1, y1 = o.x+o.staticInfo.x0+o.staticInfo.l/2, o.y+o.staticInfo.y0-o.staticInfo.w/2
	}
	return o.x_y(x1, y1)
}

// 局部旋轉
func (o object) LocalRotation() base.Angle {
	return base.NewAngle(int16(o.staticInfo.rotation), 0)
}

// 世界旋轉
func (o object) WorldRotation() base.Angle {
	return o.rotation
}

// 旋轉
func (o object) Rotation() base.Angle {
	sr := o.staticInfo.rotation
	rotation := base.AngleAdd(o.rotation, base.NewAngle(int16(sr), 0))
	rotation.Normalize()
	return rotation
}

// 設置陣營
func (o *object) SetCamp(camp CampType) {
	o.currentCamp = camp
}

// 重置陣營
func (o *object) RestoreCamp() {
	o.currentCamp = o.staticInfo.camp
}

// 添加組件
func (o *object) AddComp(comp IComponent) {
	o.components = append(o.components, comp)
	if comp.Name() == "Collider" {
		o.colliderComp = comp.(*ColliderComp)
	}
}

// 去除組件
func (o *object) RemoveComp(name string) {
	for i, c := range o.components {
		if c.Name() == name {
			o.components = append(o.components[:i], o.components[i+1:]...)
			break
		}
	}
}

// 獲取組件
func (o *object) GetComp(name string) IComponent {
	for _, c := range o.components {
		if c.Name() == name {
			return c
		}
	}
	return nil
}

// 是否擁有組件
func (o object) HasComp(name string) bool {
	for _, c := range o.components {
		if c.Name() == name {
			return true
		}
	}
	return false
}

// 設置碰撞處理函數
func (o *object) SetCollisionHandle(handle func(IMovableObject, *CollisionInfo)) {
	if o.colliderComp != nil {
		o.colliderComp.SetCollisionHandle(handle)
	}
}

// 獲得碰撞器組件
func (o *object) GetColliderComp() *ColliderComp {
	return o.colliderComp
}

// 静态物体
type StaticObject struct {
	object
}

// 创建静态物体
func NewStaticObject() *StaticObject {
	obj := &StaticObject{}
	//obj.Init(instId, info)
	return obj
}

// 初始化
func (o *StaticObject) Init(instId uint32, info *ObjStaticInfo) {
	o.object.Init(instId, info)
}

// 反初始化
func (o *StaticObject) Uninit() {
	o.object.Uninit()
}

// 距離的平方
func SquareOfDistance(obj1, obj2 IObject) int64 {
	x1, y1 := obj1.Pos()
	x2, y2 := obj2.Pos()
	return (int64(x1 - x2)) ^ 2 + int64((y1 - y2)) ^ 2
}
