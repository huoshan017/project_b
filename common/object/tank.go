package object

import (
	"project_b/common/base"
	"project_b/common/time"
	"project_b/log"
	"unsafe"
)

// 车辆
type Vehicle struct {
	MovableObject
}

// 创建车辆
func NewVehicle(instId uint32, staticInfo *ObjStaticInfo) *Vehicle {
	o := &Vehicle{}
	o.Init(instId, staticInfo)
	return o
}

// 初始化
func (v *Vehicle) Init(instId uint32, staticInfo *ObjStaticInfo) {
	v.MovableObject.Init(instId, staticInfo)
	v.setSuper(v)
}

// 反初始化
func (v *Vehicle) Uninit() {
	v.MovableObject.Uninit()
}

// 護盾
type Shield struct {
	tank        *Tank
	staticInfo  *TankShieldStaticInfo
	startTime   time.CustomTime
	pause       bool
	isEffective bool
}

// 創建護盾
func NewShield(tank *Tank, staticInfo *TankShieldStaticInfo) *Shield {
	shield := &Shield{
		tank:       tank,
		staticInfo: staticInfo,
	}
	if staticInfo.Duration > 0 {
		shield.startTime = time.Now()
	}
	shield.isEffective = true
	return shield
}

// 更新
func (s *Shield) Update(tick time.Duration) {
	if s.pause {
		return
	}
	if s.staticInfo.Duration > 0 && time.Since(s.startTime) > s.staticInfo.Duration {
		s.isEffective = false
		return
	}
}

// 暫停
func (s *Shield) Pause() {
	s.pause = true
}

// 恢復
func (s *Shield) Resume() {
	s.pause = false
}

// 是否有效果
func (s *Shield) IsEffective() bool {
	return s.isEffective
}

// 坦克
type Tank struct {
	Vehicle
	level                      int32
	changeEvent                base.Event
	fireTime, fireIntervalTime time.CustomTime
	shellFireCount             int8
	shield                     *Shield
	addShieldEvent             base.Event
	cancelShieldEvent          base.Event
}

// 创建坦克
func NewTank(instId uint32, staticInfo *TankStaticInfo) *Tank {
	tank := &Tank{}
	tank.Init(instId, &staticInfo.ObjStaticInfo)
	return tank
}

// 靜態配置
func (t *Tank) TankStaticInfo() *TankStaticInfo {
	return (*TankStaticInfo)(unsafe.Pointer(t.staticInfo))
}

// 初始化
func (t *Tank) Init(instId uint32, staticInfo *ObjStaticInfo) {
	t.Vehicle.Init(instId, staticInfo)
	t.level = t.TankStaticInfo().Level
	t.setSuper(t)
}

