package services

import "strings"

type serverServices struct {
	unitGroups []*UnitGroup
	extHandle  func(u *Unit, ext ...interface{}) error
	extRemove  func(path string)
	*metaServices
	*handleHook
	*serverHook
}

func newServerServices(v func(u *Unit, ext ...interface{}) error, remove func(path string)) *serverServices {

	s := &serverServices{
		unitGroups:   make([]*UnitGroup, 0, 1),
		handleHook:   newHandleHook(),
		metaServices: newService(),
		serverHook:   newServerHook(),
		extHandle:    v,
		extRemove:    remove,
	}
	if s.extHandle == nil {
		s.extHandle = func(u *Unit, ext ...interface{}) error { return nil }
		s.extRemove = func(path string) {}
	}
	return s
}

//Register 注册服务
func (s *serverServices) Register(group string, name string, h interface{}, ext ...interface{}) {
	groups, err := reflectHandle(name, h)
	if err != nil {
		panic(err)
	}
	if err := s.addGroup(groups, group, ext...); err != nil {
		panic(err)
	}
}
func (s *serverServices) handleExt(u *Unit, ext ...interface{}) error {
	if s.extHandle == nil {
		return nil
	}
	return s.extHandle(u, ext...)
}

//Remove 移除已注册的服务
func (s *serverServices) Remove(path string) {
	for _, g := range s.unitGroups {
		if strings.EqualFold(g.Path, path) {
			for _, u := range g.Services {
				s.handleHook.Remove(u.Service)
				s.metaServices.Remove(u.Service)
			}
			s.extRemove(path)
		}
	}
}

//addGroup 添加服务注册
func (s *serverServices) addGroup(g *UnitGroup, group string, ext ...interface{}) error {
	s.unitGroups = append(s.unitGroups, g)
	for _, u := range g.Services {

		//添加预处理函数
		if err := s.handleHook.AddHandling(u.Service, u.GetHandlings()...); err != nil {
			return err
		}

		//执行处理扩展函数
		if err := s.handleExt(u, ext...); err != nil {
			return err
		}

		//添加服务
		if err := s.metaServices.AddHanler(u.Service, group, u.Handle, u.rawUnit); err != nil {
			return err

		}

		//添加后处理函数
		if err := s.handleHook.AddHandled(u.Service, u.GetHandleds()...); err != nil {
			return err
		}

		//添加降级函数
		if err := s.metaServices.AddFallback(u.Service, u.Fallback); err != nil {
			return err

		}
	}

	//添加关闭函数
	if err := s.handleHook.AddClosingHanle(g.Closing); err != nil {
		return err
	}

	return nil
}
