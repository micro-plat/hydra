package render

//Option 渲染可选参数
type Option func(*Render)

//TmpltOption 模板可选参数
type TmpltOption func(*Tmplt)

//WithTmplt 添加模板
func WithTmplt(path string, content string, opts ...TmpltOption) Option {
	return func(a *Render) {
		if _, ok := a.Tmplts[path]; !ok {
			a.Tmplts[path] = &Tmplt{Content: content}
		}
		for _, opt := range opts {
			opt(a.Tmplts[path])
		}
	}
}

//WithDisable 禁用配置
func WithDisable() Option {
	return func(a *Render) {
		a.Disable = true
	}
}

//WithEnable 启用配置
func WithEnable() Option {
	return func(a *Render) {
		a.Disable = false
	}
}

//WithStatus 添加状态码模板
func WithStatus(templt string) TmpltOption {
	return func(a *Tmplt) {
		a.Status = templt
	}
}

//WithContentType 添加content type模板
func WithContentType(templt string) TmpltOption {
	return func(a *Tmplt) {
		a.ContentType = templt
	}
}
