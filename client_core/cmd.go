package client_core

import (
	"project_b/common/log"
	"project_b/common/object"
)

// 命令编码
type CmdCode int32

// 命令枚举
const (
	CMD_NONE               CmdCode = iota
	CMD_MOVE               CmdCode = 1   // 移动
	CMD_STOP_MOVE          CmdCode = 2   // 停止移动
	CMD_CHANGE_TANK        CmdCode = 3   // 改变坦克
	CMD_RESTORE_TANK       CmdCode = 4   // 恢复坦克
	CMD_FIRE               CmdCode = 5   // 開炮
	CMD_ROTATE             CmdCode = 6   // 旋轉
	CMD_REVIVE             CmdCode = 7   // 復活
	CMD_RESTART            CmdCode = 100 // 重新开始
	CMD_RELEASE_SMALL_BALL CmdCode = 999 // 釋放小球 測試用
)

// 命令处理管理器
type CmdHandleManager struct {
	net     *NetClient
	logic   *GameLogic
	handles map[CmdCode]func(args ...any)
}

// 创建命令处理管理器
func CreateCmdHandleManager(net *NetClient, logic *GameLogic) *CmdHandleManager {
	m := &CmdHandleManager{
		net:   net,
		logic: logic,
	}
	handles := map[CmdCode]func(args ...any){
		CMD_MOVE:               m.handleMove,
		CMD_STOP_MOVE:          m.handleStopMove,
		CMD_CHANGE_TANK:        m.handleChangeTank,
		CMD_RESTORE_TANK:       m.handleRestoreTank,
		CMD_FIRE:               m.handleTankFire,
		CMD_ROTATE:             m.handleTankRotate,
		CMD_REVIVE:             m.handleTankRevive,
		CMD_RELEASE_SMALL_BALL: m.handleTankReleaseSurroundObj,
	}
	m.handles = handles
	return m
}

func (m *CmdHandleManager) Add(cmd CmdCode, handle func(args ...any)) {
	m.handles[cmd] = handle
}

// 处理命令
func (m *CmdHandleManager) Handle(cmd CmdCode, args ...any) {
	if handle, o := m.handles[cmd]; o {
		handle(args...)
	} else {
		Log().Warn("not found handle for cmd %v", cmd)
	}
}

// 移动命令
func (m *CmdHandleManager) handleMove(args ...any) {
	dir := args[0].(object.Direction)
	orientation := object.Dir2Orientation(dir)
	m.logic.MyPlayerTankMove(orientation)
	//log.Debug("handleMove dir %v", dir)
}

// 停止移动命令
func (m *CmdHandleManager) handleStopMove(args ...any) {
	m.logic.MyPlayerTankStopMove()
	//log.Debug("handleStopMove")
}

// 改变坦克命令
func (m *CmdHandleManager) handleChangeTank(args ...any) {
	err := m.net.SendTankChangeReq()
	if err != nil {
		log.Warn("send tank change req err: %v", err)
	}
}

// 恢复坦克命令
func (m *CmdHandleManager) handleRestoreTank(args ...any) {
	err := m.net.SendTankRestoreReq()
	if err != nil {
		log.Warn("send tank restore req err: %v", err)
	}
}

// 坦克開火
func (m *CmdHandleManager) handleTankFire(args ...any) {
	shellId := args[0].(int)
	m.logic.MyPlayerTankFire(int32(shellId))
}

// 坦克釋放小球
func (m *CmdHandleManager) handleTankReleaseSurroundObj(args ...any) {
	m.logic.MyPlayerTankReleaseSurroundObj()
}

// 坦克旋轉
func (m *CmdHandleManager) handleTankRotate(args ...any) {
	angle := args[0].(int)
	m.logic.MyPlayerTankRotate(int32(angle))
}

// 坦克復活
func (m *CmdHandleManager) handleTankRevive(args ...any) {
	m.logic.MyPlayerTankRevive()
}
