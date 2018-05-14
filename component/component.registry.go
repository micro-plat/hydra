package component

import (
	"fmt"
	"reflect"

	"github.com/micro-plat/hydra/context"
)

var _ IServiceRegistry = &ServiceRegistry{}
var _ IComponentHandler = &ServiceRegistry{}

type IComponentRegistry interface {
	IComponentHandler
	IServiceRegistry
}

type IComponentHandler interface {
	GetServices() map[string]map[string]interface{}
	GetHandlings() []ServiceFunc
	GetHandleds() []ServiceFunc
	GetInitializings() []ComponentFunc
	GetClosings() []ComponentFunc
}

//IServiceRegistry 服务注册接口
type IServiceRegistry interface {
	//Customer 添加自定义服务
	Customer(group string, name string, h interface{})
	//Micro 添加微服务（api,rpc）
	Micro(name string, h interface{})
	//Autoflow 添加自动流程(mqc,cron)
	Autoflow(name string, h interface{})
	//Page 添加web页面服务(web)
	Page(name string, h interface{}, pages ...string)

	//Fallback 默认降级函数
	Fallback(name string, h interface{})

	//Get RESTful GET请求服务
	Get(name string, h interface{})

	//Get RESTful GET请求服务的降级服务
	GetFallback(name string, h interface{})

	//Post RESTful POST请求服务
	Post(name string, h interface{})
	//PostFallback RESTful POST请求服务的降级服务
	PostFallback(name string, h interface{})

	//Delete RESTful DELETE请求服务
	Delete(name string, h interface{})
	//DeleteFallback RESTful DELETE请求服务的降级服务
	DeleteFallback(name string, h interface{})

	//Put RESTful PUT请求服务
	Put(name string, h interface{})
	//PutFallback RESTful PUT请求服务的降级服务
	PutFallback(name string, h interface{})

	//Initializing 初始化
	Initializing(c func(IContainer) error)

	//Closing 关闭组件
	Closing(c func(IContainer) error)
	//Handling 每个请求的预处理函数
	Handling(h func(c *context.Context) (rs interface{}))

	//Handled 请求后处理函数
	Handled(h func(c *context.Context) (rs interface{}))
}

//ServiceRegistry 服务注册组件
type ServiceRegistry struct {
	services          map[string]map[string]interface{}
	handlingFuncs     []ServiceFunc
	handledFuncs      []ServiceFunc
	initializingFuncs []ComponentFunc
	closingFuncs      []ComponentFunc
	pages             map[string][]string
}

//NewServiceRegistry 创建ServiceRegistry
func NewServiceRegistry() *ServiceRegistry {
	return &ServiceRegistry{
		handlingFuncs:     make([]ServiceFunc, 0, 1),
		handledFuncs:      make([]ServiceFunc, 0, 1),
		initializingFuncs: make([]ComponentFunc, 0, 1),
		closingFuncs:      make([]ComponentFunc, 0, 1),
		services:          make(map[string]map[string]interface{}),
		pages:             make(map[string][]string),
	}
}

func (s *ServiceRegistry) isConstructor(h interface{}) bool {
	fv := reflect.ValueOf(h)
	tp := reflect.TypeOf(h)
	return fv.Kind() == reflect.Func && tp.NumIn() <= 2 && tp.NumOut() >= 1 && tp.NumOut() <= 2
}
func (s *ServiceRegistry) isHandler(h interface{}) bool {
	fv := reflect.ValueOf(h)
	tp := reflect.TypeOf(h)
	return fv.Kind() == reflect.Func && tp.NumIn() == 4 && tp.NumOut() == 1
}
func (s *ServiceRegistry) add(group string, name string, h interface{}) {
	g, ok := s.services[group]
	if !ok {
		s.services[group] = make(map[string]interface{})
		g = s.services[group]
	}
	if _, ok := g[name]; ok {
		panic(fmt.Sprintf("服务已存在:%s %s", group, name))
	}
	g[name] = h

}

//Customer 自定义服务
func (s *ServiceRegistry) Customer(group string, name string, h interface{}) {
	if _, ok := s.services[group][name]; ok {
		return
	}
	if s.isConstructor(h) {
		s.add(group, name, h)
		return
	}
	if !s.isHandler(h) {
		panic("不是有效的服务类型")
	}
	switch v := h.(type) {
	case ServiceFunc:
		s.add(group, name, v)
	case Handler:
		s.add(group, name, v)
	default:
		panic("不是有效的服务类型")
	}

}

//Micro 微服务
func (s *ServiceRegistry) Micro(name string, h interface{}) {
	s.Customer(MicroService, name, h)
}

//Autoflow 流程服务
func (s *ServiceRegistry) Autoflow(name string, h interface{}) {
	s.Customer(AutoflowService, name, h)
}

//Page 页面服务
func (s *ServiceRegistry) Page(name string, h interface{}, pages ...string) {
	s.Customer(PageService, name, h)
	s.pages[name] = pages
}

