package cache

import (
	"fmt"
	"strings"
)

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

//CacheResover 定义配置文件转换方法
type CacheResover interface {
	Resolve(address []string, conf string) (ICache, error)
}

var cacheResolvers = make(map[string]CacheResover)

//Register 注册配置文件适配器
func Register(proto string, resolver CacheResover) {
	if resolver == nil {
		panic("mq: Register adapter is nil")
	}
	if _, ok := cacheResolvers[proto]; ok {
		panic("mq: Register called twice for adapter " + proto)
	}
	cacheResolvers[proto] = resolver
}

//NewCache 根据适配器名称及参数返回配置处理器
func NewCache(address string, conf string) (ICache, error) {
	proto, addrs, err := getNames(address)
	if err != nil {
		return nil, err
	}
	resolver, ok := cacheResolvers[proto]
	if !ok {
		return nil, fmt.Errorf("cache: unknown adapter name %q (forgotten import?)", proto)
	}
	return resolver.Resolve(addrs, conf)
}

func getNames(address string) (proto string, raddr []string, err error) {
	addr := strings.SplitN(address, "://", 2)
	if len(addr[0]) == 0 {
		return "", nil, fmt.Errorf("cache地址配置错误%s，格式:memcached://192.168.0.1:11211", addr)
	}
	proto = addr[0]
	var rightAddr string
	if len(addr) > 1 {
		rightAddr = addr[1]
	}
	raddr = strings.Split(rightAddr, ",")
	return
}
