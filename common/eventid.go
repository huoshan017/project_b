package common

import "project_b/common/base"

const (
	/* 游戏逻辑事件 */
	EventIdTankMove     = base.EventId(200)  // 移动事件
	EventIdTankStopMove = base.EventId(201)  // 停止移动事件
	EventIdTankSetPos   = base.EventId(202)  // 设置坐标事件
	EventIdTankChange   = base.EventId(1000) // 改变坦克
	EventIdTankRestore  = base.EventId(1001) // 恢复坦克
)
