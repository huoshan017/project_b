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

// 移動相交
func (ab *AABB) MoveIntersect(dx, dy int32, aabb *AABB) bool {
	return !(ab.Left+dx >= aabb.Right || ab.Right+dx <= aabb.Left || ab.Top+dy <= aabb.Bottom || ab.Bottom+dy >= aabb.Top)
}

// 碰撞組件
type ColliderComp struct {
	collisionHandle func(IMovableObject, *CollisionInfo)
	obj             IObject
}

// 組件名稱
func (c ColliderComp) Name() string {
	return "Collider"
}

// 獲得AABB
func (c ColliderComp) GetAABB() AABB {
	var (
		s    int32
		w    = c.obj.Width()
		l    = c.obj.Length()
		x, y = c.obj.Pos()
	)
	if w < l {
		s = w
	} else {
		s = l
	}
	return AABB{
		Left:   x - s/2,
		Bottom: y - s/2,
		Right:  x + s/2,
		Top:    y + s/2,
	}
}

// 注冊碰撞事件處理
func (c *ColliderComp) SetCollisionHandle(handle func(IMovableObject, *CollisionInfo)) {
	c.collisionHandle = handle
}

// 執行
func (c *ColliderComp) CallCollisionEventHandle(mobj IMovableObject, ci *CollisionInfo) {
	if c.collisionHandle != nil {
		c.collisionHandle(mobj, ci)
	}
}