//Fallback 降级服务
func (s *ServiceRegistry) Fallback(name string, h interface{}) {
	if s.isConstructor(h) {
		s.add(MicroService, name, h)
		return
	}
	if !s.isHandler(h) {
		panic("不是有效的服务类型")
	}
	if f, ok := h.(FallbackServiceFunc); ok {
		s.add(MicroService, name, f)
	}
	panic("不是有效的FallbackServiceFunc")
}

//Get get请求
func (s *ServiceRegistry) Get(name string, h interface{}) {
	if s.isConstructor(h) {
		s.add(MicroService, name, h)
		return
	}
	if !s.isHandler(h) {
		panic("不是有效的服务类型")
	}
	if f, ok := h.(GetHandler); ok {
		s.add(MicroService, name, f)
	}
	panic("不是有效的GetHandler")
}

//Post post请求
func (s *ServiceRegistry) Post(name string, h interface{}) {
	if s.isConstructor(h) {
		s.add(MicroService, name, h)
		return
	}
	if !s.isHandler(h) {
		panic("不是有效的服务类型")
	}
	if f, ok := h.(PostHandler); ok {
		s.add(MicroService, name, f)
	}
	panic("不是有效的PostHandler")
}

//Delete delete请求
func (s *ServiceRegistry) Delete(name string, h interface{}) {
	if s.isConstructor(h) {
		s.add(MicroService, name, h)
		return
	}
	if !s.isHandler(h) {
		panic("不是有效的服务类型")
	}
	if f, ok := h.(DeleteHandler); ok {
		s.add(MicroService, name, f)
	}
	panic("不是有效的DeleteHandler")
}

//Put put请求
func (s *ServiceRegistry) Put(name string, h interface{}) {
	if s.isConstructor(h) {
		s.add(MicroService, name, h)
		return
	}
	if !s.isHandler(h) {
		panic("不是有效的服务类型")
	}
	if f, ok := h.(PutHandler); ok {
		s.add(MicroService, name, f)
	}
	panic("不是有效的PutHandler")
}

//GetFallback get降级请求
func (s *ServiceRegistry) GetFallback(name string, h interface{}) {
	if s.isConstructor(h) {
		s.add(MicroService, name, h)
		return
	}
	if !s.isHandler(h) {
		panic("不是有效的服务类型")
	}
	if f, ok := h.(GetFallbackHandler); ok {
		s.add(MicroService, name, f)
	}
	panic("不是有效的GetFallback")
}

//PostFallback post降级请求
func (s *ServiceRegistry) PostFallback(name string, h interface{}) {
	if s.isConstructor(h) {
		s.add(MicroService, name, h)
		return
	}
	if !s.isHandler(h) {
		panic("不是有效的服务类型")
	}
	if f, ok := h.(PostFallbackHandler); ok {
		s.add(MicroService, name, f)
	}
	panic("不是有效的PostFallbackHandler")
}

//DeleteFallback delete降级请求
func (s *ServiceRegistry) DeleteFallback(name string, h interface{}) {
	if s.isConstructor(h) {
		s.add(MicroService, name, h)
		return
	}
	if !s.isHandler(h) {
		panic("不是有效的服务类型")
	}
	if f, ok := h.(DeleteFallbackHandler); ok {
		s.add(MicroService, name, f)
	}
	panic("不是有效的DeleteFallbackHandler")
}

//PutFallback put降级请求
func (s *ServiceRegistry) PutFallback(name string, h interface{}) {
	if s.isConstructor(h) {
		s.add(MicroService, name, h)
		return
	}
	if !s.isHandler(h) {
		panic("不是有效的服务类型")
	}
	if f, ok := h.(PutFallbackHandler); ok {
		s.add(MicroService, name, f)
	}
	panic("不是有效的PutFallbackHandler")
}

//Handling 预处理程序
func (s *ServiceRegistry) Handling(h func(c *context.Context) (rs interface{})) {
	s.handlingFuncs = append(s.handlingFuncs, h)
}

//Handled 处理程序
func (s *ServiceRegistry) Handled(h func(c *context.Context) (rs interface{})) {
	s.handledFuncs = append(s.handledFuncs, h)
}

//Initializing 初始化
func (s *ServiceRegistry) Initializing(c func(c IContainer) error) {
	s.initializingFuncs = append(s.initializingFuncs, c)
}

//Closing 关闭组件
func (s *ServiceRegistry) Closing(c func(c IContainer) error) {
	s.closingFuncs = append(s.closingFuncs, c)
}

func (s *ServiceRegistry) GetServices() map[string]map[string]interface{} {
	return s.services
}
func (s *ServiceRegistry) GetHandlings() []ServiceFunc {
	return s.handlingFuncs
}
func (s *ServiceRegistry) GetHandleds() []ServiceFunc {
	return s.handledFuncs
}
func (s *ServiceRegistry) GetInitializings() []ComponentFunc {
	return s.initializingFuncs
}
func (s *ServiceRegistry) GetClosings() []ComponentFunc {
	return s.closingFuncs
}