// 反初始化
func (t *Tank) Uninit() {
	t.shield = nil
	t.Vehicle.Uninit()
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
func (t *Tank) Change(info *TankStaticInfo) {
	t.ChangeStaticInfo(&info.ObjStaticInfo)
	t.SetCurrentSpeed(info.Speed())
	t.changeEvent.Call(info, t.level)
}

// 还原
func (t *Tank) Restore() {
	if t.pause {
		return
	}
	t.RestoreStaticInfo()
	t.changeEvent.Call(t.staticInfo, t.level)
}

// 注册变化事件
func (t *Tank) RegisterChangeEventHandle(handle func(args ...any)) {
	t.changeEvent.Register(handle)
}

// 注销变化事件
func (t *Tank) UnregisterChangeEventHandle(handle func(args ...any)) {
	t.changeEvent.Unregister(handle)
}

// 注冊加護盾事件
func (t *Tank) RegisterAddShieldEventHandle(handle func(args ...any)) {
	t.addShieldEvent.Register(handle)
}

// 注銷加護盾事件
func (t *Tank) UnregisterAddShieldEventHandle(handle func(args ...any)) {
	t.addShieldEvent.Unregister(handle)
}

// 注冊取消護盾事件
func (t *Tank) RegisterCancelShieldEventHandle(handle func(args ...any)) {
	t.cancelShieldEvent.Register(handle)
}

// 注銷取消護盾事件
func (t *Tank) UnregisterCancelShieldEventHandle(handle func(args ...any)) {
	t.cancelShieldEvent.Unregister(handle)
}

// 檢測是否可以開炮
func (t *Tank) CheckAndFire(newShellFunc func(*ShellStaticInfo) *Shell, shellInfo *ShellStaticInfo) *Shell {
	var (
		shell      *Shell
		staticInfo = t.TankStaticInfo()
	)
	// 先檢測炮彈冷卻時間
	if t.fireTime.IsZero() || time.Since(t.fireTime) >= time.Duration(staticInfo.ShellConfig.Cooldown)*time.Millisecond {
		shell = newShellFunc(shellInfo)
		t.fireTime = time.Now()
		t.shellFireCount = 1
	}
	// 再檢測一次發射中的炮彈間隔
	if t.TankStaticInfo().ShellConfig.AmountFireOneTime > 1 && t.shellFireCount < t.TankStaticInfo().ShellConfig.AmountFireOneTime {
		if t.fireIntervalTime.IsZero() || time.Since(t.fireIntervalTime) >= time.Duration(t.TankStaticInfo().ShellConfig.IntervalInFire)*time.Millisecond {
			if shell == nil {
				shell = newShellFunc(shellInfo)
			}
			t.fireIntervalTime = time.Now()
			t.shellFireCount += 1
		}
	}
	if shell != nil {
		x, y := t.shellLaunchPos(shell)
		shell.SetPos(x, y)
		shell.SetCamp(t.currentCamp)
		shell.SetCurrentSpeed(shell.speed)
	}
	return shell
}

// 移動
func (t *Tank) Move(dir base.Angle) {
	if t.pause {
		return
	}
	t.moveDir = dir
	if t.moveDir != t.Rotation() || t.state == rotating {
		t.state = rotating
		log.Debug("@@@ tank %v rotating", t.instId)
		return
	}
	t.MovableObject.Move(dir)
}

// 停止
func (t *Tank) Stop() {
	if t.pause {
		return
	}
	if t.state == rotating {
		t.state = stopped
		x, y := t.Pos()
		t.stopEvent.Call(Pos{X: x, Y: y}, t.moveDir, t.CurrentSpeed())
		log.Debug("@@@ tank %v rotating => stopped", t.instId)
		return
	}
	t.MovableObject.Stop()
}

// 炮彈更新
func (t *Tank) Update(tick time.Duration) {
	if t.pause {
		return
	}

	if t.checkRotateState(tick) {
		return
	}
	t.MovableObject.Update(tick)
}

// 加護盾
func (t *Tank) AddShield(staticInfo *TankShieldStaticInfo) {
	if t.shield == nil {
		t.shield = NewShield(t, staticInfo)
		t.addShieldEvent.Call()
	}
}

// 取消護盾
func (t *Tank) CancelShield() {
	if t.shield != nil {
		t.shield = nil // todo 暫時先用垃圾回收處理，有需要的時候再考慮對象池或者復用
		t.cancelShieldEvent.Call()
	}
}

// 是否有護盾
func (t *Tank) HasShield() bool {
	return t.shield != nil
}

// 炮彈發射口
func (t *Tank) shellLaunchPos(shell *Shell) (int32, int32) {
	vp := t.TankStaticInfo().ShellLaunchPos
	x, y := t.Pos()
	x1, y1 := x+vp.X()+shell.Length()>>1, y+vp.Y()
	return base.Rotate(x1, y1, x, y, t.Rotation())
}

// 檢測旋轉狀態
func (t *Tank) checkRotateState(tick time.Duration) bool {
	if t.state != rotating {
		return false
	}

	tickMinutes := time.Duration(t.TankStaticInfo().SteeringAngularVelocity) * tick / time.Second
	var tickAngle base.Angle
	tickAngle.Set(int16(tickMinutes))
	tickAngle.Normalize()

	angleDiff := base.AngleSub(t.moveDir, t.Rotation())

	// 把角度差限制在[-180, 180]範圍内
	if angleDiff.Greater(base.PiAngle()) {
		angleDiff.Add(base.TwoPiAngle().Negative())
	} else if angleDiff.Less(base.PiAngle().Negative()) {
		angleDiff.Add(base.TwoPiAngle())
	}

	// 角度差的絕對值小於等於tick時間内的角度變化
	if (!angleDiff.IsNegative() && angleDiff.LessEqual(tickAngle)) || (angleDiff.IsNegative() && angleDiff.GreaterEqual(tickAngle.Negative())) {
		t.RotateTo(t.moveDir)
		t.state = toMove
		return false
	}

	// 角度差大於tick時間的角度變化量
	if angleDiff.Greater(tickAngle) {
		t.Rotate(tickAngle)
	} else {
		t.Rotate(tickAngle.Negative())
	}
	t.state = rotating
	return true
}
