package common

import (
	"project_b/common/base"
	"project_b/common/object"
	"testing"
)

func TestGameLogic_PlayerTankChange(t *testing.T) {
	type fields struct {
		eventMgr base.IEventManager
		sceneMap *SceneMap
		state    int32
		mapIndex int32
	}
	type args struct {
		uid        uint64
		staticInfo *object.TankStaticInfo
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GameLogic{
				eventMgr: tt.fields.eventMgr,
				sceneMap: tt.fields.sceneMap,
				state:    tt.fields.state,
				mapIndex: tt.fields.mapIndex,
			}
			if got := g.PlayerTankChange(tt.args.uid, tt.args.staticInfo); got != tt.want {
				t.Errorf("GameLogic.PlayerTankChange() = %v, want %v", got, tt.want)
			}
		})
	}
}
