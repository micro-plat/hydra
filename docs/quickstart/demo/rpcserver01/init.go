package main

import (
	"github.com/micro-plat/hydra/component"
	"github.com/micro-plat/hydra/quickstart/demo/rpcserver01/services/order"
)

//init 检查应用程序配置文件，并根据配置初始化服务
func (rpc *rpcserver) init() {
	rpc.config()
	rpc.handling()

	rpc.Initializing(func(c component.IContainer) error {
		//检查db配置是否正确
		// if _, err := c.GetDB(); err != nil {
		// 	return err
		// }

		return nil
	})

	//服务注册
	rpc.Micro("/order", order.NewOrderHandler)
}
