package main

import (
	"project_b/client/base"
	core "project_b/client_core"
	"project_b/common/object"
	"project_b/common/time"

	"github.com/hajimehoshi/ebiten/v2"
)

// 可播放接口
type IPlayable interface {
	Init()
	Uninit()
	Reset(object.IObject)
	Play()
	Stop()
	Update(*ebiten.Image)
}

// 可播放对象
type PlayableObject struct {
	IPlayable
	op   *ebiten.DrawImageOptions
	obj  object.IObject
	anim base.SpriteAnim
}

// 创建可播放对象
func NewPlayableObject(obj object.IObject, spriteConfig *base.SpriteAnimConfig) *PlayableObject {
	if spriteConfig == nil {
		panic("spriteConfig nil !!!")
	}

	op := &ebiten.DrawImageOptions{}
	x, y := obj.Pos()
	op.GeoM.Translate(float64(x), float64(y))
	return &PlayableObject{
		obj:  obj,
		op:   op,
		anim: *base.NewSpriteAnim(spriteConfig),
	}
}

// 初始化
func (po *PlayableObject) Init() {

}

// 反初始化
func (po *PlayableObject) Uninit() {

}

// 重置对象
func (po *PlayableObject) ResetObj(obj object.IObject) {
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

// 可移动物体的播放对象，有四个方向的动画
type PlayableMoveObject struct {
	IPlayable
	op         *ebiten.DrawImageOptions
	mobj       object.IMovableObject
	anims      []*base.SpriteAnim
	isMoving   bool             // 是否正在移动
	moveDir    object.Direction // 移动方向
	currSpeed  float64          // 当前速度
	updateTime time.CustomTime  // 更新时间点
	dx, dy     float64          // 目标点坐标，负数表示已经到达过该点
}

// 创建可移动物体的播放对象
func NewPlayableMoveObject(mobj object.IMovableObject, animConfig *ObjectAnimConfig) *PlayableMoveObject {
	op := &ebiten.DrawImageOptions{}
	x, y := mobj.Pos()
	op.GeoM.Translate(float64(x), float64(y))
	pobj := &PlayableMoveObject{
		op:      op,
		mobj:    mobj,
		moveDir: mobj.Dir(),
	}
	pobj.changeAnim(animConfig)
	return pobj
}

// 改变动画
func (po *PlayableMoveObject) changeAnim(animConfig *ObjectAnimConfig) {
	po.anims = []*base.SpriteAnim{
		nil,
		base.NewSpriteAnim(animConfig.AnimConfig[object.DirLeft]),
		base.NewSpriteAnim(animConfig.AnimConfig[object.DirRight]),
		base.NewSpriteAnim(animConfig.AnimConfig[object.DirUp]),
		base.NewSpriteAnim(animConfig.AnimConfig[object.DirDown]),
	}
}

// 初始化
func (po *PlayableMoveObject) Init() {
	// 注册移动停止更新事件
	po.mobj.RegisterMoveEventHandle(po.onEventMove)
	po.mobj.RegisterStopMoveEventHandle(po.onEventStopMove)
	po.mobj.RegisterUpdateEventHandle(po.onEventUpdate)
}

// 反初始化
func (po *PlayableMoveObject) Uninit() {
	// 注销移动停止更新事件
	po.mobj.UnregisterMoveEventHandle(po.onEventMove)
	po.mobj.UnregisterStopMoveEventHandle(po.onEventStopMove)
	po.mobj.UnregisterUpdateEventHandle(po.onEventUpdate)
}

// 播放
func (po *PlayableMoveObject) Play() {
	po.anims[po.moveDir].Play()
}

// 停止
func (po *PlayableMoveObject) Stop() {
	po.anims[po.moveDir].Stop()
}

// 更新
func (po *PlayableMoveObject) Update(screen *ebiten.Image) {
	if po.isMoving {
		var duration time.Duration
		now := time.Now()
		duration = now.Sub(po.updateTime)
		po.updateTime = now
		dx := po.op.GeoM.Element(0, 2)
		dy := po.op.GeoM.Element(1, 2)
		d := po.currSpeed * float64(duration) / float64(time.Second)
		switch po.moveDir {
		case object.DirLeft:
			po.op.GeoM.SetElement(0, 2, dx-d)
		case object.DirRight:
			po.op.GeoM.SetElement(0, 2, dx+d)
		case object.DirUp:
			po.op.GeoM.SetElement(1, 2, dy-d)
		case object.DirDown:
			po.op.GeoM.SetElement(1, 2, dy+d)
		default:
			return
		}
	}
	po.anims[po.moveDir].Update(screen, po.op)
}

// 移动事件处理
func (po *PlayableMoveObject) onEventMove(args ...interface{}) {
	pos := args[0].(object.Pos)
	dir := args[1].(object.Direction)
	speed := args[2].(float64)
	po.onupdate(pos, dir, speed)
	po.updateTime = core.GetSyncServTime() //args[2].(time.CustomTime)
	po.isMoving = true
	po.Play()
}

// 更新数据事件处理
func (po *PlayableMoveObject) onEventUpdate(args ...interface{}) {
	pos := args[0].(object.Pos)
	dir := args[1].(object.Direction)
	speed := args[2].(float64)
	po.onupdate(pos, dir, speed)
}

// 停止移动事件处理
func (po *PlayableMoveObject) onEventStopMove(args ...interface{}) {
	pos := args[0].(object.Pos)
	dir := args[1].(object.Direction)
	speed := args[2].(float64)
	po.onupdate(pos, dir, speed)
	po.isMoving = false
	po.Stop()
}

func (po *PlayableMoveObject) onupdate(pos object.Pos, dir object.Direction, speed float64) {
	po.dx = pos.X
	po.dy = pos.Y
	po.op.GeoM.SetElement(0, 2, po.dx)
	po.op.GeoM.SetElement(1, 2, po.dy)
	po.moveDir = dir
	po.currSpeed = speed
}

// 坦克播放对象
type PlayableTank struct {
	PlayableMoveObject
	tankObj object.ITank
}

// 创建坦克播放对象
func NewPlayableTank(tank object.ITank, animConfig *ObjectAnimConfig) *PlayableTank {
	pt := &PlayableTank{
		PlayableMoveObject: *NewPlayableMoveObject(tank, animConfig),
		tankObj:            tank,
	}
	return pt
}

// 初始化
func (pt *PlayableTank) Init() {
	pt.PlayableMoveObject.Init()
	pt.tankObj.RegisterChangeEventHandle(pt.onChange)
}

// 反初始化
func (pt *PlayableTank) Uninit() {
	pt.PlayableMoveObject.Uninit()
	pt.tankObj.UnregisterChangeEventHandle(pt.onChange)
}

// 变化事件
func (pt *PlayableTank) onChange(args ...interface{}) {
	info := args[0].(*object.ObjStaticInfo)
	level := args[1].(int32)
	pt.currSpeed = (float64)(info.Speed())
	pt.changeAnim(GetTankAnimConfig(info.Id(), level))
	pt.Play()
}
