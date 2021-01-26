package services

import (
	"errors"
	"fmt"
	"strings"

	"github.com/micro-plat/hydra/conf/app"
	"github.com/micro-plat/hydra/conf/server/router"
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/global"
)

const defHandling = "Handling"
const defHandler = "Handle"
const defHandled = "Handled"
const defFallback = "Fallback"
const defClose = "Close"

//IService 服务注册接口
type IService interface {
	//Micro 注册为微服务，包括api,web,rpc,websocket
	Micro(name string, h interface{}, r ...router.Option)

	//Flow 注册为流程,包括mqc,cron
	Flow(name string, h interface{})

	//API 注册为http api服务
	API(name string, h interface{}, r ...router.Option)

	//Web 注册为web服务
	Web(name string, h interface{}, r ...router.Option)

	//RPC 注册为RPC服务
	RPC(name string, h interface{}, r ...router.Option)

	//WS 注册为websocket服务
	WS(name string, h interface{}, r ...router.Option)
	//MQC 注册为消息消费服务
	MQC(name string, h interface{}, queues ...string)

	//CRON 注册为定时任务服务
	CRON(name string, h interface{}, crons ...string)

	//custome 注册为自定义服务器的服务
	Custom(tp string, name string, h interface{}, ext ...interface{})

	//RegisterServer 注册新的服务器类型
	RegisterServer(tp string, f ...func(g *Unit, ext ...interface{}) error)

	//OnSetup 服务器初始化勾子，服务器配置初始化前执行
	OnSetup(h func(app.IAPPConf) error, tps ...string)

	//OnStarting 服务器启动前勾子，服务器启动前执行
	OnStarting(h func(app.IAPPConf) error, tps ...string)

	//OnStarted 服务器启动完成勾子，服务器启动后执行
	OnStarted(h func(app.IAPPConf) error, tps ...string)

	//OnClosing 服务器关闭勾子，服务器关闭后执行
	OnClosing(h func(app.IAPPConf) error, tps ...string)

	//OnHandleExecuting Handle勾子，Handle执行前执行
	OnHandleExecuting(h context.Handler, tps ...string)

	//OnHandleExecuted Handle勾子，Handle执行后执行
	OnHandleExecuted(h context.Handler, tps ...string)
}

//Def 服务注册管理
var Def = New()

//New 构建服务组件
func New() *regist {
	return &regist{
		servers: make(map[string]*serverServices),
		caches:  make(map[string]map[string]interface{}),
	}
}

//regist  本地服务
type regist struct {
	servers map[string]*serverServices
	caches  map[string]map[string]interface{}
}

//Micro 注册为微服务包括api,web,rpc
func (s *regist) Micro(name string, h interface{}, r ...router.Option) {
	s.API(name, h, r...)
	s.Web(name, h, r...)
	s.WS(name, h, r...)
	s.RPC(name, h, r...)
}

//Flow 注册为流程服务，包括mqc,cron
func (s *regist) Flow(name string, h interface{}) {
	s.Custom(global.MQC, name, h)
	s.Custom(global.CRON, name, h)
}

//API 注册为API服务
func (s *regist) API(name string, h interface{}, ext ...router.Option) {
	v := make([]interface{}, 0, len(ext))
	for _, e := range ext {
		v = append(v, e)
	}
	s.Custom(global.API, name, h, v...)
}

//Web 注册为web服务
func (s *regist) Web(name string, h interface{}, ext ...router.Option) {
	v := make([]interface{}, 0, len(ext))
	for _, e := range ext {
		v = append(v, e)
	}
	s.Custom(global.Web, name, h, v...)
}

//RPC 注册为rpc服务
func (s *regist) RPC(name string, h interface{}, ext ...router.Option) {
	v := make([]interface{}, 0, len(ext))
	for _, e := range ext {
		v = append(v, e)
	}
	s.Custom(global.RPC, name, h, v...)
}

//WS 注册为websocket服务
func (s *regist) WS(name string, h interface{}, ext ...router.Option) {
	v := make([]interface{}, 0, len(ext))
	for _, e := range ext {
		v = append(v, e)
	}
	s.Custom(global.WS, name, h, v...)
}

//MQC 注册为消息队列服务
func (s *regist) MQC(name string, h interface{}, queues ...string) {
	v := make([]interface{}, 0, len(queues))
	for _, e := range queues {
		v = append(v, e)
	}
	s.Custom(global.MQC, name, h, v...)
}

//CRON 注册为定时任务服务
func (s *regist) CRON(name string, h interface{}, crons ...string) {
	v := make([]interface{}, 0, len(crons))
	for _, e := range crons {
		v = append(v, e)
	}
	s.Custom(global.CRON, name, h, v...)
}

