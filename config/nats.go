package config

import (
	"os"
	"sync"
)

type NatsConf struct {
	Host     string
	Username string
	Password string
}

var natsConfInstance *NatsConf
var natsConfOnce sync.Once

func GetNatsConf() *NatsConf {
	natsConfOnce.Do(func() {
		natsConfInstance = &NatsConf{
			Host:     os.Getenv("NATS_HOST"),
			Username: os.Getenv("NATS_USERNAME"),
			Password: os.Getenv("NATS_PASSWORD"),
		}
	})
	return natsConfInstance
}
