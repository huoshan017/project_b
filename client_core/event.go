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
	EventIdPlayerEnterGame          = base.EventId(100) // 进入游戏
	EventIdPlayerEnterGameCompleted = base.EventId(101) // 进入游戏完成
	EventIdPlayerExitGame           = base.EventId(102) // 离开游戏
)
