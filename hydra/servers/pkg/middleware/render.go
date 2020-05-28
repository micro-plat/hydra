package middleware

import "github.com/micro-plat/lib4go/types"

//Render 响应结果输出组件
func Render() Handler {
	return func(ctx IMiddleContext) {

		ctx.Next()

		//加载渲染配置
		status := ctx.ServerConf().GetStatusRenderConf()
		content := ctx.ServerConf().GetContentRenderConf()
		if status.Disable && content.Disable {
			return
		}
		ctx.Response().AddSpecial("render")
		s, c := ctx.Response().GetResponse()
		var err error
		if !status.Disable {
			st, err := status.Get(nil, ctx.Response().GetRaw())
			if err != nil || types.GetInt(st, 0) == 0 {
				ctx.Log().Error("渲染响应状态码出错:", err)
				return
			}
			s = types.GetInt(st)

		}
		if !content.Disable {
			if c, err = content.Get(nil, ctx.Response().GetRaw()); err != nil {
				ctx.Log().Error("渲染响应内容出错:", err)
				return
			}
		}
		ctx.Response().Render(s, c)
	}
}
