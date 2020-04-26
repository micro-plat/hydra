package component

import (
	"fmt"
	"strings"

	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/lib4go/types"
)

//NewRPCCtxSerivce 构建RPC转发请求
func NewRPCCtxSerivce(rpcServiceName string, rpcInput ...func(ctx *context.Context) interface{}) ServiceFunc {
	return func(ctx *context.Context) (rs interface{}) {
		header, _ := ctx.Request.Http.GetHeader()
		cookie, _ := ctx.Request.Http.GetCookies()
		for k, v := range cookie {
			header[k] = v
		}
		header["method"] = strings.ToUpper(ctx.Request.GetMethod())
		nheader := types.NewXMapBySMap(header)
		input := types.NewXMapByMap(ctx.Request.GetRequestMap())
		switch {
		case len(rpcInput) == 1:
			value := rpcInput[0](ctx)
			switch v := value.(type) {
			case context.IError, error:
				return value
			case map[string]string:
				input.MergeSMap(v)
			case map[string]interface{}:
				input.MergeMap(v)
			default:
				return fmt.Errorf("执NewRPCCtxSerivce服务返回的类型只支持map[string]string,map[string]interface{}")
			}

		}

		status, result, params, err := ctx.RPC.Request(rpcServiceName, nheader.ToSMap(), input.ToMap(), true)
		if err != nil {
			return err
		}
		ctx.Response.SetParams(types.GetIMap(params))
		if status != 200 {
			return context.NewError(status, result)
		}
		ctx.Response.SetJSON()
		ctx.Response.MustContent(status, result)
		return
	}

}

//NewRPCSerivce 构建RPC转发请求
func NewRPCSerivce(rpcServiceName string, rpcInput ...map[string]string) ServiceFunc {
	return func(ctx *context.Context) (rs interface{}) {
		header, _ := ctx.Request.Http.GetHeader()
		cookie, _ := ctx.Request.Http.GetCookies()
		for k, v := range cookie {
			header[k] = v
		}
		header["method"] = strings.ToUpper(ctx.Request.GetMethod())
		nheader := types.NewXMapBySMap(header)
		input := types.NewXMapByMap(ctx.Request.GetRequestMap())
		switch {
		case len(rpcInput) == 1:
			input.MergeSMap(rpcInput[0])
		case len(rpcInput) >= 2:
			nheader.MergeSMap(rpcInput[0])
			input.MergeSMap(rpcInput[1])
		}

		status, result, params, err := ctx.RPC.Request(rpcServiceName, nheader.ToSMap(), input.ToMap(), true)
		if err != nil {
			return err
		}
		ctx.Response.SetParams(types.GetIMap(params))
		if status != 200 {
			return context.NewError(status, result)
		}
		ctx.Response.SetJSON()
		ctx.Response.MustContent(status, result)
		return
	}

}
