package client_core

import (
	"project_b/common/log"
	"project_b/common/object"
)

// 命令编码
type CmdCode int32

// 命令枚举
const (
	CMD_NONE         CmdCode = iota
	CMD_MOVE         CmdCode = 1   // 移动
	CMD_STOP_MOVE    CmdCode = 2   // 停止移动
	CMD_CHANGE_TANK  CmdCode = 3   // 改变坦克
	CMD_RESTORE_TANK CmdCode = 4   // 恢复坦克
	CMD_FIRE         CmdCode = 5   // 開炮
	CMD_RESTART      CmdCode = 100 // 重新开始
)

// 命令对应处理
type Cmd2Handle struct {
	cmd    CmdCode
	handle func(args ...any)
}

// 命令处理管理器
type CmdHandleManager struct {
	net     *NetClient
	logic   *GameLogic
	handles []Cmd2Handle
}

// 创建命令处理管理器
func CreateCmdHandleManager(net *NetClient, logic *GameLogic) *CmdHandleManager {
	m := &CmdHandleManager{
		net:   net,
		logic: logic,
	}
	handles := []Cmd2Handle{
		{CMD_MOVE, m.handleMove},
		{CMD_STOP_MOVE, m.handleStopMove},
		{CMD_CHANGE_TANK, m.handleChangeTank},
		{CMD_RESTORE_TANK, m.handleRestoreTank},
	}
	m.handles = handles
	return m
}

func (m *CmdHandleManager) Add(cmd CmdCode, handle func(args ...any)) {
	m.handles = append(m.handles, Cmd2Handle{cmd: cmd, handle: handle})
}

// 处理命令
func (m *CmdHandleManager) Handle(cmd CmdCode, args ...any) {
	found := false
	for _, h := range m.handles {
		if h.cmd == cmd {
			h.handle(args...)
			found = true
			break
		}
	}
	if !found {
		Log().Warn("not found handle for cmd %v", cmd)
	}
}

// 移动命令
func (m *CmdHandleManager) handleMove(args ...any) {
	dir := args[0].(object.Direction)
	m.logic.MyPlayerTankMove(dir)
	// todo 在GameLogic中触发移动事件
	//m.game.eventMgr.InvokeEvent(EventIdTankMove)
}

// 停止移动命令
func (m *CmdHandleManager) handleStopMove(args ...any) {
	log.Debug("handleStopMove")
	m.logic.MyPlayerTankStopMove()
	// todo 在GameLogic中触发停止事件
	//m.game.eventMgr.InvokeEvent(EventIdTankStopMove)
}

// 改变坦克命令
func (m *CmdHandleManager) handleChangeTank(args ...any) {
	err := m.net.SendTankChangeReq()
	if err != nil {
		Log().Warn("send tank change req err: %v", err)
	}
}

// 恢复坦克命令
func (m *CmdHandleManager) handleRestoreTank(args ...any) {
	err := m.net.SendTankRestoreReq()
	if err != nil {
		Log().Warn("send tank restore req err: %v", err)
	}
}
