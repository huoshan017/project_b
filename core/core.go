package core

import (
	"project_b/common/base"
	"project_b/game_map"
)

type coreState int

const (
	coreStateIdle         coreState = iota
	coreStatePrepareStart coreState = 1
	coreStateRunning      coreState = 2
)

type CoreArgs struct {
	EventMgr   base.IEventManager
	PlayerNum  int32
	FrameMs    uint32
	GetMapFunc func(int32) *game_map.Config
}

type GameCore struct {
	args        CoreArgs
	inst        *Instance
	recordMgr   *RecordManager
	lastCheckMs uint32
	state       coreState
	isRecord    bool
}

func NewGameCore(args CoreArgs) *GameCore {
	inst := newInstance(args.EventMgr, args.PlayerNum, args.FrameMs)
	return &GameCore{
		args:      args,
		inst:      inst,
		recordMgr: newRecordManager(inst),
	}
}

func (core *GameCore) GetRecordMgr() *RecordManager {
	return core.recordMgr
}

func (core *GameCore) LoadMap(mapId int32) bool {
	return core.inst.Load(core.args.GetMapFunc(mapId))
}

func (core *GameCore) LoadRecordByIndex(index int32) bool {
	core.recordMgr.Select(index)
	record := core.recordMgr.SelectedRecord()
	if record == nil {
		return false
	}
	core.isRecord = false
	return core.inst.LoadRecord(core.args.GetMapFunc(record.mapId), record)
}

func (core *GameCore) RegisterEvent(eid base.EventId, handle func(...any)) {
	core.inst.RegisterEvent(eid, handle)
}

func (core *GameCore) UnregisterEvent(eid base.EventId, handle func(...any)) {
	core.inst.UnregisterEvent(eid, handle)
}

func (core *GameCore) RegisterPlayerEvent(playerId uint64, eid base.EventId, handle func(...any)) {
	core.inst.RegisterPlayerEvent(playerId, eid, handle)
}

func (core *GameCore) UnregisterPlayerEvent(playerId uint64, eid base.EventId, handle func(...any)) {
	core.inst.UnregisterPlayerEvent(playerId, eid, handle)
}

func (core *GameCore) Start(playerIdList []uint64, isRecord bool) bool {
	if !core.inst.CheckAndStart(playerIdList) {
		return false
	}
	core.state = coreStatePrepareStart
	core.isRecord = isRecord
	return true
}

func (core *GameCore) LoadRecordStart(index int32) bool {
	if !core.LoadRecordByIndex(index) {
		return false
	}
	return core.Start(nil, false)
}

func (core *GameCore) Restart() bool {
	core.checkRecord()
	if !core.inst.Restart() {
		return false
	}
	core.state = coreStatePrepareStart
	return true
}

func (core *GameCore) Pause() {
	core.inst.Pause()
}

func (core *GameCore) Resume() {
	core.inst.Resume()
}

func (core *GameCore) PushFramePlayerCmd(frameNum uint32, playerId uint64, cmdData *CmdData) bool {
	return core.inst.PushFrame(frameNum, playerId, cmdData.cmd, cmdData.args)
}

func (core *GameCore) PushSyncPlayerCmd(playerId uint64, cmdData *CmdData) bool {
	return core.inst.PushFrame(core.inst.GetFrame(), playerId, cmdData.cmd, cmdData.args)
}

func (core *GameCore) Update(ms uint32) {
	if core.state == coreStatePrepareStart {
		core.lastCheckMs = ms
		core.state = coreStateRunning
	}
	if core.state != coreStateRunning {
		return
	}
	usedMs := ms - core.lastCheckMs
	for usedMs >= uint32(core.args.FrameMs) {
		core.inst.UpdateFrame()
		usedMs -= core.args.FrameMs
		core.lastCheckMs += core.args.FrameMs
	}
}

func (core *GameCore) End() {
	core.checkRecord()
	core.inst.Unload()
	core.state = coreStateIdle
}

func (core *GameCore) GetFrame() uint32 {
	return core.inst.logic.GetCurrFrame()
}

func (core *GameCore) UsedMs() uint32 {
	return core.inst.logic.GetCurrFrame() * core.args.FrameMs
}

func (core *GameCore) checkRecord() {
	if !core.isRecord {
		return
	}
	mapConfig := core.inst.logic.World().GetMapConfig()
	core.recordMgr.Save(mapConfig.Name, Record{
		mapId: mapConfig.Id, frameList: core.inst.frameList, playerIdList: core.inst.playerIdList, frameNum: core.inst.GetFrame(), frameMs: core.inst.frameMs,
	})
}
