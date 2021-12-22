package main

import (
	"project_b/common"
	"project_b/common/object"
	custom_time "project_b/common/time"
	"project_b/common_data"
	"project_b/game_proto"
	"project_b/utils"
	"time"

	"google.golang.org/protobuf/reflect/protoreflect"
)

// 局部玩家数据
type playerData struct {
	tank        *object.Tank
	sessHandler *GameMsgHandler
	pid         uint64
	account     string
}

// 发送消息
func (d *playerData) send(msgid uint32, msg protoreflect.ProtoMessage) error {
	return d.sessHandler.Send(msgid, msg)
}

// 发送错误
func (d *playerData) sendError(err game_proto.ErrorId) error {
	return d.sessHandler.SendError(err)
}

// 游戏逻辑线程
type GameLogicThread struct {
	common.MsgLogicProc
	gameLogic *common.GameLogic
}

// 创建游戏逻辑线程
func CreateGameLogicThread() *GameLogicThread {
	t := &GameLogicThread{
		MsgLogicProc: *common.CreateMsgLogicProc(),
		gameLogic:    common.NewGameLogic(),
	}
	t.registerHandles()
	return t
}

// 初始化
func (t *GameLogicThread) registerHandles() {
	t.SetTickHandle(t.onTick, time.Duration(common_data.GameLogicTick))
	t.RegisterHandle(uint32(game_proto.MsgPlayerTankMoveReq_Id), t.onPlayerTankMoveReq)
	t.RegisterHandle(uint32(game_proto.MsgPlayerTankStopMoveReq_Id), t.onPlayerTankStopMoveReq)
	t.RegisterHandle(uint32(game_proto.MsgPlayerTankUpdatePosReq_Id), t.onPlayerTankUpdatePosReq)
	t.RegisterHandle(uint32(game_proto.MsgPlayerChangeTankReq_Id), t.onPlayerTankChange)
	t.RegisterHandle(uint32(game_proto.MsgPlayerRestoreTankReq_Id), t.onPlayerTankRestore)
}

// 玩家进入游戏主逻辑
func (t *GameLogicThread) PlayerEnter(pid uint64, data *playerData) {
	t.PushAgent(pid, data, func(d interface{}) error {
		pd := d.(*playerData)
		t.gameLogic.PlayerTankEnter(pid, pd.tank)
		return t.onPlayerTankEnterReq(pd)
	})
}

// 玩家离开游戏主逻辑
func (t *GameLogicThread) PlayerLeave(pid uint64) {
	d := &playerData{
		pid: pid,
	}
	t.DeleteAgent(pid, d, func(agentKey interface{}) error {
		pid := agentKey.(uint64)
		t.gameLogic.PlayerTankLeave(pid)
		var err error
		pd := t.getPlayerData(agentKey)
		if pd != nil {
			err = t.onPlayerTankLeaveReq(pid, pd.sessHandler)
		}
		return err
	})
}

// 重置玩家的会话处理器
func (t *GameLogicThread) PlayerResetHandler(pid uint64, sessHandler *GameMsgHandler, tank *object.Tank) {
	d := &playerData{
		sessHandler: sessHandler,
		pid:         pid,
	}
	t.UpdateAgent(pid, d, func(data interface{}) error {
		t.gameLogic.PlayerTankLeave(pid)
		t.gameLogic.PlayerTankEnter(pid, tank)
		return nil
	})
}

