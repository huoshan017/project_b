package main

import (
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
	CMD_RESTART      CmdCode = 100 // 重新开始
)

// 命令对应处理
type Cmd2Handle struct {
	cmd    CmdCode
	handle func(args ...interface{})
}

// 命令处理管理器
type CmdHandleManager struct {
	game    *Game
	handles []Cmd2Handle
}

// 创建命令处理管理器
func CreateCmdHandleManager(game *Game) *CmdHandleManager {
	m := &CmdHandleManager{
		game: game,
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

// 处理命令
func (m *CmdHandleManager) Handle(cmd CmdCode, args ...interface{}) {
	found := false
	for _, h := range m.handles {
		if h.cmd == cmd {
			h.handle(args...)
			found = true
			break
		}
	}
	if !found {
		getLog().Warn("not found handle for cmd %v", cmd)
	}
}

// 移动命令
func (m *CmdHandleManager) handleMove(args ...interface{}) {
	dir := args[0].(object.Direction)
	m.game.logic.PlayerTankMove(m.game.myId, dir)
	m.game.eventMgr.InvokeEvent(EventIdTankMove)
}

// 停止移动命令
func (m *CmdHandleManager) handleStopMove(args ...interface{}) {
	m.game.logic.PlayerTankStopMove(m.game.myId)
	m.game.eventMgr.InvokeEvent(EventIdTankStopMove)
}

// 改变坦克命令
func (m *CmdHandleManager) handleChangeTank(args ...interface{}) {
	err := m.game.net.SendTankChangeReq()
	if err != nil {
		getLog().Warn("send tank change req err: %v", err)
	}
}

// 恢复坦克命令
func (m *CmdHandleManager) handleRestoreTank(args ...interface{}) {
	err := m.game.net.SendTankRestoreReq()
	if err != nil {
		getLog().Warn("send tank restore req err: %v", err)
	}
}
