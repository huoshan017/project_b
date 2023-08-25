package core

import (
	"fmt"
	"project_b/common"
	"project_b/common/base"
	"project_b/common/object"
	"project_b/common/time"
	"project_b/game_map"
	"project_b/log"

	"github.com/huoshan017/ponu/list"
)

type instanceMode int

const (
	instanceModePlay   instanceMode = iota
	instanceModeReplay instanceMode = 1
)

type InstanceArgs struct {
	EventMgr   base.IEventManager
	PlayerNum  int32
	UpdateTick time.Duration
	SavePath   string
}

type playerData struct {
	playerId     uint64
	frameCmdList []struct {
		frame uint32
		cmd   CmdData
	}
}

type frameData struct {
	frameNum       uint32
	playerDataList []*playerData
}

func (fd *frameData) clear() {
	fd.frameNum = 0
	for _, pd := range fd.playerDataList {
		if len(pd.frameCmdList) == 0 {
			continue
		}
		clear(pd.frameCmdList)
		pd.frameCmdList = pd.frameCmdList[:0]
	}
}

type Instance struct {
	args              *InstanceArgs                        // 參數
	mode              instanceMode                         // 模式
	replay            Replay                               // 重播數據，只有在重播模式下才有用
	frameList         []*frameData                         // 幀列表
	playerIdList      []uint64                             // 玩家列表
	logic             *common.GameLogic                    // 游戲邏輯
	frameIndexInList  uint32                               // 幀列表frameList或者replay.frameList的當前索引
	frameDataFreeList *list.ListT[*frameData]              // 幀數據freelist
	recordHandle      func(*game_map.Config, []*frameData) // 錄像處理
}

func NewInstance(args *InstanceArgs) *Instance {
	return &Instance{
		args:              args,
		logic:             common.NewGameLogic(args.EventMgr),
		frameDataFreeList: list.NewListT[*frameData](),
	}
}

func (inst *Instance) Load(config *game_map.Config) bool {
	res := inst.logic.LoadScene(config)
	if res {
		inst.mode = instanceModePlay
	}
	return res
}

func (inst *Instance) Unload() {
	mapConfig := inst.logic.World().GetMapConfig()
	inst.logic.UnloadScene()
	if inst.recordHandle != nil {
		inst.recordHandle(mapConfig, inst.frameList)
	} else {
		inst.recycleFrameList()
	}
	inst.frameIndexInList = 0
}

func (inst *Instance) LoadReplay(replay Replay) bool {
	res := inst.logic.LoadScene(replay.mapConfig)
	if res {
		inst.mode = instanceModeReplay
		inst.replay = replay
	}
	return res
}

func (inst *Instance) RegisterEvent(eid base.EventId, handle func(...any)) {
	inst.logic.RegisterEvent(eid, handle)
}

func (inst *Instance) UnregisterEvent(eid base.EventId, handle func(...any)) {
	inst.logic.UnregisterEvent(eid, handle)
}

func (inst *Instance) RegisterPlayerEvent(playerId uint64, eid base.EventId, handle func(...any)) {
	inst.logic.RegisterPlayerEvent(playerId, eid, handle)
}

func (inst *Instance) UnregisterPlayerEvent(playerId uint64, eid base.EventId, handle func(...any)) {
	inst.logic.UnregisterPlayerEvent(playerId, eid, handle)
}

func (inst *Instance) CheckAndStart(playerList []uint64) bool {
	if len(playerList) != int(inst.args.PlayerNum) {
		return false
	}
	for _, pid := range playerList {
		if pid == 0 {
			panic(fmt.Sprintf("Instance Start with %v", pid))
		}
	}
	for i := 0; i < len(playerList); i++ {
		bornPos := &inst.logic.World().GetMapConfig().PlayerBornPosList[i]
		inst.logic.PlayerEnterWithStaticInfo(playerList[i], 1, 1, bornPos.X, bornPos.Y, 0)
	}
	inst.playerIdList = playerList
	inst.frameIndexInList = 0
	return true
}

func (inst *Instance) Restart() bool {
	if len(inst.playerIdList) == 0 {
		return false
	}
	if inst.recordHandle != nil {
		inst.recordHandle(inst.logic.World().GetMapConfig(), inst.frameList)
	} else {
		inst.recycleFrameList()
	}
	inst.logic.ReloadScene()
	return inst.CheckAndStart(inst.playerIdList)
}

func (inst *Instance) Pause() {
	inst.logic.Pause()
}

func (inst *Instance) Resume() {
	inst.logic.Resume()
}

