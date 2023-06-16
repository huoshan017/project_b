package game_map

import (
	"project_b/common/ds"
	"project_b/common/log"
	"project_b/common/math"
	"project_b/common/object"
)

const (
	minGridTileSize     int16 = 5               // 最小網格大小(tile為單位)
	defaultGridTileSize int16 = minGridTileSize // 缺省網格大小
	MapMaxLayer         int16 = 5               // 地圖最大層
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
	config                  *Config
	minGridTileSize         int16                     // 最小网格中的瓦片数
	gridXTiles, gridYTiles  int16                     // x轴网格瓦片数，y轴网格瓦片数
	gridLineNum, gridColNum int16                     // 網格行數和列數
	gridWidth, gridHeight   int32                     // 網格寬度高度
	objMap                  map[uint32]object.IObject // 对象map
	tiles                   [][]uint32                // 瓦片对象ID二维数组
	tile2pos                map[uint32]struct {
		line int16
		col  int16
	} // 瓦片ID对应的坐标位置
	grids           []*ds.MapListUnion[uint32, struct{}] // 網格用於碰撞檢測，提高檢測性能
	resultLayerObjs [MapMaxLayer][]uint32                // 缓存返回的结果给调用者，主要为了提高性能，减少GC
	tempKeys        []uint32                             // 临时缓存对象id，提高性能，减少GC
}

func NewMapInstance(gridTileSize int16) *MapInstance {
	if gridTileSize < minGridTileSize {
		gridTileSize = minGridTileSize
	}
	return &MapInstance{
		minGridTileSize: gridTileSize,
		objMap:          make(map[uint32]object.IObject),
		tile2pos: make(map[uint32]struct {
			line int16
			col  int16
		}),
	}
}

func (m *MapInstance) Load(config *Config) {
	m.config = config

	colNum, lineNum := int16(len(config.Layers[0])), int16(len(config.Layers))
	m.tiles = make([][]uint32, lineNum)
	for i := 0; i < len(config.Layers); i++ {
		m.tiles[i] = make([]uint32, colNum)
	}
	if colNum < m.minGridTileSize {
		m.minGridTileSize = colNum
	}
	if lineNum < m.minGridTileSize {
		m.minGridTileSize = lineNum
	}

	// 計算網格寬度和高度
	gridTileWidth, gridTileHeight := m.minGridTileSize, m.minGridTileSize
	for colNum%gridTileWidth > 0 && colNum%gridTileWidth < m.minGridTileSize {
		gridTileWidth += 1
	}
	for lineNum%gridTileHeight > 0 && lineNum%gridTileHeight < m.minGridTileSize {
		gridTileHeight += 1
	}

	// 網格x軸和y軸方向包含的tile數量
	m.gridXTiles = gridTileWidth
	m.gridYTiles = gridTileHeight

	// 創建網格
	gridColNum, gridLineNum := (colNum+gridTileWidth-1)/gridTileWidth, (lineNum+gridTileHeight-1)/gridTileHeight
	m.grids = make([]*ds.MapListUnion[uint32, struct{}], gridLineNum*gridColNum)
	for i := 0; i < int(gridLineNum*gridColNum); i++ {
		m.grids[i] = ds.NewMapListUnion[uint32, struct{}]()
	}
	m.gridLineNum = gridLineNum
	m.gridColNum = gridColNum

	// 網格寬度和高度
	m.gridWidth = int32(gridTileWidth) * m.config.TileWidth
	m.gridHeight = int32(gridTileHeight) * m.config.TileHeight
}

func (m *MapInstance) Unload() {
	m.config = nil
	m.objMap = nil
	m.tiles = nil
	m.tile2pos = nil
	m.grids = nil
	m.minGridTileSize = 0
	m.gridXTiles = 0
	m.gridYTiles = 0
	m.gridLineNum = 0
	m.gridColNum = 0
	m.gridWidth = 0
	m.gridHeight = 0
}

func (m *MapInstance) AddTile(line, col int16, tile object.IObject) {
	m.tiles[line][col] = tile.InstId()
	m.tile2pos[tile.InstId()] = struct {
		line int16
		col  int16
	}{line: line, col: col}
	m.objMap[tile.InstId()] = tile
}

