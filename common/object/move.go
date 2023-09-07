package object

import (
	"project_b/common/base"
)

// 獲得移動距離
func GetDefaultLinearDistance(obj IMovableObject, tickMs uint32) int32 {
	return (obj.CurrentSpeed() * int32(tickMs)) / 1000
}

// 默認移動，就是直綫移動
func DefaultMove(mobj IMovableObject, tickMs uint32) (int32, int32) {
	distance := GetDefaultLinearDistance(mobj, tickMs)
	dir := mobj.MoveDir()
	sn, sd := base.Sine(dir)
	cn, cd := base.Cosine(dir)
	dx := distance * cn / cd
	dy := distance * sn / sd
	x, y := mobj.Pos()
	return x + dx, y + dy
}

// 環繞物移動
func SurroundObjMove(mobj IMovableObject, tickMs uint32) (int32, int32) {
	var sobj = mobj.(*SurroundObj)
	return getSurroundObjMovedPos(sobj, tickMs, nil)
}

// 環繞物移動信息
type SurroundMoveInfo struct {
	LastCenterX, LastCenterY int32
	TurnAngle                int32 // 單位是分，1/60度
	AccumulateMs             int32
}

// 獲得環繞物移動位置
func GetSurroundObjMovedPos(sobj *SurroundObj, tickMs uint32, moveInfo *SurroundMoveInfo) (int32, int32) {
	return getSurroundObjMovedPos(sobj, tickMs, moveInfo)
}

// 環繞物移動位置
func getSurroundObjMovedPos(sobj *SurroundObj, tickMs uint32, moveInfo *SurroundMoveInfo) (x, y int32) {
	aroundCenterObj := sobj.getAroundCenterObjFunc(sobj.aroundCenterObjInstId)
	if aroundCenterObj == nil {
		return 0, 0
	}

	var (
		turnAngle    int32
		accumulateMs int32
		cx, cy       int32
	)
	if moveInfo != nil {
		turnAngle = moveInfo.TurnAngle
		accumulateMs = int32(moveInfo.AccumulateMs)
		cx, cy = moveInfo.LastCenterX, moveInfo.LastCenterY
	} else {
		turnAngle = sobj.turnAngle
		accumulateMs = int32(sobj.accumulateMs)
		cx, cy = aroundCenterObj.Pos()
	}
	accumulateMs += int32(tickMs)

	staticInfo := sobj.SurroundObjStaticInfo()
	angle := accumulateMs * staticInfo.AngularVelocity / 1000
	turnAngle += angle
	degree, minute := turnAngle/60, turnAngle%60
	if degree >= 360 {
		degree -= 360
	}
	accumulateMs -= angle * 1000 / staticInfo.AngularVelocity
	an := base.NewAngle(int16(degree), int16(minute))
	sn, sd := base.Sine(an)
	cn, cd := base.Cosine(an)
	if staticInfo.Clockwise {
		x, y = cx+staticInfo.AroundRadius*cn/cd, cy-staticInfo.AroundRadius*sn/sd
	} else {
		x, y = cx+staticInfo.AroundRadius*cn/cd, cy+staticInfo.AroundRadius*sn/sd
	}
	if moveInfo != nil {
		moveInfo.TurnAngle = degree*60 + minute
		moveInfo.AccumulateMs = accumulateMs
	} else {
		sobj.turnAngle = degree*60 + minute
		sobj.accumulateMs = accumulateMs
	}
	return x, y
}

// 跟蹤移動
func ShellTrackMove(mobj IMovableObject, tickMs uint32) (int32, int32) {
	if mobj.Subtype() != base.ObjSubtypeShell {
		return DefaultMove(mobj, tickMs)
	}

	shell := mobj.(*Shell)
	return getShellTrackMovedPos(shell, tickMs, nil)
}

// 追蹤移動信息
type TrackMoveInfo struct {
	X, Y     int32
	Rotation base.Angle
}

