package object

import (
	"fmt"
	"math"
	"project_b/common/base"
	"project_b/common/log"
	"project_b/common/time"
	"unsafe"
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
	orientation       int32          // 旋轉角度，以X軸正方向為0度，逆時針方向為旋轉正方向
	components        []IComponent   // 組件
	changedStaticInfo *ObjStaticInfo // 改变的静态常量数据
	toRecycle         bool           // 去回收
	super             IObject        // 派生類對象
	destroyedEvent    base.Event     // 銷毀事件
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
		o.AddComp(&ColliderComp{})
	}
}

// 反初始化
func (o *object) Uninit() {
	o.destroyedEvent.Call(o.instId)
	o.destroyedEvent.Clear()
	o.instId = 0
	o.ownerType = OwnerNone
	o.currentCamp = CampTypeNone
	o.staticInfo = nil
	o.x, o.y = 0, 0
	o.orientation = 0
	o.components = o.components[:0]
	o.changedStaticInfo = nil
	o.super = nil
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
func (o object) Subtype() ObjSubType {
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
		x0, y0 = o.changedStaticInfo.x0, o.changedStaticInfo.y0
	} else {
		x0, y0 = o.staticInfo.x0, o.staticInfo.y0
	}
	cos, sin := math.Sincos(float64(o.orientation) * math.Pi / 180)
	// 公式
	// x = (x1-x0)*cos(a) - (y1-y0)*sin(a) + x0
	// y = (x1-x0)*sin(a) + (y1-y0)*cos(a) + y0
	x := o.x + int32(float64(x1-x0)*cos-float64(y1-y0)*sin) + x0
	y := o.y + int32(float64(x1-x0)*sin+float64(y1-y0)*cos) + y0
	return x, y
}

// 左上坐標(相對於本地坐標系)
func (o object) LeftTop() (int32, int32) {
	var (
		x1, y1 int32
	)
	if o.changedStaticInfo != nil {
		x1, y1 = o.changedStaticInfo.l/2+o.changedStaticInfo.x0, o.changedStaticInfo.w/2+o.changedStaticInfo.y0
	} else {
		x1, y1 = o.staticInfo.l/2+o.staticInfo.x0, o.staticInfo.w/2+o.staticInfo.y0
	}
	return o.x_y(x1, y1)
}

// 左下坐标（相对于本地坐标系）
func (o object) LeftBottom() (int32, int32) {
	var (
		x1, y1 int32
	)
	if o.changedStaticInfo != nil {
		x1, y1 = o.changedStaticInfo.x0-o.changedStaticInfo.l/2, o.changedStaticInfo.y0+o.changedStaticInfo.w/2
	} else {
		x1, y1 = o.staticInfo.x0-o.staticInfo.l/2, o.staticInfo.y0+o.staticInfo.w/2
	}
	return o.x_y(x1, y1)
}

// 右下坐標 (相對於本地坐標系)
func (o object) RightBottom() (int32, int32) {
	var (
		x1, y1 int32
	)
	if o.changedStaticInfo != nil {
		x1, y1 = o.changedStaticInfo.x0-o.changedStaticInfo.l/2, o.changedStaticInfo.y0-o.changedStaticInfo.w/2
	} else {
		x1, y1 = o.staticInfo.x0-o.staticInfo.l/2, o.staticInfo.y0-o.staticInfo.w/2
	}
	return o.x_y(x1, y1)
}

// 右上坐标（相对于本地坐标系）
func (o object) RightTop() (int32, int32) {
	var (
		x1, y1 int32
	)
	if o.changedStaticInfo != nil {
		x1, y1 = o.changedStaticInfo.x0+o.changedStaticInfo.l/2, o.changedStaticInfo.y0-o.changedStaticInfo.w/2
	} else {
		x1, y1 = o.staticInfo.x0+o.staticInfo.l/2, o.staticInfo.y0-o.staticInfo.w/2
	}
	return o.x_y(x1, y1)
}

// 朝向角度
func (o object) Orientation() int32 {
	return o.orientation
}

// 原始左坐標
func (o object) OriginalLeft() int32 {
	if o.changedStaticInfo != nil {
		return o.x + o.changedStaticInfo.x0 - o.changedStaticInfo.w/2
	}
	return o.x + o.staticInfo.x0 - o.staticInfo.w/2
}

