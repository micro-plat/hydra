package zookeeper

import (
	"fmt"
	"time"

	"github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/lib4go/zk"
)

//zookeeper 基于zookeeper的注册中心
type zookeeperFactory struct {
	opts *registry.Options
}

//Create 根据配置生成zookeeper客户端
func (z *zookeeperFactory) Create(opts ...registry.Option) (registry.IRegistry, error) {

	for i := range opts {
		opts[i](z.opts)
	}

	if len(z.opts.Addrs) == 0 {
		return nil, fmt.Errorf("未指定zk服务器地址")
	}
	zclient, err := zk.NewWithLogger(z.opts.Addrs, time.Second, z.opts.Logger, zk.WithdDigest(z.opts.Auth.Username, z.opts.Auth.Password))
	if err != nil {
		return nil, err
	}
	err = zclient.Connect()
	return zclient, err
}

func init() {
	registry.Register(registry.Zookeeper, &zookeeperFactory{
		opts: &registry.Options{},
	})
}
