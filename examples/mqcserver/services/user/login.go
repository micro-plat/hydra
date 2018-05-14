package user

import (
	"github.com/micro-plat/hydra/component"
	"github.com/micro-plat/hydra/context"
)

type LoginHandler struct {
	container component.IContainer
	Name      string
}

func NewLoginHandler(container component.IContainer) (u *LoginHandler) {
	return &LoginHandler{container: container, Name: "LoginHandler"}
}
func (u *LoginHandler) Handle(ctx *context.Context) (r interface{}) {
	queue, err := ctx.GetContainer().GetQueue()
	if err != nil {
		return err
	}
	if err = queue.Push("hydra-20.mqc.001", `{"id":"1001"}`); err != nil {
		return err
	}
	return "success"
}
