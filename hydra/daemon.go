// Example of a daemon with echo service
package hydra

import (
	"os"

	"github.com/micro-plat/hydra/conf/creator"
	"github.com/micro-plat/hydra/registry"
	"github.com/urfave/cli"
	"github.com/zkfy/log"
)

var xlogger = log.New(os.Stdout, "", log.Llongcolor)

func init() {
	xlogger.SetOutputLevel(log.Ldebug)
}
func (m *MicroApp) startAction(c *cli.Context) (err error) {
	msg, err := m.service.Start()
	if err != nil {
		xlogger.Error(err)
		return err
	}
	xlogger.Info(msg)
	return nil
}
func (m *MicroApp) stopAction(c *cli.Context) (err error) {
	msg, err := m.service.Stop()
	if err != nil {
		xlogger.Error(err)
		return err
	}
	xlogger.Info(msg)
	return nil
}
func (m *MicroApp) installAction(c *cli.Context) (err error) {

	if err = m.checkInput(); err != nil {
		cli.ErrWriter.Write([]byte("  " + err.Error() + "\n\n"))
		cli.ShowCommandHelp(c, c.Command.Name)
		return err
	}
	if err = m.install(); err != nil {
		xlogger.Error(err)
		return err
	}

	//安装配置文件
	msg, err := m.service.Install(os.Args[2:]...)
	if err != nil {
		xlogger.Error(err)
		return err
	}
	xlogger.Info(msg)
	return nil
}
func (m *MicroApp) removeAction(c *cli.Context) (err error) {
	msg, err := m.service.Remove()
	if err != nil {
		xlogger.Error(err)
		return err
	}
	xlogger.Info(msg)
	return nil
}
func (m *MicroApp) statusAction(c *cli.Context) (err error) {
	msg, err := m.service.Status()
	if err != nil {
		xlogger.Error(err)
		return err
	}
	xlogger.Info(msg)
	return nil
}

func (m *MicroApp) install() (err error) {
	m.logger.PauseLogging()
	defer m.logger.StartLogging()
	//创建注册中心
	rgst, err := registry.NewRegistryWithAddress(m.RegistryAddr, m.logger)
	if err != nil {
		xlogger.Error(err)
		return err
	}

	//自动创建配置
	creator := creator.NewCreator(m.PlatName, m.SystemName, m.ServerTypes, m.ClusterName, m.Conf, m.RegistryAddr, rgst, xlogger)
	err = creator.Start()
	if err != nil {
		xlogger.Error(err)
		return err
	}
	return nil

}
