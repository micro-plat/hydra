package engines

import (
	"fmt"
	"strings"
	"time"

	"github.com/micro-plat/hydra/component"
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/lib4go/types"
)

//RPCProxy rpc 代理服务
func (r *ServiceEngine) RPCProxy() component.ServiceFunc {
	return func(ctx *context.Context) (r interface{}) {
		header, _ := ctx.Request.Http.GetHeader()
		cookie, _ := ctx.Request.Http.GetCookies()
		for k, v := range cookie {
			header[k] = v
		}
		input := ctx.Request.GetRequestMap()
		timeout := ctx.Request.Setting.GetInt("timeout", 3)
		response := ctx.RPC.AsyncRequest(ctx.Service, strings.ToUpper(ctx.Request.GetMethod()), header, input, true)
		status, result, params, err := response.Wait(time.Second * time.Duration(timeout))
		if err != nil {
			err = fmt.Errorf("rpc.proxy %v(%d)", err, status)
		}
		ctx.Response.SetParams(types.GetIMap(params))
		if err != nil {
			ctx.Response.SetStatus(status)
			return err
		}
		ctx.Response.SetStatus(status)
		return result
	}
}
