package main

import (
	"flag"
	"os"
	"project_b/common/time"
	"project_b/common_data"
	"project_b/log"

	"net/http"
	_ "net/http/pprof"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 1280
	screenHeight = 720
)

type Config struct {
	serverAddress string
	playerCount   int32
	updateTick    time.Duration
}

func main() {
	if len(os.Args) < 2 {
		return
	}

	ip_str := flag.String("ip", "", "ip set")
	flag.Parse()

	initResources()

	logger := log.InitLog("./log/client.log", 2, 100, 30, false, true, -1)

	game := NewGame(&Config{serverAddress: *ip_str, playerCount: 1, updateTick: common_data.GameLogicTick})
	err := game.Init()
	if err != nil {
		log.Error("game init err: %v", err)
		return
	}
	defer game.Uninit()
	defer logger.Sync()

	go func() {
		http.ListenAndServe("0.0.0.0:6060", nil)
	}()

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("BattleCity")
	if err := ebiten.RunGame(game); err != nil {
		log.Error("game run err: %v", err)
	}
}
