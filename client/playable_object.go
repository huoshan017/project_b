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
	anim base.SpriteAnim
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
		anim: *base.NewSpriteAnim(spriteConfig),
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
	return 0, 0
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
}

// 可移动物体的播放对象，有四个方向的动画
type PlayableMoveObject struct {
	op           *ebiten.DrawImageOptions
	mobj         object.IMovableObject
	anims        []*base.SpriteAnim
	currSpeed    int32           // 当前速度
	lastTime     time.CustomTime // 更新时间点
	interpolate  bool            // 上次是停止状态
	lastX, lastY int32
}

// 创建可移动物体的播放对象
func NewPlayableMoveObject(mobj object.IMovableObject, animConfig *MovableObjectAnimConfig) *PlayableMoveObject {
	pobj := &PlayableMoveObject{
		op:          &ebiten.DrawImageOptions{},
		mobj:        mobj,
		interpolate: mobj.IsMoving(),
	}
	pobj.changeAnim(animConfig)
	return pobj
}

// 改变动画
func (po *PlayableMoveObject) changeAnim(animConfig *MovableObjectAnimConfig) {
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
}

// 反初始化
func (po *PlayableMoveObject) Uninit() {
	// 注销移动停止更新事件
	po.mobj.UnregisterMoveEventHandle(po.onEventMove)
	po.mobj.UnregisterStopMoveEventHandle(po.onEventStopMove)
}

// 播放
func (po *PlayableMoveObject) Play() {
	po.anims[po.mobj.Dir()].Play()
}

// 停止
func (po *PlayableMoveObject) Stop() {
	po.anims[po.mobj.Dir()].Stop()
}

// 更新
func (po *PlayableMoveObject) Draw(screen *ebiten.Image, op *ebiten.DrawImageOptions) {
	op.GeoM.Concat(po.op.GeoM)
	po.anims[po.mobj.Dir()].Update(screen, op)
}

// 插值，在Draw前同步调用，得到位置插值
// ----|------------|------------|------------|------------|-------------|--------------|--------------|--------------|--------------|----
//
//	Draw      Update(Draw)    Draw         Draw         Draw       Update(Draw)      Draw           Draw           Draw        Update(Draw)
func (po *PlayableMoveObject) Interpolation() (x, y float64) {
	if !po.interpolate {
		return
	}
	if !po.mobj.IsMoving() {
		return
	}
	// 上一次Update的坐标点与当前的不一样，说明又Update了，重置LastX和LastY和开始时间
	// 所以每次Update后都要重置LastX,lastY,LastTime，是因为要从重置点开始计算位置插值
	if po.lastX != po.mobj.Left() || po.lastY != po.mobj.Bottom() {
		po.lastX = po.mobj.Left()
		po.lastY = po.mobj.Bottom()
		po.lastTime = core.GetSyncServTime()
	}
	duration := time.Since(po.lastTime)
	distance := object.GetMoveDistance(po.mobj, duration)
	switch po.mobj.Dir() {
	case object.DirLeft:
		x = -distance
	case object.DirRight:
		x = distance
	case object.DirDown:
		y = -distance
	case object.DirUp:
		y = distance
	}
	return
}

// 移动事件处理
func (po *PlayableMoveObject) onEventMove(args ...any) {
	po.Play()
	po.lastX = po.mobj.Left()
	po.lastY = po.mobj.Bottom()
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
	level := args[1].(int32)
	pt.currSpeed = (info.Speed())
	pt.changeAnim(GetTankAnimConfig(info.Id(), level))
	pt.Play()
}

type PlayableEffect struct {
	op     ebiten.DrawImageOptions
	effect object.IEffect
	anim   *base.SpriteAnim
}

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
	return 0, 0
}

func GetPlayableObject(obj object.IObject) (IPlayable, *base.SpriteAnimConfig) {
	var (
		playableObj IPlayable
		animConfig  *base.SpriteAnimConfig
	)
	switch obj.Type() {
	case object.ObjTypeStatic:
		if object.StaticObjType(obj.Subtype()) == object.StaticObjNone {
			return nil, nil
		}
		config := GetStaticObjAnimConfig(object.StaticObjType(obj.Subtype()))
		if config == nil {
			glog.Error("can't get static object anim by subtype %v", obj.Subtype())
			return nil, nil
		}
		playableObj = NewPlayableStaticObject(obj, config)
		animConfig = config.AnimConfig
	case object.ObjTypeMovable:
		if object.MovableObjType(obj.Subtype()) == object.MovableObjNone {
			return nil, nil
		}
		mobj := obj.(object.IMovableObject)
		config := GetMovableObjAnimConfig(object.MovableObjType(obj.Subtype()), mobj.Id(), mobj.Level())
		if config == nil {
			glog.Error("can't get movable object anim by subtype %v", obj.Subtype())
			return nil, nil
		}
		if obj.Subtype() == object.ObjSubTypeTank {
			playableObj = NewPlayableTank(mobj.(object.ITank), config)
		} else if obj.Subtype() == object.ObjSubTypeBullet {
			playableObj = NewPlayableMoveObject(mobj, config)
		}
		playableObj.Init()
		animConfig = config.AnimConfig[mobj.Dir()]
	}
	return playableObj, animConfig
}
