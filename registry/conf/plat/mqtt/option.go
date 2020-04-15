package mqtt

type option struct {
	Proto    string `json:"proto,omitempty"`
	UserName string `json:"userName,omitempty"`
	Password string `json:"password,omitempty"`
	Cert     string `json:"cert,omitempty"`
}

//Option 配置选项
type Option func(*option)

//WithUP 设置用户名密码
func WithUP(userName string, password string) Option {
	return func(a *option) {
		a.UserName = userName
		a.Password = password
	}
}

//WithCert 设置证书地址
func WithCert(cert string) Option {
	return func(a *option) {
		a.Cert = cert
	}
}
