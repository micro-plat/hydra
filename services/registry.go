package services

import (
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/registry"
)

const defHandler = "Handle"
const defFallback = "Fallback"

var defRequestMethod = []string{"get", "post", "put", "delete"}

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
	handlers:  make(map[string]map[string]context.IHandler),
	fallbacks: make(map[string]map[string]context.IHandler),
}

//service  本地服务
type service struct {
	handlers  map[string]map[string]context.IHandler
	fallbacks map[string]map[string]context.IHandler
	lock      sync.RWMutex
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

//register 注册服务
func (s *service) register(tp string, name string, h interface{}) {
	s.check(tp)
	handlers, fallbacks := reflectHandler(name, h)
	for k, v := range handlers {
		if _, ok := s.handlers[tp][k]; ok {
			panic(fmt.Sprintf("服务[%s]不能重复注册", k))
		}
		s.handlers[tp][k] = v
	}
	s.lock.Lock()
	defer s.lock.Unlock()
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
	if _, ok := s.handlers[tp]; !ok {
		s.handlers[tp] = make(map[string]context.IHandler)
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

func reflectHandler(name string, h interface{}) (handlers map[string]context.IHandler, fallbacks map[string]context.IHandler) {
	switch tp := h.(type) {
	case context.IHandler:
		handlers[name] = tp
		return
	default:
		obj := reflect.ValueOf(h)
		typ := reflect.TypeOf(h)
		for {
			if typ.Kind() == reflect.Ptr {
				for i := 0; i < typ.NumMethod(); i++ {

					//检查函数名称是否是以Handle,Fallback结尾
					mName := typ.Method(i).Name
					if !strings.HasSuffix(mName, defHandler) && !strings.HasSuffix(mName, defFallback) {
						continue
					}

					//处理函数名称
					var endName, handleType string
					if strings.HasSuffix(mName, defHandler) {
						handleType = defHandler
						endName = strings.ToLower(mName[0 : len(mName)-len(defHandler)])
					}
					if strings.HasSuffix(mName, defFallback) {
						handleType = defFallback
						endName = strings.ToLower(mName[0 : len(mName)-len(defFallback)])
					}
					for _, m := range defRequestMethod {
						if m == endName {
							endName = "$" + endName
							break
						}
					}

					//检查函数是否符合指定的签名格式context.IHandler
					method := obj.MethodByName(mName)
					nf, ok := method.Interface().(context.Handler)
					if !ok {
						panic("不是有效的服务类型")
					}

					//保存到缓存列表
					if handleType == defHandler {
						handlers[registry.Join(name, endName)] = nf
						continue
					}
					fallbacks[registry.Join(name, endName)] = nf
				}
			}
			break
		}

	}
	return
}
