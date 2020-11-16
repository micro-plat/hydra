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

//Etcd Etcd
const Etcd = "etcd"

//Redis redis
const Redis = "redis"

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

//NewRegistry 根据协议地址创建注册中心
func NewRegistry(address string, log logger.ILogging) (r IRegistry, err error) {
	proto, addrs, u, p, err := Parse(address)
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
		//addrs []string, userName string, password string, log logger.ILogging

		return rsvr.Create(Addrs(srvs...),
			WithAuthCreds(u, p),
			WithLogger(log),
			Domain(global.Def.PlatName))

	}, resolver, addrs, log)
	if err != nil {
		return
	}
	r = value.(IRegistry)
	return
}

//GetProto 获取协议名称
func GetProto(addr string) string {
	p, _, _, _, _ := Parse(addr)
	return p
}

//GetAddrs 获取地址信息
func GetAddrs(addr string) []string {
	_, addrs, _, _, _ := Parse(addr)
	return addrs
}

//Parse 解析地址
//如:zk://192.168.0.155:2181 或 fs://../
func Parse(address string) (proto string, raddr []string, u string, p string, err error) {
	if strings.Count(address, "://") != 1 {
		return "", nil, "", "", fmt.Errorf("%s，包含多个://。格式:[proto]://[address]", address)
	}

	addr := strings.SplitN(address, "://", 2)
	if len(addr) != 2 {
		return "", nil, "", "", fmt.Errorf("%s，必须包含://。格式:[proto]://[address]", address)
	}
	if len(addr[0]) == 0 {
		return "", nil, "", "", fmt.Errorf("%s，协议名不能为空。格式:[proto]://[address]", address)
	}
	if len(addr[1]) == 0 {
		return "", nil, "", "", fmt.Errorf("%s，地址不能为空。格式:[proto]://[address]", address)
	}
	proto = addr[0]
	raddr = strings.Split(addr[1], ",")
	var addr0 string
	u, p, addr0, err = getAddrByUserPass(raddr[0])
	raddr[0] = addr0
	return
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
func getAddrByUserPass(addr string) (u string, p string, address string, err error) {
	if !strings.Contains(addr, "@") {
		return "", "", addr, nil
	}
	addrs := strings.Split(addr, "@")
	if len(addrs) != 2 {
		return "", "", "", fmt.Errorf("地址非法%s", addr)
	}
	address = addrs[1]
	up := strings.Split(addrs[0], ":")
	switch len(up) {
	case 1:
		return up[0], up[0], address, nil
	case 2:
		return up[0], up[1], address, nil
	default:
		return "", "", "", fmt.Errorf("地址非法%s", addrs[0])
	}
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
