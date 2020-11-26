package pkgs

import (
	"fmt"
	"reflect"

	"github.com/micro-plat/cli/logs"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/hydra/hydra/cmds/pkgs/rlog"
	"github.com/micro-plat/hydra/hydra/cmds/pkgs/service"
	"github.com/micro-plat/hydra/hydra/servers"
	"github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/hydra/services"
	"github.com/urfave/cli"
)

type HydraService struct {
	service.Service
	ServiceName string
	Description string
	Arguments   []string
}

//GetService GetService
func GetService(c *cli.Context, args ...string) (hydraSrv *HydraService, err error) {
	//1. 构建服务配置
	cfg := GetSrvConfig(args...)

	//2.创建本地服务
	appSrv, err := service.New(GetSrvApp(c), cfg)
	if err != nil {
		return nil, err
	}
	return &HydraService{
		Service:     appSrv,
		ServiceName: cfg.Name,
		Description: cfg.Description,
		Arguments:   cfg.Arguments,
	}, err
}

//GetSrvConfig SrvCfg
func GetSrvConfig(args ...string) *service.Config {
	return &service.Config{
		Name:        global.Def.GetLongAppName(),
		DisplayName: global.Def.GetLongAppName(),
		Description: global.Usage,
		Arguments:   args,
	}
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
}

//Start Start
func (p *ServiceApp) Start(s service.Service) (err error) {
	err = p.run()
	return err
}

//Stop Stop
func (p *ServiceApp) Stop(s service.Service) error {

	//8. 关闭服务器释放所有资源
	global.Def.Log().Info(global.AppName, fmt.Sprintf("正在退出..."))

	if !reflect.ValueOf(p.server).IsNil() {
		//关闭服务器
		p.server.Shutdown()
	}

	//关闭各服务
	if err := services.Def.Close(); err != nil {
		global.Def.Log().Error("err:", err)
	}

	//通知关闭各组件
	global.Def.Close()

	global.Def.Log().Info(global.AppName, "已安全退出")
	return nil
}

func (p *ServiceApp) run() error {

	//1. 绑定应用程序参数
	if err := global.Def.Bind(p.c); err != nil {
		logs.Log.Error(err)
		cli.ShowCommandHelp(p.c, p.c.Command.Name)
		return nil
	}
	//3. 注册远程日志组件
	if err := rlog.Registry(global.Def.PlatName, global.Def.RegistryAddr); err != nil {
		logs.Log.Error(err)
		return nil
	}

	globalData := global.Current()
	//4.创建trace性能跟踪
	if err := startTrace(globalData.GetTrace(), globalData.GetTracePort()); err != nil {
		return err
	}
	//5. 处理本地内存作为注册中心的服务发布问题
	if registry.GetProto(globalData.GetRegistryAddr()) == registry.LocalMemory {
		if err := Pub2Registry(true); err != nil {
			return err
		}
	}
	//6. 创建服务器
	server := servers.NewRspServers(globalData.GetRegistryAddr(),
		globalData.GetPlatName(), globalData.GetSysName(),
		globalData.GetServerTypes(), globalData.GetClusterName())
	if err := server.Start(); err != nil {
		return err
	}
	p.server = server
	return nil
}