func (inst *Instance) PushFrame(frameNum uint32, playerId uint64, cmd CmdCode, args []any) bool {
	if inst.mode != instanceModePlay {
		return false
	}
	if frameNum == 0 {
		frameNum = inst.logic.GetCurrFrame() + 1
	}
	if frameNum > inst.logic.GetCurrFrame()+1 {
		log.Error("core.Instance.PushFrame: push frame %v can not greater to current frame %v", frameNum, inst.logic.GetCurrFrame()+1)
		return false
	}

	var (
		fd *frameData
		l  = len(inst.frameList)
		i  int32
	)
	if l > 0 {
		for i = int32(l) - 1; i >= 0; i-- {
			if inst.frameList[i].frameNum == frameNum {
				fd = inst.frameList[i]
				break
			}
			if inst.frameList[i].frameNum < frameNum {
				break
			}
		}
	}
	if fd == nil {
		fd = inst.getAvailableFrameData()
		fd.frameNum = frameNum
		if i >= int32(l)-1 {
			inst.frameList = append(inst.frameList, fd)
		} else {
			inst.frameList = append(inst.frameList[:i+1], append([]*frameData{fd}, inst.frameList[i+1:]...)...)
		}
	}
	for i := 0; i < len(fd.playerDataList); i++ {
		if fd.playerDataList[i].playerId == 0 {
			fd.playerDataList[i].playerId = inst.playerIdList[i]
		}
		if playerId == fd.playerDataList[i].playerId {
			playerData := fd.playerDataList[i]
			playerData.frameCmdList = append(playerData.frameCmdList, struct {
				frame uint32
				cmd   CmdData
			}{frameNum, CmdData{cmd, args}})
		}
	}
	return true
}

func (inst *Instance) UpdateFrame() {
	inst.processFrameCmdList()
	inst.logic.Update(inst.args.UpdateTick)
}

func (inst *Instance) GetFrame() uint32 {
	return inst.logic.GetCurrFrame()
}

func (inst *Instance) getAvailableFrameData() *frameData {
	fd, o := inst.frameDataFreeList.PopFront()
	if !o {
		playerDataList := make([]*playerData, len(inst.playerIdList))
		fd = &frameData{frameNum: 1, playerDataList: playerDataList}
		for i := 0; i < len(playerDataList); i++ {
			playerDataList[i] = &playerData{}
			playerDataList[i].playerId = inst.playerIdList[i]
		}
	}
	return fd
}

func (inst *Instance) processFrameCmdList() {
	var frameList []*frameData
	if inst.mode == instanceModePlay {
		frameList = inst.frameList
	} else if inst.mode == instanceModeReplay {
		frameList = inst.replay.frameList
	}

	if int(inst.frameIndexInList)+1 > len(frameList) {
		return
	}

	fd := frameList[inst.frameIndexInList]
	logicFrame := inst.logic.GetCurrFrame()
	if logicFrame > fd.frameNum {
		return
	}
	for i := 0; i < len(fd.playerDataList); i++ {
		playerData := fd.playerDataList[i]
		for j := 0; j < len(playerData.frameCmdList); j++ {
			fc := &playerData.frameCmdList[j]
			inst.execCmd(fc.cmd.cmd, fc.cmd.args, playerData.playerId, i)
		}
	}
	inst.frameIndexInList += 1
}

func (inst *Instance) execCmd(cmdCode CmdCode, cmdArgs []any, playerId uint64, playerIndex int) {
	switch cmdCode {
	case CMD_TANK_MOVE:
		dir := cmdArgs[0].(object.Direction)
		orientation := object.Dir2Orientation(dir)
		inst.logic.PlayerTankMove(playerId, orientation)
	case CMD_TANK_STOP:
		inst.logic.PlayerTankStopMove(playerId)
	case CMD_TANK_FIRE:
		inst.logic.PlayerTankFire(playerId)
	case CMD_TANK_ADD_SHELL:
		inst.logic.PlayerTankAddNewShell(playerId, int32(cmdArgs[0].(int)))
	case CMD_TANK_SWITCH_SHELL:
		inst.logic.PlayerTankSwitchShell(playerId)
	case CMD_TANK_SHIELD:
		inst.logic.PlayerTankShield(playerId)
	case CMD_TANK_RESPAWN:
		bornPos := &inst.logic.World().GetMapConfig().PlayerBornPosList[playerIndex]
		inst.logic.PlayerTankRespawn(playerId, 1, 1, bornPos.X, bornPos.Y, 0)
	case CMD_TANK_CHANGE:
		inst.logic.PlayerTankChange(playerId, nil)
	case CMD_TANK_RESTORE:
		inst.logic.PlayerTankRestore(playerId)
	}
}

func (inst *Instance) setRecordHandle(handle func(*game_map.Config, []*frameData)) {
	inst.recordHandle = handle
}

func (inst *Instance) recycleFrameList() {
	for _, f := range inst.frameList {
		f.clear()
		inst.frameDataFreeList.PushBack(f)
	}
	clear(inst.frameList)
	inst.frameList = inst.frameList[:0]
}