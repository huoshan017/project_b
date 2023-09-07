package main

import (
	"project_b/common/base"
	"project_b/game_map"
)

const (
	defaultNearPlane = 1000 // 默認相機到眎口的距離
)

type mapInfo struct {
	config       *game_map.Config
	cameraPos    base.Pos
	cameraHeight int32
	cameraFov    int32
}

var mapInfoArray = map[int32]mapInfo{
	1: {
		config:       game_map.MapConfigArray[1],
		cameraPos:    base.Pos{X: 2000, Y: 2100},
		cameraHeight: 8000,
		cameraFov:    90,
	},
	2: {
		config:       game_map.MapConfigArray[2],
		cameraPos:    base.Pos{X: 6000, Y: 2000},
		cameraHeight: 9000,
		cameraFov:    90,
	},
	3: {
		config:       game_map.MapConfigArray[3],
		cameraPos:    base.Pos{X: 2000, Y: 2000},
		cameraHeight: 7000,
		cameraFov:    90,
	},
}
