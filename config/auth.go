package config

import (
	"os"
	"strconv"
	"sync"
)

type AuthConf struct {
	CertPath    string
	KeyPath     string
	Issue       string
	Audience    string
	ExpiresIn   int
	RedisPrefix string
}

var authConfInstance *AuthConf
var authConfOnce sync.Once

func GetAuthConf() *AuthConf {
	authConfOnce.Do(func() {
		expiresIn := 72
		if os.Getenv("JWT_EXPIRES_IN") != "" {
			value, err := strconv.Atoi(os.Getenv("JWT_EXPIRES_IN"))
			if err == nil {
				expiresIn = value
			}
		}
		redisPrefix := "tokens"
		if os.Getenv("JWT_REDIS_PREFIX") != "" {
			redisPrefix = os.Getenv("JWT_REDIS_PREFIX")
		}
		authConfInstance = &AuthConf{
			CertPath:    os.Getenv("JWT_CERT_PATH"),
			KeyPath:     os.Getenv("JWT_KEY_PATH"),
			Issue:       os.Getenv("JWT_ISSUE"),
			Audience:    os.Getenv("JWT_AUDIENCE"),
			RedisPrefix: redisPrefix,
			ExpiresIn:   expiresIn,
		}
	})
	return authConfInstance
}
