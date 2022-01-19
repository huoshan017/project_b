package main

import (
	"bytes"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

type tileOpCache struct {
	img *ebiten.Image
	op  *ebiten.DrawImageOptions
}

type Map struct {
	mapInfo
	TilesImg        *ebiten.Image
	worldImg        *ebiten.Image
	tilesCacheArray [][]*tileOpCache // 缓存TilesImg生成的SubImage和Op
}

func (m *Map) Load(mapIndex int32) bool {
	d := mapInfoArray[mapIndex]
	img, _, err := image.Decode(bytes.NewReader(d.tilesImageData))
	if err != nil {
		log.Fatal(err)
	}

	m.mapInfo = d
	m.TilesImg = ebiten.NewImageFromImage(img)
	m.worldImg = ebiten.NewImage(m.config.Width, m.config.Height)
	return true
}

// 绘制到一个矩形窗口，这里是逻辑坐标
func (m *Map) Draw(left, top, width, height int32) {
	tileSize := m.tileSize
	worldSizeX := m.config.Width / tileSize
	worldSizeY := m.config.Height / tileSize

	if m.tilesCacheArray == nil {
		m.tilesCacheArray = make([][]*tileOpCache, len(m.config.Layers))
	}

	for i, l := range m.config.Layers {
		if m.tilesCacheArray[i] == nil {
			m.tilesCacheArray[i] = make([]*tileOpCache, len(l))
		}
		for j, t := range l {
			tc := m.tilesCacheArray[i][j]
			if tc == nil {
				sx := (t % m.tileXNum) * tileSize
				sy := (t / m.tileXNum) * tileSize
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64((j%worldSizeX)*tileSize), float64((j/worldSizeY)*tileSize))
				tc = &tileOpCache{
					img: m.TilesImg.SubImage(image.Rect(sx, sy, sx+tileSize, sy+tileSize)).(*ebiten.Image),
					op:  op,
				}
				m.tilesCacheArray[i][j] = tc
			}
			m.worldImg.DrawImage(tc.img, tc.op)
		}
	}
}

func (m *Map) GetImage() *ebiten.Image {
	return m.worldImg
}
