package render

//Option jwt配置选项
type Option func(Tmplt)

//WithStatus 设置状态码对应的模板
func WithStatus(tmplt string) Option {
	return func(a Tmplt) {
		a["render/status"] = tmplt
	}
}

//WithContent 设置响应内容的模板
func WithContent(tmplt string) Option {
	return func(a Tmplt) {
		a["render/content"] = tmplt
	}
}
