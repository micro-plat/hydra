package tests

import (
	"os"
	"time"

	"github.com/micro-plat/hydra"
	c "github.com/micro-plat/hydra/conf/server/cron"
	"github.com/micro-plat/hydra/hydra/servers/cron"
	"github.com/micro-plat/hydra/hydra/servers/http"
)

func startServer() {
	app := hydra.NewApp(
		hydra.WithPlatName("hydratest"),
		hydra.WithSystemName("test"),
		hydra.WithServerTypes(http.API, cron.CRON),
		hydra.WithDebug(),
		hydra.WithClusterName("t"),
		hydra.WithRegistry("lm://."),
	)
	hydra.Conf.API(":9091")

	app.CRON("/cron", cronTest)
	app.API("/api", api)

	os.Args = []string{"startserver", "run"}
	go app.Start()
	time.Sleep(time.Second * 2)
}

func api(ctx hydra.IContext) interface{} {
	return "success"
}
func cronTest(ctx hydra.IContext) interface{} {
	return "success"
}

func init() {
	hydra.OnReady(func() {
		hydra.Conf.CRON(c.WithMasterSlave(), c.WithTrace())
	})
}
