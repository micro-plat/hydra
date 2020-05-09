package registry

import (
	"fmt"

	"strings"

	"github.com/micro-plat/lib4go/concurrent/cmap"
	"github.com/micro-plat/lib4go/logger"
	"github.com/micro-plat/lib4go/registry"
)

//IRegistry 注册中心接口
type IRegistry interface {
	WatchChildren(path string) (data chan registry.ChildrenWatcher, err error)
	WatchValue(path string) (data chan registry.ValueWatcher, err error)
	GetChildren(path string) (paths []string, version int32, err error)
	GetValue(path string) (data []byte, version int32, err error)
	CreatePersistentNode(path string, data string) (err error)
	CreateTempNode(path string, data string) (err error)
	CreateSeqNode(path string, data string) (rpath string, err error)
	Update(path string, data string, version int32) (err error)
	Delete(path string) error
	Exists(path string) (bool, error)
	Close() error
}

//IFactory 注册中心构建器
type IFactory interface {
	Create(addrs []string, userName string, password string, log logger.ILogging) (IRegistry, error)
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
		return nil, fmt.Errorf("registry: unknown adapter name %q (forgotten import?)", proto)
	}
	key := fmt.Sprintf("%s_%s", proto, strings.Join(addrs, "_"))
	_, value, err := registryMap.SetIfAbsentCb(key, func(input ...interface{}) (interface{}, error) {
		rsvr := input[0].(IFactory)
		srvs := input[1].([]string)
		log := input[2].(logger.ILogging)
		return rsvr.Create(srvs, u, p, log)
	}, resolver, addrs, log)
	if err != nil {
		return
	}
	r = value.(IRegistry)
	return
}

//Parse 解析地址
//如:zk://192.168.0.155:2181 或 fs://../
func Parse(address string) (proto string, raddr []string, u string, p string, err error) {
	addr := strings.SplitN(address, "://", 2)
	if len(addr) != 2 {
		return "", nil, "", "", fmt.Errorf("%s错误，必须包含://", address)
	}
	if len(addr[0]) == 0 {
		return "", nil, "", "", fmt.Errorf("%s错误，协议名不能为空", address)
	}
	if len(addr[1]) == 0 {
		return "", nil, "", "", fmt.Errorf("%s错误，地址不能为空", address)
	}
	proto = addr[0]
	raddr = strings.Split(addr[1], ",")
	var addr0 string
	u, p, addr0, err = getAddrByUserPass(raddr[0])
	raddr[0] = addr0
	return
}

//Join 地址连接
func Join(elem ...string) string {
	var builder strings.Builder
	builder.WriteString("/")
	for _, v := range elem {
		if v == "/" || v == "\\" {
			continue
		}
		builder.WriteString(strings.Trim(v, "/"))
		builder.WriteString("/")
	}
	return strings.TrimSuffix(builder.String(), "/")
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
