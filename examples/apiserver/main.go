package main

import (
	"fmt"
	"time"

	"github.com/micro-plat/hydra"
	"github.com/micro-plat/hydra/servers/http"
	"github.com/micro-plat/lib4go/errs"
)

func main() {
	app := hydra.NewApp(
		hydra.WithServerTypes(http.API),
		hydra.WithUsage("一键加油-apiserver"),
		hydra.WithDebug(),
	)

	app.API("/order/request/:tp", request)
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
		return errs.NewError(201, "无需处理")
	case "12":
		if err := ctx.Request().Check("order_id"); err != nil {
			return err
		}
		return ctx.Request().GetString("order_id")
	case "13":
		return hydra.Application.PlatName
	default:
		return fmt.Errorf("值错误，请传入1-12")
	}
}
