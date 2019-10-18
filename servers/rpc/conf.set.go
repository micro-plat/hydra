package rpc

import (
	"fmt"

	"github.com/asaskevich/govalidator"
	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/servers"
	"github.com/micro-plat/hydra/servers/pkg/middleware"
	"github.com/micro-plat/lib4go/types"
)

type ISetMetric interface {
	SetMetric(*conf.Metric) error
}

//SetMetric 设置metric
func SetMetric(set ISetMetric, cnf conf.IServerConf) (enable bool, err error) {
	//设置静态文件路由
	var metric conf.Metric
	_, err = cnf.GetSubObject("metric", &metric)
	if err != nil && err != conf.ErrNoSetting {
		return false, err
	}
	if err == conf.ErrNoSetting {
		metric.Disable = true
	} else {
		if b, err := govalidator.ValidateStruct(&metric); !b {
			err = fmt.Errorf("metric配置有误:%v", err)
			return false, err
		}
	}
	err = set.SetMetric(&metric)
	return !metric.Disable && err == nil, err
}

type ISetStatic interface {
	SetStatic(static *conf.Static) error
}

//SetStatic 设置static
func SetStatic(set ISetStatic, cnf conf.IServerConf) (enable bool, err error) {
	//设置静态文件路由
	var static conf.Static
	_, err = cnf.GetSubObject("static", &static)
	if err != nil && err != conf.ErrNoSetting {
		return false, err
	}
	if err == conf.ErrNoSetting {
		static.Disable = true
	} else {
		if b, err := govalidator.ValidateStruct(&static); !b {
			err = fmt.Errorf("static配置有误:%v", err)
			return false, err
		}
	}
	err = set.SetStatic(&static)
	return !static.Disable && err == nil, err
}

//ISetRouterHandler 设置路由列表
type ISetRouterHandler interface {
	SetRouters([]*conf.Router) error
}

func getRouters(services map[string][]string) conf.Routers {
	routers := conf.Routers{}

	if len(services) == 0 {
		routers.Routers = make([]*conf.Router, 0, 1)
		routers.Routers = append(routers.Routers, &conf.Router{Action: []string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS"}, Name: "/*name", Service: "/@name", Engine: "*"})
		return routers
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
	return routers
}
func SetRouters(engine servers.IRegistryEngine, cnf conf.IServerConf, set ISetRouterHandler, ext map[string]interface{}) (enable bool, err error) {
	var routers conf.Routers
	if _, err = cnf.GetSubObject("router", &routers); err == conf.ErrNoSetting || len(routers.Routers) == 0 {
		// routers = conf.Routers{}
		// routers.Routers = make([]*conf.Router, 0, 1)
		// routers.Routers = append(routers.Routers, &conf.Router{Action: []string{"GET", "POST", "PUT", "DELETE", "HEAD"}, Name: "/*name", Service: "/@name", Engine: "*"})
		routers = getRouters(engine.GetServices())
	}

	if err != conf.ErrNoSetting && err != nil {
		err = fmt.Errorf("路由:%v", err)
		return false, err
	}
	if b, err := govalidator.ValidateStruct(&routers); !b {
		err = fmt.Errorf("router配置有误:%v", err)
		return false, err
	}
	for _, router := range routers.Routers {
		if len(router.Action) == 0 {
			router.Action = []string{"GET", "POST", "PUT", "DELETE", "HEAD"}
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
		router.Handler = middleware.ContextHandler(engine, router.Name, router.Engine, router.Service, router.Setting, types.Copy(ext, "path", router.Service))
	}
	err = set.SetRouters(routers.Routers)
	return len(routers.Routers) > 0 && err == nil, err
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

//---------------------------------------------------------------------------
//-------------------------------header---------------------------------------
//---------------------------------------------------------------------------

//ISetHeaderHandler 设置header
type ISetHeaderHandler interface {
	SetHeader(conf.Headers) error
}

//SetHeaders 设置header
func SetHeaders(set ISetHeaderHandler, cnf conf.IServerConf) (enable bool, err error) {
	//设置通用头信息
	var header conf.Headers
	_, err = cnf.GetSubObject("header", &header)
	if err != nil && err != conf.ErrNoSetting {
		err = fmt.Errorf("header配置有误:%v", err)
		return false, err
	}
	err = set.SetHeader(header)
	return len(header) > 0 && err == nil, err
}

//---------------------------------------------------------------------------
//-------------------------------host---------------------------------------
//---------------------------------------------------------------------------

//ISetHosts 设置hosts
type ISetHosts interface {
	SetHosts(conf.Hosts) error
}

//SetHosts 设置hosts
func SetHosts(set ISetHosts, cnf conf.IServerConf) (enable bool, err error) {
	var hosts conf.Hosts
	hosts = cnf.GetStrings("host")
	return len(hosts) > 0, set.SetHosts(hosts)
}

//---------------------------------------------------------------------------
//-------------------------------jwt---------------------------------------
//---------------------------------------------------------------------------

//ISetJwtAuth 设置jwt
type ISetJwtAuth interface {
	SetJWT(*conf.JWTAuth) error
}

//SetJWT 设置jwt
func SetJWT(set ISetJwtAuth, cnf conf.IServerConf) (enable bool, err error) {
	//设置jwt安全认证参数
	var auths conf.Authes
	if _, err := cnf.GetSubObject("auth", &auths); err != nil && err != conf.ErrNoSetting {
		err = fmt.Errorf("jwt配置有误:%v", err)
		return false, err
	}
	if auths.JWT != nil {
		if b, err := govalidator.ValidateStruct(auths.JWT); !b {
			err = fmt.Errorf("jwt配置有误:%v", err)
			return false, err
		}
		err = set.SetJWT(auths.JWT)
		return err == nil && !auths.JWT.Disable, err
	}
	return false, nil
}

//---------------------------------------------------------------------------
//-------------------------------fixed-secret---------------------------------------
//---------------------------------------------------------------------------

//CheckFixedSecret 设置FixedSecret
func CheckFixedSecret(cnf conf.IServerConf) (enable bool, err error) {
	//设置fixedSecret安全认证参数
	var auths conf.Authes
	if _, err := cnf.GetSubObject("auth", &auths); err != nil && err != conf.ErrNoSetting {
		err = fmt.Errorf("fixed-secret配置有误:%v", err)
		return false, err
	}
	if auths.FixedScret != nil {
		if b, err := govalidator.ValidateStruct(auths.FixedScret); !b {
			err = fmt.Errorf("fixed-secret配置有误:%v", err)
			return false, err
		}
		return !auths.FixedScret.Disable, nil
	}
	return false, nil
}

//---------------------------------------------------------------------------
//-------------------------------remote-auth---------------------------------------
//---------------------------------------------------------------------------

//CheckRemoteAuth 检查是否设置remote-auth
func CheckRemoteAuth(cnf conf.IServerConf) (enable bool, err error) {
	//设置Remote安全认证参数
	var auths conf.Authes
	if _, err := cnf.GetSubObject("auth", &auths); err != nil && err != conf.ErrNoSetting {
		err = fmt.Errorf("remote-auth配置有误:%v", err)
		return false, err
	}
	count := 0
	for _, auth := range auths.RemotingServiceAuths {
		if b, err := govalidator.ValidateStruct(auth); !b {
			err = fmt.Errorf("remote-auth配置有误:%v", err)
			return false, err
		}
		if !auth.Disable {
			count++
		}
	}
	return count > 0, nil
}
