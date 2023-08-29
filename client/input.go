package main

import (
	"project_b/client_base"
	"project_b/common/object"
	"project_b/core"

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

// 按下键位映射命令
var keyPressed2CmdMap = map[ebiten.Key]*core.CmdData{
	ebiten.KeyA:        core.NewCmdData(core.CMD_TANK_MOVE, []int64{int64(object.DirLeft)}),
	ebiten.KeyD:        core.NewCmdData(core.CMD_TANK_MOVE, []int64{int64(object.DirRight)}),
	ebiten.KeyW:        core.NewCmdData(core.CMD_TANK_MOVE, []int64{int64(object.DirUp)}),
	ebiten.KeyS:        core.NewCmdData(core.CMD_TANK_MOVE, []int64{int64(object.DirDown)}),
	ebiten.KeyJ:        core.NewCmdData(core.CMD_TANK_FIRE, []int64{}),
	ebiten.KeyUp:       core.NewCmdData(CMD_CAMERA_UP, []int64{10}),
	ebiten.KeyDown:     core.NewCmdData(CMD_CAMERA_DOWN, []int64{-10}),
	ebiten.KeyLeft:     core.NewCmdData(CMD_CAMERA_LEFT, []int64{-10}),
	ebiten.KeyRight:    core.NewCmdData(CMD_CAMERA_RIGHT, []int64{10}),
	ebiten.KeyPageUp:   core.NewCmdData(CMD_CAMERA_HEIGHT, []int64{10}),
	ebiten.KeyPageDown: core.NewCmdData(CMD_CAMERA_HEIGHT, []int64{-10}),
}

// 释放键位映射命令
var keyReleased2CmdMap = map[ebiten.Key]*core.CmdData{
	ebiten.KeyA:     core.NewCmdData(core.CMD_TANK_STOP, nil),
	ebiten.KeyD:     core.NewCmdData(core.CMD_TANK_STOP, nil),
	ebiten.KeyW:     core.NewCmdData(core.CMD_TANK_STOP, nil),
	ebiten.KeyS:     core.NewCmdData(core.CMD_TANK_STOP, nil),
	ebiten.Key1:     core.NewCmdData(core.CMD_TANK_ADD_SHELL, []int64{1}),
	ebiten.Key2:     core.NewCmdData(core.CMD_TANK_ADD_SHELL, []int64{2}),
	ebiten.Key3:     core.NewCmdData(core.CMD_TANK_ADD_SHELL, []int64{3}),
	ebiten.Key4:     core.NewCmdData(core.CMD_TANK_ADD_SHELL, []int64{4}),
	ebiten.Key5:     core.NewCmdData(core.CMD_TANK_ADD_SHELL, []int64{5}),
	ebiten.Key6:     core.NewCmdData(core.CMD_TANK_ADD_SHELL, []int64{6}),
	ebiten.Key7:     core.NewCmdData(core.CMD_TANK_ADD_SHELL, []int64{7}),
	ebiten.Key8:     core.NewCmdData(core.CMD_TANK_ADD_SHELL, []int64{8}),
	ebiten.Key9:     core.NewCmdData(core.CMD_TANK_ADD_SHELL, []int64{9}),
	ebiten.KeyQ:     core.NewCmdData(core.CMD_TANK_SWITCH_SHELL, []int64{}),
	ebiten.KeyC:     core.NewCmdData(core.CMD_TANK_CHANGE, []int64{}),
	ebiten.KeyR:     core.NewCmdData(core.CMD_TANK_RESTORE, []int64{}),
	ebiten.KeyEqual: core.NewCmdData(core.CMD_RELEASE_SMALL_BALL, []int64{}),
	ebiten.KeyU:     core.NewCmdData(core.CMD_TANK_SHIELD, []int64{}),
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

var keyIndex2CmdMap = map[int32]*core.CmdData{
	0: core.NewCmdData(core.CMD_TANK_MOVE, []int64{int64(object.DirLeftUp)}),
	1: core.NewCmdData(core.CMD_TANK_MOVE, []int64{int64(object.DirLeftDown)}),
	2: core.NewCmdData(core.CMD_TANK_MOVE, []int64{int64(object.DirRightUp)}),
	3: core.NewCmdData(core.CMD_TANK_MOVE, []int64{int64(object.DirRightDown)}),
}

// 輸入管理器
type InputMgr struct {
	game          client_base.IGame
	gameCore      *core.GameCore
	pressedKeys   []ebiten.Key
	releasedKeys  []ebiten.Key
	keyPressMap   map[ebiten.Key]struct{}
	cmd2KeyHandle map[core.CmdCode]func([]int64)
	wheelHandle   func(xoffset, yoffset float64)
}

// 創建輸入管理器
func NewInputMgr(game client_base.IGame, gameCore *core.GameCore) *InputMgr {
	return &InputMgr{
		game:          game,
		gameCore:      gameCore,
		keyPressMap:   make(map[ebiten.Key]struct{}),
		cmd2KeyHandle: make(map[core.CmdCode]func([]int64)),
	}
}

// 添加按鍵處理器
func (im *InputMgr) AddKeyHandle(cc core.CmdCode, handle func([]int64)) {
	im.cmd2KeyHandle[cc] = handle
}

// 設置滾輪滾動處理器
func (im *InputMgr) SetWheelHandle(handle func(xoffset, yoffset float64)) {
	im.wheelHandle = handle
}

// 處理輸入
func (im *InputMgr) HandleInput() {
	var (
		cmdData *core.CmdData
		o       bool
	)
	// 處理鍵釋放
	/*for k := range im.keyPressMap {
		if inpututil.IsKeyJustReleased(k) {
			cmdData, o = keyReleased2CmdMap[k]
			if !o {
				continue
			}
			im.inst.PushFrame(im.inst.GetFrame(), im.game.GetGameData().MyId, cmdData.Cmd(), cmdData.Args())
			delete(im.keyPressMap, k)		}
	}*/

	if len(im.releasedKeys) > 0 {
		clear(im.releasedKeys)
	}
	im.releasedKeys = inpututil.AppendJustReleasedKeys(im.releasedKeys[:0])
	for _, k := range im.releasedKeys {
		if cmdData, o := keyReleased2CmdMap[k]; o {
			im.gameCore.PushSyncPlayerCmd(im.game.GetGameData().MyId, cmdData)
		}
	}

	clear(im.keyPressMap)

	// 處理鍵按下
	if len(im.pressedKeys) > 0 {
		clear(im.pressedKeys)
	}
	im.pressedKeys = inpututil.AppendPressedKeys(im.pressedKeys[:0])
	for _, pk := range im.pressedKeys {
		im.keyPressMap[pk] = struct{}{}
	}

	var keyUsed []ebiten.Key
	var kiList []ComboKeyIndex
	for k := range im.keyPressMap {
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
						im.gameCore.PushSyncPlayerCmd(im.game.GetGameData().MyId, cmdData)
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
			if handle := im.cmd2KeyHandle[cmdData.Cmd()]; handle != nil {
				handle(cmdData.Args())
				continue
			}
			im.gameCore.PushSyncPlayerCmd(im.game.GetGameData().MyId, cmdData)
		}
	}

	xo, yo := ebiten.Wheel()
	if im.wheelHandle != nil {
		im.wheelHandle(xo, yo)
	}
}
