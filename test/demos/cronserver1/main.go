package main

import (
	"fmt"
	"time"

	"github.com/micro-plat/hydra"
	"github.com/micro-plat/hydra/conf/vars/queue/queueredis"
	"github.com/micro-plat/hydra/conf/vars/redis"
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/hydra/servers/cron"
	"github.com/micro-plat/hydra/hydra/servers/mqc"
)

var app = hydra.NewApp(
	hydra.WithServerTypes(cron.CRON, mqc.MQC),
	hydra.WithPlatName("hydra-examples"),
	hydra.WithSystemName("cronserver01"),
	hydra.WithClusterName("t1"))

func main() {

	const redisServer = "192.168.0.101"

	hydra.Conf.Vars().Redis("0.101", redis.New(nil, redis.WithAddrs(redisServer)))
	hydra.Conf.Vars().Queue().Redis("mqcqueue", queueredis.New(queueredis.WithConfigName("0.101")))
	hydra.Conf.MQC("redis://mqcqueue")
	hydra.Conf.API(":59001")

	app.CRON("/testcron/proc1", &cronService{}, "@every 10s")
	app.CRON("/testcron/proc2", func(ctx context.IContext) interface{} {
		fmt.Println("ONCE---")
		return nil
	}, "@once")
	app.CRON("/testcron/proc3", func(ctx context.IContext) interface{} {
		fmt.Println("Now")
		return nil
	}, "@now")

	app.CRON("/testcron/proc4", func(ctx context.IContext) interface{} {
		fmt.Println("time")
		return nil
	}, "0 12 * * ?")

	app.MQC("/testmqc/proc1", &mqcService{}, "service:queue1")
	//app.MQC("/testmqc/proc2", objService{}, "service:queue2")
	app.MQC("/testmqc/proc3", NewObjNoneError, "service:queue3")
	app.MQC("/testmqc/proc4", NewObjWithError, "service:queue4")
	app.MQC("/testmqc/proc5", FuncService, "service:queue5")

	app.Micro("/order/request", request)

	app.Start()
}

func request(ctx context.IContext) (r interface{}) {
	queueName := ctx.Request().GetString("queuename")

	queueClient := hydra.C.Queue().GetRegularQueue()
	data := fmt.Sprintf(`{"name":"%d"}`, time.Now().Unix())
	queueClient.Send(fmt.Sprintf("service:%s", queueName), data)

	return data
}
