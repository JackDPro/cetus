package config

import (
	"log"
	"os"
	"strconv"
	"sync"
)

type Redis struct {
	Host     string
	Password string
	Port     string
	DB       int
	PoolSize int
}

var RedisConfigInstance *Redis
var RedisConfigOnce sync.Once

func GetRedisConfig() *Redis {
	RedisConfigOnce.Do(func() {
		var err error
		host := "127.0.0.1"
		if os.Getenv("REDIS_HOST") != "" {
			host = os.Getenv("REDIS_HOST")
		}
		port := "6379"
		if os.Getenv("REDIS_PORT") != "" {
			port = os.Getenv("REDIS_PORT")
		}
		database := 0
		if os.Getenv("REDIS_DATABASE") != "" {
			database, err = strconv.Atoi(os.Getenv("REDIS_DATABASE"))
			if err != nil {
				database = 0
				log.Printf("[error] redis database is not int \n")
			}
		}
		password := "password"
		if os.Getenv("REDIS_PASSWORD") != "" {
			password = os.Getenv("REDIS_PASSWORD")
		}

		RedisConfigInstance = &Redis{
			Host:     host,
			Password: password,
			Port:     port,
			DB:       database,
			PoolSize: 10,
		}
	})
	return RedisConfigInstance
}
