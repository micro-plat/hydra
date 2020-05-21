package services

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/micro-plat/hydra/application"
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/registry/conf/server"
	"github.com/micro-plat/hydra/registry/conf/server/router"
)

const defHandling = "Handling"
const defHandler = "Handle"
const defHandled = "Handled"
const defFallback = "Fallback"
const defClose = "Close"

//IServiceRegistry 服务注册接口
type IServiceRegistry interface {
	Micro(name string, h interface{})
	Flow(name string, h interface{})
	API(name string, h interface{})
	Web(name string, h interface{})
	WS(name string, h interface{})
	MQC(name string, h interface{})
	CRON(name string, h interface{})
	OnServerStarting(h func(server.IServerConf) error, tps ...string)
	OnServerClosing(h func(server.IServerConf) error, tps ...string)
	OnHandleExecuting(h context.Handler, tps ...string)
	OnHandleExecuted(h context.Handler, tps ...string)
}

//Registry 服务注册管理
var Registry = &regist{
	registRouters: make(map[string]*httpServices),
}

//regist  本地服务
type regist struct {
	registRouters map[string]*httpServices
	lock          sync.RWMutex
}

//OnServerStarting 处理服务器启动
func (s *regist) OnServerStarting(h func(server.IServerConf) error, tps ...string) {
	if len(tps) == 0 {
		tps = application.DefApp.ServerTypes
	}
	for _, typ := range tps {
		if err := s.get(typ).Starting(h); err != nil {
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
		if err := s.get(typ).Closing(h); err != nil {
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
		if err := s.get(typ).GlobalHandling(h); err != nil {
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
		if err := s.get(typ).GlobalHandled(h); err != nil {
			panic(fmt.Errorf("%s OnHandleExecuted %v", typ, err))
		}
	}
}

//Micro 注册为微服务包括api,web,rpc
func (s *regist) Micro(name string, h interface{}) {
	s.register("api", name, h)
	s.register("web", name, h)
	s.register("rpc", name, h)
}

//Flow 注册为流程服务，包括mqc,cron
func (s *regist) Flow(name string, h interface{}) {
	s.register("mqc", name, h)
	s.register("cron", name, h)
}

//API 注册为API服务
func (s *regist) API(name string, h interface{}) {
	s.register("api", name, h)
}

//Web 注册为web服务
func (s *regist) Web(name string, h interface{}) {
	s.register("web", name, h)
}

//WS 注册为websocket服务
func (s *regist) WS(name string, h interface{}) {
	s.register("ws", name, h)
}

//MQC 注册为消息队列服务
func (s *regist) MQC(name string, h interface{}) {
	s.register("mqc", name, h)
}

//CRON 注册为定时任务服务
func (s *regist) CRON(name string, h interface{}) {
	s.register("cron", name, h)
}

//CRONBy 根据cron表达式，服务名称，服务处理函数注册cron服务
func (s *regist) CRONBy(cron string, name string, h interface{}) {
	CRON.Add(cron, name)
	s.register("cron", name, h)
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

func (s *regist) register(tp string, name string, h interface{}) (err error) {
	//检查参数
	if tp == "" || h == nil {
		return fmt.Errorf("注册对象不能为空")
	}
	//输入参数为函数
	current := s.get(tp)
	if vv, ok := h.(func(context.IContext) interface{}); ok {
		if err := current.AddHanler(name, "", context.Handler(vv)); err != nil {
			return err
		}
	}

	//检查输入的注册服务必须为struct
	typ := reflect.TypeOf(h)
	val := reflect.ValueOf(h)
	if val.Kind() != reflect.Ptr {
		return fmt.Errorf("只能接收引用类型; 实际是 %s", val.Kind())
	}

	//reflect所有函数，检查函数签名
	for i := 0; i < typ.NumMethod(); i++ {

		//检查函数参数是否符合接口要求
		mName := typ.Method(i).Name
		method := val.MethodByName(mName)

		//处理服务关闭函数
		if mName == defClose {
			err = current.AddCloser(name, method.Interface())
		}
		if err != nil {
			panic(err)
		}

		//处理handling,handle,handled,fallback
		nfx, ok := method.Interface().(func(context.IContext) interface{})
		if !ok {
			continue
		}

		//检查函数名是否符合特定要求
		var nf context.Handler = nfx
		switch {
		case strings.HasSuffix(mName, defHandling):
			endName := strings.ToLower(mName[0 : len(mName)-len(defHandling)])
			err = current.AddHandling(name, endName, nf)
		case strings.HasSuffix(mName, defHandler):
			endName := strings.ToLower(mName[0 : len(mName)-len(defHandler)])
			err = current.AddHanler(name, endName, nf)
		case strings.HasSuffix(mName, defHandled):
			endName := strings.ToLower(mName[0 : len(mName)-len(defHandled)])
			err = current.AddHandled(name, endName, nf)
		case strings.HasSuffix(mName, defFallback):
			endName := strings.ToLower(mName[0 : len(mName)-len(defFallback)])
			err = current.AddFallback(name, endName, nf)
		}
		if err != nil {
			panic(err)
		}
	}

	return nil

}

func (s *regist) get(tp string) *httpServices {
	if _, ok := s.registRouters[tp]; !ok {
		s.registRouters[tp] = newServices()
	}
	return s.registRouters[tp]
}

//GetServices 获取指定服务器已注册的服务
func (s *regist) GetRouters(serverType string) (*router.Routers, error) {
	return s.get(serverType).GetRouters()
}

//DoStart 执行服务启动函数
func (s *regist) DoStart(c server.IServerConf) error {
	handlers := s.get(c.GetMainConf().GetServerType()).GetStartingHandles()
	for _, h := range handlers {
		if err := h(c); err != nil {
			return err
		}
	}
	return nil
}

//DoClose 执行服务关闭函数
func (s *regist) DoClose(c server.IServerConf) error {
	handlers := s.get(c.GetMainConf().GetServerType()).GetClosingHandles()
	for _, h := range handlers {
		if err := h(c); err != nil {
			return err
		}
	}
	return nil
}

//GetClosers 获取资源释放函数
func (s *regist) Close() error {
	var sb strings.Builder
	for _, r := range s.registRouters {
		for _, cs := range r.GetClosers() {
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
