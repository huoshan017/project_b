package main

//#include "./net_client_capi.h"
import "C"
import (
	"project_b/client_core"
	"project_b/common/base"
	"project_b/game_proto"
)

//export net_client_new
func net_client_new(address *C.char) C.net_handle_t {
	c := client_core.CreateNetClient(C.GoString(address))
	if c == nil {
		return -1
	}
	id := NewObjectId(c)
	return C.net_handle_t(id)
}

//export net_client_delete
func net_client_delete(h C.net_handle_t) {
	ObjectFree(ObjectId(h))
}

//export net_client_send_login_req
func net_client_send_login_req(h C.net_handle_t, account *C.char, password *C.char) C.int {
	client := ObjectGet(ObjectId(h)).(*client_core.NetClient)
	if client == nil {
		return -1
	}
	err := client.SendLoginReq(C.GoString(account), C.GoString(password))
	if err != nil {
		return -2
	}
	return C.int(0)
}

//export net_client_send_game_enter_req
func net_client_send_game_enter_req(h C.net_handle_t, account *C.char, session_token *C.char) C.int {
	client := ObjectGet(ObjectId(h)).(*client_core.NetClient)
	if client == nil {
		return -1
	}
	err := client.SendEnterGameReq(C.GoString(account), C.GoString(session_token))
	if err != nil {
		return -2
	}
	return 0
}

//export net_client_send_time_sync_req
func net_client_send_time_sync_req(h C.net_handle_t) C.int {
	client := ObjectGet(ObjectId(h)).(*client_core.NetClient)
	if client == nil {
		return -1
	}
	err := client.SendTimeSyncReq()
	if err != nil {
		return -2
	}
	return 0
}

//export net_client_send_tank_move_req
func net_client_send_tank_move_req(h C.net_handle_t /*dir*/, orientation C.int) C.int {
	client := ObjectGet(ObjectId(h)).(*client_core.NetClient)
	if client == nil {
		return -1
	}
	err := client.SendTankMoveReq( /*object.Direction(dir)*/ int32(orientation))
	if err != nil {
		return -2
	}
	return 0
}

//export net_client_send_tank_update_pos_req
func net_client_send_tank_update_pos_req(h C.net_handle_t, state C.int, x, y /*dir*/, orientation, speed C.int) C.int {
	client := ObjectGet(ObjectId(h)).(*client_core.NetClient)
	if client == nil {
		return -1
	}
	err := client.SendTankUpdatePosReq(game_proto.MovementState(state), base.Pos{X: int32(x), Y: int32(y)} /*object.Direction(dir)*/, int32(orientation), int32(speed))
	if err != nil {
		return -2
	}
	return 0
}

//export net_client_send_tank_stop_move_req
func net_client_send_tank_stop_move_req(h C.net_handle_t) C.int {
	client := ObjectGet(ObjectId(h)).(*client_core.NetClient)
	if client == nil {
		return -1
	}
	err := client.SendTankStopMoveReq()
	if err != nil {
		return -2
	}
	return 0
}

//export net_client_send_tank_change_req
func net_client_send_tank_change_req(h C.net_handle_t) C.int {
	client := ObjectGet(ObjectId(h)).(*client_core.NetClient)
	if client == nil {
		return -1
	}
	err := client.SendTankChangeReq()
	if err != nil {
		return -2
	}
	return 0
}

//export net_client_send_tank_restore_req
func net_client_send_tank_restore_req(h C.net_handle_t) C.int {
	client := ObjectGet(ObjectId(h)).(*client_core.NetClient)
	if client == nil {
		return -1
	}
	err := client.SendTankRestoreReq()
	if err != nil {
		return -2
	}
	return 0
}

func main() {}
