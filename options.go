package hydra

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/lib4go/types"
)

//Option 配置选项
type Option func()

//WithRegistry 设置注册中心地址
func WithRegistry(addr string) Option {
	return func() {
		global.Def.RegistryAddr = addr
	}
}

//WithPlatName 设置平台名称
func WithPlatName(platName string) Option {
	return func() {
		global.Def.PlatName = platName
	}
}

//WithSystemName 设置系统名称
func WithSystemName(sysName string) Option {
	return func() {
		global.Def.SysName = sysName
	}
}

//WithServerTypes 设置系统类型
func WithServerTypes(serverType ...string) Option {
	return func() {
		for _, s := range serverType {
			if !types.StringContains(global.Def.ServerTypes, s) {
				global.Def.ServerTypes = append(global.Def.ServerTypes, s)
			}
		}
	}
}

//WithClusterName 设置集群名称
func WithClusterName(clusterName string) Option {
	return func() {
		global.Def.ClusterName = clusterName
	}
}

//WithName 设置系统全名 格式:/[platName]/[sysName]/[typeName]/[clusterName]
func WithName(name string) Option {
	return func() {
		global.Def.Name = name
	}
}

//WithDebug 设置dubug模式
func WithDebug() Option {
	return func() {
		global.IsDebug = true
	}
}

//WithVersion 设置当前软件版本号
func WithVersion(v string) Option {
	return func() {
		global.Version = v
	}
}

//WithUsage 设置使用说明
func WithUsage(usage string) Option {
	return func() {
		global.Usage = fmt.Sprintf("%s(%s)", filepath.Base(os.Args[0]), usage)
	}
}

//WithRunFlag 添加run命令扩展参数
func WithRunFlag(name string, usage string) Option {
	return func() {
		global.RunCli.AddFlag(name, usage)
	}
}

//WithConfFlag 添加conf命令扩展参数
func WithConfFlag(name string, usage string) Option {
	return func() {
		global.ConfCli.AddFlag(name, usage)
	}
}
 