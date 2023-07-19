package object

import (
	"project_b/common/base"
	"project_b/common/log"
	"project_b/common/time"
	"unsafe"
)

// 可移动物体状态
type moveObjectState int32

const (
	stopped  = moveObjectState(0)
	rotating = moveObjectState(1)
	toMove   = moveObjectState(2)
	isMoving = moveObjectState(3)
	toStop   = moveObjectState(4)
)

// 可移动的物体
type MovableObject struct {
	object
	speed          int32           // 当前移动速度（米/秒）
	lastX, lastY   int32           // 上次更新的位置
	lastTick       time.Duration   // 上次tick花費時間
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
	o.speed = staticInfo.speed
	o.setSuper(o)
}

// 反初始化
func (o *MovableObject) Uninit() {
	o.speed = 0
	o.lastX = 0
	o.lastY = 0
	o.lastTick = 0
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

// 设置当前速度
func (o *MovableObject) SetCurrentSpeed(speed int32) {
	o.speed = speed
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
func (o *MovableObject) Rotate(angle base.Angle) {
	o.orientation.Add(angle)
}

// 逆時針旋轉到
func (o *MovableObject) RotateTo(angle base.Angle) {
	o.orientation = angle
}

// 移动
func (o *MovableObject) Move( /*dir Direction*/ dir base.Angle) {
	if o.state == stopped {
		d := GetDefaultLinearDistance(o, o.lastTick)
		o.orientation = dir
		v := dir.DistanceToVec2(d)
		if !o.checkMove( /*dir, */ v.X(), v.Y()) {
			return
		}
		o.state = toMove
		log.Debug("@@@ object %v stopped => to move", o.instId)
	}
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
	o.lastTick = tick

	if o.state == stopped {
		return
	}

	if o.state == toMove {
		o.state = isMoving
		o.moveEvent.Call(Pos{X: o.lastX, Y: o.lastY}, o.orientation /*o.dir */, o.CurrentSpeed())
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
		if o.checkMove( /*o.dir, */ x-ox, y-oy) {
			o.SetPos(x, y)
		}
	} else {
		o.SetPos(x, y)
	}

	if o.state == isMoving {
		o.updateEvent.Call(Pos{X: x, Y: y}, o.orientation, o.CurrentSpeed())
	} else if o.state == toStop {
		o.state = stopped
		o.stopEvent.Call(Pos{X: x, Y: y}, o.orientation, o.CurrentSpeed())
		log.Debug("@@@ object %v to stop => stopped", o.instId)
	}
}

func (o *MovableObject) checkMove( /*dir Direction, */ dx, dy int32) bool {
	var (
		isMove, isCollision bool = true, false
		resObj              IObject
	)

	o.checkMoveEvent.Call(o.instId /*dir, */, dx, dy, &isMove, &isCollision, &resObj)
	if isCollision {
		comp := o.GetComp("Collider")
		if comp != nil {
			collisionComp := comp.(*ColliderComp)
			collisionComp.CallCollisionEventHandle(o.super, resObj)
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
