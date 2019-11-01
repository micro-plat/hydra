package component

import (
	"errors"
	"fmt"
	"sync"

	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/registry"
)

const (
	//MicroService 微服务
	// MicroService = "__micro_"

	//APIService http api 服务
	APIService = "__api_"

	//RPCService rpc 服务
	RPCService = "__rpc_"

	//WSService websocket流程服务
	WSService = "__websocket_"

	//FlowService 自动流程服务
	// FlowService = "__flow_"

	//PageService 页面服务
	PageService = "__page_"

	//MQCService mqc 服务
	MQCService = "__mqc_"

	//CRONService cron服务
	CRONService = "__cron_"

	//CustomerService 自定义服务
	CustomerService = "__customer_"
)

var _ IComponent = &StandardComponent{}

//ErrNotFoundFallbackService 未找到服务
var ErrNotFoundFallbackService = errors.New("未实现降级服务")

//StandardComponent 标准组件
type StandardComponent struct {
	HandlerCache     map[string]map[string]interface{} //handler缓存
	Container        IContainer
	Name             string                            //组件名称
	funcs            map[string]map[string]interface{} //每个分组对应的服务及处理程序
	Handlers         map[string]interface{}            //每个服务对应的处理程序
	FallbackHandlers map[string]interface{}            //每个服务对应的降级处理程序
	Services         []string                          //所有服务
	GroupServices    map[string][]string               //每个分组包含的服务
	ServiceGroup     map[string][]string               //每个服务对应的分组
	ServiceTags      map[string][]string               //每个服务对应的页面
	CloseHandler     []interface{}                     //用于关闭所有handler
	metaData         map[string]interface{}
	metaLock         sync.RWMutex
}

//NewStandardComponent 构建标准组件
func NewStandardComponent(componentName string, c IContainer) *StandardComponent {
	r := &StandardComponent{Name: componentName, Container: c}
	r.HandlerCache = make(map[string]map[string]interface{})
	r.funcs = make(map[string]map[string]interface{})
	r.Handlers = make(map[string]interface{})
	r.FallbackHandlers = make(map[string]interface{})
	r.GroupServices = make(map[string][]string)
	r.ServiceGroup = make(map[string][]string)
	r.Services = make([]string, 0, 2)
	r.ServiceTags = make(map[string][]string)
	r.CloseHandler = make([]interface{}, 0, 2)
	return r
}

func (m *StandardComponent) GetMeta(key string) interface{} {
	m.metaLock.RLocker().Lock()
	defer m.metaLock.RLocker().Unlock()

	data := m.metaData[key]
	return data
}
func (m *StandardComponent) SetMeta(key string, value interface{}) {
	m.metaLock.Lock()
	defer m.metaLock.Unlock()
	if m.metaData == nil {
		m.metaData = make(map[string]interface{})
	}
	m.metaData[key] = value
}

//AddRPCProxy 添加RPC代理
func (r *StandardComponent) AddRPCProxy(h interface{}) {
	r.addService(APIService, "__rpc_", h)
	r.addService(RPCService, "__rpc_", h)
	r.addService(WSService, "__rpc_", h)
	r.addService(PageService, "__rpc_", h)
	r.addService(MQCService, "__rpc_", h)
	r.addService(CRONService, "__rpc_", h)
}

//IsCustomerService 是否是指定的分组服务
func (r *StandardComponent) IsCustomerService(service string, group ...string) bool {
	groups := r.GetGroups(service)
	for _, v := range groups {
		for _, g := range group {
			if v == g {
				return true
			}
		}
	}
	return false
}

//GetComponent 获取当前组件
func (r *StandardComponent) GetComponent() IComponent {
	return r
}

//GetServices 获取组件提供的所有服务
func (r *StandardComponent) GetServices() []string {
	return r.Services
}

//GetGroupServices 根据分组获取服务
func (r *StandardComponent) GetGroupServices(group ...string) []string {
	srvs := make([]string, 0, 4)
	for _, g := range group {
		srvs = append(srvs, r.GroupServices[g]...)
	}
	return srvs
}

//GetGroups 获取服务的分组列表
func (r *StandardComponent) GetGroups(service string) []string {
	return r.ServiceGroup[service]
}

//GetTags 获取服务的tag列表
func (r *StandardComponent) GetTags(service string) []string {
	return r.ServiceTags[service]
}

//GetFallbackHandlers 获取fallback处理程序
func (r *StandardComponent) GetFallbackHandlers() map[string]interface{} {
	return r.FallbackHandlers
}

