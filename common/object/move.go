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
func DefaultMove(mobj IMovableObject, tick time.Duration, distance float64) (int32, int32) {
	var (
		fx, fy float64
	)
	x, y := mobj.Pos()
	switch mobj.Dir() {
	case DirLeft:
		fx = float64(x)
		fx -= distance
		fy = float64(y)
	case DirRight:
		fx = float64(x)
		fx += distance
		fy = float64(y)
	case DirUp:
		fy = float64(y)
		fy += distance
		fx = float64(x)
	case DirDown:
		fy = float64(y)
		fy -= distance
		fx = float64(x)
	default:
		panic("invalid direction")
	}
	return int32(fx), int32(fy)
}

func SurroundObjMove(mobj IMovableObject, tick time.Duration) (int32, int32) {
	b := mobj.(*SurroundObj)
	staticInfo := b.SurroundObjStaticInfo()
	b.accumulateTime += tick
	angle := int32(b.accumulateTime / (time.Duration(staticInfo.AngularVelocity) * time.Millisecond))
	if angle >= 1 {
		b.turnAngle += angle
		if b.turnAngle >= 360 {
			b.turnAngle -= 360
		}
		turnRadian := float64(b.turnAngle) * math.Pi / 180
		centerObj := b.getCenterObj()
		if centerObj == nil {
			return 0, 0
		}
		x, y := centerObj.Pos()
		dx := float64(staticInfo.AroundRadius) * math.Cos(turnRadian)
		dy := float64(staticInfo.AroundRadius) * math.Sin(turnRadian)
		if staticInfo.Clockwise {
			b.SetPos(int32(float64(x)+dx), int32(float64(y)-dy))
		} else {
			b.SetPos(int32(float64(x)+dx), int32(float64(y)+dy))
		}
		b.accumulateTime -= time.Duration(angle) * time.Duration(staticInfo.AngularVelocity) * time.Millisecond
	}
	return b.Pos()
}
