package main

import (
	"fmt"
	client_base "project_b/client/base"
	core "project_b/client_core"
	"project_b/common/base"
	"project_b/common/time"
	"project_b/common_data"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/huoshan017/gsnet/options"
)

type GameState int

const (
	GameStateMainMenu GameState = iota
	GameStateInGame
	GameStateOver
)

type GameData struct {
	state GameState
	myId  uint64 // 我的ID
	myAcc string // 我的帐号
}

type Game struct {
	//---------------------------------------
	// 逻辑
	net           *core.NetClient        // 网络模块
	msgHandler    *core.MsgHandler       // 消息处理器
	logic         *core.GameLogic        // 游戏逻辑
	cmdMgr        *core.CmdHandleManager // 命令处理管理器
	playerMgr     *core.CPlayerManager   // 玩家管理器
	eventMgr      *base.EventManager     // 游戏事件管理器，向上层逻辑传递事件
	lastCheckTime time.CustomTime        // 上次检测时间

	//---------------------------------------
	// 表现相关

	viewport      *client_base.Viewport // 视口
	camera        *client_base.Camera   // 摄像机
	playableScene *PlayableScene        // 场景绘制
	uiMgr         *UIManager            // UI管理
	eventHandles  *EventHandles         // 事件处理
	gameData      GameData              // 其他游戏数据
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
	g.uiMgr = NewUIMgr(g)
	g.playableScene = CreatePlayableScene()
	g.camera = client_base.CreateCamera(g.viewport, conf.cameraFov)
	g.eventHandles = CreateEventHandles(g.net, g.logic, g.playableScene, &g.gameData)
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
	//g.logic.LoadMap(mapIdList[g.logic.MapIndex()])
	g.gameData.state = 0
}

// 当前模式
func (g *Game) GetState() GameState {
	return g.gameData.state
}

// 事件管理器
func (g *Game) EventMgr() base.IEventManager {
	return g.eventMgr
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

	switch g.gameData.state {
	case GameStateMainMenu:
	case GameStateInGame:
		if !g.logic.IsStart() {
			g.loadMap()
			g.logic.Start()
		} else {
			g.update()
		}
		g.handleInput()
	case GameStateOver:
		g.restart()
		g.gameData.state = GameStateMainMenu
	}
	g.uiMgr.Update()
	return nil
}

// 绘制
func (g *Game) Draw(screen *ebiten.Image) {
	if g.gameData.state == GameStateInGame && g.logic.IsStart() {
		// 画场景
		g.camera.Draw(screen)
		// 显示帧数
		ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f", ebiten.ActualTPS()))
	}
	// 画UI
	g.uiMgr.Draw(screen)
}

// 载入地图
func (g *Game) loadMap() {
	mapId := mapIdList[g.logic.MapIndex()]
	if g.logic.LoadMap(mapId) {
		g.camera.SetScene(g.playableScene)
		mapInfo := mapInfoArray[mapId]
		g.camera.MoveTo(mapInfo.cameraPos.X, mapInfo.cameraPos.Y)
		g.camera.ChangeHeight(mapInfo.cameraHeight)
	}
}

// 更新
func (g *Game) update() {
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

// 处理输入
func (g *Game) handleInput() {
	pressedKey := inpututil.AppendPressedKeys(nil)
	for _, pk := range pressedKey {
		if c, o := keyPressed2CmdMap[pk]; o {
			g.cmdMgr.Handle(c.cmd, c.args...)
		}
	}
	for k, cmd := range keyReleased2CmdMap {
		if inpututil.IsKeyJustReleased(k) {
			g.cmdMgr.Handle(cmd)
		}
	}
}

func (g *Game) initInputHandles() {
	g.cmdMgr.Add(CMD_CAMERA_UP, func(args ...any) {
		g.camera.Move(0, 10)
	})
	g.cmdMgr.Add(CMD_CAMERA_DOWN, func(args ...any) {
		g.camera.Move(0, -10)
	})
	g.cmdMgr.Add(CMD_CAMERA_LEFT, func(args ...any) {
		g.camera.Move(-10, 0)
	})
	g.cmdMgr.Add(CMD_CAMERA_RIGHT, func(args ...any) {
		g.camera.Move(10, 0)
	})
}
