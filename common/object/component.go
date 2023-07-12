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
func (ab *AABB) MoveIntersect(moveDir Direction, aabb *AABB) bool {
	switch moveDir {
	case DirLeft:
		return !(ab.Left-1 >= aabb.Right || ab.Right-1 <= aabb.Left || ab.Top <= aabb.Bottom || ab.Bottom >= aabb.Top)
	case DirRight:
		return !(ab.Left+1 >= aabb.Right || ab.Right+1 <= aabb.Left || ab.Top <= aabb.Bottom || ab.Bottom >= aabb.Top)
	case DirUp:
		return !(ab.Left >= aabb.Right || ab.Right <= aabb.Left || ab.Top+1 <= aabb.Bottom || ab.Bottom+1 >= aabb.Top)
	case DirDown:
		return !(ab.Left >= aabb.Right || ab.Right <= aabb.Left || ab.Top-1 <= aabb.Bottom || ab.Bottom-1 >= aabb.Top)
	default:
		return false
	}
}

// 移動
func (ab *AABB) Move(dir Direction, dx, dy int32) {
	switch dir {
	case DirLeft:
		ab.Left += dx
		ab.Right += dx
	case DirRight:
		ab.Left += dx
		ab.Right += dx
	case DirUp:
		ab.Top += dy
		ab.Bottom += dy
	case DirDown:
		ab.Top += dy
		ab.Bottom += dy
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
	//left, bottom := obj.LeftBottom()
	//right, top := obj.RightTop()
	return AABB{
		Left:   obj.OriginalLeft(),
		Bottom: obj.OriginalBottom(),
		Right:  obj.OriginalRight(),
		Top:    obj.OriginalTop(),
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
