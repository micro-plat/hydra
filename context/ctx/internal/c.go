package internal

import "github.com/micro-plat/hydra/context"
import "github.com/micro-plat/hydra/components"
import "fmt"

//CallRPC RPC调用
func CallRPC(ctx context.IContext, service string) interface{} {
	response, err := components.Def.RPC().GetRegularRPC().Swap(service, ctx)
	if err != nil {
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
