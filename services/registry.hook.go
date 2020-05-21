package services

import (
	"fmt"

	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/registry/conf/server"
)

type hook struct {
	handlingGroup  map[string][]string
	handledGroup   map[string][]string
	starting       func(server.IServerConf) error
	closing        func(server.IServerConf) error
	globalHandling context.IHandler
	routerHandling map[string]context.IHandler
	handlings      map[string][]context.IHandler
	globalHandled  context.IHandler
	routerHandled  map[string]context.IHandler
	handleds       map[string][]context.IHandler
	fallbacks      map[string]context.IHandler
	closers        []func() error
}

func newHook() *hook {
	return &hook{
		handlingGroup:  make(map[string][]string),
		handledGroup:   make(map[string][]string),
		routerHandling: make(map[string]context.IHandler),
		handlings:      make(map[string][]context.IHandler),
		routerHandled:  make(map[string]context.IHandler),
		handleds:       make(map[string][]context.IHandler),
		fallbacks:      make(map[string]context.IHandler),
		closers:        make([]func() error, 0, 1),
	}
}

//Starting 设置启动服务
func (s *hook) Starting(h func(server.IServerConf) error) error {
	if h == nil {
		return fmt.Errorf("启动服务不能为空")
	}
	if s.starting != nil {
		return fmt.Errorf("启动服务不能重复注册")
	}
	s.starting = h
	return nil
}

//Closing 设置关闭服务
func (s *hook) Closing(h func(server.IServerConf) error) error {
	if h == nil {
		return fmt.Errorf("关闭服务不能为空")
	}
	if s.closing != nil {
		return fmt.Errorf("启动服务不能重复注册")
	}
	s.closing = h
	return nil
}
func (s *hook) GlobalHandling(h context.IHandler) error {
	if h == nil {
		return fmt.Errorf("不是有效的服务处理函数")
	}
	if s.globalHandling != nil {
		return fmt.Errorf("不能重复注册")
	}
	s.globalHandling = h
	return nil
}
func (s *hook) GlobalHandled(h context.IHandler) error {
	if h == nil {
		return fmt.Errorf("不是有效的服务处理函数")
	}
	if s.globalHandled != nil {
		return fmt.Errorf("不能重复注册")
	}
	s.globalHandled = h
	return nil
}
func (s *hook) AddHandling(path string, name string, h context.IHandler) error {

	//获取当前handling的服务名称
	_, srvs, _ := getPaths(path, name)
	if h == nil {
		return fmt.Errorf("%s不是有效的服务处理函数", srvs)
	}

	//初始化当前服务对应的handling
	if _, ok := s.handlings[srvs]; !ok {
		s.handlings[srvs] = make([]context.IHandler, 0, 0)
		if s.globalHandling != nil {
			s.handlings[srvs] = append(s.handlings[srvs], s.globalHandling)
		}
	}

	//广播分组下所有服务
	if name == "" {
		s.routerHandling[path] = h
		for _, v := range s.handlingGroup[path] {
			s.handlings[v] = append([]context.IHandler{h}, s.handlings[v]...)
		}
	}
	//添加到当前处理handling
	s.handlings[srvs] = append(s.handlings[srvs], h)

	return nil
}

func (s *hook) dealHandling(path string, service string) {
	s.handlingGroup[path] = append(s.handlingGroup[path], service)
	if v := s.routerHandling[path]; v != nil {
		s.handlings[service] = append(s.handlings[service], v)
	}
}
func (s *hook) dealHandled(path string, service string) {
	s.handledGroup[path] = append(s.handledGroup[path], service)
	if v := s.routerHandled[path]; v != nil {
		s.handleds[service] = append(s.handleds[service], v)
	}
}

func (s *hook) AddHandled(path string, name string, h context.IHandler) error {
	//获取当前handled的服务名称
	_, srvs, _ := getPaths(path, name)
	if h == nil {
		return fmt.Errorf("%s不是有效的服务处理函数", srvs)
	}

	//初始化当前服务对应的handled
	if _, ok := s.handleds[srvs]; !ok {
		s.handleds[srvs] = make([]context.IHandler, 0, 0)
		if s.globalHandled != nil {
			s.handleds[srvs] = append(s.handleds[srvs], s.globalHandling)
		}

	}

	//广播分组下所有服务
	if name == "" {
		s.routerHandled[path] = h
		for _, v := range s.handledGroup[path] {
			s.handleds[v] = append(s.handleds[v], h)
		}
	}
	//添加到当前处理handled
	s.handleds[srvs] = append([]context.IHandler{h}, s.handleds[srvs]...)

	return nil
}
func (s *hook) AddFallback(path string, name string, h context.IHandler) error {
	_, srvs, _ := getPaths(path, name)
	if h == nil {
		return fmt.Errorf("%s不是有效的服务处理函数", srvs)
	}
	s.fallbacks[srvs] = h
	return nil
}

//AddCloser 添加关闭函数
func (s *hook) AddCloser(path string, h interface{}) error {
	if vv, ok := h.(func() error); ok {
		s.closers = append(s.closers, vv)
		return nil
	}
	if vv, ok := h.(func()); ok {
		s.closers = append(s.closers, func() error {
			vv()
			return nil
		})
		return nil
	}
	return fmt.Errorf("%s提供的close签名类型不支持", path)
}

func (s *hook) GetHandlings(service string) (h []context.IHandler) {
	h = make([]context.IHandler, 0, 1)
	if s.globalHandling != nil {
		h = append(h, s.globalHandling)
	}
	return append(h, s.handlings[service]...)
}

func (s *hook) GetHandleds(service string) (h []context.IHandler) {
	h = make([]context.IHandler, 0, 1)
	h = append(h, s.handleds[service]...)
	if s.globalHandled != nil {
		h = append(h, s.globalHandled)
	}
	return h
}
func (s *hook) GetFallback(service string) (h context.IHandler, ok bool) {
	h, ok = s.fallbacks[service]
	return
}

func (s *hook) GetStartingHandles() []func(server.IServerConf) error {
	if s.starting == nil {
		return nil
	}
	return []func(server.IServerConf) error{s.starting}
}
func (s *hook) GetClosingHandles() []func(server.IServerConf) error {
	if s.closing == nil {
		return nil
	}
	return []func(server.IServerConf) error{s.closing}
}

//GetClosers 获取提供有资源释放的服务
func (s *hook) GetClosers() []func() error {
	return s.closers
}
