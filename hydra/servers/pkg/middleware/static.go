package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/micro-plat/hydra/services"
)

//Static 静态文件处理插件
func Static() Handler {
	return func(ctx IMiddleContext) {
		//查询静态文件中是否存在
		static, err := ctx.APPConf().GetStaticConf()
		if err != nil {
			ctx.Response().Abort(http.StatusNotExtended, err)
			return
		}
		if static.Disable {
			ctx.Next()
			return
		}

		//处理option请求
		var rpath = ctx.Request().Path().GetRequestPath()
		var method = ctx.Request().Path().GetMethod()

		//是option则处理业务逻辑
		if doOption(ctx, static.Has(rpath)) {
			return
		}

		//优先后端服务调用
		var routerPath = ctx.GetRouterPath()
		if services.Def.Has(ctx.APPConf().GetServerConf().GetServerType(), routerPath, method) {
			ctx.Next()
			return
		}

		//检查请求类型是否为允许的类型
		if !static.AllowRequest(method) {
			ctx.Next()
			return
		}

		//读取静态文件
		ctx.Response().AddSpecial("static")
		fs, p, err := static.Get(rpath)
		if err != nil || fs == nil {
			ctx.Response().Abort(http.StatusNotFound, fmt.Errorf("文件不存在%s", rpath))
			return
		}

		//写入到响应流
		if strings.HasSuffix(p, ".gz") {
			ctx.Response().Header("Content-Encoding", "gzip")
		}
		ctx.Response().File(p, fs)
		return
	}
}
