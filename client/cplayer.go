package main

import (
	"project_b/common"
	"project_b/game_proto"
	"project_b/utils"
)

type CPlayer struct {
	common.Player
	net *NetClient
}

func NewCPlayer(acc string, id uint64, net *NetClient) *CPlayer {
	p := &CPlayer{
		Player: *common.NewPlayer(acc, id),
		net:    net,
	}
	return p
}

func (p *CPlayer) InitTankFromProto(tankProtoInfo *game_proto.TankInfo) {
	tank := utils.TankProtoInfo2Info(p.Id(), tankProtoInfo)
	p.SetTank(tank)
}
