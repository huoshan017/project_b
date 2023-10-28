package common

import "project_b/common/base"

const (
	EventIdBeforeMapLoad   = base.EventId(150) // 地图逻辑载入前
	EventIdMapLoaded       = base.EventId(151) // 地图逻辑载入完成
	EventIdBeforeMapUnload = base.EventId(152) // 地图卸载前
	EventIdMapUnloaded     = base.EventId(153) // 地图卸载后
	EventIdEnterGame       = base.EventId(200) // 进入游戏
	EventIdExitGame        = base.EventId(201) // 退出游戲
	EventIdGameStart       = base.EventId(202) // 游戏开始
	EventIdGamePause       = base.EventId(203) // 游戏暂停
	EventIdGameOver        = base.EventId(204) // 游戏结束
	EventIdGameStatistic   = base.EventId(205) // 结算
	/* 游戏逻辑事件 */
	EventIdTankStartMove = base.EventId(200)  // 坦克开始移动事件
	EventIdTankStopMove  = base.EventId(201)  // 坦克停止移动事件
	EventIdTankShoot     = base.EventId(202)  // 坦克开炮事件
	EventIdHitObject     = base.EventId(203)  // 击中物体事件
	EventIdExplode       = base.EventId(204)  // 爆炸事件
	EventIdChangeTank    = base.EventId(1000) // 改变坦克
	EventIdRestoreTank   = base.EventId(1001) // 恢复坦克
	EventIdTankDestroy   = base.EventId(2000) // 坦克被擊毀
	EventIdTankRevive    = base.EventId(3000) // 坦克復活
)
