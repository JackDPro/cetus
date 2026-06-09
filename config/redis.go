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
	Prefix   string
}

// Key 给 key 添加全局前缀，前缀为空时返回原 key（不带分隔符）。
// 用于多个项目共用同一个 redis 库时做命名空间隔离。
func (r *Redis) Key(key string) string {
	if r.Prefix == "" {
		return key
	}
	return r.Prefix + ":" + key
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
		// 全局前缀，可不填写，不填写时所有 key 不带前缀
		prefix := os.Getenv("REDIS_PREFIX")

		RedisConfigInstance = &Redis{
			Host:     host,
			Password: password,
			Port:     port,
			DB:       database,
			PoolSize: 10,
			Prefix:   prefix,
		}
	})
	return RedisConfigInstance
}
