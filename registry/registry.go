package registry

import (
	"fmt"

	"strings"

	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/lib4go/concurrent/cmap"
	"github.com/micro-plat/lib4go/logger"
	"github.com/micro-plat/lib4go/registry"
)

//LocalMemory 本地内存模式
const LocalMemory = "lm"

//Zookeeper zk
const Zookeeper = "zk"

//FileSystem 本地文件系统
const FileSystem = "fs"

//Consul Consul
const Consul = "consul"

//Redis redis
const Redis = "redis"

//Mysql mysql
const Mysql = "mysql"

//IRegistry 注册中心接口
type IRegistry interface {
	WatchChildren(path string) (data chan registry.ChildrenWatcher, err error)
	WatchValue(path string) (data chan registry.ValueWatcher, err error)
	GetChildren(path string) (paths []string, version int32, err error)
	GetValue(path string) (data []byte, version int32, err error)
	CreatePersistentNode(path string, data string) (err error)
	CreateTempNode(path string, data string) (err error)
	CreateSeqNode(path string, data string) (rpath string, err error)
	Update(path string, data string) (err error)
	Delete(path string) error
	Exists(path string) (bool, error)
	Close() error
}

//IFactory 注册中心构建器
type IFactory interface {
	Create(...Option) (IRegistry, error)
}

var registryMap = cmap.New(2)
var registries = make(map[string]IFactory)

//Register 添加注册中心工厂对象
func Register(name string, builder IFactory) {
	if builder == nil {
		panic("registry: Register adapter is nil")
	}
	if _, ok := registries[name]; ok {
		panic("registry: Register called twice for adapter " + name)
	}
	registries[name] = builder
}

//CreateRegistry 创建新的的注册中心
func CreateRegistry(address string, log logger.ILogging) (r IRegistry, err error) {
	proto, addrs, u, p, mt, err := Parse(address)
	if err != nil {
		return nil, err
	}
	resolver, ok := registries[proto]
	if !ok {
		return nil, fmt.Errorf("不支持的协议类型[%s]", proto)
	}
	return resolver.Create(Addrs(addrs...),
		WithAuthCreds(u, p),
		WithLogger(log),
		WithDomain(global.Def.PlatName), WithMetadata(mt))

}

//Support 检查注册中心地址是否支持
func Support(address string) bool {
	proto, _, _, _, _, err := Parse(address)
	if err != nil {
		fmt.Println(err)
		return false
	}
	_, ok := registries[proto]
	return ok
}

//GetCurrent 获取当前注册中心
func GetCurrent() IRegistry {
	r, _ := GetRegistry(global.Def.RegistryAddr, global.Def.Log())
	return r
}

//GetRegistry 获取缓存的注册中心
func GetRegistry(address string, log logger.ILogging) (r IRegistry, err error) {
	proto, addrs, u, p, mt, err := Parse(address)
	if err != nil {
		return nil, err
	}
	resolver, ok := registries[proto]
	if !ok {
		return nil, fmt.Errorf("不支持的协议类型[%s]", proto)
	}
	key := fmt.Sprintf("%s_%s", proto, strings.Join(addrs, "_"))
	_, value, err := registryMap.SetIfAbsentCb(key, func(input ...interface{}) (interface{}, error) {
		rsvr := input[0].(IFactory)
		srvs := input[1].([]string)
		log := input[2].(logger.ILogging)

		return rsvr.Create(Addrs(srvs...),
			WithAuthCreds(u, p),
			WithLogger(log),
			WithDomain(global.Def.PlatName), WithMetadata(mt))

	}, resolver, addrs, log)
	if err != nil {
		return
	}
	r = value.(IRegistry)
	return
}

//GetProto 获取协议名称
func GetProto(addr string) string {
	p, _, _, _, _, _ := Parse(addr)
	return p
}

//GetAddrs 获取地址信息
func GetAddrs(addr string) []string {
	_, addrs, _, _, _, _ := Parse(addr)
	return addrs
}

//Parse 解析地址
//如:zk://192.168.0.155:2181 或 fs://../
func Parse(address string) (proto string, raddr []string, u string, p string, mt map[string]string, err error) {

	//=========================================检查地址是否合法=======================================
	//必须包含协议，地址两部份
	if strings.Count(address, "://") != 1 {
		return "", nil, "", "", nil, fmt.Errorf("%s，包含多个://。格式:[proto]://[address]", address)
	}

	//检查协议与地址是否合法
	addr := strings.SplitN(address, "://", 2)
	if len(addr) != 2 {
		return "", nil, "", "", nil, fmt.Errorf("%s，必须包含://。格式:[proto]://[address]", address)
	}
	if len(addr[0]) == 0 {
		return "", nil, "", "", nil, fmt.Errorf("%s，协议名不能为空。格式:[proto]://[address]", address)
	}
	if len(addr[1]) == 0 {
		return "", nil, "", "", nil, fmt.Errorf("%s，地址不能为空。格式:[proto]://[address]", address)
	}

	//=========================================处理用户密码与地址=======================================
	proto = addr[0]
	originAddr := addr[1]
	at := strings.Split(originAddr, "@") //检查是否带有用户名密码信息
	if len(at) > 2 {
		return "", nil, "", "", nil, fmt.Errorf("不能包含多个@符号%s", address)
	}
	ups := ""
	currentAddr := at[0]
	if len(at) == 2 {
		ups = at[0]
		currentAddr = at[1]
	}

	//获取用户名密码
	u, p, err = getUP(ups)
	if err != nil {
		return
	}

	//获取地址
	mt = make(map[string]string)
	mt["db"], raddr, err = getAddr(currentAddr)
	if err != nil {
		return
	}
	return
}
func getUP(up string) (u string, p string, err error) {
	if len(up) == 0 {
		return "", "", nil
	}
	ups := strings.Split(up, ":")
	switch len(ups) {
	case 1:
		return ups[0], "", nil
	case 2:
		return ups[0], ups[1], nil
	default:
		return "", "", fmt.Errorf(`地址错误，不能包含多个":"(%s)`, up)
	}

}
func getAddr(originAddrs string) (d string, addrs []string, err error) {
	naddrs := strings.Split(originAddrs, "#")
	switch len(naddrs) {
	case 1:
		return "", strings.Split(naddrs[0], ","), nil
	case 2:
		return naddrs[0], strings.Split(naddrs[1], ","), nil
	default:
		return "", nil, fmt.Errorf(`地址错误，不能包含多个"#"(%s)`, originAddrs)
	}
}

//Format 格式化注册中心地址
func Format(ele string) string {
	return Join(ele)
}

//Join 地址连接
func Join(elem ...string) string {
	var builder strings.Builder
	builder.WriteString("/")
	for _, v := range elem {
		if v == "/" || v == "\\" || strings.TrimSpace(v) == "" {
			continue
		}
		builder.WriteString(strings.Trim(v, "/"))
		builder.WriteString("/")
	}
	return strings.TrimSuffix(builder.String(), "/")
}

//Split 将路径分隔为多段数组
func Split(path string) []string {
	return strings.Split(Trim(path), "/")
}

//Trim 去掉前后的""/"
func Trim(l string) string {
	return strings.Trim(l, "/")
}

//Close 关闭注册中心的服务
func Close() {
	registryMap.RemoveIterCb(func(key string, value interface{}) bool {
		if v, ok := value.(IRegistry); ok {
			v.Close()
		}
		return true
	})
}
