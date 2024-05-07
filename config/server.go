package config

import (
	"os"
	"strconv"
	"sync"
)

type ServerConf struct {
	HttpPort uint64
	GrpcPort uint64
}

var serverConfigInstance *ServerConf
var serverConfigOnce sync.Once

func GetApiConfig() *ServerConf {
	serverConfigOnce.Do(func() {
		var httpPort uint64 = 80
		var grpcPort uint64 = 50051
		if os.Getenv("SERVER_HTTP_PORT") != "" {
			httpPort, _ = strconv.ParseUint(os.Getenv("SERVER_HTTP_PORT"), 10, 64)
		}
		if os.Getenv("SERVER_GRPC_PORT") != "" {
			grpcPort, _ = strconv.ParseUint(os.Getenv("SERVER_GRPC_PORT"), 10, 64)
		}
		serverConfigInstance = &ServerConf{
			HttpPort: httpPort,
			GrpcPort: grpcPort,
		}
	})
	return serverConfigInstance
}
