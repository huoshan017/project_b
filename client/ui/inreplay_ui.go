package ui

import (
	"fmt"
	"math"
	"project_b/client_base"
	"project_b/common/utils"
	"project_b/core"

	"github.com/inkyblackness/imgui-go/v4"
)

type PlayRecordUI struct {
	uiBase
	record *core.Record
}

func (ui *PlayRecordUI) Init(game client_base.IGame) {
	ui.uiBase.Init(game)
	ui.record = ui.game.GameCore().GetRecordMgr().SelectedRecord()
}

func (ui *PlayRecordUI) Update() {
}

func (ui *PlayRecordUI) DrawFrame() {
	w, h := ui.game.ScreenWidthHeight()
	var s = imgui.Vec2{X: float32(w) / 3, Y: float32(h) / 12}
	imgui.SetNextWindowSize(s)
	var pos = imgui.Vec2{X: float32(w)/2 - float32(s.X)/2, Y: float32(h) - float32(s.Y)}
	imgui.SetNextWindowPos(pos)
	imgui.BeginV("Record Progress", nil, imgui.WindowFlagsNoDecoration|imgui.WindowFlagsNoMove)
	ut := ui.game.GameCore().GetFrame() * ui.record.FrameMs()
	tt := ui.record.FrameNum() * ui.record.FrameMs()
	utbuf := utils.FormatSeconds2TimeString(ut)
	ttbuf := utils.FormatSeconds2TimeString(tt)
	imgui.Text(fmt.Sprintf("play time: %v / %v", utbuf, ttbuf))
	var progress float32 = float32(ui.game.GameCore().GetFrame()) / float32(ui.record.FrameNum())
	var buf = fmt.Sprintf("%v/%v", ui.game.GameCore().GetFrame(), ui.record.FrameNum())
	imgui.ProgressBarV(progress, imgui.Vec2{X: -math.SmallestNonzeroFloat32, Y: 0}, buf)
	imgui.End()
}

// InReplayUI
type InReplayUI struct {
	pauseMenu  PauseMenuUI
	debug      DebugUI
	playRecord PlayRecordUI
}

// InReplay.Init
func (ui *InReplayUI) Init(game client_base.IGame) {
	ui.pauseMenu.Init(game, getPauseMenuIdNodeTree(&ui.pauseMenu))
	ui.debug.Init(game)
	ui.playRecord.Init(game)
}

// InReplay.Update
func (ui *InReplayUI) Update() {
	ui.pauseMenu.Update()
	ui.debug.Update()
	ui.playRecord.Update()
}

// InReplay.GenFrame
func (ui *InReplayUI) DrawFrame() {
	ui.pauseMenu.DrawFrame()
	ui.debug.DrawFrame()
	ui.playRecord.DrawFrame()
}
