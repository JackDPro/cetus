package controller

import (
	"errors"
	"github.com/JackDPro/cetus/config"
	"github.com/JackDPro/cetus/model"
	"github.com/gin-gonic/gin"
	"strconv"
)

type ProbeController struct {
	Controller
}

func NewProbeController() *ProbeController {
	return &ProbeController{}
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
		ctr.ResponseBadRequest(c, 1, "request bad test")
	case 500:
		ctr.ResponseInternalError(c, 1, "request internal error test", errors.New("internal error"))
	default:
		ctr.ResponseItem(c, &probe)
	}
}