// 坦克进入
func (t *GameLogicThread) onPlayerTankEnterReq(pd *playerData) error {
	// 同步其他玩家和敌人
	var ack game_proto.MsgPlayerEnterGameAck
	tankList := t.gameLogic.GetPlayerTankList()
	gslog.Debug("!!!! tankList: %+v", tankList)
	for _, tank := range tankList {
		if pd.pid == tank.PlayerId { // 自己
			ack.SelfTankInfo = &game_proto.PlayerAccountTankInfo{}
			ack.SelfTankInfo.Account = pd.account
			ack.SelfTankInfo.PlayerId = pd.pid
			ack.SelfTankInfo.TankInfo = &game_proto.TankInfo{}
			utils.TankInfo2ProtoInfo(tank.Tank, ack.SelfTankInfo.TankInfo)
		} else { // 别人
			playerTankInfo := &game_proto.PlayerAccountTankInfo{}
			p := t.getPlayerData(tank.PlayerId)
			if p == nil {
				gslog.Warn("not found player data by pid %v", tank.PlayerId)
				continue
			}
			playerTankInfo.Account = p.account
			playerTankInfo.PlayerId = p.pid
			playerTankInfo.TankInfo = &game_proto.TankInfo{}
			utils.TankInfo2ProtoInfo(p.tank, playerTankInfo.TankInfo)
			ack.OtherPlayerTankInfoList = append(ack.OtherPlayerTankInfoList, playerTankInfo)
		}
	}
	err := pd.send(uint32(game_proto.MsgPlayerEnterGameAck_Id), &ack)
	if err != nil {
		return err
	}

	var ntf game_proto.MsgPlayerEnterGameFinishNtf
	ntf.ServerTime, err = time.Now().MarshalBinary()
	if err != nil {
		return err
	}

	err = pd.send(uint32(game_proto.MsgPlayerEnterGameFinishNtf_Id), &ntf)
	if err != nil {
		return err
	}

	// 同步给其他玩家
	if t.GetAgentCountNoLock() > 1 {
		var sync game_proto.MsgPlayerEnterGameSync
		sync.TankInfo = &game_proto.PlayerAccountTankInfo{}
		sync.TankInfo.Account = pd.account
		sync.TankInfo.PlayerId = pd.pid
		sync.TankInfo.TankInfo = &game_proto.TankInfo{}
		utils.TankInfo2ProtoInfo(pd.tank, sync.TankInfo.TankInfo)
		err = t.broadcastMsgExceptPlayer(uint32(game_proto.MsgPlayerEnterGameSync_Id), &sync, pd.pid)
	}
	return err
}

// 坦克离开
func (t *GameLogicThread) onPlayerTankLeaveReq(pid uint64, sessHandler *GameMsgHandler) error {
	var err error
	if t.GetAgentCountNoLock() > 1 {
		var sync game_proto.MsgPlayerExitGameSync
		sync.PlayerId = pid
		err = t.broadcastMsgExceptPlayer(uint32(game_proto.MsgPlayerExitGameSync_Id), &sync, pid)
	}
	return err
}

// tick处理
func (t *GameLogicThread) onTick(tick time.Duration) {
	t.gameLogic.Update(custom_time.Duration(tick))
}

// 坦克移动
func (t *GameLogicThread) onPlayerTankMoveReq(key common.AgentKey, msg common.MsgData) error {
	pd := t.getPlayerData(key)
	if pd == nil {
		gslog.Fatal("player %v not found in GameLogicThread", pd.pid)
		return nil
	}
	m := msg.(*game_proto.MsgPlayerTankMoveReq)

	// 检测移动数据的合法性，计算当前位置
	t.gameLogic.PlayerTankMove(pd.pid, object.Direction(m.MoveInfo.Direction))

	var ack game_proto.MsgPlayerTankMoveAck
	err := pd.send(uint32(game_proto.MsgPlayerTankMoveAck_Id), &ack)
	if err != nil {
		return err
	}
	if t.GetAgentCountNoLock() > 0 {
		var sync game_proto.MsgPlayerTankMoveSync
		sync.PlayerId = pd.pid
		sync.MoveInfo = m.MoveInfo
		err = t.broadcastMsgExceptPlayer(uint32(game_proto.MsgPlayerTankMoveSync_Id), &sync, pd.pid)
	}
	return err
}

// 坦克移动设置位置
func (t *GameLogicThread) onPlayerTankUpdatePosReq(key common.AgentKey, msg common.MsgData) error {
	pd := t.getPlayerData(key)
	if pd == nil {
		gslog.Error("player %v not found in GameLogicThread", pd.pid)
		return nil
	}
	m := msg.(*game_proto.MsgPlayerTankUpdatePosReq)

	switch m.State {
	case game_proto.MovementState_StartMove:
	case game_proto.MovementState_Moving:
	case game_proto.MovementState_ToStop:
	}

	var ack game_proto.MsgPlayerTankUpdatePosAck
	ack.State = m.State
	ack.MoveInfo = m.MoveInfo
	err := pd.send(uint32(game_proto.MsgPlayerTankUpdatePosAck_Id), &ack)
	if err != nil {
		return err
	}

	if t.GetAgentCountNoLock() > 0 {
		var sync game_proto.MsgPlayerTankUpdatePosSync
		sync.PlayerId = pd.pid
		sync.State = m.State
		sync.MoveInfo = m.MoveInfo
		err = t.broadcastMsgExceptPlayer(uint32(game_proto.MsgPlayerTankUpdatePosSync_Id), &sync, pd.pid)
	}
	return err
}

