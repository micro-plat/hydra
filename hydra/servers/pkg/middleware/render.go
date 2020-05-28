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

		enable, status, content, err := render.Get(ctx.Request().Path().GetPath(), ctx.Funcs(), ctx.Response().GetRaw())
		if !enable {
			return
		}
		ctx.Response().AddSpecial("render")
		if err != nil {
			ctx.Log().Error("渲染响应结果出错:", err)
			return
		}
		if status == 0 {
			status = ctx.Response().GetStatusCode()
		}
		ctx.Response().Render(status, content)
	}
}
