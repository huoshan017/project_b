package common

import (
	"fmt"
	"project_b/common/base"
	"project_b/common/ds"
	"project_b/common/math"
	"project_b/common/object"
	"project_b/game_map"
	"project_b/log"

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
	gridTileSize            int16                                           // 網格瓦片數
	minGridTileSize         int16                                           // 最小网格中的瓦片数
	gridLineNum, gridColNum int16                                           // 網格行數和列數
	gridWidth, gridHeight   int32                                           // 網格寬度高度
	sobjs                   *ds.MapListUnion[uint32, object.IStaticObject]  // 对象map
	mobjs                   *ds.MapListUnion[uint32, object.IMovableObject] // 移動對象map
	grids                   []gridObjList                                   // 網格用於碰撞檢測，提高檢測性能
	mobj2GridIndex          *ds.MapListUnion[uint32, int32]                 // 移動對象與網格索引的映射
	resultLayerObjs         [MapMaxLayer]*heap.BinaryHeapKV[uint32, int32]  // 缓存返回的结果给调用者，主要为了提高性能，减少GC
	resultMovableObjList    []uint32                                        // 緩存搜索的可移動物體的結果
	checkCollisionObjList   []object.IObject                                // 檢測碰撞物體列表
	nineSquared             [3][3]int32                                     // 九宮格
	resultObjList           []uint32                                        // 緩存物體id列表
}

func NewGridMap(gridTileSize int16) *GridMap {
	gmap := &GridMap{
		gridTileSize:    gridTileSize,
		minGridTileSize: gridTileSize,
		sobjs:           ds.NewMapListUnion[uint32, object.IStaticObject](),
		mobjs:           ds.NewMapListUnion[uint32, object.IMovableObject](),
		mobj2GridIndex:  ds.NewMapListUnion[uint32, int32](),
	}
	return gmap
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
	gridTileNum4Width, gridTileNum4Height := m.minGridTileSize, m.minGridTileSize
	for colNum%gridTileNum4Width > 0 && colNum%gridTileNum4Width < m.minGridTileSize {
		gridTileNum4Width += 1
	}
	for lineNum%gridTileNum4Height > 0 && lineNum%gridTileNum4Height < m.minGridTileSize {
		gridTileNum4Height += 1
	}

	// 創建網格
	gridColNum, gridLineNum := (colNum+gridTileNum4Width-1)/gridTileNum4Width, (lineNum+gridTileNum4Height-1)/gridTileNum4Height
	m.grids = make([]gridObjList, gridLineNum*gridColNum)
	m.gridLineNum = gridLineNum
	m.gridColNum = gridColNum

	// 網格寬度和高度
	m.gridWidth = int32(gridTileNum4Width) * m.config.TileWidth
	m.gridHeight = int32(gridTileNum4Height) * m.config.TileHeight

	// 創建結果層
	for i := 0; i < len(m.resultLayerObjs); i++ {
		m.resultLayerObjs[i] = heap.NewMaxBinaryHeapKV[uint32, int32]()
	}
}

func (m *GridMap) Unload() {
	m.ClearObjsData()
	m.config = nil
	m.grids = nil
	m.minGridTileSize = m.gridTileSize
	m.gridLineNum = 0
	m.gridColNum = 0
	m.gridWidth = 0
	m.gridHeight = 0
}

func (m *GridMap) ClearObjsData() {
	for i := 0; i < len(m.grids); i++ {
		m.grids[i].clear()
	}
	m.sobjs.Clear()
	m.mobjs.Clear()
	m.mobj2GridIndex.Clear()
	for i := 0; i < len(m.resultLayerObjs); i++ {
		m.resultLayerObjs[i].Clear()
	}
	if len(m.resultMovableObjList) > 0 {
		clear(m.resultMovableObjList)
		m.resultMovableObjList = m.resultMovableObjList[:0]
	}
	if len(m.checkCollisionObjList) > 0 {
		clear(m.checkCollisionObjList)
		m.checkCollisionObjList = m.checkCollisionObjList[:0]
	}
}

