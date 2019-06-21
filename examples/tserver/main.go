package main

import (
	"github.com/micro-plat/hydra/component"
	"github.com/micro-plat/hydra/context"

	"github.com/micro-plat/hydra/hydra"
)

func main() {
	app := hydra.NewApp(
		hydra.WithPlatName("myplat"),
		hydra.WithSystemName("helloserver"),
		hydra.WithServerTypes("api"), //服务器类型为http api
		hydra.WithDebug(),
	)

	app.Micro("/order", NewOrderHandler)
	app.Start()
}

type order struct {
	OrderNo   string `json:"order_no" xml:"order_no"`
	ProductID int    `json:"pid" xml:"pid"`
	Amount    int    `json:"amount" xml:"amount"`
}

type OrderHandler struct {
	container component.IContainer
}

func NewOrderHandler(container component.IContainer) (u *OrderHandler) {
	return &OrderHandler{container: container}
}
func (u *OrderHandler) RequestHandle(ctx *context.Context) (r interface{}) {
	tp := ctx.Request.GetString("result_type")
	switch tp {
	case "xml":
		ctx.Response.SetXML()
	case "json":
		ctx.Response.SetJSON()
	default:
	}
	return &order{
		OrderNo:   "89769037987666",
		ProductID: 5989234,
		Amount:    10,
	}

}
