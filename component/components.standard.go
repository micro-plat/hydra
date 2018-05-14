package component

import (
	"errors"
	"fmt"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/micro-plat/hydra/context"
)

const (
	//MicroService 微服务
	MicroService = "__micro_"
	//AutoflowService 自动流程服务
	AutoflowService = "__autoflow_"

	//PageService 页面服务
	PageService = "__page_"

	//CustomerService 自定义服务
	CustomerService = "__customer_"
)

var _ IComponent = &StandardComponent{}

//ErrNotFoundService 未找到服务
var ErrNotFoundService = errors.New("未找到服务")

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
	ServicePages     map[string][]string               //每个服务对应的页面
	CloseHandler     []interface{}                     //用于关闭所有handler
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
	r.ServicePages = make(map[string][]string)
	r.CloseHandler = make([]interface{}, 0, 2)
	return r
}

//AddMicroService 添加微服务(供http,rpc方式调用)
func (r *StandardComponent) AddMicroService(service string, h interface{}) {
	r.addService(MicroService, service, h)
}

//AddAutoflowService 添加自动流程服务(供cron，mq方式调用)
func (r *StandardComponent) AddAutoflowService(service string, h interface{}) {
	r.addService(AutoflowService, service, h)
}

//AddPageService 添加页面服务
func (r *StandardComponent) AddPageService(service string, h interface{}, pages ...string) {
	r.addService(PageService, service, h)
	r.ServicePages[service] = pages
}

//AddRPCProxy 添加RPC代理
func (r *StandardComponent) AddRPCProxy(h interface{}) {
	r.addService(MicroService, "__rpc_", h)
	r.addService(AutoflowService, "__rpc_", h)
}

//AddTagPageService 添加带有标签的页面服务
func (r *StandardComponent) AddTagPageService(service string, h interface{}, pages ...string) {
	r.addService(PageService, service, h)
	r.ServicePages[service] = pages
}

//AddCustomerService 添加自定义分组服务
func (r *StandardComponent) AddCustomerService(service string, h interface{}, groupNames ...string) {
	if len(groupNames) == 0 {
		panic(fmt.Sprintf("服务:%s未指定分组名称", service))
	}
	for _, group := range groupNames {
		r.addService(group, service, h)
	}
}

//IsMicroService 是否是微服务
func (r *StandardComponent) IsMicroService(service string) bool {
	return r.IsCustomerService(service, MicroService)
}

//IsAutoflowService 是否是自动流程服务
func (r *StandardComponent) IsAutoflowService(service string) bool {
	return r.IsCustomerService(service, AutoflowService)
}

