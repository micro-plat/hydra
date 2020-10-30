package main

import (
	"github.com/micro-plat/hydra"

	"github.com/micro-plat/hydra/conf/server/api"
	ccron "github.com/micro-plat/hydra/conf/server/cron"
	"github.com/micro-plat/hydra/conf/server/queue"
	crpc "github.com/micro-plat/hydra/conf/server/rpc"
	"github.com/micro-plat/hydra/conf/server/static"
	"github.com/micro-plat/hydra/conf/server/task"
	vqueue "github.com/micro-plat/hydra/conf/vars/queue"
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/hydra/servers/cron"
	"github.com/micro-plat/hydra/hydra/servers/http"
	"github.com/micro-plat/hydra/hydra/servers/mqc"
	"github.com/micro-plat/hydra/hydra/servers/rpc"

	_ "github.com/micro-plat/hydra/components/pkgs/mq/lmq"
)

var app = hydra.NewApp(
	hydra.WithServerTypes(http.API, rpc.RPC, cron.CRON, mqc.MQC),
	hydra.WithPlatName("taosytest"),
	hydra.WithSystemName("test-confcache"),
	hydra.WithClusterName("taosy1"),
	hydra.WithRegistry("zk://192.168.0.101"),
)

func init() {
	hydra.Conf.API(":8070", api.WithTimeout(10, 10)).Static(static.WithArchive("taosy-test"))
	hydra.Conf.RPC(":8888", crpc.WithTimeout(10, 10))
	hydra.Conf.CRON(ccron.WithTimeout(10)).Task(task.NewTask("@every 10s", "/taosy/testcron"))
	hydra.Conf.MQC("lmq://mqtest").Queue(queue.NewQueue("queue1", "/service1"))
	hydra.Conf.Vars().Queue("mqtest", vqueue.New("lmq://", []byte{}))
	app.API("/taosy/testapi", func(ctx context.IContext) (r interface{}) {
		ctx.Log().Info("api 接口服务测试")
		return "success"
	})
	app.RPC("/taosy/testrpc", func(ctx context.IContext) (r interface{}) {
		ctx.Log().Info("rpc 接口服务测试")
		return "success"
	})
	app.CRON("/taosy/testcron", func(ctx context.IContext) (r interface{}) {
		ctx.Log().Info("cron 接口服务测试")
		return "success"
	})

	app.MQC("/service1", func(ctx context.IContext) (r interface{}) {
		ctx.Log().Info("mqc 接口服务测试")
		return "success"
	})
}

func main() {
	app.Start()
}
