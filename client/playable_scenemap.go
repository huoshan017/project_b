ackage main

import (
	"project_b/client/base"
	"project_b/common"
	"project_b/common/math"
	"project_b/common/object"

	"github.com/hajimehoshi/ebiten/v2"
)

type objOpCache struct {
	op                      *ebiten.DrawImageOptions
	playableObj             IPlayable
	frameWidth, frameHeight int32
}

/**
 * 可绘制场景，实现base.IPlayableScene接口
 */
type PlayableSceneMap struct {
	sceneMap     *common.SceneMap
	camera       *base.Camera
	viewport     *base.Viewport
	playableObjs map[uint32]*objOpCache
}

/**
 * 创建可绘制场景
 */
func CreatePlayableSceneMap(viewport *base.Viewport) *PlayableSceneMap {
	return &PlayableSceneMap{
		viewport:     viewport,
		playableObjs: make(map[uint32]*objOpCache),
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
		if layerObjs[i].Length() == 0 {
			continue
		}
		// todo 這裏正確的做法是根據obj的邏輯距離由遠到近畫出來
		id, _, o := layerObjs[i].Get()
		for o {
			obj := s.sceneMap.GetObj(id)
			if obj == nil {
				continue
			}
			s.drawObj(obj, dstImage)
			id, _, o = layerObjs[i].Get()
		}
	}
}

func (s *PlayableSceneMap) drawObj(obj object.IObject, dstImage *ebiten.Image) {
	tc := s.playableObjs[obj.InstId()]
	if tc == nil {
		playableObj, animConfig := GetPlayableObject(obj)
		tc = &objOpCache{
			playableObj: playableObj,
			op:          &ebiten.DrawImageOptions{},
			frameWidth:  int32(animConfig.FrameWidth),
			frameHeight: int32(animConfig.FrameHeight),
		}
		s.playableObjs[obj.InstId()] = tc
	}

	tc.op.GeoM.Reset()
	// tile本地坐標到世界坐標的縮放
	sx := obj.Width() / tc.frameWidth
	sy := obj.Height() / tc.frameHeight
	tc.op.GeoM.Scale(float64(sx), float64(sy))

	// 插值
	x, y := obj.Pos()
	mapConfig := s.sceneMap.GetMapConfig()
	mapWidth := mapConfig.TileWidth * int32(len(mapConfig.Layers[0]))
	mapHeight := mapConfig.TileHeight * int32(len(mapConfig.Layers))
	var dx, dy float64
	if x >= mapConfig.X && x <= mapConfig.X+mapWidth-obj.Width() && y >= mapConfig.Y && y <= mapConfig.Y+mapHeight-obj.Height() {
		dx, dy = tc.playableObj.Interpolation()
	}
	// todo 注意这里，i是y轴方向，j是x轴方向
	// 由於世界坐標Y軸與屏幕坐標Y軸方向相反，所以變換左上角和右下角的世界坐標到屏幕坐標
	lx, ly := s.camera.World2Screen(obj.Left()+int32(dx), obj.Top()+int32(dy))
	rx, ry := s.camera.World2Screen(obj.Right()+int32(dx), obj.Bottom()+int32(dy))
	scalex := float64(rx-lx) / float64(obj.Width())
	scaley := float64(ry-ly) / float64(obj.Height())
	tc.op.GeoM.Scale(scalex, scaley)
	tc.op.GeoM.Translate(float64(lx), float64(ly))
	// 判断是否插值
	tc.playableObj.Draw(dstImage, tc.op)
}
