package filesystem

import (
	"github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/lib4go/logger"
)

//zkRegistry 基于zookeeper的注册中心
type fileBuilder struct {
}

//Build 根据配置生成文件系统注册中心
func (z *fileBuilder) Build(addrs []string, u string, p string, log logger.ILogging) (registry.IRegistry, error) {
	prefix := "."
	if len(servers) > 0 {
		prefix = servers[0]
	}
	client, err := newFileSystem(prefix)
	if err != nil {
		return nil, err
	}
	client.Start()
	return client, nil
}

func init() {
	registry.Register("fs", &fileBuilder{})
}
