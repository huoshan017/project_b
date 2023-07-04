package game_map

import "project_b/common/object"

type BotInfo struct {
	TankId int32
	Camp   object.CampType
	Pos    object.Pos
	Level  int8
}

// 地图配置，坐标系x轴朝上y轴朝右
type Config struct {
	Id                 int32                 // Id
	Layers             [][]int16             // 地图数据
	X, Y               int32                 // 左下角坐标
	TileWidth          int32                 // 瓦片寬度
	TileHeight         int32                 // 瓦片高度
	PlayerTankInitData object.TankStaticInfo // 玩家坦克配置信息
	PlayerTankInitRect object.Rect           // 玩家坦克出现位置范围矩形
	PlayerMaxCount     int32                 // 最大玩家数
	BotInfoList        []BotInfo             // bot列表
}
