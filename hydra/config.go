package hydra

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/registry"
	"github.com/urfave/cli"
)

func (m *MicroApp) queryConfigAction(c *cli.Context) (err error) {
	if err := m.checkInput(); err != nil {
		m.xlogger.Error(err)
		cli.ShowCommandHelp(c, c.Command.Name)
		return nil
	}

	m.logger.PauseLogging()
	defer m.logger.StartLogging()
	//创建注册中心
	rgst, err := registry.NewRegistryWithAddress(m.RegistryAddr, m.logger)
	if err != nil {
		m.xlogger.Error(err)
		return err
	}
	for _, tp := range m.ServerTypes {
		mainPath := filepath.Join("/", m.PlatName, m.SystemName, tp, m.ClusterName, "conf")
		buffer, version, err := rgst.GetValue(mainPath)
		if err != nil {
			return err
		}
		sc, err := conf.NewServerConf(mainPath, buffer, version, rgst)
		if err != nil {
			return err
		}
		fmt.Println(getPrintNode(mainPath, true))
		sc.IterSubConf(func(k string, conf *conf.JSONConf) bool {
			fmt.Println(getPrintNode(k, false))
			return true
		})
		fmt.Println(getPrintNode(mainPath, true))
		sc.IterSubConf(func(k string, conf *conf.JSONConf) bool {
			fmt.Println(getPrintNode(k, false))
			return true
		})
		sc.IterVarConf(func(k string, conf *conf.JSONConf) bool {
			fmt.Println(getPrintNode(k, false))
			return true
		})
	}

	return nil
}
func getPrintNode(path string, f bool) string {
	p := strings.Trim(path, "/")
	ps := strings.Split(p, "/")
	buff := bytes.NewBufferString("")
	if !f {
		for c := 0; c < len(ps); c++ {
			buff.WriteString("--")
		}
		buff.WriteString(ps[len(ps)-1])
		return buff.String()
	}
	for i, v := range ps {
		for c := 0; c < i; c++ {
			buff.WriteString("--")
		}
		buff.WriteString(v)
	}
	return buff.String()
}
