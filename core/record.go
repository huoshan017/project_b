package core

import (
	"fmt"
	"io/ioutil"
	"os"
	"project_b/common/time"
	"project_b/log"

	"github.com/golang/protobuf/proto"
)

type Record struct {
	name         string
	mapId        int32
	frameList    []*frameData
	playerIdList []uint64
	frameNum     uint32
}

func (r *Record) Clear() {
	r.name = ""
	r.mapId = 0
	r.frameList = nil
	r.playerIdList = nil
	r.frameNum = 0
}

func (r Record) MapId() int32 {
	return r.mapId
}

func (r Record) FrameNum() uint32 {
	return r.frameNum
}

type RecordManager struct {
	inst               *Instance
	recordNameList     []string
	recordFileNameList []string
	sel                int32
	savePath           string
	loaded             bool
}

func NewRecordManager(inst *Instance) *RecordManager {
	rm := &RecordManager{
		inst: inst,
		sel:  -1,
	}
	rm.genSavePath()
	return rm
}

func (rm *RecordManager) SetRecord() {
	rm.inst.setRecordHandle(rm.Save)
}

func (rm *RecordManager) LoadRecords() {
	if rm.loaded {
		return
	}
	fsList, err := ioutil.ReadDir(rm.savePath)
	if err != nil {
		log.Error("RecordManager.LoadRecords ioutil.ReadDir err: %v", err)
		return
	}
	var record Record
	for i := 0; i < len(fsList); i++ {
		if !rm.read(fsList[i].Name(), &record) {
			log.Error("RecordManager read file %v failed", fsList[i].Name())
			continue
		}
		rm.recordFileNameList = append(rm.recordFileNameList, fsList[i].Name())
		rm.recordNameList = append(rm.recordNameList, record.name)
	}
	rm.loaded = true
}

func (rm *RecordManager) Save(mapName string, record Record) {
	record.name = fmt.Sprintf("Record.%v: %v", mapName, time.Now().String())
	rm.persistance(&record)
	rm.recordNameList = append(rm.recordNameList, record.name)
}

func (rm *RecordManager) Delete(index int32) bool {
	if int(index) >= len(rm.inst.frameList) {
		return false
	}
	rm.recordFileNameList = append(rm.recordFileNameList[:index], rm.recordFileNameList[index+1:]...)
	rm.recordNameList = append(rm.recordNameList[:index], rm.recordNameList[index+1:]...)
	return true
}

func (rm *RecordManager) Select(index int32) {
	if int(index) <= len(rm.recordNameList) {
		rm.sel = index
	}
}

func (rm *RecordManager) SelectedRecord() (Record, bool) {
	if rm.sel < 0 {
		panic("not selected record")
	}
	name := rm.recordFileNameList[rm.sel]
	var record Record
	if !rm.read(name, &record) {
		log.Error("RecordManager read selected save %v failed", name)
		return record, false
	}
	return record, true
}

func (rm *RecordManager) RecordNameList() []string {
	var nameList []string
	for i := 0; i < len(rm.recordNameList); i++ {
		nameList = append(nameList, rm.recordNameList[i])
	}
	return nameList
}

func (rm *RecordManager) GetRecordCount() int32 {
	return int32(len(rm.recordNameList))
}

func (rm *RecordManager) GetRecordName(index int32) string {
	return rm.recordNameList[index]
}

func (rm *RecordManager) genSavePath() {
	dir, err := os.Getwd()
	if err != nil {
		log.Error("Record persistance os.Getwd err: %v", err)
		return
	}

	savePath := dir + "/" + rm.inst.args.SavePath
	err = os.MkdirAll(savePath, os.ModePerm)
	if err != nil {
		log.Error("Record persistance make dir err: %v", err)
		return
	}
	err = os.Chmod(savePath, os.ModePerm)
	if err != nil {
		log.Error("Record persistance chmod err: %v", err)
		return
	}

	rm.savePath = savePath
}

func (rm *RecordManager) read(fileName string, record *Record) bool {
	filePath := fmt.Sprintf("%v/%v", rm.savePath, fileName)
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Error("RecordManager read file err: %v", err)
		return false
	}
	var pbr PbRecord
	err = proto.Unmarshal(data, &pbr)
	if err != nil {
		log.Error("RecordManager unmarshal err: %v", err)
		return false
	}
	unserializeRecord(&pbr, record)
	return true
}

func (rm *RecordManager) persistance(record *Record) {
	var pbr PbRecord
	serializeRecord(record, &pbr)
	data, err := proto.Marshal(&pbr)
	if err != nil {
		log.Error("Record persistance marshal err: %v", err)
		return
	}

	fileName := fmt.Sprintf("%v.record", time.Now().Unix())
	filePath := fmt.Sprintf("%v/%v", rm.savePath, fileName)
	var f *os.File
	f, err = os.Create(filePath)
	if err != nil {
		log.Error("Record persistance create err: %v", err)
		return
	}
	defer f.Close()

	_, err = f.Write(data)
	if err != nil {
		log.Error("Record persistance write err: %v", err)
		return
	}

	f.Sync()

	rm.recordFileNameList = append(rm.recordFileNameList, fileName)
}

func serializeRecord(record *Record, pbr *PbRecord) {
	pbr.Name = record.name
	pbr.MapId = record.mapId
	pbr.FrameNum = record.frameNum
	var frameList []*PbFrameData
	for i := 0; i < len(record.frameList); i++ {
		var playerFrameList []*PbPlayerFrame
		fd := record.frameList[i]
		for j := 0; j < len(fd.playerDataList); j++ {
			var cmdList []*PbCmdData
			pd := fd.playerDataList[j]
			for k := 0; k < len(pd.cmdList); k++ {
				cmd := pd.cmdList[k]
				cmdList = append(cmdList, &PbCmdData{Code: int32(cmd.cmd), Args: cmd.args})
			}
			playerFrameList = append(playerFrameList, &PbPlayerFrame{PlayerId: pd.playerId, CmdList: cmdList})
		}
		frameList = append(frameList, &PbFrameData{FrameNum: fd.frameNum, PlayerList: playerFrameList})
	}
	pbr.FrameList = frameList
}

func unserializeRecord(pbr *PbRecord, record *Record) {
	record.name = pbr.Name
	record.mapId = pbr.MapId
	record.frameNum = pbr.FrameNum
	var (
		frameList    []*frameData
		playerIdList []uint64
	)
	for i := 0; i < len(pbr.FrameList); i++ {
		var playerFrameList []*playerData
		fd := pbr.FrameList[i]
		for j := 0; j < len(fd.PlayerList); j++ {
			var cmdList []CmdData
			pd := fd.PlayerList[j]
			for k := 0; k < len(pd.CmdList); k++ {
				cmd := pd.CmdList[k]
				cmdList = append(cmdList, CmdData{cmd: CmdCode(cmd.GetCode()), args: cmd.GetArgs()})
			}
			if len(playerIdList) < len(fd.PlayerList) {
				playerIdList = append(playerIdList, pd.PlayerId)
			}
			playerFrameList = append(playerFrameList, &playerData{playerId: pd.PlayerId, cmdList: cmdList})
		}
		frameList = append(frameList, &frameData{frameNum: fd.FrameNum, playerDataList: playerFrameList})
	}
	record.playerIdList = playerIdList
	record.frameList = frameList
}
