package conf

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/conf/app"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/hydra/hydra/cmds/pkgs"
	"github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/lib4go/types"
	"github.com/urfave/cli"
	"github.com/zkfy/log"
)

func showNow(c *cli.Context) (err error) {
	//1. 绑定应用程序参数
	global.Current().Log().Pause()
	if err := global.Def.Bind(c); err != nil {
		cli.ShowCommandHelp(c, c.Command.Name)
		return err
	}

	//2. 处理本地内存作为注册中心的服务发布问题
	if registry.GetProto(global.Current().GetRegistryAddr()) == registry.LocalMemory {
		if err := pkgs.Pub2Registry(true, ""); err != nil {
			return err
		}
	}

	//3. 显示配置
	return showConf(global.Current().GetRegistryAddr(),
		global.Current().GetPlatName(),
		global.Current().GetSysName(),
		global.Current().GetServerTypes(),
		global.Current().GetClusterName())
}

func showConf(addrs string, plat string, sysName string, types []string, cluster string) error {
	s := newShow(addrs, plat, sysName, types, cluster)
	return s.Show()

}

type show struct {
	subs    map[string]interface{}
	vars    map[string]interface{}
	nodes   [][]byte
	print   func(v ...interface{})
	plat    string
	sysName string
	types   []string
	cluster string
	addr    string
	rgst    registry.IRegistry
}

func newShow(addr string, plat string, sysName string, types []string, cluster string) *show {
	return &show{
		subs:    make(map[string]interface{}),
		vars:    make(map[string]interface{}),
		print:   log.New(os.Stdout, "", log.Llongcolor).Info,
		nodes:   make([][]byte, 0, 1),
		addr:    addr,
		plat:    plat,
		sysName: sysName,
		types:   types,
		cluster: cluster,
	}
}
func (s *show) Show() error {

	s.rgst = registry.GetCurrent()
	if err := s.printMainConf(); err != nil {
		return err
	}
	if err := s.printVarConf(); err != nil {
		return err
	}
	if extNode != "" {
		if err := s.printAnyConf(registry.Join(extNode)); err != nil {
			return err
		}
	}
	return s.readPrint()
}

func (s *show) readPrint() error {
	for {
		fmt.Print("请输入数字序号 > ")
		var value string
		fmt.Scan(&value)
		if strings.ToUpper(value) == "Q" {
			return nil
		}
		nv := types.GetInt(value, -1) - 1
		if nv > len(s.nodes)-1 || nv < 0 {
			s.print("输入的数字无效:", value)
			continue
		}
		content := s.nodes[nv]
		data := map[string]interface{}{}
		if err := json.Unmarshal(content, &data); err != nil {
			s.print(string(content))
			continue
		}
		buff, err := json.MarshalIndent(data, "", "    ")
		if err != nil {
			s.print(string(content))
			continue
		}
		s.print(string(buff))
	}
}

func (s *show) printMainConf() error {

	for _, tp := range s.types {
		sc, err := app.NewAPPConfBy(s.plat, s.sysName, tp, s.cluster, s.rgst)
		if err != nil {
			return err
		}
		s.getNodes(sc.GetServerConf().GetSubConfPath("main"), sc.GetServerConf().GetMainConf(), s.subs)
		sc.GetServerConf().Iter(func(path string, v *conf.RawConf) bool {
			npath := sc.GetServerConf().GetSubConfPath(path)
			s.getNodes(npath, v, s.subs)
			return true
		})
	}
	s.printNodes(s.subs, 0)
	return nil
}
func (s *show) printVarConf() error {
	if len(s.types) == 0 {
		return nil
	}
	sc, err := app.NewAPPConfBy(s.plat, s.sysName, s.types[0], s.cluster, s.rgst)
	if err != nil {
		return err
	}
	sc.GetVarConf().Iter(func(path string, v *conf.RawConf) bool {
		npath := sc.GetVarConf().GetVarPath(path)
		s.getNodes(npath, v, s.vars)
		return true
	})

	s.printNodes(s.vars, 0)
	return nil
}
func (s *show) printAnyConf(root string) error {

	nodes := make(map[string]interface{})
	getNode(root, root, nodes)

	s.printNodes(map[string]interface{}{root: nodes}, 0)
	return nil
}
func getNode(path string, p string, node map[string]interface{}) {
	r := registry.GetCurrent()
	paths, _, _ := r.GetChildren(path)
	if len(paths) > 0 {
		for _, p := range paths {
			nnode := make(map[string]interface{})
			getNode(registry.Join(path, p), p, nnode)
			node[p] = nnode
		}
	} else {
		v, _, _ := r.GetValue(path)
		node[p] = v
	}
}
func (s *show) getNodes(path string, v *conf.RawConf, input map[string]interface{}) {
	li := strings.SplitN(strings.Trim(path, "/"), "/", 2)
	if len(li) == 1 {
		input[li[0]] = v.GetRaw()
		return
	}
	if len(li) > 1 {
		if np, ok := input[li[0]]; !ok {
			nmap := make(map[string]interface{})
			input[li[0]] = nmap
			s.getNodes(li[1], v, nmap)
		} else {
			switch c := np.(type) {
			case map[string]interface{}:
				s.getNodes(li[1], v, c)

			}

		}
	}
}
func (s *show) printNodes(nodes map[string]interface{}, index int) {
	for k, v := range nodes {

		switch c := v.(type) {
		case map[string]interface{}:
			s.print(fmt.Sprintf("%s└─%s", strings.Repeat("  ", index), k))
			s.printNodes(c, index+1)
		case []byte:
			s.nodes = append(s.nodes, c)
			s.print(fmt.Sprintf("%s└─%s[%d]", strings.Repeat("  ", index), k, len(s.nodes)))
		}
	}
}
