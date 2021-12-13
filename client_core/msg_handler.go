package client_core

import (
	"project_b/common"
	"project_b/common/base"
	"project_b/common/object"
	custom_time "project_b/common/time"

	"time"

	//"project_b/common/time"
	"project_b/common_data"
	"project_b/game_proto"

	"github.com/huoshan017/gsnet"
	"google.golang.org/protobuf/proto"
)

type MsgHandler struct {
	net       *NetClient
	logic     *common.GameLogic
	playerMgr *CPlayerManager
	invoker   base.IEventInvoker
}

func CreateMsgHandler(net *NetClient, logic *common.GameLogic, playerMgr *CPlayerManager, invoker base.IEventInvoker) *MsgHandler {
	return &MsgHandler{
		net:       net,
		logic:     logic,
		playerMgr: playerMgr,
		invoker:   invoker,
	}
}

func (h *MsgHandler) Init() {
	h.setNetEventHandles()
	h.registerNetMsgHandles()
}

func (h *MsgHandler) setNetEventHandles() {
	h.net.SetConnectHandle(h.onConnect)
	h.net.SetDisconnectHandle(h.onDisconnect)
	h.net.SetTickHandle(h.onTick)
	h.net.SetErrorHandle(h.onError)
}

func (h *MsgHandler) onConnect(sess gsnet.ISession) {
	Log().Info("connected")
}

func (h *MsgHandler) onDisconnect(sess gsnet.ISession, err error) {
	Log().Info("disconnected")
}

func (h *MsgHandler) onTick(sess gsnet.ISession, tick time.Duration) {

}

func (h *MsgHandler) onError(err error) {
	Log().Info("get error: %v", err)
}

func (h *MsgHandler) registerNetMsgHandles() {
	h.net.RegisterHandle(uint32(game_proto.MsgAccountLoginGameAck_Id), h.onLoginAck)
	h.net.RegisterHandle(uint32(game_proto.MsgPlayerEnterGameAck_Id), h.onPlayerEnterGameAck)
	h.net.RegisterHandle(uint32(game_proto.MsgPlayerEnterGameFinishNtf_Id), h.onPlayerEnterGameFinishNtf)
	h.net.RegisterHandle(uint32(game_proto.MsgPlayerExitGameAck_Id), h.onPlayerExitGameAck)
	h.net.RegisterHandle(uint32(game_proto.MsgTimeSyncAck_Id), h.onTimeSyncAck)
	h.net.RegisterHandle(uint32(game_proto.MsgPlayerChangeTankAck_Id), h.onPlayerTankChangeAck)
	h.net.RegisterHandle(uint32(game_proto.MsgPlayerChangeTankSync_Id), h.onPlayerTankChangeSync)
	h.net.RegisterHandle(uint32(game_proto.MsgPlayerRestoreTankAck_Id), h.onPlayerTankRestoreAck)
	h.net.RegisterHandle(uint32(game_proto.MsgPlayerRestoreTankSync_Id), h.onPlayerTankRestoreSync)
	h.net.RegisterHandle(uint32(game_proto.MsgPlayerEnterGameSync_Id), h.onPlayerEnterGameSync)
	h.net.RegisterHandle(uint32(game_proto.MsgPlayerExitGameSync_Id), h.onPlayerExitGameSync)
	h.net.RegisterHandle(uint32(game_proto.MsgPlayerTankMoveAck_Id), h.onPlayerTankMoveAck)
	h.net.RegisterHandle(uint32(game_proto.MsgPlayerTankStopMoveAck_Id), h.onPlayerTankStopMoveAck)
	h.net.RegisterHandle(uint32(game_proto.MsgPlayerTankMoveSync_Id), h.onPlayerTankMoveSync)
	h.net.RegisterHandle(uint32(game_proto.MsgPlayerTankStopMoveSync_Id), h.onPlayerTankStopMoveSync)
	h.net.RegisterHandle(uint32(game_proto.MsgErrorAck_Id), h.onErrorAck)
}

// 登录处理
func (h *MsgHandler) onLoginAck(sess gsnet.ISession, data []byte) error {
	var ack game_proto.MsgAccountLoginGameAck
	err := proto.Unmarshal(data, &ack)
	if err != nil {
		return err
	}

	if ack.Result != 0 {
		Log().Warn("Account %v login result: %v", ack.Account, ack.Result)
		return nil
	}

	// todo 这里只是走一个登录流程，具体的登录逻辑等有了登录服务器再搞
	// 直接发进入游戏的消息
	err = h.net.SendEnterGameReq(string(ack.Account), string(ack.SessionToken))

	Log().Info("Account %v login", ack.Account)

	return err
}

// 进入游戏处理
func (h *MsgHandler) onPlayerEnterGameAck(sess gsnet.ISession, data []byte) error {
	var ack game_proto.MsgPlayerEnterGameAck
	err := proto.Unmarshal(data, &ack)
	if err != nil {
		return err
	}

	// 自己
	h.doPlayerEnter(ack.SelfTankInfo, true)

	// 其他玩家
	for i := 0; i < len(ack.OtherPlayerTankInfoList); i++ {
		tankInfo := ack.OtherPlayerTankInfoList[i]
		h.doPlayerEnter(tankInfo, false)
	}

	Log().Info("my player entered game")

	return nil
}

