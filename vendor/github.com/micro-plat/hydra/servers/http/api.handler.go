package http

import (
	"fmt"
	x "net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/servers"
	"github.com/micro-plat/hydra/servers/http/middleware"
)

func (s *ApiServer) getHandler(routers []*conf.Router) (x.Handler, error) {
	if !servers.IsDebug {
		gin.SetMode(gin.ReleaseMode)
	}
	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(middleware.Logging(s.conf)) //记录请求日志
	engine.Use(middleware.Recovery())
	engine.Use(s.option.metric.Handle())        //生成metric报表
	engine.Use(middleware.Host(s.conf))         // 检查主机头是否合法
	engine.Use(middleware.Static(s.conf))       //处理静态文件
	engine.Use(middleware.AjaxRequest(s.conf))  //过滤非ajax请求
	engine.Use(middleware.JwtAuth(s.conf))      //jwt安全认证
	engine.Use(middleware.CircuitBreak(s.conf)) //服务熔断配置
	//engine.Use(middleware.Body())               //处理请求form
	engine.Use(middleware.APIResponse(s.conf)) //处理返回值
	engine.Use(middleware.Header(s.conf))      //设置请求头
	engine.Use(middleware.JwtWriter(s.conf))   //设置jwt回写
	err := setRouters(engine, routers)
	return engine, err
}
func setRouters(engine *gin.Engine, routers []*conf.Router) (err error) {
	defer func() {
		if err1 := recover(); err1 != nil {
			err = fmt.Errorf("%v", err1)
		}
	}()
	for _, router := range routers {
		for _, method := range router.Action {
			engine.Handle(strings.ToUpper(method), router.Name, router.Handler.(gin.HandlerFunc))
		}
	}
	return nil
}

type Routers struct {
	routers []*conf.Router
}

func GetRouters() *Routers {
	return &Routers{
		routers: make([]*conf.Router, 0, 2),
	}

}
func (r *Routers) Get() []*conf.Router {
	return r.routers
}
func (r *Routers) Route(method string, name string, f interface{}) {
	r.routers = append(r.routers,
		&conf.Router{
			Name:    name,
			Action:  strings.Split(method, ","),
			Engine:  "*",
			Service: name,
			Handler: middleware.ContextHandler(f, name, "*", name, nil), //??
		})
}
