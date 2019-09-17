package main

import "github.com/micro-plat/hydra/hydra"

type flowserver struct {
	*hydra.MicroApp
}

func main() {
	app := &flowserver{
		hydra.NewApp(
			hydra.WithPlatName("mall"),
			hydra.WithSystemName("flowserver"),
			hydra.WithServerTypes("mqc")),
	}
	app.init()
	app.Start()
}
