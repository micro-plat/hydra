package services

import (
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/hydra/registry/conf/server/router"
)

var defRequestMethod = []string{"get", "post", "put", "delete"}

type services struct {
	routers   *router.Routers
	routerMap map[string][]string
	handlers  map[string]context.IHandler
}

func newServices() *services {
	return &services{
		routers:   router.NewRouters(),
		routerMap: make(map[string][]string),
		handlers:  make(map[string]context.IHandler),
	}
}
func (s *services) AddHanler(path string, name string, h context.IHandler) {
	methods, tp := getMethods(name)

	switch tp {
	case 1:
		s.routers.Append(path, path, methods...)
		s.routerMap[path] = append(s.routerMap[path], path)
		s.handlers[path] = h
	case 2:
		sv := registry.Join(path, "$"+name)
		s.routers.Append(path, sv, methods...)
		s.routerMap[path] = append(s.routerMap[path], path)
		s.handlers[sv] = h

	case 3:
		path := registry.Join(path, name)
		s.routers.Append(path, path, methods...)
		s.routerMap[path] = append(s.routerMap[path], path)
		s.handlers[path] = h
	}

}
func getMethods(name string) ([]string, int) {
	if name == "" {
		return []string{}, 1
	}
	for _, m := range defRequestMethod {
		if m == name {
			return []string{m}, 2
		}
	}
	return defRequestMethod, 3
}
