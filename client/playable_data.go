package main

import (
	"project_b/client/images"
	client_base "project_b/client_base"
	"project_b/common/base"
)

type frameConfig struct {
	frameNum       int32
	dirMap         map[base.Direction]int32
	frameLevelList [][]int32
}

var (
	mapTileAnimOriginalConfig     map[base.StaticObjType]*client_base.SpriteAnimConfig // 地图瓦片动画初始配置
	mapTileFrameConfig            map[base.StaticObjType]*frameConfig                  // 地图瓦片动画初始配置
	tankAnimOriginalConfig        map[int32]*client_base.SpriteAnimConfig              // 坦克动画初始配置map
	tankFrameConfig               map[int32]*frameConfig                               // 帧配置map
	bulletAnimOriginalConfig      map[int32]*client_base.SpriteAnimConfig              // 子弹动画初始配置
	bulletFrameConfig             *frameConfig                                         // 子弹帧配置
	effectAnimOriginalConfig      map[int32]*client_base.SpriteAnimConfig              // 爆炸動畫初始配置
	surroundObjAnimOriginalConfig map[int32]*client_base.SpriteAnimConfig              // 小球動畫初始配置
	itemObjAnimOriginalConfig     map[base.ItemObjType]*client_base.SpriteAnimConfig   // 物品動畫初始配置
	shieldAnimOriginalConfig      *client_base.SpriteAnimConfig                        // 護盾動畫初始配置
)