// 进入游戏结束处理
func (h *MsgHandler) onPlayerEnterGameFinishNtf(sess gsnet.ISession, data []byte) error {
	var ntf game_proto.MsgPlayerEnterGameFinishNtf
	err := proto.Unmarshal(data, &ntf)
	if err != nil {
		return err
	}

	// 向上层传递事件
	h.invoker.InvokeEvent(EventIdPlayerEnterGameCompleted)

	Log().Info("my player entered game completed")

	return nil
}

// 退出游戏处理
func (h *MsgHandler) onPlayerExitGameAck(sess gsnet.ISession, data []byte) error {
	h.doPlayerExit(h.myId())

	Log().Info("my player exited game")

	return nil
}

// 时间同步处理
func (h *MsgHandler) onTimeSyncAck(sess gsnet.ISession, data []byte) error {
	var ack game_proto.MsgTimeSyncAck
	err := proto.Unmarshal(data, &ack)
	if err != nil {
		return err
	}

	var st custom_time.CustomTime
	err = st.UnmarshalBinary(ack.ServerTime)
	if err != nil {
		return err
	}

	now := custom_time.Now()
	SetSyncRecvAndServerTime(now, st)

	if IsTimeSyncEnd() {
		h.invoker.InvokeEvent(EventIdTimeSyncEnd)
	} else {
		h.invoker.InvokeEvent(EventIdTimeSync)
	}

	gslog.Info("time sync client send time: %v, server time: %v, client recv time : %v, delay: %+v", GetSyncSendTime(), st, now, GetNetworkDelay())

	return nil
}

// 玩家坦克进入同步
func (h *MsgHandler) onPlayerEnterGameSync(sess gsnet.ISession, data []byte) error {
	var sync game_proto.MsgPlayerEnterGameSync
	err := proto.Unmarshal(data, &sync)
	if err != nil {
		return err
	}

	h.doPlayerEnter(sync.TankInfo, false)

	Log().Info("sync player (account: %v, player_id: %v) entered game", sync.TankInfo.Account, sync.TankInfo.PlayerId)

	return nil
}

// 玩家坦克退出同步
func (h *MsgHandler) onPlayerExitGameSync(sess gsnet.ISession, data []byte) error {
	var sync game_proto.MsgPlayerExitGameSync
	err := proto.Unmarshal(data, &sync)
	if err != nil {
		return err
	}

	h.doPlayerExit(sync.PlayerId)

	Log().Info("sync player (player_id: %v) exited game", sync.PlayerId)

	return nil
}

// 玩家改变坦克处理
func (h *MsgHandler) onPlayerTankChangeAck(sess gsnet.ISession, data []byte) error {
	var ack game_proto.MsgPlayerChangeTankAck
	err := proto.Unmarshal(data, &ack)
	if err != nil {
		return err
	}

	if h.doTankChange(h.myId(), ack.ChangedTankInfo.Id) {
		Log().Info("my player changed tank to %v", ack.ChangedTankInfo.Id)
	}

	return nil
}

// 玩家改变坦克同步
func (h *MsgHandler) onPlayerTankChangeSync(sess gsnet.ISession, data []byte) error {
	var sync game_proto.MsgPlayerChangeTankSync
	err := proto.Unmarshal(data, &sync)
	if err != nil {
		return err
	}

	if h.doTankChange(sync.ChangedTankInfo.PlayerId, sync.ChangedTankInfo.TankInfo.Id) {
		Log().Info("sync player %v change tank to %v", sync.ChangedTankInfo.PlayerId, sync.ChangedTankInfo.TankInfo.Id)
	}

	return nil
}

// 玩家恢复坦克处理
func (h *MsgHandler) onPlayerTankRestoreAck(sess gsnet.ISession, data []byte) error {
	var ack game_proto.MsgPlayerRestoreTankAck
	err := proto.Unmarshal(data, &ack)
	if err != nil {
		return err
	}

	if h.doTankRestore(h.myId(), ack.TankId) {
		Log().Info("my player restored tank to %v", ack.TankId)
	}

	return nil
}

// 玩家恢复坦克同步处理
func (h *MsgHandler) onPlayerTankRestoreSync(sess gsnet.ISession, data []byte) error {
	var sync game_proto.MsgPlayerRestoreTankSync
	err := proto.Unmarshal(data, &sync)
	if err != nil {
		return err
	}

	if h.doTankRestore(sync.PlayerId, sync.TankId) {
		Log().Info("player %v restore tank to %v", sync.PlayerId, sync.TankId)
	}

	return nil
}

