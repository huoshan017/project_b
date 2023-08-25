package ui

import (
	"project_b/client/images"
	"project_b/client_core"

	"github.com/inkyblackness/imgui-go/v4"
)

// loginUI
type loginUI struct {
	uiBase
	inputAccount  string
	inputPassword string
	toLogin       bool
}

// loginUI.Update
func (ui *loginUI) Update() {
	if ui.toLogin {
		ui.game.EventMgr().InvokeEvent(client_core.EventIdOpLogin, ui.inputAccount, ui.inputPassword)
		ui.toLogin = false
	}
}

// loginUI.DrawFrame
func (ui *loginUI) DrawFrame() {
	w, h := ui.game.ScreenWidthHeight()
	imgui.SetNextWindowPos(imgui.Vec2{X: 0, Y: 0})
	imgui.SetNextWindowSize(imgui.Vec2{X: float32(w), Y: float32(h)})
	imgui.BeginV("login", nil, imgui.WindowFlagsNoDecoration|imgui.WindowFlagsNoMove)
	bounds := images.GetTitleImg().Bounds()
	imgui.SetCursorPos(imgui.Vec2{X: float32(w/2 - int32(bounds.Dx())/2), Y: float32(h/2 - int32(bounds.Dy()))})
	imgui.Image(titleImgId, imgui.Vec2{X: float32(bounds.Dx()), Y: float32(bounds.Dy())})
	imgui.SetCursorPos(imgui.Vec2{X: float32(w/2 - int32(bounds.Dx())/2), Y: float32(h/2 + int32(bounds.Dy()/10))})
	imgui.Text("Account")
	imgui.SetCursorPos(imgui.Vec2{X: float32(w/2 - int32(bounds.Dx())/2), Y: float32(h/2 + int32(bounds.Dy()/4))})
	imgui.SetNextItemWidth(float32(bounds.Dx()))
	imgui.InputText("##Account", &ui.inputAccount)
	imgui.SetCursorPos(imgui.Vec2{X: float32(w/2 - int32(bounds.Dx())/2), Y: float32(h/2 + int32(bounds.Dy()*3/7))})
	imgui.Text("Password")
	imgui.SetCursorPos(imgui.Vec2{X: float32(w/2 - int32(bounds.Dx())/2), Y: float32(h/2 + int32(bounds.Dy()*4/7))})
	imgui.SetNextItemWidth(float32(bounds.Dx()))
	imgui.InputText("##Password", &ui.inputPassword)
	imgui.SetCursorPos(imgui.Vec2{X: float32(w/2 - int32(bounds.Dx())/20), Y: float32(h/2 + int32(bounds.Dy()*6/7))})
	if imgui.Button("login") {
		ui.toLogin = true
	}
	imgui.End()
}
