package main

import (
	"github.com/micro-plat/hydra/hydra"
)

func main() {
	app := hydra.NewApp(hydra.WithPlatName("hydrav4"),
		hydra.WithSystemName("collector"),
		hydra.WithServerTypes("api"),
		hydra.WithDebug())
	app.Start()
}
