package config

import (
	"os"
	"sync"
)

type App struct {
	Name         string
	Debug        bool
	LogDriver    string
	LogLevel     string
	LogFormat    string
	Version      string
	DataRoot     string
	ConfigPath   string
	PublicResUrl string
}

var appInstance *App
var appOnce sync.Once

func GetAppConfig() *App {
	appOnce.Do(func() {
		debug := false
		if os.Getenv("APP_DEBUG") == "true" {
			debug = true
		}
		appInstance = &App{
			Name:         os.Getenv("APP_NAME"),
			Debug:        debug,
			LogDriver:    os.Getenv("LOG_DRIVER"),
			LogLevel:     os.Getenv("LOG_LEVEL"),
			LogFormat:    os.Getenv("LOG_FORMAT"),
			DataRoot:     os.Getenv("APP_DATA_ROOT"),
			Version:      "dev",
			ConfigPath:   "",
			PublicResUrl: os.Getenv("APP_PUBLIC_RES_URL"),
		}
	})
	return appInstance
}
