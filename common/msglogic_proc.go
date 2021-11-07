package common

import (
	"errors"
	"sync"
	"sync/atomic"
)

const (
	GameLogicProcDefaultMsgListLength = 100
)

var (
	ErrMsgLogicProcNoFoundHandle = errors.New("common: not found handle for message logic proc")
)

type AgentKey interface{}
type MsgData interface{}

// 消息结构
type msgData struct {
	key   AgentKey
	msgid uint32
	msg   interface{}
}

// 代理数据
type agentData struct {
	typ      int32 // 1. 添加  2. 删除  3. 更新
	key      AgentKey
	data     interface{}
	onHandle func(interface{}) error
}

// 消息逻辑处理器，单线程结构，非线程安全
type MsgLogicProc struct {
	msgList     chan *msgData
	msgDataPool *sync.Pool
	agentCh     chan *agentData
	agentMap    map[AgentKey]interface{}
	handleMap   map[uint32]func(AgentKey, MsgData) error
	closeCh     chan struct{}
	closed      int32
	errHandle   func(err error)
}

// 创建消息逻辑处理
func CreateMsgLogicProc() *MsgLogicProc {
	return &MsgLogicProc{
		msgList: make(chan *msgData, GameLogicProcDefaultMsgListLength),
		msgDataPool: &sync.Pool{
			New: func() interface{} {
				return &msgData{}
			},
		},
		agentCh:   make(chan *agentData),
		agentMap:  make(map[AgentKey]interface{}),
		closeCh:   make(chan struct{}),
		handleMap: make(map[uint32]func(AgentKey, MsgData) error),
	}
}

// 注册消息处理器
func (p *MsgLogicProc) RegisterHandle(msgid uint32, handle func(key AgentKey, msg MsgData) error) {
	if atomic.LoadInt32(&p.closed) > 0 {
		return
	}
	p.handleMap[msgid] = handle
}

// 设置错误处理函数
func (p *MsgLogicProc) SetErrorHandle(handle func(err error)) {
	p.errHandle = handle
}

// 压入代理
func (t *MsgLogicProc) PushAgent(key AgentKey, data interface{}, handle func(interface{}) error) {
	if atomic.LoadInt32(&t.closed) > 0 {
		return
	}
	select {
	case <-t.closeCh:
		close(t.agentCh)
		atomic.StoreInt32(&t.closed, 1)
	case t.agentCh <- &agentData{typ: 1, key: key, data: data, onHandle: handle}:
	}
}

// 删除代理
func (t *MsgLogicProc) DeleteAgent(key AgentKey, data interface{}, handle func(interface{}) error) {
	if atomic.LoadInt32(&t.closed) > 0 {
		return
	}
	select {
	case <-t.closeCh:
		atomic.StoreInt32(&t.closed, 1)
	case t.agentCh <- &agentData{typ: 2, key: key, data: data, onHandle: handle}:
	}
}

// 更新代理
func (t *MsgLogicProc) UpdateAgent(key AgentKey, data interface{}, handle func(interface{}) error) {
	if atomic.LoadInt32(&t.closed) > 0 {
		return
	}
	select {
	case <-t.closeCh:
		atomic.StoreInt32(&t.closed, 1)
	case t.agentCh <- &agentData{typ: 3, key: key, data: data, onHandle: handle}:
	}
}

// 获得代理数量
func (t *MsgLogicProc) GetAgentCountNoLock() int32 {
	return int32(len(t.agentMap))
}

// 获得代理map
func (t *MsgLogicProc) GetAgentMapNoLock() map[AgentKey]interface{} {
	return t.agentMap
}

// 无锁获得代理
func (t *MsgLogicProc) GetAgentNoLock(key AgentKey) interface{} {
	return t.agentMap[key]
}

// 压入玩家消息数据
func (t *MsgLogicProc) PushMsg(key AgentKey, msgid uint32, msg interface{}) {
	if atomic.LoadInt32(&t.closed) > 0 {
		return
	}
	select {
	case <-t.closeCh:
		close(t.msgList)
		atomic.StoreInt32(&t.closed, 1)
	default:
		m := t.msgDataPool.Get().(*msgData)
		m.key = key
		m.msgid = msgid
		m.msg = msg
		t.msgList <- m
	}
}

// 运行
func (p *MsgLogicProc) Run() {
	if atomic.LoadInt32(&p.closed) > 0 {
		return
	}
	loop := true
	for loop {
		select {
		case d, o := <-p.msgList:
			if !o {
				loop = false
			} else {
				// 处理玩家消息
				err := p.handlePlayerMsg(d.key, d.msgid, d.msg)
				if err != nil && p.errHandle != nil {
					p.errHandle(err)
				}
				// 回收消息结构
				p.msgDataPool.Put(d)
			}
		case a, o := <-p.agentCh:
			if !o {
				loop = false
			} else {
				var err error
				if a.typ == 1 {
					// 添加代理
					p.agentMap[a.key] = a.data
					err = a.onHandle(a.data)
				} else if a.typ == 2 {
					// 先处理回调，再删除代理
					err = a.onHandle(a.key)
					delete(p.agentMap, a.key)
				} else if a.typ == 3 {
					// 更新代理
					p.agentMap[a.key] = a.data
					err = a.onHandle(a.data)
				}
				if err != nil && p.errHandle != nil {
					p.errHandle(err)
				}
			}
		case <-p.closeCh:
			loop = false
		}
	}
	atomic.CompareAndSwapInt32(&p.closed, 0, 1)
}

// 关闭
func (t *MsgLogicProc) Close() {
	close(t.closeCh)
}

// 处理玩家消息
func (t *MsgLogicProc) handlePlayerMsg(key AgentKey, msgid uint32, msg MsgData) error {
	h, o := t.handleMap[msgid]
	if !o {
		return ErrMsgLogicProcNoFoundHandle
	}
	return h(key, msg)
}
