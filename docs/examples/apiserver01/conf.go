package main

import (
	"github.com/micro-plat/hydra"
	"github.com/micro-plat/hydra/conf/server/api"
)

func init() {
	hydra.Conf.OnReady(func() {
		hydra.Conf.API(":8080", api.WithTrace())
	})
}
