package main

import (
	"math"
	"project_b/client/base"
	"project_b/common"
	pmath "project_b/common/math"
	"project_b/common/object"

	"github.com/hajimehoshi/ebiten/v2"
)

type objOpCache struct {
	op                      ebiten.DrawImageOptions
	playable                IPlayable
	frameWidth, frameHeight int32
}

/**
 * 可绘制场景，实现base.IPlayableScene接口
 */
type PlayableScene struct {
	scene           *common.SceneLogic
	camera          *base.Camera
	viewport        *base.Viewport
	playableObjs    map[uint32]*objOpCache
	playableEffects map[uint32]*objOpCache
}

/**
 * 创建可绘制场景
 */
func CreatePlayableScene(viewport *base.Viewport) *PlayableScene {
	return &PlayableScene{
		viewport:        viewport,
		playableObjs:    make(map[uint32]*objOpCache),
		playableEffects: make(map[uint32]*objOpCache),
	}
}

/**
 * 载入地图
 */
func (s *PlayableScene) SetScene(scene *common.SceneLogic) {
	mapInfo := mapInfoArray[scene.GetMapId()]
	s.camera = base.CreateCamera(s.viewport, mapInfo.cameraFov, defaultCamera2ViewportDistance)
	s.CameraMoveTo(mapInfo.cameraPos.X, mapInfo.cameraPos.Y)
	s.CameraSetHeight(mapInfo.cameraHeight)
	s.scene = scene

	s.scene.RegisterStaticObjRemovedHandle(s.onStaticObjRemovedHandle)
	s.scene.RegisterTankRemovedHandle(s.onTankRemovedHandle)
	s.scene.RegisterBulletRemovedHandle(s.onBulletRemovedHandle)
	s.scene.RegisterSurroundObjRemovedHandle(s.onSurroundObjRemovedHandle)
	s.scene.RegisterEffectRemovedHandle(s.onEffectRemovedHandle)
}

/**
 * 卸载地图
 */
func (s *PlayableScene) UnloadScene() {
	s.scene.UnregisterEffectRemovedHandle(s.onEffectRemovedHandle)
	s.scene.UnregisterSurroundObjRemovedHandle(s.onSurroundObjRemovedHandle)
	s.scene.UnregisterBulletRemovedHandle(s.onBulletRemovedHandle)
	s.scene.UnregisterTankRemovedHandle(s.onTankRemovedHandle)
	s.scene.UnregisterStaticObjRemovedHandle(s.onStaticObjRemovedHandle)
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
		playableObj, animConfig := GetPlayableObject(obj)
		tc = &objOpCache{
			playable:    playableObj,
			op:          ebiten.DrawImageOptions{},
			frameWidth:  int32(animConfig.FrameWidth),
			frameHeight: int32(animConfig.FrameHeight),
		}
		tc.op.ColorScale.SetA(0)
		s.playableObjs[obj.InstId()] = tc
	}

	s._draw(tc, obj.OriginalLeft(), obj.OriginalTop(), obj.OriginalRight(), obj.OriginalBottom(), obj.Width(), obj.Length(), obj.Orientation(), dstImage)
}

func (s *PlayableScene) drawEffect(effect object.IEffect, dstImage *ebiten.Image) {
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

	cx, cy := effect.Pos()
	left, top, right, bottom := cx-effect.Width()/2, cy+effect.Height()/2, cx+effect.Width()/2, cy-effect.Height()/2
	s._draw(tc, left, top, right, bottom, effect.Width(), effect.Height(), 0, dstImage)
}

func (s *PlayableScene) _draw(tc *objOpCache, left, top, right, bottom, width, length int32, orientation int32, dstImage *ebiten.Image) {
	// 移動插值
	dx, dy := tc.playable.Interpolation()
	left += int32(dx)
	right += int32(dx)
	top += int32(dy)
	bottom += int32(dy)
	// todo 注意这里，i是y轴方向，j是x轴方向
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
	// 旋轉
	if orientation > 0 {
		tc.op.GeoM.Translate(-float64(dw)/2, -float64(dh)/2)
		tc.op.GeoM.Rotate(-float64(orientation) * math.Pi / 180.0)
		tc.op.GeoM.Translate(float64(dw)/2, float64(dh)/2)
	}
	tc.op.GeoM.Translate(float64(lx), float64(ly))
	// 判断是否插值
	tc.playable.Draw(dstImage, &tc.op)
}

func (s *PlayableScene) onTankRemovedHandle(args ...any) {
	tank := args[0].(*object.Tank)
	pobj := s.playableObjs[tank.InstId()]
	if pobj != nil {
		pobj.playable.Uninit()
		delete(s.playableObjs, tank.InstId())
	}
}

func (s *PlayableScene) onBulletRemovedHandle(args ...any) {
	bullet := args[0].(*object.Bullet)
	pobj := s.playableObjs[bullet.InstId()]
	if pobj != nil {
		pobj.playable.Uninit()
		delete(s.playableObjs, bullet.InstId())
	}
}

func (s *PlayableScene) onSurroundObjRemovedHandle(args ...any) {
	ball := args[0].(*object.SurroundObj)
	bobj := s.playableObjs[ball.InstId()]
	if bobj != nil {
		bobj.playable.Uninit()
		delete(s.playableObjs, ball.InstId())
	}
}

func (s *PlayableScene) onStaticObjRemovedHandle(args ...any) {
	robj := args[0].(object.IObject)
	// todo 希望做成對象池可以復用這部分内存
	pobj := s.playableObjs[robj.InstId()]
	if pobj != nil {
		pobj.playable.Uninit()
		delete(s.playableObjs, robj.InstId())
	}
}

func (s *PlayableScene) onEffectRemovedHandle(args ...any) {
	effect := args[0].(object.IEffect)
	// todo 希望做成對象池可以復用這部分内存
	peffect := s.playableEffects[effect.InstId()]
	if peffect != nil {
		peffect.playable.Uninit()
		delete(s.playableEffects, effect.InstId())
	}
}
