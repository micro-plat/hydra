package creator

import (
	"fmt"
	"path/filepath"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/engines"
	"github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/lib4go/logger"
)

//Creator 配置文件创建器
type Creator struct {
	registry     registry.IRegistry
	registryAddr string
	logger       *logger.Logger
	showTitle    bool
	binder       IBinder
	platName     string
	systemName   string
	serverTypes  []string
	clusterName  string
	customer     func() error
}

//NewCreator 配置文件创建器
func NewCreator(platName string, systemName string, serverTypes []string, clusterName string, binder IBinder, registryAddr string, rgst registry.IRegistry, logger *logger.Logger) (w *Creator) {
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

//Start 扫描并绑定所有参数
func (c *Creator) Start() (err error) {
	input := c.binder.GetInput()
	for k, v := range input {
		fmt.Printf("请输入%s:", v)
		var value string
		fmt.Scan(&value)
		//c.subConfParamsForTranslate[nodeName][p] = value
		c.binder.SetParam(k, value)
	}

	for _, tp := range c.serverTypes {
		mainPath := filepath.Join("/", c.platName, c.systemName, tp, c.clusterName, "conf")
		//检查主配置
		ok, err := c.registry.Exists(c.getRealMainPath(mainPath))
		if err != nil {
			return err
		}
		if ok {
			continue
		}
		if c.binder.GetMainConfScanNum(tp) > 0 {
			if !c.checkContinue() {
				return nil
			}
		}
		if err := c.binder.ScanMainConf(mainPath, tp); err != nil {
			return err
		}

		content := c.binder.GetMainConf(tp)
		if err := c.createMainConf(mainPath, content); err != nil {
			return err
		}
		c.logger.Print("创建配置:", mainPath)
	}
	//检查子配置
	for _, tp := range c.serverTypes {
		mainPath := filepath.Join("/", c.platName, c.systemName, tp, c.clusterName, "conf")
		subNames := c.binder.GetSubConfNames(tp)
		for _, subName := range subNames {
			ok, err := c.registry.Exists(filepath.Join(mainPath, subName))
			if err != nil {
				return err
			}
			if ok {
				continue
			}
			if c.binder.GetSubConfScanNum(tp, subName) > 0 {
				if !c.checkContinue() {
					return nil
				}
			}
			if err := c.binder.ScanSubConf(mainPath, tp, subName); err != nil {
				return err
			}

			path := filepath.Join("/", mainPath, subName)
			content := c.binder.GetSubConf(tp, subName)
			if err := c.createConf(path, content); err != nil {
				return err
			}
			c.logger.Print("创建配置:", path)
		}
	}

	//检查平台配置
	varNames := c.binder.GetVarConfNames()
	for _, varName := range varNames {
		ok, err := c.registry.Exists(filepath.Join("/", c.platName, "var", varName))
		if err != nil {
			return err
		}
		if ok {
			continue
		}
		if c.binder.GetVarConfScanNum(varName) > 0 {
			if !c.checkContinue() {
				return nil
			}
		}
		if err := c.binder.ScanVarConf(c.platName, varName); err != nil {
			return err
		}
		path := filepath.Join("/", c.platName, "var", varName)
		content := c.binder.GetVarConf(varName)
		if err := c.createConf(path, content); err != nil {
			return err
		}
		c.logger.Print("创建配置:", path)
	}
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
		mainPath := filepath.Join("/", c.platName, c.systemName, tp, c.clusterName, "conf")
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
	return filepath.Join(path, extPath)
}
func (c *Creator) createMainConf(path string, data string) error {
	if data == "" {
		data = "{}"
	}
	rpath := c.getRealMainPath(path)
	return c.registry.CreatePersistentNode(rpath, data)
}
func (c *Creator) checkContinue() bool {
	if !c.showTitle {
		c.showTitle = true
	} else {
		return true
	}
	var index string
	fmt.Print("当前服务有一些参数未配置，立即配置(y|N):")
	fmt.Scan(&index)
	if index != "y" && index != "Y" && index != "yes" && index != "YES" {
		return false
	}
	return true
}
