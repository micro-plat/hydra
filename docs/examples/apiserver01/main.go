package main

import (
	"fmt"
	"time"

	"github.com/micro-plat/hydra"
	"github.com/micro-plat/hydra/conf/server/router"
	"github.com/micro-plat/hydra/hydra/servers/http"
	"github.com/micro-plat/lib4go/errs"
)

//服务器各种返回结果
func main() {
	app := hydra.NewApp(
		hydra.WithServerTypes(http.API),
		hydra.WithUsage("apiserver"),
		hydra.WithDebug(),
	)

	app.API("/order/request/:tp", request)
	app.API("/order/encoding/:tp", request, router.WithEncoding("gbk"))
	app.Start()
}

func request(ctx hydra.IContext) (r interface{}) {
	ctx.Log().Info("------request------")
	switch ctx.Request().Param("tp") { //从路由配置中获取参数值 ctx.Request.Param.Get...
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
		fmt.Println("name:", ctx.Request().GetString("name"))
		// fmt.Println(ctx.Request().GetBody())
		return ctx.Request().GetString("name")
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
