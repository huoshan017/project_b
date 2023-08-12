package ui

import (
	_ "image/png"
	"project_b/client/images"
	"project_b/client_base"
	"project_b/client_core"
	"project_b/common"

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

// InWorldUI
type InWorldUI struct {
	revive    PopupReviveUI
	pauseMenu PauseMenuUI
	debug     DebugUI
}

// InWorld.Init
func (ui *InWorldUI) Init(game client_base.IGame) {
	ui.revive.Init(game)
	ui.pauseMenu.Init(game)
	ui.debug.Init(game)
}

// InWorld.Update
func (ui *InWorldUI) Update() {
	ui.revive.Update()
	ui.pauseMenu.Update()
	ui.debug.Update()
}

// InWorld.GenFrame
func (ui *InWorldUI) DrawFrame() {
	ui.revive.DrawFrame()
	ui.pauseMenu.DrawFrame()
	ui.debug.DrawFrame()
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
func (ui *PopupReviveUI) DrawFrame() {
	if !ui.popup {
		return
	}
	cx, cy := ui.GetScreenCenterPos()
	imgui.SetNextWindowPos(imgui.Vec2{X: float32(cx) - ui.s.X/2, Y: float32(cy) - ui.s.Y/2})
	imgui.BeginV("revive", nil, imgui.WindowFlagsNoDecoration|imgui.WindowFlagsNoMove)
	imgui.Text("Revive/Exit")
	if imgui.Button("Revive") {
		ui.toRevive = true
	} else if imgui.Button("Exit") {
		ui.toExit = true
	}
	ui.s = imgui.WindowSize()
	imgui.End()
}

// PauseMenuUI
type PauseMenuUI struct {
	PopupBase
	resume bool
}

// PauseMenuUI.Init
func (ui *PauseMenuUI) Init(game client_base.IGame) {
	ui.uiBase.Init(game)
}

// PauseMenuUI.Update
func (ui *PauseMenuUI) Update() {
	if inpututil.IsKeyJustReleased(ebiten.KeyEscape) || ui.resume {
		if ui.popup {
			ui.game.GameLogic().Resume()
			ui.pop(false)
			ui.resume = false
		} else {
			ui.game.GameLogic().Pause()
			ui.pop(true)
		}
	}
}

// PauseMenuUI.GenFrame
func (ui *PauseMenuUI) DrawFrame() {
	if !ui.popup {
		return
	}
	cx, cy := ui.GetScreenCenterPos()
	imgui.SetNextWindowPos(imgui.Vec2{X: float32(cx) - ui.s.X/2, Y: float32(cy) - ui.s.Y/2})
	imgui.BeginV("popup_menu", nil, imgui.WindowFlagsNoDecoration|imgui.WindowFlagsNoMove)
	imgui.SameLine()
	imgui.Text("Pause")
	if imgui.Button("Resume") {
		ui.resume = true
	}
	imgui.Button("Options")
	imgui.Button("Restart")
	imgui.Button("Exit")
	ui.s = imgui.WindowSize()
	imgui.End()
}

// DebugUI
type DebugUI struct {
	PopupBase
	showMapGridSelected                      bool // 地圖網格顯示
	showTankBBSelected                       bool // 坦克包圍盒顯示
	showTankAABBSelected                     bool // 坦克AABB顯示
	showShellBBSelected                      bool // 炮彈包圍盒顯示
	showShellAABBSelected                    bool // 炮彈AABB顯示
	showTankCollisionDetectionRegionSelected bool // 坦克碰撞檢測區域顯示
}

// DebugUI.Init
func (ui *DebugUI) Init(game client_base.IGame) {
	ui.uiBase.Init(game)
}

// DebugUI.Update
func (ui *DebugUI) Update() {
	if inpututil.IsKeyJustReleased(ebiten.KeyF1) {
		if ui.popup {
			ui.pop(false)
		} else {
			ui.pop(true)
		}
		return
	}
	debug := ui.game.Debug()
	if ui.showMapGridSelected {
		debug.ShowMapGrid()
	} else {
		debug.HideMapGrid()
	}
	if ui.showTankBBSelected {
		debug.ShowTankBoundingbox()
	} else {
		debug.HideTankBoundingbox()
	}
	if ui.showTankAABBSelected {
		debug.ShowTankAABB()
	} else {
		debug.HideTankAABB()
	}
	if ui.showShellBBSelected {
		debug.ShowShellBoundingbox()
	} else {
		debug.HideShellBoundingbox()
	}
	if ui.showShellAABBSelected {
		debug.ShowShellAABB()
	} else {
		debug.HideShellAABB()
	}
	if ui.showTankCollisionDetectionRegionSelected {
		debug.ShowTankCollisionDetectionRegion()
	} else {
		debug.HideTankCollisionDetectionRegion()
	}
}

// DebugUI.DrawFrame
func (ui *DebugUI) DrawFrame() {
	if !ui.popup {
		return
	}
	imgui.Begin("debug_ui")
	imgui.Checkbox("show map grid", &ui.showMapGridSelected)
	imgui.Checkbox("show tank boundingbox", &ui.showTankBBSelected)
	imgui.Checkbox("show tank AABB", &ui.showTankAABBSelected)
	imgui.Checkbox("show shell boundingbox", &ui.showShellBBSelected)
	imgui.Checkbox("show shell AABB", &ui.showShellAABBSelected)
	imgui.Checkbox("show tank collision detection region", &ui.showTankCollisionDetectionRegionSelected)
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
	ui.initTextures()
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
		ui.login.DrawFrame()
	case client_base.GameStateInWorld:
		ui.inWorld.DrawFrame()
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
