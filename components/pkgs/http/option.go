package http

import (
	"encoding/json"
	"time"

	"github.com/micro-plat/hydra/context"
)

//conf 配置信息
type conf struct {
	RequestId         string              `json:"-"`
	apmCtx            context.IAPMContext `json:"-"`
	ConnectionTimeout time.Duration       `json:"ctime"`
	RequestTimeout    time.Duration       `json:"rtime"`
	Certs             []string            `json:"certs"`
	Ca                string              `json:"ca"`
	Proxy             string              `json:"proxy"`
	Keepalive         bool                `json:"keepalive"`
	Trace             bool                `json:"trace"`
}

//Option 配置选项
type Option func(*conf)

//WithConnTimeout 设置请求超时时长
func WithConnTimeout(tm time.Duration) Option {
	return func(o *conf) {
		o.ConnectionTimeout = tm
	}
}

//WithRequestTimeout 设置请求超时时长
func WithRequestTimeout(tm time.Duration) Option {
	return func(o *conf) {
		o.RequestTimeout = tm
	}
}

//WithCert 设置请求证书
func WithCert(cerfile string, key string) Option {
	return func(o *conf) {
		o.Certs = []string{cerfile, key}
	}
}

//WithCa 设置ca证书
func WithCa(cafile string) Option {
	return func(o *conf) {
		o.Ca = cafile
	}
}

//WithProxy 使用代理地址
func WithProxy(proxy string) Option {
	return func(o *conf) {
		o.Proxy = proxy
	}
}

//WithKeepalive 设置keep alive
func WithKeepalive(keepalive bool) Option {
	return func(o *conf) {
		o.Keepalive = keepalive
	}
}

//WithRequestID 设置请求编号
func WithRequestID(requestID string) Option {
	return func(o *conf) {
		o.RequestId = requestID
	}
}

//WithRaw 根据json串设置配置信息
func WithRaw(raw []byte) (Option, error) {
	c := &conf{}
	if err := json.Unmarshal(raw, c); err != nil {
		panic(err)
	}
	return func(o *conf) {
		o = c
	}, nil
}
