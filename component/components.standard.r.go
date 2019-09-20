package component

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/registry"
)

//AddCustomerService 添加自定义分组服务
func (r *StandardComponent) AddCustomerService(service string, h interface{}, groupName string, tags ...string) {
	r.addService(groupName, service, h)
	r.ServiceTags[service] = tags
}

//addService 添加服务处理程序
func (r *StandardComponent) addService(group string, service string, h interface{}) {
	r.addToCache(group, service, h)
	r.register(group, service, h)
	return
}
func (r *StandardComponent) addToCache(group string, service string, handler interface{}) {
	if _, ok := r.HandlerCache[group]; !ok {
		r.HandlerCache[group] = make(map[string]interface{})
	}
	if _, ok := r.HandlerCache[group][service]; !ok {
		r.HandlerCache[group][service] = handler
	}
}

func (r *StandardComponent) register(group string, name string, h interface{}) {
	for _, v := range r.GroupServices[group] {
		if v == name {
			panic(fmt.Sprintf("多次注册服务:%s:%v", name, r.GroupServices[group]))
		}
	}

	//注册get,post,put,delete,handle服务
	found := false
	hasHandle := false
	switch h.(type) {
	case Handler:
		r.registerAddService(name,name, group, h)
		found = true
		hasHandle = true
	}
	switch handler := h.(type) {
	case GetHandler:
		var f ServiceFunc = handler.GetHandle
		r.registerAddService(name,registry.Join(name, "$get"), group, f)
		r.registerAddService(name,registry.Join(name, "get", "$get"), group, f)
		if !hasHandle {
			r.registerAddService(name,name, group, f)
		}
		found = true
	}
	switch handler := h.(type) {
	case HeadHandler:
		var f ServiceFunc = handler.HeadHandle
		r.registerAddService(name,registry.Join(name, "$head"), group, f)
		found = true
	}
	switch handler := h.(type) {
	case PostHandler:
		var f ServiceFunc = handler.PostHandle
		r.registerAddService(name,registry.Join(name, "$post"), group, f)
		found = true
	}
	switch handler := h.(type) {
	case PutHandler:
		var f ServiceFunc = handler.PutHandle
		r.registerAddService(name,registry.Join(name, "$put"), group, f)
		found = true
	}
	switch handler := h.(type) {
	case DeleteHandler:
		var f ServiceFunc = handler.DeleteHandle
		r.registerAddService(name,registry.Join(name, "$delete"), group, f)
		found = true
	}

	obj := reflect.ValueOf(h)
	var t = reflect.TypeOf(h)
	for {
		if t.Kind() == reflect.Ptr {
			for i := 0; i < t.NumMethod(); i++ {
				mName := t.Method(i).Name
				if !strings.HasSuffix(mName, "Handle") || strings.EqualFold(mName, "Handle") {
					continue
				}
				method := obj.MethodByName(mName)
				nf, ok := method.Interface().(func(*context.Context) interface{})
				if !ok {
					panic("不是有效的服务类型")
				}
				var f ServiceFunc = nf
				endName := strings.ToLower(mName[0 : len(mName)-6])
				if endName == "get" || endName == "post" || endName == "put" || endName == "delete" {
					endName = "$" + endName
				}
				r.registerAddService(name,registry.Join(name, endName), group, f)
				found = true
			}
		}
		break
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
		name := registry.Join(name, "$get")
		var f FallbackServiceFunc = handler.GetFallback
		if _, ok := r.FallbackHandlers[name]; !ok {
			r.FallbackHandlers[name] = f
		}
	}

	//post降级服务
	switch handler := h.(type) {
	case PostFallbackHandler:
		name := registry.Join(name, "$post")
		var f FallbackServiceFunc = handler.PostFallback
		if _, ok := r.FallbackHandlers[name]; !ok {
			r.FallbackHandlers[name] = f
		}
	}

	//put降级服务
	switch handler := h.(type) {
	case PutFallbackHandler:
		name := registry.Join(name, "$put")
		var f FallbackServiceFunc = handler.PutFallback
		if _, ok := r.FallbackHandlers[name]; !ok {
			r.FallbackHandlers[name] = f
		}
	}

	//delete降级服务
	switch handler := h.(type) {
	case DeleteFallbackHandler:
		name := registry.Join(name, "$delete")
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

func (r *StandardComponent) registerAddService(rname string,name string, group string, handler interface{}) {
	
	if _,ok:=r.ServiceTags[name];!ok&&rname!=name{
		r.ServiceTags[name]=r.ServiceTags[rname]
	}
	
	_, hok := r.Handlers[name]
	if !hok {
		r.Handlers[name] = handler
	}
	if strings.HasPrefix(name, "__") {
		return
	}
	if _, ok := r.GroupServices[group]; !ok {
		r.GroupServices[group] = make([]string, 0, 2)
	}
	// if !hok {
	r.GroupServices[group] = append(r.GroupServices[group], name)
	// }

	r.Services = append(r.Services, name)

	if _, ok := r.ServiceGroup[name]; !ok {
		r.ServiceGroup[name] = make([]string, 0, 2)
	}
	// if _, ok := r.ServiceGroup[name]; !ok {
	r.ServiceGroup[name] = append(r.ServiceGroup[name], group)
	// }

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
		if tp.In(0).Name() != "IContainer" {
			return h, nil
		}
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
			// if h, ok := r.Handlers[name]; ok {
			// 	r.register(group, name, h)

			// 	continue
			// }
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
