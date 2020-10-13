package main

import (
	"github.com/micro-plat/hydra"
	_ "github.com/micro-plat/hydra/components/pkgs/mq/redis"
	"github.com/micro-plat/hydra/conf/server/api"
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/hydra/servers/http"
)

/*
	1.服务设置为stop,还是能够正常启动服务,设置无效;
	2.WithHost的设置没有效果;
	3.DNS配置也不知道有什么作用;
	4.路由注册冲突,/taosy/testapi/:name和/taosy/testapi/dsfd/dddd启动服务会程序崩溃;
*/

var app = hydra.NewApp(
	hydra.WithServerTypes(http.API),
	hydra.WithPlatName("taosy-api-test"),
	hydra.WithSystemName("test-api"),
	hydra.WithClusterName("taosyapi"),
	hydra.WithRegistry("redis://192.168.5.79:6379"),
	// hydra.WithRegistry("redis://192.168.0.111:6379,192.168.0.112:6379,192.168.0.113:6379,192.168.0.114:6379,192.168.0.115:6379,192.168.0.116:6379"),
	// hydra.WithRegistry("zk://192.168.0.101:2181"),
)

func init() {
	app.OnHandleExecuting(func(ctx context.IContext) (r interface{}) {
		ctx.Log().Info("执行前处理")
		ctx.Log().Info("执行路径:", ctx.Request().Path().GetRouter().Service)

		return
	})

	serverApi := hydra.Conf.API(":8070", api.WithTimeout(10, 10), api.WithDisable(), api.WithTrace(), api.WithHost("192.168.5.107"), api.WithDNS(""))
	//serverApi.APM("skywalking", apm.WithDisable())
	app.API("/taosy/testapi/:name", func(ctx context.IContext) (r interface{}) {
		ctx.Log().Info("api 接口服务测试")
		ctx.Log().Info("name :", ctx.Request().Param("name"))
		return "success"
	})

	app.API("/taosy/testapi", func(ctx context.IContext) (r interface{}) {
		ctx.Log().Info("api 接口服务测试")
		ctx.Log().Info("name :", ctx.Request().Param("name"))
		return "success"
	})

	app.API("/taosy/jwttest", func(ctx context.IContext) (r interface{}) {
		ctx.Log().Info("jwt 接口服务测试")
		return "success"
	})
}

func main() {
	app.Start()
}
