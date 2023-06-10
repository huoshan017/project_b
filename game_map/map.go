package game_map

import (
	"project_b/common/object"
)

const (
	minGridSize     = 5
	defaultGridSize = minGridSize
)

// 地图配置，坐标系x轴朝上y轴朝右
type Config struct {
	Id                 int32                // Id
	Layers             [][]int16            // 地图数据
	X, Y               int32                // 左下角坐标
	TileWidth          int32                // 瓦片寬度
	TileHeight         int32                // 瓦片高度
	PlayerTankInitData object.ObjStaticInfo // 玩家坦克配置信息
	PlayerTankInitRect object.Rect          // 玩家坦克出现位置范围矩形
	PlayerMaxCount     int32                // 最大玩家数
}

// 地圖實例
type MapInstance struct {
	config                *Config
	minGridSize           int32
	gridWidth, gridHeight int32
}

func NewMapInstance(gridSize int32) *MapInstance {
	if gridSize < minGridSize {
		gridSize = minGridSize
	}
	return &MapInstance{
		minGridSize: gridSize,
	}
}

func (m *MapInstance) Load(config *Config) {
	tw, th := int32(len(config.Layers[0])), int32(len(config.Layers))
	if tw < m.minGridSize {
		m.minGridSize = tw
	}
	if th < m.minGridSize {
		m.minGridSize = th
	}

	gridWidth, gridHeight := m.minGridSize, m.minGridSize
	for tw%gridWidth > 0 && tw%gridWidth < m.minGridSize {
		gridWidth += 1
	}
	for th%gridHeight > 0 && th%gridHeight < m.minGridSize {
		gridHeight += 1
	}
	m.gridWidth = gridWidth
	m.gridHeight = gridHeight

	m.config = config
}
