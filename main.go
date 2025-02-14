package main

import (
	"fmt"
	"github.com/JackDPro/cetus/config"
	"github.com/JackDPro/cetus/controller"
	"github.com/JackDPro/cetus/provider"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.New()
	router.Use(gin.Recovery())

	corsConf := cors.DefaultConfig()
	corsConf.AllowAllOrigins = true
	corsConf.AllowHeaders = []string{"Authorization", "Accept-Language"}
	router.Use(cors.New(corsConf))

	probeCtr := controller.NewProbeController()

	router.GET("/probe", probeCtr.Show)

	conf := config.GetApiConfig()
	addr := fmt.Sprintf("0.0.0.0:%d", conf.HttpPort)

	provider.GetLogger().Info("api start", "address", addr)
	err := router.Run(addr)
	if err != nil {
		panic(err)
	}
}
