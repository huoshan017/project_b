package log

import (
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
