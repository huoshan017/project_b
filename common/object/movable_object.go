package object

import (
	"project_b/common/base"
	"project_b/log"
	"unsafe"
)

// 可移动物体状态
type moveObjectState int32

const (
	stopped moveObjectState = iota
	rotating
	toMove
	isMoving
	toStop
)

// 可移动的物体
type MovableObject struct {
	object
	moveDir        base.Angle      // 移動的方向角度
	speed          int32           // 当前移动速度（米/秒）
	lastX, lastY   int32           // 上次更新的位置
	state, pstate  moveObjectState // 移动状态和上一狀態
	mySuper        IMovableObject  // 父類
	checkMoveEvent base.Event      // 檢查坐標事件
	moveEvent      base.Event      // 移动事件
	stopEvent      base.Event      // 停止事件
	updateEvent    base.Event      // 更新事件
	pauseEvent     base.Event      // 暫停事件
	resumeEvent    base.Event      // 恢復事件
	pause          bool            // 是否暫停
	collisionInfo  CollisionInfo   // 碰撞信息
}

// 创建可移动物体
func NewMovableObject() *MovableObject {
	o := &MovableObject{}
	return o
}

// 在子類中創建可移動物體
func NewMovableObjectWithSuper(super IMovableObject) *MovableObject {
	mobj := NewMovableObject()
	mobj.mySuper = super
	return mobj
}

// 初始化
func (o *MovableObject) Init(instId uint32, staticInfo *ObjStaticInfo) {
	o.object.Init(instId, staticInfo)
	o.speed = staticInfo.speed
}

