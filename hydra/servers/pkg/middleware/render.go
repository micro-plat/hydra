package middleware

//Render 响应结果输出组件
func Render() Handler {
	return func(ctx IMiddleContext) {

		ctx.Next()

		//加载渲染配置
		render := ctx.ServerConf().GetRenderConf()
		if render.Disable {
			return
		}

		enable, status, ctp, content, err := render.Get(ctx.Request().Path().GetRequestPath(), ctx.TmplFuncs(), ctx.Response().GetRaw())
		if !enable {
			return
		}
		ctx.Response().AddSpecial("render")
		if err != nil {
			ctx.Log().Error("渲染响应结果出错:", err)
			return
		}
		ctx.Response().WriteFinal(status, content, ctp)
	}
}
