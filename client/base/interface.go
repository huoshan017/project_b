package base

import (
	"project_b/common/math"
	"project_b/common/object"

	"github.com/hajimehoshi/ebiten/v2"
)

/**
 * 可播放接口
 */
type IPlayable interface {
	// 初始化
	Init()
	// 反初始化
	Uninit()
	// 重置对象
	Reset(object.IObject)
	// 播放
	Play()
	// 停止
	Stop()
	// 绘制
	Draw(*ebiten.Image, *ebiten.DrawImageOptions)
}

/**
 * 场景显示接口
 */
type IPlayableScene interface {
	// 设置视口
	// @param viewport 视口
	SetViewport(viewport *Viewport)
	// 绘制矩形范围内的可播放对象
	// @param sceneRect 表示场景的绘制逻辑坐标区域
	// @param options 从世界空间到视图到投影空间的变换矩阵
	// @param sceneImage 表示绘制的目标屏幕
	Draw(sceneRect *math.Rect, options *ebiten.DrawImageOptions, sceneImage *ebiten.Image)
}
