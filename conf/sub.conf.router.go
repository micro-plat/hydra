package conf

type Routers struct {
	Routers []*Router         `json:"routers,omitempty"`
	RPCS    []*Router         `json:"rpcs,omitempty"`
	Setting map[string]string `json:"args,omitempty"`
}

type Router struct {
	Name    string            `json:"name,omitempty" valid:"ascii,required"`
	Action  []string          `json:"action,omitempty" valid:"uppercase,in(GET|POST|PUT|DELETE|HEAD|TRACE|OPTIONS)"`
	Engine  string            `json:"-" valid:"-"`
	Service string            `json:"service" valid:"ascii,required"`
	Setting map[string]string `json:"args,omitempty"`
	Disable bool              `json:"disable,omitempty"`
	Handler interface{}       `json:"-" valid:"-"`
}

//NewRouters 构建路由
func NewRouters() *Routers {
	return &Routers{
		Routers: make([]*Router, 0),
		RPCS:    make([]*Router, 0),
	}
}

//Append 添加路由信息
func (h *Routers) Append(name string, service string) *Routers {
	h.Routers = append(h.Routers, &Router{
		Name:    name,
		Service: service,
	})
	return h
}

//AppendRPCProxy 添加proxy信息
func (h *Routers) AppendRPCProxy(name string, service string, settings ...map[string]string) *Routers {
	router := &Router{Name: name, Service: service}
	if len(settings) > 0 {
		router.Setting = settings[0]
	}
	h.RPCS = append(h.RPCS, router)
	return h
}

//AppendWithAction 添加路由信息
func (h *Routers) AppendWithAction(name string, service string, action ...string) *Routers {
	h.Routers = append(h.Routers, &Router{
		Name:    name,
		Service: service,
		Action:  action,
	})
	return h
}
