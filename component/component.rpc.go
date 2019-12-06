package component

import (
	"strings"

	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/lib4go/types"
)

//NewRPCSerivce 构建RPC转发请求
func NewRPCSerivce(rpcServiceName string, rpcInput ...map[string]string) ServiceFunc {
	return func(ctx *context.Context) (rs interface{}) {
		header, _ := ctx.Request.Http.GetHeader()
		cookie, _ := ctx.Request.Http.GetCookies()
		for k, v := range cookie {
			header[k] = v
		}
		header["method"] = strings.ToUpper(ctx.Request.GetMethod())
		input := types.NewXMapByMap(ctx.Request.GetRequestMap())
		if len(rpcInput) > 0 {
			input.MergeSMap(rpcInput[0])
		}

		status, result, params, err := ctx.RPC.Request(rpcServiceName, header, input.ToMap(), true)
		if err != nil {
			return err
		}
		ctx.Response.SetParams(types.GetIMap(params))
		if status != 200 {
			return context.NewError(status, result)
		}
		ctx.Response.MustContent(status, result)
		return
	}

}
