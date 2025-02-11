package test

import (
	"github.com/JackDPro/cetus/config"
	"github.com/JackDPro/cetus/provider"
	"testing"
)

func TestLogger(t *testing.T) {
	appConf := config.GetAppConfig()
	logger := provider.GetLogger()
	logger.Info("hello jack")
	logger.Infow("info", "foo", "bar")
	logger.Infof("info %s", "foo")
	logger.Infof("info %v", appConf)
}
