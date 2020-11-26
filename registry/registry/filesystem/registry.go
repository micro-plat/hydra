package filesystem

import (
	"fmt"

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

	if len(z.opts.Addrs) <= 0 {
		return nil, fmt.Errorf("FS注册中心，需要指定一个本地目录地址：%v", z.opts.Addrs)
	}
	if len(z.opts.Addrs) > 1 {
		return nil, fmt.Errorf("FS注册中心，只允许传递一个本地目录地址：%v", z.opts.Addrs)
	}
	rootDir := z.opts.Addrs[0]
	client, err := NewFileSystem(rootDir)
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
