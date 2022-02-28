package common_data

import (
	"project_b/common/object"
)

var (
	// 静态物体类型列表
	StaticObjectTypeList = []object.StaticObjType{
		object.StaticObjBrick,
		object.StaticObjIron,
		object.StaticObjGrass,
		object.StaticObjWater,
		object.StaticObjIce,
		object.StaticObjHome,
		object.StaticObjHomeDestroyed,
	}

	// 静态物体配置
	StaticObjectConfigData = map[object.StaticObjType]*object.ObjStaticInfo{
		// 砖
		object.StaticObjBrick: object.NewObjStaticInfo(int32(object.StaticObjBrick), object.ObjTypeStatic, object.ObjSubTypeBrick, 0, 0, 320, 320, 0, object.DirNone),
		// 铁
		object.StaticObjIron: object.NewObjStaticInfo(int32(object.StaticObjIron), object.ObjTypeStatic, object.ObjSubTypeIron, 0, 0, 320, 320, 0, object.DirNone),
		// 草
		object.StaticObjGrass: object.NewObjStaticInfo(int32(object.StaticObjGrass), object.ObjTypeStatic, object.ObjSubTypeGrass, 0, 0, 320, 320, 0, object.DirNone),
		// 水
		object.StaticObjWater: object.NewObjStaticInfo(int32(object.StaticObjWater), object.ObjTypeStatic, object.ObjSubTypeWater, 0, 0, 320, 320, 0, object.DirNone),
		// 冰
		object.StaticObjIce: object.NewObjStaticInfo(int32(object.StaticObjIce), object.ObjTypeStatic, object.ObjSubTypeIce, 0, 0, 320, 320, 0, object.DirNone),
		// 基地
		object.StaticObjHome: object.NewObjStaticInfo(int32(object.StaticObjHome), object.ObjTypeStatic, object.ObjSubTypeHome, 0, 0, 320, 320, 0, object.DirNone),
		// 摧毁的基地
		object.StaticObjHomeDestroyed: object.NewObjStaticInfo(int32(object.StaticObjHomeDestroyed), object.ObjTypeStatic, object.ObjSubTypeHomeDestroyed, 0, 0, 320, 320, 0, object.DirNone),
	}

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

	// 坦克静态配置数据，全部是逻辑数据
	TankConfigData = map[int32]*object.ObjStaticInfo{
		1: object.NewObjStaticInfo(1, object.ObjTypeMovable, object.ObjSubTypeTank, 0, 0, 280, 280, 250, object.DirUp),
		2: object.NewObjStaticInfo(2, object.ObjTypeMovable, object.ObjSubTypeTank, 0, 0, 280, 280, 250, object.DirUp),

		1000: object.NewObjStaticInfo(1000, object.ObjTypeMovable, object.ObjSubTypeTank, 0, 0, 280, 280, 250, object.DirUp),
		1001: object.NewObjStaticInfo(1001, object.ObjTypeMovable, object.ObjSubTypeTank, 0, 0, 280, 280, 650, object.DirUp),
		1002: object.NewObjStaticInfo(1002, object.ObjTypeMovable, object.ObjSubTypeTank, 0, 0, 280, 280, 300, object.DirUp),
	}

	// 玩家坦克配置信息
	PlayerTankInitData = TankConfigData[1]
	// 玩家坦克出现位置范围矩形，坐标和位置都是逻辑数据
	PlayerTankInitRect = object.Rect{LeftBottom: object.Pos{X: 0, Y: 0}, RightTop: object.Pos{X: 2000, Y: 2000}}

	// 同伴坦克配置信息
	TeammateTankInfoData = TankConfigData[2]
	// 同伴坦克初始位置，逻辑数据
	TeammateTankInfoPos = object.Pos{X: 1500, Y: 1000}

	// 敌人坦克配置信息
	EnemyTankInitData = map[int32]*object.ObjStaticInfo{
		1000: TankConfigData[1000],
		1001: TankConfigData[1001],
	}
)