package main

import "github.com/micro-plat/hydra/hydra"

type rpcserver struct {
	*hydra.MicroApp
}

func main() {
	app := &rpcserver{
		hydra.NewApp(
			hydra.WithPlatName("mall"),
			hydra.WithSystemName("rpcserver"),
			hydra.WithServerTypes("rpc")),
	}
	app.init()
	app.Start()
}
