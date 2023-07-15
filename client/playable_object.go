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
	// 初始化
	Init()
	// 反初始化
	Uninit()
	// 播放
	Play()
	// 停止
	Stop()
	// 绘制
	Draw(*ebiten.Image, *ebiten.DrawImageOptions)
	// 插值
	Interpolation() (float64, float64)
}

// 可播放对象
type PlayableObject struct {
	op   *ebiten.DrawImageOptions
	obj  object.IObject
	anim *base.SpriteAnim
}

// 创建可播放对象
func NewPlayableObject(obj object.IObject, spriteConfig *base.SpriteAnimConfig) *PlayableObject {
	if spriteConfig == nil {
		panic("spriteConfig nil !!!")
	}

	op := &ebiten.DrawImageOptions{}
	return &PlayableObject{
		obj:  obj,
		op:   op,
		anim: base.NewSpriteAnim(spriteConfig),
	}
}

// 初始化
func (po *PlayableObject) Init() {

}

// 反初始化
func (po *PlayableObject) Uninit() {

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
func (po *PlayableObject) Draw(screen *ebiten.Image, op *ebiten.DrawImageOptions) {
	// 顯示根據邏輯數據插值
	op.GeoM.Concat(po.op.GeoM)
	po.anim.Update(screen, op)
}

func (po *PlayableObject) Interpolation() (float64, float64) {
	x, y := po.obj.Pos()
	return float64(x), float64(y)
}

// 可播放的静态对象
type PlayableStaticObject struct {
	PlayableObject
}

// 创建静态物体的播放对象
func NewPlayableStaticObject(sobj object.IStaticObject, animConfig *StaticObjectAnimConfig) *PlayableStaticObject {
	playable := &PlayableStaticObject{
		PlayableObject: *NewPlayableObject(sobj.(object.IObject), animConfig.AnimConfig),
	}
	playable.Play()
	return playable
}

type IPlayableMovableObject interface {
	IPlayable
	LastInterpolation() (float64, float64)
}

// 可移动物体的播放对象，有四个方向的动画
type PlayableMoveObject struct {
	op                                     *ebiten.DrawImageOptions
	mobj                                   object.IMovableObject
	anim                                   *base.SpriteAnim
	currSpeed                              int32           // 当前速度
	lastTime                               time.CustomTime // 更新时间点
	interpolate                            bool            // 上次是停止状态
	lastX, lastY                           int32           // 上次物體位置
	lastInterpolationX, lastInterpolationY float64         // 上次插值位置
}

// 创建可移动物体的播放对象
func NewPlayableMoveObject(mobj object.IMovableObject, animConfig *MovableObjectAnimConfig) *PlayableMoveObject {
	x, y := mobj.Pos()
	pobj := &PlayableMoveObject{
		op:                 &ebiten.DrawImageOptions{},
		mobj:               mobj,
		anim:               base.NewSpriteAnim(animConfig.AnimConfig /*[object.DirRight]*/),
		interpolate:        mobj.IsMoving(),
		lastX:              x,
		lastY:              y,
		lastInterpolationX: float64(x),
		lastInterpolationY: float64(y),
	}
	return pobj
}

// 初始化
func (po *PlayableMoveObject) Init() {
	// 注册移动停止更新事件
	po.mobj.RegisterMoveEventHandle(po.onEventMove)
	po.mobj.RegisterStopMoveEventHandle(po.onEventStopMove)
}

// 反初始化
func (po *PlayableMoveObject) Uninit() {
	// 注销移动停止更新事件
	po.mobj.UnregisterMoveEventHandle(po.onEventMove)
	po.mobj.UnregisterStopMoveEventHandle(po.onEventStopMove)
}

// 播放
func (po *PlayableMoveObject) Play() {
	po.anim.Play()
}

// 停止
func (po *PlayableMoveObject) Stop() {
	po.anim.Stop()
}

// 更新
func (po *PlayableMoveObject) Draw(screen *ebiten.Image, op *ebiten.DrawImageOptions) {
	op.GeoM.Concat(po.op.GeoM)
	po.anim.Update(screen, op)
}

// 上次插值位置
func (po *PlayableMoveObject) LastInterpolation() (float64, float64) {
	return po.lastInterpolationX, po.lastInterpolationY
}

// 插值，在Draw前同步调用，得到位置插值
// ----|------------|------------|------------|------------|-------------|--------------|--------------|--------------|--------------|----
//
//	Draw      Update(Draw)    Draw         Draw         Draw       Update(Draw)      Draw           Draw           Draw        Update(Draw)
func (po *PlayableMoveObject) Interpolation() (x, y float64) {
	cx, cy := po.mobj.Pos()
	if !po.interpolate {
		return float64(cx), float64(cy)
	}
	if !po.mobj.IsMoving() {
		return float64(cx), float64(cy)
	}
	// 上一次Update的坐标点与当前的不一样，说明又Update了，重置LastX和LastY和开始时间
	// 所以每次Update后都要重置LastX,lastY,LastTime，是因为要从重置点开始计算位置插值
	if po.lastX != cx || po.lastY != cy {
		po.lastX = cx
		po.lastY = cy
		po.lastTime = core.GetSyncServTime()
	}
	duration := time.Since(po.lastTime)
	nx, ny := object.DefaultMove(po.mobj, duration)
	po.lastInterpolationX, po.lastInterpolationY = float64(nx), float64(ny)
	return po.lastInterpolationX, po.lastInterpolationY
}

// 移动事件处理
func (po *PlayableMoveObject) onEventMove(args ...any) {
	po.Play()
	po.lastX, po.lastY = po.mobj.Pos()
	po.lastTime = core.GetSyncServTime()
	po.interpolate = true
}

// 停止移动事件处理
func (po *PlayableMoveObject) onEventStopMove(args ...any) {
	po.Stop()
	po.interpolate = false
}

// 坦克播放对象
type PlayableTank struct {
	*PlayableMoveObject
}

// 创建坦克播放对象
func NewPlayableTank(tank object.ITank, animConfig *MovableObjectAnimConfig) *PlayableTank {
	return &PlayableTank{
		PlayableMoveObject: NewPlayableMoveObject(tank, animConfig),
	}
}

// 初始化
func (pt *PlayableTank) Init() {
	pt.PlayableMoveObject.Init()
	pt.mobj.(object.ITank).RegisterChangeEventHandle(pt.onChange)
}

// 反初始化
func (pt *PlayableTank) Uninit() {
	pt.PlayableMoveObject.Uninit()
	pt.mobj.(object.ITank).UnregisterChangeEventHandle(pt.onChange)
}

// 变化事件
func (pt *PlayableTank) onChange(args ...any) {
	info := args[0].(*object.ObjStaticInfo)
	pt.currSpeed = (info.Speed())
	pt.Play()
}

// 環繞物體播放對象
type PlayableSurroundObj struct {
	*PlayableMoveObject
	sobj                    object.ISurroundObject
	lastMoveInfo            object.SurroundMoveInfo
	playableAroundCenterObj IPlayable
}

// 創建環繞物體播放對象
func NewPlayableSurroundObj(sobj object.ISurroundObject, animConfig *MovableObjectAnimConfig, playableAroundCenterObj IPlayable) *PlayableSurroundObj {
	pobj := &PlayableSurroundObj{
		PlayableMoveObject:      NewPlayableMoveObject(sobj, animConfig),
		sobj:                    sobj,
		playableAroundCenterObj: playableAroundCenterObj,
	}
	pobj.interpolate = true
	pobj.lastTime = core.GetSyncSendTime()
	cx, cy := playableAroundCenterObj.Interpolation()
	pobj.lastMoveInfo.LastCenterX, pobj.lastMoveInfo.LastCenterY = int32(cx), int32(cy)
	return pobj
}

// 初始化
func (ps *PlayableSurroundObj) Init() {
	ps.sobj.RegisterLateUpdateEventHandle(ps.onEventLateUpdate)
}

// 反初始化
func (ps *PlayableSurroundObj) Uninit() {
	ps.sobj.UnregisterLateUpdateEventHandle(ps.onEventLateUpdate)
}

// 插值
func (ps *PlayableSurroundObj) Interpolation() (x, y float64) {
	nx, ny := ps.sobj.Pos()
	if !ps.interpolate {
		return float64(nx), float64(ny)
	}
	if !ps.mobj.IsMoving() {
		return float64(nx), float64(ny)
	}
	duration := time.Since(ps.lastTime)
	ps.lastTime = time.Now()
	var interpolateX, interpolateY float64
	if pobj, o := ps.playableAroundCenterObj.(IPlayableMovableObject); !o {
		interpolateX, interpolateY = ps.playableAroundCenterObj.Interpolation()
	} else {
		interpolateX, interpolateY = pobj.LastInterpolation()
	}
	ps.lastMoveInfo.LastCenterX, ps.lastMoveInfo.LastCenterY = int32(interpolateX), int32(interpolateY)
	nx, ny = object.GetSurroundObjMovedPos(ps.sobj.(*object.SurroundObj), duration, &ps.lastMoveInfo)
	return float64(nx), float64(ny)
}

// 更新事件處理
func (ps *PlayableSurroundObj) onEventLateUpdate(args ...any) {
	ps.lastTime = core.GetSyncServTime()
	ps.lastMoveInfo.TurnAngle = args[0].(int32)
	ps.lastMoveInfo.AccumulateTime = args[1].(time.Duration)
}

// 可播放效果
type PlayableEffect struct {
	op     ebiten.DrawImageOptions
	effect object.IEffect
	anim   *base.SpriteAnim
}

// 創建可播放效果
func NewPlayableEffect(effect object.IEffect, animConfig *base.SpriteAnimConfig) *PlayableEffect {
	return &PlayableEffect{
		effect: effect,
		anim:   base.NewSpriteAnim(animConfig),
	}
}

// 初始化
func (po *PlayableEffect) Init() {

}

// 反初始化
func (po *PlayableEffect) Uninit() {

}

// 播放
func (po *PlayableEffect) Play() {
	po.anim.Play()
}

// 停止
func (po *PlayableEffect) Stop() {
	po.anim.Stop()
}

// 更新
func (po *PlayableEffect) Draw(screen *ebiten.Image, op *ebiten.DrawImageOptions) {
	// 顯示根據邏輯數據插值
	op.GeoM.Concat(po.op.GeoM)
	po.anim.Update(screen, op)
}

// 插值
func (po *PlayableEffect) Interpolation() (float64, float64) {
	x, y := po.effect.Pos()
	return float64(x), float64(y)
}
