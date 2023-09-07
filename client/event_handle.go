package main

import (
	"reflect"

	"project_b/client_base"
	"project_b/client_core"
	"project_b/common"
	"project_b/common/base"
	"project_b/common/object"
	"project_b/core"
	"project_b/game_proto"
	"project_b/log"
)

type event2Handle struct {
	eid    base.EventId
	handle func(args ...any)
}

type EventHandles struct {
	net           *client_core.NetClient
	gameCore      *core.GameCore
	playableScene *PlayableScene
	gameData      *client_base.GameData
	// --------------------------------------
	// 事件处理
	gameEvent2Handles   []event2Handle // 游戏事件
	playerEvent2Handles []event2Handle // 玩家事件
}

// 创建EventHandles
func CreateEventHandles(net *client_core.NetClient, gameCore *core.GameCore, playableScene *PlayableScene, gameData *client_base.GameData) *EventHandles {
	eh := &EventHandles{
		net:           net,
		gameCore:      gameCore,
		playableScene: playableScene,
		gameData:      gameData,
	}
	return eh
}

// 初始化
func (g *EventHandles) Init() {
	g.registerEvents()
}

// 反初始化
func (g *EventHandles) Uninit() {
	g.unregisterEvents()
}

// 注册事件
// todo 发送协议的事件处理和设备输入的事件处理最好分开，方便做逻辑和显示分离
func (g *EventHandles) registerEvents() {
	// 玩家事件处理
	g.playerEvent2Handles = []event2Handle{
		{common.EventIdBeforeMapLoad, g.onEventBeforeMapLoad},                            // 地圖載入前
		{common.EventIdMapLoaded, g.onEventMapLoaded},                                    // 地圖載入完成
		{common.EventIdBeforeMapUnload, g.onEventBeforeMapUnload},                        // 地圖卸載前
		{common.EventIdMapUnloaded, g.onEventMapUnloaded},                                // 地圖卸載後
		{client_core.EventIdOpLogin, g.onEventReqLogin},                                  // 请求登录
		{client_core.EventIdOpEnterGame, g.onEventReqEnterGame},                          // 请求进入游戏
		{client_core.EventIdTimeSync, g.onEventTimeSync},                                 // 同步时间
		{client_core.EventIdTimeSyncEnd, g.onEventTimeSyncEnd},                           // 同步时间结束
		{client_core.EventIdPlayerEnterGame, g.onEventPlayerEnterGame},                   // 进入游戏
		{client_core.EventIdPlayerEnterGameCompleted, g.onEventPlayerEnterGameCompleted}, // 进入游戏完成
		{client_core.EventIdPlayerExitGame, g.onEventPlayerExitGame},                     // 退出游戏
	}

	// 游戏事件处理
	g.gameEvent2Handles = []event2Handle{
		{common.EventIdTankMove, g.onEventTankMove},         // 处理坦克移动
		{common.EventIdTankStopMove, g.onEventTankStopMove}, // 处理坦克停止移动
		{common.EventIdTankSetPos, g.onEventTankSetPos},     // 处理坦克位置更新
		{common.EventIdTankChange, g.onEventTankChange},     // 处理坦克变化
		{common.EventIdTankRestore, g.onEventTankRestore},   // 处理坦克恢复
	}

	for _, e2h := range g.playerEvent2Handles {
		g.gameCore.RegisterEvent(e2h.eid, e2h.handle)
	}
}

// 注销事件
func (g *EventHandles) unregisterEvents() {
	for _, e2h := range g.playerEvent2Handles {
		g.gameCore.UnregisterEvent(e2h.eid, e2h.handle)
	}
}

/*
*
请求登录
args[0]: account string
args[1]: password string
*/
func (g *EventHandles) onEventReqLogin(args ...any) {
	var account string
	var password string
	var o bool
	a := args[0]
	p := args[1]
	account, o = a.(string)
	if !o {
		t := reflect.TypeOf(a)
		log.Warn("account type must string on req login, this is %v", t)
		return
	}
	password, o = p.(string)
	if !o {
		log.Warn("password type must string on req login")
		return
	}
	err := g.net.SendLoginReq(account, password)
	if err != nil {
		log.Warn("send login req err: %v", err)
		return
	}
	g.gameData.MyAcc = account
	log.Info("handle event: account %v password %v send login req", account, password)
}

/*
*
请求进入游戏
args[0]: account string
*/
func (g *EventHandles) onEventReqEnterGame(args ...any) {
	var account, sessionToken string
	var o bool
	account, o = args[0].(string)
	if !o {
		log.Warn("account type must string on req enter game")
		return
	}
	err := g.net.SendEnterGameReq(account, sessionToken)
	if err != nil {
		log.Warn("send enter game err: %v", err)
		return
	}
	log.Info("handle event: account %v send enter game req", account)
}

