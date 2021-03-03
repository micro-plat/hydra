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
