package mqtt

import "encoding/json"

//Option 配置选项
type Option func(*MQTT)

//WithUP 设置用户名密码
func WithUP(userName string, password string) Option {
	return func(a *MQTT) {
		a.UserName = userName
		a.Password = password
	}
}

//WithCert 设置证书地址
func WithCert(cert string) Option {
	return func(a *MQTT) {
		a.Cert = cert
	}
}

//WithDialTimeout 设置连接超时
func WithDialTimeout(timeout int64) Option {
	return func(a *MQTT) {
		a.DialTimeout = timeout
	}
}

//WithRaw 通过json原串初始化
func WithRaw(raw string) Option {
	return func(o *MQTT) {
		if err := json.Unmarshal([]byte(raw), o); err != nil {
			panic(err)
		}
	}
}
