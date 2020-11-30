package main

import (
	"fmt"
	"reflect"
	"time"

	"github.com/micro-plat/hydra"
	"github.com/micro-plat/hydra/conf/vars/queue/lmq"
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/hydra/servers/http"
	"github.com/micro-plat/hydra/hydra/servers/mqc"
)

var app = hydra.NewApp(
	hydra.WithServerTypes(http.API, mqc.MQC),
	hydra.WithPlatName("hydra-examples"),
	hydra.WithSystemName("mqcserver02"),
	hydra.WithClusterName("t1"))

func main() {

	o := objService{}
	fmt.Println(reflect.TypeOf(o).NumField())
	fmt.Println(reflect.TypeOf(o).NumMethod())
	fmt.Println(reflect.TypeOf(o).Kind())

	const redisServer = "192.168.0.101"

	hydra.Conf.Vars().Queue().LMQ("lmqqueue", lmq.New())
	hydra.Conf.MQC("redis://redisqueue")
	hydra.Conf.API(":59002")

	//app.MQC("/testmqc/proc1", &objService{}, "mqc2:service:queue1")
	app.MQC("/testmqc/proc2", objService{}, "mqc2:service:queue2")
	app.MQC("/testmqc/proc3", NewObjNoneError, "mqc2:service:queue3")
	app.MQC("/testmqc/proc4", NewObjWithError, "mqc2:service:queue4")
	app.MQC("/testmqc/proc5", FuncService, "mqc2:service:queue5")

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
