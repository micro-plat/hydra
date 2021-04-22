package conf

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	logs "github.com/lib4dev/cli/logger"
	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/conf/app"
	"github.com/micro-plat/hydra/conf/pkgs/security"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/hydra/hydra/cmds/pkgs"
	"github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/lib4go/types"
	"github.com/urfave/cli"
)

func exportNow(c *cli.Context) (err error) {
	//1. 绑定应用程序参数
	global.Current().Log().Pause()
	if err := global.Def.Bind(c); err != nil {
		cli.ShowCommandHelp(c, c.Command.Name)
		return err
	}

	//2. 处理本地内存作为注册中心的服务发布问题
	if registry.GetProto(global.Current().GetRegistryAddr()) == registry.LocalMemory {
		if err := pkgs.Pub2Registry(true, nil); err != nil {
			return err
		}
	}

	//3. 导出配置
	return exportConf(
		global.Current().GetPlatName(),
		global.Current().GetSysName(),
		global.Current().GetServerTypes(),
		global.Current().GetClusterName(),
	)
}

func exportConf(plat string, sysName string, types []string, cluster string) error {
	s := newExport(plat, sysName, types, cluster)
	return s.Export()

}

type export struct {
	confs   map[string]interface{}
	plat    string
	sysName string
	types   []string
	cluster string
	rgst    registry.IRegistry
	cover   bool
	encrypt bool
}

func newExport(plat string, sysName string, types []string, cluster string) *export {
	return &export{
		confs:   make(map[string]interface{}),
		plat:    plat,
		sysName: sysName,
		types:   types,
		cluster: cluster,
		cover:   coverConfIfExists,
		encrypt: confEncrypt,
	}
}
func (s *export) Export() error {
	s.rgst = registry.GetCurrent()
	if err := s.getMainConf(); err != nil {
		return err
	}
	if err := s.getVarConf(); err != nil {
		return err
	}
	root := types.GetString(confExportPath, "./")
	if err := s.writeConf(filepath.Join(root, "conf.json")); err != nil {
		return err
	}
	return nil
}

func (s *export) getMainConf() error {
	for _, tp := range s.types {
		sc, err := app.NewAPPConfBy(s.plat, s.sysName, tp, s.cluster, s.rgst)
		if err != nil {
			return err
		}
		s.getNodes(sc.GetServerConf().GetServerPath(), sc.GetServerConf().GetMainConf(), s.confs)
		sc.GetServerConf().Iter(func(path string, v *conf.RawConf) bool {
			npath := sc.GetServerConf().GetSubConfPath(path)
			s.getNodes(npath, v, s.confs)
			return true
		})
	}
	return nil
}

func (s *export) getVarConf() error {
	if len(s.types) == 0 {
		return nil
	}
	sc, err := app.NewAPPConfBy(s.plat, s.sysName, s.types[0], s.cluster, s.rgst)
	if err != nil {
		return err
	}
	sc.GetVarConf().Iter(func(path string, v *conf.RawConf) bool {
		npath := sc.GetVarConf().GetVarPath(path)
		s.getNodes(npath, v, s.confs)
		return true
	})
	return nil
}

func (s *export) getNodes(path string, v *conf.RawConf, input map[string]interface{}) {
	fmt.Println("path:", path)
	if s.encrypt {
		input[path] = security.Encrypt(v.GetRaw())
		return
	}
	t := make(map[string]interface{})
	json.Unmarshal(v.GetRaw(), &t)
	input[path] = t
}

func (s *export) writeConf(path string) error {
	content, _ := json.Marshal(s.confs)
	//生成文件
	fs, err := create(path, s.cover)
	if err != nil {
		return err
	}
	defer fs.Close()
	logs.Log.Info("生成文件:", path)
	fs.WriteString(string(content))
	return nil
}

//pathExists 文件是否存在
func pathExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

//create 创建文件，文件夹 存在时写入则覆盖
func create(path string, cover bool) (file *os.File, err error) {
	dir := filepath.Dir(path)
	if !pathExists(dir) {
		if err = os.MkdirAll(dir, os.ModePerm); err != nil {
			return nil, fmt.Errorf("创建文件夹%s失败:%v", path, err)
		}
	}

	var srcf *os.File
	if !pathExists(path) {
		srcf, err = os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.ModePerm)
		if err != nil {
			return nil, fmt.Errorf("无法打开文件:%s(err:%v)", path, err)
		}
		return srcf, nil

	}
	if !cover {
		return nil, fmt.Errorf("文件:%s已经存在", path)
	}
	srcf, err = os.OpenFile(path, os.O_TRUNC|os.O_RDWR, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("无法打开文件:%s(err:%v)", path, err)
	}
	return srcf, nil

}
