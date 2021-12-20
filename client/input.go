package main

import (
	"project_b/client_core"
	"project_b/common/object"

	"github.com/hajimehoshi/ebiten/v2"
)

type KeyCmdData struct {
	cmd  client_core.CmdCode
	args []interface{}
}

// 按下键位映射命令
var keyPressed2CmdMap = map[ebiten.Key]*KeyCmdData{
	ebiten.KeyA: {cmd: client_core.CMD_MOVE, args: []interface{}{object.DirLeft}},
	ebiten.KeyD: {cmd: client_core.CMD_MOVE, args: []interface{}{object.DirRight}},
	ebiten.KeyW: {cmd: client_core.CMD_MOVE, args: []interface{}{object.DirUp}},
	ebiten.KeyS: {cmd: client_core.CMD_MOVE, args: []interface{}{object.DirDown}},
}

// 释放键位映射命令
var keyReleased2CmdMap = map[ebiten.Key]client_core.CmdCode{
	ebiten.KeyA: client_core.CMD_STOP_MOVE,
	ebiten.KeyD: client_core.CMD_STOP_MOVE,
	ebiten.KeyW: client_core.CMD_STOP_MOVE,
	ebiten.KeyS: client_core.CMD_STOP_MOVE,
	ebiten.Key1: client_core.CMD_CHANGE_TANK,
	ebiten.Key2: client_core.CMD_RESTORE_TANK,
}
