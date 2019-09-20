// +build !prod

package main

func (rpc *rpcserver) config() {
	rpc.IsDebug = true
	rpc.Conf.RPC.SetMainConf(`{"address":":9090","trace":true}`)
	rpc.Conf.Plat.SetVarConf("db", "db", `{			
			"provider":"mysql",
			"connString":"mrss:123456@tcp(192.168.0.36)/mrss?charset=utf8",
			"maxOpen":20,
			"maxIdle":10,
			"lifeTime":600		
	}`)

}
