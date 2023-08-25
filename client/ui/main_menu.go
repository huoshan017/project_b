package ui

import (
	"project_b/client/images"
	"project_b/client_base"

	"github.com/inkyblackness/imgui-go/v4"
)

const (
	menuNone menuItemId = iota

	menuSingleplay_missions = 10
	menuSingleplay_replays  = 11

	menuMultiplay_local    = 20
	menuMultiplay_internet = 21

	menuSettings_video    = 30
	menuSettings_audio    = 31
	menuSettings_gameplay = 32
	menuSettings_controls = 33
)

// 主菜單id節點樹
var getMainMenuIdNodeTree = func(menuUI *MainMenuUI) []menuIdNode {
	return []menuIdNode{
		{id: menuNone, name: "Single Player", itemList: []menuIdNode{
			{id: menuSingleplay_missions, name: "Missions", exec: menuUI.toMissionsUI},
			{id: menuSingleplay_replays, name: "Replays", exec: menuUI.toReplaysUI},
			{id: menuNone, name: "Back", exec: menuUI.back},
		}},
		{id: menuNone, name: "Multiple Players", itemList: []menuIdNode{
			{id: menuMultiplay_local, name: "Local", itemList: nil},
			{id: menuMultiplay_internet, name: "Internet", itemList: nil},
			{id: menuNone, name: "Back", exec: menuUI.back},
		}},
		{id: menuNone, name: "Settings", itemList: []menuIdNode{
			{id: menuSettings_video, name: "Video", itemList: nil},
			{id: menuSettings_audio, name: "Audio", itemList: nil},
			{id: menuSettings_gameplay, name: "Gameplay", itemList: nil},
			{id: menuSettings_controls, name: "Controls", itemList: nil},
			{id: menuNone, name: "Back", exec: menuUI.back},
		}},
		{id: menuNone, name: "Exit", itemList: nil},
	}
}

type MainMenuUI struct {
	Menu
	missionsSubUI   MissionsSubUI
	replayListSubUI ReplaysSubUI
}

func (ui *MainMenuUI) Init(game client_base.IGame, menuIdNodeList []menuIdNode) {
	ui.Menu.Init(game, menuIdNodeList)
	ui.missionsSubUI.Init(game)
	ui.replayListSubUI.Init(game)
}

func (ui *MainMenuUI) Update() {
	ui.Menu.update()
	switch ui.currSelected {
	case menuSingleplay_missions:
		ui.missionsSubUI.Update()
	case menuSingleplay_replays:
		ui.replayListSubUI.Update()
	case menuMultiplay_local:
	case menuMultiplay_internet:
	case menuSettings_audio:
	case menuSettings_video:
	case menuSettings_gameplay:
	case menuSettings_controls:
	}
}

func (ui *MainMenuUI) DrawFrame() {
	w, h := ui.game.ScreenWidthHeight()
	imgui.SetNextWindowPos(imgui.Vec2{X: 0, Y: 0})
	ui.s = imgui.Vec2{X: float32(w), Y: float32(h)}
	imgui.SetNextWindowSize(ui.s)
	imgui.BeginV("Main Menu", nil, imgui.WindowFlagsNoDecoration|imgui.WindowFlagsNoMove)
	// background image
	bounds := images.GetTitleImg().Bounds()
	var left = float32(w/2 - int32(bounds.Dx())/2)
	var top = float32(h/2 - int32(bounds.Dy()))
	if ui.currSelected == menuNone {
		imgui.SetCursorPos(imgui.Vec2{X: left, Y: top})
		imgui.Image(titleImgId, imgui.Vec2{X: float32(bounds.Dx()), Y: float32(bounds.Dy())})
		localPos := imgui.Vec2{X: float32(w/2 - int32(bounds.Dx())/4), Y: float32(h/2 + int32(bounds.Dy()/5))}
		posInterval := imgui.Vec2{X: 0, Y: float32(int32(bounds.Dy() / 5))}
		buttonSize := imgui.Vec2{X: float32(bounds.Dx()) / 2, Y: float32(bounds.Dy() / 6)}
		ui.Menu.draw(localPos, posInterval, buttonSize)
	} else {
		switch ui.currSelected {
		case menuSingleplay_missions:
			ui.missionsSubUI.SetRect(left, top, ui.s.X/5, ui.s.Y/3)
			ui.missionsSubUI.DrawFrame()
		case menuSingleplay_replays:
			ui.replayListSubUI.SetRect(left, top, ui.s.X/5, ui.s.Y/3)
			ui.replayListSubUI.DrawFrame()
		case menuMultiplay_local:
		case menuMultiplay_internet:
		case menuSettings_audio:
		case menuSettings_video:
		case menuSettings_gameplay:
		case menuSettings_controls:
		}
	}
	imgui.End()
}

func (ui *MainMenuUI) toMissionsUI() {
	ui.missionsSubUI.RegisterBackEventHandle(ui.back2Menu)
}

func (ui *MainMenuUI) toReplaysUI() {
	ui.replayListSubUI.RegisterBackEventHandle(ui.back2Menu)
}

func (ui *MainMenuUI) back2Menu(...any) {
	ui.prevSelected = ui.currSelected
	ui.currSelected = menuNone
	if ui.prevSelected == menuSingleplay_missions {
		ui.missionsSubUI.UnregisterBackEventHandle(ui.back2Menu)
	} else if ui.prevSelected == menuSingleplay_replays {
		ui.replayListSubUI.UnregisterBackEventHandle(ui.back2Menu)
	}
}
