package main

import (
	"fmt"
	"time"

	"github.com/micro-plat/hydra"
	_ "github.com/micro-plat/hydra/components/caches/cache/redis"
	"github.com/micro-plat/hydra/conf/server/acl/blacklist"
	"github.com/micro-plat/hydra/conf/server/api"
	ccron "github.com/micro-plat/hydra/conf/server/cron"
	"github.com/micro-plat/hydra/conf/server/queue"
	crpc "github.com/micro-plat/hydra/conf/server/rpc"
	"github.com/micro-plat/hydra/conf/server/static"
	"github.com/micro-plat/hydra/conf/server/task"
	cacheredis "github.com/micro-plat/hydra/conf/vars/cache/redis"
	queueredis "github.com/micro-plat/hydra/conf/vars/queue/redis"

	varredis "github.com/micro-plat/hydra/conf/vars/redis"

	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/hydra/servers/cron"
	"github.com/micro-plat/hydra/hydra/servers/http"
	"github.com/micro-plat/hydra/hydra/servers/mqc"
	"github.com/micro-plat/hydra/hydra/servers/rpc"
)

var app = hydra.NewApp(
	hydra.WithServerTypes(http.API, http.Web, rpc.RPC, cron.CRON, mqc.MQC),
	hydra.WithPlatName("xxtest"),
	hydra.WithSystemName("apiserver"),
	hydra.WithClusterName("c"),
	hydra.WithRegistry("zk://192.168.0.101"),
	//hydra.WithRegistry("lm://."),
)

func init() {
	hydra.Conf.Vars().Redis("5.79", varredis.New([]string{"192.168.5.79:6379"}))
	hydra.Conf.Vars().Cache().Redis("xxx", cacheredis.New(cacheredis.WithConfigName("5.79")))
	//hydra.Conf.Vars().Queue().Redis("xxx", queueredis.New(queueredis.WithAddrs("192.168.5.79:6379")))
	hydra.Conf.Vars().Queue().Redis("xxx", queueredis.New(queueredis.WithConfigName("5.79")))

	hydra.Conf.Web(":8071").Static(static.WithArchive("taosy-test"))
	hydra.Conf.API(":8070", api.WithTimeout(10, 10)).BlackList(blacklist.WithEnable(), blacklist.WithIP("192.168.0.101"))
	hydra.Conf.RPC(":8888", crpc.WithTimeout(10, 10))
	hydra.Conf.CRON(ccron.WithTimeout(10)).Task(task.NewTask("@every 3s", "/taosy/testcron"))
	hydra.Conf.MQC("redis://xxx").Queue(queue.NewQueue("queue1", "/service1")).Queue(queue.NewQueue("queue2", "/service1"))
	app.API("/taosy/testapi", func(ctx context.IContext) (r interface{}) {
		ctx.Log().Info("api 接口服务测试")

		key := "test:1:2:3"
		cache, err := hydra.C.Cache().GetCache("xxx")
		fmt.Println(err)
		err1 := cache.Set(key, "111", -1)
		fmt.Println(err1)
		val, err2 := cache.Get(key)
		fmt.Println(val, err2)
		return "success"
	})
	app.RPC("/taosy/testrpc", func(ctx context.IContext) (r interface{}) {
		ctx.Log().Info("rpc 接口服务测试")
		return "success"
	})
	app.CRON("/taosy/testcron", func(ctx context.IContext) (r interface{}) {
		ctx.Log().Info("cron 接口服务测试")

		q, err := hydra.C.Queue().GetQueue("xxx")
		ctx.Log().Error("C.Queue().GetQueue:", err)
		q.Push("queue1", fmt.Sprintf(`{"mqcv":"%s"}`, time.Now().Format("20060102150405")))

		return "success"
	})

	app.MQC("/service1", func(ctx context.IContext) (r interface{}) {
		ctx.Log().Info("mqc 接口服务测试,", ctx.Request().GetString("mqcv"))
		return "success"
	})
	app.RPC("/rpcs/sss", &sss{})
}

func main() {

	app.Start()

	// serverMaps := cmap.New(1)
	// serverMaps.Set("1", "2")
	// fmt.Println(len(serverMaps))
	// serverMaps.Set("2", "2")
	// fmt.Println(len(serverMaps))
}

type sss struct{}

func (s *sss) Handle(ctx context.IContext) (r interface{}) {
	ctx.Log().Info("mqc 接口服务测试,", ctx.Request().GetString("mqcv"))
	return "success"
}

func (s *sss) QueryHandle(ctx context.IContext) (r interface{}) {
	ctx.Log().Info("mqc 接口服务测试,", ctx.Request().GetString("mqcv"))
	return "success"
}
