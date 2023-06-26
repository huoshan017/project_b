package object

import "project_b/common/time"

type IRecycle interface {
	ToRecycle()
	IsRecycle() bool
}

// 物体接口
type IObject interface {
	IRecycle
	Init(uint32, *ObjStaticInfo)    // 初始化
	Uninit()                        // 反初始化
	InstId() uint32                 // 实例id
	Id() int32                      // 注意：这是配置id
	Type() ObjectType               // 类型
	Subtype() ObjSubType            // 子类型
	OwnerType() ObjOwnerType        // 所有者类型
	StaticInfo() *ObjStaticInfo     // 靜態信息
	Pos() (x, y int32)              // 位置
	SetPos(x, y int32)              // 设置位
	Center() (x, y int32)           // 中心點坐標
	Width() int32                   // 宽度
	Height() int32                  // 高度
	Left() int32                    // 左坐标
	Right() int32                   // 右坐标
	Top() int32                     // 上坐标
	Bottom() int32                  // 下坐标
	Orientation() int32             // 朝向角度
	Update(tick time.Duration)      // 更新
	AddComp(comp IComponent)        // 添加組件
	RemoveComp(name string)         // 去除組件
	GetComp(name string) IComponent // 獲取組件
	HasComp(name string) bool       // 是否擁有組件
}

// 静态物体接口
type IStaticObject interface {
	IObject
}

// 可移动的物体接口
type IMovableObject interface {
	IObject
	Level() int32          // 等级
	Dir() Direction        // 方向
	Speed() int32          // 速度
	CurrentSpeed() int32   // 当前速度
	Move(Direction)        // 移动
	Stop()                 // 停止
	IsMoving() bool        // 是否在移动
	LastPos() (x, y int32) // 上次Update時位置

	// ----------------------------------
	// 事件接口
	RegisterCheckMoveEventHandle(handle func(args ...any))   // 注冊檢查坐標事件
	UnregisterCheckMoveEventHandle(handle func(args ...any)) // 注銷檢查坐標事件
	RegisterMoveEventHandle(handle func(args ...any))        // 注册移动事件
	UnregisterMoveEventHandle(handle func(args ...any))      // 注销移动事件
	RegisterStopMoveEventHandle(handle func(args ...any))    // 注册停止移动事件
	UnregisterStopMoveEventHandle(handle func(args ...any))  // 注销停止移动事件
	RegisterUpdateEventHandle(handle func(args ...any))      // 注册更新事件
	UnregisterUpdateEventHandle(handle func(args ...any))    // 注销更新事件
}

// 车辆接口
type IVehicle interface {
	IMovableObject
}

// 坦克接口
type ITank interface {
	IVehicle
	Change(info *TankStaticInfo)
	Restore()

	// ---------------------------------
	// 事件接口
	RegisterChangeEventHandle(handle func(args ...any))   // 注册变化事件
	UnregisterChangeEventHandle(handle func(args ...any)) // 注销变化事件
}
