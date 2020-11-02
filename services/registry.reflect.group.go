package services

import (
	"strings"

	"github.com/micro-plat/hydra/conf/server/router"
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/registry"
)

type handlerType int

const (
	handle handlerType = iota + 1
	handling
	handled
	fallback
)

type UnitGroup struct {
	Path     string
	Handling context.IHandler
	Handled  context.IHandler
	Closing  interface{}
	Services map[string]*Unit
}

func newUnitGroup(path string) *UnitGroup {
	return &UnitGroup{
		Path:     path,
		Services: make(map[string]*Unit),
	}
}

//AddHandling 添加预处理函数
func (g *UnitGroup) AddHandling(name,hName string, h context.IHandler) {
	if name == "" {
		g.Handling = h
		return
	}
	//@bugfix liujinyin 修改注册对象时候，包含Handle,Handing,Handled,Fallback 丢失Path造成的错误提醒“重复注册问题”
	g.storeService(name,hName, h, handling)
}

//AddHandled 添加后处理函数
func (g *UnitGroup) AddHandled(name,hName string, h context.IHandler) {
	if name == "" {
		g.Handled = h
		return
	}
	g.storeService(name,hName, h, handled)
}

//AddHandle 添加处理函数
func (g *UnitGroup) AddHandle(name,hName string, h context.IHandler) {
	g.storeService(name,hName, h, handle)
}

//AddFallback 添加降级函数
func (g *UnitGroup) AddFallback(name,hName string, h context.IHandler) {
	g.storeService(name,hName, h, fallback)
}

func (g *UnitGroup) storeService(name,hName string, handler context.IHandler, htype handlerType) {
	//@bugfix liujinyin 修改注册对象时候，包含Handle,Handing,Handled,Fallback 丢失Path造成的错误提醒“重复注册问题”
	path, service, actions := g.getPaths(g.Path, name,hName)
	unit, ok := g.Services[service]
	if !ok {
		unit = &Unit{Group: g, Path: path, Service: service}
		g.Services[service] = unit
	}

	switch htype {
	case handling:
		unit.Handling = handler
	case handle:
		unit.Actions = actions
		unit.Handle = handler
	case handled:
		unit.Handled = handler
	case fallback:
		unit.Fallback = handler
	default:
	}
}

func (g *UnitGroup) getPaths(path , name, hName string) (rpath string, service string, action []string) {

	//@todo path 路径里面对最后一个*进行处理
   if hName!=""{
	   lastIndex:=strings.LastIndex(path,"*")
	   if lastIndex > -1{
		path=path[:lastIndex-1]+hName+path[lastIndex+1:]
	   }
   }

	//作为func注册的服务，只支持GET，POST
	
	if name == "" {
		return path, path, []string{}
	}

	//RESTful
	for _, m := range router.Methods {
		if strings.EqualFold(m, name) {
			return path, registry.Join(path, "$"+name), []string{m}
		}
	}

	//非RESTful
	return registry.Join(path, name), registry.Join(path, name), router.DefMethods
}
