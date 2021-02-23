package services

type serverServices struct {
	extHandle func(u *Unit, ext ...interface{}) error
	*metaServices
	*handleHook
	*serverHook
}

func newServerServices(v func(u *Unit, ext ...interface{}) error) *serverServices {
	return &serverServices{
		handleHook:   newHandleHook(),
		metaServices: newService(),
		serverHook:   new(serverHook),
		extHandle:    v,
	}
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

//addGroup 添加服务注册
func (s *serverServices) addGroup(g *UnitGroup, group string, ext ...interface{}) error {
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
