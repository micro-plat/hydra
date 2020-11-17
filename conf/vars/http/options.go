package http

import (
	"encoding/json"
	"fmt"
)

//Option 配置选项
type Option func(*HTTPConf)

//WithConnTimeout 设置请求超时时长
func WithConnTimeout(tm int) Option {
	return func(o *HTTPConf) {
		o.ConnectionTimeout = tm
	}
}

//WithRequestTimeout 设置请求超时时长
func WithRequestTimeout(tm int) Option {
	return func(o *HTTPConf) {
		o.RequestTimeout = tm
	}
}

//WithCert 设置请求证书
func WithCert(cerfile string, key string) Option {
	return func(o *HTTPConf) {
		o.Certs = []string{cerfile, key}
	}
}

//WithCa 设置ca证书
func WithCa(cafile string) Option {
	return func(o *HTTPConf) {
		o.Ca = cafile
	}
}

//WithProxy 使用代理地址
func WithProxy(proxy string) Option {
	return func(o *HTTPConf) {
		o.Proxy = proxy
	}
}

//WithKeepalive 设置keep alive
func WithKeepalive(keepalive bool) Option {
	return func(o *HTTPConf) {
		o.Keepalive = keepalive
	}
}

//WithRaw 根据json串设置配置信息
func WithRaw(raw []byte) Option {
	c := &HTTPConf{}
	if err := json.Unmarshal(raw, c); err != nil {
		panic(fmt.Errorf("http配置节点解析异常,%v", err))
	}
	return func(o *HTTPConf) {
		o = c
	}
}
