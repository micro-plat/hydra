package main

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/micro-plat/hydra"
	"github.com/micro-plat/hydra/hydra/servers/http"
	"github.com/micro-plat/hydra/hydra/servers/rpc"
	"github.com/micro-plat/lib4go/errs"
)

var buff []byte
var err error

//服务注册与系统勾子函数
func main() {
	app := hydra.NewApp(
		hydra.WithServerTypes(rpc.RPC, http.API),
	)
	buff, err = ioutil.ReadFile("/home/yanglei/abc.doc")
	if err != nil {
		fmt.Println(err)
	}
	app.API("/request", request)
	app.RPC("/rpc", rpcRequest)
	app.Start()
}

func request(ctx hydra.IContext) (r interface{}) {
	ctx.Log().Info("id:", ctx.Request().GetString("id"))

	response, err := hydra.Component.RPC().GetRegularRPC().Request(ctx.Context(), "/rpc", buff)

	if err != nil {
		return err
	}
	return response.Result
}

func rpcRequest(ctx hydra.IContext) (r interface{}) {
	ctx.Log().Info("------request------")
	switch ctx.Request().GetString("tp") { //从路由配置中获取参数值 ctx.Request.Param.Get...
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
		return hydra.Global.PlatName
	}
}
