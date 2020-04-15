package rpc

//Server rpc server配置信息
type Server struct {
	Address string `json:"address,omitempty" valid:"dialstring"`
	*option
}

//New 构建rpc server配置信息
func New(address string, opts ...Option) *Server {
	a := &Server{
		Address: address,
		option:  &option{},
	}
	for _, opt := range opts {
		opt(a.option)
	}
	return a
}
