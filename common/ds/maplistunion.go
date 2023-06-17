package ds

type Pair[Key comparable, Value any] struct {
	Key   Key
	Value Value
}

// map和list组合体
type MapListUnion[Key comparable, Value any] struct {
	key2index map[Key]int32      // 关键字到索引
	list      []Pair[Key, Value] // 保存kv对
}

func NewMapListUnion[Key comparable, Value any]() *MapListUnion[Key, Value] {
	return &MapListUnion[Key, Value]{
		list:      make([]Pair[Key, Value], 0),
		key2index: make(map[Key]int32),
	}
}

func (l *MapListUnion[Key, Value]) Count() int32 {
	return int32(len(l.list))
}

func (l *MapListUnion[Key, Value]) Add(k Key, v Value) {
	if l.Exists(k) {
		return
	}
	// 把值追加到队列尾部
	l.list = append(l.list, Pair[Key, Value]{Key: k, Value: v})
	// 建立key跟索引的映射
	l.key2index[k] = int32(len(l.list) - 1)
}

func (l *MapListUnion[Key, Value]) Exists(k Key) bool {
	_, o := l.key2index[k]
	return o
}

func (l *MapListUnion[Key, Value]) Get(k Key) (Value, bool) {
	// 先获得索引
	idx, o := l.key2index[k]
	if !o {
		var v Value
		return v, false
	}
	// 再取值
	return l.list[idx].Value, true
}

func (l *MapListUnion[Key, Value]) Remove(k Key) Value {
	idx, o := l.key2index[k]
	if !o {
		var v Value
		return v
	}
	// 需要删除的值
	result := l.list[idx]
	// 最后一个索引
	le := int32(len(l.list) - 1)
	// 被删除的不是队列的最后一个
	if idx != le {
		// 挪到需要删除的那个值的位置上
		l.list[idx] = l.list[le]
		// 更新索引
		l.key2index[l.list[idx].Key] = idx
	}
	// 更新列表
	l.list = l.list[:le]
	// 删除key与索引的映射
	delete(l.key2index, k)
	return result.Value
}

func (l *MapListUnion[Key, Value]) GetByIndex(idx int32) (key Key, value Value) {
	kv := l.list[idx]
	return kv.Key, kv.Value
}

func (l *MapListUnion[Key, Value]) GetValues() []Value {
	var lis []Value
	for _, kv := range l.list {
		lis = append(lis, kv.Value)
	}
	return lis
}

func (l *MapListUnion[Key, Value]) GetKeys(keys []Key) []Key {
	for _, kv := range l.list {
		keys = append(keys, kv.Key)
	}
	return keys
}

func (l *MapListUnion[Key, Value]) GetList() []Pair[Key, Value] {
	return l.list
}
