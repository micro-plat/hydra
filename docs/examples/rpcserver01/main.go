package main

import (
	"fmt"
	"time"

	"github.com/micro-plat/hydra"
	//_ "github.com/micro-plat/hydra/components/pkgs/apm/skywalking"
	"github.com/micro-plat/hydra/components/pkgs/apm/apmtypes"
	crpc "github.com/micro-plat/hydra/components/rpcs/rpc"
	"github.com/micro-plat/hydra/hydra/servers/http"
	"github.com/micro-plat/hydra/hydra/servers/rpc"
	"github.com/micro-plat/lib4go/errs"
	"github.com/micro-plat/lib4go/types"
)

var buff []byte
var err error

type logstruct struct {
	ServerIP string `json:"server-ip"`
	Time     string `json:"time"`
	Level    string `json:"level"`
	Session  string `json:"session"`
}

//服务注册与系统勾子函数
func main() {
	app := hydra.NewApp(
		//hydra.WithPlatName("taobao"),
		//hydra.WithSystemName("apiserver"),
		hydra.WithDebug(),
		hydra.WithClusterName("t"),
		hydra.WithServerTypes(rpc.RPC, http.API),
		hydra.WithAPM(apmtypes.SkyWalking),
		hydra.WithPlatName("test"),
		hydra.WithSystemName("rpcserver01"),
	)
	hydra.Conf.RPC(":8281")
	app.API("/request", request)
	app.RPC("/rpc", rpcRequest)
	app.RPC("/rpc/log", log)
	app.Start()
}

func request(ctx hydra.IContext) (r interface{}) {
	ctx.Log().Info("id:", ctx.Request().GetString("id"))
	ctx.Log().Debug(
		hydra.G.GetClusterName(),
		hydra.G.GetPlatName(),
		hydra.G.GetSysName(),
	)
	response, err := hydra.C.RPC().GetRegularRPC().Request(ctx.Context(), "/rpc@test_debug", map[string]string{
		"tp": ctx.Request().GetString("id"),
	}, crpc.WithContentType("application/json"))

	if err != nil {
		return err
	}
	return response.Result
}
func log(ctx hydra.IContext) (r interface{}) {
	fmt.Println(ctx.Request().GetBodyMap())
	fmt.Println(ctx.Request().GetBody())
	fmt.Println("content:", ctx.Request().GetString("server-ip"))
	// list := make([]*logstruct, 0, 1)
	// err := ctx.Request().Bind(&list)
	return nil
}

func rpcRequest(ctx hydra.IContext) (r interface{}) {
	ctx.Log().Info("------request------")
	ctx.Log().Info(ctx.Request().GetMap())
	var bdmp types.XMap
	bdmp, err := ctx.Request().GetBodyMap()
	ctx.Log().Info(ctx.Request().GetBody())

	ctx.Log().Error(err)
	switch bdmp.GetString("tp") { //从路由配置中获取参数值 ctx.Request.Param.Get...
	case "1":
		return "success"
	case "2":
		return 100
	case "3":
		return time.Now().String()
	case "4":
		return float32(100.20)
	case "5":
		return true
	case "6":
		type order struct {
			ID string `json:"id"`
		}
		type result struct {
			Name   string   `json:"name"`
			Age    int      `json:"age"`
			Orders []*order `json:"order"`
		}
		return &result{Name: "colin", Age: 8, Orders: []*order{&order{ID: "897776666"}}}
	case "7":
		return `{"name":"colin","age":8}`
	case "8":
		return map[string]string{
			"order": "123456",
		}
	case "9":
		return map[string]interface{}{
			"product": map[string]string{
				"price": "100",
			},
		}
	case "10":
		return `<?xml version='1.0'?><xml><name>colin</name><age>8</age></xml>`
	case "11":
		ctx.Response().ContentType("application/xml; charset=UTF-8")
		type order struct {
			ID string `json:"id" xml:"id"`
		}
		type result struct {
			Name   string   `json:"name" xml:"name"`
			Age    int      `json:"age" xml:"age"`
			Orders []*order `json:"orders" xml:"orders"`
		}
		return &result{Name: "colin", Age: 8, Orders: []*order{&order{ID: "897776666"}}}
	case "12":
		return errs.NewError(201, "无需处理")
	case "13":
		if err := ctx.Request().Check("order_id"); err != nil {
			return err
		}
		return ctx.Request().GetString("order_id")
	case "14":
		ctx.Log().Info(ctx.Request().GetBody())
		return fmt.Sprintf(`<?xml version='1.0'?><xml><name>%s</name><age>8</age></xml>`, ctx.Request().GetString("name"))
	case "15":
		ctx.Log().Info(ctx.Request().GetBody())
		r, err := ctx.Request().GetBody()
		if err != nil {
			return err
		}
		return r
	default:
		return hydra.G.PlatName
	}
}
