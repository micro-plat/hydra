package application

import (
	"fmt"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/registry/conf/server"
	"github.com/micro-plat/lib4go/logger"
)

//DefApp 默认app
var DefApp = &application{
	log: logger.New("hydra"),
}

type application struct {

	//registryAddr 集群地址
	RegistryAddr string `json:"registryAddr" valid:"ascii,required"`

	//PlatName 平台名称
	PlatName string `json:"platName" valid:"ascii,required"`

	//SysName 系统名称
	SysName string `json:"sysName" valid:"ascii,required"`

	//ServerTypes 服务器类型
	ServerTypes []string `json:"serverTypes" valid:"in(api|web|rpc|ws|mqc|cron),required"`

	//ServerTypeNames 服务类型名称
	ServerTypeNames string

	//ClusterName 集群名称
	ClusterName string `json:"clusterName" valid:"ascii,required"`

	//Name 服务器请求名称
	Name string

	//Trace 显示请求与响应信息
	Trace string `valid:"in(cpu|mem|block|mutex|web)"`

	//isClose 是否关闭当前应用程序
	isClose bool

	//log 日志管理
	log logger.ILogger

	//close 关闭通道
	close chan struct{}
}

func (m *application) Bind() (err error) {
	if m.ServerTypeNames != "" {
		m.ServerTypes = strings.Split(m.ServerTypeNames, "-")
	}
	if m.Name != "" {
		m.PlatName, m.SysName, m.ServerTypes, m.ClusterName, err = parsePath(m.Name)
		if err != nil {
			return err
		}
	}
	_, err = govalidator.ValidateStruct(m)
	if err != nil {
		return err
	}
	if IsDebug {
		m.PlatName += "_debug"
	}
	return nil
}

func (a *application) GetCMD() string {
	return ""
}

//GetHandler 获取服务对应的处理函数
func (a *application) GetHandler(tp string, service string) context.IHandler {
	return nil
}

//Server 获取服务器配置信息
func (a *application) Server(tp string) server.IServerConf {
	s, err := server.Cache.GetServerConf(tp)
	if err == nil {
		return s
	}
	panic(fmt.Errorf("[%s]服务器未启动:%w", tp, err))
}

//CurrentContext 获取当前请求上下文
func (a *application) CurrentContext() context.IContext {
	return nil
}

//GetRegistryAddr 注册中心
func (a *application) GetRegistryAddr() string {
	return a.RegistryAddr
}

//GetPlatName 平台名称
func (a *application) GetPlatName() string {
	return a.PlatName
}

//GetSysName 系统名称
func (a *application) GetSysName() string {
	return a.SysName
}

//GetServerTypes 服务器类型
func (a *application) GetServerTypes() []string {
	return a.ServerTypes
}

//GetClusterName 集群名称
func (a *application) GetClusterName() string {
	return a.ClusterName
}

//GetTrace 显示请求与响应信息
func (a *application) GetTrace() string {
	return a.Trace
}

//ClosingNotify 获取系统关闭通知
func (a *application) ClosingNotify() chan struct{} {
	return a.close
}

//Log 获取日志组件
func (a *application) Log() logger.ILogger {
	return a.log
}

//Close 显示请求与响应信息
func (a *application) Close() {
	a.isClose = true
	close(a.close)
}
func parsePath(p string) (platName string, systemName string, serverTypes []string, clusterName string, err error) {
	fs := strings.Split(strings.Trim(p, "/"), "/")
	if len(fs) != 4 {
		err := fmt.Errorf("系统名称错误，格式:/[platName]/[sysName]/[typeName]/[clusterName]")
		return "", "", nil, "", err
	}
	serverTypes = strings.Split(fs[2], "-")
	platName = fs[0]
	systemName = fs[1]
	clusterName = fs[3]
	return
}
