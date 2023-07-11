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
	level           int32
	changeEvent     base.Event
	fireTime        time.CustomTime
	bulletFireCount int8
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
func (t *Tank) CheckAndFire(newBulletFunc func(*BulletStaticInfo) *Bullet, bulletInfo *BulletStaticInfo) *Bullet {
	var (
		bullet     *Bullet
		staticInfo = t.TankStaticInfo()
	)
	// 先檢測炮彈冷卻時間
	if t.fireTime.IsZero() || time.Since(t.fireTime) >= time.Duration(staticInfo.BulletConfig.Cooldown)*time.Millisecond {
		bullet = newBulletFunc(bulletInfo)
		t.fireTime = time.Now()
		t.bulletFireCount = 1
	}
	// 再檢測一次發射中的炮彈間隔
	/*if t.bulletConfig.AmountFireOneTime > 1 && t.bulletFireCount < t.bulletConfig.AmountFireOneTime {
		if t.fireIntervalTime.IsZero() || time.Since(t.fireIntervalTime) >= time.Duration(t.bulletConfig.IntervalInFire)*time.Millisecond {
			if bullet == nil {
				bullet = newBulletFunc(bulletInfo)
			}
			t.fireIntervalTime = time.Now()
			t.bulletFireCount += 1
		}
	}*/
	if bullet != nil {
		cx, cy := t.Pos()
		tl := t.Length()
		bl := bullet.Length()
		switch t.dir {
		case DirLeft:
			bullet.SetPos(cx-tl/2-bl/2-1, cy)
		case DirRight:
			bullet.SetPos(cx+tl/2+bl/2+1, cy)
		case DirUp:
			bullet.SetPos(cx, cy+tl/2+bl/2+1)
		case DirDown:
			bullet.SetPos(cx, cy-tl/2-bl/2-1)
		}
		bullet.SetCamp(t.currentCamp)
		bullet.SetCurrentSpeed(bulletInfo.speed)
	}
	return bullet
}
