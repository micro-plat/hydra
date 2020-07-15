package middleware

import (
	"fmt"
	"net/http"
	"net/http/httputil"
)

//Gray 灰度配置
func Gray() Handler {
	return func(ctx IMiddleContext) {
		gray := ctx.ServerConf().GetGray()
		if gray.Disable {
			ctx.Next()
			return
		}

		//检查当前请求是否需要进行灰度
		need, err := gray.Check(ctx.TmplFuncs(), nil)
		if err != nil {
			ctx.Response().AddSpecial("gray")
			ctx.Response().Abort(http.StatusBadGateway, err)
			return
		}
		if !need {
			ctx.Next()
			return
		}

		//获取当前http信息
		ctx.Response().AddSpecial("gray")
		req, resp := ctx.GetHttpReqResp()
		if req == nil || resp == nil {
			ctx.Response().Abort(http.StatusBadGateway, fmt.Errorf("只有api,web服务器支持灰度配置"))
			return
		}

		//获取服务器列表
		url, err := gray.Next()
		if err != nil {
			ctx.Response().Abort(http.StatusBadGateway, fmt.Errorf("无法获取上游服务器地址:%w", err))
			return
		}

		//转到上游
		proxy := httputil.NewSingleHostReverseProxy(url)
		proxy.ServeHTTP(resp, req)
		ctx.Response().Stop(http.StatusUseProxy)
	}
}