// 反初始化
func (o *MovableObject) Uninit() {
	o.speed = 0
	o.lastX = 0
	o.lastY = 0
	o.state = stopped
	o.pstate = stopped
	o.checkMoveEvent.Clear()
	o.moveEvent.Clear()
	o.stopEvent.Clear()
	o.updateEvent.Clear()
	o.object.Uninit()
	o.collisionInfo.Clear()
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

// 運動方向
func (o MovableObject) MoveDir() base.Angle {
	return o.moveDir
}

// 当前速度
func (o MovableObject) CurrentSpeed() int32 {
	return o.speed
}

// 逆時針旋轉
func (o *MovableObject) Rotate(angle base.Angle) {
	if o.pause {
		return
	}
	o.rotation.Add(angle)
}

// 逆時針旋轉到
func (o *MovableObject) RotateTo(angle base.Angle) {
	if o.pause {
		return
	}
	angle.Sub(base.NewAngle(int16(o.staticInfo.rotation), 0))
	o.rotation = angle
}

// 朝向向量
func (o *MovableObject) Forward() base.Vec2 {
	cx0, cy0 := o.Pos()
	cx1, cy1 := cx0+o.Length()/2, cy0
	cx2, cy2 := base.Rotate(cx1, cy1, cx0, cy0, o.Rotation())
	return base.NewVec2(cx2-cx0, cy2-cy0)
}

// 移动
func (o *MovableObject) Move(dir base.Angle) {
	if o.pause {
		return
	}
	o.moveDir = dir
	if o.state == stopped {
		var tick uint32 = 100
		d := GetDefaultLinearDistance(o, tick)
		v := dir.DistanceToVec2(d)
		if !o.checkMove(v.X(), v.Y(), false, nil, nil) {
			return
		}
		o.setState(toMove)
		log.Debug("@@@ object %v stopped => to move, moveDir %v, rotation %v, currMs %v", o.instId, o.moveDir, o.Rotation(), o.CurrMs())
	}
}

// 停止
func (o *MovableObject) Stop() {
	if o.pause {
		return
	}
	// 准备移动则直接停止
	if o.state == toMove {
		o.setState(stopped)
		log.Debug("@@@ object %v to move => stopped, moveDir %v, rotation %v, currMs %v", o.instId, o.moveDir, o.Rotation(), o.CurrMs())
		return
	}
	// 正在移动则准备停止
	if o.state == isMoving {
		o.setState(toStop)
		log.Debug("@@@ object %v moving => to stop, moveDir %v, rotation %v, currMs %v", o.instId, o.moveDir, o.Rotation(), o.CurrMs())
		return
	}
}

// 立即停止
func (o *MovableObject) StopNow() {
	if o.pause {
		return
	}
	if o.state == toMove || o.state == isMoving {
		o.setState(stopped)
		if o.pstate == toMove {
			log.Debug("@@@ stop now!!! object %v to move => stopped, moveDir %v, rotation %v, currMs %v", o.instId, o.moveDir, o.Rotation(), o.CurrMs())
		} else {
			log.Debug("@@@ stop now!!! object %v moving => stopped, moveDir %v, rotation %v, currMs %v", o.instId, o.moveDir, o.Rotation(), o.CurrMs())
		}
		return
	}
}

// 是否正在移动
func (o *MovableObject) IsMoving() bool {
	if o.pause {
		return false
	}
	return o.state == toMove || o.state == isMoving || o.state == toStop
}

// 設置位置
func (o *MovableObject) SetPos(x, y int32) {
	o.lastX, o.lastY = o.Pos()
	o.x, o.y = x, y
}

// 上次Update的位置
func (o MovableObject) LastPos() (int32, int32) {
	return o.lastX, o.lastY
}

// 暫停
func (o *MovableObject) Pause() {
	o.pause = true
	o.pauseEvent.Call()
}

// 繼續
func (o *MovableObject) Resume() {
	o.pause = false
	o.resumeEvent.Call()
}

// 更新
func (o *MovableObject) Update(tickMs uint32) {
	if o.pause {
		return
	}

	o.object.Update(tickMs)

	// 每次Update都要更新lastX和lastY
	o.lastX, o.lastY = o.Pos()

	if o.state == stopped {
		return
	}

	if o.state == toMove {
		o.setState(isMoving)
		o.moveEvent.Call(base.Pos{X: o.lastX, Y: o.lastY}, o.moveDir, o.CurrentSpeed())
		log.Debug("@@@ object %v to move => moving, moveDir %v, rotation %v, currMs %v", o.instId, o.moveDir, o.Rotation(), o.CurrMs())
		return
	}

	var x, y int32
	if o.MovableObjStaticInfo().MoveFunc != nil {
		x, y = o.MovableObjStaticInfo().MoveFunc(o.mySuper, tickMs)
	} else {
		x, y = DefaultMove(o.mySuper, tickMs)
	}

	ox, oy := o.Pos()
	o.checkMove(x-ox, y-oy, true, func() {
		o.SetPos(x, y)
	}, func() {
		o.setState(toStop)
	})

	if o.state == isMoving {
		o.updateEvent.Call(base.Pos{X: x, Y: y}, o.moveDir, o.CurrentSpeed())
	} else if o.state == toStop {
		o.setState(stopped)
		o.stopEvent.Call(base.Pos{X: x, Y: y}, o.moveDir, o.CurrentSpeed())
		log.Debug("@@@ object %v to stop => stopped, moveDir %v, rotation %v, currMs %v", o.instId, o.moveDir, o.Rotation(), o.CurrMs())
	}
}

// 檢測是否移動
func (o *MovableObject) checkMove(dx, dy int32, update bool, canMoveFunc func(), cantMoveFunc func()) bool {
	o.collisionInfo.Clear()
	o.checkMoveEvent.Call(o.instId, dx, dy, &o.collisionInfo)
	if o.collisionInfo.Result != CollisionNone {
		if o.collisionInfo.Result == CollisionAndBlock {
			if update {
				o.SetPos(o.collisionInfo.MovingObjPos.X, o.collisionInfo.MovingObjPos.Y)
			}
			if cantMoveFunc != nil {
				cantMoveFunc()
			}
		} else {
			if canMoveFunc != nil {
				canMoveFunc()
			}
		}
		if o.colliderComp != nil {
			o.colliderComp.CallCollisionEventHandle(o.mySuper, &o.collisionInfo)
		}
	} else {
		if canMoveFunc != nil {
			canMoveFunc()
		}
	}
	return o.collisionInfo.Result != CollisionAndBlock
}

// 設置狀態
func (o *MovableObject) setState(newState moveObjectState) {
	o.pstate = o.state
	o.state = newState
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

// 注冊暫停事件
func (o *MovableObject) RegisterPauseEventHandle(handle func(args ...any)) {
	o.pauseEvent.Register(handle)
}

// 注銷暫停事件
func (o *MovableObject) UnregisterPauseEventHandle(handle func(args ...any)) {
	o.pauseEvent.Unregister(handle)
}

// 注冊恢復事件
func (o *MovableObject) RegisterResumeEventHandle(handle func(args ...any)) {
	o.resumeEvent.Register(handle)
}

// 注銷恢復事件
func (o *MovableObject) UnregisterResumeEventHandle(handle func(args ...any)) {
	o.resumeEvent.Unregister(handle)
}