// 初始化动画初始配置
func initAnimOriginalConfigs() {
	mapTileAnimOriginalConfig = map[base.StaticObjType]*client_base.SpriteAnimConfig{
		// 砖块
		base.StaticObjBrick: {
			Image: images.GetTileImg(), FrameWidth: 32, FrameHeight: 32, FramePosList: []client_base.SpriteIndex{{X: 0, Y: 0}},
		},
		// 铁
		base.StaticObjIron: {
			Image: images.GetTileImg(), FrameWidth: 32, FrameHeight: 32, FramePosList: []client_base.SpriteIndex{{X: 1, Y: 0}},
		},
		// 樹
		base.StaticObjTree: {
			Image: images.GetTileImg(), FrameWidth: 32, FrameHeight: 32, FramePosList: []client_base.SpriteIndex{{X: 2, Y: 0}},
		},
		// 水
		base.StaticObjWater: {
			Image: images.GetTileImg(), PlayInterval: 200, FrameWidth: 32, FrameHeight: 32, FramePosList: []client_base.SpriteIndex{{X: 3, Y: 0}, {X: 4, Y: 0}},
		},
		// 基地
		base.StaticObjHome: {
			Image: images.GetTileImg(), FrameWidth: 32, FrameHeight: 32, FramePosList: []client_base.SpriteIndex{{X: 5, Y: 0}},
		},
		// 摧毁的基地
		base.StaticObjRuins: {
			Image: images.GetTileImg(), FrameWidth: 32, FrameHeight: 32, FramePosList: []client_base.SpriteIndex{{X: 6, Y: 0}},
		},
	}
	mapTileFrameConfig = map[base.StaticObjType]*frameConfig{
		// 砖块
		base.StaticObjBrick: {
			frameNum:       1,
			frameLevelList: [][]int32{{0}},
		},
		// 铁
		base.StaticObjIron: {
			frameNum:       1,
			frameLevelList: [][]int32{{0}},
		},
		// 樹
		base.StaticObjTree: {
			frameNum:       1,
			frameLevelList: [][]int32{{0}},
		},
		// 水
		base.StaticObjWater: {
			frameNum:       2,
			frameLevelList: [][]int32{{0, 1}},
		},
		// 基地
		base.StaticObjHome: {
			frameNum:       1,
			frameLevelList: [][]int32{{0}},
		},
		// 被摧毁的基地
		base.StaticObjRuins: {
			frameNum:       1,
			frameLevelList: [][]int32{{0}},
		},
	}
	tankAnimOriginalConfig = map[int32]*client_base.SpriteAnimConfig{
		// 玩家坦克动画配置
		1: {
			Image: images.GetPlayer1Img(), PlayInterval: 200, FrameWidth: 28, FrameHeight: 28,
			FramePosList: []client_base.SpriteIndex{
				{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 2, Y: 0}, {X: 3, Y: 0}, {X: 4, Y: 0}, {X: 5, Y: 0}, {X: 6, Y: 0}, {X: 7, Y: 0},
				{X: 0, Y: 1}, {X: 1, Y: 1}, {X: 2, Y: 1}, {X: 3, Y: 1}, {X: 4, Y: 1}, {X: 5, Y: 1}, {X: 6, Y: 1}, {X: 7, Y: 1},
				{X: 0, Y: 2}, {X: 1, Y: 2}, {X: 2, Y: 2}, {X: 3, Y: 2}, {X: 4, Y: 2}, {X: 5, Y: 2}, {X: 6, Y: 2}, {X: 7, Y: 2},
				{X: 0, Y: 3}, {X: 1, Y: 3}, {X: 2, Y: 3}, {X: 3, Y: 3}, {X: 4, Y: 3}, {X: 5, Y: 3}, {X: 6, Y: 3}, {X: 7, Y: 3},
			},
		},
		// 同伴坦克动画配置
		2: {
			Image: images.GetPlayer2Img(), PlayInterval: 200, FrameWidth: 28, FrameHeight: 28,
			FramePosList: []client_base.SpriteIndex{
				{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 2, Y: 0}, {X: 3, Y: 0}, {X: 4, Y: 0}, {X: 5, Y: 0}, {X: 6, Y: 0}, {X: 7, Y: 0},
				{X: 0, Y: 1}, {X: 1, Y: 1}, {X: 2, Y: 1}, {X: 3, Y: 1}, {X: 4, Y: 1}, {X: 5, Y: 1}, {X: 6, Y: 1}, {X: 7, Y: 1},
				{X: 0, Y: 2}, {X: 1, Y: 2}, {X: 2, Y: 2}, {X: 3, Y: 2}, {X: 4, Y: 2}, {X: 5, Y: 2}, {X: 6, Y: 2}, {X: 7, Y: 2},
				{X: 0, Y: 3}, {X: 1, Y: 3}, {X: 2, Y: 3}, {X: 3, Y: 3}, {X: 4, Y: 3}, {X: 5, Y: 3}, {X: 6, Y: 3}, {X: 7, Y: 3},
			},
		},
		// 轻型坦克动画配置
		1000: {
			Image: images.GetEnemyImg(), PlayInterval: 200, FrameWidth: 28, FrameHeight: 28,
			FramePosList: []client_base.SpriteIndex{
				{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 2, Y: 0}, {X: 3, Y: 0},
				{X: 0, Y: 1}, {X: 1, Y: 1}, {X: 2, Y: 1}, {X: 3, Y: 1},
				{X: 0, Y: 2}, {X: 1, Y: 2}, {X: 2, Y: 2}, {X: 3, Y: 2},
				{X: 0, Y: 3}, {X: 1, Y: 3}, {X: 2, Y: 3}, {X: 3, Y: 3},
			},
		},
		// 快速坦克动画配置
		1001: {
			Image: images.GetEnemyImg(), PlayInterval: 200, FrameWidth: 28, FrameHeight: 28,
			FramePosList: []client_base.SpriteIndex{
				{X: 4, Y: 0}, {X: 5, Y: 0}, {X: 6, Y: 0}, {X: 7, Y: 0},
				{X: 4, Y: 1}, {X: 5, Y: 1}, {X: 6, Y: 1}, {X: 7, Y: 1},
				{X: 4, Y: 2}, {X: 5, Y: 2}, {X: 6, Y: 2}, {X: 7, Y: 2},
				{X: 4, Y: 3}, {X: 5, Y: 3}, {X: 6, Y: 3}, {X: 7, Y: 3},
			},
		},
		// 重型坦克动画配置
		1002: {
			Image: images.GetEnemyImg(), PlayInterval: 200, FrameWidth: 28, FrameHeight: 28,
			FramePosList: []client_base.SpriteIndex{
				{X: 0, Y: 4}, {X: 1, Y: 4}, {X: 2, Y: 4}, {X: 3, Y: 4}, {X: 4, Y: 4}, {X: 5, Y: 4}, {X: 6, Y: 4}, {X: 7, Y: 4},
				{X: 0, Y: 5}, {X: 1, Y: 5}, {X: 2, Y: 5}, {X: 3, Y: 5}, {X: 4, Y: 5}, {X: 5, Y: 5}, {X: 6, Y: 5}, {X: 7, Y: 5},
				{X: 0, Y: 6}, {X: 1, Y: 6}, {X: 2, Y: 6}, {X: 3, Y: 6}, {X: 4, Y: 6}, {X: 5, Y: 6}, {X: 6, Y: 6}, {X: 7, Y: 6},
				{X: 0, Y: 7}, {X: 1, Y: 7}, {X: 2, Y: 7}, {X: 3, Y: 7}, {X: 4, Y: 7}, {X: 5, Y: 7}, {X: 6, Y: 7}, {X: 7, Y: 7},
			},
		},
	}

	tankFrameConfig = map[int32]*frameConfig{
		1: {
			frameNum: 2,
			dirMap: map[base.Direction]int32{
				base.DirUp:    0,
				base.DirDown:  2,
				base.DirLeft:  3,
				base.DirRight: 1,
			},
			frameLevelList: [][]int32{
				{0, 1, 8, 9, 16, 17, 24, 25},   // 1级
				{2, 3, 10, 11, 18, 19, 26, 27}, // 2级
				{4, 5, 12, 13, 20, 21, 28, 29}, // 3级
				{6, 7, 14, 15, 22, 23, 30, 31}, // 4级
			},
		},
		2: {
			frameNum: 2,
			dirMap: map[base.Direction]int32{
				base.DirUp:    0,
				base.DirDown:  2,
				base.DirLeft:  3,
				base.DirRight: 1,
			},
			frameLevelList: [][]int32{
				{0, 1, 8, 9, 16, 17, 24, 25},   // 1级
				{2, 3, 10, 11, 18, 19, 26, 27}, // 2级
				{4, 5, 12, 13, 20, 21, 28, 29}, // 3级
				{6, 7, 14, 15, 22, 23, 30, 31}, // 4级
			},
		},
		1000: {
			frameNum: 2,
			dirMap: map[base.Direction]int32{
				base.DirUp:    0,
				base.DirDown:  2,
				base.DirLeft:  3,
				base.DirRight: 1,
			},
			frameLevelList: [][]int32{
				{0, 1, 4, 5, 8, 9, 12, 13},   // 灰
				{2, 3, 6, 7, 10, 11, 14, 15}, // 红
			},
		},
		1001: {
			frameNum: 2,
			dirMap: map[base.Direction]int32{
				base.DirUp:    0,
				base.DirDown:  2,
				base.DirLeft:  3,
				base.DirRight: 1,
			},
			frameLevelList: [][]int32{
				{0, 1, 4, 5, 8, 9, 12, 13},   // 灰
				{2, 3, 6, 7, 10, 11, 14, 15}, // 红
			},
		},
		1002: {
			frameNum: 2,
			dirMap: map[base.Direction]int32{
				base.DirUp:    0,
				base.DirDown:  2,
				base.DirLeft:  3,
				base.DirRight: 1,
			},
			frameLevelList: [][]int32{
				{0, 1, 8, 9, 16, 17, 24, 25},   // 绿
				{2, 3, 10, 11, 18, 19, 26, 27}, // 黄
				{4, 5, 12, 13, 20, 21, 28, 29}, // 灰
				{6, 7, 14, 15, 22, 23, 30, 31}, // 红
			},
		},
	}

	// 子彈動畫原始配置
	bulletAnimOriginalConfig = map[int32]*client_base.SpriteAnimConfig{
		1: {
			Image: images.GetBulletImg(), FrameWidth: 8, FrameHeight: 8, FramePosList: []client_base.SpriteIndex{{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 2, Y: 0}, {X: 3, Y: 0}},
		},
		2: {
			Image: images.GetShellImg(), FrameWidth: 8, FrameHeight: 16, FramePosList: []client_base.SpriteIndex{{X: 0, Y: 0}},
		},
	}

	// 炮彈幀配置
	bulletFrameConfig = &frameConfig{
		frameNum: 1,
		dirMap: map[base.Direction]int32{
			base.DirUp:    0,
			base.DirDown:  2,
			base.DirLeft:  3,
			base.DirRight: 1,
		},
		frameLevelList: [][]int32{{0, 1, 2, 3}},
	}

	// 效果動畫配置
	effectAnimOriginalConfig = map[int32]*client_base.SpriteAnimConfig{
		1: {
			Image: images.GetExplode1Img(), FrameWidth: 28, FrameHeight: 28, FramePosList: []client_base.SpriteIndex{{X: 0, Y: 0}},
		},
		2: {
			Image: images.GetExplode2Img(), FrameWidth: 64, FrameHeight: 64, FramePosList: []client_base.SpriteIndex{{X: 0, Y: 0}},
		},
	}

	// 環繞物動畫配置
	surroundObjAnimOriginalConfig = map[int32]*client_base.SpriteAnimConfig{
		1: {
			Image: images.GetSmallBallImg(), FrameWidth: 8, FrameHeight: 8, FramePosList: []client_base.SpriteIndex{{X: 0, Y: 0}},
		},
	}

	// 物品動畫配置
	itemObjAnimOriginalConfig = map[base.ItemObjType]*client_base.SpriteAnimConfig{
		base.ItemObjRewardLife: {
			Image: images.GetBonusImg(), FrameWidth: 32, FrameHeight: 28, FramePosList: []client_base.SpriteIndex{{X: 0, Y: 0}},
		},
		base.ItemObjFrozen: {
			Image: images.GetBonusImg(), FrameWidth: 32, FrameHeight: 28, FramePosList: []client_base.SpriteIndex{{X: 30, Y: 0}},
		},
		base.ItemObjReinforcement: {
			Image: images.GetBonusImg(), FrameWidth: 32, FrameHeight: 28, FramePosList: []client_base.SpriteIndex{{X: 60, Y: 0}},
		},
		base.ItemObjBomb: {
			Image: images.GetBonusImg(), FrameWidth: 32, FrameHeight: 28, FramePosList: []client_base.SpriteIndex{{X: 90, Y: 0}},
		},
		base.ItemObjSelfUpgrade: {
			Image: images.GetBonusImg(), FrameWidth: 32, FrameHeight: 28, FramePosList: []client_base.SpriteIndex{{X: 120, Y: 0}},
		},
		base.ItemObjShield: {
			Image: images.GetBonusImg(), FrameWidth: 32, FrameHeight: 28, FramePosList: []client_base.SpriteIndex{{X: 150, Y: 0}},
		},
	}

	// 護盾動畫配置
	shieldAnimOriginalConfig = &client_base.SpriteAnimConfig{
		Image: images.GetShieldImg(), FrameWidth: 32, FrameHeight: 32, PlayInterval: 100,
		FramePosList: []client_base.SpriteIndex{{X: 0, Y: 0}, {X: 0, Y: 1}},
	}
}

