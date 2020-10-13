package main

import (
	"fmt"

	"github.com/micro-plat/hydra"
	_ "github.com/micro-plat/hydra/components/pkgs/mq/redis"
	rpcc "github.com/micro-plat/hydra/components/rpcs/rpc"
	"github.com/micro-plat/hydra/conf/server/api"
	"github.com/micro-plat/hydra/conf/server/cron"
	"github.com/micro-plat/hydra/conf/server/mqc"
	squeue "github.com/micro-plat/hydra/conf/server/queue"
	"github.com/micro-plat/hydra/conf/server/task"
	"github.com/micro-plat/hydra/conf/vars/db"
	"github.com/micro-plat/hydra/conf/vars/db/oracle"
	"github.com/micro-plat/hydra/context"
	scron "github.com/micro-plat/hydra/hydra/servers/cron"
	"github.com/micro-plat/hydra/hydra/servers/http"
	smqc "github.com/micro-plat/hydra/hydra/servers/mqc"
	"github.com/micro-plat/hydra/hydra/servers/rpc"
)

var app = hydra.NewApp(
	hydra.WithServerTypes(http.API, rpc.RPC, smqc.MQC, scron.CRON),
	hydra.WithPlatName("taosytest"),
	hydra.WithSystemName("test-rgtredis"),
	hydra.WithClusterName("taosy"),
	hydra.WithRegistry("redis://192.168.5.79:6379"),
	// hydra.WithRegistry("redis://192.168.0.111:6379,192.168.0.112:6379,192.168.0.113:6379,192.168.0.114:6379,192.168.0.115:6379,192.168.0.116:6379"),
	// hydra.WithRegistry("zk://192.168.0.101:2181"),
)

func init() {

	hydra.Conf.Vars().DB("taosy_db", oracle.New("connstring", db.WithConnect(10, 10, 10)))
	// hydra.Conf.Vars().Cache("cache", credis.New("192.168.5.79:6379", credis.WithDbIndex(0), credis.WithPoolSize(10), credis.WithTimeout(10, 10, 10)))
	// hydra.Conf.Vars().Queue("queue", qredis.New("192.168.5.79:6379", qredis.WithDbIndex(0), qredis.WithPoolSize(10), qredis.WithTimeout(10, 10, 10)))
	hydra.Conf.RPC(":8071")
	queues := &squeue.Queues{}
	queues = queues.Append(squeue.NewQueue("queuename1", "/testmqc"))
	mqser := hydra.Conf.MQC("redis://queue", mqc.WithTrace(), mqc.WithTimeout(10))
	// mqser.Sub("server", `{"proto":"redis","addrs":["192.168.5.79:6379"],"db":0,"dial_timeout":10,"read_timeout":10,"write_time":10,"pool_size":10}`)
	mqser.Queue(queues.Queues...)
	tasks := task.Tasks{}
	tasks.Append(task.NewTask(fmt.Sprintf("@every %ds", 10), "/testcron"))
	hydra.Conf.CRON(cron.WithEnable(), cron.WithTrace(), cron.WithTimeout(10), cron.WithSharding(1)).Task(tasks.Tasks...)
	hydra.Conf.API(":8070", api.WithTimeout(10, 10), api.WithEnable())
	app.API("/taosy/testapi", func(ctx context.IContext) (r interface{}) {
		ctx.Log().Info("api 接口服务测试")
		reqID := ctx.User().GetRequestID()
		request := hydra.C.RPC().GetRegularRPC()
		response, err := request.Request(ctx.Context(), "/taosy/testrpc", nil, rpcc.WithXRequestID(reqID))
		if err != nil {
			return err
		}
		ctx.Log().Info("rpc response.Status", response.Status)
		return nil
	})

	app.RPC("/taosy/testrpc", func(ctx context.IContext) (r interface{}) {
		ctx.Log().Info("rpc 接口服务测试")
		ctx.Log().Info("------------发送mq-----------")
		queueObj := hydra.C.Queue().GetRegularQueue()
		if err := queueObj.Push("queuename1", `{"taosy":"123456"}`); err != nil {
			ctx.Log().Error("发送队列报错")
			return
		}
		return nil
	})

	app.MQC("/testmqc", func(ctx context.IContext) (r interface{}) {
		ctx.Log().Info("mqc-----接口服务测试")
		ctx.Log().Info("---------------:", ctx.Request().GetString("taosy"))
		return nil
	})

	app.CRON("/testcron", func(ctx context.IContext) (r interface{}) {
		ctx.Log().Info("cron 接口服务测试")
		return nil
	})
}

func main() {
	app.Start()
}
