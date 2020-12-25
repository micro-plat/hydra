package internal

import (
	"fmt"

	"github.com/micro-plat/hydra/components"
	"github.com/micro-plat/hydra/context"
)

//CallRPC RPC调用
func CallRPC(ctx context.IContext, service string) interface{} {
	response, err := components.Def.RPC().GetRegularRPC().Swap(service, ctx)
	if err != nil {
		ctx.Log().Errorf("调用RPC服务出错:%+v", err)
		return err
	}
	if !response.IsSuccess() {
		return fmt.Errorf(response.Result)
	}
	if v, err := response.GetResult(); err == nil {
		return v
	}
	return response.Result
}
