package main

import (
	"github.com/micro-plat/hydra"
	"github.com/micro-plat/hydra/conf/server/cron"
)

func init() {
	hydra.OnReady(func() {
		hydra.Conf.CRON(cron.WithP2P(), cron.WithTrace())
	})
}
