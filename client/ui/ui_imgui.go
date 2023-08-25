package ui

import (
	_ "image/png"
	"project_b/client/images"
	"project_b/client_base"
	"project_b/common"
	"project_b/core"

	"github.com/gabstv/ebiten-imgui/renderer"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/inkyblackness/imgui-go/v4"
)

var (
	titleImgId imgui.TextureID = 100 // 主界面
)

// uiBase
type uiBase struct {
	game client_base.IGame
	s    imgui.Vec2
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
		ui.game.Inst().PushFrame(ui.game.Inst().GetFrame(), ui.game.GetGameData().MyId, core.CMD_TANK_RESPAWN, []any{common.TankTypePlayer})
		ui.toRevive = false
		ui.pop(false)
	} else if ui.toExit {
		ui.game.EventMgr().InvokeEvent(common.EventIdExitGame)
		ui.toExit = false
		ui.pop(false)
	}
}

// PopupReviveUI.DrawFrame
func (ui *PopupReviveUI) DrawFrame() {
	if !ui.popup {
		return
	}
	cx, cy := ui.GetScreenCenterPos()
	imgui.SetNextWindowPos(imgui.Vec2{X: float32(cx) - ui.s.X/2, Y: float32(cy) - ui.s.Y/2})
	imgui.BeginV("revive", nil, imgui.WindowFlagsNoDecoration|imgui.WindowFlagsNoMove)
	imgui.Text("Revive/Exit?")
	if imgui.Button("Revive") {
		ui.toRevive = true
	} else if imgui.Button("Cancel") {
		ui.toExit = true
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
	mainMenu  MainMenuUI
	inWorld   InWorldUI
	inReplay  InReplayUI
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
	ui.initTextures()
	ui.login.Init(ui.game)
	ui.mainMenu.Init(ui.game, getMainMenuIdNodeTree(&ui.mainMenu))
	ui.inWorld.Init(ui.game)
	ui.inReplay.Init(ui.game)
}

func (ui *ImguiManager) Update() {
	ui.renderMgr.Update(1.0 / 60.0)
	switch ui.game.GetState() {
	case client_base.GameStateBeforeLogin:
		ui.login.Update()
	case client_base.GameStateMainMenu:
		ui.mainMenu.Update()
	case client_base.GameStateInWorld:
		ui.inWorld.Update()
	case client_base.GameStateInReplay:
		ui.inReplay.Update()
	}
}

func (ui *ImguiManager) Draw(screen *ebiten.Image) {
	var draw bool = true
	ui.renderMgr.BeginFrame()
	switch ui.game.GetState() {
	case client_base.GameStateBeforeLogin:
		ui.login.DrawFrame()
	case client_base.GameStateMainMenu:
		ui.mainMenu.DrawFrame()
	case client_base.GameStateInWorld:
		ui.inWorld.DrawFrame()
	case client_base.GameStateInReplay:
		ui.inReplay.DrawFrame()
	default:
		draw = false
	}
	ui.renderMgr.EndFrame()
	if draw {
		ui.renderMgr.Draw(screen)
	}
}

func (ui *ImguiManager) initTextures() {
	ui.renderMgr.Cache.SetTexture(titleImgId, images.GetTitleImg())
}
