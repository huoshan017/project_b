package main

import (
	"project_b/client/base"
	"project_b/common/object"
)

type frameConfig struct {
	frameNum       int32
	dirMap         map[object.Direction]int32
	frameLevelList [][]int32
}

var (
	mapTileAnimOriginalConfig     map[object.StaticObjType]*base.SpriteAnimConfig // 地图瓦片动画初始配置
	mapTileFrameConfig            map[object.StaticObjType]*frameConfig           // 地图瓦片动画初始配置
	tankAnimOriginalConfig        map[int32]*base.SpriteAnimConfig                // 坦克动画初始配置map
	tankFrameConfig               map[int32]*frameConfig                          // 帧配置map
	bulletAnimOriginalConfig      *base.SpriteAnimConfig                          // 子弹动画初始配置
	bulletFrameConfig             *frameConfig                                    // 子弹帧配置
	effectAnimOriginalConfig      map[int32]*base.SpriteAnimConfig                // 爆炸動畫初始配置
	surroundObjAnimOriginalConfig map[int32]*base.SpriteAnimConfig                // 小球動畫初始配置
)

// 初始化动画初始配置
func initAnimOriginalConfigs() {
	mapTileAnimOriginalConfig = map[object.StaticObjType]*base.SpriteAnimConfig{
		// 砖块
		object.StaticObjBrick: {
			Image: tile_img, FrameWidth: 32, FrameHeight: 32, FramePosList: []base.SpriteIndex{{X: 0, Y: 0}},
		},
		// 铁
		object.StaticObjIron: {
			Image: tile_img, FrameWidth: 32, FrameHeight: 32, FramePosList: []base.SpriteIndex{{X: 1, Y: 0}},
		},
		// 樹
		object.StaticObjTree: {
			Image: tile_img, FrameWidth: 32, FrameHeight: 32, FramePosList: []base.SpriteIndex{{X: 2, Y: 0}},
		},
		// 水
		object.StaticObjWater: {
			Image: tile_img, PlayInterval: 200, FrameWidth: 32, FrameHeight: 32, FramePosList: []base.SpriteIndex{{X: 3, Y: 0}, {X: 4, Y: 0}},
		},
		// 基地
		object.StaticObjHome: {
			Image: tile_img, FrameWidth: 32, FrameHeight: 32, FramePosList: []base.SpriteIndex{{X: 5, Y: 0}},
		},
		// 摧毁的基地
		object.StaticObjRuins: {
			Image: tile_img, FrameWidth: 32, FrameHeight: 32, FramePosList: []base.SpriteIndex{{X: 6, Y: 0}},
		},
	}
	mapTileFrameConfig = map[object.StaticObjType]*frameConfig{
		// 砖块
		object.StaticObjBrick: {
			frameNum:       1,
			frameLevelList: [][]int32{{0}},
		},
		// 铁
		object.StaticObjIron: {
			frameNum:       1,
			frameLevelList: [][]int32{{0}},
		},
		// 樹
		object.StaticObjTree: {
			frameNum:       1,
			frameLevelList: [][]int32{{0}},
		},
		// 水
		object.StaticObjWater: {
			frameNum:       2,
			frameLevelList: [][]int32{{0, 1}},
		},
		// 基地
		object.StaticObjHome: {
			frameNum:       1,
			frameLevelList: [][]int32{{0}},
		},
		// 被摧毁的基地
		object.StaticObjRuins: {
			frameNum:       1,
			frameLevelList: [][]int32{{0}},
		},
	}
	tankAnimOriginalConfig = map[int32]*base.SpriteAnimConfig{
		// 玩家坦克动画配置
		1: {
			Image: player1_img, PlayInterval: 200, FrameWidth: 28, FrameHeight: 28,
			FramePosList: []base.SpriteIndex{
				{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 2, Y: 0}, {X: 3, Y: 0}, {X: 4, Y: 0}, {X: 5, Y: 0}, {X: 6, Y: 0}, {X: 7, Y: 0},
				{X: 0, Y: 1}, {X: 1, Y: 1}, {X: 2, Y: 1}, {X: 3, Y: 1}, {X: 4, Y: 1}, {X: 5, Y: 1}, {X: 6, Y: 1}, {X: 7, Y: 1},
				{X: 0, Y: 2}, {X: 1, Y: 2}, {X: 2, Y: 2}, {X: 3, Y: 2}, {X: 4, Y: 2}, {X: 5, Y: 2}, {X: 6, Y: 2}, {X: 7, Y: 2},
				{X: 0, Y: 3}, {X: 1, Y: 3}, {X: 2, Y: 3}, {X: 3, Y: 3}, {X: 4, Y: 3}, {X: 5, Y: 3}, {X: 6, Y: 3}, {X: 7, Y: 3},
			},
		},
		// 同伴坦克动画配置
		2: {
			Image: player2_img, PlayInterval: 200, FrameWidth: 28, FrameHeight: 28,
			FramePosList: []base.SpriteIndex{
				{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 2, Y: 0}, {X: 3, Y: 0}, {X: 4, Y: 0}, {X: 5, Y: 0}, {X: 6, Y: 0}, {X: 7, Y: 0},
				{X: 0, Y: 1}, {X: 1, Y: 1}, {X: 2, Y: 1}, {X: 3, Y: 1}, {X: 4, Y: 1}, {X: 5, Y: 1}, {X: 6, Y: 1}, {X: 7, Y: 1},
				{X: 0, Y: 2}, {X: 1, Y: 2}, {X: 2, Y: 2}, {X: 3, Y: 2}, {X: 4, Y: 2}, {X: 5, Y: 2}, {X: 6, Y: 2}, {X: 7, Y: 2},
				{X: 0, Y: 3}, {X: 1, Y: 3}, {X: 2, Y: 3}, {X: 3, Y: 3}, {X: 4, Y: 3}, {X: 5, Y: 3}, {X: 6, Y: 3}, {X: 7, Y: 3},
			},
		},
		// 轻型坦克动画配置
		1000: {
			Image: enemy_img, PlayInterval: 200, FrameWidth: 28, FrameHeight: 28,
			FramePosList: []base.SpriteIndex{
				{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 2, Y: 0}, {X: 3, Y: 0},
				{X: 0, Y: 1}, {X: 1, Y: 1}, {X: 2, Y: 1}, {X: 3, Y: 1},
				{X: 0, Y: 2}, {X: 1, Y: 2}, {X: 2, Y: 2}, {X: 3, Y: 2},
				{X: 0, Y: 3}, {X: 1, Y: 3}, {X: 2, Y: 3}, {X: 3, Y: 3},
			},
		},
		// 快速坦克动画配置
		1001: {
			Image: enemy_img, PlayInterval: 200, FrameWidth: 28, FrameHeight: 28,
			FramePosList: []base.SpriteIndex{
				{X: 4, Y: 0}, {X: 5, Y: 0}, {X: 6, Y: 0}, {X: 7, Y: 0},
				{X: 4, Y: 1}, {X: 5, Y: 1}, {X: 6, Y: 1}, {X: 7, Y: 1},
				{X: 4, Y: 2}, {X: 5, Y: 2}, {X: 6, Y: 2}, {X: 7, Y: 2},
				{X: 4, Y: 3}, {X: 5, Y: 3}, {X: 6, Y: 3}, {X: 7, Y: 3},
			},
		},
		// 重型坦克动画配置
		1002: {
			Image: enemy_img, PlayInterval: 200, FrameWidth: 28, FrameHeight: 28,
			FramePosList: []base.SpriteIndex{
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
			dirMap: map[object.Direction]int32{
				object.DirUp:    0,
				object.DirDown:  2,
				object.DirLeft:  3,
				object.DirRight: 1,
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
			dirMap: map[object.Direction]int32{
				object.DirUp:    0,
				object.DirDown:  2,
				object.DirLeft:  3,
				object.DirRight: 1,
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
			dirMap: map[object.Direction]int32{
				object.DirUp:    0,
				object.DirDown:  2,
				object.DirLeft:  3,
				object.DirRight: 1,
			},
			frameLevelList: [][]int32{
				{0, 1, 4, 5, 8, 9, 12, 13},   // 灰
				{2, 3, 6, 7, 10, 11, 14, 15}, // 红
			},
		},
		1001: {
			frameNum: 2,
			dirMap: map[object.Direction]int32{
				object.DirUp:    0,
				object.DirDown:  2,
				object.DirLeft:  3,
				object.DirRight: 1,
			},
			frameLevelList: [][]int32{
				{0, 1, 4, 5, 8, 9, 12, 13},   // 灰
				{2, 3, 6, 7, 10, 11, 14, 15}, // 红
			},
		},
		1002: {
			frameNum: 2,
			dirMap: map[object.Direction]int32{
				object.DirUp:    0,
				object.DirDown:  2,
				object.DirLeft:  3,
				object.DirRight: 1,
			},
			frameLevelList: [][]int32{
				{0, 1, 8, 9, 16, 17, 24, 25},   // 绿
				{2, 3, 10, 11, 18, 19, 26, 27}, // 黄
				{4, 5, 12, 13, 20, 21, 28, 29}, // 灰
				{6, 7, 14, 15, 22, 23, 30, 31}, // 红
			},
		},
	}

	bulletAnimOriginalConfig = &base.SpriteAnimConfig{
		Image: bullet_img, FrameWidth: 8, FrameHeight: 8, FramePosList: []base.SpriteIndex{{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 2, Y: 0}, {X: 3, Y: 0}},
	}

	bulletFrameConfig = &frameConfig{
		frameNum: 1,
		dirMap: map[object.Direction]int32{
			object.DirUp:    0,
			object.DirDown:  2,
			object.DirLeft:  3,
			object.DirRight: 1,
		},
		frameLevelList: [][]int32{{0, 1, 2, 3}},
	}

	effectAnimOriginalConfig = map[int32]*base.SpriteAnimConfig{
		1: {
			Image: explode1_img, FrameWidth: 28, FrameHeight: 28, FramePosList: []base.SpriteIndex{{X: 0, Y: 0}},
		},
		2: {
			Image: explode2_img, FrameWidth: 64, FrameHeight: 64, FramePosList: []base.SpriteIndex{{X: 0, Y: 0}},
		},
	}

	surroundObjAnimOriginalConfig = map[int32]*base.SpriteAnimConfig{
		1: {
			Image: smallball_img, FrameWidth: 8, FrameHeight: 8, FramePosList: []base.SpriteIndex{{X: 0, Y: 0}},
		},
	}
}

// 静态物体动画配置
type StaticObjectAnimConfig struct {
	Type       object.ObjectType
	Subtype    object.ObjSubType
	AnimConfig *base.SpriteAnimConfig
}

// 移动物体动画配置
type MovableObjectAnimConfig struct {
	Type       object.ObjectType
	Subtype    object.ObjSubType
	Id         int32
	Level      int32
	AnimConfig *base.SpriteAnimConfig // 使用物体方向枚举做索引
}

var (
	staticObjectAnimConfigList   map[object.StaticObjType]*StaticObjectAnimConfig // 静态物体动画配置
	moveableObjectAnimConfigList []*MovableObjectAnimConfig                       // 移动物体动画配置
)

// 初始化动画配置
func initAnimConfigs() {
	staticObjectAnimConfigList = map[object.StaticObjType]*StaticObjectAnimConfig{
		object.StaticObjBrick: {
			Type: object.ObjTypeStatic, Subtype: object.ObjSubTypeBrick,
			AnimConfig: createStaticObjAnimConfig(object.StaticObjBrick),
		},
		object.StaticObjIron: {
			Type: object.ObjTypeStatic, Subtype: object.ObjSubTypeIron,
			AnimConfig: createStaticObjAnimConfig(object.StaticObjIron),
		},
		object.StaticObjTree: {
			Type: object.ObjTypeStatic, Subtype: object.ObjSubTypeTree,
			AnimConfig: createStaticObjAnimConfig(object.StaticObjTree),
		},
		object.StaticObjWater: {
			Type: object.ObjTypeStatic, Subtype: object.ObjSubTypeWater,
			AnimConfig: createStaticObjAnimConfig(object.StaticObjWater),
		},
		object.StaticObjHome: {
			Type: object.ObjTypeStatic, Subtype: object.ObjSubTypeHome,
			AnimConfig: createStaticObjAnimConfig(object.StaticObjHome),
		},
		object.StaticObjRuins: {
			Type: object.ObjTypeStatic, Subtype: object.ObjSubTypeRuins,
			AnimConfig: createStaticObjAnimConfig(object.StaticObjRuins),
		},
	}

	moveableObjectAnimConfigList = []*MovableObjectAnimConfig{
		{
			Type: object.ObjTypeMovable, Subtype: object.ObjSubTypeTank, Id: 1, Level: 1,
			AnimConfig://[]*base.SpriteAnimConfig{
			// 空左右上下
			//nil,
			//createTankAnimConfig(1, 1, object.DirLeft),
			createTankAnimConfig(1, 1, object.DirRight),
			//createTankAnimConfig(1, 1, object.DirUp),
			//createTankAnimConfig(1, 1, object.DirDown),
			//},
		}, // 玩家坦克一级
		{
			Type: object.ObjTypeMovable, Subtype: object.ObjSubTypeTank, Id: 1, Level: 2,
			AnimConfig://[]*base.SpriteAnimConfig{
			//nil,
			//createTankAnimConfig(1, 2, object.DirLeft),
			createTankAnimConfig(1, 2, object.DirRight),
			//createTankAnimConfig(1, 2, object.DirUp),
			//createTankAnimConfig(1, 2, object.DirDown),
			//},
		}, // 玩家坦克二级
		{
			Type: object.ObjTypeMovable, Subtype: object.ObjSubTypeTank, Id: 1, Level: 3,
			AnimConfig://[]*base.SpriteAnimConfig{
			//nil,
			//createTankAnimConfig(1, 3, object.DirLeft),
			createTankAnimConfig(1, 3, object.DirRight),
			//createTankAnimConfig(1, 3, object.DirUp),
			//createTankAnimConfig(1, 3, object.DirDown),
			//},
		}, // 玩家坦克三级
		{
			Type: object.ObjTypeMovable, Subtype: object.ObjSubTypeTank, Id: 1, Level: 4,
			AnimConfig://[]*base.SpriteAnimConfig{
			//nil,
			//createTankAnimConfig(1, 4, object.DirLeft),
			createTankAnimConfig(1, 4, object.DirRight),
			//createTankAnimConfig(1, 4, object.DirUp),
			//createTankAnimConfig(1, 4, object.DirDown),
			//},
		}, // 玩家坦克四级
		{
			Type: object.ObjTypeMovable, Subtype: object.ObjSubTypeTank, Id: 2, Level: 1,
			AnimConfig://[]*base.SpriteAnimConfig{
			//nil,
			//createTankAnimConfig(2, 1, object.DirLeft),
			createTankAnimConfig(2, 1, object.DirRight),
			//createTankAnimConfig(2, 1, object.DirUp),
			//createTankAnimConfig(2, 1, object.DirDown),
			//},
		}, // 同伴坦克一级
		{
			Type: object.ObjTypeMovable, Subtype: object.ObjSubTypeTank, Id: 2, Level: 2,
			AnimConfig://[]*base.SpriteAnimConfig{
			//nil,
			//createTankAnimConfig(2, 2, object.DirLeft),
			createTankAnimConfig(2, 2, object.DirRight),
			//createTankAnimConfig(2, 2, object.DirUp),
			//createTankAnimConfig(2, 2, object.DirDown),
			//},
		}, // 同伴坦克二级
		{
			Type: object.ObjTypeMovable, Subtype: object.ObjSubTypeTank, Id: 2, Level: 3,
			AnimConfig://[]*base.SpriteAnimConfig{
			//nil,
			//createTankAnimConfig(2, 3, object.DirLeft),
			createTankAnimConfig(2, 3, object.DirRight),
			//createTankAnimConfig(2, 3, object.DirUp),
			//createTankAnimConfig(2, 3, object.DirDown),
			//},
		}, // 同伴坦克三级
		{
			Type: object.ObjTypeMovable, Subtype: object.ObjSubTypeTank, Id: 2, Level: 4,
			AnimConfig://[]*base.SpriteAnimConfig{
			//nil,
			//createTankAnimConfig(2, 4, object.DirLeft),
			createTankAnimConfig(2, 4, object.DirRight),
			//createTankAnimConfig(2, 4, object.DirUp),
			//createTankAnimConfig(2, 4, object.DirDown),
			//},
		}, // 同伴坦克四级
		{
			Type: object.ObjTypeMovable, Subtype: object.ObjSubTypeTank, Id: 1000, Level: 1,
			AnimConfig://[]*base.SpriteAnimConfig{
			//nil,
			//createTankAnimConfig(1000, 1, object.DirLeft),
			createTankAnimConfig(1000, 1, object.DirRight),
			//createTankAnimConfig(1000, 1, object.DirUp),
			//createTankAnimConfig(1000, 1, object.DirDown),
			//},
		}, // 轻型坦克
		{
			Type: object.ObjTypeMovable, Subtype: object.ObjSubTypeTank, Id: 1000, Level: 2,
			AnimConfig://[]*base.SpriteAnimConfig{
			//nil,
			//createTankAnimConfig(1000, 2, object.DirLeft),
			createTankAnimConfig(1000, 2, object.DirRight),
			//createTankAnimConfig(1000, 2, object.DirUp),
			//createTankAnimConfig(1000, 2, object.DirDown),
			//},
		}, // 轻型坦克红色
		{
			Type: object.ObjTypeMovable, Subtype: object.ObjSubTypeTank, Id: 1001, Level: 1,
			AnimConfig://[]*base.SpriteAnimConfig{
			//nil,
			//createTankAnimConfig(1001, 1, object.DirLeft),
			createTankAnimConfig(1001, 1, object.DirRight),
			//createTankAnimConfig(1001, 1, object.DirUp),
			//createTankAnimConfig(1001, 1, object.DirDown),
			//},
		}, // 快速坦克
		{
			Type: object.ObjTypeMovable, Subtype: object.ObjSubTypeTank, Id: 1001, Level: 2,
			AnimConfig://[]*base.SpriteAnimConfig{
			//nil,
			//createTankAnimConfig(1001, 2, object.DirLeft),
			createTankAnimConfig(1001, 2, object.DirRight),
			//createTankAnimConfig(1001, 2, object.DirUp),
			//createTankAnimConfig(1001, 2, object.DirDown),
			//},
		}, // 快速坦克红
		{
			Type: object.ObjTypeMovable, Subtype: object.ObjSubTypeTank, Id: 1002, Level: 1,
			AnimConfig://[]*base.SpriteAnimConfig{
			//nil,
			//createTankAnimConfig(1002, 1, object.DirLeft),
			createTankAnimConfig(1002, 1, object.DirRight),
			//createTankAnimConfig(1002, 1, object.DirUp),
			//createTankAnimConfig(1002, 1, object.DirDown),
			//},
		}, // 重型坦克
		{
			Type: object.ObjTypeMovable, Subtype: object.ObjSubTypeTank, Id: 1002, Level: 2,
			AnimConfig://[]*base.SpriteAnimConfig{
			//nil,
			//createTankAnimConfig(1002, 2, object.DirLeft),
			createTankAnimConfig(1002, 2, object.DirRight),
			//createTankAnimConfig(1002, 2, object.DirUp),
			//createTankAnimConfig(1002, 2, object.DirDown),
			//},
		}, // 重型坦克绿
		{
			Type: object.ObjTypeMovable, Subtype: object.ObjSubTypeTank, Id: 1002, Level: 3,
			AnimConfig://[]*base.SpriteAnimConfig{
			//nil,
			//createTankAnimConfig(1002, 3, object.DirLeft),
			createTankAnimConfig(1002, 3, object.DirRight),
			//createTankAnimConfig(1002, 3, object.DirUp),
			//createTankAnimConfig(1002, 3, object.DirDown),
			//},
		}, // 重型坦克黄
		{
			Type: object.ObjTypeMovable, Subtype: object.ObjSubTypeTank, Id: 1002, Level: 4,
			AnimConfig://[]*base.SpriteAnimConfig{
			//nil,
			//createTankAnimConfig(1002, 4, object.DirLeft),
			createTankAnimConfig(1002, 4, object.DirRight),
			//createTankAnimConfig(1002, 4, object.DirUp),
			//createTankAnimConfig(1002, 4, object.DirDown),
			//},
		}, // 重型坦克红
		{
			Type: object.ObjTypeMovable, Subtype: object.ObjSubTypeBullet, Id: 1, Level: 0,
			AnimConfig://[]*base.SpriteAnimConfig{
			//nil,
			//createBulletAnimConfig(object.DirLeft),
			createBulletAnimConfig(object.DirRight),
			//createBulletAnimConfig(object.DirUp),
			//createBulletAnimConfig(object.DirDown),
			//},
		}, // 子弹
		{
			Type: object.ObjTypeMovable, Subtype: object.ObjSubTypeSurroundObj, Id: 1, Level: 0,
			AnimConfig://[]*base.SpriteAnimConfig{
			getSurroundObjAnimConfig(1),
			//},
		}, // 環繞物體
	}
}

func initResources() {
	initImages()
	initAnimOriginalConfigs()
	initAnimConfigs()
}

// 根据静态物体类型创建瓦片动画配置
func createStaticObjAnimConfig(objStaticType object.StaticObjType) *base.SpriteAnimConfig {
	ac := mapTileAnimOriginalConfig[objStaticType]
	fc := mapTileFrameConfig[objStaticType]
	anim := &base.SpriteAnimConfig{
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
func createTankAnimConfig(id int32, level int32, dir object.Direction) *base.SpriteAnimConfig {
	anim := &base.SpriteAnimConfig{}
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
func createBulletAnimConfig(dir object.Direction) *base.SpriteAnimConfig {
	anim := &base.SpriteAnimConfig{
		Image:       bulletAnimOriginalConfig.Image,
		FrameWidth:  bulletAnimOriginalConfig.FrameWidth,
		FrameHeight: bulletAnimOriginalConfig.FrameHeight,
	}
	dirIdx := bulletFrameConfig.dirMap[dir]
	anim.FramePosList = []base.SpriteIndex{bulletAnimOriginalConfig.FramePosList[dirIdx]}
	return anim
}

// 获得可移动物体的动画配置
func GetMovableObjAnimConfig(movableObjType object.MovableObjType, id int32, level int32) *MovableObjectAnimConfig {
	var animConfig *MovableObjectAnimConfig
	for _, c := range moveableObjectAnimConfigList {
		if c.Subtype == object.ObjSubType(movableObjType) && c.Id == id && c.Level == level {
			animConfig = c
			break
		}
	}
	return animConfig
}

// 获得坦克动画配置
func GetTankAnimConfig(id int32, level int32) *MovableObjectAnimConfig {
	return GetMovableObjAnimConfig(object.MovableObjTank, id, level)
}

// 获得子弹动画配置
func GetBulletAnimConfig(id int32) *MovableObjectAnimConfig {
	return GetMovableObjAnimConfig(object.MovableObjBullet, id, 0)
}

// 获得静态物体动画配置
func GetStaticObjAnimConfig(staticObjType object.StaticObjType) *StaticObjectAnimConfig {
	return staticObjectAnimConfigList[staticObjType]
}

// 获得砖块动画配置
func GetBrickAnimConfig() *StaticObjectAnimConfig {
	return GetStaticObjAnimConfig(object.StaticObjBrick)
}

// 获得铁块动画配置
func GetIronAnimConfig() *StaticObjectAnimConfig {
	return GetStaticObjAnimConfig(object.StaticObjIron)
}

// 获得樹木动画配置
func GetTreeAnimConfig() *StaticObjectAnimConfig {
	return GetStaticObjAnimConfig(object.StaticObjTree)
}

// 获得水动画配置
func GetWaterAnimConfig() *StaticObjectAnimConfig {
	return GetStaticObjAnimConfig(object.StaticObjWater)
}

// 获得基地动画配置
func GetHomeAnimConfig() *StaticObjectAnimConfig {
	return GetStaticObjAnimConfig(object.StaticObjHome)
}

// 获得廢墟动画配置
func GetRuinsAnimConfig() *StaticObjectAnimConfig {
	return GetStaticObjAnimConfig(object.StaticObjRuins)
}

// 獲得效果動畫配置
func getEffectAnimConfig(id int32) *base.SpriteAnimConfig {
	return effectAnimOriginalConfig[id]
}

// 獲得小球動畫配置
func getSurroundObjAnimConfig(id int32) *base.SpriteAnimConfig {
	return surroundObjAnimOriginalConfig[id]
}
