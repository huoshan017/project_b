package ui

import (
	"project_b/client_base"

	"github.com/inkyblackness/imgui-go/v4"
)

const (
	selQueryNone = iota
	selQueryOk
	selQueryCancel
)

type QueryDialog struct {
	uiBase
	exec     func()
	queryStr string
	sel      int
}

func NewQueryDialog() *QueryDialog {
	return &QueryDialog{}
}

func (ui *QueryDialog) Init(game client_base.IGame) {
	ui.uiBase.Init(game)
}

func (ui *QueryDialog) SetQueryString(queryStr string) {
	ui.queryStr = queryStr
}

func (ui *QueryDialog) SetExec(exec func()) {
	ui.exec = exec
}

func (ui *QueryDialog) Update() {
	if ui.sel == selQueryOk {
		ui.exec()
		ui.sel = selQueryCancel
	}
}

func (ui *QueryDialog) DrawFrame() {
	if ui.sel == selQueryCancel {
		return
	}
	w, h := ui.game.ScreenWidthHeight()
	s := imgui.Vec2{X: float32(w) / 6, Y: float32(h) / 9}
	imgui.SetNextWindowSize(s)
	pos := imgui.Vec2{X: float32(w)/2 - float32(s.X)/2, Y: float32(h)/2 - float32(s.Y)/2}
	imgui.SetNextWindowPos(pos)
	imgui.BeginV("QueryDialog", nil, imgui.WindowFlagsNoDecoration|imgui.WindowFlagsNoMove)
	imgui.Text(ui.queryStr)
	buttonSize := imgui.Vec2{X: s.X * 2 / 5, Y: s.Y * 2 / 7}
	if imgui.ButtonV("Ok", buttonSize) {
		ui.sel = selQueryOk
	}
	imgui.SameLine()
	if imgui.ButtonV("Cancel", buttonSize) {
		ui.sel = selQueryCancel
	}
	imgui.End()
}

func (ui *QueryDialog) Show() {
	ui.sel = selQueryNone
}

func (ui *QueryDialog) IsShow() bool {
	return ui.sel == selQueryNone
}
