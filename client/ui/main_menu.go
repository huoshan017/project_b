package ui

import (
	"project_b/client/images"
	"project_b/client_base"

	"github.com/inkyblackness/imgui-go/v4"
)

const (
	menuNone menuItemId = iota

	menuSingleplay_missions = 10
	menuSingleplay_records  = 11
	menuSingleplay_saves    = 12

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
			{id: menuSingleplay_records, name: "Records", exec: menuUI.toRecordsUI},
			{id: menuSingleplay_saves, name: "Saves", exec: menuUI.toSavesUI},
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
	recordListSubUI RecordsSubUI
	saveListSubUI   SavesSubUI
}

func (ui *MainMenuUI) Init(game client_base.IGame, menuIdNodeList []menuIdNode) {
	ui.Menu.Init(game, menuIdNodeList)
	ui.missionsSubUI.Init(game)
	ui.recordListSubUI.Init(game)
	ui.saveListSubUI.Init(game)
}

func (ui *MainMenuUI) Update() {
	ui.Menu.update()
	switch ui.currSelected {
	case menuSingleplay_missions:
		ui.missionsSubUI.Update()
	case menuSingleplay_records:
		ui.recordListSubUI.Update()
	case menuSingleplay_saves:
		ui.saveListSubUI.Update()
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
	if ui.currSelected == menuNone {
		// background image
		bounds := images.GetTitleImg().Bounds()
		var left = float32(w/2 - int32(bounds.Dx())/2)
		var top = float32(h/2 - int32(bounds.Dy()))
		imgui.SetCursorPos(imgui.Vec2{X: left, Y: top})
		imgui.Image(titleImgId, imgui.Vec2{X: float32(bounds.Dx()), Y: float32(bounds.Dy())})
		localPos := imgui.Vec2{X: float32(w/2 - int32(bounds.Dx())/4), Y: float32(h/2 + int32(bounds.Dy()/5))}
		posInterval := imgui.Vec2{X: 0, Y: float32(int32(bounds.Dy() / 5))}
		buttonSize := imgui.Vec2{X: float32(bounds.Dx()) / 2, Y: float32(bounds.Dy() / 6)}
		ui.Menu.draw(localPos, posInterval, buttonSize)
	} else {
		switch ui.currSelected {
		case menuSingleplay_missions:
			ww, hh := w/3, h/3
			ui.missionsSubUI.SetRect(float32(w/2-ww/2), float32(h/2-hh/2), float32(ww), float32(hh))
			ui.missionsSubUI.DrawFrame()
		case menuSingleplay_records:
			ww, hh := w/3, h/2
			ui.recordListSubUI.SetRect(float32(w/2-ww/2), float32(h/2-hh/2), float32(ww), float32(hh))
			ui.recordListSubUI.DrawFrame()
		case menuSingleplay_saves:
			ww, hh := w/3, h/2
			ui.saveListSubUI.SetRect(float32(w/2-ww/2), float32(h/2-hh/2), float32(ww), float32(hh))
			ui.saveListSubUI.DrawFrame()
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

func (ui *MainMenuUI) toRecordsUI() {
	ui.recordListSubUI.RegisterBackEventHandle(ui.back2Menu)
	ui.recordListSubUI.recordMgr.LoadRecords()
}

func (ui *MainMenuUI) toSavesUI() {
	ui.saveListSubUI.RegisterBackEventHandle(ui.back2Menu)
}

func (ui *MainMenuUI) back2Menu(...any) {
	ui.prevSelected = ui.currSelected
	ui.currSelected = menuNone
	if ui.prevSelected == menuSingleplay_missions {
		ui.missionsSubUI.UnregisterBackEventHandle(ui.back2Menu)
	} else if ui.prevSelected == menuSingleplay_records {
		ui.recordListSubUI.UnregisterBackEventHandle(ui.back2Menu)
	} else if ui.prevSelected == menuSingleplay_saves {
		ui.saveListSubUI.UnregisterBackEventHandle(ui.back2Menu)
	}
}