// 本玩家移动回应处理
func (h *MsgHandler) onPlayerTankMoveAck(sess gsnet.ISession, data []byte) error {
	var ack game_proto.MsgPlayerTankMoveAck
	err := proto.Unmarshal(data, &ack)
	if err != nil {
		return err
	}
	return nil
}

// 其他玩家移动同步处理
func (h *MsgHandler) onPlayerTankMoveSync(sess gsnet.ISession, data []byte) error {
	var sync game_proto.MsgPlayerTankMoveSync
	err := proto.Unmarshal(data, &sync)
	if err != nil {
		return err
	}

	h.logic.PlayerTankMove(sync.PlayerId, object.Direction(sync.MoveInfo.Direction))

	return nil
}

// 本玩家停止移动返回处理
func (h *MsgHandler) onPlayerTankStopMoveAck(sess gsnet.ISession, data []byte) error {
	var ack game_proto.MsgPlayerTankStopMoveAck
	err := proto.Unmarshal(data, &ack)
	if err != nil {
		return err
	}
	return nil
}

// 其他玩家停止移动同步处理
func (h *MsgHandler) onPlayerTankStopMoveSync(sess gsnet.ISession, data []byte) error {
	var sync game_proto.MsgPlayerTankMoveSync
	err := proto.Unmarshal(data, &sync)
	if err != nil {
		return err
	}

	h.logic.PlayerTankStopMove(sync.PlayerId)

	return nil
}

// 处理错误返回
func (h *MsgHandler) onErrorAck(sess gsnet.ISession, data []byte) error {
	var ack game_proto.MsgErrorAck
	err := proto.Unmarshal(data, &ack)
	if err != nil {
		return err
	}

	var s string
	switch ack.Error {
	case game_proto.ErrorId_ACCOUNT_IS_LOGGIN:
		s = "error: account is logining"
	case game_proto.ErrorId_DIFFERENT_PLAYER_ENTER_SAME_SESSION:
		s = "error: different player enter same session"
	case game_proto.ErrorId_PLAYER_CHANGE_TANK_FAILED:
		s = "error: player change tank failed"
	case game_proto.ErrorId_PLAYER_ENTER_GAME_REPEATED:
		s = "error: player enter game repeated"
	case game_proto.ErrorId_PLAYER_RESTORE_TANK_FAILED:
		s = "error: player restore tank failed"
	case game_proto.ErrorId_PLAYER_ENTERING_GAME:
		s = "error: player entering game"
	case game_proto.ErrorId_SESSION_INTERNAL_ERROR:
		s = "error: session internal error"
	case game_proto.ErrorId_ACCOUNT_NOT_FOUND:
		s = "error: account not found"
	case game_proto.ErrorId_INVALID_ACCOUNT:
		s = "error: invalid account"
	}

	gslog.Error(s)

	return nil
}

// 处理玩家进入
func (h *MsgHandler) doPlayerEnter(p *game_proto.PlayerAccountTankInfo, isMe bool) {
	// 创建玩家
	cplayer := NewCPlayer(p.Account, p.PlayerId, h.net)

	// 初始化坦克
	cplayer.InitTankFromProto(p.TankInfo)

	// 加入玩家管理器
	if isMe {
		h.playerMgr.AddMe(cplayer)
	} else {
		h.playerMgr.Add(cplayer)
	}

	// 玩家坦克进入主逻辑
	h.logic.PlayerTankEnter(p.PlayerId, cplayer.GetTank())

	// 向上层传递事件
	h.invoker.InvokeEvent(EventIdPlayerEnterGame, p.Account, p.PlayerId, cplayer.GetTank())
}

// 处理玩家退出
func (h *MsgHandler) doPlayerExit(playerId uint64) {
	// 从管理器中删除玩家
	h.playerMgr.Remove(playerId)

	// 坦克离开主逻辑
	h.logic.PlayerTankLeave(playerId)

	// 向上传递事件
	h.invoker.InvokeEvent(EventIdPlayerExitGame, playerId)
}

// 处理坦克改变
func (h *MsgHandler) doTankChange(playerId uint64, changedTankId int32) bool {
	player := h.playerMgr.Get(playerId)
	if player == nil {
		gslog.Error("not found player %v to change tank", playerId)
		return false
	}

	// 坦克改变
	h.logic.PlayerTankChange(playerId, common_data.TankConfigData[changedTankId])

	// 向上传递事件
	h.invoker.InvokeEvent(EventIdTankChange, playerId, h.logic.GetPlayerTank(playerId))

	return true
}

// 处理坦克恢复
func (h *MsgHandler) doTankRestore(playerId uint64, tankId int32) bool {
	player := h.playerMgr.Get(playerId)
	if player == nil {
		Log().Error("not found player %v to restore tank", playerId)
		return false
	}

	// 恢复坦克
	h.logic.PlayerTankRestore(playerId)

	// 向上传递事件
	h.invoker.InvokeEvent(EventIdTankChange, playerId, h.logic.GetPlayerTank(playerId))

	return true
}

// 我的玩家id
func (h *MsgHandler) myId() uint64 {
	return h.playerMgr.GetMe().Id()
}
