package object

import (
	"project_b/common/base"
	"project_b/common/time"
	"project_b/log"
	"unsafe"
)

// 炮彈
type Shell struct {
	MovableObject
	trackTargetId    uint32               // 跟蹤目標id
	searchTargetFunc func(*Shell) IObject // 搜索目標函數
	fetchTargetFunc  func(uint32) IObject // 獲得目標函數
	lateUpdateEvent  base.Event           // 后更新事件
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

// 立即移動
func (s *Shell) MoveNow(dir base.Angle) {
	if s.pause {
		return
	}
	s.moveDir = dir
	if s.state == stopped {
		tick := s.lastTick
		if tick == 0 {
			tick = 100 * time.Millisecond
		}
		d := GetDefaultLinearDistance(s, tick)
		v := dir.DistanceToVec2(d)
		if !s.checkMove(v.X(), v.Y(), false) {
			return
		}
		s.state = isMoving
		log.Debug("@@@ object %v stopped => moving", s.instId)
	}
}

// 更新
func (s *Shell) Update(tick time.Duration) {
	if s.pause {
		return
	}
	s.MovableObject.Update(tick)
	if s.state == isMoving {
		x, y := s.Pos()
		s.lateUpdateEvent.Call(x, y, s.WorldRotation())
	}
}

// 注冊后更新事件
func (s *Shell) RegisterLateUpdateEventHandle(handle func(args ...any)) {
	s.lateUpdateEvent.Register(handle)
}

// 注銷后更新事件
func (s *Shell) UnregisterLateUpdateEventHandle(handle func(args ...any)) {
	s.lateUpdateEvent.Unregister(handle)
}
