package main

import (
	"github.com/micro-plat/hydra"
	"github.com/micro-plat/hydra/components/pkgs/apm/apmtypes"
	"github.com/micro-plat/hydra/conf/server/apm"
	"github.com/micro-plat/hydra/conf/vars/rlog"
)

func init() {
	hydra.OnReady(func() {
		hydra.Conf.API(":8070").APM(apmtypes.SkyWalking, apm.WithDisable())
		hydra.Conf.Vars().RLog("/rpc/log", rlog.WithDisable())
		hydra.Conf.RPC(":8281").APM(apmtypes.SkyWalking, apm.WithDisable())

	})
}
