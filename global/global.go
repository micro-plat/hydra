package global

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/micro-plat/lib4go/logger"
	"github.com/micro-plat/lib4go/types"
	"github.com/urfave/cli"
)

var traces = []string{"cpu", "mem", "block", "mutex", "web"}

var isReady = false

//Def 默认appliction
var Def = &global{
	log:           logger.New("hydra"),
	LocalConfName: "./" + filepath.Base(os.Args[0]) + ".conf.toml",
	close:         make(chan struct{}),
}

//global 应用全局管理器
type global struct {

	//registryAddr 集群地址
	RegistryAddr string

	//PlatName 平台名称
	PlatName string

	//SysName 系统名称
	SysName string

	//ServerTypes 服务器类型
	ServerTypes []string

	//ServerTypeNames 服务类型名称
	ServerTypeNames string

	//ClusterName 集群名称
	ClusterName string

	//Name 服务器请求名称
	Name string

	//Trace 显示请求与响应信息
	Trace string

	//TracePort 当Trace为web时候，需要设置
	TracePort string

	//isClose 是否关闭当前应用程序
	isClose bool

	//log 日志管理
	log logger.ILogger

	//LocalConfName 本地配置文件名称
	LocalConfName string

	//close 关闭通道
	close chan struct{}
}

//Bind 处理app命令行参数
func (m *global) Bind(c *cli.Context) (err error) {
	//处理参数
	if err := m.check(); err != nil {
		return err
	}

	//处理服务回调
	if err := doCliCallback(c); err != nil {
		return err
	}

	//处理所有预处理函数
	if err := doReadyFuncs(); err != nil {
		return err
	}

	return nil
}

//GetLongAppName 获取包含有部分路径的app name
func (m *global) GetLongAppName(n ...string) string {
	name := types.GetStringByIndex(n, 0, AppName)
	path, _ := filepath.Abs(name)
	rname := strings.Trim(strings.Replace(path, string(filepath.Separator), "_", -1), "_")
	if len(rname) < 32 {
		return rname
	}
	return rname[len(rname)-32:]
}

//HasServerType 是否包含指定的服务类型
func (m *global) HasServerType(tp string) bool {
	for _, t := range m.ServerTypes {
		if t == tp {
			return true
		}
	}
	return false
}

//GetRegistryAddr 获取注册中心地址
func (m *global) GetRegistryAddr() string {
	return m.RegistryAddr
}

//GetPlatName 获取平台名称
func (m *global) GetPlatName() string {
	return m.PlatName
}

//GetSysName 获取系统名称
func (m *global) GetSysName() string {
	return m.SysName
}

//GetServerTypes 获取服务器类型
func (m *global) GetServerTypes() []string {
	return m.ServerTypes
}

//GetClusterName 获取集群名称
func (m *global) GetClusterName() string {
	return m.ClusterName
}

//GetTrace 获取当前启动的pprof类型
func (m *global) GetTrace() string {
	return m.Trace
}

//GetTracePort 当Trace为web时候，需要设置TracePort
func (m *global) GetTracePort() string {
	return m.TracePort
}

//ClosingNotify 获取系统关闭通知
func (m *global) ClosingNotify() chan struct{} {
	return m.close
}

//Log 获取日志组件
func (m *global) Log() logger.ILogger {
	return m.log
}

//parsePath 转换平台服务路径
func parsePath(p string) (platName string, systemName string, serverTypes []string, clusterName string, err error) {
	fs := strings.Split(strings.Trim(p, "/"), "/")
	if len(fs) != 4 {
		err := fmt.Errorf("系统名称错误，格式:/[platName]/[sysName]/[serverType]/[clusterName]")
		return "", "", nil, "", err
	}
	serverTypes = strings.Split(fs[2], "-")
	platName = fs[0]
	systemName = fs[1]
	clusterName = fs[3]
	return
}

//check 检查参数
func (m *global) check() (err error) {
	if FlagVal.RegistryAddr != "" {
		m.RegistryAddr = FlagVal.RegistryAddr
	}
	if FlagVal.Name != "" {
		m.Name = FlagVal.Name
	}
	if FlagVal.PlatName != "" {
		m.PlatName = FlagVal.PlatName
	}
	if FlagVal.SysName != "" {
		m.SysName = FlagVal.SysName
	}
	if FlagVal.ServerTypeNames != "" {
		m.ServerTypeNames = FlagVal.ServerTypeNames
	}
	if FlagVal.ClusterName != "" {
		m.ClusterName = FlagVal.ClusterName
	}
	if m.ServerTypeNames != "" {
		m.ServerTypes = strings.Split(strings.ToLower(m.ServerTypeNames), "-")
	}
	if m.Name != "" {
		m.PlatName, m.SysName, m.ServerTypes, m.ClusterName, err = parsePath(m.Name)
		if err != nil {
			return err
		}
	}
	for _, s := range m.ServerTypes {
		if !types.StringContains(ServerTypes, s) {
			return fmt.Errorf("%s不支持，只能是%v", s, ServerTypes)
		}
	}
	if m.SysName == "" {
		m.SysName = AppName
	}
	if m.RegistryAddr == "" {
		m.RegistryAddr = "lm://."
	}

	if m.RegistryAddr == "" {
		return fmt.Errorf("注册中心地址不能为空")
	}
	if m.PlatName == "" {
		return fmt.Errorf("平台名称不能为空")
	}
	if m.SysName == "" {
		return fmt.Errorf("系统名称不能为空")
	}
	if len(m.ServerTypes) == 0 {
		return fmt.Errorf("服务器类型不能为空")
	}
	if m.ClusterName == "" {
		return fmt.Errorf("集群名称不能为空")
	}
	if m.Trace != "" && !types.StringContains(traces, m.Trace) {
		return fmt.Errorf("trace名称只能是%v", traces)
	}
	//增加调试参数
	if IsDebug {
		m.PlatName += "_debug"
	}
	isReady = true
	return nil
}
