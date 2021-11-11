package services

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/micro-plat/hydra/conf/server/router"
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/registry"
)

type metaServices struct {
	pathActs   map[string][]string
	rawService map[string]*rawUnit
	groups     map[string]string
	handlers   map[string]context.IHandler
	fallbacks  map[string]context.IHandler
}

func newService() *metaServices {
	return &metaServices{
		pathActs:   make(map[string][]string),
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
	s.cachePathActs(service)
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

func (s *metaServices) cachePathActs(service string) {
	parties := strings.Split(service, "$")
	methods := router.DefMethods
	if len(parties) == 2 {
		methods = []string{parties[1]}
	}
	path := strings.TrimSuffix(parties[0], "/")

	acts, pok := s.pathActs[path]
	if !pok {
		acts = make([]string, 0)
	}
	acts = append(acts, methods...)
	s.pathActs[path] = acts
}

//Has 是否包含服务
func (s *metaServices) Has(service string) (ok bool) {
	_, ok = s.handlers[service]
	if ok {
		return true
	}
	parties := strings.Split(service, "$")
	if len(parties) != 2 {
		return false
	}
	path := strings.TrimSuffix(parties[0], "/")

	acts, pok := s.pathActs[path]
	if !pok {
		return false
	}
	method := parties[1]
	if strings.EqualFold(method, http.MethodOptions) {
		return true
	}
	for i := range acts {
		if strings.EqualFold(acts[i], method) {
			return true
		}
	}
	return false
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

//GetFallback 获取服务对应的降级函数
func (s *metaServices) GetFallback(service string) (h context.IHandler, ok bool) {
	h, ok = s.fallbacks[service]
	return
}

//GetFallback 获取服务对应的降级函数
func (s *metaServices) Remove(service string) {
	delete(s.handlers, service)
	delete(s.rawService, service)
	delete(s.groups, service)
	delete(s.fallbacks, service)

	parties := strings.Split(service, "$")
	for _, m := range s.pathActs[parties[0]] {
		rservice := registry.Join(parties[0], m)
		delete(s.handlers, rservice)
		delete(s.rawService, rservice)
		delete(s.groups, rservice)
		delete(s.fallbacks, rservice)
	}
}
