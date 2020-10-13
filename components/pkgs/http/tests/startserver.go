package tests

import (
	"fmt"
	"os"
	"time"

	"github.com/micro-plat/hydra"
	h "github.com/micro-plat/hydra/components/pkgs/http"
	"github.com/micro-plat/hydra/hydra/servers/http"
)

func startServer() {
	app := hydra.NewApp(
		hydra.WithPlatName("hjtest"),
		hydra.WithSystemName("test"),
		hydra.WithServerTypes(http.API),
		hydra.WithDebug(),
		hydra.WithClusterName("t"),
		hydra.WithRegistry("lm://."),
	)
	hydra.Conf.API(":9091")
	app.API("/client", client)

	os.Args = []string{"startserver", "run"}
	go app.Start()
	time.Sleep(time.Second * 2)
}

func client(ctx hydra.IContext) interface{} {
	raw := ctx.Request().GetString("raw")
	o, _ := h.WithRaw([]byte(raw))
	fmt.Println("RAW:", raw)
	_, err := h.NewClient(o)
	if (err != nil) != false {
		return fmt.Errorf("NewClient() error = %v, wantErr %v", err, false)
	}
	return "success"
}
