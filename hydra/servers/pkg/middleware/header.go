package middleware

import "net/http"

var originName = "Origin"

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

		//3. 处理响应header参数
		origin := ctx.Request().Path().GetHeader(originName)
		hds := headers.GetHeaderByOrigin(origin)
		for k, v := range hds {
			ctx.Response().Header(k, v)
		}

		//2. 业务处理
		ctx.Next()
	}
}
