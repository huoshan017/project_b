package object

import (
	"project_b/common/base"
	"project_b/common/time"
	"project_b/log"
	"unsafe"
)

// 炮彈
type Shell struct {
	*MovableObject
	trackTargetId    uint32               // 跟蹤目標id
	searchTargetFunc func(*Shell) IObject // 搜索目標函數
	fetchTargetFunc  func(uint32) IObject // 獲得目標函數
	lateUpdateEvent  base.Event           // 后更新事件
}

// 创建炮彈
func NewShell() *Shell {
	s := &Shell{}
	s.MovableObject = NewMovableObjectWithSuper(s)
	return s
}

// 通過子類創建炮彈
func NewShellWithSuper(super IMovableObject) *Shell {
	s := &Shell{}
	s.MovableObject = NewMovableObjectWithSuper(super)
	return s
}

// 靜態配置
func (s *Shell) ShellStaticInfo() *ShellStaticInfo {
	return (*ShellStaticInfo)(unsafe.Pointer(s.staticInfo))
}

// 初始化
func (s *Shell) Init(instId uint32, staticInfo *ObjStaticInfo) {
	s.MovableObject.Init(instId, staticInfo)
}

// 反初始化
func (s *Shell) Uninit() {
	s.trackTargetId = 0
	s.searchTargetFunc = nil
	s.fetchTargetFunc = nil
	s.MovableObject.Uninit()
}

// 設置搜索目標函數
func (s *Shell) SetSearchTargetFunc(f func(*Shell) IObject) {
	if s.ShellStaticInfo().TrackTarget {
		s.searchTargetFunc = f
	}
}

// 設置獲取目標函數
func (s *Shell) SetFetchTargetFunc(f func(uint32) IObject) {
	if s.ShellStaticInfo().TrackTarget {
		s.fetchTargetFunc = f
	}
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
		if !s.checkMove(v.X(), v.Y(), false, nil, nil) {
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
