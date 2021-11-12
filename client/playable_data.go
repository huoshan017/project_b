package main

import (
	client_base "project_b/client/base"
	"project_b/common/object"
)

var (
	PlayableLerpT = 0.1 // 插值系数
)

type frameConfig struct {
	frameNum       int32
	dirMap         map[object.Direction]int32
	frameLevelList [][]int32
}

var (
	tankAnimConfig  map[int32]*client_base.SpriteAnimConfig
	tankFrameConfig map[int32]*frameConfig
)

func initAnims() {
	tankAnimConfig = map[int32]*client_base.SpriteAnimConfig{
		// 玩家坦克动画配置
		1: {
			Image: player1_img, PlayInterval: 200, FrameWidth: 28, FrameHeight: 28,
			FramePosList: []client_base.SpriteIndex{
				{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 2, Y: 0}, {X: 3, Y: 0}, {X: 4, Y: 0}, {X: 5, Y: 0}, {X: 6, Y: 0}, {X: 7, Y: 0},
				{X: 0, Y: 1}, {X: 1, Y: 1}, {X: 2, Y: 1}, {X: 3, Y: 1}, {X: 4, Y: 1}, {X: 5, Y: 1}, {X: 6, Y: 1}, {X: 7, Y: 1},
				{X: 0, Y: 2}, {X: 1, Y: 2}, {X: 2, Y: 2}, {X: 3, Y: 2}, {X: 4, Y: 2}, {X: 5, Y: 2}, {X: 6, Y: 2}, {X: 7, Y: 2},
				{X: 0, Y: 3}, {X: 1, Y: 3}, {X: 2, Y: 3}, {X: 3, Y: 3}, {X: 4, Y: 3}, {X: 5, Y: 3}, {X: 6, Y: 3}, {X: 7, Y: 3},
			},
		},
		// 同伴坦克动画配置
		2: {
			Image: player2_img, PlayInterval: 200, FrameWidth: 28, FrameHeight: 28,
			FramePosList: []client_base.SpriteIndex{
				{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 2, Y: 0}, {X: 3, Y: 0}, {X: 4, Y: 0}, {X: 5, Y: 0}, {X: 6, Y: 0}, {X: 7, Y: 0},
				{X: 0, Y: 1}, {X: 1, Y: 1}, {X: 2, Y: 1}, {X: 3, Y: 1}, {X: 4, Y: 1}, {X: 5, Y: 1}, {X: 6, Y: 1}, {X: 7, Y: 1},
				{X: 0, Y: 2}, {X: 1, Y: 2}, {X: 2, Y: 2}, {X: 3, Y: 2}, {X: 4, Y: 2}, {X: 5, Y: 2}, {X: 6, Y: 2}, {X: 7, Y: 2},
				{X: 0, Y: 3}, {X: 1, Y: 3}, {X: 2, Y: 3}, {X: 3, Y: 3}, {X: 4, Y: 3}, {X: 5, Y: 3}, {X: 6, Y: 3}, {X: 7, Y: 3},
			},
		},
		// 轻型坦克动画配置
		1000: {
			Image: enemy_img, PlayInterval: 200, FrameWidth: 28, FrameHeight: 28,
			FramePosList: []client_base.SpriteIndex{
				{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 2, Y: 0}, {X: 3, Y: 0},
				{X: 0, Y: 1}, {X: 1, Y: 1}, {X: 2, Y: 1}, {X: 3, Y: 1},
				{X: 0, Y: 2}, {X: 1, Y: 2}, {X: 2, Y: 2}, {X: 3, Y: 2},
				{X: 0, Y: 3}, {X: 1, Y: 3}, {X: 2, Y: 3}, {X: 3, Y: 3},
			},
		},
		// 快速坦克动画配置
		1001: {
			Image: enemy_img, PlayInterval: 200, FrameWidth: 28, FrameHeight: 28,
			FramePosList: []client_base.SpriteIndex{
				{X: 4, Y: 0}, {X: 5, Y: 0}, {X: 6, Y: 0}, {X: 7, Y: 0},
				{X: 4, Y: 1}, {X: 5, Y: 1}, {X: 6, Y: 1}, {X: 7, Y: 1},
				{X: 4, Y: 2}, {X: 5, Y: 2}, {X: 6, Y: 2}, {X: 7, Y: 2},
				{X: 4, Y: 3}, {X: 5, Y: 3}, {X: 6, Y: 3}, {X: 7, Y: 3},
			},
		},
		// 重型坦克动画配置
		1002: {
			Image: enemy_img, PlayInterval: 200, FrameWidth: 28, FrameHeight: 28,
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

func init() {
	initImages()
	initAnims()
}

// 根据等级和方向创建坦克动画
func CreateTankAnimConfig(id int32, level int32, dir object.Direction) *client_base.SpriteAnimConfig {
	anim := &client_base.SpriteAnimConfig{}
	anim0 := tankAnimConfig[id]
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
