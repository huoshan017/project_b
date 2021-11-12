package object

import "time"

// 物体接口
type IObject interface {
	Id() int32                 // 注意：这是配置id
	Type() ObjectType          // 类型
	Subtype() ObjSubType       // 子类型
	SetPos(x, y float64)       // 设置位置
	Pos() (x, y float64)       // 位置
	Width() int32              // 宽度
	Height() int32             // 高度
	Left() float64             // 左坐标
	Right() float64            // 右坐标
	Top() float64              // 上坐标
	Bottom() float64           // 下坐标
	Update(tick time.Duration) // 更新
}

// 可移动的物体接口
type IMovableObject interface {
	IObject
	Dir() Direction        // 方向
	Speed() float32        // 速度
	CurrentSpeed() float32 // 当前速度
	Move(Direction)        // 移动
	IsMove() bool          // 是否在移动
}
