package object

import (
	"project_b/common/base"
	"unsafe"
)

// 環繞運動物體
type SurroundObj struct {
	*MovableObject
	aroundCenterObjInstId  uint32               // 環繞物體實例id
	getAroundCenterObjFunc func(uint32) IObject // 獲得環繞物體函數
	turnAngle              int32                // 轉過角度
	accumulateMs           int32                // 纍計時間(毫秒)
	lateUpdateEvent        base.Event           // 后更新事件
}

func NewSurroundObj() *SurroundObj {
	obj := &SurroundObj{}
	obj.MovableObject = NewMovableObjectWithSuper(obj)
	return obj
}

func (b *SurroundObj) Init(instId uint32, staticInfo *ObjStaticInfo) {
	b.MovableObject.Init(instId, staticInfo)
}

func (b *SurroundObj) Uninit() {
	b.aroundCenterObjInstId = 0
	b.getAroundCenterObjFunc = nil
	b.turnAngle = 0
	b.accumulateMs = 0
	b.MovableObject.Uninit()
}

func (b *SurroundObj) SetAroundCenterObject(instId uint32, getFunc func(uint32) IObject) {
	b.aroundCenterObjInstId = instId
	b.getAroundCenterObjFunc = getFunc
	centerObj := b.getCenterObj()
	if centerObj == nil {
		return
	}
	ox, oy := centerObj.Pos()
	b.SetPos(ox+b.SurroundObjStaticInfo().AroundRadius, oy)
}

func (b *SurroundObj) GetAroundCenterObject() IObject {
	return b.getCenterObj()
}

func (b *SurroundObj) SurroundObjStaticInfo() *SurroundObjStaticInfo {
	return (*SurroundObjStaticInfo)(unsafe.Pointer(b.staticInfo))
}

func (b *SurroundObj) getCenterObj() IObject {
	if b.aroundCenterObjInstId <= 0 || b.getAroundCenterObjFunc == nil {
		return nil
	}

	aroundObj := b.getAroundCenterObjFunc(b.aroundCenterObjInstId)
	// 找不到環繞物體時，説明環繞物體被銷毀
	if aroundObj == nil || aroundObj.IsRecycle() {
		b.ToRecycle()
		return nil
	}
	return aroundObj
}

func (b *SurroundObj) Update(tickMs uint32) {
	if b.getCenterObj() == nil {
		return
	}
	b.MovableObject.Update(tickMs)
	if b.state == isMoving {
		b.lateUpdateEvent.Call(b.turnAngle, b.accumulateMs)
	}
}

func (b *SurroundObj) RegisterLateUpdateEventHandle(handle func(args ...any)) {
	b.lateUpdateEvent.Register(handle)
}

func (b *SurroundObj) UnregisterLateUpdateEventHandle(handle func(args ...any)) {
	b.lateUpdateEvent.Unregister(handle)
}
