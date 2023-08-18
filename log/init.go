package log

import (
	"github.com/pkg/errors"
	"go.uber.org/zap/zapcore"
)

var log *Logger

func InitLog(fileName string, maxSize, maxBackups, maxAge int, compress, consoleOutput bool, logLevel int) *Logger {
	log = NewWithConfig(&LogConfig{
		Filename:      fileName,      // "./log/client.log",
		MaxSize:       maxSize,       // 2,
		MaxBackups:    maxBackups,    // 100,
		MaxAge:        maxAge,        // 30,
		Compress:      compress,      // false,
		ConsoleOutput: consoleOutput, // true,
	}, zapcore.Level(logLevel))
	return log
}

func Debug(msg string, args ...interface{}) {
	log.Debug(msg, args...)
}

func Info(msg string, args ...interface{}) {
	log.Info(msg, args...)
}

func Warn(msg string, args ...interface{}) {
	log.Warn(msg, args...)
}

func Error(msg string, args ...interface{}) {
	log.Error(msg, args...)
}

func DPanic(msg string, args ...interface{}) {
	log.DPanic(msg, args...)
}

func Panic(msg string, args ...interface{}) {
	log.Panic(msg, args...)
}

func Fatal(msg string, args ...interface{}) {
	log.Fatal(msg, args...)
}

func Stack(err interface{}) {
	er := errors.Errorf("%v", err)
	log.Error("%+v", er)
}
