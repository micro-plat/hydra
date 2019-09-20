package order

import (
	"github.com/micro-plat/hydra/component"
	"github.com/micro-plat/hydra/context"
)

type QueryHandler struct {
	container component.IContainer
}

func NewQueryHandler(container component.IContainer) (u *QueryHandler) {
	return &QueryHandler{container: container}
}
func (u *QueryHandler) Handle(ctx *context.Context) (r interface{}) {
	status, r, params, err := ctx.RPC.Request("/order/bind", "GET", nil, nil, true)
	ctx.Log.Info("handle:", status, r, params, err)
	// ctx.Response.SetHeaders(params)
	ctx.Response.MustContent(status, r)

	return

}
