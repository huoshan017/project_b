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

type InstanceArgs struct {
	EventMgr   base.IEventManager
	PlayerNum  int32
	UpdateTick time.Duration
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
	args              *InstanceArgs     // 參數
	frameList         []*frameData      // 幀列表
	currFrameIndex    uint32            // 當前幀索引
	playerIdList      []uint64          // 玩家列表
	logic             *common.GameLogic // 游戲邏輯
	frameDataFreeList *list.ListT[*frameData]
}

func NewInstance(args *InstanceArgs) *Instance {
	return &Instance{
		args:              args,
		logic:             common.NewGameLogic(args.EventMgr),
		frameDataFreeList: list.NewListT[*frameData](),
	}
}

func (inst *Instance) LoadScene(config *game_map.Config) bool {
	return inst.logic.LoadScene(config)
}

func (inst *Instance) UnloadScene() {
	inst.logic.UnloadScene()
	for _, f := range inst.frameList {
		f.clear()
		inst.frameDataFreeList.PushBack(f)
	}
	inst.currFrameIndex = 0
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
	inst.playerIdList = playerList
	fd, o := inst.frameDataFreeList.PopFront()
	if o {
		fd.frameNum = 1
		inst.frameList = []*frameData{fd}
	} else {
		playerDataList := make([]*playerData, len(playerList))
		for i := 0; i < len(playerDataList); i++ {
			playerDataList[i] = &playerData{}
			playerDataList[i].playerId = playerList[i]
			bornPos := &inst.logic.CurrentScene().GetMapConfig().PlayerBornPosList[i]
			inst.logic.PlayerEnterWithStaticInfo(playerList[i], 1, 1, bornPos.X, bornPos.Y, 0)
		}
		inst.frameList = []*frameData{{frameNum: 1, playerDataList: playerDataList}}
	}

	return true
}

func (inst *Instance) Restart() bool {
	if len(inst.playerIdList) == 0 {
		return false
	}
	inst.currFrameIndex = 0
	inst.frameList = nil
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
	if frameNum == 0 {
		frameNum = inst.currFrameIndex + 1
	}
	if frameNum > inst.currFrameIndex+1 {
		log.Error("core.Instance.PushFrame: push frame %v can not greater to current frame %v", frameNum, inst.currFrameIndex+1)
		return false
	}
	fd := inst.frameList[inst.currFrameIndex]
	for i := 0; i < len(fd.playerDataList); i++ {
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

func (inst *Instance) processFrameCmdList() {
	fd := inst.frameList[inst.currFrameIndex]
	for i := 0; i < len(fd.playerDataList); i++ {
		playerData := fd.playerDataList[i]
		for j := 0; j < len(playerData.frameCmdList); j++ {
			fc := &playerData.frameCmdList[j]
			inst.execCmd(fc.cmd.cmd, fc.cmd.args, playerData.playerId, i)
		}
	}
	inst.currFrameIndex += 1
	playerDataList := make([]*playerData, len(inst.playerIdList))
	for i := 0; i < len(playerDataList); i++ {
		playerDataList[i] = &playerData{}
		playerDataList[i].playerId = inst.playerIdList[i]
	}
	inst.frameList = append(inst.frameList, &frameData{frameNum: 1, playerDataList: playerDataList})
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
		bornPos := &inst.logic.CurrentScene().GetMapConfig().PlayerBornPosList[playerIndex]
		inst.logic.PlayerTankRespawn(playerId, 1, 1, bornPos.X, bornPos.Y, 0)
	case CMD_TANK_CHANGE:
		inst.logic.PlayerTankChange(playerId, nil)
	case CMD_TANK_RESTORE:
		inst.logic.PlayerTankRestore(playerId)
	}
}
