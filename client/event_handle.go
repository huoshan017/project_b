package main

import (
	"reflect"

	"project_b/client/core"
	"project_b/common/object"
)

// todo 发送协议的事件处理和设备输入的事件处理最好分开，方便做逻辑和显示分离
// 注册事件
func (g *Game) registerEvents() {
	// 请求登录
	g.eventMgr.RegisterEvent(core.EventIdOpLogin, g.onEventReqLogin)
	// 请求进入游戏
	g.eventMgr.RegisterEvent(core.EventIdOpEnterGame, g.onEventReqEnterGame)

	// 同步时间
	g.eventMgr.RegisterEvent(core.EventIdTimeSync, g.onEventTimeSync)
	// 同步时间结束
	g.eventMgr.RegisterEvent(core.EventIdTimeSyncEnd, g.onEventTimeSyncEnd)

	// 进入游戏
	g.eventMgr.RegisterEvent(core.EventIdPlayerEnterGame, g.onEventPlayerEnterGame)
	// 进入游戏完成
	g.eventMgr.RegisterEvent(core.EventIdPlayerEnterGameCompleted, g.onEventPlayerEnterGameCompleted)
	// 退出游戏
	g.eventMgr.RegisterEvent(core.EventIdPlayerExitGame, g.onEventPlayerExitGame)

	// 坦克移动
	g.eventMgr.RegisterEvent(core.EventIdTankMove, g.onEventTankMove)
	// 坦克停止移动
	g.eventMgr.RegisterEvent(core.EventIdTankStopMove, g.onEventTankStopMove)
	// 坦克改变
	g.eventMgr.RegisterEvent(core.EventIdTankChange, g.onEventTankChange)
	// 坦克恢复
	g.eventMgr.RegisterEvent(core.EventIdTankRestore, g.onEventTankRestore)
}

// 注销事件
func (g *Game) unregisterEvents() {
	g.eventMgr.UnregisterEvent(core.EventIdOpLogin, g.onEventReqLogin)
	g.eventMgr.UnregisterEvent(core.EventIdOpEnterGame, g.onEventReqEnterGame)
	g.eventMgr.UnregisterEvent(core.EventIdPlayerEnterGame, g.onEventPlayerEnterGame)
	g.eventMgr.UnregisterEvent(core.EventIdPlayerEnterGameCompleted, g.onEventPlayerEnterGameCompleted)
	g.eventMgr.UnregisterEvent(core.EventIdPlayerExitGame, g.onEventPlayerExitGame)
	g.eventMgr.UnregisterEvent(core.EventIdTankMove, g.onEventTankMove)
	g.eventMgr.UnregisterEvent(core.EventIdTankStopMove, g.onEventTankStopMove)
	g.eventMgr.UnregisterEvent(core.EventIdTankChange, g.onEventTankChange)
	g.eventMgr.UnregisterEvent(core.EventIdTankRestore, g.onEventTankRestore)
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
		core.Log().Warn("account type must string on req login, this is %v", t)
		return
	}
	password, o = p.(string)
	if !o {
		core.Log().Warn("password type must string on req login")
		return
	}
	err := g.net.SendLoginReq(account, password)
	if err != nil {
		core.Log().Warn("send login req err: %v", err)
		return
	}
	g.myAcc = account
	core.Log().Info("handle event: account %v password %v send login req", account, password)
}

// 请求进入游戏
func (g *Game) onEventReqEnterGame(args ...interface{}) {
	var account, sessionToken string
	var o bool
	account, o = args[0].(string)
	if !o {
		core.Log().Warn("account type must string on req enter game")
		return
	}
	err := g.net.SendEnterGameReq(account, sessionToken)
	if err != nil {
		glog.Warn("send enter game err: %v", err)
		return
	}
	glog.Info("handle event: account %v send enter game req", account)
}

// 处理玩家进入游戏事件
func (g *Game) onEventPlayerEnterGame(args ...interface{}) {
	if len(args) < 3 {
		glog.Warn("onEventEnterGame event args length cant less than 3")
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
		g.state = GameStateInGame
		glog.Info("handle event: my player (account: %v, uid: %v) entered game, tank %v", account, uid, *tank)
	} else {
		glog.Info("handle event: player (account: %v, uid: %v) entered game, tank %v", account, uid, *tank)
	}
}

// 处理进入游戏完成
func (g *Game) onEventPlayerEnterGameCompleted(args ...interface{}) {
	// 准备同步服务器时间
	if err := g.net.SendTimeSyncReq(); err != nil {
		glog.Error("handle event: send time sync request err: %v", err)
		return
	}
	glog.Info("handle event: my player (account: %v, uid: %v) enter game finished", g.myAcc, g.myId)
}

// 处理玩家离开游戏事件
func (g *Game) onEventPlayerExitGame(args ...interface{}) {
	if len(args) < 1 {
		glog.Warn("onEventPlayerExitGame event args length cant less 1")
		return
	}

	uid := args[0].(uint64)

	// 从播放管理器中删除
	g.playableMgr.RemovePlayerTankPlayable(uid)

	glog.Info("handle event: player (uid: %v) exited game", uid)
}

// 处理时间同步事件
func (g *Game) onEventTimeSync(args ...interface{}) {
	if err := g.net.SendTimeSyncReq(); err != nil {
		glog.Error("handle event: send time sync request err: %v", err)
		return
	}
	glog.Info("handle event: time sync")
}

// 处理时间同步结束事件
func (g *Game) onEventTimeSyncEnd(args ...interface{}) {
	glog.Info("handle event: time sync end")
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
		glog.Error("onEventTankChange event need 3 args")
		return
	}
	pid := args[0].(uint64)
	tank := args[1].(*object.Tank)
	glog.Info("handle event: player %v changed tank to %v", pid, tank.Id())
}

// 处理坦克恢复事件
func (g *Game) onEventTankRestore(args ...interface{}) {
	if len(args) < 2 {
		glog.Error("onEventTankRestore event need 3 args")
		return
	}
	pid := args[0].(uint64)
	tank := args[1].(*object.Tank)
	glog.Info("handle event: player %v restore tank id to %v", pid, tank.Id())
}
