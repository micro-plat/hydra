package main

import (
	"fmt"
	"time"

	"github.com/micro-plat/hydra"
	"github.com/micro-plat/hydra/conf/server/queue"
	"github.com/micro-plat/hydra/conf/vars/queue/queueredis"
	"github.com/micro-plat/hydra/conf/vars/redis"
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/hydra/servers/http"
	"github.com/micro-plat/hydra/hydra/servers/mqc"
)

var app = hydra.NewApp(
	hydra.WithServerTypes(http.API, mqc.MQC),
	hydra.WithPlatName("hydra-examples"),
	hydra.WithSystemName("mqcserver01"),
	hydra.WithClusterName("t1"))

func main() {

	const redisServer = "192.168.0.101"

	hydra.Conf.Vars().Redis("0.101", redis.New(nil, redis.WithAddrs(redisServer)))
	hydra.Conf.Vars().Queue().Redis("redisqueue", queueredis.New(queueredis.WithConfigName("0.101")))
	hydra.Conf.MQC("redis://redisqueue").Queue(queue.NewQueue("service:queue1", "/test/mqc/proc1"))
	hydra.Conf.API(":59001")

	app.MQC("/testmqc/proc1", &objService{}, "service:queue1")
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
	queueClient.Push(fmt.Sprintf("service:%s", queueName), data)

	return data
}
