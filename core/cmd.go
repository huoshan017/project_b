package core

// 命令编码
type CmdCode int32

// 命令枚举
const (
	CMD_NONE               CmdCode = iota
	CMD_TANK_MOVE          CmdCode = 1   // 移动
	CMD_TANK_STOP          CmdCode = 2   // 停止移动
	CMD_TANK_CHANGE        CmdCode = 3   // 改变坦克
	CMD_TANK_RESTORE       CmdCode = 4   // 恢复坦克
	CMD_TANK_FIRE          CmdCode = 5   // 開炮
	CMD_TANK_RESPAWN       CmdCode = 7   // 復活
	CMD_TANK_SHIELD        CmdCode = 8   // 添加或取消護盾
	CMD_RESTART            CmdCode = 100 // 重新开始
	CMD_RELEASE_SMALL_BALL CmdCode = 999 // 釋放小球 測試用
)

// 命令結構
type CmdData struct {
	cmd  CmdCode
	args []any
}

// 創建命令數據
func NewCmdData(cmd CmdCode, args []any) *CmdData {
	return &CmdData{cmd: cmd, args: args}
}

func (cd CmdData) Cmd() CmdCode {
	return cd.cmd
}

func (cd CmdData) Args() []any {
	return cd.args
}
