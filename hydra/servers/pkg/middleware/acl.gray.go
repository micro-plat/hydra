package middleware

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"strings"

	"github.com/micro-plat/hydra/conf/server/acl/gray"
)

//Gray 灰度配置
func Gray() Handler {
	return func(ctx IMiddleContext) {
		gray, err := ctx.APPConf().GetGrayConf()
		if err != nil {
			ctx.Response().Abort(http.StatusNotExtended, err)
			return
		}
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
		useProxy(ctx, gray)

	}
}
func useProxy(ctx IMiddleContext, gray *gray.Gray) {

	//检查当前请求
	req, resp := ctx.GetHttpReqResp()
	if req == nil || resp == nil {
		panic(fmt.Errorf("只有api,web服务器支持灰度配置"))
	}

	//处理重试问题
	var num, max = 0, 10
RETRY:
	var canRetry = false
	var proxyError error
	num++

	//获取服务器列表
	url, err := gray.Next()
	if err != nil {
		ctx.Response().Abort(http.StatusBadGateway, fmt.Errorf("无法获取上游服务器地址:%w", err))
		return
	}

	//转到上游
	proxy := httputil.NewSingleHostReverseProxy(url)
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		proxyError = fmt.Errorf("远程请求出错[%d]%v %v", num, url, err)
		if strings.Contains(err.Error(), "connect: connection refused") && num <= max {
			canRetry = true
		}
	}

	//处理代理服务
	proxy.ServeHTTP(resp, req)

	//处理重试问题
	if canRetry {
		goto RETRY
	}
	if proxyError != nil {
		ctx.Response().Abort(http.StatusBadGateway, proxyError)
		return
	}
	ctx.Response().Stop(http.StatusUseProxy)
}
