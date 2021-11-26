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
	tankAnimOriginalConfig map[int32]*base.SpriteAnimConfig // 坦克动画初始配置map
	tankFrameConfig        map[int32]*frameConfig           // 帧配置map
	bulletAnimConfig       map[int32]*base.SpriteAnimConfig // 子弹动画配置map
)

// 初始化动画初始配置
func initAnimOriginalConfigs() {
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
}

// 物体动画配置
type ObjectAnimConfig struct {
	Index      int32
	Type       object.ObjectType
	Subtype    object.ObjSubType
	Id         int32
	Level      int32
	AnimConfig []*base.SpriteAnimConfig // 使用物体方向枚举做索引
}

var (
	moveableObjectAnimConfigList []*ObjectAnimConfig // 移动物体动画配置
)

// 初始化动画配置
func initAnimConfigs() {
	moveableObjectAnimConfigList = []*ObjectAnimConfig{
		&ObjectAnimConfig{Index: 1, Type: object.ObjTypeMovable, Subtype: object.ObjSubTypeTank, Id: 1, Level: 1,
			AnimConfig: []*base.SpriteAnimConfig{
				// 空左右上下
				nil,
				CreateTankAnimConfig(1, 1, object.DirLeft),
				CreateTankAnimConfig(1, 1, object.DirRight),
				CreateTankAnimConfig(1, 1, object.DirUp),
				CreateTankAnimConfig(1, 1, object.DirDown),
			}}, // 玩家坦克一级
		&ObjectAnimConfig{Index: 2, Type: object.ObjTypeMovable, Subtype: object.ObjSubTypeTank, Id: 1, Level: 2,
			AnimConfig: []*base.SpriteAnimConfig{
				nil,
				CreateTankAnimConfig(1, 2, object.DirLeft),
				CreateTankAnimConfig(1, 2, object.DirRight),
				CreateTankAnimConfig(1, 2, object.DirUp),
				CreateTankAnimConfig(1, 2, object.DirDown),
			}}, // 玩家坦克二级
		&ObjectAnimConfig{Index: 3, Type: object.ObjTypeMovable, Subtype: object.ObjSubTypeTank, Id: 1, Level: 3,
			AnimConfig: []*base.SpriteAnimConfig{
				nil,
				CreateTankAnimConfig(1, 3, object.DirLeft),
				CreateTankAnimConfig(1, 3, object.DirRight),
				CreateTankAnimConfig(1, 3, object.DirUp),
				CreateTankAnimConfig(1, 3, object.DirDown),
			}}, // 玩家坦克三级
		&ObjectAnimConfig{Index: 4, Type: object.ObjTypeMovable, Subtype: object.ObjSubTypeTank, Id: 1, Level: 4,
			AnimConfig: []*base.SpriteAnimConfig{
				nil,
				CreateTankAnimConfig(1, 4, object.DirLeft),
				CreateTankAnimConfig(1, 4, object.DirRight),
				CreateTankAnimConfig(1, 4, object.DirUp),
				CreateTankAnimConfig(1, 4, object.DirDown),
			}}, // 玩家坦克四级
		&ObjectAnimConfig{Index: 5, Type: object.ObjTypeMovable, Subtype: object.ObjSubTypeTank, Id: 2, Level: 1,
			AnimConfig: []*base.SpriteAnimConfig{
				nil,
				CreateTankAnimConfig(2, 1, object.DirLeft),
				CreateTankAnimConfig(2, 1, object.DirRight),
				CreateTankAnimConfig(2, 1, object.DirUp),
				CreateTankAnimConfig(2, 1, object.DirDown),
			}}, // 同伴坦克一级
		&ObjectAnimConfig{Index: 6, Type: object.ObjTypeMovable, Subtype: object.ObjSubTypeTank, Id: 2, Level: 2,
			AnimConfig: []*base.SpriteAnimConfig{
				nil,
				CreateTankAnimConfig(2, 2, object.DirLeft),
				CreateTankAnimConfig(2, 2, object.DirRight),
				CreateTankAnimConfig(2, 2, object.DirUp),
				CreateTankAnimConfig(2, 2, object.DirDown),
			}}, // 同伴坦克二级
		&ObjectAnimConfig{Index: 7, Type: object.ObjTypeMovable, Subtype: object.ObjSubTypeTank, Id: 2, Level: 3,
			AnimConfig: []*base.SpriteAnimConfig{
				nil,
				CreateTankAnimConfig(2, 3, object.DirLeft),
				CreateTankAnimConfig(2, 3, object.DirRight),
				CreateTankAnimConfig(2, 3, object.DirUp),
				CreateTankAnimConfig(2, 3, object.DirDown),
			}}, // 同伴坦克三级
		&ObjectAnimConfig{Index: 8, Type: object.ObjTypeMovable, Subtype: object.ObjSubTypeTank, Id: 2, Level: 4,
			AnimConfig: []*base.SpriteAnimConfig{
				nil,
				CreateTankAnimConfig(2, 4, object.DirLeft),
				CreateTankAnimConfig(2, 4, object.DirRight),
				CreateTankAnimConfig(2, 4, object.DirUp),
				CreateTankAnimConfig(2, 4, object.DirDown),
			}}, // 同伴坦克四级
		&ObjectAnimConfig{Index: 9, Type: object.ObjTypeMovable, Subtype: object.ObjSubTypeTank, Id: 1000, Level: 1,
			AnimConfig: []*base.SpriteAnimConfig{
				nil,
				CreateTankAnimConfig(1000, 1, object.DirLeft),
				CreateTankAnimConfig(1000, 1, object.DirRight),
				CreateTankAnimConfig(1000, 1, object.DirUp),
				CreateTankAnimConfig(1000, 1, object.DirDown),
			}}, // 轻型坦克
		&ObjectAnimConfig{Index: 10, Type: object.ObjTypeMovable, Subtype: object.ObjSubTypeTank, Id: 1000, Level: 2,
			AnimConfig: []*base.SpriteAnimConfig{
				nil,
				CreateTankAnimConfig(1000, 2, object.DirLeft),
				CreateTankAnimConfig(1000, 2, object.DirRight),
				CreateTankAnimConfig(1000, 2, object.DirUp),
				CreateTankAnimConfig(1000, 2, object.DirDown),
			}}, // 轻型坦克红色
		&ObjectAnimConfig{Index: 11, Type: object.ObjTypeMovable, Subtype: object.ObjSubTypeTank, Id: 1001, Level: 1,
			AnimConfig: []*base.SpriteAnimConfig{
				nil,
				CreateTankAnimConfig(1001, 1, object.DirLeft),
				CreateTankAnimConfig(1001, 1, object.DirRight),
				CreateTankAnimConfig(1001, 1, object.DirUp),
				CreateTankAnimConfig(1001, 1, object.DirDown),
			}}, // 快速坦克
		&ObjectAnimConfig{Index: 12, Type: object.ObjTypeMovable, Subtype: object.ObjSubTypeTank, Id: 1001, Level: 2,
			AnimConfig: []*base.SpriteAnimConfig{
				nil,
				CreateTankAnimConfig(1001, 2, object.DirLeft),
				CreateTankAnimConfig(1001, 2, object.DirRight),
				CreateTankAnimConfig(1001, 2, object.DirUp),
				CreateTankAnimConfig(1001, 2, object.DirDown),
			}}, // 快速坦克红
		&ObjectAnimConfig{Index: 13, Type: object.ObjTypeMovable, Subtype: object.ObjSubTypeTank, Id: 1002, Level: 1,
			AnimConfig: []*base.SpriteAnimConfig{
				nil,
				CreateTankAnimConfig(1002, 1, object.DirLeft),
				CreateTankAnimConfig(1002, 1, object.DirRight),
				CreateTankAnimConfig(1002, 1, object.DirUp),
				CreateTankAnimConfig(1002, 1, object.DirDown),
			}}, // 重型坦克
		&ObjectAnimConfig{Index: 14, Type: object.ObjTypeMovable, Subtype: object.ObjSubTypeTank, Id: 1002, Level: 2,
			AnimConfig: []*base.SpriteAnimConfig{
				nil,
				CreateTankAnimConfig(1002, 2, object.DirLeft),
				CreateTankAnimConfig(1002, 2, object.DirRight),
				CreateTankAnimConfig(1002, 2, object.DirUp),
				CreateTankAnimConfig(1002, 2, object.DirDown),
			}}, // 重型坦克绿
		&ObjectAnimConfig{Index: 15, Type: object.ObjTypeMovable, Subtype: object.ObjSubTypeTank, Id: 1002, Level: 3,
			AnimConfig: []*base.SpriteAnimConfig{
				nil,
				CreateTankAnimConfig(1002, 3, object.DirLeft),
				CreateTankAnimConfig(1002, 3, object.DirRight),
				CreateTankAnimConfig(1002, 3, object.DirUp),
				CreateTankAnimConfig(1002, 3, object.DirDown),
			}}, // 重型坦克黄
		&ObjectAnimConfig{Index: 16, Type: object.ObjTypeMovable, Subtype: object.ObjSubTypeTank, Id: 1002, Level: 4,
			AnimConfig: []*base.SpriteAnimConfig{
				nil,
				CreateTankAnimConfig(1002, 4, object.DirLeft),
				CreateTankAnimConfig(1002, 4, object.DirRight),
				CreateTankAnimConfig(1002, 4, object.DirUp),
				CreateTankAnimConfig(1002, 4, object.DirDown),
			}}, // 重型坦克红
	}
}

func init() {
	initImages()
	initAnimOriginalConfigs()
	initAnimConfigs()
}

// 根据等级和方向创建坦克动画配置
func CreateTankAnimConfig(id int32, level int32, dir object.Direction) *base.SpriteAnimConfig {
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

// 获得物体动画配置
func GetObjectAnimConfig(typ object.ObjectType, subtype object.ObjSubType, id int32, level int32) *ObjectAnimConfig {
	var objAnimConfig *ObjectAnimConfig
	switch typ {
	case object.ObjTypeMovable:
		for _, c := range moveableObjectAnimConfigList {
			if c.Subtype == subtype && c.Id == id && c.Level == level {
				objAnimConfig = c
				break
			}
		}
	case object.ObjTypeStatic:
	}
	return objAnimConfig
}

// 获得坦克动画配置
func GetTankAnimConfig(id int32, level int32) *ObjectAnimConfig {
	return GetObjectAnimConfig(object.ObjTypeMovable, object.ObjSubTypeTank, id, level)
}