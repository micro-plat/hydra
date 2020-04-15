package header

//Headers http头信息
type Headers = option

//NewHeader 构建请求的头配置
func NewHeader(opts ...Option) Headers {
	h := newOption()
	for _, opt := range opts {
		opt(h)
	}
	return h
}
