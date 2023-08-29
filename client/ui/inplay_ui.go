package ui

import (
	"fmt"
	"project_b/client_base"
	"project_b/common/utils"

	"github.com/inkyblackness/imgui-go/v4"
)

type PlayUI struct {
	uiBase
}

func (ui *PlayUI) Init(game client_base.IGame) {
	ui.uiBase.Init(game)
}

func (ui *PlayUI) Update() {

}

func (ui *PlayUI) DrawFrame() {
	w, h := ui.game.ScreenWidthHeight()
	// frame text
	var s = imgui.Vec2{X: float32(w) / 6, Y: float32(h) / 30}
	imgui.SetNextWindowSize(s)
	var pos = imgui.Vec2{X: float32(w)/2 - float32(s.X)/2, Y: float32(h / 50)}
	imgui.SetNextWindowPos(pos)
	imgui.BeginV("Frame Progress", nil, imgui.WindowFlagsNoDecoration|imgui.WindowFlagsNoMove)
	ts := utils.FormatSeconds2TimeString(ui.game.GameCore().UsedMs())
	var buf = fmt.Sprintf("time : %v", ts)
	imgui.Text(buf)
	imgui.End()
}
