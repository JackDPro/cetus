package controller

import (
	"errors"
	"github.com/JackDPro/cetus/config"
	"github.com/JackDPro/cetus/model"
	"github.com/gin-gonic/gin"
	"strconv"
)

type ProbeController struct {
}

func NewProbeController() *ProbeController {
	return &ProbeController{}
}

func (ctr *ProbeController) Index(c *gin.Context) {
	conf := config.GetAppConfig()
	probe1 := model.Probe{
		AppName:    conf.Name,
		AppVersion: conf.Version,
		ConfigPath: conf.ConfigPath,
	}
	probe2 := model.Probe{
		AppName:    conf.Name,
		AppVersion: conf.Version,
		ConfigPath: conf.ConfigPath,
	}
	items := []model.Probe{probe1, probe2}
	ResponseCollection(c, items, nil)
}

func (ctr *ProbeController) Show(c *gin.Context) {
	statusStr := c.Request.URL.Query().Get("status")
	status := 200
	if statusStr != "" {
		status, _ = strconv.Atoi(statusStr)
	}
	conf := config.GetAppConfig()
	probe := model.Probe{
		AppName:    conf.Name,
		AppVersion: conf.Version,
		ConfigPath: conf.ConfigPath,
	}
	switch status {
	case 400:
		ResponseBadRequest(c, 1, "request bad test")
	case 500:
		ResponseInternalError(c, 1, "request internal error test", errors.New("internal error"))
	default:
		ResponseItem(c, &probe)
	}
}
