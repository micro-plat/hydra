package services

import (
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/registry"
)

type unitGroup struct {
	Path     string
	Handling context.IHandler
	Handled  context.IHandler
	Closing  interface{}
	Services map[string]*unit
}

func newUnitGroup(path string) *unitGroup {
	return &unitGroup{
		Path:     path,
		Services: make(map[string]*unit),
	}
}
func (g *unitGroup) AddHandling(name string, h context.IHandler) {
	if name == "" {
		g.Handling = h
		return
	}

	_, service, _ := g.getPaths(g.Path, name)
	if _, ok := g.Services[service]; ok {
		g.Services[service].Handling = h
		return
	}
	g.Services[service] = &unit{group: g, service: service, Handling: h}
}
func (g *unitGroup) AddHandled(name string, h context.IHandler) {
	if name == "" {
		g.Handled = h
		return
	}
	_, service, _ := g.getPaths(g.Path, name)
	if _, ok := g.Services[service]; ok {
		g.Services[service].Handled = h
		return
	}
	g.Services[service] = &unit{group: g, service: service, Handled: h}
}
func (g *unitGroup) AddHandle(name string, h context.IHandler) {

	path, service, actions := g.getPaths(g.Path, name)

	if _, ok := g.Services[service]; ok {
		g.Services[service].Handle = h
		return
	}
	g.Services[service] = &unit{group: g, path: path, service: service, actions: actions, Handle: h}
}
func (g *unitGroup) AddFallback(name string, h context.IHandler) {
	_, service, _ := g.getPaths(g.Path, name)
	if _, ok := g.Services[service]; ok {
		g.Services[service].Fallback = h
		return
	}
	g.Services[service] = &unit{group: g, service: service, Fallback: h}
}
func (g *unitGroup) getPaths(path string, name string) (rpath string, service string, action []string) {
	if name == "" {
		return path, path, []string{}
	}
	for _, m := range defRequestMethod {
		if m == name {
			return path, registry.Join(path, "$"+name), []string{m}
		}
	}
	return registry.Join(path, name), registry.Join(path, name), defRequestMethod
}
