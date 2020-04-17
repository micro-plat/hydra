package http

import (
	"fmt"

	"github.com/asaskevich/govalidator"
	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/servers"
	"github.com/micro-plat/hydra/servers/http/middleware"
	"github.com/micro-plat/lib4go/types"
)

//waitRemoveDir 等待移除的静态文件
var waitRemoveDir = make([]string, 0, 1)

//ISetRouterHandler 设置路由列表
type ISetRouterHandler interface {
	SetRouters([]*conf.Router) error
}

func getDefaultRouters(services map[string][]string) []*conf.Router {
	routers := conf.Routers{}

	if len(services) == 0 {
		routers.Routers = make([]*conf.Router, 0, 1)
		routers.Routers = append(routers.Routers, &conf.Router{Action: []string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS"}, Name: "/*name", Service: "/@name", Engine: "*"})
		return routers.Routers
	}
	routers.Routers = make([]*conf.Router, 0, len(services))
	for name, actions := range services {
		router := &conf.Router{
			Action:  actions, //[]string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS"},
			Name:    name,
			Service: name,
			Engine:  "*",
		}
		router.Action = append(router.Action, "OPTIONS")
		routers.Routers = append(routers.Routers, router)
	}
	return routers.Routers
}

//SetHttpRouters 设置路由
func SetHttpRouters(engine servers.IRegistryEngine, set ISetRouterHandler, cnf conf.IServerConf) (enable bool, err error) {

	var routers conf.Routers
	if _, err = cnf.GetSubObject("router", &routers); err == conf.ErrNoSetting || len(routers.Routers) == 0 {
		routers.Routers = getDefaultRouters(engine.GetServices()) //添加默认路由
	}
	if err != nil && err != conf.ErrNoSetting {
		err = fmt.Errorf("路由:%v", err)
		return false, err
	}
	if b, err := govalidator.ValidateStruct(&routers); !b {
		err = fmt.Errorf("router配置有误:%v", err)
		return false, err
	}

	//处理RPC代理服务
	nRouters := make([]*conf.Router, 0, len(routers.RPCS)+len(routers.Routers))
	for _, proxy := range routers.RPCS {
		if len(proxy.Action) == 0 {
			proxy.Action = []string{"GET", "POST"}
		}
		proxy.Engine = "rpc"
		nRouters = append(nRouters, proxy)
	}

	//处理路由
	for _, router := range routers.Routers {
		if len(router.Action) == 0 {
			router.Action = []string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS"}
		}
		router.Engine = types.GetString(router.Engine, "*")
		nRouters = append(nRouters, router)
	}

	//构建路由处理函数
	for _, router := range nRouters {
		if router.Setting == nil {
			router.Setting = make(map[string]string)
		}
		for k, v := range routers.Setting {
			if _, ok := router.Setting[k]; !ok {
				router.Setting[k] = v
			}
		}
		router.Handler = middleware.ContextHandler(engine, router.Name, router.Engine, router.Service, router.Setting)

	}
	err = set.SetRouters(nRouters)
	return len(nRouters) > 0 && err == nil, err
}

//ISetCircuitBreaker 设置CircuitBreaker
type ISetCircuitBreaker interface {
	CloseCircuitBreaker() error
	SetCircuitBreaker(*conf.CircuitBreaker) error
}

//SetCircuitBreaker 设置熔断配置
func SetCircuitBreaker(set ISetCircuitBreaker, cnf conf.IServerConf) (enable bool, err error) {
	//设置CircuitBreaker
	var breaker conf.CircuitBreaker
	if _, err = cnf.GetSubObject("circuit", &breaker); err == conf.ErrNoSetting || breaker.Disable {
		return false, set.CloseCircuitBreaker()
	}
	if err != nil {
		return false, err
	}
	if b, err := govalidator.ValidateStruct(&breaker); !b {
		err = fmt.Errorf("circuit配置有误:%v", err)
		return false, err
	}
	err = set.SetCircuitBreaker(&breaker)
	return err == nil && !breaker.Disable, err
}
