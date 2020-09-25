package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/micro-plat/hydra"
	"github.com/micro-plat/hydra/conf/server/api"
	"github.com/micro-plat/hydra/conf/server/apm"

	varapm "github.com/micro-plat/hydra/conf/vars/apm"
	"github.com/micro-plat/hydra/context"

	crpc "github.com/micro-plat/hydra/components/rpcs/rpc"

	"github.com/micro-plat/hydra/components/pkgs/apm/apmtypes"
	xhttp "github.com/micro-plat/hydra/hydra/servers/http"
	"github.com/micro-plat/hydra/hydra/servers/rpc"
)

func StartServer() {

	hydra.Conf.Vars().APM(varapm.New(apmtypes.SkyWalking, []byte(`
	{
		"server_address":"192.168.106.160:11800",
		"instance_props": {"x": "1", "y": "2"}
	}`),
	))
	//hydra.Conf.Vars().APM(varapm.New("skywalking", `{"server_address":"192.168.106.160:11800"}`))
	hydra.Conf.API(":8070", api.WithHeaderReadTimeout(30), api.WithTimeout(30, 30)).APM("skywalking", apm.WithEnable())
	hydra.Conf.RPC(":8071").APM(apmtypes.SkyWalking, apm.WithEnable())

	app := hydra.NewApp(
		hydra.WithServerTypes(rpc.RPC, xhttp.API),
		hydra.WithPlatName("test"),
		hydra.WithSystemName("test-performance"),
		hydra.WithClusterName("t"),
		hydra.WithRegistry("lm://."),
	)

	app.API("/entrance", func(ctx context.IContext) (r interface{}) {
		fmt.Println("api:inbound")
		caseVal, ok := ctx.Request().Get("case")

		mp, err := ctx.Request().GetMap()
		fmt.Println(mp, err)

		if !ok {
			return "err"
		}
		reqID := ctx.User().GetRequestID()
		fmt.Println("api:inbound", caseVal)
		switch caseVal {
		case "1":
			request := hydra.C.RPC().GetRegularRPC()
			response, err := request.Request(ctx.Context(), "/getgrpc", nil, crpc.WithXRequestID(reqID))
			if err != nil {
				return err
			}
			return response.Result

		default:

		}

		return nil
	})

	app.RPC("/getgrpc", func(ctx context.IContext) (r interface{}) {
		return ctx.User().GetRequestID()
	})

	os.Args = []string{"test", "run"}
	go app.Start()
	time.Sleep(time.Second * 2)

}

func buildRequest(caseval string, opts ...string) (bytes []byte, err error) {

	v := strings.Join(opts, "&")
	if len(v) > 0 {
		v = "&" + v
	}
	v = fmt.Sprintf("case=%s%s", caseval, v)
	fmt.Println(fmt.Sprintf("Sprintf:%s", v))

	req, err := http.NewRequest("GET", "http://localhost:8070/entrance?"+v, nil)
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

func main() {

	StartServer()
	N := 1
	for i := 0; i < N; i++ {
		bytes, err := buildRequest("1", fmt.Sprintf("time=%d", time.Now().Nanosecond()))
		if err != nil {
			panic(err)
			return
		}
		fmt.Println("result:", string(bytes))
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	time.Sleep(time.Second * 10)
}