// 坦克停止移动
func (t *GameLogicThread) onPlayerTankStopMoveReq(key common.AgentKey, msg common.MsgData) error {
	pd := t.getPlayerData(key)
	if pd == nil {
		gslog.Fatal("player %v not found in GameLogicThread", pd.pid)
		return nil
	}

	t.gameLogic.PlayerTankStopMove(pd.pid)

	var ack game_proto.MsgPlayerTankStopMoveAck
	err := pd.send(uint32(game_proto.MsgPlayerTankStopMoveAck_Id), &ack)
	if err != nil {
		return err
	}
	if t.GetAgentCountNoLock() > 0 {
		var sync game_proto.MsgPlayerTankStopMoveSync
		sync.PlayerId = pd.pid
		err = t.broadcastMsgExceptPlayer(uint32(game_proto.MsgPlayerTankStopMoveSync_Id), &sync, pd.pid)
	}
	return err
}

// 玩家坦克改变
func (t *GameLogicThread) onPlayerTankChange(key common.AgentKey, msg common.MsgData) error {
	pd := t.getPlayerData(key)
	if pd == nil {
		gslog.Fatal("player %v not found in GameLogicThread", pd.pid)
		return nil
	}

	var ack game_proto.MsgPlayerChangeTankAck
	req := msg.(*game_proto.MsgPlayerChangeTankReq)
	res := t.gameLogic.PlayerTankChange(pd.pid, common_data.TankConfigData[req.TankId])
	if !res {
		gslog.Error("player %v change tank error", pd.pid)
		return pd.sendError(game_proto.ErrorId_PLAYER_CHANGE_TANK_FAILED)
	}

	tank := t.gameLogic.GetPlayerTank(pd.pid)
	ack.ChangedTankInfo = &game_proto.TankInfo{}
	utils.TankInfo2ProtoInfo(tank, ack.ChangedTankInfo)
	gslog.Info("player %v changed tank to %v", pd.pid, tank.Id())
	err := pd.send(uint32(game_proto.MsgPlayerChangeTankAck_Id), &ack)
	if err != nil {
		return err
	}

	// 同步给其他玩家
	if t.GetAgentCountNoLock() > 1 {
		x, y := tank.Pos()
		var sync = game_proto.MsgPlayerChangeTankSync{
			ChangedTankInfo: &game_proto.PlayerTankInfo{
				PlayerId: pd.pid,
				TankInfo: &game_proto.TankInfo{
					Id:        tank.Id(),
					Level:     tank.Level(),
					Direction: int32(tank.Dir()),
					CurrSpeed: tank.CurrentSpeed(),
					CurrPos:   &game_proto.Pos{X: x, Y: y},
				},
			},
		}
		err = t.broadcastMsgExceptPlayer(uint32(game_proto.MsgPlayerChangeTankSync_Id), &sync, pd.pid)
	}
	return err
}

// 玩家坦克恢复
func (t *GameLogicThread) onPlayerTankRestore(key common.AgentKey, msg common.MsgData) error {
	var ack game_proto.MsgPlayerRestoreTankAck

	pd := t.getPlayerData(key)
	tankId := t.gameLogic.PlayerTankRestore(pd.pid)

	if tankId <= 0 {
		gslog.Info("player %v restore tank failed", pd.pid)
		return pd.sendError(game_proto.ErrorId_PLAYER_RESTORE_TANK_FAILED)
	}

	ack.TankId = tankId
	gslog.Info("player %v restored tank", pd.pid)
	err := pd.send(uint32(game_proto.MsgPlayerRestoreTankAck_Id), &ack)

	if t.GetAgentCountNoLock() > 1 {
		var sync = game_proto.MsgPlayerRestoreTankSync{
			PlayerId: pd.pid,
			TankId:   tankId,
		}
		err = t.broadcastMsgExceptPlayer(uint32(game_proto.MsgPlayerRestoreTankSync_Id), &sync, pd.pid)
	}

	return err
}

// 内部函数，获得playerData数据
func (t *GameLogicThread) getPlayerData(key common.AgentKey) *playerData {
	pd := t.GetAgentNoLock(key).(*playerData)
	return pd
}

// 广播消息
func (t *GameLogicThread) broadcastMsg(msgid uint32, msg protoreflect.ProtoMessage) error {
	return t.broadcastMsgExceptPlayer(msgid, msg, 0)
}

// 广播消息除了某玩家
func (t *GameLogicThread) broadcastMsgExceptPlayer(msgid uint32, msg protoreflect.ProtoMessage, uid uint64) error {
	var err error
	players := t.GetAgentMapNoLock()
	for _, d := range players {
		pd := d.(*playerData)
		if uid == 0 || pd.pid == uid {
			continue
		}
		err = pd.send(msgid, msg)
		if err != nil {
			break
		}
	}
	return err
}
