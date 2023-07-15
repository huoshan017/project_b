package object

import (
	"project_b/common/base"
	"project_b/common/time"
)

// 獲得移動距離
func GetDefaultLinearDistance(obj IMovableObject, duration time.Duration) float64 {
	return float64(int64(obj.CurrentSpeed())*int64(duration)) / float64(time.Second)
}

// 默認移動，就是直綫移動
func DefaultMove(mobj IMovableObject, tick time.Duration) (int32, int32) {
	var (
		fx, fy float64
	)
	distance := GetDefaultLinearDistance(mobj, tick)
	orientation := mobj.Orientation()
	if orientation == 0 {
		fx = distance
	} else if orientation == 90 {
		fy = distance
	} else if orientation == 180 {
		fx = -distance
	} else if orientation == 270 {
		fy = -distance
	}
	x, y := mobj.Pos()
	return x + int32(fx), y + int32(fy)
}

func SurroundObjMove(mobj IMovableObject, tick time.Duration) (int32, int32) {
	var (
		x, y float64
		sobj = mobj.(*SurroundObj)
	)
	x, y = getSurroundObjMovedPos(sobj, tick, nil)
	return int32(x), int32(y)
}

type SurroundMoveInfo struct {
	LastCenterX, LastCenterY int32
	TurnAngle                int32 // 單位是分，1/60度
	AccumulateTime           time.Duration
}

func GetSurroundObjMovedPos(sobj *SurroundObj, tick time.Duration, moveInfo *SurroundMoveInfo) (int32, int32) {
	x, y := getSurroundObjMovedPos(sobj, tick, moveInfo)
	return int32(x), int32(y)
}

func getSurroundObjMovedPos(sobj *SurroundObj, tick time.Duration, moveInfo *SurroundMoveInfo) (x, y float64) {
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
	s := base.Sine(an)
	c := base.Cosine(an)
	if staticInfo.Clockwise {
		x, y = float64(cx+staticInfo.AroundRadius*int32(c.Numerator)/int32(c.Denominator)), float64(cy-staticInfo.AroundRadius*int32(s.Numerator)/int32(s.Denominator))
	} else {
		x, y = float64(cx+staticInfo.AroundRadius*int32(c.Numerator)/int32(c.Denominator)), float64(cy+staticInfo.AroundRadius*int32(s.Numerator)/int32(s.Denominator))
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