func (m *GridMap) GetGridWidthHeight() (int32, int32) {
	return m.gridWidth, m.gridHeight
}

func (m *GridMap) AddTile(line, col int16, tile object.IObject) {
	m.AddObj(tile)
}

func (m *GridMap) RemoveTile(instId uint32) {
	m.RemoveObj(instId)
}

func (m *GridMap) AddObj(obj object.IObject) {
	instId := obj.InstId()
	if obj.Type() == base.ObjTypeMovable {
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
	if !m.grids[index].addObj(instId) {
		panic(fmt.Sprintf("GridMap: obj %v already exists in grid %v", instId, index))
	}
	if obj.Type() == base.ObjTypeMovable {
		m.mobjs.Add(instId, obj.(object.IMovableObject))
		m.mobj2GridIndex.Add(instId, index)
	} else {
		m.sobjs.Add(instId, obj.(object.IStaticObject))
	}
	log.Info("GridMap: obj %v(type:%v, subtype:%v) add to grid(index: %v)", instId, obj.Type(), obj.Subtype(), index)
}

func (m *GridMap) RemoveObj(instId uint32) {
	var (
		obj object.IObject
		o   bool
	)
	obj, o = m.mobjs.Get(instId)
	if o {
		m.mobjs.Remove(instId)
		index, o := m.mobj2GridIndex.Get(instId)
		if !o {
			log.Error("GridMap: obj %v cant get grid index to remove", instId)
			return
		}
		if !m.grids[index].removeObj(instId) {
			panic(fmt.Sprintf("GridMap.RemoveObj: remove obj %v from grid %v failed", instId, index))
		}
		m.mobj2GridIndex.Remove(instId)
		log.Info("GridMap: obj %v(type:%v, subtype:%v) remove from grid(index: %v)", instId, obj.Type(), obj.Subtype(), index)
	} else {
		obj, o = m.sobjs.Get(instId)
		if !o {
			return
		}
		m.sobjs.Remove(instId)
	}
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
	lastIndex, o := m.mobj2GridIndex.Get(obj.InstId())
	if !o {
		log.Warn("GridMap: obj %v not found grid index to update", obj.InstId())
		return
	}
	if index != lastIndex {
		if !m.grids[lastIndex].removeObj(obj.InstId()) {
			panic(fmt.Sprintf("GridMap.UpdateMovable: remove obj %v from grid %v failed", obj.InstId(), lastIndex))
		}
		if !m.grids[index].addObj(obj.InstId()) {
			panic(fmt.Sprintf("GridMap.UpdateMovable: add obj %v to grid %v failed", obj.InstId(), index))
		}
		m.mobj2GridIndex.Set(obj.InstId(), index)
		if obj.Subtype() == base.ObjSubtypeTank {
			log.Info("GridMap: obj %v(type:%v, subtype:%v) remove from grid(%v), add to grid(%v)", obj.InstId(), obj.Type(), obj.Subtype(), lastIndex, index)
		}
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

func (m *GridMap) PointInObjList(pos base.Pos) []uint32 {
	if len(m.resultObjList) > 0 {
		clear(m.resultObjList)
		m.resultObjList = m.resultObjList[:0]
	}
	index := m.posGridIndex(pos.X, pos.Y)
	objIds := m.grids[index]
	for i := 0; i < len(objIds); i++ {
		obj, o := m.sobjs.Get(objIds[i])
		if !o {
			obj, o = m.mobjs.Get(objIds[i])
			if !o {
				continue
			}
		}
		if object.IsPointInObject(pos, obj) {
			m.resultObjList = append(m.resultObjList, objIds[i])
		}
	}
	return m.resultObjList
}

func (m *GridMap) GetMovableObjListWithRangeAndSubtype(rect *math.Rect, subtype base.ObjSubtype) []uint32 {
	if len(m.resultMovableObjList) > 0 {
		m.resultMovableObjList = m.resultMovableObjList[:0]
	}
	var x, y, w, h = rect.X(), rect.Y(), rect.W(), rect.H()
	if x < m.config.X {
		x = m.config.X
	}
	if x > m.config.X+int32(len(m.config.Layers[0]))*m.config.TileWidth {
		x = m.config.X + int32(len(m.config.Layers[0]))*m.config.TileWidth
	}
	if y < m.config.Y {
		y = m.config.Y
	}
	if y > m.config.Y+int32(len(m.config.Layers))*m.config.TileHeight {
		y = m.config.Y + int32(len(m.config.Layers))*m.config.TileHeight
	}
	lx, by, rx, ty := m.gridBoundsBy(x, y, x+w, y+h)
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
					if subtype != base.ObjSubtypeNone && obj.Subtype() == subtype {
						m.resultMovableObjList = append(m.resultMovableObjList, obj.InstId())
					}
				}
			}
		}
	}
	return m.resultMovableObjList
}

