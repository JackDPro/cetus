package provider

import (
	"github.com/JackDPro/cetus/config"
	"github.com/redis/go-redis/v9"
	"sync"
)

var redisInstance *redis.Client
var redisOnce sync.Once

func GetRedisClient() *redis.Client {
	redisOnce.Do(func() {
		conf := config.GetRedisConfig()
		redisInstance = redis.NewClient(&redis.Options{
			Addr:         conf.Host + ":" + conf.Port,
			Password:     conf.Password,
			DB:           conf.DB,
			PoolSize:     conf.PoolSize,
			MinIdleConns: 0,
			PoolTimeout:  0,
			TLSConfig:    nil,
			Limiter:      nil,
		})
	})
	return redisInstance
}
