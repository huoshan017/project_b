package base

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

// 矩形区域，(x, y)表示左下角，w宽度，h高度
type Rect struct {
	x, y, w, h int32
}

func NewRect(x, y, w, h int32) *Rect {
	return &Rect{
		x: x, y: y, w: w, h: h,
	}
}

func (r Rect) X() int32 {
	return r.x
}

func (r Rect) Y() int32 {
	return r.y
}

func (r Rect) W() int32 {
	return r.w
}

func (r Rect) H() int32 {
	return r.h
}

func (r *Rect) SetX(x int32) {
	r.x = x
}

func (r *Rect) SetY(y int32) {
	r.y = y
}

func (r *Rect) SetW(w int32) {
	r.w = w
}

func (r *Rect) SetH(h int32) {
	r.h = h
}

func (r *Rect) Move(x, y int32) {
	r.x += x
	r.y += y
}

func (r *Rect) MoveTo(x, y int32) {
	r.x = x
	r.y = y
}

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
	distance2Viewport     float64                  // distance to viewport
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
func CreateCamera(viewport *Viewport, fov int32, distance2Viewport float64) *Camera {
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
	c.distance2Viewport = distance2Viewport //float64(c.viewport.w) / 2 * c.tanHalfFov
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
	if float64(height) < c.distance2Viewport {
		height = int32(c.distance2Viewport)
	}
	c.height = height
}

/**
 * 改变高度
 * delta 高度变化
 */
func (c *Camera) ChangeHeight(delta int32) {
	c.height += delta
	if float64(c.height) < c.distance2Viewport {
		c.height = int32(c.distance2Viewport)
	}
}

/**
 * 屏幕坐標到世界坐標
 * x, y 屏幕坐標
 */
func (c *Camera) Screen2World(x, y int32) (int32, int32) {
	var halfViewportWidthLength = c.tanHalfFov * c.distance2Viewport
	//var halfWorldSectionWidth = halfViewportWidthLength * float64(c.height) / c.distance2Viewport
	var halfViewportHeightLength = halfViewportWidthLength * float64(c.viewport.h) / float64(c.viewport.w)
	//var halfWorldSectionHeight = halfViewportHeightLength * float64(c.height) / c.distance2Viewport
	// (vx - c.wx) / halfViewportWidthLength = (x - c.viewport.w/2) / (c.viewport.w/2)
	// (vx - c.wx) / (wx - c.wx) = c.distance2Viewport / c.height
	// todo 世界坐標系和屏幕坐標系Y軸的正方向是相反的
	// (c.wy - vy) / halfViewportHeightLength = (y - c.viewport.h/2) / (c.viewport.h/2)
	// (c.wy - vy) / (c.wy - wy) = c.distance2Viewport / c.height
	deltax := halfViewportWidthLength * float64(x-c.viewport.w/2) / float64(c.viewport.w/2)
	wx := float64(c.wx) + deltax*float64(c.height)/c.distance2Viewport
	deltay := halfViewportHeightLength * float64(y-c.viewport.h/2) / float64(c.viewport.h/2)
	wy := float64(c.wy) - deltay*float64(c.height)/c.distance2Viewport
	return int32(wx), int32(wy)
}

/**
 * 世界坐標到屏幕坐標
 * wx, wy 世界坐標
 */
func (c *Camera) World2Screen(wx, wy int32) (int32, int32) {
	var halfViewportWidthLength = c.tanHalfFov * c.distance2Viewport
	var halfViewportHeightLength = halfViewportWidthLength * float64(c.viewport.h) / float64(c.viewport.w)
	//var halfWorldSectionWidth = halfViewportWidthLength * float64(c.height) / c.distance2Viewport
	//var halfWorldSectionHeight = halfViewportHeightLength * float64(c.height) / c.distance2Viewport
	deltax := float64(wx-c.wx) * c.distance2Viewport / float64(c.height)
	x := deltax*float64(c.viewport.w/2)/halfViewportWidthLength + float64(c.viewport.w/2)
	deltay := float64(c.wy-wy) * c.distance2Viewport / float64(c.height)
	y := deltay*float64(c.viewport.h/2)/halfViewportHeightLength + float64(c.viewport.h/2)
	return int32(x), int32(y)
}

/**
 * 屏幕區域到世界區域
 */
func (c *Camera) Screen2WorldRect() Rect {
	lx, ly := c.Screen2World(0, 0)
	rx, ry := c.Screen2World(c.viewport.w, c.viewport.h)
	return Rect{lx, ly, rx - lx, ry - ly}
}

/**
 * 世界區域到屏幕區域
 */
func (c *Camera) World2ScreenRect(wr Rect) Rect {
	lwx, lwy := c.World2Screen(wr.x, wr.y)
	rwx, rwy := c.World2Screen(wr.x+wr.w, wr.y+wr.h)
	return Rect{x: lwx, y: lwy, w: rwx - lwx, h: rwy - lwy}
}

/**
 * 绘制到屏幕
 * @param screen  目標屏幕
 */
/*func (c *Camera) Draw(screen *ebiten.Image) {
	var x, y, w, h float64
	// 计算绘制范围
	// (c.wx-x)/c.height = tan(c.fov*math.Pi/180/2)
	// (c.wy-y)/(c.wx-x) = c.viewport.h/c.viewport.w
	// distance2Viewport/c.height = c.viewport.w/w
	// w/h = c.viewport.w/c.viewport.h
	// 左下角坐标
	x = float64(c.wx) - c.tanHalfFov*float64(c.height)
	y = float64(c.wy) - float64(c.viewport.h)*(float64(c.wx)-x)/float64(c.viewport.w)
	// 宽和高
	w = float64(c.viewport.w) * float64(c.height) / c.distance2Viewport
	h = w * float64(c.viewport.h) / float64(c.viewport.w)
	// 更新变换矩阵
	c.updateViewProjectionMatrix()
	// 把绘制范围内的场景经过变换后绘制到场景视图
	c.scene.Draw(NewRect(int32(x), int32(y), int32(w), int32(h)), c.viewProjectionOptions, c.sceneImage)
	screen.Fill(color.Black)
	// 绘制到屏幕
	screen.DrawImage(c.sceneImage, c.viewProjectionOptions)
}*/

/**
 * view和projection变换矩陣
 */
/*func (c *Camera) updateViewProjectionMatrix() {
	// 1. (x',y')表示世界空间坐标(x,y)压缩到视口大小区域上的坐标点，如下
	//    (x'-c.wx)/(x-c.wx) = (y'-c.wy)/(y-c.wy) = c.distance2Viewport/c.height
	//    x' = c.wx + (x-c.wx)*c.distance2Viewport/c.height
	//    y' = c.wy + (y-c.wy)*c.distance2Viewport/c.height
	// 2. 或者先变换到相机空间，再压缩到视口区域，最后投影平移到视口空间，如下
	//    x' = x-c.wx    y' = y-c.wy
	//    x' = (x-c.wx)*c.distance2Viewport/c.height   y' = (y-c.wy)*c.distance2Viewport/c.height

	c.viewProjectionOptions.GeoM.Reset()
	// 世界空间到相机空间，相机坐标为原点
	c.viewProjectionOptions.GeoM.Translate(float64(-c.wx), float64(-c.wy))
	// 相机空间缩放到视口区域，相机坐标原点不变
	c.viewProjectionOptions.GeoM.Scale(c.distance2Viewport/float64(c.height), c.distance2Viewport/float64(c.height))
	// 再平移到视口空间 (x軸向上y軸向右)
	c.viewProjectionOptions.GeoM.Translate(float64(c.viewport.x/2), float64(c.viewport.y/2))
}*/
