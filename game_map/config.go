package game_map

import (
	"project_b/common/base"
	"project_b/common/object"
)

type BotInfo struct {
	TankId int32
	Camp   base.CampType
	Pos    base.Pos
	Level  int32
}

// 地图配置，坐标系x轴朝上y轴朝右
type Config struct {
	Id                 int32                 // Id
	Name               string                // 地圖名
	Layers             [][]int16             // 地图数据
	X, Y               int32                 // 左下角坐标
	TileWidth          int32                 // 瓦片寬度
	TileHeight         int32                 // 瓦片高度
	PlayerTankInitData object.TankStaticInfo // 玩家坦克配置信息
	PlayerTankInitRect base.Rect             // 玩家坦克出现位置范围矩形
	PlayerMaxCount     int32                 // 最大玩家数
	PlayerBornPosList  []base.Pos            // 玩家出生點列表
	BotInfoList        []BotInfo             // bot列表
}
