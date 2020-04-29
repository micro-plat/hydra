package context

//Handler 业务处理Handler
type Handler func(IContext) interface{}

//Handle 处理业务流程
func (h Handler) Handle(c IContext) interface{} {
	return h(c)
}

//IHandler 业务处理接口
type IHandler interface {
	Handle(IContext) interface{}
}
