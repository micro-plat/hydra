package middleware

import (
	"net/http"
)

//Render 响应结果输出组件
func Render() Handler {
	return func(ctx IMiddleContext) {

		ctx.Next()

		//加载渲染配置
		render, err := ctx.APPConf().GetRenderConf()
		if err != nil {
			ctx.Response().Abort(http.StatusNotExtended, err)
			return
		}
		if render.Disable {
			return
		}

		enable, status, ctp, content, err := render.Get(ctx.Request().Path().GetRequestPath(), nil, ctx.Response().GetRaw())
		if !enable {
			return
		}
		if err != nil {
			ctx.Log().Error("渲染响应结果出错:", err)
			return
		}

		ctx.Response().AddSpecial("render")
		ctx.Response().WriteFinal(status, content, ctp)
	}
}
