package ui

import (
	"project_b/client_base"

	"github.com/inkyblackness/imgui-go/v4"
)

type menuItemId int32

type menuIdNode struct {
	id       menuItemId
	name     string
	itemList []menuIdNode
	exec     func()
}

type menuItem struct {
	size     imgui.Vec2
	localPos imgui.Vec2
	itemList []menuItem
}

type Menu struct {
	uiBase
	menuIdNodeList   []menuIdNode
	buttonSize       imgui.Vec2
	currSelected     menuItemId
	prevSelected     menuItemId
	selItemIndexList []int32
	menuItemTree     []menuItem
	needUpdate       bool
	execFunc         func()
}

func (ui *Menu) Init(game client_base.IGame, menuIdNodeList []menuIdNode) {
	ui.uiBase.Init(game)
	ui.initMenuTree(menuIdNodeList, &ui.menuItemTree)
	ui.menuIdNodeList = menuIdNodeList
	ui.needUpdate = true
}

func (ui *Menu) initMenuTree(menuIdNodeList []menuIdNode, menuItemList *[]menuItem) {
	*menuItemList = make([]menuItem, len(menuIdNodeList))
	for i := 0; i < len(menuIdNodeList); i++ {
		node := menuIdNodeList[i]
		if len(node.itemList) > 0 {
			ui.initMenuTree(node.itemList, &(*menuItemList)[i].itemList)
		}
	}
}

func (ui *Menu) updateMenuPosSize(menuItemList []menuItem, localPos, posInterval imgui.Vec2) {
	if len(menuItemList) == 0 {
		return
	}
	for i := 0; i < len(menuItemList); i++ {
		pi := posInterval.Times(float32(i))
		menuItemList[i].localPos = localPos.Plus(pi)
		menuItemList[i].size = ui.buttonSize
		if len(menuItemList[i].itemList) > 0 {
			ui.updateMenuPosSize(menuItemList[i].itemList, localPos, posInterval)
		}
	}
}

func (ui *Menu) update() {
	if ui.execFunc != nil {
		ui.execFunc()
		ui.execFunc = nil
	}
}

func (ui *Menu) draw(localPos, posInterval, buttonSize imgui.Vec2) {
	// prepare menu data
	if ui.needUpdate {
		ui.buttonSize = buttonSize
		ui.updateMenuPosSize(ui.menuItemTree, localPos, posInterval)
		ui.needUpdate = false
	}

	// draw menu
	menuItemList := ui.menuItemTree
	menuIdNodeList := ui.menuIdNodeList
	for i := 0; i < len(ui.selItemIndexList); i++ {
		menuItemList = menuItemList[ui.selItemIndexList[i]].itemList
		menuIdNodeList = menuIdNodeList[ui.selItemIndexList[i]].itemList
	}
	ui.drawMenuList(menuItemList, menuIdNodeList)
}

func (ui *Menu) drawMenuList(menuList []menuItem, menuIdList []menuIdNode) {
	for i := 0; i < len(menuList); i++ {
		menuIdItem := &menuIdList[i]
		imgui.SetCursorPos(menuList[i].localPos)
		if imgui.ButtonV(menuIdItem.name, menuList[i].size) {
			if menuIdItem.itemList == nil {
				if menuIdItem.exec != nil {
					ui.execFunc = menuIdItem.exec
				}
				ui.prevSelected = ui.currSelected
				ui.currSelected = menuIdList[i].id
			} else {
				ui.selItemIndexList = append(ui.selItemIndexList, int32(i))
			}
		}
	}
}

func (ui *Menu) back() {
	l := len(ui.selItemIndexList)
	ui.selItemIndexList = ui.selItemIndexList[:l-1]
}
