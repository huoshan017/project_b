package main

import (
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

func (c *NetClient) SendTankMoveReq() error {
	var req game_proto.MsgPlayerTankMoveReq
	req.MoveInfo = &game_proto.TankMoveInfo{}
	data, err := proto.Marshal(&req)
	if err != nil {
		return err
	}
	return c.Send(uint32(game_proto.MsgPlayerTankMoveReq_Id), data)
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
