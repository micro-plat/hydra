package middleware

import (
	"github.com/micro-plat/hydra/services"
)

func fallback(ctx IMiddleContext, service string) bool {
	if ok := ctx.Request().Path().AllowFallback(); !ok {
		return false
	}

	//获取处理服务
	ctx.Response().AddSpecial("fallback")
	fallback, ok := services.Def.GetFallback(ctx.APPConf().GetServerConf().GetServerType(), service)
	if !ok {
		return false
	}
	result := fallback.Handle(ctx)
	//处理输出, 只会将业务处理结果进行输出---------------
	ctx.Response().WriteAny(result)
	return true
}
