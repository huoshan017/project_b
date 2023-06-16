package main

import (
	custom_time "project_b/common/time"
	"time"

	gsnet_msg "github.com/huoshan017/gsnet/msg"
)

type GameMsgHandler struct {
	service               *GameService
	sess                  *gsnet_msg.MsgSession
	lastCheckDiscDuration custom_time.Duration
	msgid2HandleMap       map[gsnet_msg.MsgIdType]func(*gsnet_msg.MsgSession, any) error
}

func CreateGameMsgHandler(owner *GameService) *GameMsgHandler {
	h := &GameMsgHandler{}
	h.service = owner
	return h
}

// 连接事件
func (h *GameMsgHandler) OnConnected(sess *gsnet_msg.MsgSession) {
	gslog.Info("new session %v connected", sess.GetId())
}

func (h *GameMsgHandler) OnReady(sess *gsnet_msg.MsgSession) {
	// 连接后把会话缓存起来
	h.sess = sess
	gslog.Info("session %v ready", sess.GetId())
}

// 断开事件
func (h *GameMsgHandler) OnDisconnected(sess *gsnet_msg.MsgSession, err error) {
	if h.sess != nil && sess != h.sess {
		panic("sess must same to OnConnect")
	}
	h.afterPlayerDisconnect(sess)
	gslog.Info("session %v disconnected", sess.GetId())
	h.sess = nil
}

// 定时器事件
func (h *GameMsgHandler) OnTick(sess *gsnet_msg.MsgSession, tick time.Duration) {
	if sess != h.sess {
		panic("sess must same to OnConnect")
	}
	h.lastCheckDiscDuration += custom_time.Duration(tick)
	// 0.5秒检测一次
	if h.lastCheckDiscDuration >= custom_time.Duration(time.Millisecond)*500 {
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

func (h *GameMsgHandler) OnMsgHandle(sess *gsnet_msg.MsgSession, msgid gsnet_msg.MsgIdType, msgobj any) error {
	return h.getMsgId2HandleMap()[msgid](sess, msgobj)
}
