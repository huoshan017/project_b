package main

import (
	"project_b/game_map"

	"github.com/hajimehoshi/ebiten/v2/examples/resources/images"
)

const (
	defaultTileSize = 16
)

type mapInfo struct {
	config         game_map.Config
	tileSize       int
	tileXNum       int
	tilesImageData []byte
}

var mapInfoArray = []mapInfo{
	{
		config:         game_map.MapConfigArray[0],
		tileSize:       defaultTileSize,
		tileXNum:       25,
		tilesImageData: images.Tiles_png,
	},
}
