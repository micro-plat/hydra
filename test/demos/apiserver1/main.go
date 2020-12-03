package main

import (
	"github.com/micro-plat/hydra"
	_ "github.com/micro-plat/hydra/components/caches/cache/redis"
	"github.com/micro-plat/hydra/conf/server/api"
	"github.com/micro-plat/hydra/conf/server/router"
	"github.com/micro-plat/hydra/hydra/servers/http"
)

var app = hydra.NewApp(
	hydra.WithServerTypes(http.API),
	hydra.WithPlatName("httpserver"),
	hydra.WithSystemName("apiserver"),
	hydra.WithClusterName("t"),
	hydra.WithRegistry("lm://."),
)

func init() {
	hydra.Conf.API(":8070", api.WithTimeout(10, 10))
	app.API("/taosy/testapi/gbk", &apiGetgbk{}, router.WithEncoding("gbk"))
	app.API("/taosy/testapi/gb2312", &apiGetgb2312{}, router.WithEncoding("gb2312"))
	app.API("/taosy/testapi", funcAPI1)
	app.API("/taosy/testapi/*", &apiGet{})
	app.API("/taosy/testapi/sddd/:xxx", funcAPI2)
}

func main() {
	app.Start()
}
