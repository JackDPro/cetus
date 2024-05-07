package provider

import (
	"cetus/config"
	"fmt"
	kitLog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"log"
	"os"
	"sync"
)

var loggerInstance kitLog.Logger
var loggerOnce sync.Once

func GetLogger() kitLog.Logger {
	loggerOnce.Do(func() {
		conf := config.GetAppConfig()
		writer := os.Stderr
		var err error
		if conf.LogDriver == "file" && conf.DataRoot != "" {
			filePath := fmt.Sprintf("%s/logs/api.log", conf.DataRoot)
			writer, err = os.OpenFile(filePath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
			if err != nil {
				log.Println("open log file failed", err)
				writer = os.Stderr
			}
		}
		loggerInstance = kitLog.NewLogfmtLogger(writer)
		if conf.LogFormat == "json" {
			loggerInstance = kitLog.NewJSONLogger(writer)
		}
		option := level.AllowAll()
		switch conf.LogLevel {
		case "debug":
			option = level.AllowDebug()
		case "info":
			option = level.AllowInfo()
		case "warning":
			option = level.AllowWarn()
		case "error":
			option = level.AllowError()
		case "none":
			option = level.AllowNone()
		}
		loggerInstance = level.NewFilter(loggerInstance, option)
		loggerInstance = kitLog.With(loggerInstance, "ts", kitLog.DefaultTimestamp, "line", kitLog.DefaultCaller)
	})
	return loggerInstance
}
