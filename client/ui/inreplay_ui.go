package ui

import (
	"project_b/client_base"
	"project_b/core"

	"github.com/inkyblackness/imgui-go/v4"
)

type PlayRecordUI struct {
	uiBase
	record core.Record
}

func (ui *PlayRecordUI) Init(game client_base.IGame) {
	ui.uiBase.Init(game)
}

func (ui *PlayRecordUI) Update() {
}

func (ui *PlayRecordUI) DrawFrame() {
	w, h := ui.game.ScreenWidthHeight()
	s := imgui.Vec2{X: float32(w) / 5, Y: float32(h) / 20}
	imgui.SetNextWindowSize(s)
	pos := imgui.Vec2{X: float32(w)/2 - float32(s.X)/2, Y: float32(h) - float32(s.Y)}
	imgui.SetNextWindowPos(pos)
	imgui.BeginV("Record Progress", nil, imgui.WindowFlagsNoDecoration|imgui.WindowFlagsNoMove)
	imgui.Text("play time")

	var o bool
	ui.record, o = ui.game.RecordMgr().SelectedRecord()
	if !o {
		imgui.End()
		return
	}
	var progress float32 = float32(ui.game.Inst().GetFrame()) / float32(ui.record.FrameNum())
	imgui.ProgressBar(progress)
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
