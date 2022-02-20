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
	viewport            *base.Viewport
	playerTankPlayables map[uint64]*PlayableTank
	enemyTankPlayables  map[int32]*PlayableTank
	playableMap         *PlayableMap
}

/**
 * 创建可绘制场景
 */
func CreatePlayableScene() *PlayableScene {
	return &PlayableScene{
		playerTankPlayables: make(map[uint64]*PlayableTank),
		enemyTankPlayables:  make(map[int32]*PlayableTank),
		playableMap:         CreatePlayableMap(),
	}
}

/**
 * 使用视口创建可绘制场景
 */
func CreatePlayableSceneWithViewport(viewport *base.Viewport) *PlayableScene {
	ps := CreatePlayableScene()
	ps.viewport = viewport
	return ps
}

/**
 * 设置相机
 */
func (s *PlayableScene) SetViewport(viewport *base.Viewport) {
	s.viewport = viewport
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
func (s *PlayableScene) Draw(rect *base.Rect, op *ebiten.DrawImageOptions, dstImage *ebiten.Image) {
	// 先绘制地图
	s.playableMap.Draw(rect, op, dstImage)
	// 再绘制其他物体
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
