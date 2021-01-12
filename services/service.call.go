package services

import (
	"net/http"

	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/lib4go/errs"
)

func (s *regist) Call(ctx context.IContext, service string) (result interface{}) {
	//获取处理服务
	h, ok := Def.GetHandler(ctx.APPConf().GetServerConf().GetServerType(), service)
	if !ok {
		ctx.Response().AddSpecial("hdl")
		return errs.NewErrorf(http.StatusNotFound, "未找到服务:%s", service)
	}

	//预处理,用户资源检查，发生错误后不再执行业务处理-------
	globalHandlings := Def.GetHandleExecutings(ctx.APPConf().GetServerConf().GetServerType())
	for _, h := range globalHandlings {
		result := h.Handle(ctx)
		if err := errs.GetError(result); err != nil {
			return result
		}
	}

	handlings := Def.GetHandlings(ctx.APPConf().GetServerConf().GetServerType(), service)
	for _, h := range handlings {
		result := h.Handle(ctx)
		if err := errs.GetError(result); err != nil {
			return result
		}
	}

	//业务处理----------------------------------
	result = h.Handle(ctx)

	//后处理，处理资源回收，无论业务处理返回什么结果都会执行--
	handleds := Def.GetHandleds(ctx.APPConf().GetServerConf().GetServerType(), service)
	for _, h := range handleds {
		hresult := h.Handle(ctx)
		if err := errs.GetError(hresult); err != nil {
			ctx.Log().Error("后处理发生错误　err:", err)
		}
	}

	//后处理，处理资源回收，无论业务处理返回什么结果都会执行--
	globalHandleds := Def.GetHandleExecuted(ctx.APPConf().GetServerConf().GetServerType())
	for _, h := range globalHandleds {
		hresult := h.Handle(ctx)
		if err := errs.GetError(hresult); err != nil {
			ctx.Log().Error("后处理发生错误　err:", err)
		}
	}

	//处理输出, 只会将业务处理结果进行输出---------------
	return result
}
