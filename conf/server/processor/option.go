package processor

import "strings"

//Option 配置选项
type Option func(*Processor)

//WithServicePrefix 服务前缀
func WithServicePrefix(prefix string) Option {
	return func(a *Processor) {
		a.ServicePrefix = "/" + strings.TrimSuffix(strings.TrimPrefix(prefix, "/"), "/")
	}
}

//WithEnableEncryption 启用加密设置
func WithEnableEncryption() Option {
	return func(a *Processor) {
		a.EnableEncryption = true
	}
}
