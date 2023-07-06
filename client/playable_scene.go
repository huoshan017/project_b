package main

import (
	"project_b/client/base"
	"project_b/common"
	"project_b/common/log"
	"project_b/common/math"
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
	s.scene.RegisterEffectRemovedHandle(s.onEffectRemovedHandle)
}

/**
 * 卸载地图
 */
func (s *PlayableScene) UnloadScene() {
	s.scene.UnregisterEffectRemovedHandle(s.onEffectRemovedHandle)
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
	var rect = math.NewRectObj(lx, ly, rx-lx, ry-ly)
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

	s._draw(tc, obj.Left(), obj.Top(), obj.Right(), obj.Bottom(), obj.Width(), obj.Height(), dstImage)
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

	cx, cy := effect.Center()
	left, top, right, bottom := cx-effect.Width()/2, cy+effect.Height()/2, cx+effect.Width()/2, cy-effect.Height()/2
	s._draw(tc, left, top, right, bottom, effect.Width(), effect.Height(), dstImage)
}

func (s *PlayableScene) _draw(tc *objOpCache, left, top, right, bottom, width, height int32, dstImage *ebiten.Image) {
	tc.op.GeoM.Reset()
	// tile本地坐標到世界坐標的縮放
	sx := width / tc.frameWidth
	sy := height / tc.frameHeight
	tc.op.GeoM.Scale(float64(sx), float64(sy))

	// 移動物體插值
	mapConfig := s.scene.GetMapConfig()
	mapWidth := mapConfig.TileWidth * int32(len(mapConfig.Layers[0]))
	mapHeight := mapConfig.TileHeight * int32(len(mapConfig.Layers))
	var dx, dy float64
	if left >= mapConfig.X && left <= mapConfig.X+mapWidth-width && bottom >= mapConfig.Y && bottom <= mapConfig.Y+mapHeight-height {
		dx, dy = tc.playable.Interpolation()
		left += int32(dx)
		right += int32(dx)
		top += int32(dy)
		bottom += int32(dy)
	}

	// todo 注意这里，i是y轴方向，j是x轴方向
	// 由於世界坐標Y軸與屏幕坐標Y軸方向相反，所以變換左上角和右下角的世界坐標到屏幕坐標
	lx, ly := s.camera.World2Screen(left, top)
	rx, ry := s.camera.World2Screen(right, bottom)
	scalex := float64(rx-lx) / float64(width)
	scaley := float64(ry-ly) / float64(height)
	tc.op.GeoM.Scale(scalex, scaley)
	tc.op.GeoM.Translate(float64(lx), float64(ly))
	// 判断是否插值
	tc.playable.Draw(dstImage, &tc.op)
}

func (s *PlayableScene) onTankRemovedHandle(args ...any) {
	tank := args[0].(*object.Tank)
	delete(s.playableObjs, tank.InstId())
}

func (s *PlayableScene) onBulletRemovedHandle(args ...any) {
	bullet := args[0].(*object.Bullet)
	pobj := s.playableObjs[bullet.InstId()]
	if pobj == nil {
		log.Debug("playable bullet %v not found", bullet.InstId())
	} else {
		pobj.playable.Uninit()
		delete(s.playableObjs, bullet.InstId())
		log.Debug("playable bullet %v removed", bullet.InstId())
	}
}

func (s *PlayableScene) onStaticObjRemovedHandle(args ...any) {
	robj := args[0].(object.IObject)
	// 刪除map中的playable，讓之後的GC回收
	// todo 希望做成對象池可以復用這部分内存
	pobj := s.playableObjs[robj.InstId()]
	if pobj != nil {
		pobj.playable.Uninit()
		delete(s.playableObjs, robj.InstId())
	}
}

func (s *PlayableScene) onEffectRemovedHandle(args ...any) {
	effect := args[0].(object.IEffect)
	// 刪除map中的playable，讓之後的GC回收
	// todo 希望做成對象池可以復用這部分内存
	peffect := s.playableEffects[effect.InstId()]
	if peffect != nil {
		peffect.playable.Uninit()
		delete(s.playableEffects, effect.InstId())
	}
}
