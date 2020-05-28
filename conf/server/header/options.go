package header

import (
	"net/http"
	"strings"
)

var allowHeader = []string{"X-Add-Delay", "X-Request-Id", "X-Requested-With", "Content-Type", "Authorization-Jwt", "Origin", "Accept"}
var exposeHeader = []string{"Authorization-Jwt"}
var allMethods = []string{http.MethodHead, http.MethodOptions, http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete}

//Option 配置选项
type Option func(Headers)

//WithCrossDomain 添加跨域配置
func WithCrossDomain(host ...string) Option {
	return func(a Headers) {
		origin := "*"
		if len(host) > 0 {
			origin = strings.Join(host, ",")
		}
		a["Access-Control-Allow-Credentials"] = "true"
		a["Access-Control-Allow-Origin"] = origin
		a["Access-Control-Allow-Methods"] = strings.Join(allMethods, ",")
		a["Access-Control-Allow-Headers"] = strings.Join(allowHeader, ",")
		a["Access-Control-Expose-Headers"] = strings.Join(exposeHeader, ",")
	}
}

//WithAllowMethods 设置允许的请求类型
func WithAllowMethods(method ...string) Option {
	return func(a Headers) {
		a["Access-Control-Allow-Methods"] = strings.ToUpper(strings.Join(method, ","))
	}
}

//WithAllowHeaders 设置允许请求的头信息
func WithAllowHeaders(header ...string) Option {
	return func(a Headers) {
		a["Access-Control-Allow-Headers"] = strings.Join(append(allowHeader, header...), ",")
	}
}

//WithExposeHeaders 设置允许导出的头信息
func WithExposeHeaders(header ...string) Option {
	return func(a Headers) {
		a["Access-Control-Expose-Headers"] = strings.Join(append(exposeHeader, header...), ",")
	}
}

//WithHeader 设置其它头信息
func WithHeader(kv ...string) Option {
	return func(a Headers) {
		l := len(kv)
		for i := 0; i < len(kv)/2 && i < l-1; i++ {
			a[kv[i*2]] = kv[i*2+1]
		}
	}
}
