package ui

import "project_b/client_base"

// InWorldUI
type InWorldUI struct {
	revive    PopupReviveUI
	pauseMenu PauseMenuUI
	inPlay    PlayUI
	debug     DebugUI
}

// InWorld.Init
func (ui *InWorldUI) Init(game client_base.IGame) {
	ui.revive.Init(game)
	ui.pauseMenu.Init(game, getPauseMenuIdNodeTree(&ui.pauseMenu))
	ui.inPlay.Init(game)
	ui.debug.Init(game)
}

// InWorld.Update
func (ui *InWorldUI) Update() {
	ui.revive.Update()
	ui.pauseMenu.Update()
	ui.inPlay.Update()
	ui.debug.Update()
}

// InWorld.GenFrame
func (ui *InWorldUI) DrawFrame() {
	ui.revive.DrawFrame()
	ui.pauseMenu.DrawFrame()
	ui.inPlay.DrawFrame()
	ui.debug.DrawFrame()
}
