package client_core

import (
	"fmt"
	"project_b/common/base"
	"project_b/common/object"
	custom_time "project_b/common/time"
	"project_b/core"
	"project_b/game_map"
	"project_b/log"
	"unsafe"

	"time"

	"project_b/common_data"
	"project_b/game_proto"

	gsnet_msg "github.com/huoshan017/gsnet/msg"
)

type MsgHandler struct {
	net       *NetClient
	inst      *core.Instance
	playerMgr *CPlayerManager
	invoker   base.IEventInvoker
}

func CreateMsgHandler(net *NetClient, playerMgr *CPlayerManager, invoker base.IEventInvoker) *MsgHandler {
	return &MsgHandler{
		net:       net,
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

func (h *MsgHandler) onConnect(sess *gsnet_msg.MsgSession) {
	log.Info("connected")
}

func (h *MsgHandler) onDisconnect(sess *gsnet_msg.MsgSession, err error) {
	log.Info("disconnected")
}

func (h *MsgHandler) onTick(sess *gsnet_msg.MsgSession, tick time.Duration) {

}

func (h *MsgHandler) onError(err error) {
	log.Info("get error: %v", err)
}

func (h *MsgHandler) registerNetMsgHandles() {
	h.net.RegisterHandle(gsnet_msg.MsgIdType(game_proto.MsgAccountLoginGameAck_Id), h.onLoginAck)
	h.net.RegisterHandle(gsnet_msg.MsgIdType(game_proto.MsgPlayerEnterGameAck_Id), h.onPlayerEnterGameAck)
	h.net.RegisterHandle(gsnet_msg.MsgIdType(game_proto.MsgPlayerEnterGameFinishNtf_Id), h.onPlayerEnterGameFinishNtf)
	h.net.RegisterHandle(gsnet_msg.MsgIdType(game_proto.MsgPlayerExitGameAck_Id), h.onPlayerExitGameAck)
	h.net.RegisterHandle(gsnet_msg.MsgIdType(game_proto.MsgTimeSyncAck_Id), h.onTimeSyncAck)
	h.net.RegisterHandle(gsnet_msg.MsgIdType(game_proto.MsgPlayerChangeTankAck_Id), h.onPlayerTankChangeAck)
	h.net.RegisterHandle(gsnet_msg.MsgIdType(game_proto.MsgPlayerChangeTankSync_Id), h.onPlayerTankChangeSync)
	h.net.RegisterHandle(gsnet_msg.MsgIdType(game_proto.MsgPlayerRestoreTankAck_Id), h.onPlayerTankRestoreAck)
	h.net.RegisterHandle(gsnet_msg.MsgIdType(game_proto.MsgPlayerRestoreTankSync_Id), h.onPlayerTankRestoreSync)
	h.net.RegisterHandle(gsnet_msg.MsgIdType(game_proto.MsgPlayerEnterGameSync_Id), h.onPlayerEnterGameSync)
	h.net.RegisterHandle(gsnet_msg.MsgIdType(game_proto.MsgPlayerExitGameSync_Id), h.onPlayerExitGameSync)
	h.net.RegisterHandle(gsnet_msg.MsgIdType(game_proto.MsgPlayerTankMoveAck_Id), h.onPlayerTankMoveAck)
	h.net.RegisterHandle(gsnet_msg.MsgIdType(game_proto.MsgPlayerTankStopMoveAck_Id), h.onPlayerTankStopMoveAck)
	h.net.RegisterHandle(gsnet_msg.MsgIdType(game_proto.MsgPlayerTankMoveSync_Id), h.onPlayerTankMoveSync)
	h.net.RegisterHandle(gsnet_msg.MsgIdType(game_proto.MsgPlayerTankStopMoveSync_Id), h.onPlayerTankStopMoveSync)
	h.net.RegisterHandle(gsnet_msg.MsgIdType(game_proto.MsgPlayerTankUpdatePosAck_Id), h.onPlayerTankUpdatePosAck)
	h.net.RegisterHandle(gsnet_msg.MsgIdType(game_proto.MsgPlayerTankUpdatePosSync_Id), h.onPlayerTankUpdatePosSync)
	h.net.RegisterHandle(gsnet_msg.MsgIdType(game_proto.MsgErrorAck_Id), h.onErrorAck)
}

// 登录处理
func (h *MsgHandler) onLoginAck(sess *gsnet_msg.MsgSession, msg any) error {
	ack, o := msg.(*game_proto.MsgAccountLoginGameAck)
	if !o {
		log.Warn("can't get to type *game_proto.MsgAccountLoginGameAck")
		return nil
	}

	if ack.Result != 0 {
		log.Warn("Account %v login result: %v", ack.Account, ack.Result)
		return nil
	}

	// todo 这里只是走一个登录流程，具体的登录逻辑等有了登录服务器再搞
	// 直接发进入游戏的消息
	err := h.net.SendEnterGameReq(string(ack.Account), string(ack.SessionToken))

	log.Info("Account %v login", ack.Account)

	return err
}

// 进入游戏处理
func (h *MsgHandler) onPlayerEnterGameAck(sess *gsnet_msg.MsgSession, msg any) error {
	ack, o := msg.(*game_proto.MsgPlayerEnterGameAck)
	if !o {
		log.Warn("can't get to type *game_proto.MsgPlayerEnterGameAck")
		return nil
	}

	// 载入地图
	config := game_map.MapConfigArray[ack.MapId]
	if !h.inst.Load(config) {
		log.Error("load map %v error", ack.MapId)
		return fmt.Errorf("load map %v failed", ack.MapId)
	}

	// 自己
	h.doPlayerEnter(ack.SelfTankInfo, true)

	// 其他玩家
	for i := 0; i < len(ack.OtherPlayerTankInfoList); i++ {
		tankInfo := ack.OtherPlayerTankInfoList[i]
		h.doPlayerEnter(tankInfo, false)
	}

	log.Info("my player entered game")

	return nil
}

// 进入游戏结束处理
func (h *MsgHandler) onPlayerEnterGameFinishNtf(sess *gsnet_msg.MsgSession, msg any) error {
	_, o := msg.(*game_proto.MsgPlayerEnterGameFinishNtf)
	if !o {
		log.Warn("cant transfer to type *game_proto.MsgPlayerEnterGameFinishNtf")
		return nil
	}

	// 向上层传递事件
	h.invoker.InvokeEvent(EventIdPlayerEnterGameCompleted)

	log.Info("my player entered game completed")

	return nil
}

// 退出游戏处理
func (h *MsgHandler) onPlayerExitGameAck(sess *gsnet_msg.MsgSession, msg any) error {
	h.doPlayerExit(h.myId())

	log.Info("my player exited game")

	return nil
}

// 时间同步处理
func (h *MsgHandler) onTimeSyncAck(sess *gsnet_msg.MsgSession, msg any) error {
	ack, o := msg.(*game_proto.MsgTimeSyncAck)
	if !o {
		log.Warn("cant transfer to type *game_proto.MsgTimeSyncAck")
		return nil
	}

	var st custom_time.CustomTime
	err := st.UnmarshalBinary(ack.ServerTime)
	if err != nil {
		return err
	}

	now := custom_time.Now()
	core.SetSyncRecvAndServerTime(now, st)

	if core.IsTimeSyncEnd() {
		h.invoker.InvokeEvent(EventIdTimeSyncEnd)
	} else {
		h.invoker.InvokeEvent(EventIdTimeSync)
	}

	log.Info("time sync client send time: %v, server time: %v, client recv time : %v, delay: %+v", core.GetSyncSendTime(), st, now, core.GetNetworkDelay())

	return nil
}

// 玩家坦克进入同步
func (h *MsgHandler) onPlayerEnterGameSync(sess *gsnet_msg.MsgSession, msg any) error {
	sync, o := msg.(*game_proto.MsgPlayerEnterGameSync)
	if !o {
		log.Warn("cant transfer to type *game_proto.MsgPlayerEnterGameSync")
		return nil
	}

	h.doPlayerEnter(sync.TankInfo, false)

	log.Info("sync player (account: %v, player_id: %v) entered game", sync.TankInfo.Account, sync.TankInfo.PlayerId)

	return nil
}

// 玩家坦克退出同步
func (h *MsgHandler) onPlayerExitGameSync(sess *gsnet_msg.MsgSession, msg any) error {
	sync, o := msg.(*game_proto.MsgPlayerExitGameSync)
	if !o {
		log.Warn("cant transfer to type *game_proto.MsgPlayerExitGameSync")
		return nil
	}

	h.doPlayerExit(sync.PlayerId)

	log.Info("sync player (player_id: %v) exited game", sync.PlayerId)

	return nil
}

// 玩家改变坦克处理
func (h *MsgHandler) onPlayerTankChangeAck(sess *gsnet_msg.MsgSession, msg any) error {
	ack, o := msg.(*game_proto.MsgPlayerChangeTankAck)
	if !o {
		log.Warn("cant transfer to type *game_proto.MsgPlayerChangeTankAck")
		return nil
	}

	if h.doTankChange(h.myId(), ack.ChangedTankInfo.Id) {
		log.Info("my player changed tank to %v", ack.ChangedTankInfo.Id)
	}

	return nil
}

// 玩家改变坦克同步
func (h *MsgHandler) onPlayerTankChangeSync(sess *gsnet_msg.MsgSession, msg any) error {
	sync, o := msg.(*game_proto.MsgPlayerChangeTankSync)
	if !o {
		log.Warn("cant transfer to type *game_proto.MsgPlayerChangeTankSync")
		return nil
	}

	if h.doTankChange(sync.ChangedTankInfo.PlayerId, sync.ChangedTankInfo.TankInfo.Id) {
		log.Info("sync player %v change tank to %v", sync.ChangedTankInfo.PlayerId, sync.ChangedTankInfo.TankInfo.Id)
	}

	return nil
}

// 玩家恢复坦克处理
func (h *MsgHandler) onPlayerTankRestoreAck(sess *gsnet_msg.MsgSession, msg any) error {
	ack, o := msg.(*game_proto.MsgPlayerRestoreTankAck)
	if !o {
		log.Warn("cant transfer to type *game_proto.MsgPlayerRestoreTankAck")
		return nil
	}

	if h.doTankRestore(h.myId(), ack.TankId) {
		log.Info("my player restored tank to %v", ack.TankId)
	}

	return nil
}

// 玩家恢复坦克同步处理
func (h *MsgHandler) onPlayerTankRestoreSync(sess *gsnet_msg.MsgSession, msg any) error {
	sync, o := msg.(*game_proto.MsgPlayerRestoreTankSync)
	if !o {
		log.Warn("cant transfer to type *game_proto.MsgPlayerRestoreTankSync")
		return nil
	}

	if h.doTankRestore(sync.PlayerId, sync.TankId) {
		log.Info("player %v restore tank to %v", sync.PlayerId, sync.TankId)
	}

	return nil
}

// 本玩家移动回应处理
func (h *MsgHandler) onPlayerTankMoveAck(sess *gsnet_msg.MsgSession, msg any) error {
	_, o := msg.(*game_proto.MsgPlayerTankMoveAck)
	if !o {
		log.Warn("cant transfer to type *game_proto.MsgPlayerTankMoveAck")
		return nil
	}
	return nil
}

// 其他玩家移动同步处理
func (h *MsgHandler) onPlayerTankMoveSync(sess *gsnet_msg.MsgSession, msg any) error {
	sync, o := msg.(*game_proto.MsgPlayerTankMoveSync)
	if !o {
		log.Warn("cant transfer to type *game_proto.MsgPlayerTankMoveSync")
		return nil
	}

	orientation := object.Dir2Orientation(object.Direction(sync.MoveInfo.Direction))
	h.inst.PushFrame(h.inst.GetFrame(), sync.PlayerId, core.CMD_TANK_MOVE, []int64{int64(orientation)})

	log.Debug("Player %v move sync", sync.PlayerId)

	return nil
}

// 本玩家停止移动返回处理
func (h *MsgHandler) onPlayerTankStopMoveAck(sess *gsnet_msg.MsgSession, msg any) error {
	_, o := msg.(*game_proto.MsgPlayerTankStopMoveAck)
	if !o {
		log.Warn("cant transfer to type *game_proto.MsgPlayerTankStopMoveAck")
		return nil
	}
	return nil
}

// 其他玩家停止移动同步处理
func (h *MsgHandler) onPlayerTankStopMoveSync(sess *gsnet_msg.MsgSession, msg any) error {
	sync, o := msg.(*game_proto.MsgPlayerTankMoveSync)
	if !o {
		log.Warn("cant transfer to type *game_proto.MsgPlayerTankMoveSync")
		return nil
	}

	h.inst.PushFrame(h.inst.GetFrame(), sync.PlayerId, core.CMD_TANK_STOP, nil)

	log.Debug("Player %v stop move sync", sync.PlayerId)

	return nil
}

// 本玩家的坦克位置更新返回
func (h *MsgHandler) onPlayerTankUpdatePosAck(sess *gsnet_msg.MsgSession, msg any) error {
	ack, o := msg.(*game_proto.MsgPlayerTankUpdatePosAck)
	if !o {
		log.Warn("cant transfer to type *game_proto.MsgPlayerTankUpdatePosAck")
		return nil
	}

	switch ack.State {
	case game_proto.MovementState_StartMove:
	case game_proto.MovementState_Moving:
	case game_proto.MovementState_ToStop:
	}

	//gslog.Debug("My tank update pos ack: %v", &ack)

	return nil
}

// 其他玩家坦克位置更新同步
func (h *MsgHandler) onPlayerTankUpdatePosSync(sess *gsnet_msg.MsgSession, msg any) error {
	return nil
}

// 处理错误返回
func (h *MsgHandler) onErrorAck(sess *gsnet_msg.MsgSession, msg any) error {
	ack, o := msg.(*game_proto.MsgErrorAck)
	if !o {
		log.Warn("cant transfer to type *game_proto.MsgErrorAck")
		return nil
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

	log.Error(s)

	return nil
}

// 处理玩家进入
func (h *MsgHandler) doPlayerEnter(p *game_proto.PlayerAccountTankInfo, isMe bool) {
	// 创建玩家
	cplayer := NewCPlayer(p.Account, p.PlayerId, h.net)

	// 加入玩家管理器
	if isMe {
		h.playerMgr.AddMe(cplayer)
	} else {
		h.playerMgr.Add(cplayer)
	}

	// 玩家坦克进入主逻辑
	h.inst.CheckAndStart(h.playerMgr.GetPlayerList())

	// 向上层传递事件
	h.invoker.InvokeEvent(EventIdPlayerEnterGame, p.Account, p.PlayerId)
}

// 处理玩家退出
func (h *MsgHandler) doPlayerExit(playerId uint64) {
	// 从管理器中删除玩家
	h.playerMgr.Remove(playerId)

	// 坦克离开主逻辑
	//h.logic.PlayerLeave(playerId)

	// 向上传递事件
	h.invoker.InvokeEvent(EventIdPlayerExitGame, playerId)
}

// 处理坦克改变
func (h *MsgHandler) doTankChange(playerId uint64, changedTankId int32) bool {
	player := h.playerMgr.Get(playerId)
	if player == nil {
		log.Error("not found player %v to change tank", playerId)
		return false
	}

	// 坦克改变
	h.inst.PushFrame(h.inst.GetFrame(), playerId, core.CMD_TANK_CHANGE, []int64{int64(uintptr(unsafe.Pointer(common_data.TankConfigData[changedTankId])))})

	// 向上传递事件
	//h.invoker.InvokeEvent(common.EventIdTankChange, playerId, h.logic.GetPlayerTank(playerId))

	return true
}

// 处理坦克恢复
func (h *MsgHandler) doTankRestore(playerId uint64, tankId int32) bool {
	player := h.playerMgr.Get(playerId)
	if player == nil {
		log.Error("not found player %v to restore tank", playerId)
		return false
	}

	// 恢复坦克
	h.inst.PushFrame(h.inst.GetFrame(), playerId, core.CMD_TANK_RESTORE, nil)

	// 向上传递事件
	//h.invoker.InvokeEvent(common.EventIdTankChange, playerId, h.logic.GetPlayerTank(playerId))

	return true
}

// 我的玩家id
func (h *MsgHandler) myId() uint64 {
	return h.playerMgr.GetMe().Id()
}
