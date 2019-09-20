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
	//检查用户名密码是否正确
	ctx.Response.SetJWT(map[string]interface{}{
		"id": 11000,
	})
	return "ok"
}
