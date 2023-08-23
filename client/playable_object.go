package main

import (
	"project_b/client_base"
	common_base "project_b/common/base"
	"project_b/common/effect"
	"project_b/common/object"
	"project_b/common/time"
	"project_b/core"

	"github.com/hajimehoshi/ebiten/v2"
)

type Transform struct {
	tx, ty   float64
	rotation common_base.Angle
	sx, sy   float64
}

func NewTransform() Transform {
	return Transform{
		sx: 1,
		sy: 1,
	}
}

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
	Interpolation(*Transform)
}

// 可播放对象
type PlayableObject struct {
	op   *ebiten.DrawImageOptions
	obj  object.IObject
	anim *client_base.SpriteAnim
}

// 创建可播放对象
func NewPlayableObject(obj object.IObject, spriteConfig *client_base.SpriteAnimConfig) *PlayableObject {
	if spriteConfig == nil {
		panic("spriteConfig nil !!!")
	}

	op := &ebiten.DrawImageOptions{}
	return &PlayableObject{
		obj:  obj,
		op:   op,
		anim: client_base.NewSpriteAnim(spriteConfig),
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

// 插值
func (po *PlayableObject) Interpolation(transform *Transform) {
	x, y := po.obj.Pos()
	transform.tx, transform.ty = float64(x), float64(y)
	transform.sx, transform.sy = 1, 1
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

// 可播放的物品對象
type PlayableItemObject struct {
	PlayableObject
}

// 創建可播放的物品對象
func NewPlayableItemObject(iobj object.IItemObject, animConfig *client_base.SpriteAnimConfig) *PlayableItemObject {
	playable := &PlayableItemObject{
		PlayableObject: *NewPlayableObject(iobj.(object.IObject), animConfig),
	}
	return playable
}

// 可移動物體播放對象接口
type IPlayableMovableObject interface {
	IPlayable
	LastInterpolation() (float64, float64)
}

// 可移动物体的播放对象，有四个方向的动画
type PlayableMoveObject struct {
	op                *ebiten.DrawImageOptions
	mobj              object.IMovableObject
	anim              *client_base.SpriteAnim
	currSpeed         int32           // 当前速度
	lastTime          time.CustomTime // 更新时间点
	interpolate       bool            // 上次是停止状态
	lastX, lastY      int32           // 上次物體位置
	lastInterpolation Transform       // 上次插值位置
}

// 创建可移动物体的播放对象
func NewPlayableMoveObject(mobj object.IMovableObject, animConfig *MovableObjectAnimConfig) *PlayableMoveObject {
	x, y := mobj.Pos()
	pobj := &PlayableMoveObject{
		op:          &ebiten.DrawImageOptions{},
		mobj:        mobj,
		anim:        client_base.NewSpriteAnim(animConfig.AnimConfig /*[object.DirRight]*/),
		interpolate: mobj.IsMoving(),
		lastX:       x,
		lastY:       y,
	}
	pobj.lastInterpolation.tx = float64(x)
	pobj.lastInterpolation.ty = float64(y)
	pobj.lastInterpolation.sx, pobj.lastInterpolation.sy = 1, 1
	pobj.lastInterpolation.rotation = mobj.WorldRotation()
	return pobj
}

// 初始化
func (po *PlayableMoveObject) Init() {
	// 注册移动停止更新事件
	po.mobj.RegisterMoveEventHandle(po.onEventMove)
	po.mobj.RegisterStopMoveEventHandle(po.onEventStopMove)
	po.mobj.RegisterPauseEventHandle(po.onEventPause)
	po.mobj.RegisterResumeEventHandle(po.onEventResume)
}

// 反初始化
func (po *PlayableMoveObject) Uninit() {
	// 注销移动停止更新事件
	po.mobj.UnregisterMoveEventHandle(po.onEventMove)
	po.mobj.UnregisterStopMoveEventHandle(po.onEventStopMove)
	po.mobj.UnregisterPauseEventHandle(po.onEventPause)
	po.mobj.UnregisterResumeEventHandle(po.onEventResume)
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
func (po *PlayableMoveObject) LastInterpolation() *Transform {
	return &po.lastInterpolation
}

// 插值，在Draw前同步调用，得到位置插值
// ----|------------|------------|------------|------------|-------------|--------------|--------------|--------------|--------------|----
//
//	Draw      Update(Draw)    Draw         Draw         Draw       Update(Draw)      Draw           Draw           Draw        Update(Draw)
func (po *PlayableMoveObject) Interpolation(transform *Transform) {
	transform.sx, transform.sy = 1, 1
	transform.rotation = po.mobj.WorldRotation()

	cx, cy := po.mobj.Pos()
	if !po.interpolate {
		transform.tx, transform.ty = float64(cx), float64(cy)
		return
	}

	if !po.mobj.IsMoving() {
		transform.tx, transform.ty = float64(cx), float64(cy)
		return
	}
	// 上一次Update的坐标点与当前的不一样，说明又Update了，重置LastX和LastY和开始时间
	// 所以每次Update后都要重置LastX,lastY,LastTime，是因为要从重置点开始计算位置插值
	if po.lastX != cx || po.lastY != cy {
		po.lastX = cx
		po.lastY = cy
		po.lastTime = time.Now() //core.GetSyncServTime()
	}
	duration := time.Since(po.lastTime)
	nx, ny := object.DefaultMove(po.mobj, duration)
	transform.tx, transform.ty = float64(nx), float64(ny)
	po.lastInterpolation = *transform
}

// 移动事件处理
func (po *PlayableMoveObject) onEventMove(args ...any) {
	po.Play()
	po.lastX, po.lastY = po.mobj.Pos()
	po.lastTime = time.Now() //core.GetSyncServTime()
	po.interpolate = true
}

// 停止移动事件处理
func (po *PlayableMoveObject) onEventStopMove(args ...any) {
	po.Stop()
	po.interpolate = false
}

// 暫停事件處理
func (po *PlayableMoveObject) onEventPause(args ...any) {
	po.interpolate = false
}

// 恢復事件處理
func (po *PlayableMoveObject) onEventResume(args ...any) {
	po.interpolate = true
	po.lastTime = time.Now()
}

// 炮彈播放對象
type PlayableShell struct {
	*PlayableMoveObject
	shell    object.IShell
	moveInfo object.TrackMoveInfo
	updated  bool
}

// 創建炮彈播放對象
func NewPlayableShell(shell object.IShell, animConfig *MovableObjectAnimConfig) *PlayableShell {
	playable := &PlayableShell{
		PlayableMoveObject: NewPlayableMoveObject(shell, animConfig),
		shell:              shell,
	}
	playable.lastTime = time.Now() //core.GetSyncServTime()
	return playable
}

// 初始化
func (ps *PlayableShell) Init() {
	ps.PlayableMoveObject.Init()
	ps.shell.RegisterLateUpdateEventHandle(ps.onEventLateUpdate)
}

// 反初始化
func (ps *PlayableShell) Uninit() {
	ps.PlayableMoveObject.Uninit()
	ps.shell.UnregisterLateUpdateEventHandle(ps.onEventLateUpdate)
}

// 插值
func (ps *PlayableShell) Interpolation(transform *Transform) {
	if !ps.shell.ShellStaticInfo().TrackTarget {
		ps.PlayableMoveObject.Interpolation(transform)
		return
	}

	transform.rotation = ps.shell.WorldRotation()
	transform.sx, transform.sy = 1, 1
	nx, ny := ps.shell.Pos()
	if !ps.interpolate {
		transform.tx, transform.ty = float64(nx), float64(ny)
		return
	}

	if !ps.mobj.IsMoving() {
		transform.tx, transform.ty = float64(nx), float64(ny)
		return
	}

	if ps.updated {
		nx, ny = ps.moveInfo.X, ps.moveInfo.Y
		ps.updated = false
	} else {
		ps.moveInfo.Rotation = ps.shell.WorldRotation()
		duration := time.Since(ps.lastTime)
		nx, ny = object.GetShellTrackMovedPos(ps.shell, duration, &ps.moveInfo)
	}
	transform.tx, transform.ty = float64(nx), float64(ny)
	transform.rotation = ps.moveInfo.Rotation
	ps.lastInterpolation = *transform
}

// 后更新事件處理
func (ps *PlayableShell) onEventLateUpdate(args ...any) {
	ps.moveInfo.X, ps.moveInfo.Y = args[0].(int32), args[1].(int32)
	ps.moveInfo.Rotation = args[2].(common_base.Angle)
	ps.updated = true
	ps.lastTime = time.Now() //core.GetSyncServTime()
}

// 坦克播放对象
type PlayableTank struct {
	*PlayableMoveObject
	shieldAnim *client_base.SpriteAnim
}

// 创建坦克播放对象
func NewPlayableTank(tank object.ITank, animConfig *MovableObjectAnimConfig) *PlayableTank {
	return &PlayableTank{
		PlayableMoveObject: NewPlayableMoveObject(tank, animConfig),
		shieldAnim:         client_base.NewSpriteAnim(getShieldAnimConfig()),
	}
}

// 初始化
func (pt *PlayableTank) Init() {
	pt.PlayableMoveObject.Init()
	tank := pt.mobj.(object.ITank)
	tank.RegisterChangeEventHandle(pt.onChange)
	tank.RegisterAddShieldEventHandle(pt.onAddShield)
	tank.RegisterCancelShieldEventHandle(pt.onCancelShield)
}

// 反初始化
func (pt *PlayableTank) Uninit() {
	pt.PlayableMoveObject.Uninit()
	tank := pt.mobj.(object.ITank)
	tank.UnregisterChangeEventHandle(pt.onChange)
	tank.UnregisterAddShieldEventHandle(pt.onAddShield)
	tank.UnregisterCancelShieldEventHandle(pt.onCancelShield)
}

// 播放
func (pt *PlayableTank) Play() {
	pt.PlayableMoveObject.Play()
}

// 停止播放
func (pt *PlayableTank) Stop() {
	pt.PlayableMoveObject.Stop()
}

// 更新
func (pt *PlayableTank) Draw(screen *ebiten.Image, op *ebiten.DrawImageOptions) {
	// 顯示根據邏輯數據插值
	op.GeoM.Concat(pt.op.GeoM)
	pt.anim.Update(screen, op)
	tank := pt.mobj.(object.ITank)
	if tank.HasShield() {
		var tmpOp ebiten.DrawImageOptions
		tmpOp.GeoM.Translate(-2, -2)
		tmpOp.GeoM.Concat(op.GeoM)
		tmpOp.ColorScale.SetA(0)
		pt.shieldAnim.Update(screen, &tmpOp)
	}
}

// 变化事件
func (pt *PlayableTank) onChange(args ...any) {
	info := args[0].(*object.ObjStaticInfo)
	pt.currSpeed = (info.Speed())
	pt.Play()
}

// 加護盾事件
func (pt *PlayableTank) onAddShield(args ...any) {
	pt.shieldAnim.Play()
}

// 取消護盾事件
func (pt *PlayableTank) onCancelShield(args ...any) {

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
	var transform Transform
	playableAroundCenterObj.Interpolation(&transform)
	pobj.lastMoveInfo.LastCenterX, pobj.lastMoveInfo.LastCenterY = int32(transform.tx), int32(transform.ty)
	return pobj
}

// 初始化
func (ps *PlayableSurroundObj) Init() {
	ps.PlayableMoveObject.Init()
	ps.sobj.RegisterLateUpdateEventHandle(ps.onEventLateUpdate)
}

// 反初始化
func (ps *PlayableSurroundObj) Uninit() {
	ps.PlayableMoveObject.Uninit()
	ps.sobj.UnregisterLateUpdateEventHandle(ps.onEventLateUpdate)
}

// 插值
func (ps *PlayableSurroundObj) Interpolation(transform *Transform) {
	transform.rotation = ps.mobj.WorldRotation()
	transform.sx, transform.sy = 1, 1

	nx, ny := ps.sobj.Pos()
	if !ps.interpolate {
		transform.tx, transform.ty = float64(nx), float64(ny)
		return
	}
	if !ps.mobj.IsMoving() {
		transform.tx, transform.ty = float64(nx), float64(ny)
		return
	}
	duration := time.Since(ps.lastTime)
	ps.lastTime = time.Now() //core.GetSyncServTime()
	var interpolateX, interpolateY float64
	if pobj, o := ps.playableAroundCenterObj.(IPlayableMovableObject); !o {
		var t Transform
		ps.playableAroundCenterObj.Interpolation(&t)
		interpolateX, interpolateY = t.tx, t.ty
	} else {
		interpolateX, interpolateY = pobj.LastInterpolation()
	}
	ps.lastMoveInfo.LastCenterX, ps.lastMoveInfo.LastCenterY = int32(interpolateX), int32(interpolateY)
	nx, ny = object.GetSurroundObjMovedPos(ps.sobj.(*object.SurroundObj), duration, &ps.lastMoveInfo)
	transform.tx, transform.ty = float64(nx), float64(ny)
	ps.lastInterpolation = *transform
}

// 更新事件處理
func (ps *PlayableSurroundObj) onEventLateUpdate(args ...any) {
	ps.lastTime = time.Now() //core.GetSyncServTime()
	ps.lastMoveInfo.TurnAngle = args[0].(int32)
	ps.lastMoveInfo.AccumulateTime = args[1].(time.Duration)
}

// 可播放效果
type PlayableEffect struct {
	op     ebiten.DrawImageOptions
	effect effect.IEffect
	anim   *client_base.SpriteAnim
}

// 創建可播放效果
func NewPlayableEffect(effect effect.IEffect, animConfig *client_base.SpriteAnimConfig) *PlayableEffect {
	return &PlayableEffect{
		effect: effect,
		anim:   client_base.NewSpriteAnim(animConfig),
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
func (po *PlayableEffect) Interpolation(transform *Transform) {
	x, y := po.effect.Pos()
	transform.tx, transform.ty = float64(x), float64(y)
	transform.sx, transform.sy = 1, 1
}
