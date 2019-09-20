package main

import (
	"github.com/micro-plat/hydra/component"
	"github.com/micro-plat/hydra/quickstart/demo/apiserver10/services/order"
)

//init 检查应用程序配置文件，并根据配置初始化服务
func (api *apiserver) init() {
	api.Initializing(func(c component.IContainer) error {
		//检查db配置是否正确
		if _, err := c.GetDB(); err != nil {
			return err
		}
		//检查消息队列配置

		//拉取应用程序配置

		return nil
	})

	//服务注册
	api.Micro("/order/request", order.NewRequestHandler)

}
