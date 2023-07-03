package base

type ObjectPool[T any] struct {
	pool  []*T
	ptMap map[*T]struct{}
}

func NewObjectPool[T any]() *ObjectPool[T] {
	return &ObjectPool[T]{
		ptMap: make(map[*T]struct{}),
	}
}

func (op *ObjectPool[T]) Get() *T {
	var pt *T
	l := len(op.pool)
	if l > 0 {
		pt = op.pool[l-1]
		op.pool = op.pool[:l-1]
	} else {
		var t T
		pt = &t
	}
	op.ptMap[pt] = struct{}{}
	return pt
}

func (op *ObjectPool[T]) Put(pt *T) bool {
	if _, o := op.ptMap[pt]; o {
		return false
	}
	op.pool = append(op.pool, pt)
	delete(op.ptMap, pt)
	return true
}
