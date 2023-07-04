package common

import (
	"math"
	"project_b/common/base"
	"project_b/common/ds"
	"project_b/common/log"
	"project_b/common/object"
	"project_b/common/time"
	"project_b/common_data"
	"project_b/game_map"
)

const (
	defaultLogicFrameMax = math.MaxInt32
)

// 基於SceneLogic增加了玩家(Player)概念的游戲邏輯
type GameLogic struct {
	eventMgr    base.IEventManager               // 事件管理
	state       int32                            // 0 未开始  1. 运行中
	mapIndex    int32                            // 地图索引
	frame       int32                            // 帧序号，每Update一次加1
	maxFrame    int32                            // 最大帧序号
	scene       *SceneLogic                      // 場景圖
	player2Tank *ds.MapListUnion[uint64, uint32] // 玩家與坦克之間對應關係
	tank2Player *ds.MapListUnion[uint32, uint64] // 坦克到玩家的對應關係
	botMgr      *BotManager                      // bot管理器
	tank2Bot    *ds.MapListUnion[uint32, int32]  // 坦克到bot的對應關係
}

// 创建游戏逻辑
func NewGameLogic(eventMgr base.IEventManager) *GameLogic {
	gl := &GameLogic{}
	if eventMgr == nil {
		eventMgr = base.NewEventManager()
	}
	gl.eventMgr = eventMgr
	gl.scene = NewSceneLogic(gl.eventMgr)
	gl.player2Tank = ds.NewMapListUnion[uint64, uint32]()
	gl.tank2Player = ds.NewMapListUnion[uint32, uint64]()
	gl.botMgr = NewBotManager()
	gl.tank2Bot = ds.NewMapListUnion[uint32, int32]()
	return gl
}

// 載入場景地圖
func (g *GameLogic) LoadScene(config *game_map.Config) bool {
	g.eventMgr.InvokeEvent(EventIdBeforeMapLoad)
	loaded := g.scene.LoadMap(config)
	if !loaded {
		log.Error("GameLogic.LoadScene mapid(%v) failed", config.Id)
		return false
	}
	g.createBots(config)
	g.scene.RegisterTankAddedHandle(g.onTankCreated)
	g.scene.RegisterTankRemovedHandle(g.onTankDestroyed)
	g.eventMgr.InvokeEvent(EventIdMapLoaded, g.scene)
	return loaded
}

// 卸載場景圖
func (g *GameLogic) UnloadScene() {
	g.eventMgr.InvokeEvent(EventIdBeforeMapUnload)
	g.scene.UnregisterTankAddedHandle(g.onTankCreated)
	g.scene.UnregisterTankRemovedHandle(g.onTankDestroyed)
	g.scene.UnloadMap()
	g.player2Tank.Clear()
	g.tank2Player.Clear()
	g.botMgr.Clear()
	g.tank2Bot.Clear()
	g.eventMgr.InvokeEvent(EventIdMapUnloaded)
}

// 場景圖
func (g *GameLogic) CurrentScene() *SceneLogic {
	return g.scene
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
	g.botMgr.Update(tick)
	g.scene.Update(tick)
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
	g.scene.RegisterEvent(eid, handle)
}

// 注销场景事件
func (g *GameLogic) UnregisterSceneEvent(eid base.EventId, handle func(args ...interface{})) {
	g.scene.UnregisterEvent(eid, handle)
}

// 注册坦克事件
func (g *GameLogic) RegisterPlayerSceneEvent(uid uint64, eid base.EventId, handle func(args ...any)) {
	tankId, o := g.player2Tank.Get(uid)
	if !o {
		return
	}
	g.scene.RegisterTankEvent(tankId, eid, handle)
}

// 注销坦克事件
func (g *GameLogic) UnregisterPlayerSceneEvent(uid uint64, eid base.EventId, handle func(args ...interface{})) {
	tankId, o := g.player2Tank.Get(uid)
	if !o {
		return
	}
	g.scene.UnregisterTankEvent(tankId, eid, handle)
}

// 获得玩家坦克
func (g *GameLogic) GetPlayerTank(uid uint64) *object.Tank {
	tankId, o := g.player2Tank.Get(uid)
	if !o {
		return nil
	}
	return g.scene.GetTank(tankId)
}

// 获得玩家坦克列表
func (g *GameLogic) GetPlayerTankList() []PlayerTankKV {
	var kvs []PlayerTankKV
	lis := g.player2Tank.GetList()
	for _, v := range lis {
		tankId, o := g.player2Tank.Get(v.Key)
		if !o {
			continue
		}
		tank := g.scene.GetTank(tankId)
		if tank == nil {
			continue
		}
		kvs = append(kvs, PlayerTankKV{v.Key, tank})
	}
	return kvs
}

