package main

import (
	core "project_b/client_core"
	"project_b/common/object"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
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
	args []any
}

// 按下键位映射命令
var keyPressed2CmdMap = map[ebiten.Key]*KeyCmdData{
	ebiten.KeyA:        {cmd: core.CMD_MOVE, args: []any{object.DirLeft}},
	ebiten.KeyD:        {cmd: core.CMD_MOVE, args: []any{object.DirRight}},
	ebiten.KeyW:        {cmd: core.CMD_MOVE, args: []any{object.DirUp}},
	ebiten.KeyS:        {cmd: core.CMD_MOVE, args: []any{object.DirDown}},
	ebiten.KeyJ:        {cmd: core.CMD_FIRE, args: []any{1}},
	ebiten.KeyI:        {cmd: core.CMD_FIRE, args: []any{2}},
	ebiten.KeyK:        {cmd: core.CMD_ROTATE, args: []any{1}},
	ebiten.KeyL:        {cmd: core.CMD_ROTATE, args: []any{-1}},
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
	ebiten.KeyC: core.CMD_CHANGE_TANK,
	ebiten.KeyR: core.CMD_RESTORE_TANK,
	ebiten.Key1: core.CMD_RELEASE_SMALL_BALL,
	ebiten.KeyU: core.CMD_SHIELD,
}

type ComboKeyIndex struct {
	otherKey ebiten.Key
	index    int32
}

var key2ComboKeyIndex = map[ebiten.Key][]ComboKeyIndex{
	ebiten.KeyA: {
		{otherKey: ebiten.KeyW, index: 0},
		{otherKey: ebiten.KeyS, index: 1},
	},
	ebiten.KeyD: {
		{otherKey: ebiten.KeyW, index: 2},
		{otherKey: ebiten.KeyS, index: 3},
	},
	ebiten.KeyW: {
		{otherKey: ebiten.KeyA, index: 0},
		{otherKey: ebiten.KeyD, index: 2},
	},
	ebiten.KeyS: {
		{otherKey: ebiten.KeyA, index: 1},
		{otherKey: ebiten.KeyD, index: 3},
	},
}

var keyIndex2CmdMap = map[int32]*KeyCmdData{
	0: {cmd: core.CMD_MOVE, args: []any{object.DirLeftUp}},
	1: {cmd: core.CMD_MOVE, args: []any{object.DirLeftDown}},
	2: {cmd: core.CMD_MOVE, args: []any{object.DirRightUp}},
	3: {cmd: core.CMD_MOVE, args: []any{object.DirRightDown}},
}

// 輸入管理器
type InputMgr struct {
	cmdMgr      *core.CmdHandleManager
	pressedKeys []ebiten.Key
	keyPressMap map[ebiten.Key]struct{}
}

// 創建輸入管理器
func NewInputMgr(cmdMgr *core.CmdHandleManager) *InputMgr {
	return &InputMgr{
		cmdMgr:      cmdMgr,
		keyPressMap: make(map[ebiten.Key]struct{}),
	}
}

// 處理輸入
func (im *InputMgr) HandleInput() {
	var (
		cmd core.CmdCode
		o   bool
	)
	// 處理鍵釋放
	for k, _ := range im.keyPressMap {
		if inpututil.IsKeyJustReleased(k) {
			cmd, o = keyReleased2CmdMap[k]
			if !o {
				continue
			}
			im.cmdMgr.Handle(cmd)
			delete(im.keyPressMap, k)
			log.Debug("key %v released", k)
		}
	}

	clear(im.keyPressMap)

	// 處理鍵按下
	im.pressedKeys = inpututil.AppendPressedKeys(im.pressedKeys[:0])
	if len(im.pressedKeys) > 0 {
		log.Debug("pressed key list %v", im.pressedKeys)
	}
	for _, pk := range im.pressedKeys {
		im.keyPressMap[pk] = struct{}{}
	}

	var keyUsed []ebiten.Key
	var kiList []ComboKeyIndex
	var cmdData *KeyCmdData
	for k, _ := range im.keyPressMap {
		var used bool
		for i := 0; i < len(keyUsed); i++ {
			if k == keyUsed[i] {
				used = true
				break
			}
		}
		if used {
			continue
		}

		// 先處理組合key
		kiList, o = key2ComboKeyIndex[k]
		if o {
			var hasCombo bool
			for _, ki := range kiList {
				// 有組合key
				if _, o = im.keyPressMap[ki.otherKey]; o {
					cmdData, o = keyIndex2CmdMap[ki.index]
					if o && cmdData != nil {
						im.cmdMgr.Handle(cmdData.cmd, cmdData.args...)
					}
					keyUsed = append(keyUsed, ki.otherKey)
					hasCombo = true
					break
				}
			}
			if hasCombo {
				continue
			}
		}

		cmdData, o = keyPressed2CmdMap[k]
		if o && cmdData != nil {
			im.cmdMgr.Handle(cmdData.cmd, cmdData.args...)
		}
	}
}
