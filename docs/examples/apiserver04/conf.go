package main

import (
	"github.com/micro-plat/hydra"
	"github.com/micro-plat/hydra/conf/vars/db/oracle"
	"github.com/micro-plat/hydra/conf/vars/queue/lmq"
)

func init() {
	hydra.OnReady(func() {
		hydra.Conf.API(":8080").Metric("http://192.168.106.219:8086", "hydra", "@every 5s")
		hydra.Conf.Vars().DB("db", oracle.New("hydra/hydra")).Queue("queue", lmq.New())
	})
}
