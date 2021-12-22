package object

import (
	"fmt"
	"project_b/common/base"
	"project_b/common/time"
)

// 坐标位置
type Pos struct {
	X, Y int32 // 注意：x轴向右，y轴向上 为正方向
}

// 矩形
type Rect struct {
	LeftTop     Pos // 左上
	RightBottom Pos // 右下
}

// 物体结构
type object struct {
	instId            uint64         // 实例id
	staticInfo        *ObjStaticInfo // 静态常量数据
	x, y              int32          // 指本地坐标系在父坐标系的坐标，如果父坐标系是世界坐标系，x、y就是世界坐标
	changedStaticInfo *ObjStaticInfo // 改变的静态常量数据
}

// 初始化
func (o *object) Init(instId uint64, staticInfo *ObjStaticInfo) {
	o.instId = instId
	o.staticInfo = staticInfo
}

// 反初始化
func (o *object) Uninit() {

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
func (o object) InstId() uint64 {
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
func (o object) Subtype() ObjSubType {
	return o.staticInfo.subType
}

// 位置
func (o object) Pos() (int32, int32) {
	return o.x, o.y
}

// 坐标位置，相对于父坐标系
func (o *object) SetPos(x, y int32) {
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
func (o object) Left() int32 {
	if o.changedStaticInfo != nil {
		return o.x + int32(o.changedStaticInfo.x0)
	}
	return o.x + int32(o.staticInfo.x0)
}

// 右侧坐标（相对于父坐标系）
func (o object) Right() int32 {
	if o.changedStaticInfo != nil {
		return o.Left() + int32(o.changedStaticInfo.w)
	}
	return o.Left() + int32(o.staticInfo.w)
}

// 顶部坐标（相对于父坐标系）
func (o object) Top() int32 {
	if o.changedStaticInfo != nil {
		return o.Bottom() + int32(o.changedStaticInfo.h)
	}
	return o.Bottom() + int32(o.staticInfo.h)
}

// 底部坐标（相对于父坐标系）
func (o object) Bottom() int32 {
	if o.changedStaticInfo != nil {
		return o.y + int32(o.changedStaticInfo.y0)
	}
	return o.y + int32(o.staticInfo.y0)
}

// 更新
func (o object) Update() {

}

type moveObjectState int32

const (
	stopped  = moveObjectState(0)
	toMove   = moveObjectState(1)
	isMoving = moveObjectState(2)
	toStop   = moveObjectState(3)
)

// 可移动的物体
type MovableObject struct {
	object
	dir   Direction       // 方向
	speed int32           // 当前移动速度（米/秒）
	state moveObjectState // 移动状态
	//moveDataList []*moveData     // 移动数据队列
	moveEvent   *base.Event // 移动事件
	stopEvent   *base.Event // 停止事件
	updateEvent *base.Event // 更新事件
}

// 创建可移动物体
func NewMovableObject(instId uint64, staticInfo *ObjStaticInfo) *MovableObject {
	o := &MovableObject{
		object:      object{instId: instId, staticInfo: staticInfo},
		dir:         staticInfo.dir,
		speed:       staticInfo.speed,
		moveEvent:   base.NewEvent(),
		stopEvent:   base.NewEvent(),
		updateEvent: base.NewEvent(),
	}
	return o
}

// 初始化
func (o *MovableObject) Init(instId uint64, staticInfo *ObjStaticInfo) {
	o.object.Init(instId, staticInfo)
	o.dir = staticInfo.dir
	o.speed = staticInfo.speed
	o.moveEvent = base.NewEvent()
	o.stopEvent = base.NewEvent()
	o.updateEvent = base.NewEvent()
}

// 反初始化
func (o *MovableObject) Uninit() {

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
func (o *MovableObject) SetCurrentSpeed(speed int32) {
	o.speed = speed
}

// 当前方向
func (o MovableObject) Dir() Direction {
	return o.dir
}

// 配置速度
func (o MovableObject) Speed() int32 {
	if o.changedStaticInfo != nil {
		return o.changedStaticInfo.speed
	}
	return o.staticInfo.speed
}

// 当前速度
func (o MovableObject) CurrentSpeed() int32 {
	return o.speed
}

// 移动
func (o *MovableObject) Move(dir Direction) {
	if dir < DirMin || dir > DirMax {
		str := fmt.Sprintf("invalid object direction %v", dir)
		panic(str)
	}
	if o.state == stopped {
		o.state = toMove
		o.dir = dir
	}
}

// 停止
func (o *MovableObject) Stop() {
	// 准备移动则直接停止
	if o.state == toMove {
		o.state = stopped
		return
	}
	// 正在移动则准备停止
	if o.state == isMoving {
		o.state = toStop
		return
	}
}

// 更新
func (o *MovableObject) Update(tick time.Duration) {
	if o.state == stopped {
		return
	}

	if o.state == toMove {
		o.state = isMoving
		// args[0]: object.Pos
		// args[1]: object.Direction
		// args[2]: int32
		o.moveEvent.Call(Pos{X: o.x, Y: o.y}, o.dir, o.CurrentSpeed())
		return
	}

	distance := (float64(o.CurrentSpeed()) * float64(tick) / float64(time.Second))
	switch o.dir {
	case DirLeft:
		x := float64(o.x)
		x -= distance
		o.x = int32(x)
	case DirRight:
		x := float64(o.x)
		x += distance
		o.x = int32(x)
	case DirUp:
		y := float64(o.y)
		y -= distance
		o.y = int32(y)
	case DirDown:
		y := float64(o.y)
		y += distance
		o.y = int32(y)
	default:
		panic("invalid direction")
	}

	if o.state == isMoving {
		// args[0]: object.Pos
		// args[1]: object.Direction
		// args[2]: int32
		o.updateEvent.Call(Pos{X: o.x, Y: o.y}, o.dir, o.CurrentSpeed())
	} else if o.state == toStop {
		o.state = stopped
		// args[0]: object.Pos
		// args[1]: object.Direction
		// args[2]: int32
		o.stopEvent.Call(Pos{X: o.x, Y: o.y}, o.dir, o.CurrentSpeed())
	}
}

// 注册移动事件
func (o *MovableObject) RegisterMoveEventHandle(handle func(args ...interface{})) {
	o.moveEvent.Register(handle)
}

// 注销移动事件
func (o *MovableObject) UnregisterMoveEventHandle(handle func(args ...interface{})) {
	o.moveEvent.Unregister(handle)
}

// 注册停止移动事件
func (o *MovableObject) RegisterStopMoveEventHandle(handle func(args ...interface{})) {
	o.stopEvent.Register(handle)
}

// 注销停止移动事件
func (o *MovableObject) UnregisterStopMoveEventHandle(handle func(args ...interface{})) {
	o.stopEvent.Unregister(handle)
}

// 注册更新事件
func (o *MovableObject) RegisterUpdateEventHandle(handle func(args ...interface{})) {
	o.updateEvent.Register(handle)
}

// 注销更新事件
func (o *MovableObject) UnregisterUpdateEventHandle(handle func(args ...interface{})) {
	o.updateEvent.Unregister(handle)
}

// 车辆
type Vehicle struct {
	MovableObject
}

// 创建车辆
func NewVehicle(instId uint64, staticInfo *ObjStaticInfo) *Vehicle {
	o := &Vehicle{
		MovableObject: *NewMovableObject(instId, staticInfo),
	}
	return o
}

// 初始化
func (v *Vehicle) Init(instId uint64, staticInfo *ObjStaticInfo) {
	v.MovableObject.Init(instId, staticInfo)
}

// 反初始化
func (v *Vehicle) Uninit() {

}

// 坦克
type Tank struct {
	Vehicle
	level       int32
	changeEvent *base.Event
}

// 创建坦克
func NewTank(instId uint64, staticInfo *ObjStaticInfo) *Tank {
	return &Tank{
		Vehicle:     *NewVehicle(instId, staticInfo),
		level:       1,
		changeEvent: base.NewEvent(),
	}
}

// 初始化
func (t *Tank) Init(instId uint64, staticInfo *ObjStaticInfo) {
	t.Vehicle.Init(instId, staticInfo)
	t.level = 1
	t.changeEvent = base.NewEvent()
}

// 反初始化
func (t *Tank) Uninit() {

}

// 等级
func (t Tank) Level() int32 {
	return t.level
}

// 设置等级
func (t *Tank) SetLevel(level int32) {
	t.level = level
}

// 变化
func (t *Tank) Change(info *ObjStaticInfo) {
	t.ChangeStaticInfo(info)
	t.SetCurrentSpeed(info.Speed())
	t.changeEvent.Call(info, t.level)
}

// 还原
func (t *Tank) Restore() {
	t.RestoreStaticInfo()
	t.changeEvent.Call(t.staticInfo, t.level)
}

// 注册变化事件
func (t *Tank) RegisterChangeEventHandle(handle func(args ...interface{})) {
	t.changeEvent.Register(handle)
}

// 注销变化事件
func (t *Tank) UnregisterChangeEventHandle(handle func(args ...interface{})) {
	t.changeEvent.Unregister(handle)
}
