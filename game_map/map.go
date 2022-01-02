package game_map

import (
	"project_b/common/object"
)

type Config struct {
	Layers             [][]int              // 地图数据
	Width              int                  // 宽度
	Height             int                  // 长度
	PlayerTankInitData object.ObjStaticInfo // 玩家坦克配置信息
	PlayerTankInitRect object.Rect          // 玩家坦克出现位置范围矩形
	PlayerMaxCount     int                  // 最大玩家数
}
