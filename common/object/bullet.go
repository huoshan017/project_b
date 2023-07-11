package object

import "unsafe"

// 子弹
type Bullet struct {
	MovableObject
}

// 创建车辆
func NewBullet(instId uint32, staticInfo *BulletStaticInfo) *Bullet {
	o := &Bullet{}
	o.Init(instId, &staticInfo.ObjStaticInfo)
	return o
}

// 靜態配置
func (b *Bullet) BulletStaticInfo() *BulletStaticInfo {
	return (*BulletStaticInfo)(unsafe.Pointer(b.staticInfo))
}

// 初始化
func (b *Bullet) Init(instId uint32, staticInfo *ObjStaticInfo) {
	b.MovableObject.Init(instId, staticInfo)
	b.setSuper(b)
}

// 反初始化
func (b *Bullet) Uninit() {
	b.MovableObject.Uninit()
}
