package ui

import (
	"project_b/client_base"
	"project_b/client_core"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/inkyblackness/imgui-go/v4"
)

// 暫停菜單id節點樹
var getPauseMenuIdNodeTree = func(ui *PauseMenuUI) []menuIdNode {
	return []menuIdNode{
		{id: menuNone, name: "Resume", itemList: nil, exec: ui.resume},
		{id: menuNone, name: "Settings", itemList: []menuIdNode{
			{id: menuSettings_controls, name: "Controls", itemList: nil},
			{id: menuSettings_gameplay, name: "Gameplay", itemList: nil},
			{id: menuNone, name: "Back", itemList: nil, exec: ui.back},
		}},
		{id: menuNone, name: "Restart", exec: ui.popQueryRestartDialog},
		{id: menuNone, name: "Quit To MainMenu", exec: ui.popQueryQuit2MainMenuDialog},
	}
}

// PauseMenuUI
type PauseMenuUI struct {
	Menu
	isPaused    bool
	queryDialog *QueryDialog
}

// PauseMenuUI.Init
func (ui *PauseMenuUI) Init(game client_base.IGame, menuIdNodeList []menuIdNode) {
	ui.Menu.Init(game, menuIdNodeList)
}

// PauseMenuUI.Update
func (ui *PauseMenuUI) Update() {
	if inpututil.IsKeyJustReleased(ebiten.KeyEscape) {
		if ui.queryDialog != nil && ui.queryDialog.IsShow() {
			return
		}
		if ui.isPaused {
			ui.resume()
		} else {
			ui.pause()
		}
	}
	if ui.isPaused {
		ui.Menu.update()
	}
	if ui.queryDialog != nil {
		ui.queryDialog.Update()
	}
}

// PauseMenuUI.GenFrame
func (ui *PauseMenuUI) DrawFrame() {
	if !ui.isPaused {
		return
	}

	w, h := ui.game.ScreenWidthHeight()
	s := imgui.Vec2{X: float32(w) / 5, Y: float32(h) * 2 / 7}
	imgui.SetNextWindowSize(s)
	pos := imgui.Vec2{X: float32(w)/2 - float32(s.X)/2, Y: float32(h)/2 - float32(s.Y)/2}
	imgui.SetNextWindowPos(pos)
	imgui.BeginV("Paused Menu", nil, imgui.WindowFlagsNoDecoration|imgui.WindowFlagsNoMove)
	imgui.Text("Paused")
	if ui.currSelected == menuNone {
		buttonSize := imgui.Vec2{X: s.X * 3 / 4, Y: s.Y / 6}
		localPos := imgui.Vec2{X: s.X/2 - buttonSize.X/2, Y: s.Y / 7}
		posInterval := imgui.Vec2{X: 0, Y: s.Y / 5}
		ui.Menu.draw(localPos, posInterval, buttonSize)
	} else {
		switch ui.currSelected {
		}
	}
	imgui.End()
	if ui.queryDialog != nil {
		ui.queryDialog.DrawFrame()
	}
}

func (ui *PauseMenuUI) resume() {
	ui.game.GameCore().Resume()
	ui.isPaused = false
}

func (ui *PauseMenuUI) pause() {
	ui.game.GameCore().Pause()
	ui.isPaused = true
}

func (ui *PauseMenuUI) popQueryRestartDialog() {
	if ui.queryDialog == nil {
		ui.queryDialog = NewQueryDialog()
		ui.queryDialog.Init(ui.game)
	}
	ui.queryDialog.SetExec(ui.restart)
	ui.queryDialog.SetQueryString("Is Restart?")
	ui.queryDialog.Show()
}

func (ui *PauseMenuUI) popQueryQuit2MainMenuDialog() {
	if ui.queryDialog == nil {
		ui.queryDialog = NewQueryDialog()
		ui.queryDialog.Init(ui.game)
	}
	ui.queryDialog.SetExec(ui.quit2MainMenu)
	ui.queryDialog.SetQueryString("Quit To MainMenu?")
	ui.queryDialog.Show()
}

func (ui *PauseMenuUI) restart() {
	ui.game.GameCore().Restart()
	ui.game.EventMgr().InvokeEvent(client_core.EventIdPlayerEnterGame, "", client_core.DefaultSinglePlayerId)
	ui.resume()
}

func (ui *PauseMenuUI) quit2MainMenu() {
	ui.game.GameCore().End()
	ui.game.EventMgr().InvokeEvent(client_core.EventIdPlayerExitGame, ui.game.GetGameData().MyId)
	ui.resume()
}
