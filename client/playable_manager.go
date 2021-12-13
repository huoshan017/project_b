package main

import (
	"project_b/common/object"

	"github.com/hajimehoshi/ebiten/v2"
)

// 播放管理器
type PlayableManager struct {
	playerTankPlayables map[uint64]*PlayableTank
	enemyTankPlayables  map[int32]*PlayableTank
}

// 创建播放管理器
func CreatePlayableManager() *PlayableManager {
	return &PlayableManager{
		playerTankPlayables: make(map[uint64]*PlayableTank),
		enemyTankPlayables:  make(map[int32]*PlayableTank),
	}
}

// 添加玩家坦克播放
func (m *PlayableManager) AddPlayerTankPlayable(uid uint64, tank *object.Tank) {
	if tank == nil {
		panic("player tank is nil")
	}
	_, o := m.playerTankPlayables[uid]
	if !o {
		p := NewPlayableTank(tank, GetTankAnimConfig(tank.Id(), tank.Level()))
		p.Init()
		m.playerTankPlayables[uid] = p
	}
}

// 删除玩家坦克播放（放入对应坦克id播放的freelist）
func (m *PlayableManager) RemovePlayerTankPlayable(uid uint64) {
	playable, o := m.playerTankPlayables[uid]
	if !o {
		return
	}
	playable.Uninit()
	delete(m.playerTankPlayables, uid)
}

// 开始播放玩家坦克动画
func (m *PlayableManager) PlayPlayerTankPlayable(uid uint64) {
	playable, o := m.playerTankPlayables[uid]
	if o {
		playable.Play()
	}
}

// 停止播放玩家坦克动画
func (m *PlayableManager) StopPlayerTankPlayable(uid uint64) {
	playable, o := m.playerTankPlayables[uid]
	if o {
		playable.Stop()
	}
}

// 播放玩家们的坦克动画
func (m *PlayableManager) PlayPlayersTankPlayable() {
	for uid := range m.playerTankPlayables {
		m.PlayPlayerTankPlayable(uid)
	}
}

// 停止播放玩家们的坦克动画
func (m *PlayableManager) StopPlayersTankPlayable() {
	for uid := range m.playerTankPlayables {
		m.StopPlayerTankPlayable(uid)
	}
}

// 添加敌人坦克播放
func (m *PlayableManager) AddEnemyTankPlayable(instId int32, tank *object.Tank) {
	_, o := m.enemyTankPlayables[instId]
	if !o {
		p := NewPlayableTank(tank, GetTankAnimConfig(tank.Id(), tank.Level()))
		p.Init()
		m.enemyTankPlayables[instId] = p
	}
}

// 删除玩家坦克播放
func (m *PlayableManager) RemoveEnemyTankPlayable(instId int32) {
	playable, o := m.enemyTankPlayables[instId]
	if !o {
		return
	}
	playable.Uninit()
	delete(m.enemyTankPlayables, instId)
}

// 播放敌人坦克动画
func (m *PlayableManager) PlayEnemyTankPlayable(instId int32) {
	playable, o := m.enemyTankPlayables[instId]
	if o {
		playable.Play()
	}
}

// 停止播放敌人坦克动画
func (m *PlayableManager) StopEnemyTankPlayable(id int32) {
	playable, o := m.enemyTankPlayables[id]
	if o {
		playable.Stop()
	}
}

// 播放所有敌人坦克的动画
func (m *PlayableManager) PlayEnemiesTankPlayable() {
	for instId := range m.enemyTankPlayables {
		m.PlayEnemyTankPlayable(instId)
	}
}

// 停止播放所有敌人坦克的动画
func (m *PlayableManager) StopEnemiesTankPlayable() {
	for instId := range m.enemyTankPlayables {
		m.StopEnemyTankPlayable(instId)
	}
}

// 播放所有坦克动画
func (m *PlayableManager) PlayAllTanksPlayable() {
	m.PlayPlayersTankPlayable()
	m.PlayEnemiesTankPlayable()
}

// 停止所有坦克动画
func (m *PlayableManager) StopAllTanksPlayable() {
	m.StopPlayersTankPlayable()
	m.StopEnemiesTankPlayable()
}

// 更新
func (m *PlayableManager) Update(screen *ebiten.Image) {
	m.UpdatePlayerTanksPlayable(screen)
	m.UpdatEnemyTanksPlayable(screen)
}

// 更新玩家坦克动画
func (m *PlayableManager) UpdatePlayerTanksPlayable(screen *ebiten.Image) {
	for _, p := range m.playerTankPlayables {
		p.Update(screen)
	}
}

// 更新敌人坦克动画
func (m *PlayableManager) UpdatEnemyTanksPlayable(screen *ebiten.Image) {
	for _, p := range m.enemyTankPlayables {
		p.Update(screen)
	}
}
