package redis

//Redis redis缓存配置
type Redis struct {
	Address string `json:"addrs,omitempty"`
	*option
}

//New 构建redis消息队列配置
func New(address string, opts ...Option) *Redis {
	r := &Redis{
		Address: address,
		option: &option{
			Proto:        "redis",
			DbIndex:      1,
			DialTimeout:  10,
			ReadTimeout:  10,
			WriteTimeout: 10,
			PoolSize:     10,
		},
	}
	for _, opt := range opts {
		opt(r.option)
	}
	return r
}
