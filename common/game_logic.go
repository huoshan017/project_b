package common

import (
	"math"
	"project_b/common/base"
	"project_b/common/ds"
	"project_b/common/object"
	"project_b/common/time"
	"project_b/common_data"
	"project_b/game_map"
	"project_b/log"
)

type TankType int32

const (
	TankTypePlayer TankType = iota
	TankTypeBot    TankType = 1
)

const (
	defaultLogicFrameMax = math.MaxUint32
)

type logicStateType int32

const (
	logicStateNotStart logicStateType = iota
	logicStateRunning
	logicStatePause
)

// 基於SceneLogic增加了玩家(Player)概念的游戲邏輯
type GameLogic struct {
	eventMgr    base.IEventManager               // 事件管理
	state       logicStateType                   // 0 未开始  1. 运行中
	mapIndex    int32                            // 地图索引
	frame       uint32                           // 帧序号，每Update一次加1
	maxFrame    uint32                           // 最大帧序号
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
	g.scene.RegisterObjectAddedHandle(g.onTankCreated)
	g.scene.RegisterObjectRemovedHandle(g.onTankDestroyed)
	g.eventMgr.InvokeEvent(EventIdMapLoaded, g.scene)
	return loaded
}

// 卸載場景圖
func (g *GameLogic) UnloadScene() {
	g.eventMgr.InvokeEvent(EventIdBeforeMapUnload)
	g.scene.UnregisterObjectAddedHandle(g.onTankCreated)
	g.scene.UnregisterObjectRemovedHandle(g.onTankDestroyed)
	g.scene.UnloadMap()
	g.player2Tank.Clear()
	g.tank2Player.Clear()
	g.botMgr.Clear()
	g.tank2Bot.Clear()
	g.eventMgr.InvokeEvent(EventIdMapUnloaded)
}

