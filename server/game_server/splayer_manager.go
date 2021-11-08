package main

import (
	"project_b/common"
	"strconv"
	"sync/atomic"

	cmap "github.com/orcaman/concurrent-map"
)

// key检测管理器
type KeyCheckManager struct {
	account2State cmap.ConcurrentMap
}

func NewKeyCheckManager() *KeyCheckManager {
	return &KeyCheckManager{
		account2State: cmap.New(),
	}
}

func (m *KeyCheckManager) checkAndAdd(account string) bool {
	return m.account2State.SetIfAbsent(account, 1)
}

func (m *KeyCheckManager) remove(account string) {
	m.account2State.Remove(account)
}

// 玩家管理器结构，主要用在同一玩家的查询
type SPlayerManager struct {
	idCounter  uint64
	id2Player  cmap.ConcurrentMap
	acc2Player cmap.ConcurrentMap
}

func NewSPlayerManager() *SPlayerManager {
	return &SPlayerManager{
		id2Player:  cmap.New(),
		acc2Player: cmap.New(),
	}
}

func (pm *SPlayerManager) Add(p *SPlayer) {
	// 没有设置id
	if p.Id() == 0 {
		pid := atomic.AddUint64(&pm.idCounter, 1)
		p.SetId(pid)
	}
	pm.id2Player.Set(strconv.FormatUint(p.Id(), 10), p)
	pm.acc2Player.Set(p.Account(), p)
}

func (pm *SPlayerManager) Remove(pid uint64) {
	pm.id2Player.RemoveCb(strconv.FormatUint(pid, 10), func(_ string, value interface{}, _ bool) bool {
		p, o := value.(common.IPlayer)
		if !o {
			return false
		}
		pm.acc2Player.Remove(p.Account())
		return true
	})
}

func (pm *SPlayerManager) Get(pid uint64) *SPlayer {
	p, o := pm.id2Player.Get(strconv.FormatUint(pid, 10))
	if !o {
		return nil
	}
	return p.(*SPlayer)
}

func (pm *SPlayerManager) GetByAccount(account string) *SPlayer {
	p, o := pm.acc2Player.Get(account)
	if !o {
		return nil
	}
	return p.(*SPlayer)
}

func (pm *SPlayerManager) GetNextId() uint64 {
	return atomic.AddUint64(&pm.idCounter, 1)
}

func (pm *SPlayerManager) PeekNextId() uint64 {
	return atomic.LoadUint64(&pm.idCounter) + 1
}
