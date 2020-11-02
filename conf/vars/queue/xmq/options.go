package xmq

import "encoding/json"

//Option 配置选项
type Option func(*XMQ)

//WithConfigName 设置数据库分片索引
func WithAddress(address string) Option {
	return func(a *XMQ) {
		a.Address = address
	}
}

//WithRaw 通过json原串初始化
func WithRaw(raw string) Option {
	return func(o *XMQ) {
		if err := json.Unmarshal([]byte(raw), o); err != nil {
			panic(err)
		}
	}
}