// 重載場景
func (g *GameLogic) ReloadScene() {
	g.clearBots()
	g.player2Tank.Clear()
	g.tank2Player.Clear()
	g.scene.ReloadMap()
	g.createBots(g.CurrentScene().mapConfig)
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
func (g *GameLogic) SetMaxFrame(maxFrame uint32) {
	g.maxFrame = maxFrame
}

// 当前帧
func (g *GameLogic) GetCurrFrame() uint32 {
	return g.frame
}

// 在逻辑线程中更新
func (g *GameLogic) Update(tick time.Duration) {
	if g.state == logicStateNotStart {
		g.state = logicStateRunning
	}
	if g.state != logicStateRunning {
		return
	}
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

// 暫停
func (g *GameLogic) Pause() {
	g.state = logicStatePause
	g.scene.Pause()
	g.botMgr.Pause()
}

// 繼續
func (g *GameLogic) Resume() {
	g.state = logicStateRunning
	g.scene.Resume()
	g.botMgr.Resume()
}

// 注册事件
func (g *GameLogic) RegisterEvent(eid base.EventId, handle func(args ...interface{})) {
	g.eventMgr.RegisterEvent(eid, handle)
}

// 注销事件
func (g *GameLogic) UnregisterEvent(eid base.EventId, handle func(args ...interface{})) {
	g.eventMgr.UnregisterEvent(eid, handle)
}

// 注册玩家事件
func (g *GameLogic) RegisterPlayerEvent(uid uint64, eid base.EventId, handle func(args ...any)) {
	tankId, o := g.player2Tank.Get(uid)
	if !o {
		return
	}
	g.scene.RegisterTankEvent(tankId, eid, handle)
}

// 注销玩家事件
func (g *GameLogic) UnregisterPlayerEvent(uid uint64, eid base.EventId, handle func(args ...interface{})) {
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
func (g *GameLogic) PlayerEnterWithStaticInfo(pid uint64, id int32, level int32, x, y int32, orientation int32 /*, currentSpeed int32*/) uint32 {
	tank := g.scene.NewTankWithStaticInfo(id, level, x, y /*, currentSpeed*/)
	if tank == nil {
		log.Error("player %v enter with static info to create tank failed", pid)
		return 0
	}
	angle := base.NewAngle(int16(orientation), 0)
	tank.RotateTo(angle)
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
func (g *GameLogic) PlayerTankMove(uid uint64 /*moveDir object.Direction*/, orientation int32) {
	if g.state == logicStatePause {
		return
	}
	tankId, o := g.player2Tank.Get(uid)
	if !o {
		return
	}
	g.scene.TankMove(tankId, orientation)
}

// 玩家坦克停止
func (g *GameLogic) PlayerTankStopMove(uid uint64) {
	if g.state == logicStatePause {
		return
	}
	tankId, o := g.player2Tank.Get(uid)
	if !o {
		return
	}
	g.scene.TankStopMove(tankId)
}

// 玩家坦克改变
func (g *GameLogic) PlayerTankChange(uid uint64, staticInfo *object.TankStaticInfo) bool {
	if g.state == logicStatePause {
		return false
	}
	tankId, o := g.player2Tank.Get(uid)
	if !o {
		return false
	}
	return g.scene.TankChange(tankId, staticInfo)
}

// 玩家坦克恢复
func (g *GameLogic) PlayerTankRestore(uid uint64) int32 {
	if g.state == logicStatePause {
		return 0
	}
	tankId, o := g.player2Tank.Get(uid)
	if !o {
		return 0
	}
	return g.scene.TankRestore(tankId)
}

// 玩家坦克開炮
func (g *GameLogic) PlayerTankFire(uid uint64) {
	if g.state == logicStatePause {
		return
	}
	tankId, o := g.player2Tank.Get(uid)
	if !o {
		return
	}
	g.scene.TankFire(tankId)
}

// 坦克增加新炮彈
func (g *GameLogic) PlayerTankAddNewShell(uid uint64, shellConfigId int32) {
	if g.state == logicStatePause {
		return
	}
	tankId, o := g.player2Tank.Get(uid)
	if o {
		g.scene.TankAddNewShell(tankId, shellConfigId)
	}
}

// 坦克切換炮彈
func (g *GameLogic) PlayerTankSwitchShell(uid uint64) {
	if g.state == logicStatePause {
		return
	}
	tankId, o := g.player2Tank.Get(uid)
	if o {
		g.scene.TankSwitchShell(tankId)
	}
}

// 坦克釋放環繞物體
func (g *GameLogic) PlayerTankReleaseSurroundObj(uid uint64) {
	if g.state == logicStatePause {
		return
	}
	tankId, o := g.player2Tank.Get(uid)
	if !o {
		return
	}
	g.scene.TankReleaseSurroundObj(tankId)
}

// 坦克旋轉
func (g *GameLogic) PlayerTankRotate(uid uint64, angle int32) {
	if g.state == logicStatePause {
		return
	}
	tankId, o := g.player2Tank.Get(uid)
	if !o {
		return
	}
	g.scene.TankRotate(tankId, angle)
}

// 坦克復活
func (g *GameLogic) PlayerTankRespawn(uid uint64, tankId int32, tankLevel int32, x, y int32, orientation int32 /*, currSpeed int32*/) {
	if g.state == logicStatePause {
		return
	}
	g.PlayerEnterWithStaticInfo(uid, tankId, tankLevel, x, y, orientation /*, currSpeed*/)
}

// 坦克護盾
func (g *GameLogic) PlayerTankShield(uid uint64) {
	if g.state == logicStatePause {
		return
	}
	tankId, o := g.player2Tank.Get(uid)
	if !o {
		return
	}
	g.scene.TankUnlimitedShield(tankId)
}

// 創建bot列表
func (g *GameLogic) createBots(config *game_map.Config) {
	var (
		tankBornPosList = g.scene.GetTankBornPosList()
		index           int32
	)
	for _, b := range config.BotInfoList {
		staticInfo := common_data.TankConfigData[b.TankId]
		if staticInfo == nil {
			log.Error("GameLogic.createBots tank config not found by id(%v)", b.TankId)
			continue
		}

		for index < int32(len(tankBornPosList)) && tankBornPosList[index].flag != object.BotTileFlag {
			index++
		}

		if index >= int32(len(tankBornPosList)) {
			break
		}

		x := tankBornPosList[index].x
		y := tankBornPosList[index].y
		// todo 等級從1開始
		tank := g.scene.NewTankWithStaticInfo(staticInfo.Id(), 1, x, y /*, staticInfo.Speed()*/)
		tank.SetCamp(b.Camp)
		tank.SetLevel(b.Level)
		bot := g.botMgr.NewBot(g.scene, tank.InstId())
		if bot == nil {
			log.Error("GameLogic.createBots NewBot failed")
			continue
		}
		g.tank2Bot.Add(tank.InstId(), bot.id)

		index += 1
	}
}

func (g *GameLogic) clearBots() {
	for i := int32(0); i < g.tank2Bot.Count(); i++ {
		k, v := g.tank2Bot.GetByIndex(i)
		g.botMgr.RemoveBot(v)
		g.scene.RemoveTank(k)
	}
	g.tank2Bot.Clear()
	g.botMgr.Clear()
}

// 坦克創建事件
func (g *GameLogic) onTankCreated(args ...any) {

}

// 坦克被擊毀事件函數
func (g *GameLogic) onTankDestroyed(args ...any) {
	tank, o := args[0].(*object.Tank)
	if !o {
		return
	}
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
		g.eventMgr.InvokeEvent(EventIdTankDestroy, TankTypeBot, botId)
	} else {
		g.tank2Player.Remove(instId)
		g.player2Tank.Remove(pid)
		g.eventMgr.InvokeEvent(EventIdTankDestroy, TankTypePlayer, pid)
	}
	// bot中處理坦克被擊毀
	g.botMgr.onEnemyTankDestoryed(instId)

}
