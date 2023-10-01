package weapon

import (
	"project_b/common/base"
)

type IHolder interface {
	LaunchPoint() base.Pos
	Forward() base.Vec2
	Camp() base.CampType
	InstId() uint32
}

type laserState int

const (
	laserStateIdle        laserState = iota // 空閑
	laserStateReadyEmit                     // 準備發射
	laserStateEmitting                      // 發射中
	laserStateEmitBlocked                   // 發射被阻擋
	laserStateRecharing                     // 充能中
)

type Laser struct {
	holder               IHolder
	staticInfo           *LaserStaticInfo
	effectFunc           func(*Laser, base.Pos, base.Pos) (base.Pos, bool)
	state                laserState
	toCancel             bool
	energy               int32
	currMs               uint32
	startPoint, endPoint base.Pos
}

func NewLaser(holder IHolder, staticInfo *LaserStaticInfo) *Laser {
	return &Laser{
		holder:     holder,
		staticInfo: staticInfo,
	}
}

func (l *Laser) StaticInfo() *LaserStaticInfo {
	return l.staticInfo
}

func (l *Laser) SetEffectFunc(effectFunc func(laser *Laser, start, end base.Pos) (base.Pos, bool)) {
	l.effectFunc = effectFunc
}

func (l *Laser) Emit() {
	if l.state == laserStateIdle {
		l.state = laserStateReadyEmit
	}
}

func (l *Laser) Cancel() {
	if l.state == laserStateReadyEmit || l.state == laserStateEmitting {
		l.toCancel = true
	}
}

func (l *Laser) Update(tickMs uint32) {
	if l.state == laserStateReadyEmit {
		if l.toCancel {
			l.state = laserStateIdle
			l.toCancel = false
		} else {
			// todo 激光效果
			l.checkEmitting(tickMs)
		}
	} else if l.state == laserStateEmitting {
		if l.toCancel {
			l.state = laserStateIdle
			l.toCancel = false
		} else {
			l.energy -= (l.staticInfo.CostPerSecond * int32(tickMs)) / 1000
			if l.energy < 0 {
				l.energy = 0
			}
			l.checkEmitting(tickMs)
			if l.energy == 0 {
				l.state = laserStateIdle
			}
		}
	} else {
		// 需要充能
		l.checkRecharging(tickMs)
	}
	l.currMs += tickMs
}

func (l *Laser) GetStartPoint() (base.Pos, bool) {
	var pos base.Pos
	if l.state != laserStateEmitting {
		return pos, false
	}
	pos = l.startPoint
	return pos, true
}

func (l *Laser) GetEndPoint() (base.Pos, bool) {
	var pos base.Pos
	if l.state != laserStateEmitting {
		return pos, false
	}
	pos = l.endPoint
	return pos, true
}

func (l *Laser) Camp() base.CampType {
	return l.holder.Camp()
}

func (l *Laser) Emitter() uint32 {
	return l.holder.InstId()
}

func (l *Laser) checkEmitting(tickMs uint32) {
	startPoint := l.holder.LaunchPoint()
	var endPoint base.Pos
	endPoint.X, endPoint.Y = base.DirPos(startPoint.X, startPoint.Y, l.staticInfo.Range, l.holder.Forward().ToAngle360())
	var o bool
	endPoint, o = l.effectFunc(l, startPoint, endPoint)
	if o {
		l.startPoint = startPoint
		l.endPoint = endPoint
		if l.state == laserStateReadyEmit {
			l.state = laserStateEmitting
		}
	} else {
		l.state = laserStateIdle
	}
}

func (l *Laser) checkRecharging(tickMs uint32) {
	if l.energy < l.staticInfo.Energy {
		if l.state == laserStateRecharing {
			l.energy += (l.staticInfo.ChargPerSecond * int32(tickMs)) / 1000
			if l.energy > l.staticInfo.Energy {
				l.energy = l.staticInfo.Energy
			}
			if l.energy == l.staticInfo.Energy {
				l.state = laserStateIdle
			}
		} else {
			l.state = laserStateRecharing
		}
	}
}
