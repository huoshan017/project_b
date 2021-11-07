package main

import (
	"errors"
	"flag"
	"os"

	"project_b/common/log"

	"github.com/huoshan017/gsnet"
)

var ErrKickDuplicatePlayer = errors.New("game service example: kick duplicate player")

type config struct {
	addr string
}

type GameService struct {
	net             *gsnet.Server
	loginCheckMgr   *KeyCheckManager
	enterCheckMgr   *KeyCheckManager
	playerMgr       *SPlayerManager
	gameLogicThread *GameLogicThread
}

func NewGameService() *GameService {
	return &GameService{}
}

func (s *GameService) Init(conf *config) bool {
	// 错误注册
	gsnet.RegisterNoDisconnectError(ErrKickDuplicatePlayer)
	// 创建服务
	net := gsnet.NewServer(&GameMsgHandler{}, []interface{}{s})
	// 监听
	err := net.Listen(conf.addr)
	if err != nil {
		gslog.Error("game service listen addr %v err: %v", conf.addr, err)
		return false
	}
	s.net = net
	s.loginCheckMgr = NewKeyCheckManager()
	s.enterCheckMgr = NewKeyCheckManager()
	s.playerMgr = NewSPlayerManager()
	s.gameLogicThread = CreateGameLogicThread()
	return true
}

func (s *GameService) Start() {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				gslog.Stack(err)
			}
		}()
		s.gameLogicThread.Run()
	}()
	s.net.Start()
}

var gslog *log.Logger

func main() {
	if len(os.Args) < 2 {
		gslog.Error("args num invalid")
		return
	}
	ip_str := flag.String("ip", "", "ip set")
	flag.Parse()

	gslog = log.NewWithConfig(&log.LogConfig{
		Filename:      "./log/game_server.log",
		MaxSize:       2,
		MaxBackups:    100,
		MaxAge:        30,
		Compress:      false,
		ConsoleOutput: true,
	}, log.DebugLevel)
	defer gslog.Sync()

	gameService := NewGameService()
	if !gameService.Init(&config{addr: *ip_str}) {
		return
	}

	gslog.Info("game service started")

	gameService.Start()
}