// 静态物体动画配置
type StaticObjectAnimConfig struct {
	Type       base.ObjectType
	Subtype    base.ObjSubtype
	AnimConfig *client_base.SpriteAnimConfig
}

// 移动物体动画配置
type MovableObjectAnimConfig struct {
	Type       base.ObjectType
	Subtype    base.ObjSubtype
	Id         int32
	Level      int32
	AnimConfig *client_base.SpriteAnimConfig // 使用物体方向枚举做索引
}

// 物品動畫配置
type ItemObjectAnimConfig struct {
	Type       base.ObjectType
	Subtype    base.ObjSubtype
	AnimConfig *client_base.SpriteAnimConfig
}

var (
	staticObjectAnimConfigList   map[base.StaticObjType]*StaticObjectAnimConfig // 静态物体动画配置
	moveableObjectAnimConfigList []*MovableObjectAnimConfig                     // 移动物体动画配置
)

// 初始化动画配置
func initAnimConfigs() {
	staticObjectAnimConfigList = map[base.StaticObjType]*StaticObjectAnimConfig{
		base.StaticObjBrick: {
			Type: base.ObjTypeStatic, Subtype: base.ObjSubtypeBrick,
			AnimConfig: createStaticObjAnimConfig(base.StaticObjBrick),
		},
		base.StaticObjIron: {
			Type: base.ObjTypeStatic, Subtype: base.ObjSubtypeIron,
			AnimConfig: createStaticObjAnimConfig(base.StaticObjIron),
		},
		base.StaticObjTree: {
			Type: base.ObjTypeStatic, Subtype: base.ObjSubtypeTree,
			AnimConfig: createStaticObjAnimConfig(base.StaticObjTree),
		},
		base.StaticObjWater: {
			Type: base.ObjTypeStatic, Subtype: base.ObjSubtypeWater,
			AnimConfig: createStaticObjAnimConfig(base.StaticObjWater),
		},
		base.StaticObjHome: {
			Type: base.ObjTypeStatic, Subtype: base.ObjSubtypeHome,
			AnimConfig: createStaticObjAnimConfig(base.StaticObjHome),
		},
		base.StaticObjRuins: {
			Type: base.ObjTypeStatic, Subtype: base.ObjSubtypeRuins,
			AnimConfig: createStaticObjAnimConfig(base.StaticObjRuins),
		},
	}

	moveableObjectAnimConfigList = []*MovableObjectAnimConfig{
		{
			Type: base.ObjTypeMovable, Subtype: base.ObjSubtypeTank, Id: 1, Level: 1,
			AnimConfig: createTankAnimConfig(1, 1, base.DirRight),
		}, // 玩家坦克一级
		{
			Type: base.ObjTypeMovable, Subtype: base.ObjSubtypeTank, Id: 1, Level: 2,
			AnimConfig: createTankAnimConfig(1, 2, base.DirRight),
		}, // 玩家坦克二级
		{
			Type: base.ObjTypeMovable, Subtype: base.ObjSubtypeTank, Id: 1, Level: 3,
			AnimConfig: createTankAnimConfig(1, 3, base.DirRight),
		}, // 玩家坦克三级
		{
			Type: base.ObjTypeMovable, Subtype: base.ObjSubtypeTank, Id: 1, Level: 4,
			AnimConfig: createTankAnimConfig(1, 4, base.DirRight),
		}, // 玩家坦克四级
		{
			Type: base.ObjTypeMovable, Subtype: base.ObjSubtypeTank, Id: 2, Level: 1,
			AnimConfig: createTankAnimConfig(2, 1, base.DirRight),
		}, // 同伴坦克一级
		{
			Type: base.ObjTypeMovable, Subtype: base.ObjSubtypeTank, Id: 2, Level: 2,
			AnimConfig: createTankAnimConfig(2, 2, base.DirRight),
		}, // 同伴坦克二级
		{
			Type: base.ObjTypeMovable, Subtype: base.ObjSubtypeTank, Id: 2, Level: 3,
			AnimConfig: createTankAnimConfig(2, 3, base.DirRight),
		}, // 同伴坦克三级
		{
			Type: base.ObjTypeMovable, Subtype: base.ObjSubtypeTank, Id: 2, Level: 4,
			AnimConfig: createTankAnimConfig(2, 4, base.DirRight),
		}, // 同伴坦克四级
		{
			Type: base.ObjTypeMovable, Subtype: base.ObjSubtypeTank, Id: 1000, Level: 1,
			AnimConfig: createTankAnimConfig(1000, 1, base.DirRight),
		}, // 轻型坦克
		{
			Type: base.ObjTypeMovable, Subtype: base.ObjSubtypeTank, Id: 1000, Level: 2,
			AnimConfig: createTankAnimConfig(1000, 2, base.DirRight),
		}, // 轻型坦克红色
		{
			Type: base.ObjTypeMovable, Subtype: base.ObjSubtypeTank, Id: 1001, Level: 1,
			AnimConfig: createTankAnimConfig(1001, 1, base.DirRight),
		}, // 快速坦克
		{
			Type: base.ObjTypeMovable, Subtype: base.ObjSubtypeTank, Id: 1001, Level: 2,
			AnimConfig: createTankAnimConfig(1001, 2, base.DirRight),
		}, // 快速坦克红
		{
			Type: base.ObjTypeMovable, Subtype: base.ObjSubtypeTank, Id: 1002, Level: 1,
			AnimConfig: createTankAnimConfig(1002, 1, base.DirRight),
		}, // 重型坦克
		{
			Type: base.ObjTypeMovable, Subtype: base.ObjSubtypeTank, Id: 1002, Level: 2,
			AnimConfig: createTankAnimConfig(1002, 2, base.DirRight),
		}, // 重型坦克绿
		{
			Type: base.ObjTypeMovable, Subtype: base.ObjSubtypeTank, Id: 1002, Level: 3,
			AnimConfig: createTankAnimConfig(1002, 3, base.DirRight),
		}, // 重型坦克黄
		{
			Type: base.ObjTypeMovable, Subtype: base.ObjSubtypeTank, Id: 1002, Level: 4,
			AnimConfig: createTankAnimConfig(1002, 4, base.DirRight),
		}, // 重型坦克红
		{
			Type: base.ObjTypeMovable, Subtype: base.ObjSubtypeShell, Id: 1, Level: 0,
			AnimConfig: createBulletAnimConfig(1, base.DirRight),
		}, // 子弹
		{
			Type: base.ObjTypeMovable, Subtype: base.ObjSubtypeShell, Id: 2, Level: 0,
			AnimConfig: createBulletAnimConfig(2, base.DirNone),
		}, // 追蹤子彈
		{
			Type: base.ObjTypeMovable, Subtype: base.ObjSubtypeSurroundObj, Id: 1, Level: 0,
			AnimConfig: getSurroundObjAnimConfig(1),
			//},
		}, // 環繞物體
	}
}

