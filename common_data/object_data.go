package common_data

import (
	"project_b/common/base"
	"project_b/common/object"
)

var (
	// 静态物体类型列表
	StaticObjectTypeList = []object.StaticObjType{
		object.StaticObjBrick,
		object.StaticObjIron,
		object.StaticObjTree,
		object.StaticObjWater,
		object.StaticObjIce,
		object.StaticObjHome,
		object.StaticObjRuins,
	}

	// 静态物体配置
	StaticObjectConfigData = map[object.StaticObjType]*object.ObjStaticInfo{
		// 砖
		object.StaticObjBrick: object.NewObjStaticInfo(
			int32(object.StaticObjBrick), object.ObjTypeStatic, object.ObjSubtypeBrick, object.CampTypeNone, 0, 0, 320, 320, 0, 0, 1, true),
		// 铁
		object.StaticObjIron: object.NewObjStaticInfo(
			int32(object.StaticObjIron), object.ObjTypeStatic, object.ObjSubtypeIron, object.CampTypeNone, 0, 0, 320, 320, 0, 0, 1, true),
		// 樹
		object.StaticObjTree: object.NewObjStaticInfo(
			int32(object.StaticObjTree), object.ObjTypeStatic, object.ObjSubtypeTree, object.CampTypeNone, 0, 0, 320, 320, 0, 0, 2, false),
		// 水
		object.StaticObjWater: object.NewObjStaticInfo(
			int32(object.StaticObjWater), object.ObjTypeStatic, object.ObjSubtypeWater, object.CampTypeNone, 0, 0, 320, 320, 0, 0, 0, true),
		// 冰
		object.StaticObjIce: object.NewObjStaticInfo(
			int32(object.StaticObjIce), object.ObjTypeStatic, object.ObjSubtypeIce, object.CampTypeNone, 0, 0, 320, 320, 0, 0, 0, false),
		// 基地
		object.StaticObjHome: object.NewObjStaticInfo(
			int32(object.StaticObjHome), object.ObjTypeStatic, object.ObjSubtypeHome, object.CampTypeNone, 0, 0, 320, 320, 0, 0, 1, true),
		// 廢墟
		object.StaticObjRuins: object.NewObjStaticInfo(
			int32(object.StaticObjRuins), object.ObjTypeStatic, object.ObjSubtypeRuins, object.CampTypeNone, 0, 0, 320, 320, 0, 0, 1, true),
	}

	// 炮彈靜態配置
	ShellConfigData = map[int32]*object.ShellStaticInfo{
		1: {
			MovableObjStaticInfo: object.MovableObjStaticInfo{
				ObjStaticInfo: *object.NewObjStaticInfo(1, object.ObjTypeMovable, object.ObjSubtypeShell, object.CampTypeNone, 0, 0, 80, 80, 0, 1200, 1, true),
			},
			Range:       10000,
			Damage:      100,
			BlastRadius: 500,
		},
		2: {
			MovableObjStaticInfo: object.MovableObjStaticInfo{
				ObjStaticInfo: *object.NewObjStaticInfo(2, object.ObjTypeMovable, object.ObjSubtypeShell, object.CampTypeNone, 0, 0, 80, 160, 90, 2500, 1, true),
				MoveFunc:      object.ShellTrackMove,
			},
			Range:                   100000,
			Damage:                  2000,
			BlastRadius:             1000,
			TrackTarget:             true,
			SearchTargetRadius:      4000,
			SteeringAngularVelocity: 60 * 140, // 分
		},
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
	TankConfigData = map[int32]*object.TankStaticInfo{
		1: {
			MovableObjStaticInfo: object.MovableObjStaticInfo{
				ObjStaticInfo: *object.NewObjStaticInfo(1, object.ObjTypeMovable, object.ObjSubtypeTank, object.CampTypeNone, 0, 0, 280, 280, 0, 600, 1, true),
			},
			Level:                   1,
			SteeringAngularVelocity: 220 * 60,
			ShellLaunchPos:          base.NewVec2(160, 0),
			ShellConfig:             object.TankShellConfig{ShellId: 1, AmountFireOneTime: 1, IntervalInFire: 0, Cooldown: 2000},
		},
		2: {
			MovableObjStaticInfo: object.MovableObjStaticInfo{
				ObjStaticInfo: *object.NewObjStaticInfo(2, object.ObjTypeMovable, object.ObjSubtypeTank, object.CampTypeNone, 0, 0, 280, 280, 0, 600, 1, true),
			},
			Level:                   1,
			SteeringAngularVelocity: 220 * 60,
			ShellLaunchPos:          base.NewVec2(160, 0),
			ShellConfig:             object.TankShellConfig{ShellId: 1, AmountFireOneTime: 1, IntervalInFire: 0, Cooldown: 2000},
		},
		1000: {
			MovableObjStaticInfo: object.MovableObjStaticInfo{
				ObjStaticInfo: *object.NewObjStaticInfo(1000, object.ObjTypeMovable, object.ObjSubtypeTank, object.CampTypeNone, 0, 0, 280, 280, 0, 600, 1, true),
			},
			Level:                   1,
			SteeringAngularVelocity: 240 * 60,
			ShellLaunchPos:          base.NewVec2(160, 0),
			ShellConfig:             object.TankShellConfig{ShellId: 1, AmountFireOneTime: 1, IntervalInFire: 0, Cooldown: 1500},
		},
		1001: {
			MovableObjStaticInfo: object.MovableObjStaticInfo{
				ObjStaticInfo: *object.NewObjStaticInfo(1001, object.ObjTypeMovable, object.ObjSubtypeTank, object.CampTypeNone, 0, 0, 280, 280, 0, 900, 1, true),
			},
			Level:                   1,
			SteeringAngularVelocity: 250 * 60,
			ShellLaunchPos:          base.NewVec2(165, 0),
			ShellConfig:             object.TankShellConfig{ShellId: 1, AmountFireOneTime: 2, IntervalInFire: 500, Cooldown: 1000},
		},
		1002: {
			MovableObjStaticInfo: object.MovableObjStaticInfo{
				ObjStaticInfo: *object.NewObjStaticInfo(1002, object.ObjTypeMovable, object.ObjSubtypeTank, object.CampTypeNone, 0, 0, 280, 280, 0, 650, 1, true),
			},
			Level:                   1,
			SteeringAngularVelocity: 200 * 60,
			ShellLaunchPos:          base.NewVec2(165, 0),
			ShellConfig:             object.TankShellConfig{ShellId: 1, AmountFireOneTime: 2, IntervalInFire: 500, Cooldown: 1000},
		},
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
	EnemyTankInitData = map[int32]*object.TankStaticInfo{
		1000: TankConfigData[1000],
		1001: TankConfigData[1001],
	}

	// 效果配置信息
	EffectConfigData = map[int32]*object.EffectStaticInfo{
		1: {
			Id:     1,
			Et:     object.EffectTypeRequency,
			Param:  1,
			Width:  280,
			Height: 280,
		},
		2: {
			Id:     2,
			Et:     object.EffectTypeTime,
			Param:  500,
			Width:  640,
			Height: 640,
		},
	}

	// todo 環繞物體配置信息，測試用
	SurroundObjConfigData = map[int32]*object.SurroundObjStaticInfo{
		1: {
			MovableObjStaticInfo: object.MovableObjStaticInfo{
				ObjStaticInfo: *object.NewObjStaticInfo(
					1, object.ObjTypeMovable, object.ObjSubtypeSurroundObj, object.CampTypeNone, 0, 0, 80, 80, 0, 0, 1, true),
				MoveFunc: object.SurroundObjMove,
			},
			AroundRadius:    600,
			AngularVelocity: 100 * 60, // 單位: 分(1/60度)每秒
			Clockwise:       false,
		},
	}
)
