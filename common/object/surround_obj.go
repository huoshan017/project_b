package object

import (
	"project_b/common/time"
	"unsafe"
)

// 環繞運動物體
type SurroundObj struct {
	MovableObject                         // todo 暫時作爲靜態物體處理，object系統要等完善的時候再作爲可移動物體處理
	centerObjInstId  uint32               // 環繞物體實例id
	getCenterObjFunc func(uint32) IObject // 獲得環繞物體函數
	turnAngle        int32                // 轉過角度
	accumulateTime   time.Duration        // 纍計時間
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
	b.centerObjInstId = 0
	b.getCenterObjFunc = nil
	b.turnAngle = 0
	b.accumulateTime = 0
	b.MovableObject.Uninit()
}

func (b *SurroundObj) SetCenterObject(instId uint32, getFunc func(uint32) IObject) {
	b.centerObjInstId = instId
	b.getCenterObjFunc = getFunc
	centerObj := b.getCenterObj()
	if centerObj == nil {
		return
	}
	ox, oy := centerObj.Pos()
	b.SetPos(ox, oy+b.SurroundObjStaticInfo().AroundRadius)
}

func (b *SurroundObj) GetCenterPos() (int32, int32, bool) {
	centerObj := b.getCenterObj()
	if centerObj == nil {
		return -1, -1, false
	}
	x, y := centerObj.Pos()
	return x, y, true
}

func (b *SurroundObj) SurroundObjStaticInfo() *SurroundObjStaticInfo {
	return (*SurroundObjStaticInfo)(unsafe.Pointer(b.staticInfo))
}

func (b *SurroundObj) getCenterObj() IObject {
	if b.centerObjInstId <= 0 || b.getCenterObjFunc == nil {
		return nil
	}

	aroundObj := b.getCenterObjFunc(b.centerObjInstId)
	// 找不到環繞物體時，説明環繞物體被銷毀
	if aroundObj == nil || aroundObj.IsRecycle() {
		b.ToRecycle()
		return nil
	}
	return aroundObj
}
