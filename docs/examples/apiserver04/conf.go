package main

import (
	"github.com/micro-plat/hydra"
	"github.com/micro-plat/hydra/conf/server/apm"
)

func init() {
	hydra.OnReady(func() {
		hydra.Conf.API(":8084").APM(apm.WithEnable())
		//hydra.Conf.API(":8080").Metric("http://192.168.106.219:8086", "hydra", "@every 5s")
		//hydra.Conf.Vars().DB("db", oracle.New("hydra/hydra")).Queue("queue", lmq.New())
	})
}
