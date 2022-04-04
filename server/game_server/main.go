package main

import (
	"errors"
	"flag"
	"os"

	"project_b/common/log"
	"project_b/game_proto"

	gsnet_common "github.com/huoshan017/gsnet/common"
	gsnet_msg "github.com/huoshan017/gsnet/msg"
	gsnet_server "github.com/huoshan017/gsnet/server"
)

var ErrKickDuplicatePlayer = errors.New("game service example: kick duplicate player")

type config struct {
	addr     string
	mapIndex int32
}

type GameService struct {
	server          *gsnet_msg.MsgServer
	loginCheckMgr   *KeyCheckManager
	enterCheckMgr   *KeyCheckManager
	playerMgr       *SPlayerManager
	gameLogicThread *GameLogicThread
}

func NewGameService() *GameService {
	return &GameService{}
}

func CreateGameMsgHandlerWrapper(args ...interface{}) gsnet_msg.IMsgSessionHandler {
	handler := CreateGameMsgHandler(args[0].(*GameService))
	return gsnet_msg.IMsgSessionHandler(handler)
}

func (s *GameService) Init(conf *config) bool {
	// 错误注册
	gsnet_common.RegisterNoDisconnectError(ErrKickDuplicatePlayer)
	// 创建服务
	net := gsnet_msg.NewPBMsgServer(CreateGameMsgHandlerWrapper, gsnet_msg.CreateIdMsgMapperWith(game_proto.Id2MsgMapOnServer), gsnet_server.WithNewSessionHandlerFuncArgs(s))
	// 监听
	err := net.Listen(conf.addr)
	if err != nil {
		gslog.Error("game service listen addr %v err: %v", conf.addr, err)
		return false
	}
	s.server = net
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
	s.server.Start()
}

func (s *GameService) End() {
	s.gameLogicThread.Close()
	s.server.End()
}

var gslog *log.Logger

func main() {
	gslog = log.NewWithConfig(&log.LogConfig{
		Filename:      "./log/game_server.log",
		MaxSize:       2,
		MaxBackups:    100,
		MaxAge:        30,
		Compress:      false,
		ConsoleOutput: true,
	}, log.DebugLevel)
	defer gslog.Sync()

	if len(os.Args) < 2 {
		gslog.Error("args num invalid")
		return
	}
	ip_str := flag.String("ip", "", "ip set")
	flag.Parse()

	gameService := NewGameService()
	if !gameService.Init(&config{addr: *ip_str, mapIndex: 0}) {
		return
	}
	defer gameService.End()

	gslog.Info("game service started")

	gameService.Start()
}
