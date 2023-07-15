package common

import (
	"project_b/common/ds"
	"project_b/common/log"
	"project_b/common/math"
	"project_b/common/object"
	"project_b/game_map"

	"github.com/huoshan017/ponu/heap"
)

type gridObjList []uint32

func (gl *gridObjList) addObj(instId uint32) bool {
	l := len(*gl)
	for i := 0; i < l; i++ {
		if (*gl)[i] == instId {
			return false
		}
	}
	*gl = append(*gl, instId)
	return true
}

func (gl *gridObjList) removeObj(instId uint32) bool {
	l := len(*gl)
	for i := 0; i < l; i++ {
		if (*gl)[i] == instId {
			*gl = append((*gl)[:i], (*gl)[i+1:]...)
			return true
		}
	}
	return false
}

func (gl *gridObjList) clear() {
	if len(*gl) > 0 {
		*gl = (*gl)[:0]
	}
}

// 網格地圖
type GridMap struct {
	config                  *game_map.Config
	minGridTileSize         int16                                           // 最小网格中的瓦片数
	gridXTiles, gridYTiles  int16                                           // x轴网格瓦片数，y轴网格瓦片数
	gridLineNum, gridColNum int16                                           // 網格行數和列數
	gridWidth, gridHeight   int32                                           // 網格寬度高度
	sobjs                   *ds.MapListUnion[uint32, object.IStaticObject]  // 对象map
	mobjs                   *ds.MapListUnion[uint32, object.IMovableObject] // 移動對象map
	grids                   []gridObjList                                   // 網格用於碰撞檢測，提高檢測性能
	resultLayerObjs         [MapMaxLayer]*heap.BinaryHeapKV[uint32, int32]  // 缓存返回的结果给调用者，主要为了提高性能，减少GC
	resultMovableObjList    []uint32                                        // 緩存搜索的可移動物體的結果
}

func NewGridMap(gridTileSize int16) *GridMap {
	return &GridMap{
		minGridTileSize: gridTileSize,
		sobjs:           ds.NewMapListUnion[uint32, object.IStaticObject](),
		mobjs:           ds.NewMapListUnion[uint32, object.IMovableObject](),
	}
}

func (m *GridMap) Load(config *game_map.Config) {
	m.config = config

	colNum, lineNum := int16(len(config.Layers[0])), int16(len(config.Layers))
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
	m.grids = make([]gridObjList, gridLineNum*gridColNum)
	for i := 0; i < int(gridLineNum*gridColNum); i++ {
		m.grids[i] = gridObjList{}
	}
	m.gridLineNum = gridLineNum
	m.gridColNum = gridColNum

	// 網格寬度和高度
	m.gridWidth = int32(gridTileWidth) * m.config.TileWidth
	m.gridHeight = int32(gridTileHeight) * m.config.TileHeight

	// 創建結果層
	for i := 0; i < len(m.resultLayerObjs); i++ {
		m.resultLayerObjs[i] = heap.NewMaxBinaryHeapKV[uint32, int32]()
	}
}

func (m *GridMap) Unload() {
	m.config = nil
	for i := 0; i < len(m.grids); i++ {
		m.grids[i].clear()
	}
	m.grids = m.grids[:0]
	m.sobjs.Clear()
	m.mobjs.Clear()
	m.grids = nil
	m.minGridTileSize = 0
	m.gridXTiles = 0
	m.gridYTiles = 0
	m.gridLineNum = 0
	m.gridColNum = 0
	m.gridWidth = 0
	m.gridHeight = 0
}

func (m *GridMap) AddTile(line, col int16, tile object.IObject) {
	m.AddObj(tile)
}

func (m *GridMap) RemoveTile(instId uint32) {
	m.RemoveObj(instId)
}

func (m *GridMap) AddObj(obj object.IObject) {
	instId := obj.InstId()
	if obj.Type() == object.ObjTypeMovable {
		if m.mobjs.Exists(instId) {
			return
		}
	} else {
		if m.sobjs.Exists(instId) {
			return
		}
	}

	x, y := obj.Pos()
	index := m.posGridIndex(x, y)
	m.grids[index].addObj(instId)
	log.Info("MapInstance: obj %v add to grid(index: %v)", instId, index)
	if obj.Type() == object.ObjTypeMovable {
		m.mobjs.Add(instId, obj.(object.IMovableObject))
	} else {
		m.sobjs.Add(instId, obj.(object.IStaticObject))
	}
}

func (m *GridMap) RemoveObj(instId uint32) {
	var (
		obj object.IObject
		o   bool
	)
	obj, o = m.mobjs.Get(instId)
	if o {
		m.mobjs.Remove(instId)
	} else {
		obj, o = m.sobjs.Get(instId)
		if !o {
			return
		}
		m.sobjs.Remove(instId)
	}

	x, y := obj.Pos()
	index := m.posGridIndex(x, y)
	m.grids[index].removeObj(instId)
	log.Info("MapInstance: obj %v remove from grid(index: %v)", instId, index)
}

func (m *GridMap) UpdateMovable(obj object.IMovableObject) {
	if !m.mobjs.Exists(obj.InstId()) {
		log.Warn("obj %v (type: %v, subtype: %v) not found in map", obj.InstId(), obj.Type(), obj.Subtype())
		return
	}

	x, y := obj.Pos()
	lastX, lastY := obj.LastPos()
	if x == lastX && y == lastY {
		return
	}

	index := m.posGridIndex(x, y)
	lastIndex := m.posGridIndex(lastX, lastY)

	if index != lastIndex {
		m.grids[lastIndex].removeObj(obj.InstId())
		m.grids[index].addObj(obj.InstId())
	}
}

