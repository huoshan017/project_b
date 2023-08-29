package client_base

import (
	"image"
	"project_b/common/time"

	"github.com/hajimehoshi/ebiten/v2"
)

// 坐标索引
type SpriteIndex struct {
	X int
	Y int
}

type SpriteAnimConfig struct {
	Image        *ebiten.Image
	PlayInterval int32 // 毫秒
	FramePosList []SpriteIndex
	FrameWidth   int
	FrameHeight  int
}

type SpriteAnimState int

const (
	stop    = SpriteAnimState(0)
	playing = SpriteAnimState(1)
)

// 精灵动画
type SpriteAnim struct {
	Config       *SpriteAnimConfig
	currFrame    int32
	lastUpdateMs uint32
	currState    SpriteAnimState
	images       []*ebiten.Image
}

// 创建精灵动画
func NewSpriteAnim(config *SpriteAnimConfig) *SpriteAnim {
	return &SpriteAnim{
		Config: config,
		images: make([]*ebiten.Image, len(config.FramePosList)),
	}
}

// 寬度
func (a *SpriteAnim) Width() int32 {
	return int32(a.Config.FrameWidth)
}

// 高度
func (a *SpriteAnim) Height() int32 {
	return int32(a.Config.FrameHeight)
}

// 播放
func (a *SpriteAnim) Play() {
	if a.currState != stop {
		return
	}

	a.currState = playing
	a.currFrame = 0
}

// 停止
func (a *SpriteAnim) Stop() {
	a.currState = stop
	a.currFrame = 0
}

// 更新
func (a *SpriteAnim) Update(screen *ebiten.Image, options *ebiten.DrawImageOptions) {
	// 停止模式，播放当前帧，一般是第0帧
	if a.currState == stop || a.Config.PlayInterval == 0 {
		a.draw(a.currFrame, screen, options)
		return
	}

	currMs := time.CurrentMs()
	if a.lastUpdateMs == 0 {
		a.draw(a.currFrame, screen, options)
		a.lastUpdateMs = currMs
		return
	}

	sub := int32(currMs - a.lastUpdateMs)
	if sub < a.Config.PlayInterval {
		a.draw(a.currFrame, screen, options)
		return
	}

	for int32(sub) >= a.Config.PlayInterval {
		a.currFrame += 1
		if a.currFrame >= int32(len(a.Config.FramePosList)) {
			a.currFrame = 0
		}
		sub -= a.Config.PlayInterval
	}
	a.draw(a.currFrame, screen, options)

	// 更新时间，保证该时间的增长是刷新间隔PlayerInterval的整数倍
	a.lastUpdateMs = currMs - uint32(sub)
}

// 内部函数，画到screen上
func (a *SpriteAnim) draw(frameIndex int32, screen *ebiten.Image, options *ebiten.DrawImageOptions) {
	c := a.Config
	f := c.FramePosList[frameIndex]
	x0 := f.X * c.FrameWidth
	y0 := f.Y * c.FrameHeight
	x1 := x0 + c.FrameWidth
	y1 := y0 + c.FrameHeight
	if a.images[frameIndex] == nil {
		a.images[frameIndex] = c.Image.SubImage(image.Rect(x0, y0, x1, y1)).(*ebiten.Image)
	}
	screen.DrawImage(a.images[frameIndex], options)
}
