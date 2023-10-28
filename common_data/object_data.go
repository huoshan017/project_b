package common_data

import (
	"project_b/common/base"
	"project_b/common/effect"
	"project_b/common/object"
	"project_b/common/weapon"
)

var (
	// 静态物体类型列表
	StaticObjectTypeList = []base.StaticObjType{
		base.StaticObjBrick,
		base.StaticObjIron,
		base.StaticObjTree,
		base.StaticObjWater,
		base.StaticObjIce,
		base.StaticObjHome,
		base.StaticObjRuins,
	}

	// 静态物体配置
	StaticObjectConfigData = map[base.StaticObjType]*object.ObjStaticInfo{
		// 砖
		base.StaticObjBrick: object.NewObjStaticInfo(
			int32(base.StaticObjBrick), base.ObjTypeStatic, base.ObjSubtypeBrick, base.CampTypeNone, 0, 0, 320, 320, 0, 0, 1, true),
		// 铁
		base.StaticObjIron: object.NewObjStaticInfo(
			int32(base.StaticObjIron), base.ObjTypeStatic, base.ObjSubtypeIron, base.CampTypeNone, 0, 0, 320, 320, 0, 0, 1, true),
		// 樹
		base.StaticObjTree: object.NewObjStaticInfo(
			int32(base.StaticObjTree), base.ObjTypeStatic, base.ObjSubtypeTree, base.CampTypeNone, 0, 0, 320, 320, 0, 0, 2, false),
		// 水
		base.StaticObjWater: object.NewObjStaticInfo(
			int32(base.StaticObjWater), base.ObjTypeStatic, base.ObjSubtypeWater, base.CampTypeNone, 0, 0, 320, 320, 0, 0, 0, true),
		// 冰
		base.StaticObjIce: object.NewObjStaticInfo(
			int32(base.StaticObjIce), base.ObjTypeStatic, base.ObjSubtypeIce, base.CampTypeNone, 0, 0, 320, 320, 0, 0, 0, false),
		// 基地
		base.StaticObjHome: object.NewObjStaticInfo(
			int32(base.StaticObjHome), base.ObjTypeStatic, base.ObjSubtypeHome, base.CampTypeNone, 0, 0, 320, 320, 0, 0, 1, true),
		// 廢墟
		base.StaticObjRuins: object.NewObjStaticInfo(
			int32(base.StaticObjRuins), base.ObjTypeStatic, base.ObjSubtypeRuins, base.CampTypeNone, 0, 0, 320, 320, 0, 0, 1, true),
	}

	// 炮彈靜態配置
	ShellConfigData = map[int32]*object.ShellStaticInfo{
		1: {
			MovableObjStaticInfo: object.MovableObjStaticInfo{
				ObjStaticInfo: *object.NewObjStaticInfo(1, base.ObjTypeMovable, base.ObjSubtypeShell, base.CampTypeNone, 0, 0, 80, 80, 0, 1200, 1, true),
			},
			Range:       10000,
			Damage:      100,
			BlastRadius: 500,
		},
		2: {
			MovableObjStaticInfo: object.MovableObjStaticInfo{
				ObjStaticInfo: *object.NewObjStaticInfo(2, base.ObjTypeMovable, base.ObjSubtypeShell, base.CampTypeNone, 0, 0, 80, 160, 90, 2500, 1, true),
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
				ObjStaticInfo: *object.NewObjStaticInfo(1, base.ObjTypeMovable, base.ObjSubtypeTank, base.CampTypeNone, 0, 0, 280, 280, 0, 600, 1, true),
			},
			Level:                   1,
			SteeringAngularVelocity: 460 * 60,
			ShellLaunchPos:          base.NewVec2(160, 0),
			ShellConfig:             object.TankShellConfig{ShellInfo: ShellConfigData[1], AmountFireOneTime: 1, IntervalInFireMs: 1000, CooldownMs: 3000},
			ShieldConfig:            object.TankShieldStaticInfo{},
			LaserId:                 1,
		},
		2: {
			MovableObjStaticInfo: object.MovableObjStaticInfo{
				ObjStaticInfo: *object.NewObjStaticInfo(2, base.ObjTypeMovable, base.ObjSubtypeTank, base.CampTypeNone, 0, 0, 280, 280, 0, 600, 1, true),
			},
			Level:                   1,
			SteeringAngularVelocity: 220 * 60,
			ShellLaunchPos:          base.NewVec2(160, 0),
			ShellConfig:             object.TankShellConfig{ShellInfo: ShellConfigData[1], AmountFireOneTime: 1, IntervalInFireMs: 1000, CooldownMs: 3000},
			LaserId:                 1,
		},
		1000: {
			MovableObjStaticInfo: object.MovableObjStaticInfo{
				ObjStaticInfo: *object.NewObjStaticInfo(1000, base.ObjTypeMovable, base.ObjSubtypeTank, base.CampTypeNone, 0, 0, 280, 280, 0, 600, 1, true),
			},
			Level:                   1,
			SteeringAngularVelocity: 240 * 60,
			ShellLaunchPos:          base.NewVec2(160, 0),
			ShellConfig:             object.TankShellConfig{ShellInfo: ShellConfigData[1], AmountFireOneTime: 1, IntervalInFireMs: 1500, CooldownMs: 1500},
			LaserId:                 1,
		},
		1001: {
			MovableObjStaticInfo: object.MovableObjStaticInfo{
				ObjStaticInfo: *object.NewObjStaticInfo(1001, base.ObjTypeMovable, base.ObjSubtypeTank, base.CampTypeNone, 0, 0, 280, 280, 0, 900, 1, true),
			},
			Level:                   1,
			SteeringAngularVelocity: 250 * 60,
			ShellLaunchPos:          base.NewVec2(165, 0),
			ShellConfig:             object.TankShellConfig{ShellInfo: ShellConfigData[1], AmountFireOneTime: 2, IntervalInFireMs: 500, CooldownMs: 1000},
			LaserId:                 1,
		},
		1002: {
			MovableObjStaticInfo: object.MovableObjStaticInfo{
				ObjStaticInfo: *object.NewObjStaticInfo(1002, base.ObjTypeMovable, base.ObjSubtypeTank, base.CampTypeNone, 0, 0, 280, 280, 0, 650, 1, true),
			},
			Level:                   1,
			SteeringAngularVelocity: 200 * 60,
			ShellLaunchPos:          base.NewVec2(165, 0),
			ShellConfig:             object.TankShellConfig{ShellInfo: ShellConfigData[1], AmountFireOneTime: 2, IntervalInFireMs: 500, CooldownMs: 1000},
			LaserId:                 1,
		},
	}

	// 玩家坦克配置信息
	PlayerTankInitData = TankConfigData[1]
	// 玩家坦克出现位置范围矩形，坐标和位置都是逻辑数据
	PlayerTankInitRect = base.Rect{LeftBottom: base.Pos{X: 0, Y: 0}, RightTop: base.Pos{X: 2000, Y: 2000}}

	// 同伴坦克配置信息
	TeammateTankInfoData = TankConfigData[2]
	// 同伴坦克初始位置，逻辑数据
	TeammateTankInfoPos = base.Pos{X: 1500, Y: 1000}

	// 敌人坦克配置信息
	EnemyTankInitData = map[int32]*object.TankStaticInfo{
		1000: TankConfigData[1000],
		1001: TankConfigData[1001],
	}

	// 效果配置信息
	EffectConfigData = map[int32]*effect.EffectStaticInfo{
		1: {
			Id:     1,
			Et:     effect.EffectTypeRequency,
			Param:  1,
			Width:  280,
			Height: 280,
		},
		2: {
			Id:     2,
			Et:     effect.EffectTypeTime,
			Param:  200,
			Width:  640,
			Height: 640,
		},
	}

	// todo 環繞物體配置信息，測試用
	SurroundObjConfigData = map[int32]*object.SurroundObjStaticInfo{
		1: {
			MovableObjStaticInfo: object.MovableObjStaticInfo{
				ObjStaticInfo: *object.NewObjStaticInfo(
					1, base.ObjTypeMovable, base.ObjSubtypeSurroundObj, base.CampTypeNone, 0, 0, 80, 80, 0, 0, 1, true),
				MoveFunc: object.SurroundObjMove,
			},
			AroundRadius:    600,
			AngularVelocity: 100 * 60, // 單位: 分(1/60度)每秒
			Clockwise:       false,
		},
	}

	// 坦克護盾配置信息
	TankShieldConfigData = map[int32]*object.TankShieldStaticInfo{
		1: {
			Width: 32, Length: 32,
			DurationMs: 0,
		},
		2: {
			Width: 32, Length: 32,
			DurationMs: 10,
		},
	}

	// 激光配置信息
	LaserConfigData = map[int32]*weapon.LaserStaticInfo{
		1: {
			Diameter:       1,
			Range:          2500,
			Dps:            100,
			Energy:         100,
			CostPerSecond:  40,
			ChargPerSecond: 200,
		},
	}
)
