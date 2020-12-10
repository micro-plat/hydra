package memcached

import (
	"encoding/json"
)

//Option 配置选项
type Option func(*Memcache)

//WithAddress 设置哨兵服务器
func WithAddress(address ...string) Option {
	return func(o *Memcache) {
		o.Address = append(o.Address, address...)
	}
}

//WithTimeout 设置连接超时时长
func WithTimeout(timeout int) Option {
	return func(o *Memcache) {
		o.Timeout = timeout
	}
}

//WithMaxIdleConns 设置连接空闲数
func WithMaxIdleConns(idleConns int) Option {
	return func(o *Memcache) {
		o.MaxIdleConns = idleConns
	}
}

//WithRaw 通过json原串初始化
func WithRaw(raw string) Option {
	return func(o *Memcache) {
		if err := json.Unmarshal([]byte(raw), o); err != nil {
			panic(err)
		}
	}
}
