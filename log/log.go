package log

import (
	"io"
	"os"

	"github.com/natefinch/lumberjack"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Level = zapcore.Level

const (
	InfoLevel   Level = zap.InfoLevel
	WarnLevel   Level = zap.WarnLevel
	ErrorLevel  Level = zap.ErrorLevel
	DPanicLevel Level = zap.DPanicLevel
	PanicLevel  Level = zap.PanicLevel
	FatalLevel  Level = zap.FatalLevel
	DebugLevel  Level = zap.DebugLevel
)

type Field = zap.Field

type Logger struct {
	l     *zap.SugaredLogger
	level Level
}

func (l *Logger) Debug(msg string, args ...interface{}) {
	l.l.Debugf(msg, args...)
}

func (l *Logger) Info(msg string, args ...interface{}) {
	l.l.Infof(msg, args...)
}

func (l *Logger) Warn(msg string, args ...interface{}) {
	l.l.Warnf(msg, args...)
}

func (l *Logger) Error(msg string, args ...interface{}) {
	l.l.Errorf(msg, args...)
}

func (l *Logger) DPanic(msg string, args ...interface{}) {
	l.l.DPanicf(msg, args...)
}

func (l *Logger) Panic(msg string, args ...interface{}) {
	l.l.Panicf(msg, args...)
}

func (l *Logger) Fatal(msg string, args ...interface{}) {
	l.l.Fatalf(msg, args...)
}

func (l *Logger) Stack(err interface{}) {
	er := errors.Errorf("%v", err)
	l.l.Errorf("%+v", er)
}

func (l *Logger) Sync() error {
	return l.l.Sync()
}

type LogConfig struct {
	Filename      string `json:"filename" yaml:"filename"`
	MaxSize       int    `json:"maxsize" yaml:"maxsize"`
	MaxAge        int    `json:"maxage" yaml:"maxage"`
	MaxBackups    int    `json:"maxbackups" yaml:"maxbackups"`
	LocalTime     bool   `json:"localtime" yaml:"localtime"`
	Compress      bool   `json:"compress" yaml:"compress"`
	ConsoleOutput bool   `json:"console" yaml:"console"`
}

func NewWithConfig(config *LogConfig, level Level) *Logger {
	var infoFileWriteSyncer = zapcore.AddSync(&lumberjack.Logger{
		Filename:   config.Filename,   //日志文件存放目录，如果文件夹不存在会自动创建
		MaxSize:    config.MaxSize,    //文件大小限制,单位MB
		MaxBackups: config.MaxBackups, //最大保留日志文件数量
		MaxAge:     config.MaxAge,     //日志文件保留天数
		Compress:   config.Compress,   //是否压缩处理
	})
	if config.ConsoleOutput {
		infoFileWriteSyncer = zapcore.NewMultiWriteSyncer(infoFileWriteSyncer, zapcore.AddSync(os.Stdout))
	}
	return New(infoFileWriteSyncer, level)
}

func New(writer io.Writer, level Level) *Logger {
	if writer == nil {
		panic("the writer is nil")
	}

	cfg := zap.NewProductionEncoderConfig()
	cfg.EncodeTime = zapcore.RFC3339TimeEncoder
	cfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(cfg),
		zapcore.AddSync(writer),
		zapcore.Level(level),
	)

	l := zap.New(core)
	logger := &Logger{
		l:     l.Sugar(),
		level: level,
	}
	return logger
}

var clog = New(os.Stdout, zapcore.DebugLevel)

func Debug(msg string, args ...interface{}) {
	clog.Debug(msg, args...)
}

func Info(msg string, args ...interface{}) {
	clog.Info(msg, args...)
}

func Warn(msg string, args ...interface{}) {
	clog.Warn(msg, args...)
}

func Error(msg string, args ...interface{}) {
	clog.Error(msg, args...)
}

func DPanic(msg string, args ...interface{}) {
	clog.DPanic(msg, args...)
}

func Panic(msg string, args ...interface{}) {
	clog.Panic(msg, args...)
}

func Fatal(msg string, args ...interface{}) {
	clog.Fatal(msg, args...)
}

func Stack(err interface{}) {
	er := errors.Errorf("%v", err)
	clog.Error("%+v", er)
}
