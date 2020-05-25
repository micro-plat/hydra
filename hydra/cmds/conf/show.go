package conf

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/conf/server"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/hydra/hydra/cmds/pkgs"
	"github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/lib4go/types"
	"github.com/zkfy/log"
)

func show() error {

	//1. 处理日志
	print := log.New(os.Stdout, "", log.Llongcolor).Info

	//2. 创建注册中心
	rgst, err := registry.NewRegistry(global.Current().GetRegistryAddr(), global.Current().Log())
	if err != nil {
		return err
	}

	//3. 处理本地内存作为注册中心的服务发布问题
	if registry.GetProto(global.Current().GetRegistryAddr()) == registry.LocalMemory {
		if err := pkgs.Pub2Registry(true); err != nil {
			return err
		}
	}

	queryIndex := 0
	queryList := make(map[int][]byte)
	for i, tp := range global.Current().GetServerTypes() {
		sc, err := server.NewServerConfBy(global.Current().GetPlatName(), global.Current().GetSysName(), tp, global.Current().GetClusterName(), rgst)
		if err != nil {
			return err
		}
		queryIndex++
		if i == 0 {
			print(getPrintNode(sc.GetMainConf().GetMainPath(), queryIndex, 0))
		} else {
			print(getPrintNode(sc.GetMainConf().GetMainPath(), queryIndex, 2))
		}
		queryList[queryIndex] = sc.GetMainConf().GetMainConf().GetRaw()

		sc.GetMainConf().Iter(func(k string, cn *conf.JSONConf) bool {
			queryIndex++
			print(getPrintNode(registry.Join(sc.GetMainConf().GetMainPath(), k), queryIndex, -1))
			queryList[queryIndex] = cn.GetRaw()
			return true
		})
		if i == len(global.Current().GetServerTypes())-1 {
			index := -1
			sc.GetVarConf().Iter(func(k string, cn *conf.JSONConf) bool {
				queryIndex++
				if index == -1 {
					index++
					print(getPrintNode(registry.Join(sc.GetMainConf().GetPlatName(), "var", k), queryIndex, 1))
				} else {
					print(getPrintNode(registry.Join(sc.GetMainConf().GetPlatName(), "var", k), queryIndex, -1))
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
