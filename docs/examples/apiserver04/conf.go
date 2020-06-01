package main

import (
	"github.com/micro-plat/hydra"
)

func init() {
	hydra.Conf.OnReady(func() {
		hydra.Conf.API(":8080").Metric("http://192.168.106.219:8086", "hydra", "@every 5s")
	})
}
