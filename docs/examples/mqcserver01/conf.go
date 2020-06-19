package main

import (
	"github.com/micro-plat/hydra"
	"github.com/micro-plat/hydra/conf/server/mqc"
	"github.com/micro-plat/hydra/conf/server/queue"
	"github.com/micro-plat/hydra/conf/vars/queue/lmq"
)

func init() {
	hydra.Conf.MQC("lmq://queue", mqc.WitchMasterSlave()).Queue(queue.NewQueue("order.query", "/order/request"))
	hydra.Conf.Vars().Queue("queue", lmq.New())
}
