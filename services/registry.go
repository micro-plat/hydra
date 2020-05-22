package services

import (
	"errors"
	"fmt"
	"strings"

	"github.com/micro-plat/hydra/application"
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/registry/conf/server"
)

const defHandling = "Handling"
const defHandler = "Handle"
const defHandled = "Handled"
const defFallback = "Fallback"
const defClose = "Close"

//IServiceRegistry 服务注册接口
type IServiceRegistry interface {
	Micro(name string, h interface{}, pages ...string)
	Flow(name string, h interface{})
	API(name string, h interface{}, pages ...string)
	Web(name string, h interface{}, pages ...string)
	RPC(name string, h interface{}, pages ...string)
	WS(name string, h interface{}, pages ...string)
	MQC(name string, h interface{}, queues ...string)
	CRON(name string, h interface{}, crons ...string)
	OnServerStarting(h func(server.IServerConf) error, tps ...string)
	OnServerClosing(h func(server.IServerConf) error, tps ...string)
	OnHandleExecuting(h context.Handler, tps ...string)
	OnHandleExecuted(h context.Handler, tps ...string)
}

//Registry 服务注册管理
var Registry = &regist{
	servers: make(map[string]*serverServices),
	caches:  make(map[string]map[string]interface{}),
}

//regist  本地服务
type regist struct {
	servers map[string]*serverServices
	caches  map[string]map[string]interface{}
}

//Micro 注册为微服务包括api,web,rpc
func (s *regist) Micro(name string, h interface{}, pages ...string) {
	s.get(application.API).Add(name, h, pages...)
	s.get(application.Web).Add(name, h, pages...)
	s.get(application.WS).Add(name, h, pages...)
}

//Flow 注册为流程服务，包括mqc,cron
func (s *regist) Flow(name string, h interface{}) {
	s.get(application.MQC).Add(name, h)
	s.get(application.CRON).Add(name, h)
}

//API 注册为API服务
func (s *regist) API(name string, h interface{}, pages ...string) {
	s.get(application.API).Add(name, h, pages...)
}

//Web 注册为web服务
func (s *regist) Web(name string, h interface{}, pages ...string) {
	s.get(application.Web).Add(name, h, pages...)
}

//RPC 注册为rpc服务
func (s *regist) RPC(name string, h interface{}, pages ...string) {
	s.get(application.WS).Add(name, h, pages...)
}

//WS 注册为websocket服务
func (s *regist) WS(name string, h interface{}, pages ...string) {
	s.get(application.WS).Add(name, h, pages...)
}

//MQC 注册为消息队列服务
func (s *regist) MQC(name string, h interface{}, queues ...string) {
	s.get(application.MQC).Add(name, h, queues...)
}

//CRON 注册为定时任务服务
func (s *regist) CRON(name string, h interface{}, crons ...string) {
	s.get(application.CRON).Add(name, h, crons...)
}

//OnServerStarting 处理服务器启动
func (s *regist) OnServerStarting(h func(server.IServerConf) error, tps ...string) {
	if len(tps) == 0 {
		tps = application.DefApp.ServerTypes
	}
	for _, typ := range tps {
		if err := s.get(typ).AddStarting(h); err != nil {
			panic(fmt.Errorf("%s OnServerStarting %v", typ, err))
		}
	}
}

//OnServerClosing 处理服务器关闭
func (s *regist) OnServerClosing(h func(server.IServerConf) error, tps ...string) {
	if len(tps) == 0 {
		tps = application.DefApp.ServerTypes
	}
	for _, typ := range tps {
		if err := s.get(typ).AddClosing(h); err != nil {
			panic(fmt.Errorf("%s OnServerClosing %v", typ, err))
		}
	}
}

//OnHandleExecuting 处理handling业务
func (s *regist) OnHandleExecuting(h context.Handler, tps ...string) {
	if len(tps) == 0 {
		tps = application.DefApp.ServerTypes
	}
	for _, typ := range tps {
		if err := s.get(typ).AddHandleExecuting(h); err != nil {
			panic(fmt.Errorf("%s OnHandleExecuting %v", typ, err))
		}
	}
}

//Handled 处理Handled业务
func (s *regist) OnHandleExecuted(h context.Handler, tps ...string) {
	if len(tps) == 0 {
		tps = application.DefApp.ServerTypes
	}
	for _, typ := range tps {
		if err := s.get(typ).AddHandleExecuted(h); err != nil {
			panic(fmt.Errorf("%s OnHandleExecuted %v", typ, err))
		}
	}
}

//GetHandleExecutings 获取handle预处理勾子
func (s *regist) GetHandleExecutings(serverType string) []context.IHandler {
	return s.get(serverType).GetHandleExecutings()
}

//GetHandleExecuted 获取handle后处理勾子
func (s *regist) GetHandleExecuted(serverType string) []context.IHandler {
	return s.get(serverType).GetHandleExecuteds()
}

//GetHandler 获取服务对应的处理函数
func (s *regist) GetHandler(serverType string, service string) (context.IHandler, bool) {
	return s.get(serverType).GetHandlers(service)
}

//GetHandling 获取预处理函数
func (s *regist) GetHandlings(serverType string, service string) []context.IHandler {
	return s.get(serverType).GetHandlings(service)
}

//GetHandling 获取后处理函数
func (s *regist) GetHandleds(serverType string, service string) []context.IHandler {
	return s.get(serverType).GetHandleds(service)
}

//GetFallback 获取服务对应的降级函数
func (s *regist) GetFallback(serverType string, service string) (context.IHandler, bool) {
	return s.get(serverType).GetFallback(service)
}

func (s *regist) get(tp string) *serverServices {
	if v, ok := s.servers[tp]; ok {
		return v
	}
	panic(fmt.Sprintf("不支持的服务器类型:%s", tp))
}

//OnStarting 执行服务启动函数
func (s *regist) OnStarting(c server.IServerConf) error {
	return s.get(c.GetMainConf().GetServerType()).OnStarting(c)

}

//OnClose 执行服务关闭函数
func (s *regist) OnClosing(c server.IServerConf) error {
	return s.get(c.GetMainConf().GetServerType()).OnClosing(c)
}

//GetClosers 获取资源释放函数
func (s *regist) Close() error {
	var sb strings.Builder
	for _, r := range s.servers {
		for _, cs := range r.GetClosingHandlers() {
			if err := cs(); err != nil {
				sb.WriteString(err.Error())
				sb.WriteString("\n")
			}
		}
	}
	if sb.Len() == 0 {
		return nil
	}
	return errors.New(strings.Trim(sb.String(), "\n"))
}

//-----------------------注册缓存-------------------------------------------

//init 处理服务初始化及特殊注册函数
func init() {
	Registry.servers[application.API] = newServerServices(func(g *unit, ext ...string) error {
		return API.Add(g.path, g.service, g.actions, ext...)
	})
	Registry.servers[application.Web] = newServerServices(func(g *unit, ext ...string) error {
		return WEB.Add(g.path, g.service, g.actions, ext...)
	})
	Registry.servers[application.RPC] = newServerServices(func(g *unit, ext ...string) error {
		return RPC.Add(g.path, g.service, g.actions, ext...)
	})

	Registry.servers[application.WS] = newServerServices(nil)

	Registry.servers[application.CRON] = newServerServices(func(g *unit, ext ...string) error {
		for _, t := range ext {
			CRON.Add(t, g.service)
		}
		return nil

	})
	Registry.servers[application.MQC] = newServerServices(func(g *unit, ext ...string) error {
		for _, t := range ext {
			MQC.Add(t, g.service)
		}
		return nil
	})
}
