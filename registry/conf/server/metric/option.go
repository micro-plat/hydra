package metric

//option 配置参数
type option struct {
	UserName string `json:"userName,omitempty" valid:"ascii"`
	Password string `json:"password,omitempty" valid:"ascii"`
	Disable  bool   `json:"disable,omitempty"`
}

//Option 配置选项
type Option func(*option)

//WithUPName 设置用户名密码
func WithUPName(userName string, password string) Option {
	return func(a *option) {
		a.UserName = userName
		a.Password = password
	}
}

//WithDisable 禁用配置
func WithDisable() Option {
	return func(a *option) {
		a.Disable = true
	}
}

//WithEnable 启用配置
func WithEnable() Option {
	return func(a *option) {
		a.Disable = false
	}
}
