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

func (u *OrderHandler) Handle(ctx *context.Context) (r interface{}) {
	// current, totoal := ctx.Request.GetSharding()
	return nil
}
