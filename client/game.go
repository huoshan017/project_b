package main

import (
	"fmt"
	"project_b/client/ui"
	"project_b/client_base"
	"project_b/client_core"
	core "project_b/client_core"
	"project_b/common/base"
	"project_b/common/time"
	"project_b/common_data"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/huoshan017/gsnet/options"
)

type Game struct {
	//---------------------------------------
	// 逻辑
	net           *core.NetClient        // 网络模块
	msgHandler    *core.MsgHandler       // 消息处理器
	logic         *core.GameLogic        // 游戏逻辑
	cmdMgr        *core.CmdHandleManager // 命令处理管理器
	playerMgr     *core.CPlayerManager   // 玩家管理器，這裏的玩家是指獨立於游戲邏輯GameLogic之外的登錄用戶
	eventMgr      *base.EventManager     // 游戏事件管理器，向上层逻辑传递事件
	lastCheckTime time.CustomTime        // 上次检测时间

	//---------------------------------------
	// 表现相关

	viewport      *client_base.Viewport // 视口
	playableScene *PlayableScene        // 場景圖繪製
	uiMgr         client_base.IUIMgr    // UI管理
	eventHandles  *EventHandles         // 事件处理
	inputMgr      *InputMgr             // 輸入管理器
	gameData      client_base.GameData  // 其他游戏数据
}

// 创建游戏
func NewGame(conf *Config) *Game {
	g := &Game{
		viewport: client_base.CreateViewport(0, 0, screenWidth, screenHeight),
	}
	g.net = core.CreateNetClient(conf.serverAddress, options.WithRunMode(options.RunModeOnlyUpdate), options.WithNoDelay(true))
	g.eventMgr = base.NewEventManager()
	g.logic = core.CreateGameLogic(g.eventMgr)
	g.cmdMgr = core.CreateCmdHandleManager(g.net, g.logic)
	g.playerMgr = core.CreateCPlayerManager()
	g.msgHandler = core.CreateMsgHandler(g.net, g.logic, g.playerMgr, g.eventMgr)
	g.uiMgr = ui.NewImguiManager(g) //ui.NewUIMgr(g)
	g.playableScene = CreatePlayableScene(g.viewport)
	g.eventHandles = CreateEventHandles(g.net, g.logic, g.playableScene, &g.gameData)
	g.inputMgr = NewInputMgr(g.cmdMgr)
	return g
}

// 初始化
func (g *Game) Init() error {
	g.uiMgr.Init()
	g.msgHandler.Init()
	g.eventHandles.Init()
	g.initInputHandles()
	g.restart()
	return nil
}

// 反初始化
func (g *Game) Uninit() {
	g.eventHandles.Uninit()
}

// 重新开始
func (g *Game) restart() {
	g.gameData.State = client_base.GameStateBeforeLogin
}

// 当前模式
func (g *Game) GetState() client_base.GameState {
	return g.gameData.State
}

// 獲得游戲數據
func (g *Game) GetGameData() *client_base.GameData {
	return &g.gameData
}

// 事件管理器
func (g *Game) EventMgr() base.IEventManager {
	return g.eventMgr
}

// 命令管理器
func (g *Game) CmdMgr() *client_core.CmdHandleManager {
	return g.cmdMgr
}

// 布局
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

// 更新逻辑
func (g *Game) Update() error {
	err := g.net.Update()
	if err != nil {
		return err
	}

	switch g.gameData.State {
	case client_base.GameStateBeforeLogin:
	case client_base.GameStateInLogin:
	case client_base.GameStateMainMenu:
	case client_base.GameStateEnteringWorld:
	case client_base.GameStateInWorld:
		g.update()
		g.inputMgr.HandleInput()
	case client_base.GameStateExitingWorld:
		//	g.restart()
	}
	g.uiMgr.Update()
	return nil
}

// 绘制
func (g *Game) Draw(screen *ebiten.Image) {
	if g.gameData.State == client_base.GameStateInWorld && g.logic.IsStart() {
		// 画场景
		g.playableScene.Draw(screen)
		// 显示帧数
		ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f", ebiten.ActualTPS()))
	}
	// 画UI
	g.uiMgr.Draw(screen)
}

// 屏幕寬高
func (g *Game) ScreenWidthHeight() (int32, int32) {
	return screenWidth, screenHeight
}

// 更新
func (g *Game) update() {
	if !g.logic.IsStart() {
		// g.loadMap()
		g.logic.Start()
		return
	}
	// 时间同步完成
	if core.IsTimeSyncEnd() {
		now := core.GetSyncServTime()
		if g.lastCheckTime.IsZero() {
			g.lastCheckTime = now
		} else {
			tick := now.Sub(g.lastCheckTime)
			for ; tick >= common_data.GameLogicTick; tick -= common_data.GameLogicTick {
				g.logic.Update(common_data.GameLogicTick)
				g.lastCheckTime = g.lastCheckTime.Add(common_data.GameLogicTick)
			}
		}
	}
}

func (g *Game) initInputHandles() {
	g.cmdMgr.Add(CMD_CAMERA_UP, func(...any) {
		g.playableScene.CameraMove(0, 10)
	})
	g.cmdMgr.Add(CMD_CAMERA_DOWN, func(...any) {
		g.playableScene.CameraMove(0, -10)
	})
	g.cmdMgr.Add(CMD_CAMERA_LEFT, func(...any) {
		g.playableScene.CameraMove(-10, 0)
	})
	g.cmdMgr.Add(CMD_CAMERA_RIGHT, func(...any) {
		g.playableScene.CameraMove(10, 0)
	})
	g.cmdMgr.Add(CMD_CAMERA_HEIGHT, func(args ...any) {
		delta := args[0].(int)
		g.playableScene.CameraChangeHeight(int32(delta))
	})
}
