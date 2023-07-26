package object

import (
	"project_b/common/base"
	"project_b/common/time"
)

type IRecycle interface {
	ToRecycle()
	IsRecycle() bool
}

// 物体接口
type IObject interface {
	IRecycle
	ObjStaticInfo() *ObjStaticInfo                      // 獲得靜態信息
	Init(uint32, *ObjStaticInfo)                        // 初始化
	Uninit()                                            // 反初始化
	InstId() uint32                                     // 实例id
	Id() int32                                          // 注意：这是配置id
	Type() ObjectType                                   // 类型
	Subtype() ObjSubtype                                // 子类型
	OwnerType() ObjOwnerType                            // 所有者类型
	StaticInfo() *ObjStaticInfo                         // 靜態信息
	Center() (x, y int32)                               // 中心點坐標，本地坐標係
	Pos() (x, y int32)                                  // 中心位置，世界坐標係
	SetPos(x, y int32)                                  // 设置中心位置，世界坐標係
	Width() int32                                       // 宽度
	Length() int32                                      // 長度
	LeftBottom() (int32, int32)                         // 左下坐标，相對於本地坐標系
	LeftTop() (int32, int32)                            // 左上坐標，相對於本地坐標系
	RightTop() (int32, int32)                           // 右上坐标，相對於本地坐標系
	RightBottom() (int32, int32)                        // 右下坐標，相對於本地坐標系
	WorldRotation() base.Angle                          // todo 是局部坐標系在世界坐標系中的旋轉
	LocalRotation() base.Angle                          // 局部坐標系中的旋轉
	Rotation() base.Angle                               // 最終的旋轉(x軸正向逆時針旋轉角度)，局部旋轉與世界旋轉纍加，垂直於寬(Width)，平行於長(Height)
	OriginalLeft() int32                                // 原始左坐標
	OriginalRight() int32                               // 原始右坐標
	OriginalTop() int32                                 // 原始上坐標
	OriginalBottom() int32                              // 原始下坐標
	Update(tick time.Duration)                          // 更新
	Camp() CampType                                     // 陣營
	SetCamp(CampType)                                   // 設置陣營
	RestoreCamp()                                       // 重置陣營
	AddComp(comp IComponent)                            // 添加組件
	RemoveComp(name string)                             // 去除組件
	GetComp(name string) IComponent                     // 獲取組件
	HasComp(name string) bool                           // 是否擁有組件
	RegisterDestroyedEventHandle(handle func(...any))   // 注冊銷毀事件函數
	UnregisterDestroyedEventHandle(handle func(...any)) // 注銷銷毀事件函數
}

// 静态物体接口
type IStaticObject interface {
	IObject
}

// 可移动的物体接口
type IMovableObject interface {
	IObject
	MovableObjStaticInfo() *MovableObjStaticInfo
	Level() int32              // 等级
	Speed() int32              // 速度
	MoveDir() base.Angle       // 移動方向
	CurrentSpeed() int32       // 当前速度
	Rotate(angle base.Angle)   // 旋轉，逆時針為正方向 [0, 360)
	RotateTo(angle base.Angle) // 逆時針旋轉到 [0, 360)
	Forward() base.Vec2        // 朝向向量
	Move(dir base.Angle)       // 移动
	Stop()                     // 停止
	IsMoving() bool            // 是否在移动
	LastPos() (x, y int32)     // 上次Update時位置

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

// 環繞物
type ISurroundObject interface {
	IMovableObject
	SurroundObjStaticInfo() *SurroundObjStaticInfo                                 // 靜態配置
	SetAroundCenterObject(centerObjInstId uint32, getObjFunc func(uint32) IObject) // 設置環繞中心物體
	GetAroundCenterObject() IObject                                                // 獲得中心點

	// 事件接口
	RegisterLateUpdateEventHandle(handle func(args ...any))   // 注冊后更新事件
	UnregisterLateUpdateEventHandle(handle func(args ...any)) // 注銷后更新事件
}

// 车辆接口
type IVehicle interface {
	IMovableObject
}

// 坦克接口
type ITank interface {
	IVehicle
	TankStaticInfo() *TankStaticInfo
	Change(info *TankStaticInfo)
	Restore()

	// ---------------------------------
	// 事件接口
	RegisterChangeEventHandle(handle func(args ...any))   // 注册变化事件
	UnregisterChangeEventHandle(handle func(args ...any)) // 注销变化事件
}

// 效果接口
type IEffect interface {
	InstId() uint32
	StaticInfo() *EffectStaticInfo
	SetPos(int32, int32)
	Pos() (int32, int32)
	Width() int32
	Height() int32
	Update()
	IsOver() bool
}
