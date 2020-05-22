package services

import "fmt"

type serverServices struct {
	extHandle func(u *unit, ext ...string) error
	*services
	*handleHook
	*serverHook
	caches    map[string]interface{}
	cacheExts map[string][]string
}

func newServerServices(v func(u *unit, ext ...string) error) *serverServices {
	return &serverServices{
		handleHook: newHandleHook(),
		services:   newService(),
		serverHook: new(serverHook),
		caches:     make(map[string]interface{}),
		cacheExts:  make(map[string][]string),
		extHandle:  v,
	}
}

//Load 加载所有服务
func (s *serverServices) Load() error {
	for path, v := range s.caches {
		groups, err := reflectHandle(path, v)
		if err != nil {
			return err
		}
		if err := s.Add(groups); err != nil {
			return err
		}
	}
	return nil
}
func (s *serverServices) Cache(name string, h interface{}, ext ...string) {
	if _, ok := s.caches[name]; ok {
		panic(fmt.Sprintf("服务%s不能重复注册", name))
	}
	s.caches[name] = h
	if len(ext) > 0 {
		s.cacheExts[name] = ext
	}
}
func (s *serverServices) handleExt(g *unit) error {
	if s.extHandle == nil {
		return nil
	}
	return s.extHandle(g, s.cacheExts[g.path]...)
}

//Add 添加服务注册
func (s *serverServices) Add(g *unitGroup) error {
	for _, u := range g.Services {

		//添加预处理函数
		if err := s.handleHook.AddHandling(u.service, u.GetHandlings()...); err != nil {
			return err
		}

		//执行处理扩展函数
		if err := s.handleExt(u); err != nil {
			return err
		}

		//添加服务
		if err := s.services.AddHanler(u.service, u.Handle); err != nil {
			return err

		}

		//添加后处理函数
		if err := s.handleHook.AddHandled(u.service, u.GetHandleds()...); err != nil {
			return err
		}

		//添加降级函数
		if err := s.services.AddFallback(u.service, u.Fallback); err != nil {
			return err

		}
	}

	//添加关闭函数
	if err := s.handleHook.AddClosingHanle(g.Closing); err != nil {
		return err
	}

	return nil
}
