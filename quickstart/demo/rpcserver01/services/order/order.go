package order

import (
	"github.com/micro-plat/hydra/component"
	"github.com/micro-plat/hydra/context"
)

type OrderHandler struct {
	container component.IContainer
}

func NewOrderHandler(container component.IContainer) (u *OrderHandler) {
	return &OrderHandler{
		container: container,
	}
}

//QueryHandle 充值结果查询
func (u *OrderHandler) QueryHandle(ctx *context.Context) (r interface{}) {
	ctx.Log.Info("--------------充值结果查询---------------")
	return "SUCCESS"
}
