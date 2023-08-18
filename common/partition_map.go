package common

import (
	"fmt"
	"project_b/common/ds"
	"project_b/common/math"
	"project_b/common/object"
	"project_b/game_map"
	"project_b/log"

	"github.com/huoshan017/ponu/heap"
)

const (
	minGridTileSize     int16 = 5               // 最小網格大小(tile為單位)
	defaultGridTileSize int16 = minGridTileSize // 缺省網格大小
	MapMaxLayer         int16 = 5               // 地圖最大層
)

type gridObjs struct {
	sobjs *ds.MapListUnion[uint32, struct{}]
	mobjs *ds.MapListUnion[uint32, struct{}]
}

func newGridObjs() gridObjs {
	return gridObjs{
		sobjs: ds.NewMapListUnion[uint32, struct{}](),
		mobjs: ds.NewMapListUnion[uint32, struct{}](),
	}
}

func (g *gridObjs) addObj(obj object.IObject) {
	if obj.Type() == object.ObjTypeStatic {
		g.sobjs.Add(obj.InstId(), struct{}{})
	} else if obj.Type() == object.ObjTypeMovable {
		g.mobjs.Add(obj.InstId(), struct{}{})
	} else {
		panic(fmt.Sprintf("invalid object type %v", obj.Type()))
	}
}

func (g *gridObjs) removeObj(instId uint32) {
	if _, o := g.sobjs.Get(instId); o {
		g.sobjs.Remove(instId)
	} else {
		g.mobjs.Remove(instId)
	}
}

func (g *gridObjs) getSObjs() *ds.MapListUnion[uint32, struct{}] {
	return g.sobjs
}

func (g *gridObjs) getMObjs() *ds.MapListUnion[uint32, struct{}] {
	return g.mobjs
}

// 分割地圖
type PartitionMap struct {
	config                  *game_map.Config
	minGridTileSize         int16                                           // 最小网格中的瓦片数
	gridXTiles, gridYTiles  int16                                           // x轴网格瓦片数，y轴网格瓦片数
	gridLineNum, gridColNum int16                                           // 網格行數和列數
	gridWidth, gridHeight   int32                                           // 網格寬度高度
	sobjs                   *ds.MapListUnion[uint32, object.IStaticObject]  // 对象map
	mobjs                   *ds.MapListUnion[uint32, object.IMovableObject] // 移動對象map
	grids                   []gridObjs                                      // 網格用於碰撞檢測，提高檢測性能
	resultLayerObjs         [MapMaxLayer]*heap.BinaryHeapKV[uint32, int32]  // 缓存返回的结果给调用者，主要为了提高性能，减少GC
	resultMovableObjList    []uint32                                        // 緩存搜索的可移動物體的結果
}

func NewPartitionMap(gridTileSize int16) *PartitionMap {
	if gridTileSize < minGridTileSize {
		gridTileSize = minGridTileSize
	}
	return &PartitionMap{
		minGridTileSize: gridTileSize,
		sobjs:           ds.NewMapListUnion[uint32, object.IStaticObject](),
		mobjs:           ds.NewMapListUnion[uint32, object.IMovableObject](),
	}
}

