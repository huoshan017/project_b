package client_core

import "project_b/common/log"

type Logger struct {
	log.Logger
}

var gslog *Logger

func InitLog(fileName string, maxSize, maxBackups, maxAge int, compress, consoleOutput bool, logLevel int) *Logger {
	if gslog == nil {
		gslog = &Logger {
			Logger: *log.NewWithConfig(&log.LogConfig{
				Filename:      fileName,      // "./log/client.log",
				MaxSize:       maxSize,       // 2,
				MaxBackups:    maxBackups,    // 100,
				MaxAge:        maxAge,        // 30,
				Compress:      compress,      // false,
				ConsoleOutput: consoleOutput, // true,
			}, log.DebugLevel),
		}
	}
	return gslog
}

func Log() *Logger {
	return gslog
}
