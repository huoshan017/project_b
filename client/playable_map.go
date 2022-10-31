package main

import (
	"project_b/client/base"
	"project_b/common/object"

	"github.com/hajimehoshi/ebiten/v2"
)

type tileOpCache struct {
	op          *ebiten.DrawImageOptions
	playableObj *PlayableStaticObject
}

type PlayableMap struct {
	mapInfo
	staticObjArray  [][]*object.StaticObject
	tilesImg        *ebiten.Image
	tilesCacheArray [][]*tileOpCache // 缓存TilesImg生成的SubImage和Op，提升性能
	viewport        *base.Viewport
}

// 创建可绘制的地图
func CreatePlayableMap() *PlayableMap {
	return &PlayableMap{}
}

/**
 * 载入地图
 * @param mapId  地图Id
 * @param staticObjArray  静态物体数组
 * @return 是否载入成功
 */
func (m *PlayableMap) Load(mapId int32, staticObjArray [][]*object.StaticObject) bool {
	m.mapInfo = mapInfoArray[mapId]
	m.staticObjArray = staticObjArray
	m.tilesImg = tile_img
	m.tilesCacheArray = make([][]*tileOpCache, len(m.config.Layers))
	for i := 0; i < len(m.config.Layers); i++ {
		m.tilesCacheArray[i] = make([]*tileOpCache, len(m.config.Layers[0]))
	}
	glog.Info("PlayableMap Load map %v", mapId)
	return true
}

/**
 * 设置视口
 * @param viewport  视口
 */
func (m *PlayableMap) SetViewport(viewport *base.Viewport) {
	m.viewport = viewport
}

/**
 * 把地图上一个矩形窗口范围内的playable绘制到screen
 * @param rect  可視矩形范围，世界坐标，左下角為原點
 * @param dstImage 目标屏幕
 */
func (m *PlayableMap) Draw(rect *base.Rect, op *ebiten.DrawImageOptions, dstImage *ebiten.Image) {
	// 获取绘制tiles的索引范围
	l := (rect.X() - m.config.X) / m.tileSize
	if l < 0 {
		l = 0
	}
	if l >= m.config.Width/m.tileSize {
		l = m.config.Width/m.tileSize - 1
	}
	r := (rect.X() - m.config.X + rect.W()) / m.tileSize
	if r < 0 {
		r = 0
	}
	if r >= m.config.Width/m.tileSize {
		r = m.config.Width/m.tileSize - 1
	}
	b := (rect.Y() - m.config.Y) / m.tileSize
	if b < 0 {
		b = 0
	}
	if b >= m.config.Height/m.tileSize {
		b = m.config.Height/m.tileSize - 1
	}
	t := (rect.Y() - m.config.Y + rect.H()) / m.tileSize
	if t < 0 {
		t = 0
	}
	if t >= m.config.Height/m.tileSize {
		t = m.config.Height/m.tileSize - 1
	}
	ly := (rect.X() - m.config.X) % m.tileSize
	by := (rect.Y() - m.config.Y) % m.tileSize

	// 按照世界坐標系的坐標軸方向遍歷tiles數組繪製
	// y坐标，從下到上
	for i := b; i <= t; i++ {
		// x坐标，從左到右
		for j := l; j <= r; j++ {
			// 瓦片类型
			v := m.config.Layers[i][j]
			if object.StaticObjType(v) == object.StaticObjNone {
				continue
			}
			tileAnimConfig := GetStaticObjAnimConfig(object.StaticObjType(v))
			if tileAnimConfig == nil {
				glog.Error("can't get static object anim by type %v", v)
				continue
			}
			tc := m.tilesCacheArray[i][j]
			if tc == nil {
				tc = &tileOpCache{
					playableObj: NewPlayableStaticObject(m.staticObjArray[i][j], tileAnimConfig),
					op:          &ebiten.DrawImageOptions{},
				}
				//sx := (v % m.tileXNum) * m.tileSize
				//sy := (v / m.tileXNum) * m.tileSize
				//op := &ebiten.DrawImageOptions{}
				//op.GeoM.Translate(float64((int32(j)%worldSizeX)*m.tileSize), float64((int32(j)/worldSizeY)*m.tileSize))
				//tc = &tileOpCache{
				//	img: m.tilesImg.SubImage(image.Rect(int(sx), int(sy), int(sx+m.tileSize), int(sy+m.tileSize))).(*ebiten.Image),
				//	op:  op,
				//}
				m.tilesCacheArray[i][j] = tc
			}
			tc.op.GeoM.Reset()
			// tile圖片與世界坐標尺寸的縮放比例
			tc.op.GeoM.Scale(multiplesObjLenAndDisplayLen, multiplesObjLenAndDisplayLen)
			// todo 注意这里，i是y轴方向，j是x轴方向
			tc.op.GeoM.Translate(-float64(by)+float64(j*m.tileSize), -float64(ly)+float64(i*m.tileSize))
			tc.op.GeoM.Concat(op.GeoM)
			tc.playableObj.Draw(dstImage, tc.op)
		}
	}
}
