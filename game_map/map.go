package game_map

import (
	"project_b/common/ds"
	"project_b/common/math"
	"project_b/common/object"
)

const (
	minGridSize     = 5           // 最小網格大小(tile為單位)
	defaultGridSize = minGridSize // 缺省網格大小
	MapMaxLayer     = 5           // 地圖最大層
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
	objMap                map[uint32]object.IObject
	tiles                 [][]uint32
	tile2pos              map[uint32]struct {
		x int32
		y int32
	}
	grids           [][]*ds.MapListUnion[uint32, object.IObject] // 網格用於碰撞檢測，提高檢測性能
	resultLayerObjs [MapMaxLayer][]uint32
}

func NewMapInstance(gridSize int32) *MapInstance {
	if gridSize < minGridSize {
		gridSize = minGridSize
	}
	return &MapInstance{
		minGridSize: gridSize,
		objMap:      make(map[uint32]object.IObject),
		tile2pos: make(map[uint32]struct {
			x int32
			y int32
		}),
	}
}

func (m *MapInstance) Load(config *Config) {
	tw, th := int32(len(config.Layers[0])), int32(len(config.Layers))

	m.tiles = make([][]uint32, th)
	for i := 0; i < len(config.Layers); i++ {
		m.tiles[i] = make([]uint32, tw)
	}

	if tw < m.minGridSize {
		m.minGridSize = tw
	}
	if th < m.minGridSize {
		m.minGridSize = th
	}
	// 計算網格寬度和高度
	gridWidth, gridHeight := m.minGridSize, m.minGridSize
	for tw%gridWidth > 0 && tw%gridWidth < m.minGridSize {
		gridWidth += 1
	}
	for th%gridHeight > 0 && th%gridHeight < m.minGridSize {
		gridHeight += 1
	}
	m.gridWidth = gridWidth
	m.gridHeight = gridHeight
	// 創建網格
	m.grids = make([][]*ds.MapListUnion[uint32, object.IObject], gridHeight)
	for i := 0; i < len(m.grids); i++ {
		m.grids[i] = make([]*ds.MapListUnion[uint32, object.IObject], gridWidth)
	}
	m.config = config
}

func (m *MapInstance) Unload() {
	m.config = nil
	m.objMap = nil
	m.tiles = nil
	m.tile2pos = nil
	m.grids = nil
	m.minGridSize = 0
	m.gridWidth = 0
	m.gridHeight = 0
}

func (m *MapInstance) AddTile(x, y int32, tile object.IObject) {
	m.tiles[x][y] = tile.InstId()
	m.tile2pos[tile.InstId()] = struct {
		x int32
		y int32
	}{x: x, y: y}
	m.objMap[tile.InstId()] = tile
}

func (m *MapInstance) RemoveTile(instId uint32) {
	pos, o := m.tile2pos[instId]
	if o {
		m.tiles[pos.x][pos.y] = 0
		delete(m.tile2pos, instId)
	}
	delete(m.objMap, instId)
}

func (m *MapInstance) AddObj(obj object.IObject) {
	if _, o := m.objMap[obj.InstId()]; o {
		return
	}
	m.objMap[obj.InstId()] = obj
	lx, by, rx, ty := m.objBounds(obj)
	m.grids[by][lx].Add(obj.InstId(), obj)
	if lx != rx {
		m.grids[by][rx].Add(obj.InstId(), obj)
		if by != ty {
			m.grids[ty][rx].Add(obj.InstId(), obj)
		}
	} else if by != ty {
		m.grids[ty][lx].Add(obj.InstId(), obj)
	}
}

func (m *MapInstance) RemoveObj(instId uint32) {
	var (
		obj object.IObject
		o   bool
	)
	if obj, o = m.objMap[instId]; !o {
		return
	}
	lx, by, rx, ty := m.objBounds(obj)
	m.grids[by][lx].Remove(obj.InstId())
	if lx != rx {
		m.grids[by][rx].Remove(obj.InstId())
		if by != ty {
			m.grids[ty][rx].Remove(obj.InstId())
		}
	} else if by != ty {
		m.grids[ty][lx].Remove(obj.InstId())
	}
	delete(m.objMap, instId)
}

func (m *MapInstance) GetLayerObjsWithRange(rect *math.Rect) [MapMaxLayer][]uint32 {
	// 獲得範圍内的瓦片
	tl := (rect.X() - m.config.X) / m.config.TileWidth
	if tl < 0 {
		tl = 0
	}
	tb := (rect.Y() - m.config.Y) / m.config.TileHeight
	if tb < 0 {
		tb = 0
	}
	tr := (rect.X() + rect.W() - m.config.X) / m.config.TileWidth
	if tr >= int32(len(m.config.Layers[0])) {
		tr = int32(len(m.config.Layers[0])) - 1
	}
	tt := (rect.Y() + rect.H() - m.config.Y) / m.config.TileHeight
	if tt >= int32(len(m.config.Layers)) {
		tt = int32(len(m.config.Layers)) - 1
	}

	for layer := 0; layer < len(m.resultLayerObjs); layer++ {
		m.resultLayerObjs[layer] = m.resultLayerObjs[layer][:0]
	}
	for i := tb; i <= tt; i++ {
		for j := tl; j <= tr; j++ {
			if m.tiles[i][j] <= 0 {
				continue
			}
			obj := m.objMap[m.tiles[i][j]]
			layer := obj.StaticInfo().Layer()
			m.resultLayerObjs[layer] = append(m.resultLayerObjs[layer], m.tiles[i][j])
		}
	}

	// 獲得範圍内的其他obj
	return m.resultLayerObjs
}

func (m *MapInstance) objBounds(obj object.IObject) (lx, by, rx, ty int32) {
	left, bottom := obj.Pos()
	right, top := left+obj.Width(), bottom+obj.Height()
	lx = (left - m.config.X) / m.config.TileWidth
	rx = (right - m.config.X) / m.config.TileWidth
	by = (bottom - m.config.Y) / m.config.TileHeight
	ty = (top - m.config.Y) / m.config.TileHeight
	return
}
