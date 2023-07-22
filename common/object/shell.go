package object

import "unsafe"

// 炮彈
type Shell struct {
	MovableObject
	trackTargetId    uint32
	searchTargetFunc func(*Shell) IObject
	fetchTargetFunc  func(uint32) IObject
}

// 创建炮彈
func NewShell(instId uint32, staticInfo *ShellStaticInfo) *Shell {
	o := &Shell{}
	o.Init(instId, &staticInfo.ObjStaticInfo)
	return o
}

// 靜態配置
func (b *Shell) ShellStaticInfo() *ShellStaticInfo {
	return (*ShellStaticInfo)(unsafe.Pointer(b.staticInfo))
}

// 初始化
func (b *Shell) Init(instId uint32, staticInfo *ObjStaticInfo) {
	b.MovableObject.Init(instId, staticInfo)
	b.setSuper(b)
}

// 反初始化
func (b *Shell) Uninit() {
	b.trackTargetId = 0
	b.searchTargetFunc = nil
	b.fetchTargetFunc = nil
	b.MovableObject.Uninit()
}

// 設置搜索目標函數
func (b *Shell) SetSearchTargetFunc(f func(*Shell) IObject) {
	b.searchTargetFunc = f
}

// 設置獲取目標函數
func (b *Shell) SetFetchTargetFunc(f func(uint32) IObject) {
	b.fetchTargetFunc = f
}
