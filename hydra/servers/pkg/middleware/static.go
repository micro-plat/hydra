package middleware

import (
	"fmt"
	"os"
)

//Static 静态文件处理插件
func Static() Handler {
	return func(ctx IMiddleContext) {
		static := ctx.ServerConf().GetStaticConf()
		if static.Disable || !static.AllowRequest(ctx.Request().Path().GetMethod()) {
			ctx.Next()
			return
		}
		var rpath = ctx.Request().Path().GetRequestPath()
		ok, fpath := static.IsStatic(rpath)
		if !ok {
			ctx.Next()
			return
		}

		ctx.Response().AddSpecial("static")
		finfo, err := os.Stat(fpath)
		if err != nil {
			if os.IsNotExist(err) {
				err := fmt.Errorf("找不到文件:%s", fpath)
				ctx.Response().AbortWithError(404, err)
				return
			}
			err := fmt.Errorf("%s,err:%v", fpath, err)
			ctx.Response().AbortWithError(500, err)
			return
		}
		if finfo.IsDir() {
			err := fmt.Errorf("找不到文件:%s", fpath)
			ctx.Response().AbortWithError(404, err)
			return
		}
		//文件已存在，则返回文件
		ctx.Response().File(fpath)
		return
	}
}
