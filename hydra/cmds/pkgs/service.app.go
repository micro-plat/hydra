package pkgs

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/hydra/hydra/cmds/pkgs/service"
	"github.com/micro-plat/hydra/hydra/servers"
	"github.com/urfave/cli"
)

type HydraService struct {
	service.Service
	ServiceName string
	DisplayName string
	Description string
	Arguments   []string
}

//GetService GetService
func GetService(c *cli.Context, isFixed bool, args ...string) (hydraSrv *HydraService, err error) {
	//1. 构建服务配置
	cfg := GetSrvConfig(isFixed, args...)

	//2.创建本地服务
	appSrv, err := service.New(GetSrvApp(c), cfg)
	if err != nil {
		return nil, err
	}
	return &HydraService{
		Service:     appSrv,
		ServiceName: cfg.Name,
		DisplayName: cfg.DisplayName,
		Description: cfg.Description,
		Arguments:   cfg.Arguments,
	}, err
}

//GetSrvConfig SrvCfg
func GetSrvConfig(isFixed bool, args ...string) *service.Config {
	svcName := global.AppName
	dispName := svcName
	if !isFixed {
		svcName = global.Def.GetLongAppName()
		parties := strings.Split(svcName, "_")
		dispName = fmt.Sprintf("%s(%s)", strings.Join(parties[:len(parties)-1], "_"), parties[len(parties)-1])
	}
	cfg := &service.Config{
		Name:        svcName,
		DisplayName: dispName,
		Description: global.Usage,
		Arguments:   args,
	}
	path, _ := filepath.Abs(os.Args[0])
	cfg.WorkingDirectory = filepath.Dir(path)
	return cfg
}

//GetSrvApp SrvCfg
func GetSrvApp(c *cli.Context) *ServiceApp {
	return &ServiceApp{
		c: c,
	}
}

//ServiceApp ServiceApp
type ServiceApp struct {
	c      *cli.Context
	server *servers.RspServers
	trace  itrace
}
