package hydra

import (
	"github.com/micro-plat/hydra/application"
)

//Option 配置选项
type Option func()

//WithRegistry 设置注册中心地址
func WithRegistry(addr string) Option {
	return func() {
		application.RegistryAddr = addr
	}
}

//WithPlatName 设置平台名称
func WithPlatName(platName string) Option {
	return func() {
		application.PlatName = platName
	}
}

//WithSystemName 设置系统名称
func WithSystemName(sysName string) Option {
	return func() {
		application.SysName = sysName
	}
}

//WithServerTypes 设置系统类型
func WithServerTypes(serverType ...string) Option {
	return func() {
		application.ServerTypes = serverType
	}
}

//WithClusterName 设置集群名称
func WithClusterName(clusterName string) Option {
	return func() {
		application.ClusterName = clusterName
	}
}

//WithName 设置系统全名 格式:/[platName]/[sysName]/[typeName]/[clusterName]
func WithName(name string) Option {
	return func() {
		application.Name = name
	}
}

//WithDebug 设置dubug模式
func WithDebug() Option {
	return func() {
		application.IsDebug = true
	}
}
