package main

import (
	"github.com/micro-plat/hydra"
	"github.com/micro-plat/hydra/conf/server/api"
	"github.com/micro-plat/hydra/conf/server/apm"
	varapm "github.com/micro-plat/hydra/conf/vars/apm"
)

func init() {
	hydra.OnReady(func() {
		hydra.Conf.API(":8081", api.WithHeaderReadTimeout(30), api.WithTimeout(30, 30)).
			APM("skywalking", apm.WithEnable(), apm.WithDB("db", "sup17", "common"), apm.WithCache("cache", "redis", "mem"))
		//hydra.Conf.Vars().RLog("/rpc/log@hydra", rlog.WithAll())
		hydra.Conf.Vars().APM(varapm.New("skywalking", []byte(`
		{
			"server_address":"192.168.106.160:11800",
			"instance_props": {"x": "1", "y": "2"}
		}`),
		))
	})
}
