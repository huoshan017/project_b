package object

// 物体类型
type ObjectType int

const (
	ObjTypeNone    = ObjectType(0) // 无物体
	ObjTypeStatic  = ObjectType(1) // 静止物体
	ObjTypeMovable = ObjectType(2) // 可移动物体
)

// 静止物体类型
type StaticObjType int

const (
	StaticObjNone          = StaticObjType(0) // 无
	StaticObjBrick         = StaticObjType(1) // 砖块
	StaticObjIron          = StaticObjType(2) // 铁
	StaticObjGrass         = StaticObjType(3) // 草
	StaticObjIce           = StaticObjType(4) // 冰
	StaticObjWater         = StaticObjType(5) // 水
	StaticObjHome          = StaticObjType(6) // 基地
	StaticObjHomeDestroyed = StaticObjType(7) // 被摧毁的基地
)

// 可移动物体类型
type MovableObjType int

const (
	MovableObjNone   = MovableObjType(0) // 无
	MovableObjTank   = MovableObjType(1) // 坦克
	MovableObjBullet = MovableObjType(2) // 子弹
)

// 物体子类型
type ObjSubType int

const (
	ObjSubTypeBrick         = ObjSubType(StaticObjBrick)
	ObjSubTypeIron          = ObjSubType(StaticObjIron)
	ObjSubTypeGrass         = ObjSubType(StaticObjGrass)
	ObjSubTypeIce           = ObjSubType(StaticObjIce)
	ObjSubTypeWater         = ObjSubType(StaticObjWater)
	ObjSubTypeHome          = ObjSubType(StaticObjHome)
	ObjSubTypeHomeDestroyed = StaticObjType(StaticObjHomeDestroyed)
	ObjSubTypeTank          = ObjSubType(MovableObjTank)
	ObjSubTypeBullet        = ObjSubType(MovableObjBullet)
)

// 方向
type Direction int

const (
	DirNone  = Direction(0) // 无，用于静止物体
	DirLeft  = Direction(1) // 左
	DirRight = Direction(2) // 右
	DirUp    = Direction(3) // 上
	DirDown  = Direction(4) // 下
	DirMin   = DirLeft      // 最小
	DirMax   = DirDown      // 最大
)

// 其他常量
const (
	Delta                  = 0.00001 // 浮点数精度
	DefaultMinMoveDistance = 1       // 默认最小移动距离
)
