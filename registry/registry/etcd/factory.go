package etcd

import (
	"github.com/micro-plat/hydra/registry"
)

//etcdFactory 基于etcd系统
type etcdFactory struct {
	opts *registry.Options
}

//Build 根据配置生成文件系统注册中心
func (z *etcdFactory) Create(opts ...registry.Option) (registry.IRegistry, error) {
	//addrs []string, u string, p string, log logger.ILogging
	for i := range opts {
		opts[i](z.opts)
	}
	reg := NewRegistry(opts...)
	return reg, nil
}

func init() {
	registry.Register(registry.Etcd, &etcdFactory{
		opts: &registry.Options{},
	})
}
