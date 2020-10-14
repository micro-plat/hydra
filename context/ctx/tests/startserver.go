package tests

import (
	"os"
	"time"

	"github.com/micro-plat/hydra"
	"github.com/micro-plat/hydra/hydra/servers/http"
)

func startServer() {
	app := hydra.NewApp(
		hydra.WithPlatName("hydratest"),
		hydra.WithSystemName("test"),
		hydra.WithServerTypes(http.API),
		hydra.WithDebug(),
		hydra.WithClusterName("t"),
		hydra.WithRegistry("lm://."),
	)
	hydra.Conf.API(":9091")
	app.API("/api", api)

	os.Args = []string{"startserver", "run"}
	go app.Start()
	time.Sleep(time.Second * 2)
}

func api(ctx hydra.IContext) interface{} {
	return "success"
}
