package main

import (
	"fmt"
	core "project_b/client_core"
	"project_b/common"
	"project_b/common/base"
	"project_b/common/time"
	"project_b/common_data"

	"golang.org/x/image/math/f64"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type GameState int

const (
	GameStateMainMenu GameState = iota
	GameStateInGame
	GameStateOver
)

type Game struct {
	state GameState

	//---------------------------------------
	// 逻辑
	net           *core.NetClient      // 网络模块
	msgHandler    *core.MsgHandler     // 消息处理器
	logic         *common.GameLogic    // 游戏逻辑
	playerMgr     *core.CPlayerManager // 玩家管理器
	eventMgr      base.IEventManager   // 游戏事件管理器，向上层逻辑传递事件
	lastCheckTime time.CustomTime      // 上次检测时间

	//---------------------------------------
	// 表现相关
	cmdMgr      *CmdHandleManager // 命令处理管理器
	camera      *Camera           // 摄像机
	currMap     *Map              // 当前地图资源
	uiMgr       *UIManager        // UI管理
	playableMgr *PlayableManager  // 可播放管理器
	myId        uint64            // 我的ID
	myAcc       string            // 我的帐号

	// --------------------------------------
	// 事件处理
	gameEvent2Handles   []event2Handle // 游戏事件
	playerEvent2Handles []event2Handle // 玩家事件
}

// 创建游戏
func NewGame() *Game {
	g := &Game{
		camera: &Camera{ViewPort: f64.Vec2{screenWidth, screenHeight}}, // 相机的视口范围与窗口屏幕大小一样
	}
	g.playerMgr = core.CreateCPlayerManager()
	g.eventMgr = base.NewEventManager()
	g.logic = common.NewGameLogic()
	g.cmdMgr = CreateCmdHandleManager(g)
	g.uiMgr = NewUIMgr(g)
	g.uiMgr.Init()
	g.playableMgr = CreatePlayableManager()
	return g
}

// 初始化
func (g *Game) Init(conf *Config) error {
	g.restart()
	g.net = core.CreateNetClient(conf.ServerAddress)
	g.msgHandler = core.CreateMsgHandler(g.net, g.logic, g.playerMgr, g.eventMgr)
	g.msgHandler.Init()
	g.registerEvents()
	return nil
}

// 反初始化
func (g *Game) Uninit() {
	g.unregisterEvents()
}

// 重新开始
func (g *Game) restart() {
	g.logic.LoadMap(0)
}

// 当前模式
func (g *Game) GetState() GameState {
	return g.state
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
	if !g.net.IsDisconnected() {
		err := g.net.Update()
		if err != nil {
			return err
		}
	}

	switch g.state {
	case GameStateMainMenu:
	case GameStateInGame:
		if !g.logic.IsStart() {
			g.loadMap()
			g.logic.Start()
		} else {
			// 时间同步完成
			if core.IsTimeSyncEnd() {
				now := core.GetSyncServTime() //time.Now()
				if g.lastCheckTime.IsZero() {
					g.lastCheckTime = now
				} else {
					tick := now.Sub(g.lastCheckTime)
					for ; tick >= common_data.GameLogicTick; tick -= common_data.GameLogicTick {
						g.logic.Update(common_data.GameLogicTick)
						g.lastCheckTime = g.lastCheckTime.Add(common_data.GameLogicTick)
					}
				}
				g.handleInput()
			}
		}
	case GameStateOver:
		g.restart()
		g.state = GameStateMainMenu
	}
	g.uiMgr.Update()
	return nil
}

// 绘制
func (g *Game) Draw(screen *ebiten.Image) {
	if g.state == GameStateInGame {
		// 画场景
		g.drawScene(screen)
	}

	// 画UI
	g.uiMgr.Draw(screen)
}

// 载入地图
func (g *Game) loadMap() {
	if g.currMap == nil {
		g.currMap = &Map{}
		g.currMap.Load(g.logic.MapIndex())
	}
}

// 画场景
func (g *Game) drawScene(screen *ebiten.Image) {
	if g.currMap != nil {
		// 先画地图场景
		g.currMap.Draw()
		// 再画物体
		g.playableMgr.Update(g.currMap.GetImage())
		// 渲染到屏幕
		g.camera.Render(g.currMap.GetImage(), screen)
		// 显示帧数
		ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f", ebiten.CurrentTPS()))
	}
}

// 处理输入
func (g *Game) handleInput() {
	pressedKey := inpututil.PressedKeys()
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