//Custom 自定义服务注册
func (s *regist) Custom(tp string, name string, h interface{}, ext ...interface{}) {
	//去掉两头'/'
	name = fmt.Sprintf("/%s", strings.Trim(name, "/"))
	s.get(tp).Register(name, h, ext...)
}

//RegisterServer 注册服务器
func (s *regist) RegisterServer(tp string, f ...func(g *Unit, ext ...interface{}) error) {
	if _, ok := s.servers[tp]; ok {
		panic(fmt.Errorf("服务%s已存在，不能重复注册", tp))
	}
	if len(f) > 0 {
		s.servers[tp] = newServerServices(f[0])
		return
	}
	s.servers[tp] = newServerServices(nil)
}

//OnSetup 服务器初始化勾子，服务器配置初始化前执行
func (s *regist) OnSetup(h func(app.IAPPConf) error, tps ...string) {
	if len(tps) == 0 {
		tps = global.Def.ServerTypes
	}
	for _, typ := range tps {
		if err := s.get(typ).AddSetup(h); err != nil {
			panic(fmt.Errorf("%s OnSetup %v", typ, err))
		}
	}
}

//OnStarting 处理服务器启动前
func (s *regist) OnStarting(h func(app.IAPPConf) error, tps ...string) {
	if len(tps) == 0 {
		tps = global.Def.ServerTypes
	}
	for _, typ := range tps {
		if err := s.get(typ).AddStarting(h); err != nil {
			panic(fmt.Errorf("%s OnStarting %v", typ, err))
		}
	}
}

//OnStarted 处理服务器启动后
func (s *regist) OnStarted(h func(app.IAPPConf) error, tps ...string) {
	if len(tps) == 0 {
		tps = global.Def.ServerTypes
	}
	for _, typ := range tps {
		if err := s.get(typ).AddStarted(h); err != nil {
			panic(fmt.Errorf("%s OnStarted %v", typ, err))
		}
	}
}

//OnClosing 处理服务器关闭
func (s *regist) OnClosing(h func(app.IAPPConf) error, tps ...string) {
	if len(tps) == 0 {
		tps = global.Def.ServerTypes
	}
	for _, typ := range tps {
		if err := s.get(typ).AddClosing(h); err != nil {
			panic(fmt.Errorf("%s OnClosing %v", typ, err))
		}
	}
}

//OnHandleExecuting 处理handling业务
func (s *regist) OnHandleExecuting(h context.Handler, tps ...string) {
	if len(tps) == 0 {
		tps = global.Def.ServerTypes
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
		tps = global.Def.ServerTypes
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

//GetRawPathAndTag 获取服务原始注册路径与方法名(restful服务的tag值为空)
func (s *regist) GetRawPathAndTag(serverType string, service string) (path string, tagName string, ok bool) {
	return s.get(serverType).GetRawPathAndTag(service)
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

//DoStarting 执行服务启动函数
func (s *regist) DoStarting(c app.IAPPConf) error {
	return s.get(c.GetServerConf().GetServerType()).DoStarting(c)

}

//DoStarted 执行服务启动函数
func (s *regist) DoStarted(c app.IAPPConf) error {
	return s.get(c.GetServerConf().GetServerType()).DoStarted(c)

}

//DoSetup 执行服务启动函数
func (s *regist) DoSetup(c app.IAPPConf) error {
	return s.get(c.GetServerConf().GetServerType()).DoSetup(c)

}

//DoClosing 执行服务关闭函数
func (s *regist) DoClosing(c app.IAPPConf) error {
	return s.get(c.GetServerConf().GetServerType()).DoClosing(c)
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
	Def.servers[global.API] = newServerServices(func(g *Unit, ext ...interface{}) error {
		return API.Add(g.Path, g.Service, g.Actions, ext...)
	})
	Def.servers[global.Web] = newServerServices(func(g *Unit, ext ...interface{}) error {
		return WEB.Add(g.Path, g.Service, g.Actions, ext...)
	})
	Def.servers[global.RPC] = newServerServices(func(g *Unit, ext ...interface{}) error {
		return RPC.Add(g.Path, g.Service, g.Actions, ext...)
	})

	Def.servers[global.WS] = newServerServices(func(g *Unit, ext ...interface{}) error {
		return WS.Add(g.Path, g.Service, g.Actions, ext...)
	})
	Def.servers[global.CRON] = newServerServices(func(g *Unit, ext ...interface{}) error {
		for _, t := range ext {
			CRON.Static(t.(string), g.Service)
		}
		return nil
	})
	Def.servers[global.MQC] = newServerServices(func(g *Unit, ext ...interface{}) error {
		for _, t := range ext {
			MQC.Static(t.(string), g.Service)
		}
		return nil
	})
}
