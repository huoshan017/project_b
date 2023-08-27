package ui

import (
	"project_b/client_base"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/inkyblackness/imgui-go/v4"
)

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
	debug := client_base.GetDebug()
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
