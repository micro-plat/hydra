package middleware

import "github.com/micro-plat/hydra/global"

//Recovery 用于处理请求过程中出现的非预见的错误
func Recovery() Handler {
	return func(ctx IMiddleContext) {
		defer func() {
			if err := recover(); err != nil {
				ctx.Log().Errorf("[Recovery] panic recovered:\n%s\n%s", err, global.GetStack())
				ctx.Response().Abort(500)
			}
		}()
		ctx.Next()
	}
}
