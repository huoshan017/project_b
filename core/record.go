package core

import (
	"fmt"
	"project_b/common/time"
	"project_b/game_map"
)

type recordData struct {
	name      string
	mapConfig *game_map.Config
	frameList []*frameData
}

type Replay struct {
	mapConfig *game_map.Config
	frameList []*frameData
}

func (r *Replay) Clear() {
	r.mapConfig = nil
	r.frameList = nil
}

type ReplayManager struct {
	inst       *Instance
	replayList []recordData
	sel        int32
}

func NewReplayManager(inst *Instance) *ReplayManager {
	return &ReplayManager{
		inst: inst,
		sel:  -1,
	}
}

func (rm *ReplayManager) SetRecord() {
	rm.inst.setRecordHandle(rm.Save)
}

func (rm *ReplayManager) Save(mapConfig *game_map.Config, frameList []*frameData) {
	saveName := fmt.Sprintf("%v: %v", mapConfig.Name, time.Now().GoString())
	rm.replayList = append(rm.replayList, recordData{name: saveName, mapConfig: mapConfig, frameList: frameList})
}

func (rm *ReplayManager) Delete(index int32) bool {
	if int(index) >= len(rm.inst.frameList) {
		return false
	}
	rm.replayList = append(rm.replayList[:index], rm.replayList[index+1:]...)
	return true
}

func (rm *ReplayManager) Select(index int32) {
	if int(index) <= len(rm.replayList) {
		rm.sel = index
	}
}

func (rm *ReplayManager) SelectedReplay() Replay {
	if rm.sel < 0 {
		panic("not selected replay")
	}
	replay := rm.replayList[rm.sel]
	return Replay{mapConfig: replay.mapConfig, frameList: replay.frameList}
}

func (rm *ReplayManager) RecordNameList() []string {
	var nameList []string
	for i := 0; i < len(rm.replayList); i++ {
		nameList = append(nameList, rm.replayList[i].name)
	}
	return nameList
}

func (rm *ReplayManager) GetRecordCount() int32 {
	return int32(len(rm.replayList))
}

func (rm *ReplayManager) GetRecordName(index int32) string {
	return rm.replayList[index].name
}
