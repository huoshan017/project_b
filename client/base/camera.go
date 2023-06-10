package base

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

/**
 * 世界坐标系x轴向上，y轴向右
 * 相机坐标系与世界坐标系的轴向一致
 */
type Camera struct {
	sceneImage            *ebiten.Image
	viewport              *Viewport
	scene                 IPlayableScene           // interface of Scene Playable
	wx, wy                int32                    // world cordinate
	height                int32                    // distance to world
	fov                   int32                    // field of view
	tanHalfFov            float64                  // tan value for fov
	nearPlane             float64                  // near plane
	viewportOptions       *ebiten.DrawImageOptions // 视口变换参数
	viewProjectionOptions *ebiten.DrawImageOptions // 视图投影变换参数
}

/**
 * 创建照相机
 * viewport 视口，坐标对应到屏幕坐标，区域对应屏幕的一部分或全部
 * wx wy 世界坐标
 * height 与地图的高度
 * fov 视场角，单位是角度
 */
func CreateCamera(viewport *Viewport, fov int32, nearPlane float64) *Camera {
	op := ebiten.GeoM{}
	op.Translate(float64(viewport.x), float64(viewport.y))
	c := &Camera{
		sceneImage:            ebiten.NewImage(int(viewport.w), int(viewport.h)),
		viewport:              viewport,
		fov:                   fov,
		viewportOptions:       &ebiten.DrawImageOptions{GeoM: op},
		viewProjectionOptions: &ebiten.DrawImageOptions{GeoM: ebiten.GeoM{}},
	}
	c.tanHalfFov = math.Tan(float64(c.fov) * math.Pi / 360)
	c.nearPlane = nearPlane //float64(c.viewport.w) / 2 * c.tanHalfFov
	return c
}

/**
 * 设置可视场景
 */
func (c *Camera) SetScene(scene IPlayableScene) {
	scene.SetViewport(c.viewport)
	c.scene = scene
}

/**
 * 在世界坐标系中移动
 * x, y 移动的量
 */
func (c *Camera) Move(x, y int32) {
	c.wx += x
	c.wy += y
}

/**
 * 相机移动到一个世界坐标
 * wx, wy 世界坐标
 */
func (c *Camera) MoveTo(wx, wy int32) {
	c.wx = wx
	c.wy = wy
}

/**
 * 设置高度
 * height 高度
 */
func (c *Camera) SetHeight(height int32) {
	if float64(height) < c.nearPlane {
		height = int32(c.nearPlane)
	}
	c.height = height
}

/**
 * 改变高度
 * delta 高度变化
 */
func (c *Camera) ChangeHeight(delta int32) {
	c.height += delta
	if float64(c.height) < c.nearPlane {
		c.height = int32(c.nearPlane)
	}
}

/**
 * 屏幕坐標到世界坐標
 * x, y 屏幕坐標
 */
func (c *Camera) Screen2World(x, y int32) (int32, int32) {
	var halfViewportWidthLength = c.tanHalfFov * c.nearPlane
	var halfViewportHeightLength = halfViewportWidthLength * float64(c.viewport.h) / float64(c.viewport.w)
	// (vx - c.wx) / halfViewportWidthLength = (x - c.viewport.w/2) / (c.viewport.w/2)
	// (vx - c.wx) / (wx - c.wx) = c.nearPlane / c.height
	// todo 世界坐標系和屏幕坐標系Y軸的正方向是相反的
	// (c.wy - vy) / halfViewportHeightLength = (y - c.viewport.h/2) / (c.viewport.h/2)
	// (c.wy - vy) / (c.wy - wy) = c.nearPlane / c.height
	deltax := halfViewportWidthLength * float64(2*x-c.viewport.w) / float64(c.viewport.w)
	wx := float64(c.wx) + deltax*float64(c.height)/c.nearPlane
	deltay := halfViewportHeightLength * float64(2*y-c.viewport.h) / float64(c.viewport.h)
	wy := float64(c.wy) - deltay*float64(c.height)/c.nearPlane
	return int32(wx), int32(wy)
}

/**
 * 世界坐標到屏幕坐標
 * wx, wy 世界坐標
 */
func (c *Camera) World2Screen(wx, wy int32) (int32, int32) {
	var halfViewportWidthLength = c.tanHalfFov * c.nearPlane
	var halfViewportHeightLength = halfViewportWidthLength * float64(c.viewport.h) / float64(c.viewport.w)
	deltax := float64(wx-c.wx) * c.nearPlane / float64(c.height)
	x := deltax*float64(c.viewport.w/2)/halfViewportWidthLength + float64(c.viewport.w/2)
	deltay := float64(c.wy-wy) * c.nearPlane / float64(c.height)
	y := deltay*float64(c.viewport.h/2)/halfViewportHeightLength + float64(c.viewport.h/2)
	return int32(x), int32(y)
}
