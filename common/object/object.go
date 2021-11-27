package object

import (
	"project_b/common/base"
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
	instId            uint64         // 实例id
	staticInfo        *ObjStaticInfo // 静态常量数据
	x, y              float64        // 指本地坐标系在父坐标系的坐标，如果父坐标系是世界坐标系，x、y就是世界坐标
	changedStaticInfo *ObjStaticInfo // 改变的静态常量数据
}

// 初始化
func (o *object) Init() {

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
	dir           Direction   // 方向
	speed         float32     // 当前移动速度（米/秒）
	isMoving      bool        // 是否已更新
	startMoveTime time.Time   // 移动开始时间
	stopTime      time.Time   // 停止移动时间
	moveDataList  []*moveData // 移动数据队列
	moveEvent     *base.Event // 移动事件
	stopEvent     *base.Event // 停止事件
	updateEvent   *base.Event // 更新事件
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

// 移动
func (o *MovableObject) Move(dir Direction) {
	if o.startMoveTime.IsZero() || !o.isMoving {
		o.dir = dir
		o.startMoveTime = time.Now()
		o.isMoving = true
		o.moveEvent.Call(o.startMoveTime, dir, float64(o.CurrentSpeed()))
	}
}

// 停止
func (o *MovableObject) Stop() {
	if !o.isMoving {
		return
	}
	o.stopTime = time.Now()
	o.isMoving = false
	d := mdFreeList.get()
	d.dir = o.dir
	d.speed = float64(o.CurrentSpeed())
	d.duration = o.stopTime.Sub(o.startMoveTime)
	o.moveDataList = append(o.moveDataList, d)
	o.stopEvent.Call(o.stopTime)
}

// 更新
func (o *MovableObject) Update(tick time.Duration) {
	updated := false

	if o.moveDataList != nil && len(o.moveDataList) > 0 {
		for _, md := range o.moveDataList {
			o.update(md)
			mdFreeList.put(md)
		}
		o.moveDataList = o.moveDataList[:0]
		updated = true
	}

	// 上一次移动还没停止
	if o.startMoveTime.After(o.stopTime) {
		now := time.Now()
		var md = moveData{duration: now.Sub(o.startMoveTime), speed: float64(o.CurrentSpeed()), dir: o.dir}
		o.update(&md)
		o.startMoveTime = now
		updated = true
	}

	if updated {
		o.updateEvent.Call(o.x, o.y)
	}
}

// 内部更新函数
func (o *MovableObject) update(md *moveData) {
	distance := float64(md.speed) * float64(md.duration) / float64(time.Second)
	switch md.dir {
	case DirLeft:
		o.x -= distance
	case DirRight:
		o.x += distance
	case DirUp:
		o.y -= distance
	case DirDown:
		o.y += distance
	default:
		panic("invalid direction")
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
