package object

import (
	"project_b/common/base"
	"project_b/common/log"
	"project_b/common/time"
)

// 獲得移動距離
func GetDefaultLinearDistance(obj IMovableObject, duration time.Duration) int32 {
	return int32((int64(obj.CurrentSpeed()) * int64(duration)) / int64(time.Second))
}

// 默認移動，就是直綫移動
func DefaultMove(mobj IMovableObject, tick time.Duration) (int32, int32) {
	distance := GetDefaultLinearDistance(mobj, tick)
	dir := mobj.MoveDir()
	sn, sd := base.Sine(dir)
	cn, cd := base.Cosine(dir)
	dx := distance * cn / cd
	dy := distance * sn / sd
	x, y := mobj.Pos()
	return x + dx, y + dy
}

// 環繞物移動
func SurroundObjMove(mobj IMovableObject, tick time.Duration) (int32, int32) {
	var sobj = mobj.(*SurroundObj)
	return getSurroundObjMovedPos(sobj, tick, nil)
}

// 環繞物移動信息
type SurroundMoveInfo struct {
	LastCenterX, LastCenterY int32
	TurnAngle                int32 // 單位是分，1/60度
	AccumulateTime           time.Duration
}

// 獲得環繞物移動位置
func GetSurroundObjMovedPos(sobj *SurroundObj, tick time.Duration, moveInfo *SurroundMoveInfo) (int32, int32) {
	return getSurroundObjMovedPos(sobj, tick, moveInfo)
}

func getSurroundObjMovedPos(sobj *SurroundObj, tick time.Duration, moveInfo *SurroundMoveInfo) (x, y int32) {
	aroundCenterObj := sobj.getAroundCenterObjFunc(sobj.aroundCenterObjInstId)
	if aroundCenterObj == nil {
		return 0, 0
	}

	var (
		turnAngle      int32
		accumulateTime time.Duration
		cx, cy         int32
	)
	if moveInfo != nil {
		turnAngle = moveInfo.TurnAngle
		accumulateTime = moveInfo.AccumulateTime
		cx, cy = moveInfo.LastCenterX, moveInfo.LastCenterY
	} else {
		turnAngle = sobj.turnAngle
		accumulateTime = sobj.accumulateTime
		cx, cy = aroundCenterObj.Pos()
	}
	accumulateTime += tick

	staticInfo := sobj.SurroundObjStaticInfo()
	angle := int32(accumulateTime * time.Duration(staticInfo.AngularVelocity) / time.Second)
	turnAngle += angle
	degree, minute := turnAngle/60, turnAngle%60
	if degree >= 360 {
		degree -= 360
	}
	accumulateTime -= time.Duration(angle) * time.Second / time.Duration(staticInfo.AngularVelocity)
	an := base.NewAngleObj(int16(degree), int16(minute))
	sn, sd := base.Sine(an)
	cn, cd := base.Cosine(an)
	if staticInfo.Clockwise {
		x, y = cx+staticInfo.AroundRadius*cn/cd, cy-staticInfo.AroundRadius*sn/sd
	} else {
		x, y = cx+staticInfo.AroundRadius*cn/cd, cy+staticInfo.AroundRadius*sn/sd
	}
	if moveInfo != nil {
		moveInfo.TurnAngle = degree*60 + minute
		moveInfo.AccumulateTime = accumulateTime
	} else {
		sobj.turnAngle = degree*60 + minute
		sobj.accumulateTime = accumulateTime
	}
	return x, y
}

// 跟蹤移動
func ShellTrackMove(mobj IMovableObject, tick time.Duration) (int32, int32) {
	if mobj.Subtype() != ObjSubtypeShell {
		return DefaultMove(mobj, tick)
	}

	shell := mobj.(*Shell)
	if shell.ShellStaticInfo().SteeringAngularVelocity <= 0 {
		return DefaultMove(mobj, tick)
	}

	var target IObject
	if shell.trackTargetId == 0 {
		target = shell.searchTargetFunc(shell)
		if target == nil {
			shell.trackTargetId = 0
			return DefaultMove(mobj, tick)
		}
		shell.trackTargetId = target.InstId()
	} else {
		target = shell.fetchTargetFunc(shell.trackTargetId)
		if target == nil {
			shell.trackTargetId = 0
			return DefaultMove(mobj, tick)
		}
	}

	mx, my := mobj.Pos()
	tx, ty := target.Pos()
	a := base.NewVec2(mx, my)
	b := base.NewVec2(tx, ty)
	// 目標方向向量
	targetDir := b.Sub(a)
	// todo 求炮彈的方向向量
	shellRotation := shell.Rotation()
	shellDir := shellRotation.ToVec2()
	// 求叉積確定逆時針還是順時針轉
	cross := shellDir.Cross(targetDir)
	if cross == 0 {
		return DefaultMove(shell, tick)
	}
	log.Debug("@@@@@@@@ target dir %v, shell rotation %v, shell current dir %v", targetDir, shellRotation, shellDir)
	// tick時間轉向角度
	deltaMinutes := int16(time.Duration(shell.ShellStaticInfo().SteeringAngularVelocity) * tick / time.Second)
	if deltaMinutes == 0 {
		return DefaultMove(shell, tick)
	}
	deltaAngle := base.NewAngleObj(int16(deltaMinutes/60), int16(deltaMinutes%60))
	deltaAngle.Normalize()
	log.Debug(">>>>>>>> normalized deltaAngle %v", deltaAngle)
	// 利用點積求夾角
	dot := shellDir.Dot(targetDir)
	// cos(Theta) := dot / (a.Length() * b.Length())
	theta := base.ArcCosine(dot, shellDir.Length()*targetDir.Length())

	var angle base.Angle
	if cross > 0 { // 逆時針轉
		// tick時間内的轉向角度超過了需要的角度差
		if deltaAngle.GreaterEqual(theta) {
			angle = base.AngleAdd(shellRotation, theta)
			shell.RotateTo(angle)
			shell.Move(angle)
			return DefaultMove(shell, tick)
		}
		angle = base.AngleAdd(shellRotation, deltaAngle)
		log.Debug("< 逆時針")
	} else { // 順時針轉
		if deltaAngle.Negative().LessEqual(theta.Negative()) {
			angle = base.AngleSub(shellRotation, theta)
			shell.RotateTo(angle)
			shell.Move(angle)
			return DefaultMove(shell, tick)
		}
		angle = base.AngleSub(shellRotation, deltaAngle)
		log.Debug("> 順時針")
	}

	log.Debug("!!!!!!!! rotate to angle %v, previous angle %v", angle, shellRotation)
	log.Debug("track target %v to rotate angle %v", target.InstId(), angle)

	shell.Move(angle)
	shell.RotateTo(angle)
	return DefaultMove(shell, tick)
}
