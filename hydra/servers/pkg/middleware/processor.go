package middleware

import "net/http"

//Processor 请求处理
func Processor() Handler {
	return func(ctx IMiddleContext) {
		//1. 获取jwt配置
		processorObj, err := ctx.APPConf().GetProcessorConf()
		if err != nil {
			ctx.Response().Abort(http.StatusNotExtended, err)
			return
		}
		if processorObj.ServicePrefix == "" {
			ctx.Next()
			return
		}
		ctx.ServicePrefix(processorObj.ServicePrefix)
		ctx.Next()
	}
}
