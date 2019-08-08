package conf

type Routers struct {
	Setting map[string]string `json:"args,omitempty"`
	Routers []*Router         `json:"routers,omitempty"`
}
type Router struct {
	Name    string            `json:"name,omitempty" valid:"ascii,required"`
	Action  []string          `json:"action,omitempty" valid:"uppercase,in(GET|POST|PUT|DELETE|HEAD|TRACE|OPTIONS)"`
	Engine  string            `json:"engine,omitempty" valid:"ascii,uppercase,in(*|RPC)"`
	Service string            `json:"service" valid:"ascii,required"`
	Setting map[string]string `json:"args,omitempty"`
	Disable bool              `json:"disable,omitempty"`
	Handler interface{}
}

//NewRouters 构建路由
func NewRouters() *Routers {
	return &Routers{
		Routers: make([]*Router, 0),
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

//AppendWithAction 添加路由信息
func (h *Routers) AppendWithAction(name string, service string, action ...string) *Routers {
	h.Routers = append(h.Routers, &Router{
		Name:    name,
		Service: service,
		Action:  action,
	})
	return h
}
