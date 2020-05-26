package main

import (
	"github.com/micro-plat/hydra"
	"github.com/micro-plat/hydra/context"
)

type OrderService struct {
}

func (o *OrderService) GetHandle(ctx hydra.IContext) interface{} {
	ctx.Log().Info("order.get.handle:", ctx.Request().GetString("order_no"))
	return "order.get.handle"
}

func (o *OrderService) Handling(ctx hydra.IContext) interface{} {
	ctx.Log().Info("order.Handling", ctx.Request().GetString("order_no"))
	return "order.all.handling"
}
func (o *OrderService) GetHandling(ctx hydra.IContext) interface{} {
	ctx.Log().Info("order.GetHandling", ctx.Request().GetString("order_no"))
	return "order.get.handling"
}
func (o *OrderService) GetHandled(ctx hydra.IContext) interface{} {
	ctx.Log().Info("order.GetHandled", ctx.Request().GetString("order_no"))
	return "order.all.handling"
}
func (o *OrderService) Handled(ctx hydra.IContext) interface{} {
	ctx.Log().Info("order.Handled", ctx.Request().GetString("order_no"))
	return "order.all.handling"
}
func (o *OrderService) PostHandle(ctx hydra.IContext) interface{} {
	ctx.Log().Info("order_id:", ctx.Request().GetString("order_id"))
	t, err := ctx.Request().GetBody()
	if err != nil {
		return err
	}
	return t
}
func (o *OrderService) Handle(ctx hydra.IContext) interface{} {
	ctx.Log().Info("order.handle.order_no:", ctx.Request().GetString("order_no"))
	return "order.others.handle"
}
func (o *OrderService) QueryHandle(ctx context.IContext) interface{} {
	ctx.Log().Info("order.query.handle.order_no:", ctx.Request().GetString("order_no"))
	return "order.query.handle"
}
func (o *OrderService) Close() error {
	hydra.Global.Log().Info("order.Close")
	return nil
}
