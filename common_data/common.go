package common_data

import "project_b/common/time"

var (
	GameLogicTick = time.Millisecond * 100
	// 地图Id列表
	MapIdList = []int32{1, 2, 3}
	// 默認bot搜索半徑
	DefaultSearchRadius int32 = 3000
)
