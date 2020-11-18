package filesystem

import (
	"github.com/micro-plat/hydra/registry"
)

//zkRegistry 基于zookeeper的注册中心
type fsRegistry struct {
	opts *registry.Options
}

//Resolve 根据配置生成zookeeper客户端
func (z *fsRegistry) Create(opts ...registry.Option) (registry.IRegistry, error) {
	for i := range opts {
		opts[i](z.opts)
	}
	prefix := "."
	client, err := newLocal(prefix)
	if err != nil {
		return nil, err
	}
	client.Start()
	return client, nil
}

func init() {
	registry.Register(registry.FileSystem, &fsRegistry{
		opts: &registry.Options{},
	})
}