// 新玩家进入
func (g *GameLogic) NewPlayerEnterWithPos(pid uint64, x, y int32) *object.Tank {
	tank := g.scene.NewTankWithPos(x, y)
	if tank == nil {
		log.Error("player %v enter with pos (%v, %v) to tank failed", pid, x, y)
		return nil
	}
	g.player2Tank.Add(pid, tank.InstId())
	g.tank2Player.Add(tank.InstId(), pid)
	return tank
}

// 玩家進入
func (g *GameLogic) PlayerEnterWithStaticInfo(pid uint64, id int32, level int32, x, y int32, dir object.Direction, currentSpeed int32) uint32 {
	tank := g.scene.NewTankWithStaticInfo(id, level, x, y, dir, currentSpeed)
	if tank == nil {
		log.Error("player %v enter with static info to create tank failed", pid)
		return 0
	}
	g.player2Tank.Add(pid, tank.InstId())
	g.tank2Player.Add(tank.InstId(), pid)
	return tank.InstId()
}

// 玩家坦克进入
func (g *GameLogic) PlayerEnterWithTank(uid uint64, tank *object.Tank) {
	tankId, o := g.player2Tank.Get(uid)
	if o {
		log.Warn("player %v already entered with tank %v", tankId)
		return
	}
	g.scene.AddTank(tank)
	g.player2Tank.Add(uid, tank.InstId())
	g.tank2Player.Add(tank.InstId(), uid)
}

// 玩家离开
func (g *GameLogic) PlayerLeave(pid uint64) {
	tankId, o := g.player2Tank.Get(pid)
	if !o {
		return
	}
	g.scene.RemoveTank(tankId)
	g.player2Tank.Remove(pid)
	g.tank2Player.Remove(tankId)
}

// 玩家坦克移动
func (g *GameLogic) PlayerTankMove(uid uint64, moveDir object.Direction) {
	tankId, o := g.player2Tank.Get(uid)
	if !o {
		return
	}
	g.scene.TankMove(tankId, moveDir)
}

// 玩家坦克停止
func (g *GameLogic) PlayerTankStopMove(uid uint64) {
	tankId, o := g.player2Tank.Get(uid)
	if !o {
		return
	}
	g.scene.TankStopMove(tankId)
}

// 玩家坦克改变
func (g *GameLogic) PlayerTankChange(uid uint64, staticInfo *object.TankStaticInfo) bool {
	tankId, o := g.player2Tank.Get(uid)
	if !o {
		return false
	}
	return g.scene.TankChange(tankId, staticInfo)
}

// 玩家坦克恢复
func (g *GameLogic) PlayerTankRestore(uid uint64) int32 {
	tankId, o := g.player2Tank.Get(uid)
	if !o {
		return 0
	}
	return g.scene.TankRestore(tankId)
}

// 玩家坦克開炮
func (g *GameLogic) PlayerTankFire(uid uint64) {
	tankId, o := g.player2Tank.Get(uid)
	if !o {
		return
	}
	g.scene.TankFire(tankId)
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

// 創建bot列表
func (g *GameLogic) createBots(config *game_map.Config) {
	for _, b := range config.BotInfoList {
		staticInfo := common_data.TankConfigData[b.TankId]
		if staticInfo == nil {
			log.Error("GameLogic.createBots tank config not found by id(%v)", b.TankId)
			continue
		}
		// todo 等級從1開始
		tank := g.scene.NewTankWithStaticInfo(staticInfo.Id(), 1, b.Pos.X, b.Pos.Y, staticInfo.Dir(), staticInfo.Speed())
		tank.SetCamp(b.Camp)
		bot := g.botMgr.NewBot(g.scene, tank.InstId())
		if bot == nil {
			log.Error("GameLogic.createBots NewBot failed")
			continue
		}
		g.tank2Bot.Add(tank.InstId(), bot.id)
	}
}

// 坦克創建事件
func (g *GameLogic) onTankCreated(args ...any) {

}

// 坦克被擊毀事件函數
func (g *GameLogic) onTankDestroyed(args ...any) {
	tank := args[0].(*object.Tank)
	instId := tank.InstId()
	pid, o := g.tank2Player.Get(instId)
	if !o {
		botId, o := g.tank2Bot.Get(instId)
		if !o {
			log.Warn("GameLogic.onTankDestroyed cant found tank by id %v", instId)
			return
		}
		g.tank2Bot.Remove(instId)
		g.botMgr.RemoveBot(botId)
	} else {
		g.tank2Player.Remove(instId)
		g.player2Tank.Remove(pid)
	}
	// bot中處理坦克被擊毀
	g.botMgr.onEnemyTankDestoryed(instId)
}
