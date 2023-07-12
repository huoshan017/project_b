package object

import (
	"math"
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
	switch mobj.Dir() {
	case DirLeft:
		fx = -distance
	case DirRight:
		fx = distance
	case DirUp:
		fy = distance
	case DirDown:
		fy = -distance
	default:
		panic("invalid direction")
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
	TurnAngle                int32
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
	angle := int32(accumulateTime / (time.Duration(staticInfo.AngularVelocity) * time.Millisecond))
	if angle >= 1 {
		turnAngle += angle
		if turnAngle >= 360 {
			turnAngle -= 360
		}
		accumulateTime -= time.Duration(angle) * time.Duration(staticInfo.AngularVelocity) * time.Millisecond
		if moveInfo != nil {
			moveInfo.TurnAngle = turnAngle
			moveInfo.AccumulateTime = accumulateTime
		} else {
			sobj.turnAngle = turnAngle
			sobj.accumulateTime = accumulateTime
		}
	}
	turnRadian := float64(turnAngle) * math.Pi / 180
	s, c := math.Sincos(turnRadian)
	if staticInfo.Clockwise {
		x, y = float64(cx)+float64(staticInfo.AroundRadius)*c, float64(cy)-float64(staticInfo.AroundRadius)*s
	} else {
		x, y = float64(cx)+float64(staticInfo.AroundRadius)*c, float64(cy)+float64(staticInfo.AroundRadius)*s
	}
	return x, y
}
