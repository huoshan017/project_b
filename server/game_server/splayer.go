package main

import (
	"project_b/common"
	"project_b/common_data"
	"sync/atomic"
	"unsafe"

	"github.com/huoshan017/gsnet"
)

// 服务器玩家结构
type SPlayer struct {
	common.SPlayer
	currChangedTankIdIndex int32 // 当前坦克索引
	kicker                 *SamePlayerKickHandler
}

// 创建SPlayer
func NewSPlayer(acc string, id uint64, sess gsnet.ISession) *SPlayer {
	p := &SPlayer{
		SPlayer:                *common.NewSPlayer(acc, id, sess),
		currChangedTankIdIndex: -1, // 初始化成-1，因为每次使用时肯定会加上1，所以可以保证为非负数
	}
	return p
}

// 重置会话
func (p *SPlayer) ResetSess(sess gsnet.ISession) {
	p.SPlayer.ResetSess(sess)
}

// 设置会话数据
func (p *SPlayer) SetSessData(k string, d interface{}) {
	p.SPlayer.SetSessData(k, d)
}

// 获得会话数据
func (p *SPlayer) GetSessData(k string) interface{} {
	return p.SPlayer.GetSessData(k)
}

// 获得坦克改变的id
func (p *SPlayer) GetChangeTankId() int32 {
	tank := p.GetTank()
	oid := tank.OriginId()
	var changeId int32
	for {
		p.currChangedTankIdIndex += 1
		if int(p.currChangedTankIdIndex) >= len(common_data.TankIdList) {
			p.currChangedTankIdIndex = 0
		}
		changeId = common_data.TankIdList[p.currChangedTankIdIndex]
		if changeId != oid {
			break
		}
	}
	return changeId
}

// 断开并等待结束通知
func (p *SPlayer) WaitDisconnect() {
	ptr := (*unsafe.Pointer)(unsafe.Pointer(&p.kicker))
	if atomic.CompareAndSwapPointer(ptr, nil, unsafe.Pointer(NewSamePlayerKickHandler(p))) {
		p.kicker.NotfyDisconnectAndWait()
	}
}

// 获得踢人处理器
func (p *SPlayer) GetKicker() *SamePlayerKickHandler {
	ptr := (*unsafe.Pointer)(unsafe.Pointer(&p.kicker))
	return (*SamePlayerKickHandler)(atomic.LoadPointer(ptr))
}

// 同玩家踢人处理器
type SamePlayerKickHandler struct {
	p         *SPlayer
	discCh    chan struct{} // 通知断连通道
	discEndCh chan struct{} // 断连结束通道
	closed    int32
}

// 创建同玩家踢人处理器
func NewSamePlayerKickHandler(p *SPlayer) *SamePlayerKickHandler {
	return &SamePlayerKickHandler{
		p:         p,
		discCh:    make(chan struct{}, 1),
		discEndCh: make(chan struct{}, 1),
	}
}

// 通知断连并等待
func (h *SamePlayerKickHandler) NotfyDisconnectAndWait() {
	if atomic.LoadInt32(&h.closed) > 0 {
		return
	}
	// 先通知断连
	h.discCh <- struct{}{}
	// 等待断连结束
	<-h.discEndCh
	// 关闭断连channel
	atomic.StoreInt32(&h.closed, 1)
	close(h.discCh)
}

// 检测断连通知
func (h *SamePlayerKickHandler) CheckDisconnectNotification() {
	if atomic.LoadInt32(&h.closed) > 0 {
		return
	}
	select {
	case <-h.discCh:
		if h.p.Player.IsLeft() {
			h.p.Player.Left(true)
		}
		h.p.SPlayer.Disconnect()
	default:
	}
}

// 通知断连结束, 这个必须在ISessionHandler.OnDisconnect中调用
func (h *SamePlayerKickHandler) NotifyDisconnectEnd() {
	if atomic.LoadInt32(&h.closed) > 0 {
		return
	}
	select {
	case h.discEndCh <- struct{}{}:
	default:
	}
	// 关闭通知断连结束channel
	close(h.discEndCh)
}
