package ws

import (
	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/servers"
	"github.com/micro-plat/hydra/servers/http/middleware"
)

//ISetRouterHandler 设置路由列表
type ISetRouterHandler interface {
	SetRouters([]*conf.Router) error
}

//SetHttpRouters 设置路由
func SetHttpRouters(engine servers.IRegistryEngine, set ISetRouterHandler, cnf conf.IServerConf) (enable bool, err error) {
	var routers conf.Routers
	routers = conf.Routers{}
	routers.Routers = make([]*conf.Router, 0, 1)
	routers.Routers = append(routers.Routers, &conf.Router{Action: []string{"GET"}, Name: "/*name", Service: "/@name", Engine: "*"})

	for _, router := range routers.Routers {
		if len(router.Action) == 0 {
			router.Action = []string{"GET"}
		}
		if router.Engine == "" {
			router.Engine = "*"
		}
		if router.Setting == nil {
			router.Setting = make(map[string]string)
		}
		for k, v := range routers.Setting {
			if _, ok := router.Setting[k]; !ok {
				router.Setting[k] = v
			}
		}
		router.Handler = middleware.WSContextHandler(engine, router.Name, router.Engine, router.Service, router.Setting)
	}
	err = set.SetRouters(routers.Routers)
	return len(routers.Routers) > 0 && err == nil, err
}
