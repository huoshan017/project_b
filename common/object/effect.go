package object

import "project_b/common/time"

type Effect struct {
	instId     uint32
	staticInfo *EffectStaticInfo
	effectFunc func(...any)
	args       []any
	center     Pos
	count      int32
	startTime  time.CustomTime
	isOver     bool
}

func NewEffect(instId uint32, staticInfo *EffectStaticInfo, effectFunc func(...any), args ...any) *Effect {
	effect := &Effect{}
	effect.init(instId, staticInfo, effectFunc, args...)
	return effect
}

func (e Effect) InstId() uint32 {
	return e.instId
}

func (e Effect) Width() int32 {
	return e.staticInfo.Width
}

func (e Effect) Height() int32 {
	return e.staticInfo.Height
}

func (e Effect) Center() (int32, int32) {
	return e.center.X, e.center.Y
}

func (e *Effect) SetCenter(x, y int32) {
	e.center.X, e.center.Y = x, y
}

func (e *Effect) Update() {
	if e.isOver {
		return
	}
	if e.staticInfo.Et == EffectTypeRequency {
		if e.count < e.staticInfo.Param {
			if e.effectFunc != nil {
				e.effectFunc(e.args...)
			}
			e.count += 1
		} else {
			e.isOver = true
		}
	} else if e.staticInfo.Et == EffectTypeTime {
		if e.startTime.IsZero() || time.Since(e.startTime) < time.Duration(e.staticInfo.Param) {
			if e.effectFunc != nil {
				e.effectFunc(e.args...)
			}
			e.startTime = time.Now()
		} else {
			e.isOver = true
		}
	}
}

func (e *Effect) IsOver() bool {
	return e.isOver
}

func (e Effect) StaticInfo() *EffectStaticInfo {
	return e.staticInfo
}

func (e *Effect) init(instId uint32, staticInfo *EffectStaticInfo, effectFunc func(...any), args ...any) {
	e.instId = instId
	e.staticInfo = staticInfo
	e.effectFunc = effectFunc
	e.args = append(e.args, args...)
}

func (e *Effect) uninit() {
	e.instId = 0
	e.staticInfo = nil
	e.effectFunc = nil
	e.args = nil
	e.count = 0
	e.startTime = time.CustomTime{}
	e.isOver = false
}
