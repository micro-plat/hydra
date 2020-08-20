package main

import (
	"github.com/micro-plat/hydra"
	"github.com/micro-plat/hydra/conf/server/api"
	"github.com/micro-plat/hydra/conf/server/apm"
	"github.com/micro-plat/lib4go/types"
)

func init() {
	hydra.OnReady(func() {
		hydra.Conf.API(":8081", api.WithTrace()).APM(apm.WithEnable())
		//hydra.Conf.Vars().RLog("/rpc/log@hydra", rlog.WithAll())
		hydra.Conf.Vars()["apm"] = types.XMap{
			"apm": types.XMap{
				"apmtype":             "skywalking",
				"check_interval":      1,
				"max_send_queue_size": 500000,
				"instance_props":      map[string]string{"x": "1", "y": "2"},
			},
		}
	})
}
