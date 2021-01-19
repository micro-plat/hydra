package balancer

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/lib4go/logger"

	//"google.golang.org/grpc/naming"

	"google.golang.org/grpc/attributes"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/resolver/manual"
)

/*
// Builder creates a resolver that will be used to watch name resolution updates.
type Builder interface {
	// Build creates a new resolver for the given target.
	//
	// gRPC dial calls Build synchronously, and fails if the returned error is
	// not nil.
	Build(target Target, cc ClientConn, opts BuildOptions) (Resolver, error)
	// Scheme returns the scheme supported by this resolver.
	// Scheme is defined at https://github.com/grpc/grpc/blob/master/doc/naming.md.
	Scheme() string
}
*/

//ResolverBuilder creates a resolver that will be used to watch name resolution updates.
type ResolverBuilder struct {
	address     string
	proto       string
	plat        string
	service     string
	sortPrefix  string
	caches      map[string]bool
	logger      *logger.Logger
	regst       registry.IRegistry
	orgResolver *manual.Resolver
	closeChan   chan struct{}
	isClose     bool
	lock        sync.Mutex
	onceLock    sync.Once
}

var _ resolver.Builder = &ResolverBuilder{}

//NewResolverBuilder 新建builder
func NewResolverBuilder(address, plat, service, sortPrefix string) (*ResolverBuilder, error) {
	proto, addr, err := global.ParseProto(address)
	if err != nil {
		return nil, fmt.Errorf("GRPC address:%s parse error:%+v", address, err)
	}
	logging := logger.New("rpc.resolve")

	builder := &ResolverBuilder{
		plat:       plat,
		service:    service,
		sortPrefix: sortPrefix,
		address:    address,
		proto:      proto,
		logger:     logging,
		closeChan:  make(chan struct{}),
		caches:     map[string]bool{},
	}

	addresses := []string{addr}
	//兼容直接传服务器ip来进行访问
	if len(plat) > 0 {
		regst, err := registry.GetRegistry(address, logging)
		if err != nil {
			return nil, fmt.Errorf("rpc.client.resolver target err:%v", err)
		}

		builder.regst = regst
		addresses, err = builder.getGrpcAddress()
		if err != nil {
			return nil, fmt.Errorf("rpc.client.resolver target err:%v", err)
		}
		go builder.watchChildren()
	}

	builder.buildManualResolver(proto, addresses)
	return builder, nil
}

// Build creates a new resolver for the given target.
func (b *ResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver resolver.Resolver, err error) {
	return b.orgResolver.Build(target, cc, opts)
}

// Scheme returns the scheme supported by this resolver.
func (b *ResolverBuilder) Scheme() string {
	return b.proto
}

func (b *ResolverBuilder) buildManualResolver(proto string, address []string) {
	rb := manual.NewBuilderWithScheme(proto)
	rb.ResolveNowCallback = func(o resolver.ResolveNowOptions) {}
	var grpcAddrs []resolver.Address
	for i := range address {
		grpcAddrs = append(grpcAddrs, resolver.Address{
			Addr:       address[i],
			Type:       resolver.Backend,
			Attributes: attributes.New(),
		})
		b.caches[address[i]] = true
	}

	rb.InitialState(resolver.State{Addresses: grpcAddrs})
	b.orgResolver = rb
}

func (b *ResolverBuilder) getGrpcAddress() (addrs []string, err error) {

	rpath, err := b.getRealPath()
	if err != nil {
		return []string{}, err
	}

	//获取所有rpc服务下的子节点
	children, _, err := b.regst.GetChildren(rpath)
	if err != nil {
		return []string{}, fmt.Errorf("GetChildren服务地址出错 %s %w", rpath, err)
	}

	addrs = b.extractAddrs(children)
	return
}

func (b *ResolverBuilder) extractAddrs(resp []string) []string {
	addrs := make([]string, 0, len(resp))
	for _, v := range resp {
		item := strings.SplitN(v, "_", 2)
		addrs = append(addrs, item[0])
	}
	if b.sortPrefix != "" {
		sort.Slice(addrs, func(i, j int) bool {
			return strings.HasPrefix(addrs[i], b.sortPrefix)
		})
	}
	return addrs
}

