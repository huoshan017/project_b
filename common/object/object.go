package object

import (
	"time"
)

// 坐标位置
type Pos struct {
	X, Y float64 // 注意：x轴向右，y轴向上 为正方向
}

// 矩形
type Rect struct {
	LeftTop     Pos // 左上
	RightBottom Pos // 右下
}

// 物体结构
type object struct {
	staticInfo        *ObjStaticInfo // 静态常量数据
	x, y              float64        // 指本地坐标系在父坐标系的坐标，如果父坐标系是世界坐标系，x、y就是世界坐标
	changedStaticInfo *ObjStaticInfo // 改变的静态常量数据
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
func (o object) Subtype() ObjSubType {
	return o.staticInfo.subType
}

// 位置
func (o object) Pos() (float64, float64) {
	return o.x, o.y
}

// 坐标位置，相对于父坐标系
func (o *object) SetPos(x, y float64) {
	o.x = x
	o.y = y
}

// 宽度
func (o object) Width() int32 {
	if o.changedStaticInfo != nil {
		return o.changedStaticInfo.w
	}
	return o.staticInfo.w
}

// 长度
func (o object) Height() int32 {
	if o.changedStaticInfo != nil {
		return o.changedStaticInfo.h
	}
	return o.staticInfo.h
}

// 左侧坐标（相对于父坐标系）
func (o object) Left() float64 {
	if o.changedStaticInfo != nil {
		return o.x + float64(o.changedStaticInfo.x0)
	}
	return o.x + float64(o.staticInfo.x0)
}

// 右侧坐标（相对于父坐标系）
func (o object) Right() float64 {
	if o.changedStaticInfo != nil {
		return o.Left() + float64(o.changedStaticInfo.w)
	}
	return o.Left() + float64(o.staticInfo.w)
}

// 顶部坐标（相对于父坐标系）
func (o object) Top() float64 {
	if o.changedStaticInfo != nil {
		return o.Bottom() + float64(o.changedStaticInfo.h)
	}
	return o.Bottom() + float64(o.staticInfo.h)
}

// 底部坐标（相对于父坐标系）
func (o object) Bottom() float64 {
	if o.changedStaticInfo != nil {
		return o.y + float64(o.changedStaticInfo.y0)
	}
	return o.y + float64(o.staticInfo.y0)
}

// 更新
func (o object) Update() {

}

// 可移动的物体
type MovableObject struct {
	object
	dir             Direction
	speed           float32 // 当前移动速度（米/秒）
	state           int32   // 状态    0. 停止  1. 移动
	lastUpdateTick  time.Duration
	minMoveDistance float32
}

// 创建可移动物体
func NewMovableObject(staticInfo *ObjStaticInfo) *MovableObject {
	o := &MovableObject{
		object: object{staticInfo: staticInfo},
		dir:    staticInfo.dir,
		speed:  staticInfo.speed,
	}
	return o
}

// 改变静态信息
func (o *MovableObject) ChangeStaticInfo(staticInfo *ObjStaticInfo) {
	o.object.ChangeStaticInfo(staticInfo)
	o.SetCurrentSpeed(staticInfo.Speed())
}

// 恢复静态信息
func (o *MovableObject) RestoreStaticInfo() {
	o.object.RestoreStaticInfo()
	o.SetCurrentSpeed(o.staticInfo.Speed())
}

// 设置当前方向
func (o *MovableObject) SetDir(dir Direction) {
	o.dir = dir
}

// 设置当前速度
func (o *MovableObject) SetCurrentSpeed(speed float32) {
	o.speed = speed
}

// 设置最小移动距离
func (o *MovableObject) SetMinMoveDistance(d float32) {
	o.minMoveDistance = d
}

// 当前方向
func (o MovableObject) Dir() Direction {
	return o.dir
}

// 配置速度
func (o MovableObject) Speed() float32 {
	if o.changedStaticInfo != nil {
		return o.changedStaticInfo.speed
	}
	return o.staticInfo.speed
}

// 当前速度
func (o MovableObject) CurrentSpeed() float32 {
	return o.speed
}

// 移动，指定方向和下一个点的位置，提供给使用者判断是否会碰撞
func (o *MovableObject) Move(dir Direction) {
	o.dir = dir
	o.state = 1
}

// 是否在移动
func (o *MovableObject) IsMove() bool {
	return o.state == 1
}

// 停止
func (o *MovableObject) Stop() {
	o.state = 0
}

// 更新
func (o *MovableObject) Update(tick time.Duration) {
	if o.state == 0 {
		return
	}

	if o.minMoveDistance < Delta {
		o.minMoveDistance = DefaultMinMoveDistance
	}

	o.lastUpdateTick += tick
	diffMs := float64(o.lastUpdateTick / time.Millisecond)
	distance := float64(o.speed/1000) * diffMs

	// 不够最小移动距离
	if distance < float64(o.minMoveDistance) {
		return
	}

	switch o.dir {
	case DirLeft:
		o.x -= distance
	case DirRight:
		o.x += distance
	case DirUp:
		o.y -= distance
	case DirDown:
		o.y += distance
	default:
		return
	}

	o.lastUpdateTick -= (time.Duration(diffMs) * time.Millisecond)
}

// 车辆
type Vehicle struct {
	MovableObject
}

// 创建车辆
func NewVehicle(staticInfo *ObjStaticInfo) *Vehicle {
	o := &Vehicle{
		MovableObject: *NewMovableObject(staticInfo),
	}
	return o
}

// 坦克
type Tank struct {
	Vehicle
	level int32
}

// 创建坦克
func NewTank(staticInfo *ObjStaticInfo) *Tank {
	return &Tank{
		Vehicle: *NewVehicle(staticInfo),
		level:   1,
	}
}

// 等级
func (t Tank) Level() int32 {
	return t.level
}

// 设置等级
func (t *Tank) SetLevel(level int32) {
	t.level = level
}