// 獲得炮彈追蹤移動位置
func GetShellTrackMovedPos(shell IShell, tickMs uint32, moveInfo *TrackMoveInfo) (int32, int32) {
	s := shell.(*Shell)
	return getShellTrackMovedPos(s, tickMs, moveInfo)
}

// 炮彈追蹤移動位置
func getShellTrackMovedPos(shell *Shell, tickMs uint32, moveInfo *TrackMoveInfo) (int32, int32) {
	staticInfo := shell.ShellStaticInfo()
	if !staticInfo.TrackTarget || staticInfo.SteeringAngularVelocity <= 0 {
		return DefaultMove(shell, tickMs)
	}

	var target IObject
	if shell.trackTargetId == 0 {
		if moveInfo == nil {
			return DefaultMove(shell, tickMs)
		}
		target = shell.searchTargetFunc(shell)
		if target == nil {
			return DefaultMove(shell, tickMs)
		}
		shell.trackTargetId = target.InstId()
	} else {
		target = shell.fetchTargetFunc(shell.trackTargetId)
		if target == nil {
			return DefaultMove(shell, tickMs)
		}
	}

	mx, my := shell.Pos()
	tx, ty := target.Pos()
	a := base.NewVec2(mx, my)
	b := base.NewVec2(tx, ty)
	// 目標方向向量
	targetDir := b.Sub(a)
	// 炮彈的方向向量
	shellDir := shell.Forward()
	// 求叉積確定逆時針還是順時針轉
	cross := shellDir.Cross(targetDir)
	if cross == 0 {
		return DefaultMove(shell, tickMs)
	}

	// tick時間轉向角度
	deltaMinutes := int16(shell.ShellStaticInfo().SteeringAngularVelocity * int32(tickMs) / 1000)
	if deltaMinutes == 0 {
		return DefaultMove(shell, tickMs)
	}
	// 利用點積求夾角
	dot := shellDir.Dot(targetDir)
	// cos(Theta) := dot / (a.Length() * b.Length())
	theta := base.ArcCosine(dot, shellDir.Length()*targetDir.Length())
	thetaMinutes := theta.ToMinutes()

	var angle base.Angle
	rotation := shell.Rotation()
	if cross > 0 { // 逆時針轉
		// tick時間内的轉向角度超過了需要的角度差
		if deltaMinutes >= thetaMinutes {
			angle = base.AngleAdd(rotation, theta)
		} else {
			var deltaAngle base.Angle
			deltaAngle.Set(deltaMinutes)
			angle = base.AngleAdd(rotation, deltaAngle)
		}
	} else { // 順時針轉
		if deltaMinutes >= thetaMinutes {
			angle = base.AngleSub(rotation, theta)
		} else {
			var deltaAngle base.Angle
			deltaAngle.Set(deltaMinutes)
			angle = base.AngleSub(rotation, deltaAngle)
		}
	}

	if moveInfo == nil {
		shell.Move(angle)
		shell.RotateTo(angle)
		/*if cross > 0 {
			log.Debug("< 逆時針 !!!!!!!! rotate to angle %v, previous angle %v, track target %v, tick %v", angle, *rotation, target.InstId(), tick)
		} else {
			log.Debug("> 順時針 !!!!!!!! rotate to angle %v, previous angle %v, track target %v, tick %v", angle, *rotation, target.InstId(), tick)
		}*/
		return DefaultMove(shell, tickMs)
	} else {
		/*if cross > 0 {
			log.Debug("< 逆時針 !!!!!!!! rotate to angle %v, track target %v, tick %v", angle, target.InstId(), tick)
		} else {
			log.Debug("> 順時針 !!!!!!!! rotate to angle %v, track target %v, tick %v", angle, target.InstId(), tick)
		}*/
		x, y := shell.Pos()
		return base.MovePos(x, y, angle, shell.CurrentSpeed(), tickMs)
	}
}
