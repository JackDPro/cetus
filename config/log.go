package config

import (
	"os"
	"sync"
)

type LogConf struct {
	ConsoleOut bool
	FileOut    bool
	FilePath   string
	Level      string // debug/info/warn/error
	Format     string // json/line
}

var logInstance *LogConf
var logOnce sync.Once

func GetLogConfig() *LogConf {
	logOnce.Do(func() {
		consoleOut := true
		if os.Getenv("LOG_CONSOLE_OUT") == "false" {
			consoleOut = false
		}
		fileOut := false
		if os.Getenv("LOG_FILE_OUT") == "true" && os.Getenv("LOG_FILE_PATH") != "" {
			fileOut = true
		}

		logInstance = &LogConf{
			ConsoleOut: consoleOut,
			FileOut:    fileOut,
			FilePath:   os.Getenv("LOG_FILE_PATH"),
			Level:      os.Getenv("LOG_LEVEL"),
			Format:     os.Getenv("LOG_FORMAT"),
		}
	})
	return logInstance
}
