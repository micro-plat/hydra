package main

import (
	"time"

	"github.com/micro-plat/hydra"
	_ "github.com/micro-plat/hydra/components/pkgs/mq/redis"
	"github.com/micro-plat/hydra/conf/server/api"
	"github.com/micro-plat/hydra/conf/server/apm"
	"github.com/micro-plat/hydra/conf/server/header"
	"github.com/micro-plat/hydra/conf/server/mqc"
	squeue "github.com/micro-plat/hydra/conf/server/queue"
	"github.com/micro-plat/hydra/conf/vars/queue"
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/hydra/servers/http"
	smqc "github.com/micro-plat/hydra/hydra/servers/mqc"
)

/*
	1.安装时监听队列列表不能发布到注册中心;
	2.var的queue节点,不能发布json数据到注册中心;
	3.主即节点timeout配置无效;
*/

var app = hydra.NewApp(
	hydra.WithServerTypes(http.API, smqc.MQC),
	hydra.WithPlatName("taosy-mqc-test"),
	hydra.WithSystemName("test-mqc"),
	hydra.WithClusterName("taosymqc"),
	hydra.WithRegistry("redis://192.168.5.79:6379"),
	// hydra.WithRegistry("redis://192.168.0.111:6379,192.168.0.112:6379,192.168.0.113:6379,192.168.0.114:6379,192.168.0.115:6379,192.168.0.116:6379"),
	// hydra.WithRegistry("zk://192.168.0.101:2181"),
)

func init() {

	hydra.Conf.Vars().Queue("queue", queue.New("redis", []byte("")))
	queues := &squeue.Queues{}
	queues = queues.Append(squeue.NewQueue("taosy-mqc-test:queuename1", "/taosy/testmqc"))
	mqser := hydra.Conf.MQC("redis://queue", mqc.WithTrace(), mqc.WithTimeout(10))
	// mqser.Sub("server", `{"proto":"redis","addrs":["192.168.5.79:6379"],"db":0,"dial_timeout":10,"read_timeout":10,"write_time":10,"pool_size":10}`)
	mqser.Queue(queues.Queues...)

	hydra.Conf.API(":8080", api.WithTimeout(10, 10), api.WithEnable()).APM("skywalking", apm.WithDisable()).Header(header.WithCrossDomain())
	app.API("/taosy/testapi", func(ctx context.IContext) (r interface{}) {
		ctx.Log().Info("api 接口服务测试")
		queueObj := hydra.C.Queue().GetRegularQueue("redis")
		if err := queueObj.Push("queuename1", `{"taosy":"123456"}`); err != nil {
			ctx.Log().Error("发送队列报错")
			return
		}
		return nil
	})

	app.MQC("/taosy/testmqc", func(ctx context.IContext) (r interface{}) {
		ctx.Log().Info("mqc-----接口服务测试")
		time.Sleep(time.Second * 15)
		ctx.Log().Info("---------------:", ctx.Request().GetString("taosy"))
		return nil
	})
}

func main() {
	app.Start()
}
