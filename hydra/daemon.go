// Example of a daemon with echo service
package hydra

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
)

func (m *MicroApp) startAction(c *cli.Context) (err error) {
	msg, err := m.service.Start()
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(msg)
	return nil
}
func (m *MicroApp) stopAction(c *cli.Context) (err error) {
	msg, err := m.service.Stop()
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(msg)
	return nil
}
func (m *MicroApp) installAction(c *cli.Context) (err error) {

	if err = m.checkInput(); err != nil {
		cli.ErrWriter.Write([]byte("  " + err.Error() + "\n\n"))
		cli.ShowCommandHelp(c, c.Command.Name)
		return err
	}
	//安装配置服务
	p := func(v ...interface{}) {
		fmt.Println(v...)
	}
	if err = m.install(p); err != nil {
		fmt.Println(err)
		return err
	}

	//安装配置文件
	msg, err := m.service.Install(os.Args[2:]...)
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println(msg)
	return nil
}
func (m *MicroApp) removeAction(c *cli.Context) (err error) {
	msg, err := m.service.Remove()
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(msg)
	return nil
}
func (m *MicroApp) statusAction(c *cli.Context) (err error) {
	msg, err := m.service.Status()
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(msg)
	return nil
}
