package main

import (
	"github.com/micro-plat/hydra"
	_ "github.com/micro-plat/hydra/components/caches/cache/redis"
	"github.com/micro-plat/hydra/conf/server/api"
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/hydra/servers/http"
)

var app = hydra.NewApp(
	hydra.WithServerTypes(http.API),
	hydra.WithPlatName("xxtest"),
	hydra.WithSystemName("apiserver"),
	hydra.WithClusterName("t"),
	hydra.WithRegistry("zk://192.168.0.101"),
)

func init() {
	hydra.Conf.Web(":8071", api.WithTimeout(10, 10)).Static()
	hydra.Conf.API(":8070", api.WithTimeout(10, 10))
	app.API("/taosy/testapi", &apiGet{})
	app.API("/taosy/testapi", func(ctx context.IContext) (r interface{}) {
		ctx.Log().Info("api-all 接口服务测试")
		return "success"
	})
	app.API("/taosy/testapi/*", &apiGet{})
	app.API("/taosy/testapi/sddd/:xxx", func(ctx context.IContext) (r interface{}) {
		ctx.Log().Info("api-all 接口服务测试")
		return "success"
	})
}

func main() {
	app.Start()
}

type apiGet struct{}

func (s *apiGet) GetHandle(ctx context.IContext) (r interface{}) {
	ctx.Log().Info("api-get 接口服务测试")
	return "api-get-success"
}
