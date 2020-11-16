package middleware

import (
	"fmt"
	"net/http"
	"os"
)

//Static 静态文件处理插件
func Static() Handler {
	return func(ctx IMiddleContext) {
		static, err := ctx.APPConf().GetStaticConf()
		if err != nil {
			ctx.Response().Abort(http.StatusNotExtended, err)
			return
		}
		if static.Disable {
			ctx.Next()
			return
		}
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
				ctx.Response().Abort(http.StatusNotFound, err)
				return
			}
			err := fmt.Errorf("%s,err:%v", fpath, err)
			ctx.Response().Abort(http.StatusInternalServerError, err)
			return
		}
		if finfo.IsDir() {
			err := fmt.Errorf("找不到文件:%s", fpath)
			ctx.Response().Abort(http.StatusNotFound, err)
			return
		}
		//fmt.Println("statis:", static.GetFileMap(), fpath)
		//文件已存在，则返回文件
		ctx.Response().StatusCode(200)
		if gzfile := static.GetGzFile(fpath); gzfile != "" {
			ctx.Response().File(gzfile)
			ctx.Response().Header("Content-Encoding", "gzip")
			return
		}
		ctx.Response().File(fpath)
		return
	}
}
