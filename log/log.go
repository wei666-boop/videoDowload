package log

import (
	"log"
	"os"
)

var (
	INFO    = 1
	WARNING = 2
	ERROR   = 3
)

func GetLog(path string) *log.Logger {
	logFile, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		panic("无法创建日志" + err.Error())
	}
	logger := log.New(logFile, "[下载]", log.Llongfile|log.Lmicroseconds|log.Ldate)
	return logger
}
func WriteLog(level int, msg interface{}, logger *log.Logger) {
	switch level {
	case INFO:
		logger.Println("[info]", msg)
	case WARNING:
		logger.Println("[warning]", msg)
	case ERROR:
		logger.Println("[error]", msg)
	}
}
