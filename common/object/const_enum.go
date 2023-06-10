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
	StaticObjNone          = StaticObjType(0) // 无
	StaticObjBrick         = StaticObjType(1) // 砖块
	StaticObjIron          = StaticObjType(2) // 铁
	StaticObjTree          = StaticObjType(3) // 樹
	StaticObjWater         = StaticObjType(4) // 水
	StaticObjIce           = StaticObjType(5) // 冰
	StaticObjHome          = StaticObjType(6) // 基地
	StaticObjHomeDestroyed = StaticObjType(7) // 被摧毁的基地
	StaticObjEnumMax       = 8
)

// 可移动物体类型
type MovableObjType int

const (
	MovableObjNone    = MovableObjType(0) // 无
	MovableObjTank    = MovableObjType(1) // 坦克
	MovableObjBullet  = MovableObjType(2) // 子弹
	MovableObjEnumMax = 3
)

// 物体子类型
type ObjSubType int

const (
	ObjSubTypeBrick         = ObjSubType(StaticObjBrick)
	ObjSubTypeIron          = ObjSubType(StaticObjIron)
	ObjSubTypeTree          = ObjSubType(StaticObjTree)
	ObjSubTypeIce           = ObjSubType(StaticObjIce)
	ObjSubTypeWater         = ObjSubType(StaticObjWater)
	ObjSubTypeHome          = ObjSubType(StaticObjHome)
	ObjSubTypeHomeDestroyed = ObjSubType(StaticObjHomeDestroyed)
	ObjSubTypeTank          = ObjSubType(MovableObjTank)
	ObjSubTypeBullet        = ObjSubType(MovableObjBullet)
)

// 对象所有者类型
type ObjOwnerType int

const (
	OwnerNone      = ObjOwnerType(0) // 无所有者或者系统所有
	OwnerPlayer    = ObjOwnerType(1) // 玩家
	OwnerAI        = ObjOwnerType(2) // AI
	OwnerAI4Player = ObjOwnerType(3) // 玩家拥有的AI对象
)

// 相对关系类型
type RelationType int

const (
	RelationNeutral = RelationType(0) // 中立
	RelationOwn     = RelationType(1) // 我方
	RelationEnemy   = RelationType(2) // 敌方
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
	DefaultMinMoveDistance = 1 // 默认最小移动距离
)