func (m *MapInstance) RemoveTile(instId uint32) {
	pos, o := m.tile2pos[instId]
	if o {
		m.tiles[pos.line][pos.col] = 0
		delete(m.tile2pos, instId)
	}
	delete(m.objMap, instId)
}

func (m *MapInstance) AddObj(obj object.IObject) {
	instId := obj.InstId()
	if _, o := m.objMap[instId]; o {
		return
	}
	m.objMap[instId] = obj
	lx, by, rx, ty := m.objGridBounds(obj)
	if lx <= rx && by <= ty {
		index := m.gridLineCol2Index(by, lx)
		m.grids[index].Add(instId, struct{}{})
		log.Info("MapInstance: obj %v add to grid(line: %v, col: %v, index: %v)", instId, by, lx, index)
		if lx != rx && by == ty {
			index = m.gridLineCol2Index(by, rx)
			m.grids[index].Add(instId, struct{}{})
			log.Info("MapInstance: obj %v add to grid(line: %v, col: %v, index: %v)", instId, by, rx, index)
		}
		if lx != rx && by != ty {
			index = m.gridLineCol2Index(ty, rx)
			m.grids[m.gridLineCol2Index(ty, rx)].Add(instId, struct{}{})
			log.Info("MapInstance: obj %v add to grid(line: %v, col: %v, index: %v)", instId, ty, rx, index)
		}
		if lx == rx && by != ty {
			index = m.gridLineCol2Index(ty, lx)
			m.grids[index].Add(instId, struct{}{})
			log.Info("MapInstance: obj %v add to grid(line: %v, col: %v, index: %v)", instId, ty, lx, index)
		}
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
	lx, by, rx, ty := m.objGridBounds(obj)
	if lx <= rx && by <= ty {
		index := m.gridLineCol2Index(by, lx)
		m.grids[index].Remove(instId)
		log.Info("MapInstance: obj %v remove from grid(line: %v, col: %v, index: %v)", instId, by, lx, index)
		if lx != rx && by == ty {
			index = m.gridLineCol2Index(by, rx)
			m.grids[index].Remove(instId)
			log.Info("MapInstance: obj %v remove from grid(line: %v, col: %v, index: %v)", instId, by, rx, index)
		}
		if lx != rx && by != ty {
			index = m.gridLineCol2Index(ty, rx)
			m.grids[index].Remove(instId)
			log.Info("MapInstance: obj %v remove from grid(line: %v, col: %v, index: %v)", instId, ty, rx, index)
		}
		if lx == rx && by != ty {
			index = m.gridLineCol2Index(ty, lx)
			m.grids[index].Remove(instId)
			log.Info("MapInstance: obj %v remove from grid(line: %v, col: %v, index: %v)", instId, ty, lx, index)
		}
	}
	delete(m.objMap, instId)
}

func (m *MapInstance) UpdateObj(obj object.IObject) {
	x, y := obj.Pos()
	lastX, lastY := obj.LastPos()
	if x == lastX && y == lastY {
		return
	}

	//log.Info("lastX: %v, lastY：%v | x: %v, y: %v", lastX, lastY, x, y)

	left1, bottom1, right1, top1 := m.gridBoundsBy(x, y, x+obj.Width(), y+obj.Height())
	left0, bottom0, right0, top0 := m.gridBoundsBy(lastX, lastY, lastX+obj.Width(), lastY+obj.Height())

	//log.Info("left0: %v, bottom0: %v, right0: %v, top0: %v", left0, bottom0, right0, top0)
	//log.Info("left1: %v, bottom1: %v, right1: %v, top1: %v", left1, bottom1, right1, top1)

	// ---------------|---------------|---------------|--
	//                |               |               |
	//               _|__             |               |
	//              | |  |            |               |
	//              !_|__!            |               |
	//                |               |               |
	// ---------------|---------------|---------------|--
	//           ____ |               |               |
	//          |    ||               |               |
	//          !____!|               |               |
	//                |               |               |
	//                |               |               |
	// ---------------|---------------|---------------|--

	var index int32

	if x > lastX { // 向右移動
		if left1 > left0 { // 離開 column left0
			index = m.gridLineCol2Index(bottom0, left0)
			m.grids[index].Remove(obj.InstId())
			log.Info("MapInstance.UpdateObj: obj %v remove from grid(line: %v, col: %v, index: %v)", obj.InstId(), bottom0, left0, index)
			if bottom0 != top0 { // 説明左下和左上在不同的grid中，離開 (top0, left0)
				index = m.gridLineCol2Index(top0, left0)
				m.grids[index].Remove(obj.InstId())
				log.Info("MapInstance.UpdateObj: obj %v remove from grid(line: %v, col: %v, index: %v)", obj.InstId(), top0, left0, index)
			}
		}
		if right1 > right0 { // 進入 column right1
			index = m.gridLineCol2Index(bottom1, right1)
			m.grids[index].Add(obj.InstId(), struct{}{})
			log.Info("MapInstance.UpdateObj: obj %v add to grid(line: %v, col: %v, index: %v)", obj.InstId(), bottom1, right1, index)
			if bottom1 != top1 { // 右下和右上在不同的grid中，右上進入(top1, right1)
				index = m.gridLineCol2Index(top1, right1)
				m.grids[index].Add(obj.InstId(), struct{}{})
				log.Info("MapInstance.UpdateObj: obj %v add to grid(line: %v, col: %v, index: %v)", obj.InstId(), top1, right1, index)
			}
		}
	} else if x < lastX { // 向左移動
		if left1 < left0 { //
			index = m.gridLineCol2Index(bottom1, left1)
			m.grids[index].Add(obj.InstId(), struct{}{})
			log.Info("MapInstance.UpdateObj: obj %v add to grid(line: %v, col: %v, index: %v)", obj.InstId(), bottom1, left1, index)
			if bottom1 != top1 { // 左下和左上在不同的grid中，左上進入(top1, left1)
				index = m.gridLineCol2Index(top1, left1)
				m.grids[index].Add(obj.InstId(), struct{}{})
				log.Info("MapInstance.UpdateObj: obj %v add to grid(line: %v, col: %v, index: %v)", obj.InstId(), top1, left1, index)
			}
		}
		if right1 < right0 { // 右側離開原來的grid
			index = m.gridLineCol2Index(bottom0, right0)
			m.grids[index].Remove(obj.InstId())
			log.Info("MapInstance.UpdateObj: obj %v remove from grid(line: %v, col: %v, index: %v)", obj.InstId(), bottom0, right0, index)
			if top0 != bottom0 { // 右上和右下在不同grid中，離開(top0, right0)
				index = m.gridLineCol2Index(top0, right0)
				m.grids[index].Remove(obj.InstId())
				log.Info("MapInstance.UpdateObj: obj %v remove from grid(line: %v, col: %v, index: %v)", obj.InstId(), top0, right0, index)
			}
		}
	}

	if y > lastY { // 向上移動
		if bottom1 > bottom0 { // 向上移動，左下離開(bottom0, left0)
			index = m.gridLineCol2Index(bottom0, left0)
			m.grids[index].Remove(obj.InstId())
			log.Info("MapInstance.UpdateObj: obj %v remove from grid(line: %v, col: %v, index: %v)", obj.InstId(), bottom0, left0, index)
			if left0 != right0 { // 左下和右下在不同的grid中，離開(bottom0, right0)
				index = m.gridLineCol2Index(bottom0, right0)
				m.grids[index].Remove(obj.InstId())
				log.Info("MapInstance.UpdateObj: obj %v remove from grid(line: %v, col: %v, index: %v)", obj.InstId(), bottom0, right0, index)
			}
		}
		if top1 > top0 { // 頂部進入新grid
			index = m.gridLineCol2Index(top1, left1)
			m.grids[index].Add(obj.InstId(), struct{}{})
			log.Info("MapInstance.UpdateObj: obj %v add to grid(line: %v, col: %v, index: %v)", obj.InstId(), top1, left1, index)
			if left1 != right1 { // 左上和右上在不同的grid中，右上進入(top1, right1)
				index = m.gridLineCol2Index(top1, right1)
				m.grids[index].Add(obj.InstId(), struct{}{})
				log.Info("MapInstance.UpdateObj: obj %v add to grid(line: %v, col: %v, index: %v)", obj.InstId(), top1, right1, index)
			}
		}
	} else if y < lastY { // 向下移動
		if bottom1 < bottom0 { // 向上移動，左下進入(bottom1, left1)
			index = m.gridLineCol2Index(bottom1, left1)
			m.grids[index].Add(obj.InstId(), struct{}{})
			log.Info("MapInstance.UpdateObj: obj %v add to grid(line: %v, col: %v, index: %v)", obj.InstId(), bottom1, left1, index)
			if left1 != right1 { // 左下和右上在不同的grid中，右下也要進入(bottom1, right1)
				index = m.gridLineCol2Index(bottom1, right1)
				m.grids[index].Add(obj.InstId(), struct{}{})
				log.Info("MapInstance.UpdateObj: obj %v add to grid(line: %v, col: %v, index: %v)", obj.InstId(), bottom1, right1, index)
			}
		}
		if top1 < top0 { // 頂部左上離開(top0, left0)
			index = m.gridLineCol2Index(top0, left0)
			m.grids[index].Remove(obj.InstId())
			log.Info("MapInstance.UpdateObj: obj %v remove from grid(line: %v, col: %v, index: %v)", obj.InstId(), top0, left0, index)
			if left0 != right0 { // 左上和右上在不同的grid中，右上離開(top0, right0)
				index = m.gridLineCol2Index(top0, right0)
				m.grids[index].Remove(obj.InstId())
				log.Info("MapInstance.UpdateObj: obj %v remove from grid(line: %v, col: %v, index: %v)", obj.InstId(), top0, right0, index)
			}
		}
	}
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

	// 先清空
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
	// 先获取grids
	lx, by, rx, ty := m.gridBoundsBy(rect.X(), rect.Y(), rect.X()+rect.W(), rect.Y()+rect.H())
	if rx >= lx && ty >= by {
		m.tempKeys = m.tempKeys[:0]
		for y := by; y <= ty; y++ {
			for x := lx; x <= rx; x++ {
				m.tempKeys = m.grids[m.gridLineCol2Index(y, x)].GetKeys(m.tempKeys)
			}
		}

		for _, k := range m.tempKeys {
			obj := m.objMap[k]
			if obj != nil {
				layer := obj.StaticInfo().Layer()
				m.resultLayerObjs[layer] = append(m.resultLayerObjs[layer], k)
			}
		}
	}

	return m.resultLayerObjs
}

func (m *MapInstance) lineCol2Index(line, col int16) int32 {
	return int32(line)*int32(len(m.config.Layers)) + int32(col)
}

func (m *MapInstance) index2LineCol(index int32) (int16, int16) {
	return int16(index / int32(len(m.config.Layers))), int16(index % int32(m.gridXTiles))
}

func (m *MapInstance) gridLineCol2Index(line, col int16) int32 {
	return int32(line)*int32(m.gridColNum) + int32(col)
}

func (m *MapInstance) gridIndex2LineCol(index int32) (int16, int16) {
	return int16(index / int32(m.gridColNum)), int16(index % int32(m.gridColNum))
}

func (m *MapInstance) objGridBounds(obj object.IObject) (lx, by, rx, ty int16) {
	left, bottom := obj.Pos()
	right, top := left+obj.Width(), bottom+obj.Height()
	return m.gridBoundsBy(left, bottom, right, top)
}

func (m *MapInstance) gridBoundsBy(left, bottom, right, top int32) (lx, by, rx, ty int16) {
	lx = int16((left - m.config.X) / m.gridWidth)
	rx = int16((right - m.config.X) / m.gridWidth)
	by = int16((bottom - m.config.Y) / m.gridHeight)
	ty = int16((top - m.config.Y) / m.gridHeight)
	if lx < 0 {
		lx = 0
	}
	if rx >= m.gridColNum-1 {
		rx = m.gridColNum - 1
	}
	if by < 0 {
		by = 0
	}
	if ty >= m.gridLineNum-1 {
		ty = m.gridLineNum - 1
	}
	return
}
