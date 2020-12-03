package redis

import (
	"fmt"

	"github.com/micro-plat/hydra/registry"
)

//redisFactory 基于redis系统
type redisFactory struct {
	opts *registry.Options
}

//Build 根据配置生成文件系统注册中心
func (z *redisFactory) Create(opts ...registry.Option) (registry.IRegistry, error) {
	//addrs []string, u string, p string, log logger.ILogging
	for i := range opts {
		opts[i](z.opts)
	}

	if len(z.opts.Addrs) == 0 {
		return nil, fmt.Errorf("未指定redis服务器地址")
	}

	return NewRegistry(z.opts)
}

func init() {
	registry.Register(registry.Redis, &redisFactory{
		opts: &registry.Options{},
	})
}
