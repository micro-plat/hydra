package hydra

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/lib4go/types"
	"github.com/urfave/cli"
)

func (m *MicroApp) queryConfigAction(c *cli.Context) (err error) {
	if err := m.checkInput(); err != nil {
		m.xlogger.Error(err)
		cli.ShowCommandHelp(c, c.Command.Name)
		return nil
	}
	print := m.xlogger.Info

	m.logger.PauseLogging()
	defer m.logger.StartLogging()
	//创建注册中心
	rgst, err := registry.NewRegistryWithAddress(m.RegistryAddr, m.logger)
	if err != nil {
		m.xlogger.Error(err)
		return err
	}
	queryIndex := 0
	queryList := make(map[int][]byte)
	for i, tp := range m.ServerTypes {
		mainPath := registry.Join("/", m.PlatName, m.SystemName, tp, m.ClusterName, "conf")
		buffer, version, err := rgst.GetValue(mainPath)
		if err != nil {
			return err
		}
		sc, err := conf.NewServerConf(mainPath, buffer, version, rgst)
		if err != nil {
			return err
		}
		queryIndex++
		if i == 0 {
			print(getPrintNode(mainPath, queryIndex, 0))
		} else {
			print(getPrintNode(mainPath, queryIndex, 2))
		}
		queryList[queryIndex] = buffer

		sc.IterSubConf(func(k string, cn *conf.JSONConf) bool {
			queryIndex++
			print(getPrintNode(registry.Join(mainPath, k), queryIndex, -1))
			queryList[queryIndex] = cn.GetRaw()
			return true
		})
		if i == len(m.ServerTypes)-1 {
			index := -1
			sc.IterVarConf(func(k string, cn *conf.JSONConf) bool {
				queryIndex++
				if index == -1 {
					index++
					print(getPrintNode(registry.Join(m.PlatName, "var", k), queryIndex, 1))
				} else {
					print(getPrintNode(registry.Join(m.PlatName, "var", k), queryIndex, -1))
				}
				queryList[queryIndex] = cn.GetRaw()
				return true
			})
		}
	}
	for {
		fmt.Print("请输入数字序号 > ")
		var value string
		fmt.Scan(&value)
		if strings.ToUpper(value) == "Q" {
			return nil
		}
		nv := types.GetInt(value, -1)
		content, ok := queryList[nv]
		if !ok {
			continue
		}
		data := map[string]interface{}{}
		if err := json.Unmarshal(content, &data); err != nil {
			print(string(content))
			continue
		}
		buff, err := json.MarshalIndent(data, "", "    ")
		if err != nil {
			print(string(content))
			continue
		}
		print(string(buff))
	}
}
func getPrintNode(path string, index int, f int) string {
	p := strings.Trim(path, "/")
	ps := strings.Split(p, "/")
	buff := bytes.NewBufferString("")
	switch f {
	case -1:
		for c := 0; c < len(ps)-1; c++ {
			buff.WriteString("  ")
		}
		buff.WriteString("└─")
		buff.WriteString(fmt.Sprintf("[%d]", index))
		buff.WriteString(ps[len(ps)-1])
	case 0:
		for i, v := range ps {
			for c := 0; c < i; c++ {
				buff.WriteString("  ")
			}
			if i > 0 {
				buff.WriteString("└─")
			}
			if i == len(ps)-1 {
				buff.WriteString(fmt.Sprintf("[%d]", index))
			}
			buff.WriteString(v)
			buff.WriteString("\n")
		}
	default:
		for i := f; i < len(ps); i++ {
			for c := -1; c < i-1; c++ {
				buff.WriteString("  ")
			}
			buff.WriteString("└─")
			if i == len(ps)-1 {
				buff.WriteString(fmt.Sprintf("[%d]", index))
			}
			buff.WriteString(ps[i])
			buff.WriteString("\n")
		}
	}
	return strings.Trim(buff.String(), "\n")
}
