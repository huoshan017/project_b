package main

import (
	"errors"
	"project_b/common_data"
	"project_b/game_proto"
	"project_b/utils"
	"time"

	"github.com/huoshan017/gsnet"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type GameMsgHandler struct {
	gsnet.MsgHandler
	service               *GameService
	sess                  gsnet.ISession
	lastCheckDiscDuration time.Duration
}

// 初始化
func (h *GameMsgHandler) Init(args ...interface{}) {
	h.MsgHandler = *gsnet.NewMsgHandler(&gsnet.DefaultMsgProto{})
	h.service = args[0].(*GameService)
	h.registerHandles()
}

// 连接事件
func (h *GameMsgHandler) OnConnect(sess gsnet.ISession) {
	// 连接后把会话缓存起来
	h.sess = sess
	gslog.Info("new session %v connected", sess.GetId())
}

// 断开事件
func (h *GameMsgHandler) OnDisconnect(sess gsnet.ISession, err error) {
	h.afterPlayerDisconnect(sess)
	gslog.Info("session %v disconnected", sess.GetId())
	h.sess = nil
}

// 定时器事件
func (h *GameMsgHandler) OnTick(sess gsnet.ISession, tick time.Duration) {
	h.lastCheckDiscDuration += tick
	// 0.5秒检测一次
	if h.lastCheckDiscDuration >= time.Microsecond*500 {
		p, err := h.toPlayer(sess)
		if err != nil {
			return
		}
		// 这里断开之后，后面会走到OnDisconnect中
		kicker := p.GetKicker()
		if kicker != nil {
			kicker.CheckDisconnectNotification()
		}
	}
}

// 错误事件
func (h *GameMsgHandler) OnError(err error) {
	gslog.Info("get error: %v", err)
}

// 发送协议
func (h *GameMsgHandler) Send(msgid uint32, msg protoreflect.ProtoMessage) error {
	d, e := proto.Marshal(msg)
	if e != nil {
		return e
	}
	return h.MsgHandler.Send(h.sess, msgid, d)
}

// 发送错误码
func (h *GameMsgHandler) SendError(err game_proto.ErrorId) error {
	return h.Send(uint32(game_proto.MsgErrorAck_Id), &game_proto.MsgErrorAck{
		Error: err,
	})
}

// todo 暂时在这里处理登录，有了登录服务器再说
func (h *GameMsgHandler) onPlayerLoginReq(sess gsnet.ISession, data []byte) error {
	var req game_proto.MsgAccountLoginGameReq
	err := proto.Unmarshal(data, &req)
	result := uint32(0)
	if err != nil {
		return err
	}

	var e game_proto.ErrorId
	if req.Account == "" {
		e = game_proto.ErrorId_INVALID_ACCOUNT
		gslog.Warn("Invalid empty account")
	} else if h.service.loginCheckMgr.checkAndAdd(string(req.Account)) {
		// todo 需要检测账号是否合法
		h.service.loginCheckMgr.remove(string(req.Account))
	} else {
		e = game_proto.ErrorId_ACCOUNT_IS_LOGGIN
		gslog.Warn("account %v is logining", req.Account)
	}

	if e != game_proto.ErrorId_NONE {
		return h.SendError(e)
	}

	var ack game_proto.MsgAccountLoginGameAck
	ack.Account = req.Account
	ack.Result = result
	if result == 0 {
		gslog.Info("player (account: %v, session: %v) login", req.Account, sess.GetId())
	}
	return h.Send(uint32(game_proto.MsgAccountLoginGameAck_Id), &ack)
}

// 处理进入游戏
func (h *GameMsgHandler) onPlayerEnterGameReq(sess gsnet.ISession, data []byte) error {
	var req game_proto.MsgPlayerEnterGameReq
	err := proto.Unmarshal(data, &req)
	if err != nil {
		return err
	}

	var p *SPlayer
	var e game_proto.ErrorId

	if req.Account == "" {
		e = game_proto.ErrorId_INVALID_ACCOUNT
		gslog.Warn("Invalid empty account")
	} else if !h.service.enterCheckMgr.checkAndAdd(req.Account) {
		e = game_proto.ErrorId_PLAYER_ENTERING_GAME
		gslog.Warn("already have same account entering")
	} else {
		// 检测请求的Account是否有其他goroutine在处理，保证同一时刻只有一个goroutine处理一个玩家进入
		// 先判断session中有没保存的数据
		pd := sess.GetData("player")
		// 有数据说明该会话已经存在，是客户端重复发送了进入消息
		if pd != nil {
			var o bool
			p, o = pd.(*SPlayer)
			if !o {
				e = game_proto.ErrorId_SESSION_INTERNAL_ERROR
				gslog.Error("session %v data %v convert failed", sess.GetId(), pd)
			}
			// 同一个玩家，重复发送进入游戏的消息
			if p.Account() == req.GetAccount() {
				e = game_proto.ErrorId_PLAYER_ENTER_GAME_REPEATED
				gslog.Warn("account %v send enter game repeated", req.Account)
			} else {
				e = game_proto.ErrorId_DIFFERENT_PLAYER_ENTER_SAME_SESSION
				gslog.Warn("different account enter with same session")
			}
		} else {
			var pid uint64
			p = h.service.playerMgr.GetByAccount(req.Account)
			// 该玩家正在游戏中，在另一个goroutine中
			if p != nil {
				pid = p.Id()
				// 断连等待结束，该函数只在这里使用
				p.WaitDisconnect()
				gslog.Warn("duplicate account %v, kicked another", req.Account)
			} else {
				pid = h.service.playerMgr.GetNextId()
			}
			// todo 读取数据库的逻辑（暂时不加）： 先读取redis缓存，没有则创建一个
			// 创建一个新Player
			p = NewSPlayer(req.Account, pid, sess)
			// 加入玩家管理器
			h.service.playerMgr.Add(p)
			p.SetToken(string(req.SessionToken))
			// 初始化坦克
			p.InitTank(common_data.PlayerTankInitData)
			tank := p.GetTank()
			// 随机坦克位置
			pos := utils.RandomPosInRect(common_data.PlayerTankInitRect)
			tank.SetPos(pos.X, pos.Y)
			// 把玩家加入游戏逻辑线程
			h.service.gameLogicThread.PlayerEnter(p.Id(), &playerData{tank: tank, sessHandler: h, pid: p.Id(), account: p.Account()})
			// 玩家进入
			p.Entered()
			// 把玩家对象设置到会话
			sess.SetData("player", p)
		}
		// 删除进入游戏状态
		h.service.enterCheckMgr.remove(req.Account)
	}

	if e != game_proto.ErrorId_NONE {
		return h.SendError(e)
	}

	return nil
}

// 处理退出游戏
func (h *GameMsgHandler) onPlayerExitGameReq(sess gsnet.ISession, data []byte) error {
	p, err := h.toPlayer(sess)
	if err != nil {
		return err
	}
	h.service.gameLogicThread.PlayerLeave(p.Id())
	h.service.playerMgr.Remove(p.Id())
	gslog.Info("player (account: %v, player_id: %v, session: %v) exit game", p.Account(), p.Id(), sess.GetId())
	var ack game_proto.MsgPlayerExitGameAck
	return h.Send(uint32(game_proto.MsgPlayerExitGameAck_Id), &ack)
}

// 处理改变坦克
func (h *GameMsgHandler) onPlayerChangeTankReq(sess gsnet.ISession, data []byte) error {
	p, err := h.toPlayer(sess)
	if err != nil {
		return err
	}

	var req game_proto.MsgPlayerChangeTankReq
	err = proto.Unmarshal(data, &req)
	if err != nil {
		return err
	}

	tankId := p.GetChangeTankId()
	req.TankId = tankId
	h.service.gameLogicThread.PushMsg(p.Id(), uint32(game_proto.MsgPlayerChangeTankReq_Id), &req)

	gslog.Info("player (account: %v, player_id: %v, session: %v) pushed change tank msg to game logic thread", p.Account(), p.Id(), sess.GetId())

	return nil
}

// 处理恢复坦克
func (h *GameMsgHandler) onPlayerRestoreTankReq(sess gsnet.ISession, data []byte) error {
	p, err := h.toPlayer(sess)
	if err != nil {
		return err
	}

	var req game_proto.MsgPlayerRestoreTankReq
	err = proto.Unmarshal(data, &req)
	if err != nil {
		return err
	}
	h.service.gameLogicThread.PushMsg(p.Id(), uint32(game_proto.MsgPlayerRestoreTankReq_Id), &req)

	gslog.Info("player (account: %v, player_id: %v, session: %v) pushed restore tank msg to game logic threadv", p.Account(), p.Id(), sess.GetId())

	return nil
}

// 坦克移动请求
func (h *GameMsgHandler) onPlayerTankMoveReq(sess gsnet.ISession, data []byte) error {
	p, err := h.toPlayer(sess)
	if err != nil {
		return err
	}

	var sync game_proto.MsgPlayerTankMoveReq
	err = proto.Unmarshal(data, &sync)
	if err != nil {
		return err
	}
	h.service.gameLogicThread.PushMsg(p.Id(), uint32(game_proto.MsgPlayerTankMoveReq_Id), &sync)

	return nil
}

// 内部函数，注册协议
func (s *GameMsgHandler) registerHandles() {
	s.RegisterHandle(uint32(game_proto.MsgAccountLoginGameReq_Id), s.onPlayerLoginReq)
	s.RegisterHandle(uint32(game_proto.MsgPlayerEnterGameReq_Id), s.onPlayerEnterGameReq)
	s.RegisterHandle(uint32(game_proto.MsgPlayerExitGameReq_Id), s.onPlayerExitGameReq)
	s.RegisterHandle(uint32(game_proto.MsgPlayerChangeTankReq_Id), s.onPlayerChangeTankReq)
	s.RegisterHandle(uint32(game_proto.MsgPlayerRestoreTankReq_Id), s.onPlayerRestoreTankReq)
	s.RegisterHandle(uint32(game_proto.MsgPlayerTankMoveSync_Id), s.onPlayerTankMoveReq)
}

// 会话转成玩家
func (h *GameMsgHandler) toPlayer(sess gsnet.ISession) (*SPlayer, error) {
	d := sess.GetData("player")
	if d == nil {
		return nil, errors.New("game_service: no invalid session")
	}
	p, o := d.(*SPlayer)
	if !o {
		return nil, errors.New("game_service: type cast to Player failed")
	}
	return p, nil
}

// 玩家断开后
func (h *GameMsgHandler) afterPlayerDisconnect(playerSess gsnet.ISession) {
	d := playerSess.GetData("player")
	if p, o := d.(*SPlayer); o {
		pid := p.Id()
		// 离开游戏逻辑线程
		h.service.gameLogicThread.PlayerLeave(pid)
		// 离开游戏
		p.Left(true)
		// 从管理器中删除Player
		h.service.playerMgr.Remove(pid)
		// 通知断开结束
		kicker := p.GetKicker()
		if kicker != nil {
			kicker.NotifyDisconnectEnd()
		}
		gslog.Info("player (account: %v, player_id: %v, session: %v) disconnect", p.Account(), pid, playerSess.GetId())
	}
}