func (m *PartitionMap) Load(config *game_map.Config) {
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
	m.grids = make([]gridObjs, gridLineNum*gridColNum)
	for i := 0; i < int(gridLineNum*gridColNum); i++ {
		m.grids[i] = newGridObjs()
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

func (m *PartitionMap) Unload() {
	m.config = nil
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

func (m *PartitionMap) AddTile(line, col int16, tile object.IObject) {
	m.AddObj(tile)
}

func (m *PartitionMap) RemoveTile(instId uint32) {
	m.RemoveObj(instId)
}

func (m *PartitionMap) AddObj(obj object.IObject) {
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
	lx, by, rx, ty := m.objGridBounds(obj)
	if lx <= rx && by <= ty {
		index := m.gridLineCol2Index(by, lx)
		m.grids[index].addObj(obj)
		log.Info("MapInstance: obj %v add to grid(line: %v, col: %v, index: %v)", instId, by, lx, index)
		if lx != rx && by == ty {
			index = m.gridLineCol2Index(by, rx)
			m.grids[index].addObj(obj)
			log.Info("MapInstance: obj %v add to grid(line: %v, col: %v, index: %v)", instId, by, rx, index)
		}
		if lx != rx && by != ty {
			index = m.gridLineCol2Index(ty, rx)
			m.grids[m.gridLineCol2Index(ty, rx)].addObj(obj)
			log.Info("MapInstance: obj %v add to grid(line: %v, col: %v, index: %v)", instId, ty, rx, index)
		}
		if lx == rx && by != ty {
			index = m.gridLineCol2Index(ty, lx)
			m.grids[index].addObj(obj)
			log.Info("MapInstance: obj %v add to grid(line: %v, col: %v, index: %v)", instId, ty, lx, index)
		}
	}
	if obj.Type() == object.ObjTypeMovable {
		m.mobjs.Add(instId, obj.(object.IMovableObject))
	} else {
		m.sobjs.Add(instId, obj.(object.IStaticObject))
	}
}

func (m *PartitionMap) RemoveObj(instId uint32) {
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
	lx, by, rx, ty := m.objGridBounds(obj)
	if lx <= rx && by <= ty {
		index := m.gridLineCol2Index(by, lx)
		m.grids[index].removeObj(instId)
		log.Info("MapInstance: obj %v remove from grid(line: %v, col: %v, index: %v)", instId, by, lx, index)
		if lx != rx && by == ty {
			index = m.gridLineCol2Index(by, rx)
			m.grids[index].removeObj(instId)
			log.Info("MapInstance: obj %v remove from grid(line: %v, col: %v, index: %v)", instId, by, rx, index)
		}
		if lx != rx && by != ty {
			index = m.gridLineCol2Index(ty, rx)
			m.grids[index].removeObj(instId)
			log.Info("MapInstance: obj %v remove from grid(line: %v, col: %v, index: %v)", instId, ty, rx, index)
		}
		if lx == rx && by != ty {
			index = m.gridLineCol2Index(ty, lx)
			m.grids[index].removeObj(instId)
			log.Info("MapInstance: obj %v remove from grid(line: %v, col: %v, index: %v)", instId, ty, lx, index)
		}
	}
}

func (m *PartitionMap) UpdateMovable(obj object.IMovableObject) {
	if !m.mobjs.Exists(obj.InstId()) {
		log.Warn("obj %v (type: %v, subtype: %v) not found in map", obj.InstId(), obj.Type(), obj.Subtype())
		return
	}

	x, y := obj.Pos()
	lastX, lastY := obj.LastPos()
	if x == lastX && y == lastY {
		return
	}

	left1, bottom1, right1, top1 := m.gridBoundsBy(x, y, x+obj.Width(), y+obj.Length())
	left0, bottom0, right0, top0 := m.gridBoundsBy(lastX, lastY, lastX+obj.Width(), lastY+obj.Length())

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
			m.grids[index].removeObj(obj.InstId())
			log.Info("MapInstance.UpdateObj: obj %v remove from grid(line: %v, col: %v, index: %v)", obj.InstId(), bottom0, left0, index)
			if bottom0 != top0 { // 説明左下和左上在不同的grid中，離開 (top0, left0)
				index = m.gridLineCol2Index(top0, left0)
				m.grids[index].removeObj(obj.InstId())
				log.Info("MapInstance.UpdateObj: obj %v remove from grid(line: %v, col: %v, index: %v)", obj.InstId(), top0, left0, index)
			}
		}
		if right1 > right0 { // 進入 column right1
			index = m.gridLineCol2Index(bottom1, right1)
			m.grids[index].addObj(obj)
			log.Info("MapInstance.UpdateObj: obj %v add to grid(line: %v, col: %v, index: %v)", obj.InstId(), bottom1, right1, index)
			if bottom1 != top1 { // 右下和右上在不同的grid中，右上進入(top1, right1)
				index = m.gridLineCol2Index(top1, right1)
				m.grids[index].addObj(obj)
				log.Info("MapInstance.UpdateObj: obj %v add to grid(line: %v, col: %v, index: %v)", obj.InstId(), top1, right1, index)
			}
		}
	} else if x < lastX { // 向左移動
		if left1 < left0 { //
			index = m.gridLineCol2Index(bottom1, left1)
			m.grids[index].addObj(obj)
			log.Info("MapInstance.UpdateObj: obj %v add to grid(line: %v, col: %v, index: %v)", obj.InstId(), bottom1, left1, index)
			if bottom1 != top1 { // 左下和左上在不同的grid中，左上進入(top1, left1)
				index = m.gridLineCol2Index(top1, left1)
				m.grids[index].addObj(obj)
				log.Info("MapInstance.UpdateObj: obj %v add to grid(line: %v, col: %v, index: %v)", obj.InstId(), top1, left1, index)
			}
		}
		if right1 < right0 { // 右側離開原來的grid
			index = m.gridLineCol2Index(bottom0, right0)
			m.grids[index].removeObj(obj.InstId())
			log.Info("MapInstance.UpdateObj: obj %v remove from grid(line: %v, col: %v, index: %v)", obj.InstId(), bottom0, right0, index)
			if top0 != bottom0 { // 右上和右下在不同grid中，離開(top0, right0)
				index = m.gridLineCol2Index(top0, right0)
				m.grids[index].removeObj(obj.InstId())
				log.Info("MapInstance.UpdateObj: obj %v remove from grid(line: %v, col: %v, index: %v)", obj.InstId(), top0, right0, index)
			}
		}
	}

	if y > lastY { // 向上移動
		if bottom1 > bottom0 { // 向上移動，左下離開(bottom0, left0)
			index = m.gridLineCol2Index(bottom0, left0)
			m.grids[index].removeObj(obj.InstId())
			log.Info("MapInstance.UpdateObj: obj %v remove from grid(line: %v, col: %v, index: %v)", obj.InstId(), bottom0, left0, index)
			if left0 != right0 { // 左下和右下在不同的grid中，離開(bottom0, right0)
				index = m.gridLineCol2Index(bottom0, right0)
				m.grids[index].removeObj(obj.InstId())
				log.Info("MapInstance.UpdateObj: obj %v remove from grid(line: %v, col: %v, index: %v)", obj.InstId(), bottom0, right0, index)
			}
		}
		if top1 > top0 { // 頂部進入新grid
			index = m.gridLineCol2Index(top1, left1)
			m.grids[index].addObj(obj)
			log.Info("MapInstance.UpdateObj: obj %v add to grid(line: %v, col: %v, index: %v)", obj.InstId(), top1, left1, index)
			if left1 != right1 { // 左上和右上在不同的grid中，右上進入(top1, right1)
				index = m.gridLineCol2Index(top1, right1)
				m.grids[index].addObj(obj)
				log.Info("MapInstance.UpdateObj: obj %v add to grid(line: %v, col: %v, index: %v)", obj.InstId(), top1, right1, index)
			}
		}
	} else if y < lastY { // 向下移動
		if bottom1 < bottom0 { // 向上移動，左下進入(bottom1, left1)
			index = m.gridLineCol2Index(bottom1, left1)
			m.grids[index].addObj(obj)
			log.Info("MapInstance.UpdateObj: obj %v add to grid(line: %v, col: %v, index: %v)", obj.InstId(), bottom1, left1, index)
			if left1 != right1 { // 左下和右上在不同的grid中，右下也要進入(bottom1, right1)
				index = m.gridLineCol2Index(bottom1, right1)
				m.grids[index].addObj(obj)
				log.Info("MapInstance.UpdateObj: obj %v add to grid(line: %v, col: %v, index: %v)", obj.InstId(), bottom1, right1, index)
			}
		}
		if top1 < top0 { // 頂部左上離開(top0, left0)
			index = m.gridLineCol2Index(top0, left0)
			m.grids[index].removeObj(obj.InstId())
			log.Info("MapInstance.UpdateObj: obj %v remove from grid(line: %v, col: %v, index: %v)", obj.InstId(), top0, left0, index)
			if left0 != right0 { // 左上和右上在不同的grid中，右上離開(top0, right0)
				index = m.gridLineCol2Index(top0, right0)
				m.grids[index].removeObj(obj.InstId())
				log.Info("MapInstance.UpdateObj: obj %v remove from grid(line: %v, col: %v, index: %v)", obj.InstId(), top0, right0, index)
			}
		}
	}
}

func (m *PartitionMap) GetLayerObjsWithRange(rect *math.Rect) [MapMaxLayer]*heap.BinaryHeapKV[uint32, int32] {
	// 獲得範圍内的靜止和運動的obj
	lx, by, rx, ty := m.gridBoundsBy(rect.X(), rect.Y(), rect.X()+rect.W(), rect.Y()+rect.H())
	if rx >= lx && ty >= by {
		for y := by; y <= ty; y++ {
			for x := lx; x <= rx; x++ {
				gidx := m.gridLineCol2Index(y, x)
				lis := m.grids[gidx].getMObjs().GetList()
				for i := 0; i < len(lis); i++ {
					key := lis[i].Key
					obj, o := m.mobjs.Get(key)
					if !o {
						continue
					}
					layer := obj.StaticInfo().Layer()
					_, y := obj.Pos()
					m.resultLayerObjs[layer].Set(key, y)
				}
				lis = m.grids[gidx].getSObjs().GetList()
				for i := 0; i < len(lis); i++ {
					key := lis[i].Key
					obj, o := m.sobjs.Get(key)
					if !o {
						continue
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

func (m *PartitionMap) GetMovableObjListWithRangeAndSubtype(rect *math.Rect, subtype object.ObjSubtype) []uint32 {
	if len(m.resultMovableObjList) > 0 {
		m.resultMovableObjList = m.resultMovableObjList[:0]
	}
	lx, by, rx, ty := m.gridBoundsBy(rect.X(), rect.Y(), rect.X()+rect.W(), rect.Y()+rect.H())
	if rx >= lx && ty >= by {
		for y := by; y <= ty; y++ {
			for x := lx; x <= rx; x++ {
				gidx := m.gridLineCol2Index(y, x)
				lis := m.grids[gidx].getMObjs().GetList()
				for i := 0; i < len(lis); i++ {
					key := lis[i].Key
					obj, o := m.mobjs.Get(key)
					if !o {
						continue
					}
					if subtype != object.ObjSubtypeNone && obj.Subtype() == subtype {
						m.resultMovableObjList = append(m.resultMovableObjList, obj.InstId())
					}
				}
			}
		}
	}
	return m.resultMovableObjList
}

func (m *PartitionMap) GetMovableObjListWithRange(rect *math.Rect) []uint32 {
	return m.GetMovableObjListWithRangeAndSubtype(rect, object.ObjSubtypeNone)
}

// checkMovableObjCollision 遍歷碰撞範圍内的網格檢查碰撞結果 移動之前調用
func (m *PartitionMap) CheckMovableObjCollision(obj object.IMovableObject, dir object.Direction, dx, dy int32, collisionObj *object.IObject) bool {
	// 獲取檢測碰撞範圍
	lx, by, rx, ty := m.objGridBounds(obj)
	if rx < lx || ty < by {
		return false
	}

	for y := by; y <= ty; y++ {
		for x := lx; x <= rx; x++ {
			gidx := m.gridLineCol2Index(y, x)
			lis := m.grids[gidx].getMObjs().GetList()
			for i := 0; i < len(lis); i++ {
				item := lis[i]
				obj2, o := m.mobjs.Get(item.Key)
				if !o {
					log.Warn("Collision: grid(x:%v y:%v) not found movable object %v", x, y, item.Key)
					continue
				}
				if obj2.InstId() != obj.InstId() && obj2.StaticInfo().Layer() == obj.StaticInfo().Layer() {
					if object.CheckMovingObjCollisionObj(obj, dx, dy, obj2) != object.CollisionNone {
						if collisionObj != nil {
							*collisionObj = obj2
						}
						return true
					}
				}
			}

			lis = m.grids[gidx].getSObjs().GetList()
			for i := 0; i < len(lis); i++ {
				item := lis[i]
				obj2, o := m.sobjs.Get(item.Key)
				if !o {
					log.Warn("Collision: grid(x:%v y:%v) not found static object %v", x, y, item.Key)
					continue
				}
				if object.CheckMovingObjCollisionObj(obj, dx, dy, obj2) != object.CollisionNone {
					if collisionObj != nil {
						*collisionObj = obj2
					}
					return true
				}
			}
		}
	}
	return false
}

func (m *PartitionMap) gridLineCol2Index(line, col int16) int32 {
	return int32(line)*int32(m.gridColNum) + int32(col)
}

func (m *PartitionMap) objGridBounds(obj object.IObject) (lx, by, rx, ty int16) {
	left, bottom := obj.Pos()
	right, top := left+obj.Width()-1, bottom+obj.Length()-1
	return m.gridBoundsBy(left, bottom, right, top)
}

func (m *PartitionMap) gridBoundsBy(left, bottom, right, top int32) (lx, by, rx, ty int16) {
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
