package object

import "project_b/common/time"

// 物体接口
type IObject interface {
	Init()                     // 初始化
	Uninit()                   // 反初始化
	InstId() uint64            // 实例id
	Id() int32                 // 注意：这是配置id
	Type() ObjectType          // 类型
	Subtype() ObjSubType       // 子类型
	SetPos(x, y int32)         // 设置位置
	Pos() (x, y int32)         // 位置
	Width() int32              // 宽度
	Height() int32             // 高度
	Left() int32               // 左坐标
	Right() int32              // 右坐标
	Top() int32                // 上坐标
	Bottom() int32             // 下坐标
	Update(tick time.Duration) // 更新
}

// 可移动的物体接口
type IMovableObject interface {
	IObject
	Dir() Direction      // 方向
	Speed() int32        // 速度
	CurrentSpeed() int32 // 当前速度
	Move(Direction)      // 移动

	// ----------------------------------
	// 事件接口
	RegisterMoveEventHandle(handle func(args ...interface{}))       // 注册移动事件
	UnregisterMoveEventHandle(handle func(args ...interface{}))     // 注销移动事件
	RegisterStopMoveEventHandle(handle func(args ...interface{}))   // 注册停止移动事件
	UnregisterStopMoveEventHandle(handle func(args ...interface{})) // 注销停止移动事件
	RegisterUpdateEventHandle(handle func(args ...interface{}))     // 注册更新事件
	UnregisterUpdateEventHandle(handle func(args ...interface{}))   // 注销更新事件
}

// 车辆接口
type IVehicle interface {
	IMovableObject
}

// 坦克接口
type ITank interface {
	IVehicle
	Change(info *ObjStaticInfo)
	Restore()

	// ---------------------------------
	// 事件接口
	RegisterChangeEventHandle(handle func(args ...interface{}))   // 注册变化事件
	UnregisterChangeEventHandle(handle func(args ...interface{})) // 注销变化事件
}
