package main

import (
	"reflect"

	"project_b/common/base"
	"project_b/common/object"
)

const (
	//EventIdNone = common.EventId(iota)

	/* UI操作事件 */
	EventIdOpLogin     = base.EventId(1)
	EventIdOpEnterGame = base.EventId(2)

	/* 网络协议事件 */
	// 进入游戏 参数：Account(string), PlayerId(uint64), TankInfo()
	EventIdPlayerEnterGame = base.EventId(100)
	// 进入游戏完成
	EventIdPlayerEnterGameCompleted = base.EventId(101)
	// 离开游戏
	EventIdPlayerExitGame = base.EventId(102)

	/* 游戏逻辑事件 */
	EventIdTankMove     = base.EventId(200)  // 移动事件
	EventIdTankStopMove = base.EventId(201)  // 停止移动事件
	EventIdTankMoveSync = base.EventId(202)  // 移动同步事件
	EventIdTankChange   = base.EventId(1000) // 改变坦克
	EventIdTankRestore  = base.EventId(1001) // 恢复坦克
)

// todo 发送协议的事件处理和设备输入的事件处理最好分开，方便做逻辑和显示分离
// 注册事件
func (g *Game) registerEvents() {
	// 请求登录
	g.eventMgr.RegisterEvent(EventIdOpLogin, g.onEventReqLogin)
	// 请求进入游戏
	g.eventMgr.RegisterEvent(EventIdOpEnterGame, g.onEventReqEnterGame)

	// 进入游戏
	g.eventMgr.RegisterEvent(EventIdPlayerEnterGame, g.onEventPlayerEnterGame)
	// 进入游戏完成
	g.eventMgr.RegisterEvent(EventIdPlayerEnterGameCompleted, g.onEventPlayerEnterGameCompleted)
	// 退出游戏
	g.eventMgr.RegisterEvent(EventIdPlayerExitGame, g.onEventPlayerExitGame)

	// 坦克移动
	g.eventMgr.RegisterEvent(EventIdTankMove, g.onEventTankMove)
	// 坦克停止移动
	g.eventMgr.RegisterEvent(EventIdTankStopMove, g.onEventTankStopMove)
	// 坦克改变
	g.eventMgr.RegisterEvent(EventIdTankChange, g.onEventTankChange)
	// 坦克恢复
	g.eventMgr.RegisterEvent(EventIdTankRestore, g.onEventTankRestore)
}

// 注销事件
func (g *Game) unregisterEvents() {
	g.eventMgr.UnregisterEvent(EventIdOpLogin, g.onEventReqLogin)
	g.eventMgr.UnregisterEvent(EventIdOpEnterGame, g.onEventReqEnterGame)
	g.eventMgr.UnregisterEvent(EventIdPlayerEnterGame, g.onEventPlayerEnterGame)
	g.eventMgr.UnregisterEvent(EventIdPlayerEnterGameCompleted, g.onEventPlayerEnterGameCompleted)
	g.eventMgr.UnregisterEvent(EventIdPlayerExitGame, g.onEventPlayerExitGame)
	g.eventMgr.UnregisterEvent(EventIdTankMove, g.onEventTankMove)
	g.eventMgr.UnregisterEvent(EventIdTankStopMove, g.onEventTankStopMove)
	g.eventMgr.UnregisterEvent(EventIdTankChange, g.onEventTankChange)
	g.eventMgr.UnregisterEvent(EventIdTankRestore, g.onEventTankRestore)
}

// 请求登录
func (g *Game) onEventReqLogin(args ...interface{}) {
	var account string
	var password string
	var o bool
	a := args[0]
	p := args[1]
	account, o = a.(string)
	if !o {
		t := reflect.TypeOf(a)
		getLog().Warn("account type must string on req login, this is %v", t)
		return
	}
	password, o = p.(string)
	if !o {
		getLog().Warn("password type must string on req login")
		return
	}
	err := g.net.SendLoginReq(account, password)
	if err != nil {
		getLog().Warn("send login req err: %v", err)
		return
	}
	g.myAcc = account
	getLog().Info("handle event: account %v password %v send login req", account, password)
}

// 请求进入游戏
func (g *Game) onEventReqEnterGame(args ...interface{}) {
	var account, sessionToken string
	var o bool
	account, o = args[0].(string)
	if !o {
		getLog().Warn("account type must string on req enter game")
		return
	}
	err := g.net.SendEnterGameReq(account, sessionToken)
	if err != nil {
		getLog().Warn("send enter game err: %v", err)
		return
	}
	getLog().Info("handle event: account %v send enter game req", account)
}

// 处理玩家进入游戏事件
func (g *Game) onEventPlayerEnterGame(args ...interface{}) {
	if len(args) < 3 {
		getLog().Warn("onEventEnterGame event args length cant less than 3")
		return
	}

	account := args[0].(string)
	uid := args[1].(uint64)
	tank := args[2].(*object.Tank)

	// 加入播放管理
	g.playableMgr.AddPlayerTankPlayable(uid, tank)

	if g.myAcc == account {
		g.myId = uid
		// 游戏状态
		g.mode = ModeGame
		getLog().Info("handle event: my player (account: %v, uid: %v) entered game, tank %v", account, uid, *tank)
	} else {
		getLog().Info("handle event: player (account: %v, uid: %v) entered game, tank %v", account, uid, *tank)
	}
}

// 处理进入游戏完成
func (g *Game) onEventPlayerEnterGameCompleted(args ...interface{}) {
	getLog().Info("handle event: my player (account: %v, uid: %v) enter game finished", g.myAcc, g.myId)
}

// 处理玩家离开游戏事件
func (g *Game) onEventPlayerExitGame(args ...interface{}) {
	if len(args) < 1 {
		getLog().Warn("onEventPlayerExitGame event args length cant less 1")
		return
	}

	uid := args[0].(uint64)

	// 从播放管理器中删除
	g.playableMgr.RemovePlayerTankPlayable(uid)

	getLog().Info("handle event: player (uid: %v) exited game", uid)
}

// 处理坦克移动事件
func (g *Game) onEventTankMove(args ...interface{}) {
	g.playableMgr.PlayPlayerTankPlayable(g.myId)
	//log.Printf("handle event: my player move tank")
}

// 处理坦克停止移动事件
func (g *Game) onEventTankStopMove(args ...interface{}) {
	g.playableMgr.StopPlayerTankPlayable(g.myId)
	//log.Printf("handle event: my player stop move tank")
}

// 处理坦克改变事件
func (g *Game) onEventTankChange(args ...interface{}) {
	if len(args) < 2 {
		gslog.Error("onEventTankChange event need 3 args")
		return
	}
	pid := args[0].(uint64)
	tank := args[1].(*object.Tank)
	g.playableMgr.RemovePlayerTankPlayable(pid)
	g.playableMgr.AddPlayerTankPlayable(pid, tank)
	getLog().Info("handle event: player %v changed tank to %v", pid, tank.Id())
}

// 处理坦克恢复事件
func (g *Game) onEventTankRestore(args ...interface{}) {
	if len(args) < 2 {
		gslog.Error("onEventTankRestore event need 3 args")
		return
	}
	pid := args[0].(uint64)
	tank := args[1].(*object.Tank)
	g.playableMgr.RemovePlayerTankPlayable(pid)
	g.playableMgr.AddPlayerTankPlayable(pid, tank)
	getLog().Info("handle event: player %v restore tank id to %v", pid, tank.Id())
}
