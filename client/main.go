package main

import (
	"flag"
	"log"
	"math/rand"
	"os"
	"project_b/client/base"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

const (
	screenWidth  = 640
	screenHeight = 480
)

type Config struct {
	cameraFov     int32
	serverAddress string
}

var (
	glog *base.Logger
)

func main() {
	if len(os.Args) < 2 {
		log.Printf("args num invalid")
		return
	}

	ip_str := flag.String("ip", "", "ip set")
	flag.Parse()

	initResources()

	glog = base.InitLog("./log/client.log", 2, 100, 30, false, true, 1)

	game := NewGame(&Config{cameraFov: 120, serverAddress: *ip_str})
	err := game.Init()
	if err != nil {
		glog.Error("game init err: %v", err)
		return
	}
	defer game.Uninit()
	defer glog.Sync()

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("ProjectB")
	if err := ebiten.RunGame(game); err != nil {
		glog.Error("game run err: %v", err)
	}
}