func (m *GridMap) GetMovableObjListWithRange(rect *math.Rect) []uint32 {
	return m.GetMovableObjListWithRangeAndSubtype(rect, base.ObjSubtypeNone)
}

func (m *GridMap) GetObjListWithLineSegment(start, end *base.Pos) (objList []uint32) {
	sx, sy := start.X, start.Y
	ex, ey := end.X, end.Y
	var (
		dx, dy                 int32 = ex - sx, ey - sy
		is, ie, js, je, di, dj int32
		flag                   bool
	)
	// check abs_dx/abs_dy and m.gridWidth/m.gridHeight
	// abs_dx/abs_dy >= m.gridWidth/m.gridHeight
	abs_dx, abs_dy := base.Abs(dx), base.Abs(dy)
	// 計算從start到end的綫段在水平和垂直兩個方向上的步進增量
	// 以一個方向(x軸或y軸)網格的寬度或高度爲增量遍歷經過的網格
	// 而另一個方向(y軸或x軸)遍歷時的增量按網格寬高比計算，不超
	// 過網格的高度或寬度，這樣才能保證遍歷所有經過的網格，而不
	// 遺漏掉任何網格
	if abs_dx*m.gridHeight >= abs_dy*m.gridWidth {
		if dx >= 0 {
			if dy >= 0 {
				is, ie, js, je = sx, ex, sy, ey
				dj = abs_dy * m.gridWidth / abs_dx
			} else {
				is, ie, js, je = sx, ex, ey, sy
				dj = -(abs_dy * m.gridWidth / abs_dx)
			}
			di = m.gridWidth
		} else {
			if dy >= 0 {
				is, ie, js, je = ex, sx, sy, ey
				dj = abs_dy * m.gridWidth / abs_dx
			} else {
				is, ie, js, je = ex, sx, ey, sy
				dj = -(abs_dy * m.gridWidth / abs_dx)
			}
			di = -m.gridWidth
		}
	} else { // abs_dx/abs_dy >= m.gridWidth/m.gridHeight
		if dy >= 0 {
			if dx >= 0 {
				is, ie, js, je = sy, ey, sx, ex
				dj = abs_dx * m.gridHeight / abs_dy
			} else {
				is, ie, js, je = sy, ey, ex, sx
				dj = -(abs_dx * m.gridHeight / abs_dy)
			}
			di = m.gridHeight
		} else {
			if dx >= 0 {
				is, ie, js, je = ey, sy, sx, ex
				dj = abs_dx * m.gridHeight / abs_dy
			} else {
				is, ie, js, je = ey, sy, ex, sx
				dj = -(abs_dx * m.gridHeight / abs_dy)
			}
			di = -m.gridHeight
		}
		flag = true
	}
	var (
		index int32
		i, j  = is, js
	)
	for i = is; ; i += di {
		for j = js; ; j += dj {
			if !flag {
				index = m.posGridIndex(i, j)
			} else {
				index = m.posGridIndex(j, i)
			}
			l := m.grids[index]
			for k := 0; k < len(l); k++ {
				// 判斷綫段跟物體是否有相交
				obj, o := m.sobjs.Get(l[k])
				if o && object.CheckLineSegmentIntersectObj(start, end, obj) {
					objList = append(objList, l[k])
				}
			}
			if dj == 0 || (dj > 0 && j+dj > je) || (dj < 0 && j+dj < je) {
				break
			}
		}
		if di == 0 || (di > 0 && i+di > ie) || (di < 0 && i+di < ie) {
			break
		}
	}
	return
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

func (m GridMap) pos2LineColumn(x, y int32) (int16, int16) {
	cx := (x - m.config.X) / m.gridWidth
	ly := (y - m.config.Y) / m.gridHeight
	return int16(ly), int16(cx)
}

func (m GridMap) posGridIndex(x, y int32) int32 {
	l, c := m.pos2LineColumn(x, y)
	return m.gridLineCol2Index(l, c)
}

func (m *GridMap) gridLineCol2Index(line, col int16) int32 {
	return int32(line)*int32(m.gridColNum) + int32(col)
}

func (m GridMap) gridIndex2LineColumn(index int32) (int32, int32) {
	return index / int32(m.gridColNum), index % int32(m.gridColNum)
}

func (m *GridMap) getNineSquaredGridIndexWithPos(x, y int32) {
	line, column := m.pos2LineColumn(x, y)
	for l := int16(-1); l <= 1; l++ {
		if line+l < 0 || line+l >= m.gridLineNum {
			for c := int16(-1); c <= 1; c++ {
				m.nineSquared[l-(-1)][c-(-1)] = -1
			}
			continue
		}
		for c := int16(-1); c <= 1; c++ {
			if column+c < 0 || column+c >= m.gridColNum {
				m.nineSquared[l-(-1)][c-(-1)] = -1
				continue
			}
			m.nineSquared[l-(-1)][c-(-1)] = m.gridLineCol2Index(line+l, column+c)
		}
	}
}

// 遍歷碰撞範圍内的網格檢查碰撞結果 移動之前調用
func (m *GridMap) CheckMovingObjCollision(mobj object.IMovableObject, dx, dy int32, collisionInfo *object.CollisionInfo) object.CollisionResult {
	// 是否擁有碰撞組件
	if mobj.GetColliderComp() == nil {
		return object.CollisionNone
	}

	x, y := mobj.Pos()
	x += dx
	y += dy
	var index int32

	// 九宮格
	m.getNineSquaredGridIndexWithPos(x, y)

	var (
		obj object.IObject
		o   bool
	)
	if len(m.checkCollisionObjList) > 0 {
		clear(m.checkCollisionObjList)
		m.checkCollisionObjList = m.checkCollisionObjList[:0]
	}
	for i := 0; i < len(m.nineSquared); i++ {
		for j := 0; j < len(m.nineSquared[i]); j++ {
			index = m.nineSquared[i][j]
			if index < 0 {
				continue
			}
			ids := m.grids[index]
			for n := 0; n < len(ids); n++ {
				obj, o = m.mobjs.Get(ids[n])
				if !o {
					obj, o = m.sobjs.Get(ids[n])
					if !o {
						gl, gc := m.gridIndex2LineColumn(index)
						log.Warn("Collision: grid(x:%v y:%v, index:%v) exist object %v, but not found in obj map", gc, gl, index, ids[n])
						continue
					}
				}
				if obj.InstId() != mobj.InstId() {
					cr := object.CheckMovingObjCollisionObj(mobj, dx, dy, obj)
					if cr == object.CollisionAndBlock {
						m.checkCollisionObjList = append(m.checkCollisionObjList, obj)
					} else if cr == object.CollisionOnly {
						collisionInfo.ObjList = append(collisionInfo.ObjList, obj)
					}
				}
			}
		}
	}

	if len(m.checkCollisionObjList) == 0 {
		if len(collisionInfo.ObjList) == 0 {
			return object.CollisionNone
		}
		collisionInfo.MovingObj = mobj
		collisionInfo.Result = object.CollisionOnly
		return object.CollisionOnly
	}

	return object.NarrowPhaseCheckMovingObjCollision2ObjList(mobj, dx, dy, m.checkCollisionObjList, collisionInfo)
}
