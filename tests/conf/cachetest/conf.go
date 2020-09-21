package main

import (
	"fmt"

	"github.com/micro-plat/hydra"
	"github.com/micro-plat/hydra/conf/server/mqc"
	squeue "github.com/micro-plat/hydra/conf/server/queue"
	"github.com/micro-plat/hydra/conf/server/task"
	"github.com/micro-plat/hydra/conf/vars/cache"
	"github.com/micro-plat/hydra/conf/vars/db"
	"github.com/micro-plat/hydra/conf/vars/queue"
)

func init() {
	hydra.OnReady(func() {
		hydra.Conf.Vars().DB("taosy_db", db.New("oracle", "connstring", db.WithConnect(10, 10, 10)))
		hydra.Conf.Vars().Cache("taosy_cache", cache.New("redis", []byte(`{"cache":""oracle-redis}`)))
		hydra.Conf.Vars().Queue("taosy_queue", queue.New("redis", []byte(`{"queue":""oracle-redis}`)))
		hydra.Conf.API(":8070")
		hydra.Conf.RPC(":8071")
		queues := squeue.Queues{}
		queues.Append(squeue.NewQueue("queuename1", "/testmqc"))
		hydra.Conf.MQC("redis://192.168.0.111:6379", mqc.WithTrace(), mqc.WithTimeout(10)).Queue(queues.Queues...)
		tasks := task.Tasks{}
		tasks.Append(task.NewTask(fmt.Sprintf("@every %ds", 10), "/testcron"),
			task.NewTask(fmt.Sprintf("@every %ds", 10), "/testcron1"),
			task.NewTask(fmt.Sprintf("@every %ds", 10), "/testcron2"),
		)
	})
}
