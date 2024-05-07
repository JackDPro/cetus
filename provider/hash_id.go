package provider

import (
	"cetus/config"
	"github.com/pjebs/optimus-go"
)

func Hash() optimus.Optimus {
	conf := config.GetHashIdConfig()
	return optimus.New(conf.Prime, conf.Inverse, conf.Random)
}
