package object

import (
	"project_b/common/base"
	"project_b/common/weapon"
	"project_b/log"
	"unsafe"
)

// 车辆
type Vehicle struct {
	*MovableObject
}

// 创建车辆
func NewVehicle() *Vehicle {
	o := &Vehicle{}
	o.MovableObject = NewMovableObjectWithSuper(o)
	return o
}

// 通過子類創建車輛
func NewVehicleWithSuper(super IMovableObject) *Vehicle {
	v := &Vehicle{}
	v.MovableObject = NewMovableObjectWithSuper(super)
	return v
}

// 初始化
func (v *Vehicle) Init(instId uint32, staticInfo *ObjStaticInfo) {
	v.MovableObject.Init(instId, staticInfo)
}

// 反初始化
func (v *Vehicle) Uninit() {
	v.MovableObject.Uninit()
}

// 護盾
type Shield struct {
	tank        *Tank
	staticInfo  *TankShieldStaticInfo
	startMs     uint32
	totalMs     uint32
	pause       bool
	isEffective bool
}

// 創建護盾
func NewShield(tank *Tank, staticInfo *TankShieldStaticInfo) *Shield {
	shield := &Shield{
		tank:       tank,
		staticInfo: staticInfo,
	}
	shield.isEffective = true
	return shield
}

