package main

import (
	"github.com/micro-plat/hydra"
	"github.com/micro-plat/hydra/conf/server/mqc"
	"github.com/micro-plat/hydra/conf/vars/queue/lmq"
)

func init() {
	hydra.Conf.OnReady(func() {
		hydra.Conf.MQC("lmq://queue", mqc.WitchMasterSlave())
		hydra.Conf.Vars().Queue("queue", lmq.New())
	})
}