//IsPageService 是否是页面服务
func (r *StandardComponent) IsPageService(service string) bool {
	return r.IsCustomerService(service, PageService)
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

func (r *StandardComponent) addToCache(group string, service string, handler interface{}) {
	if _, ok := r.HandlerCache[group]; !ok {
		r.HandlerCache[group] = make(map[string]interface{})
	}
	if _, ok := r.HandlerCache[group][service]; !ok {
		r.HandlerCache[group][service] = handler
	}
}

//addService 添加服务处理程序
func (r *StandardComponent) addService(group string, service string, h interface{}) {
	r.addToCache(group, service, h)
	r.register(group, service, h)
	return
}
func (r *StandardComponent) registerAddService(name string, group string, handler interface{}) {
	if _, ok := r.Handlers[name]; !ok {
		r.Handlers[name] = handler
		r.Services = append(r.Services, name)
	}
	if strings.HasPrefix(name, "__") {
		return
	}
	if _, ok := r.GroupServices[group]; !ok {
		r.GroupServices[group] = make([]string, 0, 2)
	}
	r.GroupServices[group] = append(r.GroupServices[group], name)

	if _, ok := r.ServiceGroup[name]; !ok {
		r.ServiceGroup[name] = make([]string, 0, 2)
	}
	r.ServiceGroup[name] = append(r.ServiceGroup[name], group)
}
func (r *StandardComponent) register(group string, name string, h interface{}) {
	for _, v := range r.GroupServices[group] {
		if v == name {
			panic(fmt.Sprintf("多次注册服务:%s:%v", name, r.GroupServices[group]))
		}
	}

	//注册get,post,put,delete,handle服务
	found := false
	switch handler := h.(type) {
	case GetHandler:
		var f ServiceFunc = handler.GetHandle
		r.registerAddService(filepath.Join(name, "get"), group, f)
		found = true
	}
	switch handler := h.(type) {
	case PostHandler:
		var f ServiceFunc = handler.PostHandle
		r.registerAddService(filepath.Join(name, "post"), group, f)
		found = true
	}
	switch handler := h.(type) {
	case PutHandler:
		var f ServiceFunc = handler.PutHandle
		r.registerAddService(filepath.Join(name, "put"), group, f)
		found = true
	}
	switch handler := h.(type) {
	case DeleteHandler:
		var f ServiceFunc = handler.DeleteHandle
		r.registerAddService(filepath.Join(name, "delete"), group, f)
		found = true
	}
	switch h.(type) {
	case Handler:
		r.registerAddService(name, group, h)
		found = true
	}
	if !found {
		r.checkFuncType(name, h)
		if _, ok := r.funcs[group]; !ok {
			r.funcs[group] = make(map[string]interface{})
		}
		if _, ok := r.funcs[group][name]; ok {
			panic(fmt.Sprintf("多次注册服务:%s", name))
		}
		r.funcs[group][name] = h
	}

	//close handler
	switch h.(type) {
	case CloseHandler:
		r.CloseHandler = append(r.CloseHandler, h)
	}

	//处理降级服务

	//get降级服务
	switch handler := h.(type) {
	case GetFallbackHandler:
		name := filepath.Join(name, "get")
		var f FallbackServiceFunc = handler.GetFallback
		if _, ok := r.FallbackHandlers[name]; !ok {
			r.FallbackHandlers[name] = f
		}
	}

	//post降级服务
	switch handler := h.(type) {
	case PostFallbackHandler:
		name := filepath.Join(name, "post")
		var f FallbackServiceFunc = handler.PostFallback
		if _, ok := r.FallbackHandlers[name]; !ok {
			r.FallbackHandlers[name] = f
		}
	}

	//put降级服务
	switch handler := h.(type) {
	case PutFallbackHandler:
		name := filepath.Join(name, "put")
		var f FallbackServiceFunc = handler.PutFallback
		if _, ok := r.FallbackHandlers[name]; !ok {
			r.FallbackHandlers[name] = f
		}
	}

	//delete降级服务
	switch handler := h.(type) {
	case DeleteFallbackHandler:
		name := filepath.Join(name, "delete")
		var f FallbackServiceFunc = handler.DeleteFallback
		if _, ok := r.FallbackHandlers[name]; !ok {
			r.FallbackHandlers[name] = f
		}
	}

	//通用降级服务
	switch handler := h.(type) {
	case FallbackHandler:
		if _, ok := r.FallbackHandlers[name]; !ok {
			r.FallbackHandlers[name] = handler
		}
	}
}
func (r *StandardComponent) checkFuncType(name string, h interface{}) {
	fv := reflect.ValueOf(h)
	if fv.Kind() != reflect.Func {
		panic(fmt.Sprintf("服务:%s必须为Handler,MapHandler,StandardHandler,ObjectHandler,WebHandler, Handler, MapServiceFunc, StandardServiceFunc, WebServiceFunc, ServiceFunc:%v", name, h))
	}
	tp := reflect.TypeOf(h)
	if tp.NumIn() > 2 || tp.NumOut() == 0 || tp.NumOut() > 2 {
		panic(fmt.Sprintf("服务:%s只能包含最多1个输入参数(%d)，最多2个返回值(%d)", name, tp.NumIn(), tp.NumOut()))
	}
	if tp.NumIn() == 1 {
		if tp.In(0).Name() != "IContainer" {
			panic(fmt.Sprintf("服务:%s输入参数必须为component.IContainer类型(%s)", name, tp.In(0).Name()))
		}
	}
	if tp.NumOut() == 2 {
		if tp.Out(1).Name() != "error" {
			panic(fmt.Sprintf("服务:%s的2个返回值必须为error类型", name))
		}
	}
}
func (r *StandardComponent) callFuncType(name string, h interface{}) (i interface{}, err error) {
	fv := reflect.ValueOf(h)
	tp := reflect.TypeOf(h)
	var rvalue []reflect.Value
	if tp.NumIn() == 1 {
		ivalue := make([]reflect.Value, 0, 1)
		ivalue = append(ivalue, reflect.ValueOf(r.Container))
		rvalue = fv.Call(ivalue)
	} else {
		rvalue = fv.Call(nil)
	}
	if len(rvalue) == 0 || len(rvalue) > 2 {
		panic(fmt.Sprintf("%s类型错误,返回值只能有1个(handler)或2个（Handler,error）", name))
	}
	if len(rvalue) > 1 {
		if rvalue[1].Interface() != nil {
			if err, ok := rvalue[1].Interface().(error); ok {
				return nil, err
			}
		}
	}
	return rvalue[0].Interface(), nil
}

//LoadServices 加载所有服务
func (r *StandardComponent) LoadServices() error {
	for group, v := range r.funcs {
		for name, sv := range v {
			if h, ok := r.Handlers[name]; ok {
				r.register(group, name, h)
				continue
			}
			rt, err := r.callFuncType(name, sv)
			if err != nil {
				return err
			}
			r.register(group, name, rt)
		}
		delete(r.funcs, group)
	}
	return nil
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

//GetPages 获取服务的页面列表
func (r *StandardComponent) GetPages(service string) []string {
	return r.ServicePages[service]
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
		if r, ok := r.Handlers[filepath.Join(service, method)]; ok {
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
	if r.IsPageService(c.Service) {
		c.Response.SetTextHTML()
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
	if f, ok := r.FallbackHandlers[filepath.Join(service, method)]; ok {
		return f, ok
	}
	f, ok := r.FallbackHandlers[service]
	return f, ok

}

//Fallback 降级处理
func (r *StandardComponent) Fallback(c *context.Context) (rs interface{}) {
	c.Response.SetStatus(404)
	h, ok := r.GetFallbackHandler(c.Engine, c.Service, c.Request.GetMethod())
	if !ok {
		return ErrNotFoundService
	}
	switch handler := h.(type) {
	case FallbackHandler:
		rs = handler.Fallback(c)
	default:
		rs = fmt.Errorf("未找到服务:%s", c.Service)
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
	r.ServicePages = nil
	for _, handler := range r.CloseHandler {
		h := handler.(CloseHandler)
		h.Close()
	}
	return nil
}

//GetGroupName 获取分组类型[api,rpc > micro mq,cron > autoflow, web > page,others > customer]
func GetGroupName(serverType string) []string {
	switch serverType {
	case "api", "rpc":
		return []string{MicroService}
	case "mqc", "cron":
		return []string{AutoflowService}
	case "web":
		return []string{PageService, MicroService}
	}
	return []string{CustomerService}
}
