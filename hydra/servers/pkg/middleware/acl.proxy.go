package middleware

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"strings"

	"github.com/micro-plat/hydra/conf/server/acl/proxy"
)

//Proxy 代理配置
func Proxy() Handler {
	return func(ctx IMiddleContext) {
		proxy, err := ctx.APPConf().GetProxyConf()
		if err != nil {
			ctx.Response().Abort(http.StatusNotExtended, err)
			return
		}
		if proxy.Disable {
			ctx.Next()
			return
		}

		//检查当前请求是否需要进行代理
		cluster, need, err := proxy.Check()
		if err != nil {
			ctx.Response().AddSpecial("proxy")
			ctx.Response().Abort(http.StatusBadGateway, err)
			return
		}
		if !need {
			ctx.Next()
			return
		}

		//获取当前http信息
		ctx.Response().AddSpecial("proxy")
		useProxy(ctx, cluster)
	}
}
func useProxy(ctx IMiddleContext, cluster *proxy.UpCluster) {

	//检查当前请求
	req := ctx.Request().GetHTTPRequest()
	if strings.Contains(req.Header.Get("proxy"), ctx.APPConf().GetServerConf().GetServerID()) {
		ctx.Response().Abort(http.StatusBadGateway, fmt.Errorf("服务多次经过当前服务器 %s", req.Header.Get("proxy")))
		return
	}

	req.Header.Set("proxy", fmt.Sprintf("%s|%s", req.Header.Get("proxy"), ctx.APPConf().GetServerConf().GetServerID()))
	if req.Header.Get("X-Request-Id") == "" {
		req.Header.Set("X-Request-Id", ctx.User().GetTraceID())
	}

	//处理重试问题
	var num, max = 0, 10
RETRY:

	var canRetry = false
	var proxyError error
	num++
	if num > max {
		ctx.Response().Abort(http.StatusBadGateway, fmt.Errorf("重试超过服务器限制"))
		return
	}
	//获取服务器列表
	url, err := cluster.Next()
	if err != nil {
		goto RETRY
	}

	//转到上游
	rproxy := httputil.NewSingleHostReverseProxy(url)
	rproxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		proxyError = fmt.Errorf("远程请求出错[%d]%v %v", num, url, err)
		ctx.Log().Error(proxyError)
		if strings.Contains(err.Error(), "connect: connection refused") && num <= max {
			canRetry = true
		}
	}

	//处理代理服务
	resp := ctx.Response().GetHTTPReponse()
	response := newRWriter(resp)
	rproxy.ServeHTTP(response, req)

	//处理重试问题
	if canRetry {
		goto RETRY
	}
	if proxyError != nil {
		ctx.Response().Abort(http.StatusBadGateway, proxyError)
		return
	}

	ctx.Response().Abort(response.statusCode)
	rproxy = nil
}

type rWriter struct {
	http.ResponseWriter
	statusCode int
}

func newRWriter(w http.ResponseWriter) *rWriter {
	return &rWriter{w, http.StatusOK}
}

func (lrw *rWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (lrw *rWriter) Write(b []byte) (int, error) {
	return lrw.ResponseWriter.Write(b)
}