//GetCachedHandler 获取已缓存的handler
func (r *StandardComponent) GetCachedHandler(group string, service string) interface{} {
	if srvs, ok := r.HandlerCache[group]; ok {
		return srvs[service]
	}
	return nil
}

//AddFallbackHandlers 添加降级函数
func (r *StandardComponent) AddFallbackHandlers(f map[string]interface{}) {
	for k, v := range f {
		r.FallbackHandlers[k] = v
	}
}

//Handling 每次handle执行前执行
func (r *StandardComponent) Handling(c *context.Context) (rs interface{}) {
	return nil
}

//Handled 每次handle执行后执行
func (r *StandardComponent) Handled(c *context.Context) (rs interface{}) {
	return nil
}

//GetHandler 获取服务的处理函数
func (r *StandardComponent) GetHandler(engine string, service string, method string) (interface{}, bool) {
	switch engine {
	case "rpc":
		r, ok := r.Handlers["__rpc_"]
		return r, ok
	default:
		if r, ok := r.Handlers[registry.Join(service, "$"+method)]; ok {
			return r, ok
		}
		r, ok := r.Handlers[service]
		return r, ok
	}
}

//Handle 组件服务执行
func (r *StandardComponent) Handle(c *context.Context) (rs interface{}) {
	h, ok := r.GetHandler(c.Engine, c.Service, c.Request.GetMethod())
	if !ok {
		c.Response.SetStatus(404)
		return fmt.Errorf("%s:未找到服务:%s", r.Name, c.Service)
	}
	if r.IsCustomerService(c.Service, PageService) {
		c.Response.SetHTML()
	}
	switch handler := h.(type) {
	case Handler:
		rs = handler.Handle(c)
	default:
		c.Response.SetStatus(404)
		rs = fmt.Errorf("未找到服务:%s", c.Service)
	}
	return
}

//GetFallbackHandler 获取失败降级处理函数
func (r *StandardComponent) GetFallbackHandler(engine string, service string, method string) (interface{}, bool) {
	if f, ok := r.FallbackHandlers[registry.Join(service, "$", method)]; ok {
		return f, ok
	}
	f, ok := r.FallbackHandlers[service]
	return f, ok

}

//Fallback 降级处理
func (r *StandardComponent) Fallback(c *context.Context) (rs interface{}) {
	h, ok := r.GetFallbackHandler(c.Engine, c.Service, c.Request.GetMethod())
	if !ok {
		c.Response.SetStatus(404)
		return ErrNotFoundFallbackService
	}
	switch handler := h.(type) {
	case FallbackHandler:
		rs = handler.Fallback(c)
	default:
		c.Response.SetStatus(404)
		rs = fmt.Errorf("%v:%s", ErrNotFoundFallbackService, c.Service)
	}
	return
}

//Close 卸载组件
func (r *StandardComponent) Close() error {
	r.funcs = nil
	r.Handlers = nil
	r.GroupServices = nil
	r.ServiceGroup = nil
	r.Services = nil
	r.ServiceTags = nil
	for _, handler := range r.CloseHandler {
		h := handler.(CloseHandler)
		h.Close()
	}
	return nil
}

func (r *StandardComponent) GetRegistryNames(groups ...string) map[string][]string {
	srvs := r.GetGroupServices(groups...)
	svsMap := make(map[string][]string, len(srvs)/2)
	for _, srv := range srvs {
		name, methods := getMethod(srv)
		if _, ok := svsMap[name]; !ok {
			svsMap[name] = []string{}
		}
		if len(methods) > 0 {
			for _, method := range methods {
				if !requestMethods.Contains(method) {
					panic(fmt.Sprintf("不支持的请求方式:%s(%s)", srv, method))
				}
				if !xContains(svsMap[name], method) {
					svsMap[name] = append(svsMap[name], method)
				}
			}

		} else { //默认只支持GET,POST
			if !xContains(svsMap[name], "get") {
				svsMap[name] = append(svsMap[name], "get")
			}
			if !xContains(svsMap[name], "post") {
				svsMap[name] = append(svsMap[name], "post")
			}
		}

	}
	//"HEAD", "OPTIONS"
	return svsMap
}

//GetGroupName 获取分组类型[api,rpc > micro mq,cron > flow, web > page,others > customer]
func GetGroupName(serverType string) []string {
	switch serverType {
	case "api":
		return []string{APIService}
	case "rpc":
		return []string{RPCService}
	case "mqc":
		return []string{MQCService}
	case "cron":
		return []string{CRONService}
	case "web":
		return []string{PageService, APIService}
	case "ws":
		return []string{WSService}
	}
	return []string{CustomerService}
}
