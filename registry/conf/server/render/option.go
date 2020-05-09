package render

//option 配置参数
type option struct {

	//模板内容
	Content string `json:"content,omitempty" valid:"required" toml:"content,omitempty"`

	//适用的服务
	Services []string `json:"services,omitempty" valid:"required" toml:"services,omitempty"`

	//响应状态码
	Status string `json:"status,omitempty" valid:"required" toml:"status,omitempty"`
}

//Option 配置选项
type Option func(*option)

//WithItemByService 设置普通模板
func WithItemByService(content string, service ...string) Option {
	return func(a *option) {
		if len(service) > 0 {
			a.Services = append(a.Services, service...)
		}
	}
}

//WithItemByStatus 设置状态码模板
func WithItemByStatus(statusTmpl string, content string, service ...string) Option {
	return func(a *option) {
		a.Status = statusTmpl
		a.Content = content
		if len(service) > 0 {
			a.Services = append(a.Services, service...)
		}
	}
}
