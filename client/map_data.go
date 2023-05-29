package main

import (
	"project_b/common/object"
	"project_b/game_map"
)

const (
	defaultTileSideLength          = 320  // 缺省瓦片逻辑长度, 是图片变长 x 10
	multiplesObjLenAndDisplayLen   = 10   // 与物体尺寸的倍数
	defaultCamera2ViewportDistance = 1000 // 默認相機到眎口的距離
)

type mapInfo struct {
	config       *game_map.Config
	tileSize     int32
	cameraPos    object.Pos
	cameraHeight int32
}

var mapInfoArray = map[int32]mapInfo{
	1: {
		config:       game_map.MapConfigArray[1],
		tileSize:     defaultTileSideLength,
		cameraPos:    object.Pos{X: 2000, Y: 2100},
		cameraHeight: 8000,
	},
	2: {
		config:       game_map.MapConfigArray[2],
		tileSize:     defaultTileSideLength,
		cameraPos:    object.Pos{X: 1000, Y: 1000},
		cameraHeight: 3500,
	},
	3: {
		config:       game_map.MapConfigArray[3],
		tileSize:     defaultTileSideLength,
		cameraPos:    object.Pos{X: 1000, Y: 1000},
		cameraHeight: 3500,
	},
}

// 地图Id列表
var mapIdList = []int32{1, 2, 3}
