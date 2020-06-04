package middleware

import (
	"fmt"
	"net/http"

	"github.com/micro-plat/hydra/services"
	"github.com/micro-plat/lib4go/errs"
)

//ExecuteHandler 业务处理Handler
func ExecuteHandler(service string) Handler {
	return func(ctx IMiddleContext) {

		h, ok := services.Def.GetHandler(ctx.ServerConf().GetMainConf().GetServerType(), service)
		if !ok {
			ctx.Response().AddSpecial("handler")
			ctx.Response().AbortWithError(http.StatusNotFound, fmt.Errorf("未找到服务%s", service))
			return
		}

		//预处理,用户资源检查，发生错误后不再执行业务处理-------

		globalHandlings := services.Def.GetHandleExecutings(ctx.ServerConf().GetMainConf().GetServerType())
		for _, h := range globalHandlings {
			result := h.Handle(ctx)
			if err := errs.GetError(result); err != nil {
				ctx.Log().Error("预处理发生错误 err:", err)
				ctx.Response().WriteAny(result)
				return
			}
		}

		handlings := services.Def.GetHandlings(ctx.ServerConf().GetMainConf().GetServerType(), service)
		for _, h := range handlings {
			result := h.Handle(ctx)
			if err := errs.GetError(result); err != nil {
				ctx.Log().Error("预处理发生错误 err:", err)
				ctx.Response().WriteAny(result)
				return
			}
		}

		//业务处理----------------------------------
		result := h.Handle(ctx)

		//后处理，处理资源回收，无论业务处理返回什么结果都会执行--
		handleds := services.Def.GetHandleds(ctx.ServerConf().GetMainConf().GetServerType(), service)
		for _, h := range handleds {
			hresult := h.Handle(ctx)
			if err := errs.GetError(hresult); err != nil {
				ctx.Log().Error("后处理发生错误　err:", err)
			}
		}

		//后处理，处理资源回收，无论业务处理返回什么结果都会执行--
		globalHandleds := services.Def.GetHandleExecuted(ctx.ServerConf().GetMainConf().GetServerType())
		for _, h := range globalHandleds {
			hresult := h.Handle(ctx)
			if err := errs.GetError(hresult); err != nil {
				ctx.Log().Error("后处理发生错误　err:", err)
			}
		}

		//处理输出, 只会将业务处理结果进行输出---------------
		if err := errs.GetError(result); err != nil {
			ctx.Log().Error("err:", err)
		}
		ctx.Response().WriteAny(result)
	}
}
