package main

import (
	"github.com/micro-plat/hydra"
	"github.com/micro-plat/hydra/conf/server/api"
	"github.com/micro-plat/hydra/conf/vars/rlog"
)

func init() {
	hydra.OnReady(func() {
		hydra.Conf.API(":8081", api.WithTrace())
		hydra.Conf.Vars().RLog("/rpc/log@hydra", rlog.WithAll())
	})
}
