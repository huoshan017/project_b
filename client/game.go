package main

import (
	"fmt"
	"project_b/client/ui"
	"project_b/client_base"
	"project_b/client_core"
	"project_b/common/base"
	"project_b/common/time"
	"project_b/core"
	"project_b/game_map"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/huoshan017/gsnet/options"
)

var debug = client_base.GetDebug()

type Game struct {
	conf *Config
	//---------------------------------------
	// 逻辑
	gameCore        *core.GameCore              // 游戲實例
	net             *client_core.NetClient      // 网络模块
	msgHandler      *client_core.MsgHandler     // 消息处理器
	playerMgr       *client_core.CPlayerManager // 玩家管理器，這裏的玩家是指獨立於游戲邏輯GameLogic之外的登錄用戶
	eventMgr        *base.EventManager          // 游戏事件管理器，向上层逻辑传递事件
	lastCheckTime   time.CustomTime             // 上次检测时间
	isStartInstance bool                        // 開始實例
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
		conf:     conf,
		viewport: client_base.CreateViewport(0, 0, screenWidth, screenHeight),
	}
	g.eventMgr = base.NewEventManager()
	g.gameCore = core.NewGameCore(core.CoreArgs{EventMgr: g.eventMgr, PlayerNum: conf.playerMaxCount, FrameMs: conf.frameMs, GetMapFunc: func(mapId int32) *game_map.Config {
		return game_map.MapConfigArray[mapId]
	}})
	g.net = client_core.CreateNetClient(conf.serverAddress, options.WithRunMode(options.RunModeOnlyUpdate), options.WithNoDelay(true))
	g.playerMgr = client_core.CreateCPlayerManager()
	g.msgHandler = client_core.CreateMsgHandler(g.net, g.playerMgr, g.eventMgr)
	g.uiMgr = ui.NewImguiManager(g)
	g.playableScene = CreatePlayableScene(g.viewport)
	g.eventHandles = CreateEventHandles(g.net, g.gameCore, g.playableScene, &g.gameData)
	g.inputMgr = NewInputMgr(g, g.gameCore)
	return g
}

// 初始化
func (g *Game) Init() error {
	g.uiMgr.Init()
	g.msgHandler.Init()
	g.eventHandles.Init()
	g.initInputHandles()
	g.restart()
	time.InitTime()
	return nil
}

// 反初始化
func (g *Game) Uninit() {
	g.eventHandles.Uninit()
}

// 重新开始
func (g *Game) restart() {
	g.gameData.State = client_base.GameStateMainMenu //client_base.GameStateBeforeLogin
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

// 游戲核心
func (g *Game) GameCore() *core.GameCore {
	return g.gameCore
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
	case client_base.GameStateInReplay:
		g.updateInReplay()
	case client_base.GameStateExitingWorld:
	}
	g.uiMgr.Update()
	return nil
}

// 绘制
func (g *Game) Draw(screen *ebiten.Image) {
	if g.gameData.State == client_base.GameStateInWorld || g.gameData.State == client_base.GameStateInReplay {
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
	// 时间同步完成
	/*if !core.IsTimeSyncEnd() {
		return
	}
	now := core.GetSyncServTime()*/
	g.inputMgr.HandleInput()
	g.gameCore.Update(time.CurrentMs())
}

// 重播更新
func (g *Game) updateInReplay() {
	g.inputMgr.HandleInput()
	g.gameCore.Update(time.CurrentMs())
}

func (g *Game) initInputHandles() {
	g.inputMgr.AddKeyHandle(CMD_CAMERA_UP, func(args []int64) {
		delta := args[0]
		g.playableScene.CameraMove(0, int32(delta))
	})
	g.inputMgr.AddKeyHandle(CMD_CAMERA_DOWN, func(args []int64) {
		delta := args[0]
		g.playableScene.CameraMove(0, int32(delta))
	})
	g.inputMgr.AddKeyHandle(CMD_CAMERA_LEFT, func(args []int64) {
		delta := args[0]
		g.playableScene.CameraMove(int32(delta), 0)
	})
	g.inputMgr.AddKeyHandle(CMD_CAMERA_RIGHT, func(args []int64) {
		delta := args[0]
		g.playableScene.CameraMove(int32(delta), 0)
	})
	g.inputMgr.AddKeyHandle(CMD_CAMERA_HEIGHT, func(args []int64) {
		delta := args[0]
		g.playableScene.CameraChangeHeight(-int32(delta))
	})
	g.inputMgr.SetWheelHandle(func(xoffset, yoffset float64) {
		delta := int32(yoffset * 30)
		g.playableScene.CameraChangeHeight(-delta)
	})
}
