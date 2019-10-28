package component

import (
	"fmt"
	"reflect"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/rpc"
)

var _ IServiceRegistry = &ServiceRegistry{}
var _ IComponentHandler = &ServiceRegistry{}

type IComponentRegistry interface {
	GetServices() map[string]map[string]interface{}
	GetHandlings() []ServiceFunc
	GetHandleds() []ServiceFunc
	GetInitializings() []ComponentFunc
	GetClosings() []ComponentFunc
	GetRPCTLS() map[string][]string
	IServiceRegistry
}

type IComponentHandler interface {
	GetServices() map[string]map[string]interface{}
	GetHandlings() []ServiceFunc
	GetHandleds() []ServiceFunc
	GetInitializings() []ComponentFunc
	GetClosings() []ComponentFunc
	GetTags(name string) []string
	GetDynamicQueue() chan *conf.Queue
	GetDynamicCron() chan *conf.Task
	GetRPCTLS() map[string][]string
	//GetBalancer 获取负载均衡模式
	GetBalancer() map[string]*rpc.BalancerMode
}

//IServiceRegistry 服务注册接口
type IServiceRegistry interface {
	//Customer 添加自定义服务
	Customer(group string, name string, h interface{}, tags ...string)

	//Micro 添加微服务（api,rpc）
	Micro(name string, h interface{}, tags ...string)

	//API 添加微服务
	API(name string, h interface{}, tags ...string)

	//RPC 添加微服务
	RPC(name string, h interface{}, tags ...string)

	//Flow 添加自动流程(mqc,cron)
	Flow(name string, h interface{}, tags ...string)

	//MQC 添加自动流程
	MQC(name string, h interface{}, tags ...string)

	//CRON 添加自动流程
	CRON(name string, h interface{}, tags ...string)

	//WS 添加websocket(mqc,cron)
	WS(name string, h interface{}, tags ...string)

	//Page 添加web页面服务(web)
	Web(name string, h interface{}, tags ...string)

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

	//GetDynamicQueue 获取动态队列注册消息
	GetDynamicQueue() chan *conf.Queue

	//GetDynamicCron 获取动态任务
	GetDynamicCron() chan *conf.Task

	//Closing 关闭组件
	Closing(c func(IContainer) error)
	//Handling 每个请求的预处理函数
	Handling(h func(c *context.Context) (rs interface{}))

	//Handled 请求后处理函数
	Handled(h func(c *context.Context) (rs interface{}))

	GetTags(name string) []string

	//AddRPCTLS 添加RPC安全认证证书
	AddRPCTLS(platName string, cert string, key string) error

	//SetBalancer 设置平台对应的负载均衡器 platName:平台名称 mode:rpc.RoundRobin 或rpc.LocalFirst
	SetBalancer(platName string, mode int, p ...string) error

	//GetBalancer 获取负载均衡模式
	GetBalancer() map[string]*rpc.BalancerMode
}

//ServiceRegistry 服务注册组件
type ServiceRegistry struct {
	services          map[string]map[string]interface{}
	handlingFuncs     []ServiceFunc
	handledFuncs      []ServiceFunc
	initializingFuncs []ComponentFunc
	closingFuncs      []ComponentFunc
	exts              map[string]interface{}
	tags              map[string][]string
	tls               map[string][]string
	rpcBalancers      map[string]*rpc.BalancerMode
	dynamicQueues     chan *conf.Queue
	dynamicCrons      chan *conf.Task
}

//NewServiceRegistry 创建ServiceRegistry
func NewServiceRegistry() *ServiceRegistry {
	return &ServiceRegistry{
		handlingFuncs:     make([]ServiceFunc, 0, 1),
		handledFuncs:      make([]ServiceFunc, 0, 1),
		initializingFuncs: make([]ComponentFunc, 0, 1),
		closingFuncs:      make([]ComponentFunc, 0, 1),
		services:          make(map[string]map[string]interface{}),
		exts:              make(map[string]interface{}),
		tags:              make(map[string][]string),
		rpcBalancers:      make(map[string]*rpc.BalancerMode),
		dynamicQueues:     make(chan *conf.Queue, 100),
		dynamicCrons:      make(chan *conf.Task, 100),
	}
}

