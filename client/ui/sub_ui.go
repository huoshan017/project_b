package ui

import (
	"project_b/client_base"
	"project_b/common/base"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type SubUI struct {
	uiBase
	w, h      float32
	left, top float32
	backEvent base.Event
}

func (ui *SubUI) Init(game client_base.IGame) {
	ui.uiBase.Init(game)
}

func (ui *SubUI) Update() {
	if inpututil.IsKeyJustReleased(ebiten.KeyEscape) {
		ui.Back()
		return
	}
}

func (ui *SubUI) SetRect(left, top, w, h float32) {
	ui.left = left
	ui.top = top
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
