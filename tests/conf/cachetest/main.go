package main

import (
	"time"

	"github.com/micro-plat/hydra"
	"github.com/micro-plat/hydra/conf/server"
	"github.com/micro-plat/hydra/global"
)

func main() {
	app := hydra.NewApp(
		hydra.WithServerTypes(global.API),
		hydra.WithClusterName("taosy"),
		hydra.WithDebug(),
		hydra.WithPlatName("test"),
		hydra.WithSystemName("apiserver"),
		hydra.WithRegistry("fs://../"),
	)

	app.API("/conf/cache/:testtype", request)
	app.Start()
}

func request(ctx hydra.IContext) (r interface{}) {
	ctx.Log().Info("------配置文件缓存测试------")
	switch ctx.Request().Param("testtype") {
	case "1":
		serverConf, err := server.Cache.GetServerConf(global.API)
		if err != nil {
			return
		}
		router := serverConf.GetRouterConf()
		ctx.Log().Infof("api-router信息:%s", router.String())
		varConf, err := server.Cache.GetVarConf()
		if err != nil {
			return err
		}

		confInfo, err := varConf.GetConf("taosytest", "db")
		if err != nil {
			return err
		}

		b, v, err := confInfo.GetJSON("db")
		if err != nil {
			return err
		}
		ctx.Log().Infof("db节点数据:%s \n", string(b))
		ctx.Log().Infof("db配置版本号:%d \n", v)
		return
	case "2":
		return 100
	case "3":
		return time.Now().String()
	case "4":
		return float32(100.20)
	case "5":
		return true
	case "6":
	default:
	}
	return
}
