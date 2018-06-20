package order

import (
	"time"

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
	input := make(map[string]interface{})
	if err := ctx.Request.GetJWT(&input); err != nil {
		return err
	}
	go func() {
		for {
			select {
			case <-time.After(time.Second):
				context.WSExchange.Notify(ctx.Request.GetUUID(), time.Now().Unix())
			}
		}

	}()
	return input
}
