package middleware

import (
	"fmt"
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

		rd, enable, err := render.Get(ctx.Request().Path().GetRequestPath())
		if !enable {
			return
		}
		if err != nil {
			ctx.Log().Error("render出错:", err)
			return
		}
		fmt.Println("render:", rd)

		ctx.Response().AddSpecial("render")
		ctx.Response().WriteFinal(rd.Status, rd.Content, rd.ContentType)
	}
}