func initResources() {
	images.InitImages()
	initAnimOriginalConfigs()
	initAnimConfigs()
}

// 根据静态物体类型创建瓦片动画配置
func createStaticObjAnimConfig(objStaticType base.StaticObjType) *client_base.SpriteAnimConfig {
	ac := mapTileAnimOriginalConfig[objStaticType]
	fc := mapTileFrameConfig[objStaticType]
	anim := &client_base.SpriteAnimConfig{
		Image:        ac.Image,
		FrameWidth:   ac.FrameWidth,
		FrameHeight:  ac.FrameHeight,
		PlayInterval: ac.PlayInterval,
	}
	for i := int32(0); i < int32(fc.frameNum); i++ {
		anim.FramePosList = append(anim.FramePosList, ac.FramePosList[fc.frameLevelList[0][i]])
	}
	return anim
}

// 根据等级和方向创建坦克动画配置
func createTankAnimConfig(id int32, level int32, dir base.Direction) *client_base.SpriteAnimConfig {
	anim := &client_base.SpriteAnimConfig{}
	anim0 := tankAnimOriginalConfig[id]
	frame0 := tankFrameConfig[id]
	anim.Image = anim0.Image
	anim.PlayInterval = anim0.PlayInterval
	anim.FrameWidth = anim0.FrameWidth
	anim.FrameHeight = anim0.FrameHeight
	for i := int32(0); i < int32(frame0.frameNum); i++ {
		dirIdx := frame0.dirMap[dir]
		idx := dirIdx*frame0.frameNum + i
		anim.FramePosList = append(anim.FramePosList, anim0.FramePosList[frame0.frameLevelList[level-1][idx]])
	}
	return anim
}