// 原始右坐標
func (o object) OriginalRight() int32 {
	if o.changedStaticInfo != nil {
		return o.x + o.changedStaticInfo.x0 + o.changedStaticInfo.w/2
	}
	return o.x + o.staticInfo.x0 + o.staticInfo.w/2
}

// 原始上坐標
func (o object) OriginalTop() int32 {
	if o.changedStaticInfo != nil {
		return o.y + o.changedStaticInfo.y0 + o.changedStaticInfo.l/2
	}
	return o.y + o.staticInfo.y0 + o.staticInfo.l/2
}

// 原始下坐標
func (o object) OriginalBottom() int32 {
	if o.changedStaticInfo != nil {
		return o.y + o.changedStaticInfo.y0 - o.changedStaticInfo.l/2
	}
	return o.y + o.staticInfo.y0 - o.staticInfo.l/2
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

// 設置派生類
func (o *object) setSuper(super IObject) {
	o.super = super
}

// 注冊銷毀事件處理函數
func (o *object) RegisterDestroyedEventHandle(handle func(...any)) {
	o.destroyedEvent.Register(handle)
}

// 注銷銷毀事件處理函數
func (o *object) UnregisterDestroyedEventHandle(handle func(...any)) {
	o.destroyedEvent.Unregister(handle)
}

// 静态物体
type StaticObject struct {
	object
}

// 创建静态物体
func NewStaticObject(instId uint32, info *ObjStaticInfo) *StaticObject {
	obj := &StaticObject{}
	obj.Init(instId, info)
	return obj
}

// 初始化
func (o *StaticObject) Init(instId uint32, info *ObjStaticInfo) {
	o.object.Init(instId, info)
	o.object.setSuper(o)
}

// 反初始化
func (o *StaticObject) Uninit() {
	o.object.Uninit()
}

// 更新
func (o *StaticObject) Update(tick time.Duration) {

}

// 可移动物体状态
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
	dir            Direction       // 方向
	speed          int32           // 当前移动速度（米/秒）
	lastX, lastY   int32           // 上次更新的位置
	state          moveObjectState // 移动状态
	mySuper        IMovableObject  // 父類
	checkMoveEvent base.Event      // 檢查坐標事件
	moveEvent      base.Event      // 移动事件
	stopEvent      base.Event      // 停止事件
	updateEvent    base.Event      // 更新事件
}

// 创建可移动物体
func NewMovableObject(instId uint32, staticInfo *ObjStaticInfo) *MovableObject {
	o := &MovableObject{}
	o.Init(instId, staticInfo)
	return o
}

// 初始化
func (o *MovableObject) Init(instId uint32, staticInfo *ObjStaticInfo) {
	o.object.Init(instId, staticInfo)
	o.dir = staticInfo.dir
	o.speed = staticInfo.speed
	o.setSuper(o)
}

// 反初始化
func (o *MovableObject) Uninit() {
	o.dir = DirNone
	o.speed = 0
	o.lastX = 0
	o.lastY = 0
	o.state = stopped
	o.checkMoveEvent.Clear()
	o.moveEvent.Clear()
	o.stopEvent.Clear()
	o.updateEvent.Clear()
	o.object.Uninit()
}

