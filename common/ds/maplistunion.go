package ds

type kvPair struct {
	key   interface{}
	value interface{}
}

// map和list组合体
type MapListUnion struct {
	key2index map[interface{}]int32 // 关键字到索引
	list      []kvPair              // 保存kv对
}

func NewMapListUnion() *MapListUnion {
	return &MapListUnion{
		list:      make([]kvPair, 0),
		key2index: make(map[interface{}]int32),
	}
}

func (l *MapListUnion) Count() int32 {
	return int32(len(l.list))
}

func (l *MapListUnion) Add(k interface{}, v interface{}) {
	if l.Exists(k) {
		return
	}
	// 把值追加到队列尾部
	l.list = append(l.list, kvPair{key: k, value: v})
	// 建立key跟索引的映射
	l.key2index[k] = int32(len(l.list) - 1)
}

func (l *MapListUnion) Exists(k interface{}) bool {
	_, o := l.key2index[k]
	return o
}

func (l *MapListUnion) Get(k interface{}) (interface{}, bool) {
	// 先获得索引
	idx, o := l.key2index[k]
	if !o {
		return nil, false
	}
	// 再取值
	return l.list[idx].value, true
}

func (l *MapListUnion) Remove(k interface{}) interface{} {
	idx, o := l.key2index[k]
	if !o {
		return nil
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
		l.key2index[l.list[idx].key] = idx
	}
	// 更新列表
	l.list = l.list[:le]
	// 删除key与索引的映射
	delete(l.key2index, k)
	return result.value
}

func (l *MapListUnion) GetByIndex(idx int32) (key interface{}, value interface{}) {
	kv := l.list[idx]
	return kv.key, kv.value
}

func (l *MapListUnion) GetList() []interface{} {
	var lis []interface{}
	for _, kv := range l.list {
		lis = append(lis, kv.value)
	}
	return lis
}
