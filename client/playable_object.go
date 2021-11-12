package main

import (
	client_base "project_b/client/base"
	common_object "project_b/common/object"
	"project_b/utils"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

// 可播放接口
type IPlayable interface {
	Play()
	Stop()
	Update(*ebiten.Image)
}

// 可播放对象
type PlayableObject struct {
	IPlayable
	op   *ebiten.DrawImageOptions
	obj  common_object.IObject
	anim client_base.SpriteAnim
}

// 创建可播放对象
func NewPlayableObject(obj common_object.IObject, spriteConfig *client_base.SpriteAnimConfig) *PlayableObject {
	if spriteConfig == nil {
		panic("spriteConfig nil !!!")
	}

	op := &ebiten.DrawImageOptions{}
	x, y := obj.Pos()
	op.GeoM.Translate(float64(x), float64(y))
	return &PlayableObject{
		obj:  obj,
		op:   op,
		anim: *client_base.NewSpriteAnim(spriteConfig),
	}
}

// 重置对象
func (po *PlayableObject) ResetObj(obj common_object.IObject) {
	po.obj = obj
}

// 播放
func (po *PlayableObject) Play() {
	po.anim.Play()
}

// 停止
func (po *PlayableObject) Stop() {
	po.anim.Stop()
}

// 更新
func (po *PlayableObject) Update(tick time.Duration, screen *ebiten.Image) {
	x, y := po.obj.Pos()
	x0 := po.op.GeoM.Element(0, 2)
	y0 := po.op.GeoM.Element(1, 2)
	po.op.GeoM.Translate(float64(x)-x0, float64(y)-y0)
	po.anim.Update(screen, po.op)
}

// 可移动物体的播放对象
type PlayableMoveObject struct {
	PlayableObject
	mobj     common_object.IMovableObject
	x, y     float64
	lastTick time.Duration // 上一帧tick
}

// 创建可移动物体的播放对象
func NewPlayableMoveObject(mobj common_object.IMovableObject, spriteConfig *client_base.SpriteAnimConfig) *PlayableMoveObject {
	obj := NewPlayableObject(mobj, spriteConfig)
	x, y := mobj.Pos()
	return &PlayableMoveObject{
		PlayableObject: *obj,
		mobj:           mobj,
		x:              x,
		y:              y,
	}
}

// 更新
// todo 如果一个方向上的移动到停止最后逻辑帧和渲染帧不一致的话，则再次移动或换方向时会出现被往回拉情况
func (po *PlayableMoveObject) Update(tick time.Duration, screen *ebiten.Image) {
	x, y := po.mobj.Pos()
	dx := po.op.GeoM.Element(0, 2)
	dy := po.op.GeoM.Element(1, 2)
	if utils.Float64IsEqual(x, po.x) && utils.Float64IsEqual(y, po.y) {
		if po.mobj.IsMove() {
			po.lastTick += tick
			ms := float64(po.lastTick / time.Millisecond)
			d := float64(po.mobj.CurrentSpeed()) * ms / 1000
			getLog().Debug("PlayableMoveObject currentSpeed = %v, tick = %v, lastTick = %v, d = %v", po.mobj.CurrentSpeed(), tick, po.lastTick, d)
			dir := po.mobj.Dir()
			switch dir {
			case common_object.DirLeft:
				po.op.GeoM.SetElement(0, 2, dx-d)
			case common_object.DirRight:
				po.op.GeoM.SetElement(0, 2, dx+d)
			case common_object.DirUp:
				po.op.GeoM.SetElement(1, 2, dy-d)
			case common_object.DirDown:
				po.op.GeoM.SetElement(1, 2, dy+d)
			default:
				return
			}
			po.lastTick -= time.Duration(ms) * time.Millisecond
			dx = po.op.GeoM.Element(0, 2)
			dy = po.op.GeoM.Element(1, 2)
			getLog().Debug("PlayableMoveObject after lastTick = %v, dx=%v, dy=%v", po.lastTick, dx, dy)
		}
	} else {
		getLog().Debug("PlayableMoveObject dx=%v, dy=%v, po.x=%v, po.y=%v, x=%v, y=%v", dx, dy, po.x, po.y, x, y)
		po.op.GeoM.SetElement(0, 2, x)
		po.op.GeoM.SetElement(1, 2, y)
		po.x = x
		po.y = y
		po.lastTick = 0
	}
	po.anim.Update(screen, po.op)
}
