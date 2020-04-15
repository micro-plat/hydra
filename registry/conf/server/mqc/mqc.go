package mqc

//Server mqc服务配置
type Server struct {
	*option
}

//New 构建mqc server配置，默认为对等模式
func New(opts ...Option) *Server {
	s := &Server{option: &option{}}
	for _, opt := range opts {
		opt(s.option)
	}
	return s
}
