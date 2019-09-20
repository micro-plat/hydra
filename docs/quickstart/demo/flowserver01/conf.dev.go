package main

func (flow *flowserver) config() {
	flow.IsDebug = true

	flow.Conf.MQC.SetSubConf("server", `
		{
			"proto":"redis",
			"addrs":[
					"192.168.0.111:6379",
					"192.168.0.112:6379",
					"192.168.0.113:6379",
					"192.168.0.114:6379"
			],
			"db":1,
			"dial_timeout":10,
			"read_timeout":10,
			"write_timeout":10,
			"pool_size":10
	}
	`)
	// flow.Conf.MQC.SetSubConf("queue", `{
	//     "queues":[
	//         {
	//             "queue":"coupon:base:coupon_produce",
	//             "service":"/coupon/produce"
	// 		},
	// 		{
	//             "queue":"coupon:base:down_payment",
	//             "service":"/order/pay"
	// 		}
	//     ]
	// }`)
}
