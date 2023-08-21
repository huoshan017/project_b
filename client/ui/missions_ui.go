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
	selListItem int32
	toEnter     bool
	mapNameList []string
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
}

func (ui *MissionsSubUI) DrawFrame() {
	imgui.BeginGroup()
	var item = ui.selListItem
	imgui.SetNextItemWidth(ui.w)
	if imgui.ListBox("", &item, ui.mapNameList) {
		ui.selListItem = item
	}
	if imgui.Button("Back") {
		ui.Back()
	}
	imgui.SameLine()
	if ui.selListItem >= 0 {
		if imgui.Button("Enter") {
			ui.toEnter = true
		}
	}
	imgui.EndGroup()
}

func (ui *MissionsSubUI) enterGame() {
	// 载入地图
	mapId := common_data.MapIdList[ui.selListItem]
	config := game_map.MapConfigArray[mapId]
	if !ui.game.Inst().LoadScene(config) {
		log.Error("load map %v error", mapId)
		return
	}
	ui.game.Inst().CheckAndStart([]uint64{client_core.DefaultSinglePlayerId})
	ui.game.EventMgr().InvokeEvent(client_core.EventIdPlayerEnterGame, "", client_core.DefaultSinglePlayerId)
}
