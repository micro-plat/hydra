package main

import (
	"github.com/micro-plat/hydra"
	"github.com/micro-plat/hydra/context"
)

type OrderService struct {
}

func (o *OrderService) GetHandle(ctx hydra.IContext) interface{} {
	ctx.Log().Info("order.get.handle.order_no:", ctx.Request().GetString("order_no"))
	return "success"
}
func (o *OrderService) Handle(ctx hydra.IContext) interface{} {
	ctx.Log().Info("order.handle.order_no:", ctx.Request().GetString("order_no"))
	return "success"
}
func (o *OrderService) QueryHandle(ctx context.IContext) interface{} {
	ctx.Log().Info("order.query.handle.order_no:", ctx.Request().GetString("order_no"))
	return "success"
}
