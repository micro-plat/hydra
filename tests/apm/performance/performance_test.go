package performance

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/micro-plat/hydra"
	"github.com/micro-plat/hydra/conf/server/apm"
	varapm "github.com/micro-plat/hydra/conf/var/apm"
	"github.com/micro-plat/hydra/context"

	xhttp "github.com/micro-plat/hydra/hydra/servers/http"
	"github.com/micro-plat/hydra/hydra/servers/rpc"
)

func StartServer( ) {

	hydra.Conf.Vars().APM(varapm.New("skywalking", `{"server_address":"192.168.106.160:11800"}`))
	hydra.Conf.API(":8070").APM("skywalking", apm.WithEnable())
	hydra.Conf.RPC(":8071").APM("skywalking", apm.WithEnable())

	app := hydra.NewApp(
		hydra.WithServerTypes(rpc.RPC, xhttp.API),
		hydra.WithPlatName("test"),
		hydra.WithSystemName("test-performance"),
		hydra.WithClusterName("t"),
		hydra.WithRegistry("lm://"),
	)

	app.API("/entrance", func(ctx context.IContext) (r interface{}) {
		fmt.Println("api:inbound")
		caseVal, ok := ctx.Request().Get("case")
		if !ok {
			return "err"
		}
		reqID := ctx.User().GetRequestID
		fmt.Println("api:inbound", caseVal)
		switch caseVal {
		case "1":
			request := hydra.C.RPC().GetRegularRPC()
			response, err := request.Request(ctx.Context(), "/getgrpc", nil, rpc.WithXRequestID(reqID))
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

func buildRequest(caseval string,opts...string) (bytes []byte, err error) {

	v:= strings.Join(opts,"&")
	if len(v)>0 {
		v="&"+v
	}

	reader := strings.NewReader(fmt.Sprintf("case=%s%s", caseval,v))

	req, err := http.NewRequest("get", "http://localhost:8070/entrance",reader)
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

func BenchmarkApmPerformance(b *testing.B) {
	StartServer()

	b.ResetTimer()
	b.N = 1

	for i := 0 ;i<b.N ;i++{
		bytes, err := buildRequest("1",fmt.Sprintf("time=%s",time.Now().Na))
		if err != nil {
			b.Error(err)
			return
		}
		if string(bytes) != "1" {
			t.Error("result", string(bytes))
			return
		}
	}

 

	time.Sleep(time.Minute)
}
