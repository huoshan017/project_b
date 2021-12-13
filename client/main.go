package main

import (
	"flag"
	_ "image/png"
	"log"
	"math/rand"
	"os"
	"project_b/client/core"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

const (
	screenWidth  = 1280
	screenHeight = 720
)

type Config struct {
	ServerAddress string
}

var (
	glog *core.Logger
)

func main() {
	if len(os.Args) < 2 {
		log.Printf("args num invalid")
		return
	}

	ip_str := flag.String("ip", "", "ip set")
	flag.Parse()

	glog = core.InitLog("./log/client.log", 2, 100, 30, false, true, 1)

	game := NewGame()
	err := game.Init(&Config{ServerAddress: *ip_str})
	if err != nil {
		return
	}
	defer game.Uninit()
	defer glog.Sync()

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("ProjectB")
	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
