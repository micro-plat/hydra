package hydra

import (
	"fmt"
	"os"
	"reflect"

	"github.com/micro-plat/hydra/conf/creator"

	"github.com/asaskevich/govalidator"
	"github.com/micro-plat/hydra/component"
	_ "github.com/micro-plat/hydra/hydra/impt"
	"github.com/micro-plat/hydra/hydra/rqs"
	"github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/lib4go/logger"
	"github.com/urfave/cli"
)

//MicroApp  微服务应用
type MicroApp struct {
	app    *cli.App
	logger *logger.Logger
	Conf   *creator.Binder
	hydra  *Hydra
	*option
	remoteQueryService *rqs.RemoteQueryService
	registry           registry.IRegistry
	component.IComponentRegistry
}

//NewApp 创建微服务应用
func NewApp(opts ...Option) (m *MicroApp) {
	m = &MicroApp{option: &option{}, IComponentRegistry: component.NewServiceRegistry()}
	m.Conf = creator.NewBinder()
	for _, opt := range opts {
		opt(m.option)
	}
	m.logger = logger.GetSession("hydra", logger.CreateSession())
	return m
}

//Start 启动服务器
func (m *MicroApp) Start() {
	defer logger.Close()
	m.app = m.getCliApp()
	if err := m.app.Run(os.Args); err != nil {
		return
	}
}

//Use 注册所有服务
func (m *MicroApp) Use(r func(r component.IServiceRegistry)) {
	r(m.IComponentRegistry)
}
func (m *MicroApp) action(c *cli.Context) (err error) {
	if m.remoteLogger {
		m.RemoteLogger = m.remoteLogger
	}
	if err := m.checkInput(); err != nil {
		cli.ErrWriter.Write([]byte("  " + err.Error() + "\n\n"))
		cli.ShowCommandHelp(c, c.Command.Name)
		return nil
	}
	if m.registry, err = registry.NewRegistryWithAddress(m.RegistryAddr, m.logger); err != nil {
		m.logger.Error(err)
		return err
	}

	if m.RemoteQueryService {
		m.remoteQueryService, err = rqs.NewHRemoteQueryService(m.PlatName, m.SystemName, m.ServerTypes, m.ClusterName, m.registry, VERSION)
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

	m.hydra = NewHydra(m.PlatName, m.SystemName, m.ServerTypes, m.ClusterName, m.Trace,
		m.RegistryAddr, m.Conf, m.IsDebug, m.RemoteLogger, m.IComponentRegistry)
	if err := m.hydra.Start(); err != nil {
		m.logger.Error(err)
		return err
	}
	return nil
}

func (m *MicroApp) checkInput() (err error) {
	if m.ServerTypeNames != "" && len(m.ServerTypes) == 0 {
		WithServerTypes(m.ServerTypeNames)(m.option)
	}
	if m.PlatName == "" && m.Name != "" {
		WithName(m.Name)(m.option)
	}

	if b, err := govalidator.ValidateStruct(m.option); !b {
		err = fmt.Errorf("validate(%v) %v", reflect.TypeOf(m.option), err)
		return err
	}
	return
}
func (m *MicroApp) Shutdown() {
	m.hydra.Shutdown()
}
