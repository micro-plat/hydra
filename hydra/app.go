package hydra

import (
	"fmt"
	"os"

	"github.com/asaskevich/govalidator"
	"github.com/micro-plat/hydra/component"
	"github.com/micro-plat/hydra/conf/creator"
	"github.com/micro-plat/hydra/hydra/daemon"
	_ "github.com/micro-plat/hydra/hydra/impt"
	"github.com/micro-plat/hydra/hydra/rqs"
	"github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/lib4go/logger"
	"github.com/urfave/cli"
	"github.com/zkfy/log"
)

//MicroApp  微服务应用
type MicroApp struct {
	app     *cli.App
	logger  *logger.Logger
	xlogger logger.ILogging
	Conf    *creator.Binder
	hydra   *Hydra
	*option
	remoteQueryService *rqs.RemoteQueryService
	component.IComponentRegistry
	service daemon.Daemon
	funs    *funs
	Cli     ICli
}

//NewApp 创建微服务应用
func NewApp(opts ...Option) (m *MicroApp) {
	m = &MicroApp{
		option:             &option{},
		funs:               newFuns(),
		IComponentRegistry: component.NewServiceRegistry(),
	}
	logging := log.New(os.Stdout, "", log.Llongcolor)
	logging.SetOutputLevel(log.Ldebug)
	m.xlogger = logging
	m.Cli = NewCli()
	m.Conf = creator.NewBinder(logging)
	for _, opt := range opts {
		opt(m.option)
	}
	m.logger = logger.GetSession("hydra", logger.CreateSession())
	return m
}

//Start 启动服务器
func (m *MicroApp) Start() {
	var err error
	defer logger.Close()
	m.app = m.getCliApp()
	m.service, err = daemon.New(m.app.Name, m.app.Name)
	if err != nil {
		m.logger.Error(err)
		return
	}
	if err := m.app.Run(os.Args); err != nil {
		return
	}

}

//Use 注册所有服务
func (m *MicroApp) Use(r func(r component.IServiceRegistry)) {
	r(m.IComponentRegistry)
}

func (m *MicroApp) action(c *cli.Context) (err error) {
	if err := m.checkInput(c); err != nil {
		m.xlogger.Warn(err)
		cli.ShowCommandHelp(c, c.Command.Name)
		return nil
	}

	//初始化远程日志
	if m.remoteLogger {
		m.RemoteLogger = m.remoteLogger
	}

	//启动服务查询
	if m.RemoteQueryService {
		//创建注册中心
		rgst, err := registry.NewRegistryWithAddress(m.RegistryAddr, m.logger)
		if err != nil {
			m.logger.Error(err)
			return err
		}

		m.remoteQueryService, err = rqs.NewHRemoteQueryService(m.PlatName, m.SystemName, m.ServerTypes, m.ClusterName, rgst, VERSION)
		if err != nil {
			m.logger.Error(err)
			return err
		}
		if err = m.remoteQueryService.Start(); err != nil {
			m.logger.Error(err)
			return err
		}
		m.remoteQueryService.HydraShutdown = m.Shutdown
		defer m.remoteQueryService.Shutdown()
	}

	m.hydra = NewHydra(m.app.Name, m.PlatName, m.SystemName, m.ServerTypes, m.ClusterName, m.Trace,
		m.RegistryAddr, m.IsDebug, m.RemoteLogger, m.logger, m.IComponentRegistry)

	m.run()
	return nil
}
func (m *MicroApp) run() {
	p := &once{app: m}
	p.run()
}
func (m *MicroApp) checkInput(c *cli.Context) (err error) {
	m.Cli.setContext(c)
	if m.ServerTypeNames != "" && len(m.ServerTypes) == 0 {
		WithServerTypes(m.ServerTypeNames)(m.option)
	}
	if m.PlatName == "" && m.Name != "" {
		WithName(m.Name)(m.option)
	}
	if m.IsDebug {
		m.PlatName += "_debug"
	}
	if b, err := govalidator.ValidateStruct(m.option); !b {
		err = fmt.Errorf("服务器运行缺少参数，请查看以下帮助信息")
		return err
	}

	//获取外部验证
	vds := m.Cli.getValidators(c.Command.Name)
	for _, validator := range vds {
		if err := validator(c); err != nil {
			return err
		}
	}
	return m.funs.Call()
}

//Shutdown 关闭服务
func (m *MicroApp) Shutdown() {
	if m.hydra != nil {
		m.hydra.Shutdown()
	}
}