func (s *ServiceRegistry) isConstructor(h interface{}) bool {
	fv := reflect.ValueOf(h)
	tp := reflect.TypeOf(h)
	if fv.Kind() != reflect.Func || tp.NumIn() > 1 || tp.NumOut() > 2 || tp.NumOut() == 0 {
		return false
	}
	if tp.NumIn() == 1 && tp.In(0).Name() == "IContainer" {
		return true
	}
	if tp.NumIn() == 0 {
		return true
	}
	return false
}
func (s *ServiceRegistry) isHandler(h interface{}) bool {
	fv := reflect.ValueOf(h)
	tp := reflect.TypeOf(h)
	return fv.Kind() == reflect.Func && tp.NumIn() == 1 && tp.NumOut() == 1
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
func (s *ServiceRegistry) Customer(group string, name string, h interface{}, tags ...string) {
	if _, ok := s.services[group][name]; ok {
		return
	}
	if s.isConstructor(h) {
		s.add(group, name, h)
		s.tags[name] = tags
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
		nf, ok := h.(func(*context.Context) interface{})
		if !ok {
			panic("不是有效的服务类型")
		}
		s.add(group, name, (ServiceFunc)(nf))
	}
	s.tags[name] = tags

}

//API http api服务
func (s *ServiceRegistry) API(name string, h interface{}, tags ...string) {
	s.Customer(APIService, name, h, tags...)
}

//RPC rpc服务
func (s *ServiceRegistry) RPC(name string, h interface{}, tags ...string) {
	s.Customer(RPCService, name, h, tags...)
}

//Micro 微服务
func (s *ServiceRegistry) Micro(name string, h interface{}, tags ...string) {
	s.API(name, h, tags...)
	s.RPC(name, h, tags...)
}

//MQC MQC流程服务
func (s *ServiceRegistry) MQC(name string, h interface{}, tags ...string) {
	s.Customer(MQCService, name, h, tags...)
}

//CRON Cron服务
func (s *ServiceRegistry) CRON(name string, h interface{}, tags ...string) {
	s.Customer(CRONService, name, h, tags...)
}

//Flow rpc服务
func (s *ServiceRegistry) Flow(name string, h interface{}, tags ...string) {
	s.Customer(CRONService, name, h, tags...)
	s.Customer(MQCService, name, h, tags...)
}

//WS websocket服务
func (s *ServiceRegistry) WS(name string, h interface{}, tags ...string) {
	s.Customer(WSService, name, h, tags...)
}

//Web 页面服务
func (s *ServiceRegistry) Web(name string, h interface{}, tags ...string) {
	s.Customer(PageService, name, h, tags...)
}

//Fallback 降级服务
func (s *ServiceRegistry) Fallback(name string, h interface{}) {
	if s.isConstructor(h) {
		s.add(APIService, name, h)
		return
	}
	if !s.isHandler(h) {
		panic("不是有效的服务类型")
	}
	if f, ok := h.(FallbackServiceFunc); ok {
		s.add(APIService, name, f)
	}
	panic("不是有效的FallbackServiceFunc")
}

//Get get请求
func (s *ServiceRegistry) Get(name string, h interface{}) {
	if s.isConstructor(h) {
		s.add(APIService, name, h)
		return
	}
	if !s.isHandler(h) {
		panic("不是有效的服务类型")
	}
	if f, ok := h.(GetHandler); ok {
		s.add(APIService, name, f)
	}
	panic("不是有效的GetHandler")
}

//Post post请求
func (s *ServiceRegistry) Post(name string, h interface{}) {
	if s.isConstructor(h) {
		s.add(APIService, name, h)
		return
	}
	if !s.isHandler(h) {
		panic("不是有效的服务类型")
	}
	if f, ok := h.(PostHandler); ok {
		s.add(APIService, name, f)
	}
	panic("不是有效的PostHandler")
}

//Delete delete请求
func (s *ServiceRegistry) Delete(name string, h interface{}) {
	if s.isConstructor(h) {
		s.add(APIService, name, h)
		return
	}
	if !s.isHandler(h) {
		panic("不是有效的服务类型")
	}
	if f, ok := h.(DeleteHandler); ok {
		s.add(APIService, name, f)
	}
	panic("不是有效的DeleteHandler")
}

//Put put请求
func (s *ServiceRegistry) Put(name string, h interface{}) {
	if s.isConstructor(h) {
		s.add(APIService, name, h)
		return
	}
	if !s.isHandler(h) {
		panic("不是有效的服务类型")
	}
	if f, ok := h.(PutHandler); ok {
		s.add(APIService, name, f)
	}
	panic("不是有效的PutHandler")
}

//GetFallback get降级请求
func (s *ServiceRegistry) GetFallback(name string, h interface{}) {
	if s.isConstructor(h) {
		s.add(APIService, name, h)
		return
	}
	if !s.isHandler(h) {
		panic("不是有效的服务类型")
	}
	if f, ok := h.(GetFallbackHandler); ok {
		s.add(APIService, name, f)
	}
	panic("不是有效的GetFallback")
}

//PostFallback post降级请求
func (s *ServiceRegistry) PostFallback(name string, h interface{}) {
	if s.isConstructor(h) {
		s.add(APIService, name, h)
		return
	}
	if !s.isHandler(h) {
		panic("不是有效的服务类型")
	}
	if f, ok := h.(PostFallbackHandler); ok {
		s.add(APIService, name, f)
	}
	panic("不是有效的PostFallbackHandler")
}

//DeleteFallback delete降级请求
func (s *ServiceRegistry) DeleteFallback(name string, h interface{}) {
	if s.isConstructor(h) {
		s.add(APIService, name, h)
		return
	}
	if !s.isHandler(h) {
		panic("不是有效的服务类型")
	}
	if f, ok := h.(DeleteFallbackHandler); ok {
		s.add(APIService, name, f)
	}
	panic("不是有效的DeleteFallbackHandler")
}

//PutFallback put降级请求
func (s *ServiceRegistry) PutFallback(name string, h interface{}) {
	if s.isConstructor(h) {
		s.add(APIService, name, h)
		return
	}
	if !s.isHandler(h) {
		panic("不是有效的服务类型")
	}
	if f, ok := h.(PutFallbackHandler); ok {
		s.add(APIService, name, f)
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

//GetDynamicQueue 获取动态队列
func (s *ServiceRegistry) GetDynamicQueue() chan *conf.Queue {
	return s.dynamicQueues
}

//GetDynamicCron 获取动态cron
func (s *ServiceRegistry) GetDynamicCron() chan *conf.Task {
	return s.dynamicCrons
}

//Ext 注册扩展
func (s *ServiceRegistry) Ext(name string, i interface{}) {
	s.exts[name] = i
}

//GetExt 获取扩展
func (s *ServiceRegistry) GetExt(name string) (bool, interface{}) {
	f, b := s.exts[name]
	return b, f
}

func (s *ServiceRegistry) GetServices() map[string]map[string]interface{} {
	return s.services
}
func (s *ServiceRegistry) GetTags(name string) []string {
	return s.tags[name]
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

//AddRPCTLS 添加RPC认证证书
func (s *ServiceRegistry) AddRPCTLS(platName string, cert string, key string) error {
	if cert == "" || key == "" {
		return fmt.Errorf("rpc证书文件cert:%s,key:%s不能为空", cert, key)
	}
	s.tls[platName] = []string{cert, key}
	return nil
}
func (s *ServiceRegistry) GetRPCTLS() map[string][]string {
	return s.tls
}

//SetBalancer 设置平台对应的负载均衡器 platName:平台名称 mode:rpc.RoundRobin 或rpc.LocalFirst
func (s *ServiceRegistry) SetBalancer(platName string, mode int, p ...string) error {
	if len(p) > 0 {
		s.rpcBalancers[platName] = &rpc.BalancerMode{mode, p[0]}
		return nil
	}
	s.rpcBalancers[platName] = &rpc.BalancerMode{Mode: mode}
	return nil
}
func (s *ServiceRegistry) GetBalancer() map[string]*rpc.BalancerMode {
	return s.rpcBalancers
}
