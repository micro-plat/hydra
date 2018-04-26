package mqc

import (
	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/servers"
	"github.com/micro-plat/lib4go/logger"
)

type rpcServerAdapter struct {
}

func (h *rpcServerAdapter) Resolve(c string, conf conf.IServerConf, log *logger.Logger) (servers.IRegistryServer, error) {
	return NewMqcResponsiveServer(c, conf, log)
}

func init() {
	servers.Register("mqc", &rpcServerAdapter{})
}
