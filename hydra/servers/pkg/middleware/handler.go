package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/micro-plat/hydra/components"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/hydra/services"
)

//ExecuteHandler 业务处理Handler
func ExecuteHandler(service string) Handler {
	return func(ctx IMiddleContext) {
		//检查是否被限流
		ctx.Service(service) //保存服务信息
		if ctx.Request().Path().IsLimited() {
			//降级处理
			fallback(ctx, service)
			return
		}
		//处理RPC服务调用
		if addr, ok := global.IsProto(service, global.ProtoRPC); ok {
			response, err := components.Def.RPC().GetRegularRPC().Swap(addr, ctx)
			if err != nil {
				ctx.Response().Write(response.GetStatus(), err)
				return
			}
			headers := response.GetHeaders()
			for k := range headers {
				ctx.Response().Header(k, headers.GetString(k))
			}
			ctx.Response().Write(response.GetStatus(), response.GetResult())
			return
		}

		//处理服务请求
		if procSvc(ctx, service) {
			return
		}

		//处理静态资源服务
		if procStatic(ctx, service) {
			return
		}

		ctx.Response().Abort(http.StatusNotFound, fmt.Errorf("未找到路径:%s", ctx.Request().Path().GetRequestPath()))
	}
}

//处理服务请求
func procSvc(ctx IMiddleContext, service string) (hasSvc bool) {
	//检查服务中是否包含当前请求路径
	hasSvc = services.Def.Has(ctx.APPConf().GetServerConf().GetServerType(), service)
	//处理option
	if hasSvc && checkOption(ctx) {
		return
	}

	if hasSvc {
		result := services.Def.Call(ctx, service)
		ctx.Response().WriteAny(result)
		return
	}
	return
}

//处理静态资源服务
func procStatic(ctx IMiddleContext, service string) (exists bool) {
	//处理静态文件
	exists, filePath, fs := getStatic(ctx, service)
	if exists && checkOption(ctx) {
		return
	}
	if exists {
		//写入到响应流
		if strings.HasSuffix(filePath, ".gz") {
			ctx.Response().Header("Content-Encoding", "gzip")
		}
		ctx.Response().File(filePath, fs)
		return
	}
	return
}

//checkOption 请求处理
func checkOption(ctx IMiddleContext) bool {
	if strings.ToUpper(ctx.Request().Path().GetMethod()) != http.MethodOptions {
		return false
	}
	ctx.Response().AddSpecial("opt")
	ctx.Response().Abort(http.StatusOK, nil)
	return true
}
