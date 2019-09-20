package order

import (
	"github.com/micro-plat/hydra/component"
	"github.com/micro-plat/hydra/context"
)

type RequestHandler struct {
	container component.IContainer
}

func NewRequestHandler(container component.IContainer) (u *RequestHandler, err error) {
	return &RequestHandler{container: container}, nil
}

//Handle .
func (u *RequestHandler) Handle(ctx *context.Context) (r interface{}) {
	return "success"
}