// 根据方向创建坦克动画配置
func createBulletAnimConfig(id int32, dir base.Direction) *client_base.SpriteAnimConfig {
	animConfig := bulletAnimOriginalConfig[id]
	anim := &client_base.SpriteAnimConfig{
		Image:       animConfig.Image,
		FrameWidth:  animConfig.FrameWidth,
		FrameHeight: animConfig.FrameHeight,
	}
	dirIdx, o := bulletFrameConfig.dirMap[dir]
	if !o {
		dirIdx = 0
	}
	anim.FramePosList = []client_base.SpriteIndex{animConfig.FramePosList[dirIdx]}
	return anim
}

// 获得可移动物体的动画配置
func GetMovableObjAnimConfig(movableObjType base.MovableObjType, id int32, level int32) *MovableObjectAnimConfig {
	var animConfig *MovableObjectAnimConfig
	for _, c := range moveableObjectAnimConfigList {
		if c.Subtype == base.ObjSubtype(movableObjType) && c.Id == id && c.Level == level {
			animConfig = c
			break
		}
	}
	return animConfig
}

// 获得坦克动画配置
func GetTankAnimConfig(id int32, level int32) *MovableObjectAnimConfig {
	return GetMovableObjAnimConfig(base.MovableObjTank, id, level)
}

