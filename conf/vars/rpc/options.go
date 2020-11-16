package rpc

import (
	"encoding/json"
	"fmt"
)

//Option 配置选项
type Option func(*RPCConf)

//WithConnectionTimeout 配置网络连接超时时长
func WithConnectionTimeout(t int) Option {
	return func(o *RPCConf) {
		o.ConntTimeout = t
	}
}

//WithLogger 配置日志记录器
func WithLogger(log string) Option {
	return func(o *RPCConf) {
		o.Log = log
	}
}

//WithTLS 设置TLS证书(pem,key)
func WithTLS(tls []string) Option {
	return func(o *RPCConf) {
		if len(tls) == 2 {
			o.Tls = tls
		}
	}
}

//WithRoundRobin 配置为轮询负载均衡器
func WithRoundRobin() Option {
	return func(o *RPCConf) {
		o.SortPrefix = ""
		o.RoundRobin = true
		o.LocalFirst = false
		o.LocalIP = ""
	}
}

//WithLocalFirst 配置为本地优先负载均衡器
func WithLocalFirst(local string) Option {
	return func(o *RPCConf) {
		o.SortPrefix = local
		o.LocalIP = local
		o.LocalFirst = true
		o.RoundRobin = false
	}
}

//WithRaw 根据json串设置配置信息
func WithRaw(raw []byte) Option {
	c := &RPCConf{}
	if err := json.Unmarshal(raw, c); err != nil {
		panic(fmt.Errorf("rpc配置节点解析异常,%v", err))
	}
	return func(o *RPCConf) {
		o = c
	}
}
