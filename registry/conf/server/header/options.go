package header

import "strings"

var allow = []string{"X-Add-Delay", "X-Request-Id", "X-Requested-With", "Content-Type", "hsid"}
var expose = []string{"hsid"}

type option = map[string]string

//Option 配置选项
type Option func(option)

func newOption() option {
	return map[string]string{
		"Access-Control-Allow-Credentials": "true",
		"Access-Control-Allow-Origin":      "*",
		"Access-Control-Allow-Methods":     "GET,POST,PUT,DELETE,PATCH,OPTIONS",
		"Access-Control-Allow-Headers":     strings.Join(allow, ","),
		"Access-Control-Expose-Headers":    strings.Join(expose, ","),
	}
}

//WithHosts 设置允许的主机名
func WithHosts(host ...string) Option {
	return func(a option) {
		a["Access-Control-Allow-Origin"] = strings.Join(host, ",")
	}
}

//WithMethods 设置允许的请求类型
func WithMethods(method ...string) Option {
	return func(a option) {
		a["Access-Control-Allow-Methods"] = strings.ToUpper(strings.Join(method, ","))
	}
}

//WithAllowHeaders 设置允许请求的头信息
func WithAllowHeaders(header ...string) Option {
	return func(a option) {
		a["Access-Control-Allow-Headers"] = strings.Join(append(allow, header...), ",")
	}
}

//WithExposeHeaders 设置允许导出的头信息
func WithExposeHeaders(header ...string) Option {
	return func(a option) {
		a["Access-Control-Expose-Headers"] = strings.Join(append(expose, header...), ",")
	}
}

//WithHeader 设置其它头信息
func WithHeader(kv ...string) Option {
	return func(a option) {
		l := len(kv)
		for i := 0; i < len(kv)/2 && i < l-1; i++ {
			a[kv[i*2]] = kv[i*2+1]
		}
	}
}
