package game_proto

import (
	reflect "reflect"

	"github.com/huoshan017/gsnet/msg"
)

var (
	// 客户端发出服务器处理
	Id2MsgMapOnServer = map[msg.MsgIdType]reflect.Type{

		msg.MsgIdType(MsgAccountLoginGameReq_Id):    reflect.TypeOf(&MsgAccountLoginGameReq{}),
		msg.MsgIdType(MsgTimeSyncReq_Id):            reflect.TypeOf(&MsgTimeSyncReq{}),
		msg.MsgIdType(MsgPlayerEnterGameReq_Id):     reflect.TypeOf(&MsgPlayerEnterGameReq{}),
		msg.MsgIdType(MsgPlayerExitGameReq_Id):      reflect.TypeOf(&MsgPlayerExitGameReq{}),
		msg.MsgIdType(MsgPlayerBasicInfoReq_Id):     reflect.TypeOf(&MsgPlayerBasicInfoReq{}),
		msg.MsgIdType(MsgPlayerTankMoveReq_Id):      reflect.TypeOf(&MsgPlayerTankMoveReq{}),
		msg.MsgIdType(MsgPlayerTankUpdatePosReq_Id): reflect.TypeOf(&MsgPlayerTankUpdatePosReq{}),
		msg.MsgIdType(MsgPlayerTankStopMoveReq_Id):  reflect.TypeOf(&MsgPlayerTankStopMoveReq{}),
		msg.MsgIdType(MsgPlayerChangeTankReq_Id):    reflect.TypeOf(&MsgPlayerChangeTankReq{}),
		msg.MsgIdType(MsgPlayerRestoreTankReq_Id):   reflect.TypeOf(&MsgPlayerRestoreTankReq{}),
	}

	Id2MsgMapOnClient = map[msg.MsgIdType]reflect.Type{
		// 服务器发出客户端接收
		msg.MsgIdType(MsgAccountLoginGameAck_Id):     reflect.TypeOf(&MsgAccountLoginGameAck{}),
		msg.MsgIdType(MsgTimeSyncAck_Id):             reflect.TypeOf(&MsgTimeSyncAck{}),
		msg.MsgIdType(MsgPlayerEnterGameAck_Id):      reflect.TypeOf(&MsgPlayerEnterGameAck{}),
		msg.MsgIdType(MsgPlayerExitGameAck_Id):       reflect.TypeOf(&MsgPlayerExitGameAck{}),
		msg.MsgIdType(MsgPlayerBasicInfoAck_Id):      reflect.TypeOf(&MsgPlayerBasicInfoAck{}),
		msg.MsgIdType(MsgPlayerTankMoveAck_Id):       reflect.TypeOf(&MsgPlayerTankMoveAck{}),
		msg.MsgIdType(MsgPlayerTankMoveSync_Id):      reflect.TypeOf(&MsgPlayerTankMoveSync{}),
		msg.MsgIdType(MsgPlayerTankUpdatePosAck_Id):  reflect.TypeOf(&MsgPlayerTankUpdatePosAck{}),
		msg.MsgIdType(MsgPlayerTankUpdatePosSync_Id): reflect.TypeOf(&MsgPlayerTankUpdatePosSync{}),
		msg.MsgIdType(MsgPlayerTankStopMoveAck_Id):   reflect.TypeOf(&MsgPlayerTankStopMoveAck{}),
		msg.MsgIdType(MsgPlayerTankStopMoveSync_Id):  reflect.TypeOf(&MsgPlayerTankStopMoveSync{}),
		msg.MsgIdType(MsgPlayerChangeTankAck_Id):     reflect.TypeOf(&MsgPlayerChangeTankAck{}),
		msg.MsgIdType(MsgPlayerChangeTankSync_Id):    reflect.TypeOf(&MsgPlayerChangeTankSync{}),
		msg.MsgIdType(MsgPlayerRestoreTankAck_Id):    reflect.TypeOf(&MsgPlayerRestoreTankAck{}),
		msg.MsgIdType(MsgPlayerRestoreTankSync_Id):   reflect.TypeOf(&MsgPlayerRestoreTankSync{}),
	}
)
