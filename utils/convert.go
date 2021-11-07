package utils

import (
	"project_b/common/object"
	"project_b/common_data"
	"project_b/game_proto"
)

func TankInfo2ProtoInfo(info *object.Tank, protoInfo *game_proto.TankInfo) {
	protoInfo.Id = info.Id()
	protoInfo.Level = info.Level()
	x, y := info.Pos()
	if protoInfo.CurrPos == nil {
		protoInfo.CurrPos = &game_proto.Pos{}
	}
	protoInfo.CurrPos.X = x
	protoInfo.CurrPos.Y = y
	protoInfo.Direction = int32(info.Dir())
	protoInfo.CurrSpeed = info.CurrentSpeed()
}

func TankProtoInfo2Info(protoInfo *game_proto.TankInfo) *object.Tank {
	td := common_data.TankConfigData[protoInfo.Id]
	tankObj := object.NewTank(td)
	tankObj.SetLevel(protoInfo.Level)
	tankObj.SetPos(protoInfo.CurrPos.X, protoInfo.CurrPos.Y)
	tankObj.SetDir(object.Direction(protoInfo.Direction))
	tankObj.SetCurrentSpeed(protoInfo.CurrSpeed)
	return tankObj
}
