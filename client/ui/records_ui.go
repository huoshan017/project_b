package ui

import (
	"project_b/client_base"
	"project_b/core"

	"github.com/inkyblackness/imgui-go/v4"
)

type RecordsSubUI struct {
	SubUI
	recordMgr *core.RecordManager
	selIndex  int32
	isPlay    bool
	isDelete  bool
}

func (ui *RecordsSubUI) Init(game client_base.IGame) {
	ui.SubUI.Init(game)
	ui.recordMgr = game.RecordMgr()
	ui.selIndex = -1
}

func (ui *RecordsSubUI) Update() {
	if ui.selIndex >= 0 {
		if ui.isPlay {
			ui.game.RecordMgr().Select(ui.selIndex)
			ui.game.ToReplay()
			ui.isPlay = false
		}
		if ui.isDelete {
			ui.game.RecordMgr().Delete(ui.selIndex)
			ui.isDelete = false
		}
	}
	ui.SubUI.Update()
}

func (ui *RecordsSubUI) DrawFrame() {
	x := ui.left
	y := ui.top
	imgui.SetCursorPos(imgui.Vec2{X: x, Y: y})
	imgui.SetNextItemWidth(ui.w / 3)
	var isSelected bool
	if imgui.BeginListBoxV("", imgui.Vec2{X: x, Y: ui.h * 2 / 3}) {
		for i := int32(0); i < ui.recordMgr.GetRecordCount(); i++ {
			isSelected = ui.selIndex == int32(i)
			if imgui.SelectableV(ui.recordMgr.GetRecordName(i), isSelected, 0, imgui.Vec2{}) {
				ui.selIndex = int32(i)
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
	if ui.selIndex >= 0 {
		imgui.PushItemFlag(imgui.ItemFlagsDisabled, false)
	} else {
		imgui.PushItemFlag(imgui.ItemFlagsDisabled, true)
	}
	if imgui.Button("Play") {
		ui.isPlay = true
	}
	imgui.SameLine()
	if imgui.Button("Delete") {
		ui.isDelete = true
	}
	imgui.PopItemFlag()
}
