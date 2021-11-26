package common

import (
	"project_b/common/base"
	"project_b/common/object"
	"project_b/game_map"
	"time"
)

type GameLogic struct {
	eventMgr base.IEventManager // 事件管理
	scene    *Scene             // 场景
	state    int32              // 0 未开始  1. 运行中
	mapIndex int32              // 地图索引
}

// 创建游戏逻辑
func NewGameLogic() *GameLogic {
	gl := &GameLogic{eventMgr: base.NewEventManager()}
	gl.scene = NewScene(gl.eventMgr)
	return gl
}

// 载入地图
func (g *GameLogic) LoadMap(mapIndex int32) {
	m := game_map.MapConfigArray[mapIndex]
	g.scene.LoadMap(&m)
}

// 地图索引
func (g *GameLogic) MapIndex() int32 {
	return g.mapIndex
}

// 在逻辑线程中更新
func (g *GameLogic) Update(tick time.Duration) {
	g.scene.Update(tick)
}

// 开始逻辑
func (g *GameLogic) Start() {
	g.state = 1
}

// 结束逻辑
func (g *GameLogic) End() {
	g.state = 0
}

// 是否已开始
func (g *GameLogic) IsStart() bool {
	return g.state == 1
}

// 注册事件
func (g *GameLogic) RegisterEvent(eid base.EventId, handle func(args ...interface{})) {
	g.eventMgr.RegisterEvent(eid, handle)
}

// 注销事件
func (g *GameLogic) UnregisterEvent(eid base.EventId, handle func(args ...interface{})) {
	g.eventMgr.UnregisterEvent(eid, handle)
}

// 获得玩家坦克
func (g *GameLogic) GetPlayerTank(uid uint64) *object.Tank {
	return g.scene.GetPlayerTank(uid)
}

// 获得所有玩家坦克列表
func (g *GameLogic) GetPlayerTankList() []*object.Tank {
	return g.scene.GetPlayerTankList()
}

// 获得所有玩家坦克
func (g *GameLogic) GetPlayerTanks() map[uint64]*object.Tank {
	return g.scene.GetPlayerTanks()
}

// 玩家坦克进入
func (g *GameLogic) PlayerTankEnter(uid uint64, tank *object.Tank) {
	if g.scene.GetPlayerTank(uid) == nil {
		g.scene.AddPlayerTank(uid, tank)
	}
}

// 玩家坦克离开
func (g *GameLogic) PlayerTankLeave(uid uint64) {
	g.scene.RemovePlayerTank(uid)
}

// 获得敌人坦克
func (g *GameLogic) GetEnemyTank(id int32) *object.Tank {
	return g.scene.GetEnemyTank(id)
}

// 获得所有敌人坦克列表
func (g *GameLogic) GetEnemyTankList() []*object.Tank {
	return g.scene.GetEnemyTankList()
}

// 获得所有敌人坦克
func (g *GameLogic) GetEnemyTanks() map[int32]*object.Tank {
	return g.scene.GetEnemyTanks()
}

// 玩家坦克移动
func (g *GameLogic) PlayerTankMove(uid uint64, moveDir object.Direction) {
	g.scene.PlayerTankMove(uid, moveDir)
}

// 玩家坦克停止移动
func (g *GameLogic) PlayerTankStopMove(uid uint64) {
	g.scene.PlayerTankStopMove(uid)
}

// 玩家坦克改变
func (g *GameLogic) PlayerTankChange(uid uint64, staticInfo *object.ObjStaticInfo) bool {
	return g.scene.PlayerTankChange(uid, staticInfo)
}

// 玩家坦克恢复
func (g *GameLogic) PlayerTankRestore(uid uint64) int32 {
	return g.scene.PlayerTankRestore(uid)
}
