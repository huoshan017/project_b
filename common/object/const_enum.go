package object

// 物体类型
type ObjectType int

const (
	ObjTypeNone    = ObjectType(0) // 无物体
	ObjTypeStatic  = ObjectType(1) // 静止物体
	ObjTypeMovable = ObjectType(2) // 可移动物体
	ObjTypeEnumMax = 3
)

// 静止物体类型
type StaticObjType int

const (
	StaticObjNone    = StaticObjType(0) // 无
	StaticObjBrick   = StaticObjType(1) // 砖块
	StaticObjIron    = StaticObjType(2) // 铁
	StaticObjTree    = StaticObjType(3) // 樹
	StaticObjWater   = StaticObjType(4) // 水
	StaticObjIce     = StaticObjType(5) // 冰
	StaticObjHome    = StaticObjType(6) // 基地
	StaticObjRuins   = StaticObjType(7) // 被摧毁的基地
	StaticObjEnumMax = 9
)

// 可移动物体类型
type MovableObjType int

const (
	MovableObjNone        = MovableObjType(0) // 无
	MovableObjTank        = MovableObjType(1) // 坦克
	MovableObjShell       = MovableObjType(2) // 炮弹
	MovableObjSurroundObj = MovableObjType(3) // 環繞物體，測試用
	MovableObjEnumMax     = 4
)

// 物体子类型
type ObjSubtype int

const (
	ObjSubtypeNone        = iota
	ObjSubtypeBrick       = ObjSubtype(StaticObjBrick)
	ObjSubtypeIron        = ObjSubtype(StaticObjIron)
	ObjSubtypeTree        = ObjSubtype(StaticObjTree)
	ObjSubtypeIce         = ObjSubtype(StaticObjIce)
	ObjSubtypeWater       = ObjSubtype(StaticObjWater)
	ObjSubtypeHome        = ObjSubtype(StaticObjHome)
	ObjSubtypeRuins       = ObjSubtype(StaticObjRuins)
	ObjSubtypeTank        = ObjSubtype(MovableObjTank)
	ObjSubtypeShell       = ObjSubtype(MovableObjShell)
	ObjSubtypeSurroundObj = ObjSubtype(MovableObjSurroundObj)
)

// 对象所有者类型
type ObjOwnerType int

const (
	OwnerNone       = ObjOwnerType(0) // 无所有者或者系统所有
	OwnerPlayer     = ObjOwnerType(1) // 玩家
	OwnerBOT        = ObjOwnerType(2) // BOT
	OwnerBOT4Player = ObjOwnerType(3) // 玩家拥有的BOT对象
)

// 方向
type Direction int

const (
	DirNone      = Direction(0) // 无，用于静止物体
	DirLeft      = Direction(1) // 左
	DirRight     = Direction(2) // 右
	DirUp        = Direction(3) // 上
	DirDown      = Direction(4) // 下
	DirLeftUp    = Direction(5) // 左上
	DirLeftDown  = Direction(6) // 左下
	DirRightUp   = Direction(7) // 右上
	DirRightDown = Direction(8) // 右下
	DirMin       = DirNone      // 最小
	DirMax       = DirDown      // 最大
)

// 其他常量
const (
	DefaultMinMoveDistance = 1  // 默认最小移动距离
	BotTileFlag            = -1 // bot瓦片位置標記
	PlayerTileFlag         = -2 // 玩家瓦片位置標記
)

// 效果作用類型
type EffectType int

const (
	EffectTypeTime     = iota // 時間
	EffectTypeRequency = 1    // 次數
)

// 陣營類型
type CampType int

const (
	CampTypeNone  = iota // 無陣營
	CampTypeOne   = 1    // 陣營1
	CampTypeTwo   = 2    // 陣營2
	CampTypeThree = 3    // 陣營3
	CampTypeFour  = 4    // 陣營4
)

func Dir2Orientation(dir Direction) int32 {
	switch dir {
	case DirLeft:
		return 180
	case DirRight:
		return 0
	case DirUp:
		return 90
	case DirDown:
		return 270
	case DirLeftUp:
		return 135
	case DirLeftDown:
		return 225
	case DirRightUp:
		return 45
	case DirRightDown:
		return 315
	}
	return 0
}
