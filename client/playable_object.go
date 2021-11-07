package main

import (
	client_base "project_b/client/base"
	common_object "project_b/common/object"

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
	obj  common_object.IObject
	anim client_base.SpriteAnim
}

// 创建可播放对象
func NewPlayableObject(obj common_object.IObject, spriteConfig *client_base.SpriteAnimConfig) *PlayableObject {
	if spriteConfig == nil {
		panic("spriteConfig nil !!!")
	}
	return &PlayableObject{
		obj:  obj,
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
func (po *PlayableObject) Update(screen *ebiten.Image) {
	x, y := po.obj.Pos()
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	po.anim.Update(screen, op)
}
