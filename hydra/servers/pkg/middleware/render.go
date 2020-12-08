package middleware

import (
	"net/http"
)

//Render 响应结果输出组件
func Render() Handler {
	return func(ctx IMiddleContext) {

		//render是最后根据配置修改相应结果   所以需要先执行了业务逻辑 然后在执行下面逻辑
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

		rd, enable, err := render.Get()
		if err != nil {
			ctx.Log().Error("render出错:", err)
			return
		}

		if !enable {
			return
		}
		if rd.ContentType != "" {
			ctx.Response().ContentType(rd.ContentType)
		}
		ctx.Response().AddSpecial("render")
		ctx.Response().Write(rd.Status, rd.Content)
	}
}
