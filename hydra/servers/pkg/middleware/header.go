package middleware

import (
	"net/http"

	"github.com/micro-plat/lib4go/types"
)

var originName = "Origin"
var hostName = "Host"

//Header 响应头设置
func Header() Handler {

	return func(ctx IMiddleContext) {

		//1. 获取header配置
		headers, err := ctx.APPConf().GetHeaderConf()
		if err != nil {
			ctx.Response().Abort(http.StatusNotExtended, err)
			return
		}
		if len(headers) > 0 {
			ctx.Response().AddSpecial("hdr")

		}

		//2. 处理响应header参数
		origin := ctx.Request().Headers().GetString(originName)
		hds := headers.GetHeaderByOrigin(types.GetString(origin, ctx.Response().GetHeaders().GetString(hostName)))
		for k, v := range hds {
			ctx.Response().Header(k, v)
		}

		//3. 业务处理
		ctx.Next()

	}
}
