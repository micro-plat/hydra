package ws

import (
	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/servers"
	"github.com/micro-plat/lib4go/logger"
)

type wsServerAdapter struct {
}

func (h *wsServerAdapter) Resolve(registryAddr string, conf conf.IServerConf, log *logger.Logger) (servers.IRegistryServer, error) {
	return NewWSServerResponsiveServer(registryAddr, conf, log)
}

func init() {
	servers.Register("ws", &wsServerAdapter{})
}
