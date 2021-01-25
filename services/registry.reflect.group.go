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
const (
	defaultReplaceMent = "handle"
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
func (g *UnitGroup) AddHandling(name string, h context.IHandler) {
	if name == "" {
		g.Handling = h
		return
	}
	g.storeService(name, h, handling)
}

//AddHandled 添加后处理函数
func (g *UnitGroup) AddHandled(name string, h context.IHandler) {
	if name == "" {
		g.Handled = h
		return
	}
	g.storeService(name, h, handled)
}

//AddHandle 添加处理函数
func (g *UnitGroup) AddHandle(name string, h context.IHandler) {
	g.storeService(name, h, handle)
}

//AddFallback 添加降级函数
func (g *UnitGroup) AddFallback(name string, h context.IHandler) {
	g.storeService(name, h, fallback)
}

func (g *UnitGroup) storeService(name string, handler context.IHandler, htype handlerType) {
	path, service, actions := g.getPaths(g.Path, name)
	unit, ok := g.Services[service]
	if !ok {
		unit = &Unit{Group: g, Path: path, Service: service, rawUnit: &rawUnit{RawPath: g.Path, RawMTag: name}}
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

func (g *UnitGroup) getPaths(path, name string) (rpath string, service string, action []string) {

	//rpc
	if strings.HasPrefix(name, "rpc://") {
		return path, name, []string{}
	}

	//替换注册路径中最后一个*
	lastIndex := strings.LastIndex(path, "*")
	if lastIndex > -1 {
		replacement := defaultReplaceMent
		if name != "" {
			replacement = name
		}
		path = path[:lastIndex] + replacement + path[lastIndex+1:]
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
	if lastIndex > -1 {
		return path, path, router.DefMethods
	}
	return registry.Join(path, name), registry.Join(path, name), router.DefMethods
}
