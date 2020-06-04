package middleware

import (
	"fmt"
	"net/http"
	"os"
)

//Static 静态文件处理插件
func Static() Handler {
	return func(ctx IMiddleContext) {
		static := ctx.ServerConf().GetStaticConf()

		//检查文件是否需要按静态文件处理
		var rpath = ctx.Request().Path().GetRequestPath()
		var method = ctx.Request().Path().GetMethod()
		ok, fpath := static.IsStatic(rpath, method)
		if !ok {
			ctx.Next()
			return
		}

		ctx.Response().AddSpecial("static")
		finfo, err := os.Stat(fpath)
		if err != nil {
			if os.IsNotExist(err) {
				err := fmt.Errorf("找不到文件:%s %w", fpath, err)
				ctx.Response().AbortWithError(http.StatusNotFound, err)
				return
			}
			err := fmt.Errorf("%s,err:%v", fpath, err)
			ctx.Response().AbortWithError(http.StatusInternalServerError, err)
			return
		}
		if finfo.IsDir() {
			err := fmt.Errorf("找不到文件:%s", fpath)
			ctx.Response().AbortWithError(http.StatusNotFound, err)
			return
		}
		//文件已存在，则返回文件
		ctx.Response().File(fpath)
		return
	}
}
