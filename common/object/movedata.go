package object

import "time"

// 移动数据
type moveData struct {
	duration time.Duration
	speed    float64
	dir      Direction
}

// 移动数据freelist
type moveDataFreeList struct {
	datas []*moveData
}

func createMoveDataFreeList() *moveDataFreeList {
	return &moveDataFreeList{
		datas: make([]*moveData, 0),
	}
}

func (fl *moveDataFreeList) put(d *moveData) {
	fl.datas = append(fl.datas, d)
}

func (fl *moveDataFreeList) get() *moveData {
	l := len(fl.datas)
	if l == 0 {
		return &moveData{}
	}
	d := fl.datas[l-1]
	fl.datas = fl.datas[:l-1]
	return d
}

var mdFreeList = createMoveDataFreeList()