func (m *GridMap) GetLayerObjsWithRange(rect *math.Rect) [MapMaxLayer]*heap.BinaryHeapKV[uint32, int32] {
	// 獲得範圍内的靜止和運動的obj
	// rect的範圍向外擴展半個Tile(默認一個grid就是一個Tile)
	left := rect.X() - m.config.TileWidth/2
	if left < m.config.X {
		left = m.config.X
	}
	bottom := rect.Y() - m.config.TileHeight/2
	if bottom < m.config.Y {
		bottom = m.config.Y
	}
	right := rect.X() + rect.W() + m.config.TileWidth/2
	if right > m.config.X+int32(len(m.config.Layers[0]))*m.config.TileWidth {
		right = m.config.X + int32(len(m.config.Layers[0]))*m.config.TileWidth
	}
	top := rect.Y() + rect.H() + m.config.TileHeight/2
	if top > m.config.Y+int32(len(m.config.Layers))*m.config.TileHeight {
		top = m.config.Y + int32(len(m.config.Layers))*m.config.TileHeight
	}

	lx, by, rx, ty := m.gridBoundsBy(left, bottom, right, top)
	if rx >= lx && ty >= by {
		for y := by; y <= ty; y++ {
			for x := lx; x <= rx; x++ {
				var (
					obj object.IObject
					o   bool
				)
				gidx := m.gridLineCol2Index(y, x)
				lis := m.grids[gidx]
				for i := 0; i < len(lis); i++ {
					key := lis[i]
					obj, o = m.sobjs.Get(key)
					if !o {
						obj, o = m.mobjs.Get(key)
						if !o {
							continue
						}
					}
					layer := obj.StaticInfo().Layer()
					_, y := obj.Pos()
					m.resultLayerObjs[layer].Set(key, y)
				}
			}
		}
	}

	return m.resultLayerObjs
}

func (m *GridMap) GetMovableObjListWithRangeAndSubtype(rect *math.Rect, subtype object.ObjSubType) []uint32 {
	if len(m.resultMovableObjList) > 0 {
		m.resultMovableObjList = m.resultMovableObjList[:0]
	}
	lx, by, rx, ty := m.gridBoundsBy(rect.X(), rect.Y(), rect.X()+rect.W(), rect.Y()+rect.H())
	if rx >= lx && ty >= by {
		for y := by; y <= ty; y++ {
			for x := lx; x <= rx; x++ {
				gidx := m.gridLineCol2Index(y, x)
				lis := m.grids[gidx]
				for i := 0; i < len(lis); i++ {
					key := lis[i]
					obj, o := m.mobjs.Get(key)
					if !o {
						continue
					}
					if subtype != object.ObjSubTypeNone && obj.Subtype() == subtype {
						m.resultMovableObjList = append(m.resultMovableObjList, obj.InstId())
					}
				}
			}
		}
	}
	return m.resultMovableObjList
}

func (m *GridMap) GetMovableObjListWithRange(rect *math.Rect) []uint32 {
	return m.GetMovableObjListWithRangeAndSubtype(rect, object.ObjSubTypeNone)
}

func (m *GridMap) gridLineCol2Index(line, col int16) int32 {
	return int32(line)*int32(m.gridColNum) + int32(col)
}

func (m *GridMap) gridBoundsBy(left, bottom, right, top int32) (lx, by, rx, ty int16) {
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

func (m GridMap) posGridIndex(x, y int32) int32 {
	lx := int16(x / m.gridWidth)
	by := int16(y / m.gridHeight)
	return m.gridLineCol2Index(by, lx)
}

// 遍歷碰撞範圍内的網格檢查碰撞結果 移動之前調用
func (m *GridMap) CheckMovableObjCollision(obj object.IMovableObject /*dir object.Direction, */, dx, dy float64, collisionObj *object.IObject) bool {
	// 是否擁有碰撞組件
	comp := obj.GetComp("Collider")
	if comp == nil {
		return false
	}

	// 九宮格
	var nineSquared = [3][3]int32{{-1, -1, -1}, {-1, -1, -1}, {-1, -1, -1}}
	x, y := obj.Pos()
	index := m.posGridIndex(x, y)
	nineSquared[1][1] = index

	for j := int32(-1); j <= 1; j++ {
		dy := y + j*m.gridHeight
		// y坐標範圍[Y, Y+MapHeight]
		if dy < m.config.Y || dy >= m.config.Y+m.config.TileHeight*int32(len(m.config.Layers)) {
			continue
		}
		for i := int32(-1); i <= 1; i++ {
			dx := x + i*m.gridWidth
			// x坐標範圍[X, X+MapWidth]
			if dx < m.config.X || dx >= m.config.X+m.config.TileWidth*int32(len(m.config.Layers[0])) {
				continue
			}
			nineSquared[1+j][1+i] = index + j*int32(m.gridColNum) + i
		}
	}

	var (
		obj2 object.IObject
		o    bool
	)
	for i := 0; i < len(nineSquared); i++ {
		for j := 0; j < len(nineSquared[i]); j++ {
			index = nineSquared[i][j]
			if index < 0 {
				continue
			}
			ids := m.grids[index]
			for n := 0; n < len(ids); n++ {
				obj2, o = m.mobjs.Get(ids[n])
				if !o {
					obj2, o = m.sobjs.Get(ids[n])
					if !o {
						log.Warn("Collision: grid(x:%v y:%v) not found object %v", x, y, ids[n])
						continue
					}
				}
				if obj2.InstId() != obj.InstId() {
					if checkMovableObjCollisionObj(obj, comp /*dir,*/, dx, dy, obj2) {
						if collisionObj != nil {
							*collisionObj = obj2
						}
						return true
					}
				}
			}
		}
	}
	return false
}
