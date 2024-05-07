package config

import (
	"github.com/joho/godotenv"
)

func init() {
	paths := []string{
		"./.env", "../.env",
	}
	for i := range paths {
		err := godotenv.Load(paths[i])
		if err == nil {
			GetAppConfig().ConfigPath = paths[i]
			break
		}
	}
}
