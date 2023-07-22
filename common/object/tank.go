package object

import (
	"project_b/common/base"
	"project_b/common/time"
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

// 坦克
type Tank struct {
	Vehicle
	level                      int32
	changeEvent                base.Event
	fireTime, fireIntervalTime time.CustomTime
	shellFireCount             int8
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
		x, y := t.shellLaunchPos()
		shell.SetPos(x, y)
		shell.SetCamp(t.currentCamp)
		shell.SetCurrentSpeed(shell.speed)
	}
	return shell
}

// 炮彈發射口
func (t *Tank) shellLaunchPos() (int32, int32) {
	vp := t.TankStaticInfo().ShellLaunchPos
	x, y := t.Pos()
	x1, y1 := x+vp.X(), y+vp.Y()
	rotation := t.Rotation()
	cn, cd := base.Cosine(rotation)
	sn, sd := base.Sine(rotation)
	xa := (x1-x)*cn/cd - (y1-y)*sn/sd + x
	ya := (x1-x)*sn/sd + (y1-y)*cn/cd + y
	return xa, ya
}
