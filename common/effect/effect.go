package effect

import (
	"project_b/common/object"
)

// 效果作用類型
type EffectType int

const (
	EffectTypeTime     = iota // 時間
	EffectTypeRequency = 1    // 次數
)

// 效果接口
type IEffect interface {
	InstId() uint32
	StaticInfo() *EffectStaticInfo
	SetPos(int32, int32)
	Pos() (int32, int32)
	Width() int32
	Height() int32
	Update(uint32)
	IsOver() bool
}

// 效果靜態信息
type EffectStaticInfo struct {
	Id            int32      // 配置id
	Et            EffectType // 效果類型
	Param         int32      // 參數
	Width, Height int32      // 寬高
}

type Effect struct {
	instId     uint32
	staticInfo *EffectStaticInfo
	effectFunc func(...any)
	args       []any
	center     object.Pos
	count      int32
	durationMs uint32
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

func (e Effect) Pos() (int32, int32) {
	return e.center.X, e.center.Y
}

func (e *Effect) SetPos(x, y int32) {
	e.center.X, e.center.Y = x, y
}

func (e *Effect) Update(tickMs uint32) {
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
		if int32(e.durationMs) < e.staticInfo.Param {
			if e.effectFunc != nil {
				e.effectFunc(e.args...)
			}
		} else {
			e.isOver = true
		}
		e.durationMs += tickMs
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
	e.durationMs = 0
	e.isOver = false
}
