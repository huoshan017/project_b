package object

import "project_b/common/base"

type IComponent interface {
	Name() string
}

type AABB struct {
	Left, Bottom int32
	Right, Top   int32
}

// 是否相交
func (ab *AABB) Intersect(aabb *AABB) bool {
	return !(ab.Left >= aabb.Right || ab.Right <= aabb.Left || ab.Top <= aabb.Bottom || ab.Bottom >= aabb.Top)
}

// 移動
func (ab *AABB) Move(dir Direction, distance float64) {
	switch dir {
	case DirLeft:
		ab.Left -= int32(distance)
		ab.Right -= int32(distance)
	case DirRight:
		ab.Left += int32(distance)
		ab.Right += int32(distance)
	case DirUp:
		ab.Top += int32(distance)
		ab.Bottom += int32(distance)
	case DirDown:
		ab.Top -= int32(distance)
		ab.Bottom -= int32(distance)
	}
}

// 碰撞組件
type CollisionComp struct {
	collisionEvent base.Event
}

// 組件名稱
func (c CollisionComp) Name() string {
	return "Collision"
}

// 獲得AABB
func (c CollisionComp) GetAABB(obj IObject) AABB {
	return AABB{
		Left:   obj.Left(),
		Bottom: obj.Bottom(),
		Right:  obj.Right(),
		Top:    obj.Top(),
	}
}

// 注冊碰撞事件處理
func (c *CollisionComp) RegisterCollisionEventHandle(handle func(...any)) {
	c.collisionEvent.Register(handle)
}

// 注銷碰撞事件處理
func (c *CollisionComp) UnregisterCollisionEventHandle(handle func(...any)) {
	c.collisionEvent.Unregister(handle)
}

// 執行
func (c *CollisionComp) Call(args ...any) {
	c.collisionEvent.Call(args...)
}
