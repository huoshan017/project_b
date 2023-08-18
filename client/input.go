package main

import (
	"project_b/client_base"
	"project_b/common/object"
	"project_b/core"
	"project_b/log"

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
	ebiten.KeyA:        core.NewCmdData(core.CMD_TANK_MOVE, []any{object.DirLeft}),
	ebiten.KeyD:        core.NewCmdData(core.CMD_TANK_MOVE, []any{object.DirRight}),
	ebiten.KeyW:        core.NewCmdData(core.CMD_TANK_MOVE, []any{object.DirUp}),
	ebiten.KeyS:        core.NewCmdData(core.CMD_TANK_MOVE, []any{object.DirDown}),
	ebiten.KeyJ:        core.NewCmdData(core.CMD_TANK_FIRE, []any{1}),
	ebiten.KeyI:        core.NewCmdData(core.CMD_TANK_FIRE, []any{2}),
	ebiten.KeyUp:       core.NewCmdData(CMD_CAMERA_UP, []any{10}),
	ebiten.KeyDown:     core.NewCmdData(CMD_CAMERA_DOWN, []any{-10}),
	ebiten.KeyLeft:     core.NewCmdData(CMD_CAMERA_LEFT, []any{-10}),
	ebiten.KeyRight:    core.NewCmdData(CMD_CAMERA_RIGHT, []any{10}),
	ebiten.KeyPageUp:   core.NewCmdData(CMD_CAMERA_HEIGHT, []any{10}),
	ebiten.KeyPageDown: core.NewCmdData(CMD_CAMERA_HEIGHT, []any{-10}),
}

// 释放键位映射命令
var keyReleased2CmdMap = map[ebiten.Key]core.CmdCode{
	ebiten.KeyA: core.CMD_TANK_STOP,
	ebiten.KeyD: core.CMD_TANK_STOP,
	ebiten.KeyW: core.CMD_TANK_STOP,
	ebiten.KeyS: core.CMD_TANK_STOP,
	ebiten.KeyC: core.CMD_TANK_CHANGE,
	ebiten.KeyR: core.CMD_TANK_RESTORE,
	ebiten.Key1: core.CMD_RELEASE_SMALL_BALL,
	ebiten.KeyU: core.CMD_TANK_SHIELD,
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
	0: core.NewCmdData(core.CMD_TANK_MOVE, []any{object.DirLeftUp}),
	1: core.NewCmdData(core.CMD_TANK_MOVE, []any{object.DirLeftDown}),
	2: core.NewCmdData(core.CMD_TANK_MOVE, []any{object.DirRightUp}),
	3: core.NewCmdData(core.CMD_TANK_MOVE, []any{object.DirRightDown}),
}

// 輸入管理器
type InputMgr struct {
	game        client_base.IGame
	inst        *core.Instance
	pressedKeys []ebiten.Key
	keyPressMap map[ebiten.Key]struct{}
	cmd2Handle  map[core.CmdCode]func(...any)
}

// 創建輸入管理器
func NewInputMgr(game client_base.IGame, inst *core.Instance) *InputMgr {
	return &InputMgr{
		game:        game,
		inst:        inst,
		keyPressMap: make(map[ebiten.Key]struct{}),
		cmd2Handle:  make(map[core.CmdCode]func(...any)),
	}
}

// 添加處理器
func (im *InputMgr) AddHandle(cc core.CmdCode, handle func(...any)) {
	im.cmd2Handle[cc] = handle
}

// 處理輸入
func (im *InputMgr) HandleInput() {
	var (
		cmd core.CmdCode
		o   bool
	)
	// 處理鍵釋放
	for k := range im.keyPressMap {
		if inpututil.IsKeyJustReleased(k) {
			cmd, o = keyReleased2CmdMap[k]
			if !o {
				continue
			}
			im.inst.PushFrame(0, im.game.GetGameData().MyId, cmd, nil)
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
	var cmdData *core.CmdData
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
						im.inst.PushFrame(0, im.game.GetGameData().MyId, cmdData.Cmd(), cmdData.Args())
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
			if handle := im.cmd2Handle[cmdData.Cmd()]; handle != nil {
				handle(cmdData.Args()...)
				continue
			}
			im.inst.PushFrame(0, im.game.GetGameData().MyId, cmdData.Cmd(), cmdData.Args())
		}
	}
}