// 更新
func (s *Shield) Update(tickMs uint32) {
	if s.pause {
		return
	}
	if s.startMs == 0 {
		s.startMs = tickMs
	}
	s.totalMs += tickMs
	if s.staticInfo.DurationMs > 0 && s.totalMs-s.startMs > s.staticInfo.DurationMs {
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

// 激光持有者
type laserHolder struct {
	tank *Tank
}

func (holder *laserHolder) LaunchPoint() base.Pos {
	vp := holder.tank.TankStaticInfo().ShellLaunchPos
	x, y := holder.tank.Pos()
	x1, y1 := x+vp.X(), y+vp.Y()
	return base.NewPos(base.Rotate(x1, y1, x, y, holder.tank.Rotation()))
}

// 朝向
func (holder *laserHolder) Forward() base.Vec2 {
	return holder.tank.Forward()
}

// 陣營
func (holder *laserHolder) Camp() base.CampType {
	return holder.tank.Camp()
}

// 坦克
type Tank struct {
	*Vehicle
	level                          int32
	changeEvent                    base.Event
	fireTimeMs, fireIntervalTimeMs uint32
	shellFireCount                 int8
	shellStaticInfoList            []*ShellStaticInfo
	shellIndex                     int32
	shield                         *Shield
	laser                          *weapon.Laser
	laserHolder                    laserHolder
	addShieldEvent                 base.Event
	cancelShieldEvent              base.Event
}

// 创建坦克
func NewTank() *Tank {
	tank := &Tank{}
	tank.Vehicle = NewVehicleWithSuper(tank)
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
	shellConfig := &t.TankStaticInfo().ShellConfig
	if shellConfig.ShellInfo != nil {
		t.shellStaticInfoList = []*ShellStaticInfo{shellConfig.ShellInfo}
	}
}

// 反初始化
func (t *Tank) Uninit() {
	t.level = 0
	if len(t.shellStaticInfoList) > 0 {
		t.shellStaticInfoList = t.shellStaticInfoList[:0]
	}
	t.shield = nil
	t.shellIndex = 0
	t.fireTimeMs = 0
	t.fireIntervalTimeMs = 0
	t.shellFireCount = 0
	t.changeEvent.Clear()
	t.addShieldEvent.Clear()
	t.cancelShieldEvent.Clear()
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
func (t *Tank) CheckAndFire(newShellFunc func(*ShellStaticInfo) *Shell) *Shell {
	var (
		shell      *Shell
		staticInfo = t.TankStaticInfo()
	)
	if len(t.shellStaticInfoList) == 0 {
		return nil
	}
	// 先檢測炮彈冷卻時間
	currMs := t.CurrMs()
	if t.fireTimeMs == 0 || int32(currMs-t.fireTimeMs) >= staticInfo.ShellConfig.CooldownMs {
		shell = newShellFunc(t.shellStaticInfoList[t.shellIndex])
		t.fireTimeMs = currMs
		t.shellFireCount = 1
	}
	// 再檢測一次發射中的炮彈間隔
	if t.TankStaticInfo().ShellConfig.AmountFireOneTime > 1 && t.shellFireCount < t.TankStaticInfo().ShellConfig.AmountFireOneTime {
		if t.fireIntervalTimeMs == 0 || int32(currMs-t.fireIntervalTimeMs) >= t.TankStaticInfo().ShellConfig.IntervalInFireMs {
			if shell == nil {
				shell = newShellFunc(t.shellStaticInfoList[t.shellIndex])
			}
			t.fireIntervalTimeMs = currMs
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

// 發射激光
func (t *Tank) LaunchLaser(laserStaticInfo *weapon.LaserStaticInfo) *weapon.Laser {
	if t.laser == nil {
		t.laserHolder.tank = t
		t.laser = weapon.NewLaser(&t.laserHolder, laserStaticInfo)
	}
	return t.laser
}

// 取消發射激光
func (t *Tank) CancelLaser() {
	if t.laser == nil {
		return
	}
	t.laser.Cancel()
}

// 移動
func (t *Tank) Move(dir base.Angle) {
	if t.pause {
		return
	}
	// toStop狀態不能被打斷
	if t.state == toStop {
		return
	}
	t.moveDir = dir
	if t.moveDir != t.Rotation() || t.state == rotating {
		t.setState(rotating)
		log.Debug("@@@ tank %v rotating, moveDir %v, rotation %v, currMs %v", t.instId, t.moveDir, t.Rotation(), t.CurrMs())
		return
	}
	t.MovableObject.Move(dir)
}

// 停止
func (t *Tank) Stop() {
	if t.pause {
		return
	}
	// rotating狀態不能被打斷
	if t.state == rotating {
		/*t.setState(toStop)
		x, y := t.Pos()
		t.stopEvent.Call(Pos{X: x, Y: y}, t.moveDir, t.CurrentSpeed())
		log.Debug("@@@ tank %v rotating => to stop, moveDir %v, rotation %v, currMs %v", t.instId, t.moveDir, t.Rotation(), t.CurrMs())*/
		return
	}
	t.MovableObject.Stop()
}

// 炮彈更新
func (t *Tank) Update(tickMs uint32) {
	if t.pause {
		return
	}

	if t.checkRotateState(tickMs) {
		return
	}
	t.MovableObject.Update(tickMs)
	if t.shield != nil {
		t.shield.Update(tickMs)
	}
	if t.laser != nil {
		t.laser.Update(tickMs)
	}
}

// 添加彈藥
func (t *Tank) AppendShell(shellStaticInfo *ShellStaticInfo) bool {
	for i := 0; i < len(t.shellStaticInfoList); i++ {
		if t.shellStaticInfoList[i].Id() == shellStaticInfo.Id() {
			return false
		}
	}
	t.shellStaticInfoList = append(t.shellStaticInfoList, shellStaticInfo)
	t.shellIndex = int32(len(t.shellStaticInfoList)) - 1
	return true
}

// 切換彈藥
func (t *Tank) SwitchShell() {
	if len(t.shellStaticInfoList) > 0 {
		t.shellIndex = (t.shellIndex + 1) % int32(len(t.shellStaticInfoList))
	}
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

// 激光
func (t *Tank) GetLaser() *weapon.Laser {
	return t.laser
}

// 炮彈發射點
func (t *Tank) shellLaunchPos(shell *Shell) (int32, int32) {
	vp := t.TankStaticInfo().ShellLaunchPos
	x, y := t.Pos()
	x1, y1 := x+vp.X()+shell.Length()>>1, y+vp.Y()
	return base.Rotate(x1, y1, x, y, t.Rotation())
}

// 檢測旋轉狀態
func (t *Tank) checkRotateState(tickMs uint32) bool {
	if t.state != rotating {
		return false
	}

	tickMinutes := uint32(t.TankStaticInfo().SteeringAngularVelocity) * tickMs / 1000
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
		// todo 旋轉結束設置坦克狀態為停止，如果設置為toStop,toMove,moving則會出現漂移現象
		t.setState(stopped)
		return false
	}

	// 角度差大於tick時間的角度變化量
	if angleDiff.Greater(tickAngle) {
		t.Rotate(tickAngle)
	} else {
		t.Rotate(tickAngle.Negative())
	}
	t.setState(rotating)
	return true
}