// 获得子弹动画配置
func GetBulletAnimConfig(id int32) *MovableObjectAnimConfig {
	return GetMovableObjAnimConfig(base.MovableObjShell, id, 0)
}

// 获得静态物体动画配置
func GetStaticObjAnimConfig(staticObjType base.StaticObjType) *StaticObjectAnimConfig {
	return staticObjectAnimConfigList[staticObjType]
}

// 获得砖块动画配置
func GetBrickAnimConfig() *StaticObjectAnimConfig {
	return GetStaticObjAnimConfig(base.StaticObjBrick)
}

// 获得铁块动画配置
func GetIronAnimConfig() *StaticObjectAnimConfig {
	return GetStaticObjAnimConfig(base.StaticObjIron)
}

// 获得樹木动画配置
func GetTreeAnimConfig() *StaticObjectAnimConfig {
	return GetStaticObjAnimConfig(base.StaticObjTree)
}

// 获得水动画配置
func GetWaterAnimConfig() *StaticObjectAnimConfig {
	return GetStaticObjAnimConfig(base.StaticObjWater)
}

// 获得基地动画配置
func GetHomeAnimConfig() *StaticObjectAnimConfig {
	return GetStaticObjAnimConfig(base.StaticObjHome)
}

// 获得廢墟动画配置
func GetRuinsAnimConfig() *StaticObjectAnimConfig {
	return GetStaticObjAnimConfig(base.StaticObjRuins)
}

// 獲得效果動畫配置
func getEffectAnimConfig(id int32) *client_base.SpriteAnimConfig {
	return effectAnimOriginalConfig[id]
}

// 獲得小球動畫配置
func getSurroundObjAnimConfig(id int32) *client_base.SpriteAnimConfig {
	return surroundObjAnimOriginalConfig[id]
}

// 獲得物品動畫配置
func getItemObjAnimConfig(itemType base.ItemObjType) *client_base.SpriteAnimConfig {
	return itemObjAnimOriginalConfig[itemType]
}

// 獲得護盾動畫配置
func getShieldAnimConfig() *client_base.SpriteAnimConfig {
	return shieldAnimOriginalConfig
}
