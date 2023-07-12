package object

import (
	"project_b/common/time"
	"unsafe"
)

// 環繞運動物體
type SurroundObj struct {
	MovableObject                               // todo 暫時作爲靜態物體處理，object系統要等完善的時候再作爲可移動物體處理
	aroundCenterObjInstId  uint32               // 環繞物體實例id
	getAroundCenterObjFunc func(uint32) IObject // 獲得環繞物體函數
	turnAngle              int32                // 轉過角度
	accumulateTime         time.Duration        // 纍計時間
}

func NewSurroundObj(instId uint32, staticInfo *SurroundObjStaticInfo) *SurroundObj {
	obj := &SurroundObj{}
	obj.Init(instId, &staticInfo.ObjStaticInfo)
	return obj
}

func (b *SurroundObj) Init(instId uint32, staticInfo *ObjStaticInfo) {
	b.MovableObject.Init(instId, staticInfo)
	b.setSuper(b)
}

func (b *SurroundObj) Uninit() {
	b.aroundCenterObjInstId = 0
	b.getAroundCenterObjFunc = nil
	b.turnAngle = 0
	b.accumulateTime = 0
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
