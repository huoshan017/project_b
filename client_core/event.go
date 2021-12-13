package client_core

import (
	"project_b/common/base"
)

const (
	//EventIdNone = common.EventId(iota)

	/* UI操作事件 */
	EventIdOpLogin     = base.EventId(1)
	EventIdOpEnterGame = base.EventId(2)

	/* 时间同步事件 */
	EventIdTimeSync    = base.EventId(10)
	EventIdTimeSyncEnd = base.EventId(11)

	/* 网络协议事件 */
	EventIdPlayerEnterGame          = base.EventId(100) // 进入游戏 参数：Account(string), PlayerId(uint64), TankInfo()
	EventIdPlayerEnterGameCompleted = base.EventId(101) // 进入游戏完成
	EventIdPlayerExitGame           = base.EventId(102) // 离开游戏

	/* 游戏逻辑事件 */
	EventIdTankMove     = base.EventId(200)  // 移动事件
	EventIdTankStopMove = base.EventId(201)  // 停止移动事件
	EventIdTankMoveSync = base.EventId(202)  // 移动同步事件
	EventIdTankChange   = base.EventId(1000) // 改变坦克
	EventIdTankRestore  = base.EventId(1001) // 恢复坦克
)