/*
*
处理玩家进入游戏事件
args[0]: account(string)
args[1]: uid(uint64)
*/
func (g *EventHandles) onEventPlayerEnterGame(args ...any) {
	if len(args) < 2 {
		log.Warn("onEventEnterGame event args length cant less than 3")
		return
	}

	account := args[0].(string)
	uid := args[1].(uint64)

	if g.gameData.MyAcc == account {
		g.gameData.MyId = uid
		// 游戏状态
		g.gameData.State = client_base.GameStateInWorld
		log.Info("handle event: my player (account: %v, uid: %v) entered game", account, uid)
	} else {
		log.Info("handle event: player (account: %v, uid: %v) entered game", account, uid)
	}
}

/*
*
处理进入游戏完成
*/
func (g *EventHandles) onEventPlayerEnterGameCompleted(args ...any) {
	// 准备同步服务器时间
	if err := g.net.SendTimeSyncReq(); err != nil {
		log.Error("handle event: send time sync request err: %v", err)
		return
	}

	// 注册本游戏场景事件
	for _, e2h := range g.gameEvent2Handles {
		g.gameCore.RegisterPlayerEvent(g.gameData.MyId, e2h.eid, e2h.handle)
	}

	log.Info("handle event: my player (account: %v, uid: %v) enter game finished", g.gameData.MyAcc, g.gameData.MyId)
}

/*
*
处理玩家离开游戏事件
args[0]: uid(uint64)
*/
func (g *EventHandles) onEventPlayerExitGame(args ...any) {
	// 注销本游戏场景事件
	for _, e2h := range g.gameEvent2Handles {
		g.gameCore.UnregisterPlayerEvent(g.gameData.MyId, e2h.eid, e2h.handle)
	}

	g.gameData.State = client_base.GameStateMainMenu

	log.Info("handle event: player exited game")
}

/*
*
处理时间同步事件
*/
func (g *EventHandles) onEventTimeSync(args ...any) {
	if err := g.net.SendTimeSyncReq(); err != nil {
		log.Error("handle event: send time sync request err: %v", err)
		return
	}
}

/**
 *处理时间同步结束事件
 */
func (g *EventHandles) onEventTimeSyncEnd(args ...any) {
	log.Info("handle event: time sync end")
}

/**
 * 处理载入地图前事件
 */
func (g *EventHandles) onEventBeforeMapLoad(args ...any) {

}

/**
 * 处理地图载入完成事件
 */
func (eh *EventHandles) onEventMapLoaded(args ...any) {
	currentScene := args[0].(*common.World)
	eh.playableScene.SetScene(currentScene)
}

/**
 * 处理地图卸载前事件
 */
func (eh *EventHandles) onEventBeforeMapUnload(args ...any) {

}

/**
 * 处理地图卸载后事件
 */
func (eh *EventHandles) onEventMapUnloaded(args ...any) {
	eh.playableScene.UnloadScene()
}

/*
*
处理坦克移动事件
args[0]: pos(object.Pos)
args[1]: direction(object.Direction)
args[2]: speed(int32)
*/
func (eh *EventHandles) onEventTankMove(args ...any) {
	pos := args[0].(base.Pos)
	orientation := args[1].(base.Angle)
	speed := args[2].(int32)
	err := eh.net.SendTankUpdatePosReq(game_proto.MovementState_StartMove, pos, int32(orientation.ToMinutes()) /*dir*/, speed)
	if err != nil {
		log.Error("send tank move req err: %v", err)
	}
}

/*
*
处理坦克停止移动事件
args[0]: object.Pos
args[1]: object.Direction
args[2]: int32
*/
func (eh *EventHandles) onEventTankStopMove(args ...any) {
	pos := args[0].(base.Pos)
	orientation := args[1].(base.Angle)
	speed := args[2].(int32)
	err := eh.net.SendTankUpdatePosReq(game_proto.MovementState_ToStop, pos /*dir*/, int32(orientation.ToMinutes()), speed)
	if err != nil {
		log.Error("send tank stop move req err: %v", err)
	}
}

/*
*
处理坦克设置坐标事件
args[0]: object.Pos
args[1]: object.Direction
args[2]: int32
*/
func (eh *EventHandles) onEventTankSetPos(args ...any) {
	pos := args[0].(base.Pos)
	orientation := args[1].(base.Angle)
	speed := args[2].(int32)
	err := eh.net.SendTankUpdatePosReq(game_proto.MovementState_Moving, pos /*dir*/, int32(orientation.ToMinutes()), speed)
	if err != nil {
		log.Error("send tank update pos req err: %v", err)
	}
}

/*
*
处理坦克改变事件
args[0]: uint64
args[1]: *object.Tank
*/
func (eh *EventHandles) onEventTankChange(args ...any) {
	if len(args) < 2 {
		log.Error("onEventTankChange event need 3 args")
		return
	}
	pid := args[0].(uint64)
	tank := args[1].(*object.Tank)
	log.Info("handle event: player %v changed tank to %v", pid, tank.Id())
}

/*
**
处理坦克恢复事件
args[0]: uint64
args[1]: *object.Tank
*/
func (eh *EventHandles) onEventTankRestore(args ...any) {
	if len(args) < 2 {
		log.Error("onEventTankRestore event need 2 args")
		return
	}
	pid := args[0].(uint64)
	tank := args[1].(*object.Tank)
	log.Info("handle event: player %v restore tank id to %v", pid, tank.Id())
}
