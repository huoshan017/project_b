package main

import (
	"project_b/client/base"
	"project_b/common/object"

	"github.com/hajimehoshi/ebiten/v2"
)

/**
 * 可绘制场景，实现base.IPlayableScene接口
 */
type PlayableScene struct {
	camera              *base.Camera
	viewport            *base.Viewport
	playerTankPlayables map[uint64]*PlayableTank
	enemyTankPlayables  map[int32]*PlayableTank
	playableMap         *PlayableMap
}

/**
 * 创建可绘制场景
 */
func CreatePlayableScene(viewport *base.Viewport, camera *base.Camera) *PlayableScene {
	return &PlayableScene{
		camera:              camera,
		viewport:            viewport,
		playerTankPlayables: make(map[uint64]*PlayableTank),
		enemyTankPlayables:  make(map[int32]*PlayableTank),
		playableMap:         CreatePlayableMap(camera),
	}
}

/**
 * 载入地图
 */
func (s *PlayableScene) LoadMap(mapId int32, objArray [][]*object.StaticObject) bool {
	return s.playableMap.Load(mapId, objArray)
}

/**
 * 卸载地图
 */
func (s *PlayableScene) UnloadMap() {

}

/**
 * 绘制场景
 */
func (s *PlayableScene) Draw( /*rect *base.Rect, op *ebiten.DrawImageOptions, */ dstImage *ebiten.Image) {
	// 屏幕左下角
	lx, ly := s.camera.Screen2World(0, s.viewport.H())
	// 屏幕右上角
	rx, ry := s.camera.Screen2World(s.viewport.W(), 0)
	// 繪製場景圖
	s.playableMap.Draw(base.NewRect(lx, ly, rx-lx, ry-ly), dstImage)
	// 再绘制其他物体
	//s.drawEnemyTanksPlayable(dstImage, op)
	//s.drawPlayerTanksPlayable(dstImage, op)
}

// 更新玩家坦克动画
func (m *PlayableScene) drawPlayerTanksPlayable(screen *ebiten.Image, op *ebiten.DrawImageOptions) {
	for _, p := range m.playerTankPlayables {
		p.Draw(screen, op)
	}
}

// 更新敌人坦克动画
func (m *PlayableScene) drawEnemyTanksPlayable(screen *ebiten.Image, op *ebiten.DrawImageOptions) {
	for _, p := range m.enemyTankPlayables {
		p.Draw(screen, op)
	}
}

func (m *PlayableScene) AddPlayerTankPlayable(uid uint64, tank *object.Tank) bool {
	return true
}

func (m *PlayableScene) RemovePlayerTankPlayable(uid uint64) bool {
	return true
}

func (m *PlayableScene) AddEnemyTankPlayable(id int32, tank *object.Tank) bool {
	return true
}
