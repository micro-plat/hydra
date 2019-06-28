package main

import (
	"github.com/micro-plat/hydra/context"
)

func (api *apiserver) handling() {
	api.MicroApp.Handling(func(ctx *context.Context) (rt interface{}) {
		//检验登录状态或权限
		return nil
	})
}
