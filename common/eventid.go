package common

import "project_b/common/base"

const (
	EventIdBeforeMapLoad   = base.EventId(150) // 地图逻辑载入前
	EventIdMapLoaded       = base.EventId(151) // 地图逻辑载入完成
	EventIdBeforeMapUnload = base.EventId(152) // 地图卸载前
	EventIdMapUnloaded     = base.EventId(153) // 地图卸载后
	/* 游戏逻辑事件 */
	EventIdTankMove     = base.EventId(200)  // 移动事件
	EventIdTankStopMove = base.EventId(201)  // 停止移动事件
	EventIdTankSetPos   = base.EventId(202)  // 设置坐标事件
	EventIdTankChange   = base.EventId(1000) // 改变坦克
	EventIdTankRestore  = base.EventId(1001) // 恢复坦克
)
