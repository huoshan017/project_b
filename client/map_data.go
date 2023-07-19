package main

import (
	"project_b/common/object"
	"project_b/game_map"
)

const (
	defaultNearPlane = 1000 // 默認相機到眎口的距離
)

type mapInfo struct {
	config       *game_map.Config
	cameraPos    object.Pos
	cameraHeight int32
	cameraFov    int32
}

var mapInfoArray = map[int32]mapInfo{
	1: {
		config:       game_map.MapConfigArray[1],
		cameraPos:    object.Pos{X: 2000, Y: 2100},
		cameraHeight: 8000,
		cameraFov:    90,
	},
	2: {
		config:       game_map.MapConfigArray[2],
		cameraPos:    object.Pos{X: 6000, Y: 2000},
		cameraHeight: 9000,
		cameraFov:    90,
	},
	3: {
		config:       game_map.MapConfigArray[3],
		cameraPos:    object.Pos{X: 1000, Y: 1000},
		cameraHeight: 3500,
		cameraFov:    60,
	},
}
