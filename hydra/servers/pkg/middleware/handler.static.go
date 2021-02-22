package middleware

import (
	"fmt"
	"net/http"
)

func getStatic(ctx IMiddleContext, service string) (exists bool, filePath string, fs http.FileSystem) {

	//查询静态文件中是否存在
	static, err := ctx.APPConf().GetStaticConf()
	if err != nil {
		return
	}
	if static.Disable {
		return
	}

	//检查文件是否需要按静态文件处理
	ctx.Response().AddSpecial("static")
	var rpath = ctx.Request().Path().GetRequestPath()
	var method = ctx.Request().Path().GetMethod()
	if !static.AllowRequest(method) {
		return
	}

	//读取静态文件
	fs, filePath, err = static.Get(rpath)
	if err != nil || fs == nil {
		ctx.Response().ContentType("text/plain")
		ctx.Response().Abort(http.StatusNotFound, fmt.Errorf("文件不存在:%s", rpath))
		return
	}
	exists = true
	return
}
