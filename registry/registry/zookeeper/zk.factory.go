package zookeeper

import (
	"fmt"
	"time"

	"github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/lib4go/logger"
	"github.com/micro-plat/lib4go/zk"
)

//zookeeper 基于zookeeper的注册中心
type zookeeperFactory struct {
}

//Create 根据配置生成zookeeper客户端
func (z *zookeeperFactory) Create(addrs []string, u string, p string, log logger.ILogging) (registry.IRegistry, error) {
	if len(addrs) == 0 {
		return nil, fmt.Errorf("未指定zk服务器地址")
	}
	zclient, err := zk.NewWithLogger(addrs, time.Second, log, zk.WithdDigest(u, p))
	if err != nil {
		return nil, err
	}
	err = zclient.Connect()
	return zclient, err
}

func init() {
	registry.Register(registry.Zookeeper, &zookeeperFactory{})
}
