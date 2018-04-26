package hydra

import (
	"errors"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/hydra/servers"
	"github.com/micro-plat/lib4go/logger"

	"time"
)

var (
	errServerIsExist    = errors.New("服务已存在")
	errServerIsNotExist = errors.New("服务不存在")
)

//Server hydra的单个服务器示例
type server struct {
	registry     registry.IRegistry
	cnf          conf.IServerConf
	registryAddr string

	startTime time.Time
	logger    *logger.Logger
	server    servers.IRegistryServer
}

//newServer 初始化服务器
func newServer(cnf conf.IServerConf, registryAddr string, registry registry.IRegistry) *server {
	return &server{
		cnf:          cnf,
		registryAddr: registryAddr,
		registry:     registry,
		logger:       logger.New(cnf.GetServerName()),
	}
}

//Start 启用服务器
func (h *server) Start() (err error) {
	//构建服务器
	h.server, err = servers.NewRegistryServer(h.cnf.GetServerType(), h.registryAddr, h.cnf, h.logger)
	if err != nil {
		return err
	}
	err = h.server.Start()
	if err != nil {
		return err
	}
	h.startTime = time.Now()
	return nil
}

//Notify 配置发生变化通知服务器变更
func (h *server) Notify(cnf conf.IServerConf) error {
	return h.server.Notify(cnf)
}

//GetStatus 获取当前服务状态
func (h *server) GetStatus() string {
	return h.server.GetStatus()
}

//GetAddress 获取服务器地址
func (h *server) GetAddress() string {
	return h.server.GetAddress()
}

//GetServices 获取服务列表
func (h *server) GetServices() []string {
	return h.server.GetServices()
}

//Restarted 服务器是否已重启
func (h *server) Restarted() bool {
	return h.server.Restarted()
}

//Shutdown 关闭服务器
func (h *server) Shutdown() {
	if h.server != nil {
		h.server.Shutdown()
	}
}
