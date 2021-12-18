package main

import (
	"project_b/common"
	"project_b/common/object"
	custom_time "project_b/common/time"
	"project_b/common_data"

	"time"

	"github.com/huoshan017/gproc"
)

// 游戏逻辑处理器
type GameLogicHandler struct {
	gproc.RequestHandler
	logic *common.GameLogic
}

// 创建游戏逻辑处理器
func CreateGameLogicHandler() *GameLogicHandler {
	return &GameLogicHandler{
		RequestHandler: *gproc.NewDefaultRequestHandler(),
		logic:          common.NewGameLogic(),
	}
}

// 初始化，注册消息函数
func (h *GameLogicHandler) Init() {
	// 坦克进入同步
	h.RegisterHandle("playerTankEnterSync", h.onPlayerTankEnterSync)
	// 坦克离开同步
	h.RegisterHandle("playerTankLeaveSync", h.onPlayerTankLeaveSync)
	// 坦克移动同步
	h.RegisterHandle("playerTankMoveSync", h.onPlayerTankMoveSync)
	// 改变坦克
	h.RegisterHandle("playerTankChange", h.onPlayerChangeTank)
	// tick处理
	h.SetTickHandle(h.onTick, time.Millisecond*10)
}

// 玩家坦克进入同步
func (h *GameLogicHandler) onPlayerTankEnterSync(sender gproc.ISender, args interface{}) {
	msg, o := args.(*MsgPlayerTankEnterSync)
	if !o {
		gslog.Fatal("Must msg type %v", args)
		return
	}
	h.logic.PlayerTankEnter(msg.playerId, object.NewTank(msg.playerId, common_data.PlayerTankInitData))
}

// 玩家坦克离开同步
func (h *GameLogicHandler) onPlayerTankLeaveSync(sender gproc.ISender, args interface{}) {
	msg, o := args.(*MsgPlayerTankLeaveSync)
	if !o {
		gslog.Fatal("Must msg type: MsgPlayerTankLeaveSync")
		return
	}
	h.logic.PlayerTankLeave(msg.playerId)
}

// 玩家坦克移动同步
func (h *GameLogicHandler) onPlayerTankMoveSync(sender gproc.ISender, args interface{}) {
	msg, o := args.(*MsgPlayerTankMoveSync)
	if !o {
		gslog.Fatal("Must msg type: MsgPlayerTankMoveSync")
		return
	}
	h.logic.PlayerTankMove(msg.playerId, msg.direction)
}

// 玩家改变坦克
func (h *GameLogicHandler) onPlayerChangeTank(sender gproc.ISender, args interface{}) {
	msg, o := args.(*MsgPlayerTankChange)
	if !o {
		gslog.Fatal("Must msg type: MsgPlayerChangeTank")
		return
	}
	// todo 暂时把改变坦克的静态信息结构设为nil，在上一层消息处理函数中做改变坦克的逻辑
	h.logic.PlayerTankChange(msg.playerId, nil)
}

// tick处理
func (h *GameLogicHandler) onTick(tick time.Duration) {
	h.logic.Update(custom_time.Duration(tick))
}
