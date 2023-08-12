package main

import (
	common_log "project_b/common/log"

	"go.uber.org/zap/zapcore"
)

type Logger struct {
	common_log.Logger
}

var log *Logger

func InitLog(fileName string, maxSize, maxBackups, maxAge int, compress, consoleOutput bool, logLevel int) *Logger {
	if log == nil {
		log = &Logger{
			Logger: *common_log.NewWithConfig(&common_log.LogConfig{
				Filename:      fileName,      // "./log/client.log",
				MaxSize:       maxSize,       // 2,
				MaxBackups:    maxBackups,    // 100,
				MaxAge:        maxAge,        // 30,
				Compress:      compress,      // false,
				ConsoleOutput: consoleOutput, // true,
			}, zapcore.Level(logLevel)),
		}
	}
	return log
}

func Log() *Logger {
	return log
}
