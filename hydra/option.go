package hydra

import (
	"fmt"
	"strings"
)

var (
	microServerType   = []string{"api", "rpc", "web"}
	flowServerType    = []string{"cron", "mqc"}
	wsServerType      = []string{"ws"}
	supportServerType = []string{"api", "rpc", "ws", "web", "cron", "mqc", "micro", "flow"}
)

type option struct {
	RegistryAddr       string `json:"--registry" valid:"ascii,required"`
	Name               string `json:"--name" `
	PlatName           string `json:"--plat" valid:"ascii,required"`
	SystemName         string `json:"--system" valid:"ascii,required"`
	ServerTypeNames    string `json:"--serverTypes" valid:"ascii,required"`
	ServerTypes        []string
	ClusterName        string `json:"--cluster" valid:"ascii,required"`
	IsDebug            bool
	Trace              string
	remoteLogger       bool
	RemoteLogger       bool
	RemoteQueryService bool
}

//Option 配置选项
type Option func(*option)

//WithRegistry 设置注册中心地址
func WithRegistry(addr string) Option {
	return func(o *option) {
		o.RegistryAddr = addr
	}
}

//WithPlatName 设置平台名称
func WithPlatName(platName string) Option {
	return func(o *option) {
		o.PlatName = platName
	}
}

//WithSystemName 设置系统名称
func WithSystemName(systemName string) Option {
	return func(o *option) {
		o.SystemName = systemName
	}
}

//WithServerTypes 设置系统类型
func WithServerTypes(serverType ...string) Option {
	return func(o *option) {
		o.ServerTypeNames = strings.Join(serverType, "-")
		for _, st := range serverType {
			sts, err := getServerTypes(st)
			if err != nil {
				panic(err)
			}
			o.ServerTypes = append(o.ServerTypes, sts...)
		}

	}
}

//WithClusterName 设置集群名称
func WithClusterName(clusterName string) Option {
	return func(o *option) {
		o.ClusterName = clusterName
	}
}

//WithRemoteQueryService 启动远程查询服务
func WithRemoteQueryService(remoteQueryService bool) Option {
	return func(o *option) {
		o.RemoteQueryService = remoteQueryService
	}
}

//WithName 设置系统全名 格式:/[platName]/[sysName]/[typeName]/[clusterName]
func WithName(name string) Option {
	return func(o *option) {
		o.Name = name
		var err error
		o.PlatName, o.SystemName, o.ServerTypes, o.ClusterName, err = parsePath(name)
		if err != nil {
			panic(fmt.Errorf("%s %v", name, err))
		}
		o.ServerTypeNames = strings.Join(o.ServerTypes, "-")
	}
}

//WithDebug 设置dubug模式
func WithDebug() Option {
	return func(o *option) {
		o.IsDebug = true
	}
}

//WithProduct 设置产品模式
func WithProduct() Option {
	return func(o *option) {
		o.IsDebug = false
	}
}

//WithRemoteLogger 设置产品模式
func WithRemoteLogger() Option {
	return func(o *option) {
		o.RemoteLogger = true
	}
}
func parsePath(p string) (platName string, systemName string, serverTypes []string, clusterName string, err error) {
	fs := strings.Split(strings.Trim(p, "/"), "/")
	if len(fs) != 4 {
		err := fmt.Errorf("系统名称错误，格式:/[platName]/[sysName]/[typeName]/[clusterName]")
		return "", "", nil, "", err
	}
	if serverTypes, err = getServerTypes(fs[2]); err != nil {
		return "", "", nil, "", err
	}
	platName = fs[0]
	systemName = fs[1]
	clusterName = fs[3]
	return
}

func getServerTypes(serverTypes string) ([]string, error) {
	sts := strings.Split(serverTypes, "-")
	removeRepMap := make(map[string]byte)
	for _, v := range sts {
		var ctn bool
		for _, k := range supportServerType {
			if ctn = k == v; ctn {
				break
			}
		}
		if !ctn {
			return nil, fmt.Errorf("不支持的服务器类型:%v", v)
		}
		switch v {
		case "*":
			for _, value := range microServerType {
				removeRepMap[value] = 0
			}
			for _, value := range flowServerType {
				removeRepMap[value] = 0
			}
			break
		case "micro":
			for _, value := range microServerType {
				removeRepMap[value] = 0
			}
		case "flow":
			for _, value := range flowServerType {
				removeRepMap[value] = 0
			}
		default:
			removeRepMap[v] = 0
		}
	}
	types := make([]string, 0, len(removeRepMap))
	for k := range removeRepMap {
		types = append(types, k)
	}
	return types, nil
}
