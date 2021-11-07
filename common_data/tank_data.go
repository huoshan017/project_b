package common_data

import (
	"project_b/common/object"
)

var (
	ObjectConfigData = map[int32]*object.ObjStaticInfo{}

	// 坦克id列表
	TankIdList = []int32{1, 2, 1000, 1001, 1002}

	// 坦克最大等级
	TankMaxLevelMap = map[int32]int32{
		1:    4,
		2:    4,
		1000: 2,
		1001: 2,
		1002: 4,
	}

	// 坦克静态配置数据
	TankConfigData = map[int32]*object.ObjStaticInfo{
		1: object.NewObjStaticInfo(1, object.ObjTypeMovable, object.ObjSubTypeTank, 0, 0, 28, 28, 25, object.DirUp),
		2: object.NewObjStaticInfo(2, object.ObjTypeMovable, object.ObjSubTypeTank, 0, 0, 28, 28, 25, object.DirUp),

		1000: object.NewObjStaticInfo(1000, object.ObjTypeMovable, object.ObjSubTypeTank, 0, 0, 28, 28, 25, object.DirUp),
		1001: object.NewObjStaticInfo(1001, object.ObjTypeMovable, object.ObjSubTypeTank, 0, 0, 28, 28, 65, object.DirUp),
		1002: object.NewObjStaticInfo(1002, object.ObjTypeMovable, object.ObjSubTypeTank, 0, 0, 28, 28, 30, object.DirUp),
	}

	// 玩家坦克配置信息
	PlayerTankInitData = TankConfigData[1]
	// 玩家坦克出现位置范围矩形
	PlayerTankInitRect = object.Rect{LeftTop: object.Pos{X: 0, Y: 0}, RightBottom: object.Pos{X: 200, Y: 200}}

	// 同伴坦克配置信息
	TeammateTankInfoData = TankConfigData[2]
	// 同伴坦克初始位置
	TeammateTankInfoPos = object.Pos{X: 150, Y: 100}

	// 敌人坦克配置信息
	EnemyTankInitData = map[int32]*object.ObjStaticInfo{
		1000: TankConfigData[1000],
		1001: TankConfigData[1001],
	}
)
