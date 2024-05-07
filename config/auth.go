package config

import (
	"os"
	"sync"
)

type AuthConf struct {
	CertPath string
	KeyPath  string
	Issue    string
}

var authConfInstance *AuthConf
var authConfOnce sync.Once

func GetAuthConf() *AuthConf {
	authConfOnce.Do(func() {
		authConfInstance = &AuthConf{
			CertPath: os.Getenv("JWT_CERT_PATH"),
			KeyPath:  os.Getenv("JWT_KEY_PATH"),
			Issue:    os.Getenv("JWT_ISSUE"),
		}
	})
	return authConfInstance
}
