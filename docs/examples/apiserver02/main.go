package main

import (
	"fmt"
	"time"

	"github.com/micro-plat/hydra"

	"github.com/micro-plat/hydra/components"
	"github.com/micro-plat/hydra/components/rpcs/rpc"
	 
	cfapp "github.com/micro-plat/hydra/conf/app"

	"github.com/micro-plat/hydra/conf/server/api"
	"github.com/micro-plat/hydra/hydra/servers/http"
)

//服务注册与系统勾子函数
func main() {
	hydra.Conf.API(":8082", api.WithTrace())
	app := hydra.NewApp(
		hydra.WithServerTypes(http.API),
		hydra.WithDebug(),
		hydra.WithPlatName("test"),
		hydra.WithSystemName("apiserver02"),
	)

	app.API("/order/request", request)
	app.API("/order", &OrderService{})

	app.OnStarting(func(cfapp.IAPPConf) error {
		hydra.G.Log().Info("server.OnServerStarting")
		return nil
	})

	app.OnClosing(func(cfapp.IAPPConf) error {
		hydra.G.Log().Info("server.OnServerClosing")
		return nil
	})

	app.OnHandleExecuting(func(ctx hydra.IContext) interface{} {
		ctx.Log().Info("global.OnHandleExecuting")
		return nil
	})
	app.OnHandleExecuted(func(ctx hydra.IContext) interface{} {
		ctx.Log().Info("global.OnHandleExecuted")
		return nil
	})
	app.Start()
}

func request(ctx hydra.IContext) (r interface{}) {
	time.Sleep(3 * time.Second)
	request, err := components.Def.RPC().GetRPC()
	if err != nil {
		return fmt.Errorf("RPC().GetRPC:%s", err.Error())
	}
	resp, err := request.Request(ctx.Context(), "/rpc", map[string]string{
		"tp": "1",
	}, rpc.WithContentType("application/json"))
	if err != nil {
		return fmt.Errorf("RPC.Request%s", err.Error())
	}
	fmt.Println("RPC.Header", resp.Header)
	fmt.Println("RPC.Result", resp.Result)
	fmt.Println("RPC.Status", resp.Status)

	client, err := components.Def.HTTP().GetClient()
	if err != nil {
		return fmt.Errorf("HTTP().GetClient:%s", err.Error())
	}

	content, status, err := client.Post("http://192.168.5.108:8083/index", "a=1&b=2")
	fmt.Println("HTTP1.content", content)
	fmt.Println("HTTP1.err", err)
	fmt.Println("HTTP1.Status", status)

	content, status, err = client.Post("http://192.168.5.108:8084/order/request", "a=1&b=2")
	fmt.Println("HTTP2.content", content)
	fmt.Println("HTTP2.err", err)
	fmt.Println("HTTP2.Status", status)

	return "request"
}
