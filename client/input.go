package main

import (
	core "project_b/client_core"
	"project_b/common/object"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	CMD_CAMERA_UP     = 1000
	CMD_CAMERA_DOWN   = 1001
	CMD_CAMERA_LEFT   = 1002
	CMD_CAMERA_RIGHT  = 1003
	CMD_CAMERA_HEIGHT = 1100
)

type KeyCmdData struct {
	cmd  core.CmdCode
	args []interface{}
}

// 按下键位映射命令
var keyPressed2CmdMap = map[ebiten.Key]*KeyCmdData{
	ebiten.KeyA:        {cmd: core.CMD_MOVE, args: []any{object.DirLeft}},
	ebiten.KeyD:        {cmd: core.CMD_MOVE, args: []any{object.DirRight}},
	ebiten.KeyW:        {cmd: core.CMD_MOVE, args: []any{object.DirUp}},
	ebiten.KeyS:        {cmd: core.CMD_MOVE, args: []any{object.DirDown}},
	ebiten.KeyUp:       {cmd: CMD_CAMERA_UP, args: []any{}},
	ebiten.KeyDown:     {cmd: CMD_CAMERA_DOWN, args: []any{}},
	ebiten.KeyLeft:     {cmd: CMD_CAMERA_LEFT, args: []any{}},
	ebiten.KeyRight:    {cmd: CMD_CAMERA_RIGHT, args: []any{}},
	ebiten.KeyPageUp:   {cmd: CMD_CAMERA_HEIGHT, args: []any{10}},
	ebiten.KeyPageDown: {cmd: CMD_CAMERA_HEIGHT, args: []any{-10}},
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