func (b *ResolverBuilder) getRealPath() (string, error) {
	rpath := registry.Join(b.plat, "services", "rpc", b.service, "providers")
	v, err := b.regst.Exists(rpath)
	if err != nil {
		return "", fmt.Errorf("检查新版服务地址出错 %s %w", rpath, err)
	}
	if v {
		return rpath, nil
	}

	//如果不存在  检查是否是老版本的配置   按照老版本配置查找
	if strings.Contains(b.plat, ".") {
		p := strings.Split(b.plat, ".")
		if len(p) != 2 {
			return "", fmt.Errorf("配置服务地址错误 %s", b.plat)
		}
		rpath = registry.Join(p[1], "services", "rpc", p[0], b.service, "providers")
		v, err = b.regst.Exists(rpath)
		if err != nil {
			return "", fmt.Errorf("检查老版服务地址出错 %s %w", rpath, err)
		}
		if v {
			return rpath, nil
		}
	}

	rpath = rpath[:len(rpath)-len(registry.Join(b.service, "providers"))]
	//精确路径没有找到  现在需要获取模糊匹配路径
	sp := strings.Split(strings.Trim(b.service, "/"), "/")
	if len(sp) == 0 {
		return "", fmt.Errorf("service服务路径错误：%s", b.service)
	}

	//递归获取真实的路径
	path, ok, err := b.getPath(rpath, sp, 0)
	if err != nil {
		return "", err
	}

	if !ok {
		return "", fmt.Errorf("未找到服务提供程序：%s", b.service)
	}

	rpath = registry.Join(path, "providers")
	return rpath, nil
}

func (b *ResolverBuilder) getPath(rpath string, sp []string, index int) (string, bool, error) {
	if index >= len(sp) {
		return rpath, false, nil
	}

	chilren, _, err := b.regst.GetChildren(rpath)
	if err != nil {
		return "", false, fmt.Errorf("GetChildren服务地址出错 %s %w", rpath, err)
	}

	ok := false
	for _, str := range chilren {

		if strings.Contains(str, sp[index]) || strings.HasPrefix(str, ":") {
			xpath := registry.Join(rpath, str)
			if index == len(sp)-1 {
				return xpath, true, nil
			}
			xpath, ok, err = b.getPath(xpath, sp, index+1)
			if err != nil {
				return "", false, err
			}

			if ok {
				return xpath, ok, nil
			}
		}
	}

	return "", false, nil
}

func (b *ResolverBuilder) watchChildren() {
	if len(b.plat) <= 0 {
		return
	}

	for {
		if b.isClose {
			return
		}
		realPath, err := b.getRealPath()
		if err != nil {
			b.logger.Errorf("watchChildren.获取注册中心路径:%+v", err)
			time.Sleep(time.Second)
			continue
		}
		childChan, err := b.regst.WatchChildren(realPath)
		if err != nil {
			b.logger.Errorf("watchChildren.监控路径%s:%+v", realPath, err)
			time.Sleep(time.Second)
			continue
		}
		select {
		case <-b.closeChan:
			return
		case <-childChan:
			b.updateAddrs()
		}
	}
}

func (b *ResolverBuilder) updateAddrs() {
	if len(b.plat) <= 0 || b.isClose {
		return
	}
	b.lock.Lock()
	defer b.lock.Unlock()

	address, err := b.getGrpcAddress()
	if err != nil {
		b.logger.Errorf("获取grpc地址错误:%+v", err)
		return
	}

	if !b.checkUpdate(address) {
		return
	}

	var grpcAddrs []resolver.Address
	for i := range address {
		grpcAddrs = append(grpcAddrs, resolver.Address{
			Addr:       address[i],
			Type:       resolver.Backend,
			Attributes: attributes.New(),
		})
	}
	b.orgResolver.CC.UpdateState(resolver.State{Addresses: grpcAddrs})
}

func (b *ResolverBuilder) checkUpdate(address []string) bool {
	var needUpdate = false
	if len(address) != len(b.caches) {
		needUpdate = true
	}
	newCache := make(map[string]bool)
	for i := 0; i < len(address); i++ {
		newCache[address[i]] = true
		if _, ok := b.caches[address[i]]; !ok {
			needUpdate = true
		}
	}
	b.caches = newCache
	return needUpdate
}

//Close 关闭builder
func (b *ResolverBuilder) Close() {
	b.onceLock.Do(func() {
		b.isClose = true
		close(b.closeChan)
	})
}
