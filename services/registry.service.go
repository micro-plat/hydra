package services

import (
	"fmt"

	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/registry"
)

type services struct {
	services []string
	handlers map[string]context.IHandler
}

func newService() *services {
	return &services{
		services: make([]string, 0, 1),
		handlers: make(map[string]context.IHandler),
	}
}

func (s *services) add(service string, h context.IHandler) error {
	if _, ok := s.handlers[service]; ok {
		return fmt.Errorf("服务不能重复注册，%s找到有多次注册", service)
	}
	s.handlers[service] = h
	s.services = append(s.services, service)
	return nil
}
func (s *services) remove(service string) {
	delete(s.handlers, service)
	for i, srv := range s.services {
		if srv == service {
			s.services = append(s.services[:i], s.services[i+1:]...)
			return
		}
	}
}

func (s *services) GetHandlers(service string) (h context.IHandler, ok bool) {
	h, ok = s.handlers[service]
	return
}

//GetServices 获取已注册的服务
func (s *services) GetServices() []string {
	return s.services
}

func getPaths(path string, name string) (rpath string, service string, action []string) {
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
