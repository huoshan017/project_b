package common

import (
	"math"
	"project_b/common/base"
	"project_b/common/log"
	"project_b/common/object"
	"project_b/common/time"
	"project_b/game_map"
)

const (
	defaultLogicFrameMax = math.MaxInt32
)

type GameLogic struct {
	eventMgr base.IEventManager // 事件管理
	//scene    *Scene             // 场景
	sceneMap *SceneMap // 場景圖
	state    int32     // 0 未开始  1. 运行中
	mapIndex int32     // 地图索引
	frame    int32     // 帧序号，每Update一次加1
	maxFrame int32     // 最大帧序号
}

// 创建游戏逻辑
func NewGameLogic(eventMgr base.IEventManager) *GameLogic {
	gl := &GameLogic{}
	if eventMgr == nil {
		eventMgr = base.NewEventManager()
	}
	gl.eventMgr = eventMgr
	//gl.scene = NewScene(gl.eventMgr)
	gl.sceneMap = NewSceneMap(gl.eventMgr)
	return gl
}

// 载入地图
func (g *GameLogic) LoadMap(mapId int32) bool {
	m := game_map.MapConfigArray[mapId]
	return g.sceneMap.LoadMap(m)
}

// 載入場景地圖
func (g *GameLogic) LoadSceneMap(mapId int32) bool {
	return g.sceneMap.LoadMap(game_map.MapConfigArray[mapId])
}

// 卸載場景圖
func (g *GameLogic) UnloadSceneMap() {
	g.sceneMap.UnloadMap()
}

// 場景圖
func (g *GameLogic) CurrentSceneMap() *SceneMap {
	return g.sceneMap
}

// 地图索引
func (g *GameLogic) MapIndex() int32 {
	return g.mapIndex
}

// 设置最大帧序号
func (g *GameLogic) SetMaxFrame(maxFrame int32) {
	g.maxFrame = maxFrame
}

// 当前帧
func (g *GameLogic) GetCurrFrame() int32 {
	return g.frame
}

// 在逻辑线程中更新
func (g *GameLogic) Update(tick time.Duration) {
	g.sceneMap.Update(tick)
	g.frame += 1
	if g.maxFrame > 0 {
		if g.frame >= g.maxFrame {
			g.frame = 1
		}
	} else if g.frame >= defaultLogicFrameMax {
		g.frame = 1
	}
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

// 注册场景事件
func (g *GameLogic) RegisterSceneEvent(eid base.EventId, handle func(args ...interface{})) {
	g.sceneMap.RegisterEvent(eid, handle)
}

// 注销场景事件
func (g *GameLogic) UnregisterSceneEvent(eid base.EventId, handle func(args ...interface{})) {
	g.sceneMap.UnregisterEvent(eid, handle)
}

// 注册坦克事件
func (g *GameLogic) RegisterPlayerSceneEvent(uid uint64, eid base.EventId, handle func(args ...interface{})) {
	g.sceneMap.RegisterPlayerEvent(uid, eid, handle)
}

// 注销坦克事件
func (g *GameLogic) UnregisterPlayerSceneEvent(uid uint64, eid base.EventId, handle func(args ...interface{})) {
	g.sceneMap.UnregisterPlayerEvent(uid, eid, handle)
}

// 获得玩家坦克
func (g *GameLogic) GetPlayerTank(uid uint64) *object.Tank {
	return g.sceneMap.GetPlayerTank(uid)
}

// 获得玩家坦克列表
func (g *GameLogic) GetPlayerTankList() []PlayerTankKV {
	return g.sceneMap.GetPlayerTankList()
}

// 玩家进入
func (g *GameLogic) NewPlayerEnter(pid uint64) *object.Tank {
	return g.sceneMap.NewPlayerTank(pid)
}

// 玩家坦克进入
func (g *GameLogic) PlayerEnterWithTank(uid uint64, tank *object.Tank) {
	if g.sceneMap.GetPlayerTank(uid) == nil {
		g.sceneMap.AddPlayerTank(uid, tank)
	}
}

// 玩家离开
func (g *GameLogic) PlayerLeave(pid uint64) {
	g.sceneMap.RemovePlayerTank(pid)
}

// 获得敌人坦克
func (g *GameLogic) GetEnemyTank(instId uint32) *object.Tank {
	return g.sceneMap.GetEnemyTank(instId)
}

// 玩家坦克移动
func (g *GameLogic) PlayerTankMove(uid uint64, moveDir object.Direction) {
	g.sceneMap.PlayerTankMove(uid, moveDir)
}

// 玩家坦克停止移动
func (g *GameLogic) PlayerTankStopMove(uid uint64) {
	g.sceneMap.PlayerTankStopMove(uid)
}

// 玩家坦克改变
func (g *GameLogic) PlayerTankChange(uid uint64, staticInfo *object.ObjStaticInfo) bool {
	return g.sceneMap.PlayerTankChange(uid, staticInfo)
}

// 玩家坦克恢复
func (g *GameLogic) PlayerTankRestore(uid uint64) int32 {
	return g.sceneMap.PlayerTankRestore(uid)
}

// 检测玩家
func (g *GameLogic) CheckPlayerTankStartMove(uid uint64, startPos object.Pos, dir object.Direction, speed int32) bool {
	tank := g.GetPlayerTank(uid)
	if tank == nil {
		log.Error("Cant get tank from uid %v", uid)
		return false
	}
	x, y := tank.Pos()
	// 坐标是不是一致
	if startPos.X != x || startPos.Y != y {
		return false
	}
	// 速度是不是合法
	if speed != tank.CurrentSpeed() {
		return false
	}
	g.PlayerTankMove(uid, dir)
	return true
}
