package services

import (
	"strings"

	"github.com/micro-plat/hydra/conf/server/router"
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/registry"
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
	_, service, _ := g.getPaths(g.Path, name)
	if _, ok := g.Services[service]; ok {
		g.Services[service].Handling = h
		return
	}
	g.Services[service] = &Unit{Group: g, Service: service, Handling: h}
}

//AddHandled 添加后处理函数
func (g *UnitGroup) AddHandled(name string, h context.IHandler) {
	if name == "" {
		g.Handled = h
		return
	}
	_, service, _ := g.getPaths(g.Path, name)
	if _, ok := g.Services[service]; ok {
		g.Services[service].Handled = h
		return
	}
	g.Services[service] = &Unit{Group: g, Service: service, Handled: h}
}

//AddHandle 添加处理函数
func (g *UnitGroup) AddHandle(name string, h context.IHandler) {

	path, service, actions := g.getPaths(g.Path, name)
	if _, ok := g.Services[service]; ok {
		g.Services[service].Handle = h
		return
	}
	g.Services[service] = &Unit{Group: g, Path: path, Service: service, Actions: actions, Handle: h}
}

//AddFallback 添加降级函数
func (g *UnitGroup) AddFallback(name string, h context.IHandler) {
	_, service, _ := g.getPaths(g.Path, name)
	if _, ok := g.Services[service]; ok {
		g.Services[service].Fallback = h
		return
	}
	g.Services[service] = &Unit{Group: g, Service: service, Fallback: h}
}
func (g *UnitGroup) getPaths(path string, name string) (rpath string, service string, action []string) {
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
