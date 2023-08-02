package ui

import (
	_ "image/png"
	"project_b/client_base"
	"project_b/client_core"
	"project_b/common"

	"github.com/gabstv/ebiten-imgui/renderer"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/inkyblackness/imgui-go/v4"
)

// uiBase
type uiBase struct {
	game client_base.IGame
}

// uiBase.Init
func (ui *uiBase) Init(game client_base.IGame) {
	ui.game = game
}

// uiBase.SetToScreenCenter
func (ui *uiBase) GetScreenCenterPos() (int32, int32) {
	w, h := ui.game.ScreenWidthHeight()
	return w / 2, h / 2
}

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

// loginUI.GenFrame
func (ui *loginUI) GenFrame() {
	var opened = true
	imgui.BeginV("login", &opened, imgui.WindowFlagsNoDecoration|imgui.WindowFlagsNoMove)
	imgui.Text("Account")
	imgui.InputText("##Account", &ui.inputAccount)
	imgui.Text("Password")
	imgui.InputText("##Password", &ui.inputPassword)
	if imgui.Button("login") {
		ui.toLogin = true
	}
	imgui.End()
}

// InWorldUI
type InWorldUI struct {
	revive    PopupReviveUI
	pauseMenu PauseMenuUI
}

// InWorld.Init
func (ui *InWorldUI) Init(game client_base.IGame) {
	ui.revive.Init(game)
	ui.pauseMenu.Init(game)
}

// InWorld.Update
func (ui *InWorldUI) Update() {
	ui.revive.Update()
	ui.pauseMenu.Update()
}

// InWorld.GenFrame
func (ui *InWorldUI) GenFrame() {
	ui.revive.GenFrame()
	ui.pauseMenu.GenFrame()
}

// PopupReviveUI
type PopupReviveUI struct {
	PopupBase
	toRevive bool
	toExit   bool
}

// PopupReviveUI.Init
func (ui *PopupReviveUI) Init(game client_base.IGame) {
	ui.uiBase.Init(game)
	ui.game.EventMgr().RegisterEvent(common.EventIdTankDestroy, func(args ...any) {
		tt := args[0].(common.TankType)
		if tt == common.TankTypePlayer {
			pid := args[1].(uint64)
			if pid == game.GetGameData().MyId {
				ui.pop(true)
			}
		}
	})
}

// PopupReviveUI.Update
func (ui *PopupReviveUI) Update() {
	if ui.toRevive {
		ui.game.CmdMgr().Handle(client_core.CMD_REVIVE, common.TankTypePlayer, ui.game.GetGameData().MyId)
		ui.toRevive = false
		ui.pop(false)
	} else if ui.toExit {
		ui.game.EventMgr().InvokeEvent(common.EventIdExitGame)
		ui.toExit = false
		ui.pop(false)
	}
}

// PopupReviveUI.GenFrame
func (ui *PopupReviveUI) GenFrame() {
	if !ui.popup {
		return
	}
	imgui.BeginV("revive", nil, imgui.WindowFlagsNoDecoration|imgui.WindowFlagsNoMove)
	imgui.Text("Revive/Exit")
	if imgui.Button("Revive") {
		ui.toRevive = true
	} else if imgui.Button("Exit") {
		ui.toExit = true
	}
	imgui.End()
}

// PauseMenuUI
type PauseMenuUI struct {
	PopupBase
	s imgui.Vec2
}

// PauseMenuUI.Init
func (ui *PauseMenuUI) Init(game client_base.IGame) {
	ui.uiBase.Init(game)
}

// PauseMenuUI.Update
func (ui *PauseMenuUI) Update() {
	if inpututil.IsKeyJustReleased(ebiten.KeyEscape) {
		if ui.popup {
			ui.pop(false)
		} else {
			ui.pop(true)
		}
	}
}

// PauseMenuUI.GenFrame
func (ui *PauseMenuUI) GenFrame() {
	if !ui.popup {
		return
	}
	cx, cy := ui.GetScreenCenterPos()
	imgui.SetNextWindowPos(imgui.Vec2{X: float32(cx) - ui.s.X/2, Y: float32(cy) - ui.s.Y/2})
	imgui.BeginV("popup_menu", nil, imgui.WindowFlagsNoDecoration|imgui.WindowFlagsNoMove)
	imgui.Text("Pause")
	if imgui.Button("Resume") {
		ui.pop(false)
	}
	imgui.Button("Options")
	imgui.Button("GiveUp")
	if imgui.Button("Restart") {

	}
	ui.s = imgui.WindowSize()
	imgui.End()
}

// PopupBase
type PopupBase struct {
	uiBase
	popup bool
}

// PopupBase.Init
func (ui *PopupBase) Init(game client_base.IGame) {
	ui.uiBase.Init(game)
}

// PopupBase.Update
func (ui *PopupBase) Update() {
	if inpututil.IsKeyJustReleased(ebiten.KeyEscape) {
		ui.pop(false)
	}
}

// PopupBase.pop
func (ui *PopupBase) pop(enable bool) {
	ui.popup = enable
}

type ImguiManager struct {
	game      client_base.IGame
	renderMgr *renderer.Manager
	login     loginUI
	inWorld   InWorldUI
}

func NewImguiManager(game client_base.IGame) *ImguiManager {
	return &ImguiManager{
		game: game,
	}
}

func (ui *ImguiManager) Init() {
	ui.renderMgr = renderer.New(nil)
	w, h := ui.game.ScreenWidthHeight()
	ui.renderMgr.SetDisplaySize(float32(w), float32(h))
	ui.login.Init(ui.game)
	ui.inWorld.Init(ui.game)
}

func (ui *ImguiManager) Update() {
	ui.renderMgr.Update(1.0 / 60.0)
	ui.login.Update()
	ui.inWorld.Update()
}

func (ui *ImguiManager) Draw(screen *ebiten.Image) {
	var draw bool = true
	ui.renderMgr.BeginFrame()
	switch ui.game.GetState() {
	case client_base.GameStateBeforeLogin:
		ui.login.GenFrame()
	case client_base.GameStateInWorld:
		ui.inWorld.GenFrame()
	default:
		draw = false
	}
	ui.renderMgr.EndFrame()
	if draw {
		ui.renderMgr.Draw(screen)
	}
}
