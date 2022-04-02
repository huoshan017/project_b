package main

import (
	"errors"
	"flag"
	"os"
	"time"

	"project_b/common/log"
	"project_b/game_proto"

	gsnet_common "github.com/huoshan017/gsnet/common"
	gsnet_msg "github.com/huoshan017/gsnet/msg"
)

var ErrKickDuplicatePlayer = errors.New("game service example: kick duplicate player")

type config struct {
	addr     string
	mapIndex int32
}

type GameService struct {
	net             *gsnet_msg.MsgServer
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
	gsnet_common.RegisterNoDisconnectError(ErrKickDuplicatePlayer)
	// 创建服务
	net := gsnet_msg.NewPBMsgServer(gsnet_msg.CreateIdMsgMapperWith(game_proto.Id2MsgMapOnServer))
	// 创建会话和消息处理器
	msgHandler := CreateGameMsgHandler(s)
	// 设置处理器到服务
	handles := struct {
		ConnectedHandle    func(*gsnet_msg.MsgSession)
		DisconnectedHandle func(*gsnet_msg.MsgSession, error)
		TickHandle         func(*gsnet_msg.MsgSession, time.Duration)
		ErrorHandle        func(error)
		MsgHandles         map[gsnet_msg.MsgIdType]func(*gsnet_msg.MsgSession, interface{}) error
	}{
		msgHandler.OnConnect,
		msgHandler.OnDisconnect,
		msgHandler.OnTick,
		msgHandler.OnError,
		msgHandler.getMsgId2HandleMap(),
	}
	net.SetSessionHandles(handles)
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

func (s *GameService) End() {
	s.gameLogicThread.Close()
	s.net.End()
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
