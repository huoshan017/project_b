package main

import (
	"project_b/common/object"

	"github.com/hajimehoshi/ebiten/v2"
)

type KeyCmdData struct {
	cmd  CmdCode
	args []interface{}
}

// 按下键位映射命令
var keyPressed2CmdMap = map[ebiten.Key]*KeyCmdData{
	ebiten.KeyA: {cmd: CMD_MOVE, args: []interface{}{object.DirLeft}},
	ebiten.KeyD: {cmd: CMD_MOVE, args: []interface{}{object.DirRight}},
	ebiten.KeyW: {cmd: CMD_MOVE, args: []interface{}{object.DirUp}},
	ebiten.KeyS: {cmd: CMD_MOVE, args: []interface{}{object.DirDown}},
}

// 释放键位映射命令
var keyReleased2CmdMap = map[ebiten.Key]CmdCode{
	ebiten.KeyA: CMD_STOP_MOVE,
	ebiten.KeyD: CMD_STOP_MOVE,
	ebiten.KeyW: CMD_STOP_MOVE,
	ebiten.KeyS: CMD_STOP_MOVE,
	ebiten.Key1: CMD_CHANGE_TANK,
	ebiten.Key2: CMD_RESTORE_TANK,
}
