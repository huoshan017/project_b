package object

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
type ColliderComp struct {
	collisionHandle func(...any)
}

// 組件名稱
func (c ColliderComp) Name() string {
	return "Collider"
}

// 獲得AABB
func (c ColliderComp) GetAABB(obj IObject) AABB {
	return AABB{
		Left:   obj.Left(),
		Bottom: obj.Bottom(),
		Right:  obj.Right(),
		Top:    obj.Top(),
	}
}

// 注冊碰撞事件處理
func (c *ColliderComp) SetCollisionHandle(handle func(...any)) {
	c.collisionHandle = handle
}

// 執行
func (c *ColliderComp) CallCollisionEventHandle(args ...any) {
	if c.collisionHandle != nil {
		c.collisionHandle(args...)
	}
}
