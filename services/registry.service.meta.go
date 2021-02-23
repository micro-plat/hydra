package services

import (
	"fmt"

	"github.com/micro-plat/hydra/context"
)

type metaServices struct {
	services   []string
	rawService map[string]*rawUnit
	groups     map[string]string
	handlers   map[string]context.IHandler
	fallbacks  map[string]context.IHandler
}

func newService() *metaServices {
	return &metaServices{
		services:   make([]string, 0, 1),
		rawService: make(map[string]*rawUnit),
		groups:     make(map[string]string),
		handlers:   make(map[string]context.IHandler),
		fallbacks:  make(map[string]context.IHandler),
	}
}

func (s *metaServices) AddHanler(service string, group string, h context.IHandler, r *rawUnit) error {
	if _, ok := s.handlers[service]; ok {
		return fmt.Errorf("服务不能重复注册，%s找到有多次注册%v", service, s.handlers)
	}
	s.handlers[service] = h
	s.rawService[service] = r
	s.groups[service] = group
	s.services = append(s.services, service)
	return nil
}

//AddFallback 添加
func (s *metaServices) AddFallback(service string, h context.IHandler) error {
	if h == nil {
		return nil
	}
	s.fallbacks[service] = h
	return nil
}

func (s *metaServices) remove(service string) {
	delete(s.handlers, service)
	for i, srv := range s.services {
		if srv == service {
			s.services = append(s.services[:i], s.services[i+1:]...)
			return
		}
	}
}

//Has 是否包含服务
func (s *metaServices) Has(service string) (ok bool) {
	_, ok = s.handlers[service]
	return
}

//GetHandlers 获取服务的处理对象
func (s *metaServices) GetHandlers(service string) (h context.IHandler, ok bool) {
	h, ok = s.handlers[service]
	return
}

//GetGroup 获取服务的分组信息
func (s *metaServices) GetGroup(service string) string {
	return s.groups[service]
}

//GetRawPathAndTag 获取服务原始注册路径与方法名
func (s *metaServices) GetRawPathAndTag(service string) (path string, tagName string, ok bool) {
	u, ok := s.rawService[service]
	if ok {
		return u.RawPath, u.RawMTag, true
	}
	return "", "", false
}

//GetServices 获取已注册的服务
func (s *metaServices) GetServices() []string {
	return s.services
}

//GetFallback 获取服务对应的降级函数
func (s *metaServices) GetFallback(service string) (h context.IHandler, ok bool) {
	h, ok = s.fallbacks[service]
	return
}
