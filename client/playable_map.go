package main

import (
	"project_b/client/base"
	"project_b/common/log"
	"project_b/common/math"
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
	camera          *base.Camera
}

// 创建可绘制的地图
func CreatePlayableMap(camera *base.Camera) *PlayableMap {
	return &PlayableMap{
		camera: camera,
	}
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
	log.Info("PlayableMap Load map %v", mapId)
	return true
}

/**
 * 把地图上一个矩形窗口范围内的playable绘制到screen
 * @param rect  可視矩形范围，世界坐标，左下角為原點
 * @param dstImage 目标屏幕
 */
func (m *PlayableMap) Draw(rect *math.Rect, dstImage *ebiten.Image) {
	tw, th := int32(len(m.config.Layers[0])), int32(len(m.config.Layers))
	// 获取绘制tiles的索引范围
	if rect.X()+rect.W() <= m.config.X {
		return
	}
	if rect.X() >= m.config.X+tw*m.config.TileWidth {
		return
	}
	if rect.Y() >= m.config.Y+th*m.config.TileHeight {
		return
	}
	if rect.Y()+rect.H() <= m.config.Y {
		return
	}

	var t, b, l, r int32

	if rect.X() <= m.config.X {
		l = 0
	} else {
		l = (rect.X() - m.config.X) / m.config.TileWidth
	}
	if int(l) >= len(m.config.Layers[0]) {
		return
	}

	if rect.X()+rect.W() >= m.config.X+tw*m.config.TileWidth {
		r = int32(len(m.config.Layers[0])) - 1
	} else {
		r = (rect.X() + rect.W() - m.config.X) / m.config.TileWidth
	}
	if int(r) >= len(m.config.Layers[0]) {
		r = int32(len(m.config.Layers[0])) - 1
	}

	if rect.Y() < m.config.Y {
		b = 0
	} else {
		b = (rect.Y() - m.config.Y) / m.config.TileHeight
	}
	if int(b) >= len(m.config.Layers) {
		return
	}

	if rect.Y()+rect.H() >= m.config.Y+th*m.config.TileHeight {
		t = int32(len(m.config.Layers)) - 1
	} else {
		t = (rect.Y() + rect.H() - m.config.Y) / m.config.TileHeight
	}
	if int(t) >= len(m.config.Layers) {
		t = int32(len(m.config.Layers)) - 1
	}

	log.Info("##### b(%v) t(%v) l(%v) r(%v)", b, t, l, r)

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
				m.tilesCacheArray[i][j] = tc
				tc.playableObj.Play()
			}
			tc.op.GeoM.Reset()
			// tile本地坐標到世界坐標的縮放
			sx := m.config.TileWidth / int32(tileAnimConfig.AnimConfig.FrameWidth)
			sy := m.config.TileHeight / int32(tileAnimConfig.AnimConfig.FrameHeight)
			tc.op.GeoM.Scale(float64(sx), float64(sy))
			// todo 注意这里，i是y轴方向，j是x轴方向
			// 由於世界坐標Y軸與屏幕坐標Y軸方向相反，所以變換左上角和右下角的世界坐標到屏幕坐標
			lx, ly := m.camera.World2Screen(m.config.X+j*m.config.TileWidth, m.config.Y+i*m.config.TileHeight+m.config.TileHeight)
			rx, ry := m.camera.World2Screen(m.config.X+j*m.config.TileWidth+m.config.TileWidth, m.config.Y+i*m.config.TileHeight)
			scalex := float64(rx-lx) / float64(m.config.TileWidth)
			scaley := float64(ry-ly) / float64(m.config.TileHeight)
			tc.op.GeoM.Scale(scalex, scaley)
			tc.op.GeoM.Translate(float64(lx), float64(ly))
			tc.playableObj.Draw(dstImage, tc.op)
		}
	}
}
