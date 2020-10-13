package main

import (
	"github.com/micro-plat/hydra"
	"github.com/micro-plat/hydra/conf/server/api"
 
)

func init() {
	hydra.OnReady(func() {
		hydra.Conf.API(":8081", api.WithHeaderReadTimeout(30), api.WithTimeout(30, 30)) 

	})
}
