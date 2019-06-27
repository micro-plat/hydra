package main

import (
	"github.com/micro-plat/hydra/hydra"
)

func main() {
	app := hydra.NewApp(
		hydra.WithPlatName("mall"),
		hydra.WithSystemName("apiserver02"),
		hydra.WithServerTypes("api"),
		hydra.WithDebug())

	app.Start()
}
