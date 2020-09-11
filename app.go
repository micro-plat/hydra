package hydra

import (
	"syscall"

	"github.com/micro-plat/cli"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/hydra/services"
	"github.com/micro-plat/lib4go/logger"

	_ "github.com/micro-plat/hydra/registry/watcher/wchild"
	_ "github.com/micro-plat/hydra/registry/watcher/wvalue"

	_ "github.com/micro-plat/hydra/hydra/cmds/conf"
	_ "github.com/micro-plat/hydra/hydra/cmds/install"
	_ "github.com/micro-plat/hydra/hydra/cmds/remove"
	_ "github.com/micro-plat/hydra/hydra/cmds/run"
	_ "github.com/micro-plat/hydra/hydra/cmds/update"

	_ "github.com/micro-plat/hydra/hydra/cmds/start"
	_ "github.com/micro-plat/hydra/hydra/cmds/status"
	_ "github.com/micro-plat/hydra/hydra/cmds/stop"

	_ "github.com/micro-plat/hydra/hydra/cmds/backup"
	_ "github.com/micro-plat/hydra/hydra/cmds/rollback"

	_ "github.com/micro-plat/hydra/hydra/cmds/pub"

	_ "github.com/micro-plat/hydra/registry/registry/filesystem"
	_ "github.com/micro-plat/hydra/registry/registry/localmemory"
	_ "github.com/micro-plat/hydra/registry/registry/zookeeper"
	_ "github.com/micro-plat/hydra/registry/registry/etcd"

)

//MicroApp  微服务应用
type MicroApp struct {
	app *cli.App
	services.IService
	Cli *ucli
}

//NewApp 创建微服务应用
func NewApp(opts ...Option) (m *MicroApp) {
	m = &MicroApp{
		IService: services.Def,
		Cli:      newUCli(),
	}
	for _, opt := range opts {
		opt()
	}
	return m
}

//Start 启动服务器
func (m *MicroApp) Start() {
	defer logger.Close()
	m.app = cli.New(cli.WithVersion(global.Version), cli.WithUsage(global.Usage))
	m.app.Start()
}

//Close 关闭服务器
func (m *MicroApp) Close() {
	parent := syscall.Getppid()
	syscall.Kill(parent, syscall.SIGUSR2)
}
