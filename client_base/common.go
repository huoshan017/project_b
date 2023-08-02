package client_base

import (
	"project_b/client_core"
	"project_b/common/base"

	"github.com/hajimehoshi/ebiten/v2"
)

type GameState int

const (
	GameStateBeforeLogin   GameState = iota // 登陸前
	GameStateInLogin                        // 登陸中
	GameStateMainMenu                       // 已登陸進入主菜單
	GameStateEnteringWorld                  // 正進入游戲世界
	GameStateInWorld                        // 在游戲世界中
	GameStatePopupInWorld                   // 在游戲中彈出事件
	GameStateExitingWorld                   // 正在退出游戲世界
)

type IGame interface {
	GetState() GameState
	GetGameData() *GameData
	ScreenWidthHeight() (int32, int32)
	EventMgr() base.IEventManager
	CmdMgr() *client_core.CmdHandleManager
}

type IUIMgr interface {
	Init()
	Update()
	Draw(*ebiten.Image)
}

type GameData struct {
	State GameState
	MyId  uint64 // 我的ID
	MyAcc string // 我的帐号
}
