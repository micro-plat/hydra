package creator

import (
	"fmt"
	"strings"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/engines"
	"github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/lib4go/logger"
)

//Creator 配置文件创建器
type Creator struct {
	registry     registry.IRegistry
	registryAddr string
	logger       logger.ILogging
	showTitle    bool
	binder       IBinder
	platName     string
	systemName   string
	serverTypes  []string
	clusterName  string
	customer     func() error
}

//NewCreator 配置文件创建器
func NewCreator(platName string, systemName string, serverTypes []string, clusterName string, binder IBinder, registryAddr string, rgst registry.IRegistry, logger logger.ILogging) (w *Creator) {
	w = &Creator{
		platName:     platName,
		systemName:   systemName,
		serverTypes:  serverTypes,
		clusterName:  clusterName,
		registry:     rgst,
		registryAddr: registryAddr,
		logger:       logger,
		binder:       binder,
	}
	return
}
func (c *Creator) installParams() error {
	//检查必须输入参数
	input := c.binder.GetInput()
	if len(input) > 0 {
		if !c.binder.Confirm("设置基础参数值(这些参数用于创建配置数据)?") {
			return nil
		}
	}
	for k := range input {
		if strings.HasPrefix(k, "#") {
			continue
		}
		nvalue, err := getInputValue(k, input, "")
		if err != nil {
			return err
		}
		c.binder.SetParam(k, nvalue)
	}
	return nil
}
func (c *Creator) installRegistry() error {
	//检查配置模式
	mode, cn := c.checkRegistry()
	//创建主配置
	if !cn {
		return nil
	}
	for _, tp := range c.serverTypes {
		mainPath := registry.Join("/", c.platName, c.systemName, tp, c.clusterName, "conf")
		rpath := c.getRealMainPath(mainPath)
		ok, err := c.registry.Exists(rpath)
		if err != nil {
			return err
		}
		if ok && mode == modeAuto {
			continue
		}
		pc, _, _ := c.registry.GetChildren(rpath)
		for _, v := range pc {
			c.registry.Delete(registry.Join(rpath, v))
		}
		err = c.registry.Delete(rpath)
		if mode == modeNew {
			c.logger.Info("\t\t删除配置:", rpath)
		}
		if err := c.binder.ScanMainConf(mainPath, tp); err != nil {
			return err
		}
		content := c.binder.GetMainConf(tp)
		if ok && mode == modeCover {
			if err := c.createMainConf(mainPath, content); err != nil {
				return err
			}
			c.logger.Info("\t\t修改配置:", mainPath)
		}
		if err := c.createMainConf(mainPath, content); err != nil {
			return err
		}
		c.logger.Info("\t\t创建配置:", mainPath)
	}
	//检查子配置
	for _, tp := range c.serverTypes {
		mainPath := registry.Join("/", c.platName, c.systemName, tp, c.clusterName, "conf")
		subNames := c.binder.GetSubConfNames(tp)
		for _, subName := range subNames {
			ok, err := c.registry.Exists(registry.Join(mainPath, subName))
			if err != nil {
				return err
			}
			if ok && mode == modeAuto {
				continue
			}
			//删除配置重建
			c.registry.Delete(registry.Join(mainPath, subName))
			if err := c.binder.ScanSubConf(mainPath, tp, subName); err != nil {
				return err
			}

			path := registry.Join("/", mainPath, subName)
			content := c.binder.GetSubConf(tp, subName)
			if err := c.createConf(path, content); err != nil {
				return err
			}
			c.logger.Info("\t\t创建配置:", path)
		}
	}

	//检查平台配置
	varNames := c.binder.GetVarConfNames()
	for _, varName := range varNames {
		ok, err := c.registry.Exists(registry.Join("/", c.platName, "var", varName))
		if err != nil {
			return err
		}
		if ok && mode == modeAuto {
			continue
		}
		//删除配置重建
		c.registry.Delete(registry.Join("/", c.platName, "var", varName))
		if err := c.binder.ScanVarConf(c.platName, varName); err != nil {
			return err
		}
		path := registry.Join("/", c.platName, "var", varName)
		content := c.binder.GetVarConf(varName)
		if err := c.createConf(path, content); err != nil {
			return err
		}
		c.logger.Info("\t\t创建配置:", path)
	}
	return nil
}

//Start 扫描并绑定所有参数
func (c *Creator) Start() (err error) {
	if err = c.installParams(); err != nil {
		return err
	}
	if err = c.installRegistry(); err != nil {
		return err
	}
	//执行用户自定义安装
	if err = c.customerInstall(); err != nil {
		return fmt.Errorf("安装程序执行失败:%v", err)
	}
	return nil
}

func (c *Creator) customerInstall() error {
	for _, tp := range c.serverTypes {
		installs := c.binder.GetInstallers(tp)
		if installs == nil || len(installs) == 0 {
			continue
		}
		mainPath := registry.Join("/", c.platName, c.systemName, tp, c.clusterName, "conf")
		buffer, version, err := c.registry.GetValue(mainPath)
		if err != nil {
			return err
		}
		conf, err := conf.NewServerConf(mainPath, buffer, version, c.registry)
		if err != nil {
			return err
		}
		engine, err := engines.NewServiceEngine(conf, c.registryAddr, c.logger)
		if err != nil {
			return err
		}
		defer engine.Close()
		for _, install := range installs {
			if err := install(engine); err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *Creator) createConf(path string, data string) error {
	if data == "" {
		return nil
	}
	return c.registry.CreatePersistentNode(path, data)
}

func (c *Creator) getRealMainPath(path string) string {
	extPath := ""
	if !c.registry.CanWirteDataInDir() {
		extPath = ".init"
	}
	return registry.Join(path, extPath)
}
func (c *Creator) createMainConf(path string, data string) error {
	if data == "" {
		data = "{}"
	}
	rpath := c.getRealMainPath(path)
	return c.registry.CreatePersistentNode(rpath, data)
}
func (c *Creator) updateMainConf(path string, data string) error {
	if data == "" {
		data = "{}"
	}
	rpath := c.getRealMainPath(path)
	_, v, err := c.registry.GetValue(rpath)
	if err != nil {
		return err
	}
	return c.registry.Update(rpath, data, v)
}

func (c *Creator) checkRegistry() (mode int, cn bool) {
	msg := "创建注册中心配置数据?,如存在则不修改(1),如果存在则覆盖(2),删除所有配置并重建(3),退出(n|no):"

	var value string
	fmt.Print("\t\033[;33m-> " + msg + "\033[0m")
	fmt.Scan(&value)
	nvalue := strings.ToUpper(value)
	switch nvalue {
	case "1":
		return modeAuto, true
	case "2":
		return modeCover, true
	case "3":
		return modeNew, true
	}
	return 0, false
}
