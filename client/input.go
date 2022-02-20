package main

import (
	"project_b/client/core"
	"project_b/common/object"

	"github.com/hajimehoshi/ebiten/v2"
)

type KeyCmdData struct {
	cmd  core.CmdCode
	args []interface{}
}

// 按下键位映射命令
var keyPressed2CmdMap = map[ebiten.Key]*KeyCmdData{
	ebiten.KeyA: {cmd: core.CMD_MOVE, args: []interface{}{object.DirLeft}},
	ebiten.KeyD: {cmd: core.CMD_MOVE, args: []interface{}{object.DirRight}},
	ebiten.KeyW: {cmd: core.CMD_MOVE, args: []interface{}{object.DirUp}},
	ebiten.KeyS: {cmd: core.CMD_MOVE, args: []interface{}{object.DirDown}},
}

// 释放键位映射命令
var keyReleased2CmdMap = map[ebiten.Key]core.CmdCode{
	ebiten.KeyA: core.CMD_STOP_MOVE,
	ebiten.KeyD: core.CMD_STOP_MOVE,
	ebiten.KeyW: core.CMD_STOP_MOVE,
	ebiten.KeyS: core.CMD_STOP_MOVE,
	ebiten.Key1: core.CMD_CHANGE_TANK,
	ebiten.Key2: core.CMD_RESTORE_TANK,
}
