package zookeeper

import (
	"time"

	"github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/lib4go/logger"
	"github.com/micro-plat/lib4go/zk"
)

//zkRegistry 基于zookeeper的注册中心
type zkRegistry struct {
}

//Resolve 根据配置生成zookeeper客户端
func (z *zkRegistry) Resolve(servers []string, log *logger.Logger) (registry.IRegistry, error) {
	zclient, err := zk.NewWithLogger(servers, time.Second, log)
	if err != nil {
		return nil, err
	}
	err = zclient.Connect()
	return zclient, err
}

func init() {
	registry.Register("zk", &zkRegistry{})
}
