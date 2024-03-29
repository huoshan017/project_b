package main

import (
	"image/color"
	"math"
	"project_b/client_base"
	"project_b/common"
	"project_b/common/base"
	"project_b/common/effect"
	pmath "project_b/common/math"
	"project_b/common/object"
	"project_b/log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var mainCamera *client_base.Camera

type objOpCache struct {
	op                      ebiten.DrawImageOptions
	playable                IPlayable
	frameWidth, frameHeight int32
}

/**
 * 可绘制场景，实现base.IPlayableScene接口
 */
type PlayableScene struct {
	scene           *common.World
	camera          *client_base.Camera
	viewport        *client_base.Viewport
	playableObjs    map[uint32]*objOpCache
	playableEffects map[uint32]*objOpCache
}

/**
 * 创建可绘制场景
 */
func CreatePlayableScene(viewport *client_base.Viewport) *PlayableScene {
	return &PlayableScene{
		viewport:        viewport,
		playableObjs:    make(map[uint32]*objOpCache),
		playableEffects: make(map[uint32]*objOpCache),
	}
}

/**
 * 载入地图
 */
func (s *PlayableScene) SetScene(world *common.World) {
	mapInfo := mapInfoArray[world.GetMapId()]
	s.camera = client_base.CreateCamera(s.viewport, mapInfo.cameraFov, defaultNearPlane)
	s.CameraMoveTo(world.Center() /*mapInfo.cameraPos.X, mapInfo.cameraPos.Y*/)
	s.CameraSetHeight(mapInfo.cameraHeight)
	s.scene = world
	s.scene.RegisterObjectRemovedHandle(s.onObjRemovedHandle)
	s.scene.RegisterEffectRemovedHandle(s.onEffectRemovedHandle)
	mainCamera = s.camera
}

/**
 * 卸载地图
 */
func (s *PlayableScene) UnloadScene() {
	s.scene.UnregisterEffectRemovedHandle(s.onEffectRemovedHandle)
	s.scene.UnregisterObjectRemovedHandle(s.onObjRemovedHandle)
	clear(s.playableObjs)
	clear(s.playableEffects)
	s.scene = nil
	s.camera = nil
}

/**
 * 移動相機
 */
func (s *PlayableScene) CameraMove(x, y int32) {
	s.camera.Move(x, y)
}

/**
 * 相機移到
 */
func (s *PlayableScene) CameraMoveTo(x, y int32) {
	s.camera.MoveTo(x, y)
}

/**
 * 改變相機高度
 */
func (s *PlayableScene) CameraChangeHeight(delta int32) {
	s.camera.ChangeHeight(delta)
}

/**
 * 設置相機高度
 */
func (s *PlayableScene) CameraSetHeight(height int32) {
	s.camera.SetHeight(height)
}

/**
 * 绘制场景
 */
func (s *PlayableScene) Draw(dstImage *ebiten.Image) {
	s.drawMapGrid(dstImage)
	// 屏幕左下角
	lx, ly := s.camera.Screen2World(0, s.viewport.H())
	// 屏幕右上角
	rx, ry := s.camera.Screen2World(s.viewport.W(), 0)
	// 繪製場景
	var rect = pmath.NewRectObj(lx, ly, rx-lx, ry-ly)
	layerObjs := s.scene.GetLayerObjsWithRange(&rect)
	for i := 0; i < len(layerObjs); i++ {
		if layerObjs[i].Length() == 0 {
			continue
		}
		var (
			o  = true
			id uint32
		)
		// 從大頂堆中取出obj，按邏輯距離由遠到近畫出來
		for o {
			id, _, o = layerObjs[i].Get()
			if !o {
				continue
			}
			obj := s.scene.GetObj(id)
			if obj == nil {
				continue
			}
			s.drawObj(obj, dstImage)
		}
	}
	effectList := s.scene.GetEffectListWithRange(&rect)
	if len(effectList) > 0 {
		for _, effectId := range effectList {
			effect := s.scene.GetEffect(effectId)
			if effect == nil {
				continue
			}
			s.drawEffect(effect, dstImage)
		}
	}
}

func (s *PlayableScene) drawObj(obj object.IObject, dstImage *ebiten.Image) {
	tc := s.playableObjs[obj.InstId()]
	if tc == nil {
		playableObj, animConfig := s.getPlayableObject(obj, dstImage)
		tc = &objOpCache{
			playable:    playableObj,
			op:          ebiten.DrawImageOptions{},
			frameWidth:  int32(animConfig.FrameWidth),
			frameHeight: int32(animConfig.FrameHeight),
		}
		tc.op.ColorScale.SetA(0)
		s.playableObjs[obj.InstId()] = tc
	}

	s.draw(tc, obj.Width(), obj.Length(), dstImage)
}

func (s *PlayableScene) drawEffect(effect effect.IEffect, dstImage *ebiten.Image) {
	tc := s.playableEffects[effect.InstId()]
	if tc == nil {
		animConfig := getEffectAnimConfig(effect.StaticInfo().Id)
		tc = &objOpCache{
			playable:    NewPlayableEffect(effect, animConfig),
			op:          ebiten.DrawImageOptions{},
			frameWidth:  int32(animConfig.FrameWidth),
			frameHeight: int32(animConfig.FrameHeight),
		}
		tc.op.ColorScale.SetA(0)
		s.playableEffects[effect.InstId()] = tc
	}

	s.draw(tc, effect.Width(), effect.Height(), dstImage)
}

func (s *PlayableScene) draw(tc *objOpCache, width, length int32, dstImage *ebiten.Image) {
	// 插值
	var transform Transform
	tc.playable.Interpolation(&transform)
	left := int32(transform.tx) - width/2
	right := int32(transform.tx) + width/2
	top := int32(transform.ty) + length/2
	bottom := int32(transform.ty) - length/2

	// 由於世界坐標Y軸與屏幕坐標Y軸方向相反，所以變換左上角和右下角的世界坐標到屏幕坐標
	// 遵循縮放、旋轉、平移的變換順序
	tc.op.GeoM.Reset()
	// tile本地坐標到世界坐標的縮放
	sx := width / tc.frameWidth
	sy := length / tc.frameHeight
	tc.op.GeoM.Scale(float64(sx), float64(sy))
	lx, ly := s.camera.World2Screen(left, top)
	rx, ry := s.camera.World2Screen(right, bottom)
	dw, dh := rx-lx, ry-ly
	scalex := float64(dw) / float64(width)
	scaley := float64(dh) / float64(length)
	tc.op.GeoM.Scale(scalex, scaley)

	// todo 這個旋轉指的是基於圖片的旋轉，是順時針旋轉的變化量，不同於以x軸正方向為零度，逆時針旋轉為正向的邏輯旋轉
	minutes := transform.rotation.ToMinutes()
	if minutes != 0 {
		tc.op.GeoM.Translate(-float64(dw)/2, -float64(dh)/2)
		tc.op.GeoM.Rotate(-float64(minutes) * math.Pi / (60 * 180.0))
		tc.op.GeoM.Translate(float64(dw)/2, float64(dh)/2)
	}
	tc.op.GeoM.Translate(float64(lx), float64(ly))
	tc.playable.Draw(dstImage, &tc.op)
}

func (s *PlayableScene) drawMapGrid(dstImage *ebiten.Image) {
	if !debug.IsShowMapGrid() {
		return
	}
	mw, mh := s.scene.GetMapWidthHeight()
	left, bottom := s.scene.GetMapLeftBottom()
	gw, gh := s.scene.GetGridMap().GetGridWidthHeight()
	var x, y int32
	for x = left; x <= left+mw; x += gw {
		sx0, sy0 := s.camera.World2Screen(x, bottom)
		sx1, sy1 := s.camera.World2Screen(x, bottom+mh)
		vector.StrokeLine(dstImage, float32(sx0), float32(sy0), float32(sx1), float32(sy1), 1, color.RGBA{100, 100, 100, 0}, false)
	}
	for y = bottom; y <= bottom+mh; y += gh {
		sx0, sy0 := s.camera.World2Screen(left, y)
		sx1, sy1 := s.camera.World2Screen(left+mw, y)
		vector.StrokeLine(dstImage, float32(sx0), float32(sy0), float32(sx1), float32(sy1), 1, color.RGBA{100, 100, 100, 0}, false)
	}
}

func (s *PlayableScene) getPlayableObject(obj object.IObject, dstImage *ebiten.Image) (IPlayable, *client_base.SpriteAnimConfig) {
	var (
		playableObj IPlayable
		animConfig  *client_base.SpriteAnimConfig
	)
	switch obj.Type() {
	case base.ObjTypeStatic:
		if base.StaticObjType(obj.Subtype()) == base.StaticObjNone {
			return nil, nil
		}
		config := GetStaticObjAnimConfig(base.StaticObjType(obj.Subtype()))
		if config == nil {
			log.Error("can't get static object anim by subtype %v", obj.Subtype())
			return nil, nil
		}
		playableObj = NewPlayableStaticObject(obj, config)
		animConfig = config.AnimConfig
	case base.ObjTypeItem:
		if base.ItemObjType(obj.Subtype()) == base.ItemObjNone {
			return nil, nil
		}
		config := getItemObjAnimConfig(base.ItemObjType(obj.Subtype()))
		if config == nil {
			log.Error("can't get item object anim by subtype %v", obj.Subtype())
			return nil, nil
		}
		playableObj = NewPlayableItemObject(obj, config)
		animConfig = config
	case base.ObjTypeMovable:
		if base.MovableObjType(obj.Subtype()) == base.MovableObjNone {
			return nil, nil
		}
		mobj := obj.(object.IMovableObject)
		config := GetMovableObjAnimConfig(base.MovableObjType(obj.Subtype()), mobj.Id(), mobj.Level())
		if config == nil {
			log.Error("can't get movable object anim by subtype %v", obj.Subtype())
			return nil, nil
		}
		switch obj.Subtype() {
		case base.ObjSubtypeTank:
			playableObj = NewPlayableTank(mobj.(object.ITank), config)
		case base.ObjSubtypeShell:
			playableObj = NewPlayableShell(mobj.(object.IShell), config)
		case base.ObjSubtypeSurroundObj:
			surroundObj := mobj.(object.ISurroundObject)
			aroundCenterObj := surroundObj.GetAroundCenterObject()
			// 環繞物體需要先創建被環繞物體
			cobjCache := s.playableObjs[aroundCenterObj.InstId()]
			if cobjCache == nil {
				s.drawObj(aroundCenterObj, dstImage)
				cobjCache = s.playableObjs[aroundCenterObj.InstId()]
			}
			playableObj = NewPlayableSurroundObj(surroundObj, config, cobjCache.playable)
		default:
			playableObj = NewPlayableMoveObject(mobj, config)
		}
		playableObj.Init()
		animConfig = config.AnimConfig
	}
	return playableObj, animConfig
}

func (s *PlayableScene) onObjRemovedHandle(args ...any) {
	obj := args[0].(object.IObject)
	cache := s.playableObjs[obj.InstId()]
	if cache != nil {
		cache.playable.Uninit()
		delete(s.playableObjs, obj.InstId())
	}
}

func (s *PlayableScene) onEffectRemovedHandle(args ...any) {
	effect := args[0].(effect.IEffect)
	// todo 希望做成對象池可以復用這部分内存
	peffect := s.playableEffects[effect.InstId()]
	if peffect != nil {
		peffect.playable.Uninit()
		delete(s.playableEffects, effect.InstId())
	}
}
