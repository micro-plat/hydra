package services

import (
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/hydra/servers"
)

const defHandling = "Handling"
const defHandler = "Handle"
const defHandled = "Handled"
const defFallback = "Fallback"

//IServiceRegistry 服务注册接口
type IServiceRegistry interface {
	Micro(name string, h interface{})
	Flow(name string, h interface{})
	API(name string, h interface{})
	Web(name string, h interface{})
	WS(name string, h interface{})
	MQC(name string, h interface{})
	CRON(name string, h interface{})
}

//Registry 服务注册管理
var Registry = &service{
	handlings: make(map[string]map[string]context.IHandler),
	handlers:  make(map[string]map[string]context.IHandler),
	fallbacks: make(map[string]map[string]context.IHandler),
	handleds:  make(map[string]map[string]context.IHandler),
}

//service  本地服务
type service struct {
	handlings map[string]map[string]context.IHandler
	handlers  map[string]map[string]context.IHandler
	fallbacks map[string]map[string]context.IHandler
	handleds  map[string]map[string]context.IHandler
	lock      sync.RWMutex
}

//OnHandleExecuting 处理handling业务
func (s *service) OnHandleExecuting(h context.IHandler, tps ...string) {
	if len(tps) == 0 {
		tps = servers.GetServerTypes()
	}
	for _, typ := range tps {
		s.check(typ)
		if _, ok := s.handlings[typ]["*"]; ok {
			panic(fmt.Sprintf("[%s]服务的Handling函数不能重复注册", typ))
		}
		s.handlings[typ]["*"] = h
	}
}

//Handled 处理Handled业务
func (s *service) OnHandleExecuted(h context.IHandler, tps ...string) {
	if len(tps) == 0 {
		tps = servers.GetServerTypes()
	}
	for _, typ := range tps {
		s.check(typ)
		if _, ok := s.handleds[typ]["*"]; ok {
			panic(fmt.Sprintf("[%s]服务的Handling函数不能重复注册", typ))
		}
		s.handleds[typ]["*"] = h
	}
}

//Micro 注册为微服务包括api,web,rpc
func (s *service) Micro(name string, h interface{}) {
	s.register("api", name, h)
	s.register("web", name, h)
	s.register("rpc", name, h)
}

//Flow 注册为流程服务，包括mqc,cron
func (s *service) Flow(name string, h interface{}) {
	s.register("mqc", name, h)
	s.register("cron", name, h)
}

//API 注册为API服务
func (s *service) API(name string, h interface{}) {
	s.register("api", name, h)
}

//Web 注册为web服务
func (s *service) Web(name string, h interface{}) {
	s.register("web", name, h)
}

//WS 注册为websocket服务
func (s *service) WS(name string, h interface{}) {
	s.register("ws", name, h)
}

//MQC 注册为消息队列服务
func (s *service) MQC(name string, h interface{}) {
	s.register("mqc", name, h)
}

//CRON 注册为定时任务服务
func (s *service) CRON(name string, h interface{}) {
	s.register("cron", name, h)
}

//CRONBy 根据cron表达式，服务名称，服务处理函数注册cron服务
func (s *service) CRONBy(cron string, name string, h interface{}) {
	CRON.Add(cron, name)
	s.register("cron", name, h)
}

//GetHandler 获取服务对应的处理函数
func (s *service) GetHandler(serverType string, service string, method string) (context.IHandler, bool) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	list := []string{service, fmt.Sprintf("%s@%s", service, method)}
	for _, srvs := range list {
		if h, ok := s.handlers[serverType][srvs]; ok {
			return h, true
		}
	}
	return nil, false
}

//GetHandling 获取预处理函数
func (s *service) GetHandlings(serverType string, service string, method string) []context.IHandler {
	handlings := make([]context.IHandler, 0, 1)
	if c, ok := s.handlings[serverType]["*"]; ok {
		handlings = append(handlings, c)
	}
	list := []string{service, fmt.Sprintf("%s@%s", service, method)}
	for _, srvs := range list {
		if h, ok := s.handlings[serverType][srvs]; ok {
			handlings = append(handlings, h)
		}
	}
	return handlings
}

//GetHandling 获取后处理函数
func (s *service) GetHandleds(serverType string, service string, method string) []context.IHandler {
	handleds := make([]context.IHandler, 0, 1)
	if c, ok := s.handleds[serverType]["*"]; ok {
		handleds = append(handleds, c)
	}
	list := []string{service, fmt.Sprintf("%s@%s", service, method)}
	for _, srvs := range list {
		if h, ok := s.handleds[serverType][srvs]; ok {
			handleds = append(handleds, h)
		}
	}
	return handleds
}

