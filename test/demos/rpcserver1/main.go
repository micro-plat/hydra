package main

import (
	"github.com/micro-plat/hydra"
	_ "github.com/micro-plat/hydra/components/caches/cache/redis"
	"github.com/micro-plat/hydra/conf/server/router"
	rpcConf "github.com/micro-plat/hydra/conf/server/rpc"
	"github.com/micro-plat/hydra/hydra/servers/rpc"
)

var app = hydra.NewApp(
	hydra.WithServerTypes(rpc.RPC),
	hydra.WithPlatName("hydra"),
	hydra.WithSystemName("rpcserver"),
	hydra.WithClusterName("t"),
	hydra.WithRegistry("lm://."),
)

func init() {
	hydra.Conf.RPC(":8090", rpcConf.WithEnable())
	app.RPC("/testrpc/utf8", &rpcStruct{}, router.WithEncoding("utf8"))
	app.RPC("/testrpc/gbk", &rpcStruct{}, router.WithEncoding("gbk"))
	app.RPC("/testrpc/func", rpcFunc)
	app.RPC("/testrpc/*", rpcFunc)
	app.RPC("/testrpc/struct", rpcStruct2{})
}

func main() {
	app.Start()
}
