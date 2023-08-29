package ui

import (
	"project_b/client_base"
	"project_b/client_core"
	"project_b/common_data"
	"project_b/game_map"
	"project_b/log"

	"github.com/inkyblackness/imgui-go/v4"
)

type MissionsSubUI struct {
	SubUI
	mapNameList []string
	selListItem int32
	toEnter     bool
	isRecord    bool
}

func (ui *MissionsSubUI) Init(game client_base.IGame) {
	ui.SubUI.Init(game)
	for i := 0; i < len(common_data.MapIdList); i++ {
		ui.mapNameList = append(ui.mapNameList, game_map.MapConfigArray[common_data.MapIdList[i]].Name)
	}
	ui.selListItem = -1
}

func (ui *MissionsSubUI) Update() {
	if ui.toEnter && ui.selListItem >= 0 {
		ui.enterGame()
		ui.toEnter = false
	}
	ui.SubUI.Update()
}

func (ui *MissionsSubUI) DrawFrame() {
	x := ui.left
	y := ui.top
	imgui.SetCursorPos(imgui.Vec2{X: x, Y: y})
	imgui.SetNextItemWidth(ui.w / 3)
	var isSelected bool
	if imgui.BeginListBoxV("", imgui.Vec2{X: x, Y: ui.h * 2 / 3}) {
		for i := 0; i < len(ui.mapNameList); i++ {
			isSelected = ui.selListItem == int32(i)
			if imgui.SelectableV(ui.mapNameList[i], isSelected, 0, imgui.Vec2{}) {
				ui.selListItem = int32(i)
			}
			if isSelected {
				imgui.SetItemDefaultFocus()
			}
		}
		imgui.EndListBox()
	}
	imgui.SetCursorPos(imgui.Vec2{X: x, Y: y + ui.h*3/4})
	if imgui.Button("Back") {
		ui.Back()
	}
	imgui.SameLine()
	if ui.selListItem >= 0 {
		imgui.PushItemFlag(imgui.ItemFlagsDisabled, false)
	} else {
		imgui.PushItemFlag(imgui.ItemFlagsDisabled, true)
	}
	if imgui.Button("Enter") {
		ui.toEnter = true
	}
	imgui.PopItemFlag()
	imgui.SameLine()
	imgui.Checkbox("Record Game", &ui.isRecord)
}

func (ui *MissionsSubUI) enterGame() {
	// 载入地图
	mapId := common_data.MapIdList[ui.selListItem]
	if !ui.game.GameCore().LoadMap(mapId) {
		log.Error("load map %v error", mapId)
		return
	}
	ui.game.GameCore().Start([]uint64{client_core.DefaultSinglePlayerId}, ui.isRecord)
	ui.game.EventMgr().InvokeEvent(client_core.EventIdPlayerEnterGame, "", client_core.DefaultSinglePlayerId)
}
