// +build !prod

package main

func (api *apiserver) config() {
	api.IsDebug = true
	api.Conf.API.SetMainConf(`{"address":":8090","trace":true}`)
	api.Conf.Plat.SetVarConf("db", "db", `{			
			"provider":"mysql",
			"connString":"mrss:123456@tcp(192.168.0.36)/mrss?charset=utf8",
			"maxOpen":20,
			"maxIdle":10,
			"lifeTime":600		
	}`)
}
