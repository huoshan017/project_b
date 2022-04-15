package client_core

import (
	"project_b/common/object"
	custom_time "project_b/common/time"
	"project_b/game_proto"
	"time"

	gsnet_common "github.com/huoshan017/gsnet/common"
	gsnet_msg "github.com/huoshan017/gsnet/msg"
)

type NetClient struct {
	msgClient     *gsnet_msg.MsgClient
	serverAddress string
}

func CreateNetClient(serverAddress string, options ...gsnet_common.Option) *NetClient {
	client := &NetClient{
		msgClient:     gsnet_msg.NewPBMsgClient(gsnet_msg.CreateIdMsgMapperWith(game_proto.Id2MsgMapOnClient), options...),
		serverAddress: serverAddress,
	}
	return client
}

func (c *NetClient) RegisterHandle(msgId gsnet_msg.MsgIdType, handle func(*gsnet_msg.MsgSession, interface{}) error) {
	c.msgClient.RegisterMsgHandle(msgId, handle)
}

func (c *NetClient) SetConnectHandle(handle func(*gsnet_msg.MsgSession)) {
	c.msgClient.SetConnectHandle(handle)
}

func (c *NetClient) SetDisconnectHandle(handle func(*gsnet_msg.MsgSession, error)) {
	c.msgClient.SetDisconnectHandle(handle)
}

func (c *NetClient) SetTickHandle(handle func(*gsnet_msg.MsgSession, time.Duration)) {
	c.msgClient.SetTickHandle(handle)
}

func (c *NetClient) SetErrorHandle(handle func(error)) {
	c.msgClient.SetErrorHandle(handle)
}

func (c *NetClient) IsConnecting() bool {
	return c.msgClient.IsConnecting()
}

func (c *NetClient) IsConnected() bool {
	return c.msgClient.IsConnected()
}

func (c *NetClient) IsDisconnecting() bool {
	return c.msgClient.IsDisconnecting()
}

func (c *NetClient) IsDisconnected() bool {
	return c.msgClient.IsDisconnected()
}

func (c *NetClient) Update() error {
	return c.msgClient.Update()
}

func (c *NetClient) SendLoginReq(account, password string) error {
	if c.msgClient.IsConnected() {
		gslog.Warn("already connected")
		return nil
	}
	err := c.msgClient.Connect(c.serverAddress)
	if err != nil {
		return err
	}
	var req game_proto.MsgAccountLoginGameReq
	req.Account = account
	req.Password = password
	return c.msgClient.Send(gsnet_msg.MsgIdType(game_proto.MsgAccountLoginGameReq_Id), &req)
}

func (c *NetClient) SendEnterGameReq(account, sessionToken string) error {
	var req game_proto.MsgPlayerEnterGameReq
	req.Account = account
	req.SessionToken = sessionToken
	return c.msgClient.Send(gsnet_msg.MsgIdType(game_proto.MsgPlayerEnterGameReq_Id), &req)
}

func (c *NetClient) SendTimeSyncReq() error {
	var req game_proto.MsgTimeSyncReq
	// 当前时间序列化
	now := custom_time.Now()
	td, err := now.MarshalBinary()
	if err != nil {
		return err
	}
	req.ClientTime = td
	SetSyncSendTime(now)
	return c.msgClient.Send(gsnet_msg.MsgIdType(game_proto.MsgTimeSyncReq_Id), &req)
}

func (c *NetClient) SendTankMoveReq(dir object.Direction) error {
	var req game_proto.MsgPlayerTankMoveReq
	req.MoveInfo = &game_proto.TankMoveInfo{}
	req.MoveInfo.Direction = int32(dir)
	return c.msgClient.Send(gsnet_msg.MsgIdType(game_proto.MsgPlayerTankMoveReq_Id), &req)
}

func (c *NetClient) SendTankUpdatePosReq(state game_proto.MovementState, pos object.Pos, dir object.Direction, speed int32) error {
	var req game_proto.MsgPlayerTankUpdatePosReq
	req.State = game_proto.MovementState(state)
	req.MoveInfo = &game_proto.TankMoveInfo{}
	req.MoveInfo.CurrPos = &game_proto.Pos{X: pos.X, Y: pos.Y}
	req.MoveInfo.CurrSpeed = speed
	req.MoveInfo.Direction = int32(dir)
	return c.msgClient.Send(gsnet_msg.MsgIdType(game_proto.MsgPlayerTankUpdatePosReq_Id), &req)
}

func (c *NetClient) SendTankStopMoveReq() error {
	var req game_proto.MsgPlayerTankStopMoveReq
	return c.msgClient.Send(gsnet_msg.MsgIdType(game_proto.MsgPlayerTankStopMoveReq_Id), &req)
}

func (c *NetClient) SendTankChangeReq() error {
	var req game_proto.MsgPlayerChangeTankReq
	return c.msgClient.Send(gsnet_msg.MsgIdType(game_proto.MsgPlayerChangeTankReq_Id), &req)
}

func (c *NetClient) SendTankRestoreReq() error {
	var req game_proto.MsgPlayerRestoreTankReq
	return c.msgClient.Send(gsnet_msg.MsgIdType(game_proto.MsgPlayerRestoreTankReq_Id), &req)
}
