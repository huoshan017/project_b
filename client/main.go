package main

import (
	"flag"
	_ "image/png"
	"math/rand"
	"os"
	"time"

	"project_b/common/log"

	"github.com/hajimehoshi/ebiten/v2"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

const (
	screenWidth  = 1280
	screenHeight = 720
)

var gslog *log.Logger

func getLog() *log.Logger {
	if gslog == nil {
		gslog = log.NewWithConfig(&log.LogConfig{
			Filename:      "./log/client.log",
			MaxSize:       2,
			MaxBackups:    100,
			MaxAge:        30,
			Compress:      false,
			ConsoleOutput: true,
		}, log.DebugLevel)
	}
	return gslog
}

type Config struct {
	ServerAddress string
}

func main() {
	if len(os.Args) < 2 {
		getLog().Error("args num invalid")
		return
	}
	ip_str := flag.String("ip", "", "ip set")
	flag.Parse()

	game := NewGame()
	err := game.Init(&Config{ServerAddress: *ip_str})
	if err != nil {
		return
	}
	defer game.Uninit()
	defer getLog().Sync()

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("ProjectB")
	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
