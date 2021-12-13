package main

import (
	"project_b/common/object"
	custom_time "project_b/common/time"
	"project_b/game_proto"

	"google.golang.org/protobuf/proto"

	"github.com/huoshan017/gsnet"
)

type NetClient struct {
	gsnet.MsgClient
	serverAddress string
}

func CreateNetClient(serverAddress string) *NetClient {
	client := &NetClient{}
	client.MsgClient = *gsnet.NewMsgClient()
	client.serverAddress = serverAddress
	return client
}

func (c *NetClient) SendLoginReq(account, password string) error {
	err := c.Connect(c.serverAddress)
	if err != nil {
		return err
	}
	var req game_proto.MsgAccountLoginGameReq
	req.Account = account
	req.Password = password
	data, err := proto.Marshal(&req)
	if err != nil {
		return err
	}
	return c.Send(uint32(game_proto.MsgAccountLoginGameReq_Id), data)
}

func (c *NetClient) SendEnterGameReq(account, sessionToken string) error {
	var req game_proto.MsgPlayerEnterGameReq
	req.Account = account
	req.SessionToken = sessionToken
	data, err := proto.Marshal(&req)
	if err != nil {
		return err
	}
	return c.Send(uint32(game_proto.MsgPlayerEnterGameReq_Id), data)
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
	var data []byte
	data, err = proto.Marshal(&req)
	if err != nil {
		return err
	}
	SetSyncSendTime(now)
	return c.Send(uint32(game_proto.MsgTimeSyncReq_Id), data)
}

func (c *NetClient) SendTankMoveReq(dir object.Direction) error {
	var req game_proto.MsgPlayerTankMoveReq
	req.MoveInfo = &game_proto.TankMoveInfo{}
	req.MoveInfo.Direction = int32(dir)
	data, err := proto.Marshal(&req)
	if err != nil {
		return err
	}
	return c.Send(uint32(game_proto.MsgPlayerTankMoveReq_Id), data)
}

func (c *NetClient) SendTankStopMoveReq() error {
	var req game_proto.MsgPlayerTankStopMoveReq
	data, err := proto.Marshal(&req)
	if err != nil {
		return err
	}
	return c.Send(uint32(game_proto.MsgPlayerTankStopMoveReq_Id), data)
}

func (c *NetClient) SendTankChangeReq() error {
	var req game_proto.MsgPlayerChangeTankReq
	data, err := proto.Marshal(&req)
	if err != nil {
		return err
	}
	return c.Send(uint32(game_proto.MsgPlayerChangeTankReq_Id), data)
}

func (c *NetClient) SendTankRestoreReq() error {
	var req game_proto.MsgPlayerRestoreTankReq
	data, err := proto.Marshal(&req)
	if err != nil {
		return err
	}
	return c.Send(uint32(game_proto.MsgPlayerRestoreTankReq_Id), data)
}
