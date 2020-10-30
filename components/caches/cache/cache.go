package cache

import (
	"fmt"
)

type ICacheExt interface {
	GetProto() string
	GetServers() []string
}

//ICache 缓存接口
type ICache interface {
	Get(key string) (string, error)
	Decrement(key string, delta int64) (n int64, err error)
	Increment(key string, delta int64) (n int64, err error)
	Gets(key ...string) (r []string, err error)
	Add(key string, value string, expiresAt int) error
	Set(key string, value string, expiresAt int) error
	Delete(key string) error
	Exists(key string) bool
	Delay(key string, expiresAt int) error
	Close() error
}

//Resover 定义配置文件转换方法
type Resover interface {
	Resolve(conf string) (ICache, error)
}

var cacheResolvers = make(map[string]Resover)

//Register 注册配置文件适配器
func Register(proto string, resolver Resover) {
	if resolver == nil {
		panic("mq: Register adapter is nil")
	}
	if _, ok := cacheResolvers[proto]; ok {
		panic("mq: Register called twice for adapter " + proto)
	}
	cacheResolvers[proto] = resolver
}

//New 根据适配器名称及参数返回配置处理器
func New(proto string, conf string) (ICache, error) {
	resolver, ok := cacheResolvers[proto]
	if !ok {
		return nil, fmt.Errorf("cache: unknown adapter name %q (forgotten import?)", proto)
	}
	return resolver.Resolve(conf)
}
