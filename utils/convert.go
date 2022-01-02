package utils

import (
	"project_b/common/object"
	"project_b/common_data"
	"project_b/game_proto"
)

func TankObj2ProtoInfo(obj *object.Tank, protoInfo *game_proto.TankInfo) {
	protoInfo.Id = obj.Id()
	protoInfo.InstId = obj.InstId()
	protoInfo.Level = obj.Level()
	x, y := obj.Pos()
	if protoInfo.CurrPos == nil {
		protoInfo.CurrPos = &game_proto.Pos{}
	}
	protoInfo.CurrPos.X = x
	protoInfo.CurrPos.Y = y
	protoInfo.Direction = int32(obj.Dir())
	protoInfo.CurrSpeed = obj.CurrentSpeed()
}

func TankProtoInfo2Obj(protoInfo *game_proto.TankInfo, obj *object.Tank) {
	td := common_data.TankConfigData[protoInfo.Id]
	obj.Init(protoInfo.InstId, td)
	obj.SetLevel(protoInfo.Level)
	obj.SetPos(protoInfo.CurrPos.X, protoInfo.CurrPos.Y)
	obj.SetDir(object.Direction(protoInfo.Direction))
	obj.SetCurrentSpeed(protoInfo.CurrSpeed)
}
