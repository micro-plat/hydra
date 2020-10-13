package rpcs

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/micro-plat/hydra"
	"github.com/micro-plat/hydra/context"

	xhttp "github.com/micro-plat/hydra/hydra/servers/http"
	"github.com/micro-plat/hydra/hydra/servers/rpc"
)

func StartRpc(t *testing.T) {

	hydra.OnReady(func() {
		hydra.Conf.API(":8070")
		hydra.Conf.RPC(":8071")
	})

	app := hydra.NewApp(
		hydra.WithServerTypes(rpc.RPC, xhttp.API),
		hydra.WithPlatName("test"),
		hydra.WithSystemName("test-rpc"),
		hydra.WithClusterName("t"),
		hydra.WithRegistry("lm://"),
	)

	app.API("/inbound", func(ctx context.IContext) (r interface{}) {
		fmt.Println("api:inbound")
		caseVal, ok := ctx.Request().Get("case")
		if !ok {
			return "err"
		}
		fmt.Println("api:inbound", caseVal)
		switch caseVal {
		case "1":
			request := hydra.C.RPC().GetRegularRPC()
			response, err := request.RequestByCtx("/getgrpc", ctx)
			if err != nil {
				return err
			}
			return response.Result

		default:

		}

		return nil
	})

	app.RPC("/getgrpc", func(ctx context.IContext) (r interface{}) {
		return "1"
	})

	app.RPC("/ttt2", func(ctx context.IContext) (r interface{}) {
		return nil
	})

	os.Args = []string{"test", "run"}
	go app.Start()
	time.Sleep(time.Second * 2)

}

func buildRequest(caseval string) (bytes []byte, err error) {

	req, err := http.NewRequest("get", "http://localhost:8070/inbound", strings.NewReader(fmt.Sprintf("case=%s", caseval)))
	if err != nil {
		return
	}
	req.Close = true

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=utf-8")

	client := http.DefaultClient

	response, err := client.Do(req)
	if response != nil {
		defer response.Body.Close()
	}
	if err != nil {
		return
	}
	bytes, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}
	return
}

func TestGetRegularRPC(t *testing.T) {
	StartRpc(t)

	bytes, err := buildRequest("1")
	if err != nil {
		t.Error(err)
	}
	if string(bytes) != "1" {
		t.Error("result", string(bytes))
	}
	time.Sleep(time.Minute)
}