//GetFallback 获取服务对应的降级函数
func (s *service) GetFallback(serverType string, service string, method string) (context.IHandler, bool) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	list := []string{service, fmt.Sprintf("%s@%s", service, method)}
	for _, srvs := range list {
		if h, ok := s.fallbacks[serverType][srvs]; ok {
			return h, true
		}
	}
	return nil, false
}

func reflectHandler(name string, h interface{}) (handlings map[string]context.IHandler, handlers map[string]context.IHandler, handleds map[string]context.IHandler, fallbacks map[string]context.IHandler) {
	handlings = make(map[string]context.IHandler)
	handlers = make(map[string]context.IHandler)
	handleds = make(map[string]context.IHandler)
	fallbacks = make(map[string]context.IHandler)

	if vv, ok := h.(func(context.IContext) interface{}); ok {
		handlers[name] = context.Handler(vv)
		return
	}

	switch tp := h.(type) {
	case context.Handler:
		handlers[name] = tp
		return
	default:
		obj := reflect.ValueOf(h)
		typ := reflect.TypeOf(h)

		if typ.Kind() == reflect.Ptr || typ.Kind() == reflect.Struct {
			for i := 0; i < typ.NumMethod(); i++ {

				//检查函数名称是否是以Handle,Fallback结尾
				mName := typ.Method(i).Name
				if !strings.HasSuffix(mName, defHandling) &&
					!strings.HasSuffix(mName, defHandler) &&
					!strings.HasSuffix(mName, defHandled) &&
					!strings.HasSuffix(mName, defFallback) {
					continue
				}

				//处理函数名称
				var endName, handleType string
				if strings.HasSuffix(mName, defHandler) {
					handleType = defHandler
					endName = strings.ToLower(mName[0 : len(mName)-len(defHandler)])
				} else if strings.HasSuffix(mName, defFallback) {
					handleType = defFallback
					endName = strings.ToLower(mName[0 : len(mName)-len(defFallback)])
				} else if strings.HasSuffix(mName, defHandling) {
					handleType = defHandling
					endName = strings.ToLower(mName[0 : len(mName)-len(defHandling)])
				} else if strings.HasSuffix(mName, defHandled) {
					handleType = defHandled
					endName = strings.ToLower(mName[0 : len(mName)-len(defHandled)])
				}

				//检查函数是否符合指定的签名格式context.IHandler
				method := obj.MethodByName(mName)
				nfx, ok := method.Interface().(func(context.IContext) interface{})
				if !ok {
					panic(fmt.Sprintf("%s不是有效的服务类型", mName))
				}
				// fmt.Println("ok")

				var nf context.Handler = nfx
				//保存到缓存列表
				switch handleType {
				case defHandling:
					handlings[registry.Join(name, endName)] = nf
				case defHandler:
					handlers[registry.Join(name, endName)] = nf
				case defHandled:
					handleds[registry.Join(name, endName)] = nf
				case defFallback:
					fallbacks[registry.Join(name, endName)] = nf
				}
			}

		}

	}
	return
}

//register 注册服务
func (s *service) register(tp string, name string, h interface{}) {
	s.check(tp)
	handlings, handlers, handleds, fallbacks := reflectHandler(name, h)
	s.lock.Lock()
	defer s.lock.Unlock()
	for k, v := range handlings {
		if _, ok := s.handlings[tp][k]; ok {
			panic(fmt.Sprintf("服务[%s]不能重复注册", k))
		}
		s.handlings[tp][k] = v
	}
	for k, v := range handlers {
		if _, ok := s.handlers[tp][k]; ok {
			panic(fmt.Sprintf("服务[%s]不能重复注册", k))
		}
		s.handlers[tp][k] = v
	}
	for k, v := range handleds {
		if _, ok := s.handleds[tp][k]; ok {
			panic(fmt.Sprintf("服务[%s]不能重复注册", k))
		}
		s.handleds[tp][k] = v
	}

	for k, v := range fallbacks {
		if _, ok := s.fallbacks[tp][k]; ok {
			panic(fmt.Sprintf("服务[%s]不能重复注册", k))
		}
		s.fallbacks[tp][k] = v
	}
}
func (s *service) check(tp string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if _, ok := s.handlings[tp]; !ok {
		s.handlings[tp] = make(map[string]context.IHandler)
	}
	if _, ok := s.handlers[tp]; !ok {
		s.handlers[tp] = make(map[string]context.IHandler)
	}
	if _, ok := s.handleds[tp]; !ok {
		s.handleds[tp] = make(map[string]context.IHandler)
	}
	if _, ok := s.fallbacks[tp]; !ok {
		s.fallbacks[tp] = make(map[string]context.IHandler)
	}
}
func (s *service) GetServices(serverType string) []string {
	list := make([]string, 0, len(s.handlers[serverType]))
	for k := range s.handlers[serverType] {
		list = append(list, k)
	}
	return list
}
