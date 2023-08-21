package ui

import (
	"project_b/client_base"
	"project_b/common/base"
)

type SubUI struct {
	uiBase
	w, h      float32
	backEvent base.Event
}

func (ui *SubUI) Init(game client_base.IGame) {
	ui.uiBase.Init(game)
}

func (ui *SubUI) SetWidthHight(w, h float32) {
	ui.w, ui.h = w, h
}

func (ui *SubUI) RegisterBackEventHandle(handle func(...any)) {
	ui.backEvent.Register(handle)
}

func (ui *SubUI) UnregisterBackEventHandle(handle func(...any)) {
	ui.backEvent.Unregister(handle)
}

func (ui *SubUI) Back() {
	ui.backEvent.Call()
}
