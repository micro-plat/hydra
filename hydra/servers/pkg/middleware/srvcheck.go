package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"

	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/hydra/servers/pkg/dispatcher"
)

//GinServiceExistsCheck 用于处理请求过程中出现的非预见的错误
func GinServiceExistsCheck(engine *gin.Engine) Handler {
	callback := ginServiceExistsCheck(engine)
	return func(ctx IMiddleContext) {
		if !callback(ctx.Request().Path()) {
			ctx.Response().Abort(http.StatusNotFound, fmt.Errorf("未找到服务:%s", ctx.Request().Path().GetURL()))
		}
		ctx.Next()
	}
}

//DispServiceExistsCheck 用于处理请求过程中出现的非预见的错误
func DispServiceExistsCheck(engine *dispatcher.Engine) Handler {
	callback := dispServiceExistsCheck(engine)
	return func(ctx IMiddleContext) {
		if !callback(ctx.Request().Path()) {
			ctx.Response().Abort(http.StatusNotFound, fmt.Errorf("未找到服务:%s", ctx.Request().Path().GetURL()))
		}
		ctx.Next()
	}
}

func ginServiceExistsCheck(engine *gin.Engine) func(context.IPath) bool {
	var syncMap sync.Map
	return func(path context.IPath) bool {
		key := fmt.Sprintf("%s_%s", path.GetMethod(), path.GetRequestPath())
		if val, ok := syncMap.Load(key); ok {
			return val.(bool)
		}
		routers := engine.Routes()
		for _, r := range routers {
			if checkExists(r.Method, r.Path, path.GetMethod(), path.GetRequestPath()) {
				syncMap.Store(key, true)
				return true
			}
		}
		syncMap.Store(key, false)
		return false
	}
}

func dispServiceExistsCheck(engine *dispatcher.Engine) func(context.IPath) bool {
	var syncMap sync.Map
	return func(path context.IPath) bool {
		key := fmt.Sprintf("%s_%s", path.GetMethod(), path.GetRequestPath())
		if val, ok := syncMap.Load(key); ok {
			return val.(bool)
		}
		routers := engine.Routes()
		for _, r := range routers {
			if checkExists(r.Method, r.Path, path.GetMethod(), path.GetRequestPath()) {
				syncMap.Store(key, true)
				return true
			}
		}
		syncMap.Store(key, false)
		return false
	}
}

func checkExists(routeMethod, routePath, reqMethod, reqPath string) (found bool) {
	found = false
	if !strings.EqualFold(routeMethod, reqMethod) {
		return
	}

	if routePath == reqPath {
		found = true
		return
	}
	routepts := strings.Split(routePath, "/")
	pathpts := strings.Split(reqPath, "/")
	if len(routepts) != len(pathpts) {
		return
	}
	ismatch := true
	for i := range routepts {
		if routepts[i] == pathpts[i] {
			continue
		}
		if strings.HasPrefix(routepts[i], ":") {
			continue
		}
		ismatch = false
	}
	if ismatch {
		found = true
		return
	}
	return
}
