package weapon

type LaserStaticInfo struct {
	Diameter       int32 // 直徑
	Range          int32 // 射程
	Dps            int32 // 每秒傷害
	Speed          int32 // 速度(秒)
	Energy         int32 // 能量
	CostPerSecond  int32 // 每秒消耗能量
	ChargPerSecond int32 // 每秒充能值 (充能速度)
}
