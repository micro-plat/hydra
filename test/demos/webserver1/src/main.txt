package main

import (
	"github.com/micro-plat/hydra"
	_ "github.com/micro-plat/hydra/components/caches/cache/redis"
	"github.com/micro-plat/hydra/conf/server/api"
	"github.com/micro-plat/hydra/hydra/servers/http"
)

var app = hydra.NewApp(
	hydra.WithServerTypes(http.Web),
	hydra.WithPlatName("httpserver"),
	hydra.WithSystemName("webserver"),
	hydra.WithClusterName("t"),
	hydra.WithRegistry("lm://."),
)

func init() {
	hydra.Conf.Web(":8071", api.WithTimeout(10, 10)).Static()
}

func main() {
	app.Start()
}
