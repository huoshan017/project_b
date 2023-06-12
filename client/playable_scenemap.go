package main

import (
	"project_b/client/base"
	"project_b/common"
	"project_b/common/math"
	"project_b/common/object"
	"unsafe"

	"github.com/hajimehoshi/ebiten/v2"
)

type objOpCache struct {
	op          *ebiten.DrawImageOptions
	playableObj *PlayableObject
}

/**
 * 可绘制场景，实现base.IPlayableScene接口
 */
type PlayableSceneMap struct {
	sceneMap            *common.SceneMap
	camera              *base.Camera
	viewport            *base.Viewport
	playableObjs        map[uint32]*objOpCache
	playerTankPlayables map[uint64]*PlayableTank
	enemyTankPlayables  map[int32]*PlayableTank
}

/**
 * 创建可绘制场景
 */
func CreatePlayableSceneMap(viewport *base.Viewport) *PlayableSceneMap {
	return &PlayableSceneMap{
		viewport:            viewport,
		playableObjs:        make(map[uint32]*objOpCache),
		playerTankPlayables: make(map[uint64]*PlayableTank),
		enemyTankPlayables:  make(map[int32]*PlayableTank),
	}
}

/**
 * 载入地图
 */
func (s *PlayableSceneMap) SetMap(sceneMap *common.SceneMap) {
	mapInfo := mapInfoArray[sceneMap.GetMapId()]
	s.camera = base.CreateCamera(s.viewport, mapInfo.cameraFov, defaultCamera2ViewportDistance)
	s.CameraMoveTo(mapInfo.cameraPos.X, mapInfo.cameraPos.Y)
	s.CameraSetHeight(mapInfo.cameraHeight)
	s.sceneMap = sceneMap
}

/**
 * 卸载地图
 */
func (s *PlayableSceneMap) UnloadMap() {
}

/**
 * 移動相機
 */
func (s *PlayableSceneMap) CameraMove(x, y int32) {
	s.camera.Move(x, y)
}

/**
 * 相機移到
 */
func (s *PlayableSceneMap) CameraMoveTo(x, y int32) {
	s.camera.MoveTo(x, y)
}

/**
 * 改變相機高度
 */
func (s *PlayableSceneMap) CameraChangeHeight(delta int32) {
	s.camera.ChangeHeight(delta)
}

/**
 * 設置相機高度
 */
func (s *PlayableSceneMap) CameraSetHeight(height int32) {
	s.camera.SetHeight(height)
}

/**
 * 绘制场景
 */
func (s *PlayableSceneMap) Draw(dstImage *ebiten.Image) {
	// 屏幕左下角
	lx, ly := s.camera.Screen2World(0, s.viewport.H())
	// 屏幕右上角
	rx, ry := s.camera.Screen2World(s.viewport.W(), 0)
	// 繪製場景圖
	var rect = math.NewRectObj(lx, ly, rx-lx, ry-ly)
	layerObjs := s.sceneMap.GetLayerObjsWithRange(&rect)
	for i := 0; i < len(layerObjs); i++ {
		if len(layerObjs[i]) == 0 {
			continue
		}
		for j := 0; j < len(layerObjs[i]); j++ {
			obj := s.sceneMap.GetObj(layerObjs[i][j])
			if obj == nil {
				continue
			}
			s.drawObj(obj, dstImage)
		}
	}
}

func (s *PlayableSceneMap) drawObj(obj object.IObject, dstImage *ebiten.Image) {
	var animConfig *base.SpriteAnimConfig
	switch obj.Type() {
	case object.ObjTypeStatic:
		if object.StaticObjType(obj.Subtype()) == object.StaticObjNone {
			return
		}
		config := GetStaticObjAnimConfig(object.StaticObjType(obj.Subtype()))
		if config == nil {
			glog.Error("can't get static object anim by type %v", obj.Subtype())
			return
		}
		animConfig = config.AnimConfig
	case object.ObjTypeMovable:
		if object.MovableObjType(obj.Subtype()) == object.MovableObjNone {
			return
		}
		var level int32
		mobj := obj.(*object.MovableObject)
		if mobj.Subtype() == object.ObjSubType(object.MovableObjTank) {
			tobj := (*object.Tank)(unsafe.Pointer(mobj))
			level = tobj.Level()
		}
		config := GetMovableObjAnimConfig(object.MovableObjType(obj.Subtype()), mobj.Id(), level)
		if config == nil {
			glog.Error("can't get static object anim by type %v", obj.Subtype())
			return
		}
		animConfig = config.AnimConfig[mobj.Dir()]
	default:
		return
	}
	tc := s.playableObjs[obj.InstId()]
	if tc == nil {
		tc = &objOpCache{
			playableObj: NewPlayableObject(obj, animConfig),
			op:          &ebiten.DrawImageOptions{},
		}
		tc.playableObj.Play()
		s.playableObjs[obj.InstId()] = tc
	}
	tc.op.GeoM.Reset()
	// tile本地坐標到世界坐標的縮放
	sx := obj.Width() / int32(animConfig.FrameWidth)
	sy := obj.Height() / int32(animConfig.FrameHeight)
	tc.op.GeoM.Scale(float64(sx), float64(sy))
	// todo 注意这里，i是y轴方向，j是x轴方向
	// 由於世界坐標Y軸與屏幕坐標Y軸方向相反，所以變換左上角和右下角的世界坐標到屏幕坐標
	lx, ly := s.camera.World2Screen(obj.Left(), obj.Top())
	rx, ry := s.camera.World2Screen(obj.Right(), obj.Bottom())
	scalex := float64(rx-lx) / float64(obj.Width())
	scaley := float64(ry-ly) / float64(obj.Height())
	tc.op.GeoM.Scale(scalex, scaley)
	tc.op.GeoM.Translate(float64(lx), float64(ly))
	tc.playableObj.Draw(dstImage, tc.op)
}

func (m *PlayableSceneMap) AddPlayerTankPlayable(uid uint64, tank *object.Tank) bool {
	return true
}

func (m *PlayableSceneMap) RemovePlayerTankPlayable(uid uint64) bool {
	return true
}

func (m *PlayableSceneMap) AddEnemyTankPlayable(id int32, tank *object.Tank) bool {
	return true
}
