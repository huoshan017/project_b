package main

import (
	"errors"
	"flag"
	"os"

	"project_b/game_proto"
	"project_b/log"

	gsnet_common "github.com/huoshan017/gsnet/common"
	gsnet_msg "github.com/huoshan017/gsnet/msg"
	gsnet_options "github.com/huoshan017/gsnet/options"
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

func CreateGameMsgHandlerWrapper(args ...any) gsnet_msg.IMsgSessionHandler {
	handler := CreateGameMsgHandler(args[0].(*GameService))
	return gsnet_msg.IMsgSessionHandler(handler)
}

func (s *GameService) Init(conf *config) bool {
	// 错误注册
	gsnet_common.RegisterNoDisconnectError(ErrKickDuplicatePlayer)
	// 创建服务
	net := gsnet_msg.NewProtobufMsgServer(CreateGameMsgHandlerWrapper, []any{s}, gsnet_msg.CreateIdMsgMapperWith(game_proto.Id2MsgMapOnServer), gsnet_options.WithNoDelay(true))
	// 监听
	err := net.Listen(conf.addr)
	if err != nil {
		log.Error("game service listen addr %v err: %v", conf.addr, err)
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
				log.Stack(err)
			}
		}()
		s.gameLogicThread.Run()
	}()
	s.server.Serve()
}

func (s *GameService) End() {
	s.gameLogicThread.Close()
	s.server.End()
}

func main() {
	logger := log.InitLog("./log/game_server.log", 2, 100, 30, false, true, -1)
	defer logger.Sync()

	if len(os.Args) < 2 {
		log.Error("args num invalid")
		return
	}
	ip_str := flag.String("ip", "", "ip set")
	flag.Parse()

	gameService := NewGameService()
	if !gameService.Init(&config{addr: *ip_str, mapIndex: 0}) {
		return
	}
	defer gameService.End()

	log.Info("game service started")

	gameService.Start()
}
