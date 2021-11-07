package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var (
	bonus_img    *ebiten.Image
	bore_img     *ebiten.Image
	bullet_img   *ebiten.Image
	enemy_img    *ebiten.Image
	explode1_img *ebiten.Image
	explode2_img *ebiten.Image
	flag_img     *ebiten.Image
	gameover_img *ebiten.Image
	misc_img     *ebiten.Image
	num_img      *ebiten.Image
	player1_img  *ebiten.Image
	player2_img  *ebiten.Image
	shield_img   *ebiten.Image
	splash_img   *ebiten.Image
	tile_img     *ebiten.Image
)

func initImages() {
	imgPtrList := []**ebiten.Image{
		&bonus_img, &bore_img, &bullet_img, &enemy_img, &explode1_img, &explode2_img, &flag_img, &gameover_img,
		&misc_img, &num_img, &player1_img, &player2_img, &shield_img, &splash_img, &tile_img,
	}

	imgPathList := []string{
		"png/bonus.png", "png/bore.png", "png/bullet.png", "png/enemy.png", "png/explode1.png", "png/explode2.png",
		"png/flag.png", "png/gameover.png", "png/misc.png", "png/num.png", "png/player1.png", "png/player2.png", "png/shield.png",
		"png/splash.png", "png/tile.png",
	}

	var err error
	for i := 0; i < len(imgPtrList); i++ {
		img := imgPtrList[i]
		*img, _, err = ebitenutil.NewImageFromFile(imgPathList[i])
		if err != nil {
			log.Fatal(err)
		}
	}
}
