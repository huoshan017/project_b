package ui

import "project_b/client_base"

type SavesSubUI struct {
	SubUI
}

func (ui *SavesSubUI) Init(game client_base.IGame) {
	ui.SubUI.Init(game)
}

func (ui *SavesSubUI) Update() {
	ui.SubUI.Update()
}

func (ui *SavesSubUI) DrawFrame() {

}
