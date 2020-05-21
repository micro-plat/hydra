package services

import (
	"fmt"

	"github.com/micro-plat/hydra/context"
)

type httpServices struct {
	*services
	*hook
	*apiRouter
}

func newServices() *httpServices {
	return &httpServices{
		hook:      newHook(),
		services:  newService(),
		apiRouter: newAPIRouter(),
	}
}

func (s *httpServices) AddHanler(path string, name string, h context.IHandler) error {
	rpath, service, action := getPaths(path, name)
	if h == nil {
		return fmt.Errorf("%s不是有效的服务处理函数", service)
	}

	//添加服务
	if err := s.services.add(service, h); err != nil {
		return err

	}
	//添加路由
	if err := s.apiRouter.add(rpath, service, action...); err != nil {
		s.services.remove(service)
		return err
	}

	s.dealHandling(path, service)
	s.dealHandled(path, service)

	return nil
}
