package main

import (
	"project_b/common/object"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

// 坦克包含四个方向的播放对象
type tankDirPlayable struct {
	left  IPlayable
	right IPlayable
	up    IPlayable
	down  IPlayable
}

// 坦克等级方向播放对象
type tankLevelDirPlayable struct {
	tankId         int32              // 坦克Id，用作tank对象中静态数据变化之后之前的id查询
	levelPlayables []*tankDirPlayable // 等级索引播放对象
	tank           *object.Tank       // 当前坦克对象
	lastPlayable   IPlayable          // 上一次的播放对象
}

// 当前正在播放的
func (p *tankLevelDirPlayable) getCurrPlayable() IPlayable {
	level := p.tank.Level()
	dir := p.tank.Dir()
	levelPlayable := p.levelPlayables[level-1]
	switch dir {
	case object.DirUp:
		return levelPlayable.up
	case object.DirDown:
		return levelPlayable.down
	case object.DirLeft:
		return levelPlayable.left
	case object.DirRight:
		return levelPlayable.right
	}
	return nil
}

// 停止播放
func (p *tankLevelDirPlayable) stop() {
	playable := p.getCurrPlayable()
	playable.Stop()
}

// 重置坦克
func (p *tankLevelDirPlayable) resetTank(tank *object.Tank) bool {
	if p.tank == tank || tank.Id() != p.tankId {
		return false
	}
	p.tank = tank
	for i := 0; i < len(p.levelPlayables); i++ {
	}
	return true
}

// 播放管理器
type PlayableManager struct {
	playerTankPlayables map[uint64]*PlayableMoveObject
	enemyTankPlayables  map[int32]*PlayableMoveObject
	lastCheckTime       time.Time
}

// 创建播放管理器
func CreatePlayableManager() *PlayableManager {
	return &PlayableManager{
		playerTankPlayables: make(map[uint64]*PlayableMoveObject),
		enemyTankPlayables:  make(map[int32]*PlayableMoveObject),
	}
}

// 添加玩家坦克播放
func (m *PlayableManager) AddPlayerTankLevelDirPlayable(uid uint64, tank *object.Tank) {
	if tank == nil {
		panic("player tank is nil")
	}
	_, o := m.playerTankPlayables[uid]
	if !o {
		p := NewPlayableMoveObject(tank, GetTankAnimConfig(tank.Id(), tank.Level()))
		p.Init()
		m.playerTankPlayables[uid] = p
	}
}

// 删除玩家坦克播放（放入对应坦克id播放的freelist）
func (m *PlayableManager) RemovePlayerTankLevelDirPlayable(uid uint64) {
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
		p := NewPlayableMoveObject(tank, GetTankAnimConfig(tank.Id(), tank.Level()))
		p.Init()
		m.enemyTankPlayables[instId] = p
	}
}

// 删除玩家坦克播放
func (m *PlayableManager) RemoveEnemyTankPlayable(instId int32) {
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
	var tick time.Duration
	now := time.Now()
	if !m.lastCheckTime.IsZero() {
		tick = now.Sub(m.lastCheckTime)
	}
	m.UpdatePlayerTanksPlayable(tick, screen)
	m.UpdatEnemyTanksPlayable(tick, screen)
	m.lastCheckTime = now
}

// 更新玩家坦克动画
func (m *PlayableManager) UpdatePlayerTanksPlayable(tick time.Duration, screen *ebiten.Image) {
	for _, p := range m.playerTankPlayables {
		p.Update(tick, screen)
	}
}

// 更新敌人坦克动画
func (m *PlayableManager) UpdatEnemyTanksPlayable(tick time.Duration, screen *ebiten.Image) {
	for _, p := range m.enemyTankPlayables {
		p.Update(tick, screen)
	}
}