func (o *MovableObject) MovableObjStaticInfo() *MovableObjStaticInfo {
	return (*MovableObjStaticInfo)(unsafe.Pointer(o.staticInfo))
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

// 等级
func (o MovableObject) Level() int32 {
	return 0
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

// 逆時針旋轉
func (o *MovableObject) Rotate(angle int32) {
	o.orientation += angle
	if o.orientation >= 360 {
		o.orientation -= 360
	}
}

// 逆時針旋轉到
func (o *MovableObject) RotateTo(angle int32) {
	o.orientation = angle
}

// 移动
func (o *MovableObject) Move(dir Direction) {
	if dir < DirMin || dir > DirMax {
		str := fmt.Sprintf("invalid object direction %v", dir)
		panic(str)
	}
	if o.state == stopped {
		o.dir = dir
		if !o.checkMove(dir, 0, 0) {
			return
		}
		o.state = toMove
		log.Debug("@@@ object %v stopped => to move", o.instId)
	}
}

// 移動了距離
func (o *MovableObject) MovedDistance(x, y int32) {
	o.x += x
	o.y += y
}

// 停止
func (o *MovableObject) Stop() {
	// 准备移动则直接停止
	if o.state == toMove {
		o.state = stopped
		log.Debug("@@@ object %v to move => stopped", o.instId)
		return
	}
	// 正在移动则准备停止
	if o.state == isMoving {
		o.state = toStop
		log.Debug("@@@ object %v moving => to stop", o.instId)
		return
	}
}

// 是否正在移动
func (o *MovableObject) IsMoving() bool {
	return o.state == toMove || o.state == isMoving || o.state == toStop
}

// 上次Update的位置
func (o MovableObject) LastPos() (int32, int32) {
	return o.lastX, o.lastY
}

// 更新
func (o *MovableObject) Update(tick time.Duration) {
	// 每次Update都要更新lastX和lastY
	o.lastX, o.lastY = o.Pos()

	if o.state == stopped {
		return
	}

	if o.state == toMove {
		o.state = isMoving
		// args[0]: object.Pos
		// args[1]: object.Direction
		// args[2]: int32
		o.moveEvent.Call(Pos{X: o.lastX, Y: o.lastY}, o.dir, o.CurrentSpeed())
		log.Debug("@@@ object %v to move => moving", o.instId)
		return
	}

	var x, y int32
	if o.MovableObjStaticInfo().MoveFunc != nil {
		if o.mySuper == nil {
			o.mySuper = o.super.(IMovableObject)
		}
		x, y = o.MovableObjStaticInfo().MoveFunc(o.mySuper, tick)
	} else {
		x, y = DefaultMove(o, tick)
	}

	if o.state != stopped {
		ox, oy := o.Pos()
		if o.checkMove(o.dir, float64(x-ox), float64(y-oy)) {
			o.SetPos(x, y)
		}
	} else {
		o.SetPos(x, y)
	}

	if o.state == isMoving {
		// args[0]: object.Pos
		// args[1]: object.Direction
		// args[2]: int32
		o.updateEvent.Call(Pos{X: x, Y: y}, o.dir, o.CurrentSpeed())
	} else if o.state == toStop {
		o.state = stopped
		// args[0]: object.Pos
		// args[1]: object.Direction
		// args[2]: int32
		o.stopEvent.Call(Pos{X: x, Y: y}, o.dir, o.CurrentSpeed())
		log.Debug("@@@ object %v to stop => stopped", o.instId)
	}
}

func (o *MovableObject) checkMove(dir Direction, dx, dy float64) bool {
	var (
		isMove, isCollision bool
		resObj              IObject
	)

	if dir == DirNone {
		return true
	}

	o.checkMoveEvent.Call(o.instId, dir, dx, dy, &isMove, &isCollision, &resObj)
	if isCollision {
		comp := o.GetComp("Collider")
		if comp != nil {
			collisionComp := comp.(*ColliderComp)
			collisionComp.CallCollisionEventHandle(o.super, resObj)
		}
	}
	if isMove {
		switch dir {
		case DirLeft:
			o.orientation = 180
		case DirRight:
			o.orientation = 0
		case DirUp:
			o.orientation = 90
		case DirDown:
			o.orientation = 270
		}
	}
	return isMove
}

// 注冊檢查坐標事件
func (o *MovableObject) RegisterCheckMoveEventHandle(handle func(args ...any)) {
	o.checkMoveEvent.Register(handle)
}

// 注銷檢查坐標事件
func (o *MovableObject) UnregisterCheckMoveEventHandle(handle func(args ...any)) {
	o.checkMoveEvent.Unregister(handle)
}

// 注册移动事件
func (o *MovableObject) RegisterMoveEventHandle(handle func(args ...any)) {
	o.moveEvent.Register(handle)
}

// 注销移动事件
func (o *MovableObject) UnregisterMoveEventHandle(handle func(args ...any)) {
	o.moveEvent.Unregister(handle)
}

// 注册停止移动事件
func (o *MovableObject) RegisterStopMoveEventHandle(handle func(args ...any)) {
	o.stopEvent.Register(handle)
}

// 注销停止移动事件
func (o *MovableObject) UnregisterStopMoveEventHandle(handle func(args ...any)) {
	o.stopEvent.Unregister(handle)
}

// 注册更新事件
func (o *MovableObject) RegisterUpdateEventHandle(handle func(args ...any)) {
	o.updateEvent.Register(handle)
}

// 注销更新事件
func (o *MovableObject) UnregisterUpdateEventHandle(handle func(args ...any)) {
	o.updateEvent.Unregister(handle)
}
