package mqtt

//Redis redis缓存配置
type Redis struct {
	Address string `json:"address,omitempty"`
	*option
}

//New 构建mqtt配置
func New(address string, opts ...Option) *Redis {
	r := &Redis{
		Address: address,
		option: &option{
			Proto: "mqtt",
		},
	}
	for _, opt := range opts {
		opt(r.option)
	}
	return r
}
