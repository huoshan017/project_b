package ui

import "project_b/client_base"

// InReplayUI
type InReplayUI struct {
	pauseMenu PauseMenuUI
	debug     DebugUI
}

// InReplay.Init
func (ui *InReplayUI) Init(game client_base.IGame) {
	ui.pauseMenu.Init(game, getPauseMenuIdNodeTree(&ui.pauseMenu))
	ui.debug.Init(game)
}

// InReplay.Update
func (ui *InReplayUI) Update() {
	ui.pauseMenu.Update()
	ui.debug.Update()
}

// InReplay.GenFrame
func (ui *InReplayUI) DrawFrame() {
	ui.pauseMenu.DrawFrame()
	ui.debug.DrawFrame()
}
