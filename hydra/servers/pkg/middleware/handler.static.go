package middleware

import (
	"fmt"
	"net/http"
	"strings"
)

func doStatic(ctx IMiddleContext, service string) bool {

	//查询静态文件中是否存在
	static, err := ctx.APPConf().GetStaticConf()
	if err != nil {
		return false
	}
	if static.Disable {
		return false
	}

	//检查文件是否需要按静态文件处理
	ctx.Response().AddSpecial("static")
	var rpath = ctx.Request().Path().GetRequestPath()
	var method = ctx.Request().Path().GetMethod()
	if !static.AllowRequest(method) {
		return false
	}

	//读取静态文件
	fs, p, err := static.Get(rpath)
	if err != nil || fs == nil {
		ctx.Response().Abort(http.StatusNotFound, fmt.Errorf("文件不存在%s", rpath))
		return false
	}

	//写入到响应流
	if strings.HasSuffix(p, ".gz") {
		ctx.Response().Header("Content-Encoding", "gzip")
	}
	ctx.Response().File(p, fs)
	return true
}
