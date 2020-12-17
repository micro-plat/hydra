package main

import (
	"github.com/micro-plat/hydra"
	_ "github.com/micro-plat/hydra/components/caches/cache/gocache"
	_ "github.com/micro-plat/hydra/components/caches/cache/memcached"
	_ "github.com/micro-plat/hydra/components/caches/cache/redis"

	_ "github.com/micro-plat/hydra/components/queues/mq/lmq"
	_ "github.com/micro-plat/hydra/components/queues/mq/mqtt"
	_ "github.com/micro-plat/hydra/components/queues/mq/redis"
	_ "github.com/micro-plat/hydra/components/queues/mq/xmq"

	"github.com/micro-plat/hydra/conf/server/api"
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/hydra/servers/http"
)

var app = hydra.NewApp(
	hydra.WithServerTypes(http.API),
	hydra.WithPlatName("hydratest"),
	hydra.WithSystemName("rpctestserver"),
	hydra.WithClusterName("t"),
	hydra.WithRegistry("zk://192.168.0.101"),
)

func init() {
	hydra.Conf.Vars().RPC("rpc")
	hydra.Conf.API(":8079", api.WithTimeout(10, 10))
	app.API("/rpctest/api", func(ctx context.IContext) (r interface{}) {
		ctx.Log().Info("api-testrpc 接口服务测试")
		request, err := hydra.C.RPC().GetRPC()
		if err != nil {
			ctx.Log().Error("获取GetRPC客户端异常", err)
			return
		}
		ctx.Log().Info("RPC.IP.Result-1:")
		response, err := request.Request("/rpc@192.168.0.137:18070", map[string]interface{}{"rpc": "test"})
		if err != nil {
			ctx.Log().Error("RPC.IP.Request异常", err)
			return
		}
		ctx.Log().Info("RPC.IP.Result:", response.Result)

		ctx.Log().Info("RPC.PlatName.Result-1:")
		response, err = request.Request("/rpc/handle@192.168.0.137:18070", map[string]interface{}{"rpc": "plattest"})
		if err != nil {
			ctx.Log().Error("RPC.PlatName.Request异常", err)
			return
		}
		ctx.Log().Info("RPC.PlatName.Result:", response.Result)

		ctx.Log().Info("RPC.OLD.Result-1:")
		response, err = request.Request("/rpc/param/taosy@hydratest", map[string]interface{}{"rpc": "plattest"})
		if err != nil {
			ctx.Log().Error("RPC.OLD.Request异常", err)
			return
		}
		ctx.Log().Info("RPC.OLD.Result:", response.Result)

		return
	})
}

func main() {
	app.Start()
}
