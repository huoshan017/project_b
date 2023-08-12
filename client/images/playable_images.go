package images

import (
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var (
	bonus_img     *ebiten.Image
	bore_img      *ebiten.Image
	bullet_img    *ebiten.Image
	enemy_img     *ebiten.Image
	explode1_img  *ebiten.Image
	explode2_img  *ebiten.Image
	flag_img      *ebiten.Image
	gameover_img  *ebiten.Image
	misc_img      *ebiten.Image
	num_img       *ebiten.Image
	player1_img   *ebiten.Image
	player2_img   *ebiten.Image
	shield_img    *ebiten.Image
	splash_img    *ebiten.Image
	tile_img      *ebiten.Image
	smallball_img *ebiten.Image
	shell_img     *ebiten.Image
	title_img     *ebiten.Image
)

func GetBonusImg() *ebiten.Image {
	return bonus_img
}

func GetBoreImg() *ebiten.Image {
	return bore_img
}

func GetBulletImg() *ebiten.Image {
	return bullet_img
}

func GetEnemyImg() *ebiten.Image {
	return enemy_img
}

func GetExplode1Img() *ebiten.Image {
	return explode1_img
}

func GetExplode2Img() *ebiten.Image {
	return explode2_img
}

func GetFlagImg() *ebiten.Image {
	return flag_img
}

func GetGameOverImg() *ebiten.Image {
	return gameover_img
}

func GetMiscImg() *ebiten.Image {
	return misc_img
}

func GetNumImg() *ebiten.Image {
	return num_img
}

func GetPlayer1Img() *ebiten.Image {
	return player1_img
}

func GetPlayer2Img() *ebiten.Image {
	return player2_img
}

func GetShieldImg() *ebiten.Image {
	return shield_img
}

func GetSplashImg() *ebiten.Image {
	return splash_img
}

func GetTileImg() *ebiten.Image {
	return tile_img
}

func GetSmallBallImg() *ebiten.Image {
	return smallball_img
}

func GetShellImg() *ebiten.Image {
	return shell_img
}

func GetTitleImg() *ebiten.Image {
	return title_img
}

func InitImages() {
	imgPtrList := []**ebiten.Image{
		&bonus_img,
		&bore_img,
		&bullet_img,
		&enemy_img,
		&explode1_img,
		&explode2_img,
		&flag_img,
		&gameover_img,
		&misc_img,
		&num_img,
		&player1_img,
		&player2_img,
		&shield_img,
		&splash_img,
		&tile_img,
		&smallball_img,
		&shell_img,
		&title_img,
	}

	imgPathList := []string{
		"png/bonus.png",
		"png/bore.png",
		"png/bullet.png",
		"png/enemy.png",
		"png/explode1.png",
		"png/explode2.png",
		"png/flag.png",
		"png/gameover.png",
		"png/misc.png",
		"png/num.png",
		"png/player1.png",
		"png/player2.png",
		"png/shield.png",
		"png/splash.png",
		"png/tile.png",
		"png/small_ball.png",
		"png/shell.png",
		"png/title.png",
	}

	var err error
	for i := 0; i < len(imgPtrList); i++ {
		img := imgPtrList[i]
		*img, _, err = ebitenutil.NewImageFromFile(imgPathList[i])
		if err != nil {
			log.Fatalf("load png %v err: %v", imgPathList[i], err)
		} else {
			log.Printf("load png %v done", imgPathList[i])
		}
	}
}
