package object

import (
	"fmt"
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
	orientation       int32          // 旋轉角度，以Y軸正方向為0度，逆時針方向為旋轉正方向
	components        []IComponent   // 組件
	changedStaticInfo *ObjStaticInfo // 改变的静态常量数据
	toRecycle         bool           // 去回收
	super             IObject        // 派生類對象
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
	o.toRecycle = false
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
	o.x = x
	o.y = y
}

// 中心點坐標
func (o object) Center() (x, y int32) {
	if o.changedStaticInfo != nil {
		return o.x + o.changedStaticInfo.w>>1, o.y + o.changedStaticInfo.h>>1
	}
	return o.x + o.staticInfo.w>>1, o.y + o.staticInfo.h>>1
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

// 朝向角度
func (o object) Orientation() int32 {
	return o.orientation
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
	checkMoveEvent *base.Event     // 檢查坐標事件
	moveEvent      *base.Event     // 移动事件
	stopEvent      *base.Event     // 停止事件
	updateEvent    *base.Event     // 更新事件
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
	o.checkMoveEvent = base.NewEvent()
	o.moveEvent = base.NewEvent()
	o.stopEvent = base.NewEvent()
	o.updateEvent = base.NewEvent()
	o.setSuper(o)
}

// 反初始化
func (o *MovableObject) Uninit() {
	o.object.Uninit()
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

// 移动
func (o *MovableObject) Move(dir Direction) {
	if dir < DirMin || dir > DirMax {
		str := fmt.Sprintf("invalid object direction %v", dir)
		panic(str)
	}
	if o.state == stopped {
		o.state = toMove
		o.dir = dir
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
	return o.state == isMoving || o.state == toStop
}

// 上次Update的位置
func (o MovableObject) LastPos() (int32, int32) {
	return o.lastX, o.lastY
}

// 獲得移動距離
func GetMoveDistance(obj IMovableObject, duration time.Duration) float64 {
	return float64(int64(obj.CurrentSpeed())*int64(duration)) / float64(time.Second)
}

// 更新
func (o *MovableObject) Update(tick time.Duration) {
	// 每次Update都要更新lastX和lastY
	o.lastX, o.lastY = o.x, o.y

	if o.state == stopped {
		return
	}

	if o.state == toMove {
		o.state = isMoving
		// args[0]: object.Pos
		// args[1]: object.Direction
		// args[2]: int32
		if o.moveEvent != nil {
			o.moveEvent.Call(Pos{X: o.x, Y: o.y}, o.dir, o.CurrentSpeed())
		}
		log.Debug("@@@ object %v to move => moving", o.instId)
		return
	}

	x, y := float64(o.x), float64(o.y)
	distance := GetMoveDistance(o, tick)
	switch o.dir {
	case DirLeft:
		x -= distance
	case DirRight:
		x += distance
	case DirUp:
		y += distance
	case DirDown:
		y -= distance
	default:
		panic("invalid direction")
	}

	if o.state != stopped && o.checkMoveEvent != nil {
		var (
			isMove, isCollision bool
			resObj              IObject
		)
		o.checkMoveEvent.Call(o.instId, o.dir, distance, &isMove, &isCollision, &resObj)
		if isMove {
			o.x, o.y = int32(x), int32(y)
		}
		if isCollision {
			comp := o.GetComp("Collider")
			if comp != nil {
				collisionComp := comp.(*ColliderComp)
				collisionComp.CallCollisionEventHandle(o.super, resObj)
			}
		}
	} else {
		o.x, o.y = int32(x), int32(y)
	}

	if o.state == isMoving {
		// args[0]: object.Pos
		// args[1]: object.Direction
		// args[2]: int32
		if o.updateEvent != nil {
			o.updateEvent.Call(Pos{X: o.x, Y: o.y}, o.dir, o.CurrentSpeed())
		}
	} else if o.state == toStop {
		o.state = stopped
		// args[0]: object.Pos
		// args[1]: object.Direction
		// args[2]: int32
		if o.stopEvent != nil {
			o.stopEvent.Call(Pos{X: o.x, Y: o.y}, o.dir, o.CurrentSpeed())
		}
		log.Debug("@@@ object %v to stop => stopped", o.instId)
	}
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
	bulletConfig     *TankBulletConfig
	level            int32
	changeEvent      *base.Event
	fireTime         time.CustomTime
	fireIntervalTime time.CustomTime
	bulletFireCount  int8
}

// 创建坦克
func NewTank(instId uint32, staticInfo *TankStaticInfo) *Tank {
	tank := &Tank{}
	tank.Init(instId, &staticInfo.ObjStaticInfo)
	return tank
}

// 初始化
func (t *Tank) Init(instId uint32, staticInfo *ObjStaticInfo) {
	t.Vehicle.Init(instId, staticInfo)
	tankStaticInfo := (*TankStaticInfo)(unsafe.Pointer(staticInfo))
	t.bulletConfig = &tankStaticInfo.BulletConfig
	t.level = tankStaticInfo.Level
	t.changeEvent = base.NewEvent()
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
	var bullet *Bullet
	// 先檢測炮彈冷卻時間
	if t.fireTime.IsZero() || time.Since(t.fireTime) >= time.Duration(t.bulletConfig.Cooldown)*time.Millisecond {
		bullet = newBulletFunc(bulletInfo)
		t.fireTime = time.Now()
		t.bulletFireCount = 1
	}
	// 再檢測一次發射中的炮彈間隔
	if t.bulletConfig.AmountFireOneTime > 1 && t.bulletFireCount < t.bulletConfig.AmountFireOneTime {
		if t.fireIntervalTime.IsZero() || time.Since(t.fireIntervalTime) >= time.Duration(t.bulletConfig.IntervalInFire)*time.Millisecond {
			if bullet == nil {
				bullet = newBulletFunc(bulletInfo)
			}
			t.fireIntervalTime = time.Now()
			t.bulletFireCount += 1
		}
	}
	if bullet != nil {
		switch t.dir {
		case DirLeft:
			bullet.SetPos(t.Left()-bullet.Height()-1, t.Top()-t.Width()>>1-bullet.Width()>>1)
		case DirRight:
			bullet.SetPos(t.Right()+1, t.Top()-t.Width()>>1-bullet.Width()>>1)
		case DirUp:
			bullet.SetPos(t.Left()+t.Width()>>1-bullet.Width()>>1, t.Top()+1)
		case DirDown:
			bullet.SetPos(t.Left()+t.Width()>>1-bullet.Width()>>1, t.Bottom()-bullet.Height()-1)
		}
		bullet.SetCurrentSpeed(bulletInfo.speed)
	}
	return bullet
}

// 獲得炮彈配置
func (t *Tank) GetBulletConfig() *TankBulletConfig {
	return t.bulletConfig
}

// 子弹
type Bullet struct {
	MovableObject
	info         *BulletStaticInfo
	explodeEvent base.Event
}

// 创建车辆
func NewBullet(instId uint32, staticInfo *BulletStaticInfo) *Bullet {
	o := &Bullet{}
	o.Init(instId, &staticInfo.ObjStaticInfo)
	return o
}

// 初始化
func (b *Bullet) Init(instId uint32, staticInfo *ObjStaticInfo) {
	b.MovableObject.Init(instId, staticInfo)
	b.info = (*BulletStaticInfo)(unsafe.Pointer(staticInfo))
	b.setSuper(b)
}

// 反初始化
func (b *Bullet) Uninit() {
	b.MovableObject.Uninit()
}

// 爆炸
func (b *Bullet) Explode() {
	b.explodeEvent.Call()
}

// 注冊爆炸事件
func (b *Bullet) RegisterExplodeEventHandle(handle func(...any)) {
	b.explodeEvent.Register(handle)
}

// 注銷爆炸事件
func (b *Bullet) UnregisterExplodeEventHandle(handle func(...any)) {
	b.explodeEvent.Unregister(handle)
}
