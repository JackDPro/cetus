package config

import (
	"os"
	"sync"
)

type App struct {
	Name         string
	Env          string
	Version      string
	DataRoot     string
	ConfigPath   string
	PublicResUrl string
}

var appInstance *App
var appOnce sync.Once

func GetAppConfig() *App {
	appOnce.Do(func() {
		appInstance = &App{
			Name:         os.Getenv("APP_NAME"),
			Env:          os.Getenv("APP_ENV"),
			DataRoot:     os.Getenv("APP_DATA_ROOT"),
			Version:      "dev",
			ConfigPath:   "",
			PublicResUrl: os.Getenv("APP_PUBLIC_RES_URL"),
		}
	})
	return appInstance
}
