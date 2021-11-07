package main

import (
	"fmt"
	"project_b/common/object"
	"project_b/common_data"

	"github.com/hajimehoshi/ebiten/v2"
)

// 坦克包含四个方向的播放对象
type tankDirPlayable struct {
	left  *PlayableObject
	right *PlayableObject
	up    *PlayableObject
	down  *PlayableObject
}

// 坦克等级方向播放对象
type tankLevelDirPlayable struct {
	tankId         int32              // 坦克Id，用作tank对象中静态数据变化之后之前的id查询
	levelPlayables []*tankDirPlayable // 等级索引播放对象
	tank           *object.Tank       // 当前坦克对象
	lastPlayable   *PlayableObject    // 上一次的播放对象
}

// 当前正在播放的
func (p *tankLevelDirPlayable) getCurrPlayable() *PlayableObject {
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
		d := p.levelPlayables[i]
		d.left.ResetObj(tank)
		d.right.ResetObj(tank)
		d.up.ResetObj(tank)
		d.down.ResetObj(tank)
	}
	return true
}

// 播放管理器
type PlayableManager struct {
	playerTankPlayables map[uint64]*tankLevelDirPlayable  // 唯一uid映射到*PlayableObject
	enemyTankPlayables  map[int32]*tankLevelDirPlayable   // 实例id映射到*PlayableObject
	freeTankPlayables   map[int32][]*tankLevelDirPlayable // 空闲的坦克等级方向播放对象，按照坦克配置id索引
}

// 创建播放管理器
func CreatePlayableManager() *PlayableManager {
	return &PlayableManager{
		playerTankPlayables: make(map[uint64]*tankLevelDirPlayable),
		enemyTankPlayables:  make(map[int32]*tankLevelDirPlayable),
		freeTankPlayables:   make(map[int32][]*tankLevelDirPlayable),
	}
}

// 添加玩家坦克播放
func (m *PlayableManager) AddPlayerTankLevelDirPlayable(uid uint64, tank *object.Tank) {
	if tank == nil {
		panic("player tank is nil")
	}
	_, o := m.playerTankPlayables[uid]
	if !o {
		var playables *tankLevelDirPlayable
		freeList := m.freeTankPlayables[tank.Id()]
		if len(freeList) > 0 { // 从freelist中取一个
			l := len(freeList)
			playables = freeList[l-1]
			playables.resetTank(tank)
			freeList = freeList[:l-1]
			m.freeTankPlayables[tank.Id()] = freeList
		} else {
			playables = m.newTankLevelDirPlayable(tank)
		}
		m.playerTankPlayables[uid] = playables
	}
}

// 创建坦克可播放结构
func (m *PlayableManager) newTankLevelDirPlayable(tank *object.Tank) *tankLevelDirPlayable {
	playables := &tankLevelDirPlayable{}
	tid := tank.Id()
	maxLevel := common_data.TankMaxLevelMap[tid]
	playables.levelPlayables = make([]*tankDirPlayable, maxLevel)
	for i := int32(0); i < maxLevel; i++ {
		playables.levelPlayables[i] = &tankDirPlayable{}
		for j := int32(0); j <= int32(object.DirMax-object.DirMin); j++ {
			level := i + 1
			dir := object.Direction(j + int32(object.DirMin))
			anim := CreateTankAnim(tank.Id(), level, dir)
			if anim == nil {
				str := fmt.Sprintf("new tank playable anim is nil, level %v, dir %v", level, dir)
				panic(str)
			}
			switch dir {
			case object.DirUp:
				playables.levelPlayables[i].up = NewPlayableObject(tank, anim)
			case object.DirDown:
				playables.levelPlayables[i].down = NewPlayableObject(tank, anim)
			case object.DirLeft:
				playables.levelPlayables[i].left = NewPlayableObject(tank, anim)
			case object.DirRight:
				playables.levelPlayables[i].right = NewPlayableObject(tank, anim)
			}
		}
	}
	playables.tank = tank
	playables.tankId = tank.Id()
	return playables
}

// 删除玩家坦克播放（放入对应坦克id播放的freelist）
func (m *PlayableManager) RemovePlayerTankLevelDirPlayable(uid uint64) {
	ldPlayable, o := m.playerTankPlayables[uid]
	if !o {
		return
	}
	ldPlayable.tank = nil // 把坦克指针置为空
	delete(m.playerTankPlayables, uid)

	freeList := m.freeTankPlayables[ldPlayable.tankId]
	freeList = append(freeList, ldPlayable)
	m.freeTankPlayables[ldPlayable.tankId] = freeList
}

// 开始播放玩家坦克动画
func (m *PlayableManager) PlayPlayerTankPlayable(uid uint64) {
	playables, o := m.playerTankPlayables[uid]
	if o {
		m.playTankPlayable(playables)
	}
}

// 播放坦克动画
func (m *PlayableManager) playTankPlayable(playables *tankLevelDirPlayable) {
	playable := playables.getCurrPlayable()
	if playables.lastPlayable != playable {
		if playables.lastPlayable != nil {
			playables.lastPlayable.Stop()
		}
		playables.lastPlayable = playable
	}
	playable.Play()
}

// 停止播放玩家坦克动画
func (m *PlayableManager) StopPlayerTankPlayable(uid uint64) {
	playables, o := m.playerTankPlayables[uid]
	if o {
		playables.stop()
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
		playables := m.newTankLevelDirPlayable(tank)
		m.enemyTankPlayables[instId] = playables
	}
}

// 删除玩家坦克播放
func (m *PlayableManager) RemoveEnemyTankPlayable(instId int32) {
	delete(m.enemyTankPlayables, instId)
}

// 播放敌人坦克动画
func (m *PlayableManager) PlayEnemyTankPlayable(instId int32) {
	playables, o := m.enemyTankPlayables[instId]
	if o {
		m.playTankPlayable(playables)
	}
}

// 停止播放敌人坦克动画
func (m *PlayableManager) StopEnemyTankPlayable(id int32) {
	playables, o := m.enemyTankPlayables[id]
	if o {
		playables.stop()
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
		playable := p.getCurrPlayable()
		playable.Update(screen)
	}
}

// 更新敌人坦克动画
func (m *PlayableManager) UpdatEnemyTanksPlayable(screen *ebiten.Image) {
	for _, p := range m.enemyTankPlayables {
		playable := p.getCurrPlayable()
		playable.Update(screen)
	}
}
