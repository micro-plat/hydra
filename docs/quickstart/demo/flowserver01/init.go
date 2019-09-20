package main

import "github.com/micro-plat/hydra/conf"

func (flow *flowserver) init() {
	flow.config()

	ch := flow.GetDynamicQueue()

	//订阅消息
	ch <- &conf.Queue{Queue: "mall:flow:order_pay", Service: "/order/pay"}

	//取消订阅
	// ch<- &cf.Queue{Queue: "mall:flow:order_pay",Disable:true}
}
