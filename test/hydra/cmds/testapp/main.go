package main

import (
	"github.com/micro-plat/hydra"
	varredis "github.com/micro-plat/hydra/conf/vars/redis"
	"github.com/micro-plat/hydra/hydra/servers/http"
)

var app = hydra.NewApp(
	hydra.WithServerTypes(http.API),
	hydra.WithPlatName("xxtest"),
	hydra.WithSystemName("apiserver"),
	hydra.WithClusterName("c"),
)

const OrgRedisAddr = "192.168.5.79:1000"
const ZKRegistryAddr = "zk://192.168.0.101"

func main() {
	hydra.Conf.Vars().Redis("5.79", varredis.New([]string{OrgRedisAddr}))
	hydra.Conf.API(":19003")
	app.Start()
}
