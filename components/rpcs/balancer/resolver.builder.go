package balancer

import (
	"fmt"
	"sort"
	"strings"

	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/lib4go/logger"

	//"google.golang.org/grpc/naming"

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
	address    string
	proto      string
	plat       string
	service    string
	sortPrefix string
	caches     map[string]bool
	logger     *logger.Logger

	regst       registry.IRegistry
	orgResolver *manual.Resolver
}

func NewResolverBuilder(address, plat, service, sortPrefix string) (resolver.Builder, error) {
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
	}

	builder.buildManualResolver(proto, addresses)
	return builder, nil
}

// Build creates a new resolver for the given target.
func (b *ResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver resolver.Resolver, err error) {
	//fmt.Println(" (b *ResolverBuilder) Build:", target.Scheme, target.Authority, target.Endpoint)
	return b.orgResolver.Build(target, cc, opts)
}

// Scheme returns the scheme supported by this resolver.
func (b *ResolverBuilder) Scheme() string {
	return b.proto
}

func (b *ResolverBuilder) buildManualResolver(proto string, address []string) {
	//fmt.Println("buildManualResolver:", proto)
	rb := manual.NewBuilderWithScheme(proto)

	rb.ResolveNowCallback = func(o resolver.ResolveNowOptions) {

		if len(b.plat) <= 0 {
			return
		}

		address, err := b.getGrpcAddress()
		if err != nil {
			b.logger.Errorf("getGrpcAddress:%+v", err)
			return
		}
		var needUpdate = false
		newCache := make(map[string]bool)
		for i := 0; i < len(address); i++ {
			newCache[address[i]] = true
			if _, ok := b.caches[address[i]]; !ok {
				needUpdate = true
			}
		}
		b.caches = newCache

		if !needUpdate {
			return
		}

		var grpcAddrs []resolver.Address
		for i := range address {
			grpcAddrs = append(grpcAddrs, resolver.Address{Addr: address[i], Type: resolver.Backend})
		}
		rb.CC.UpdateState(resolver.State{Addresses: grpcAddrs})
	}
	var grpcAddrs []resolver.Address
	for i := range address {
		grpcAddrs = append(grpcAddrs, resolver.Address{Addr: address[i], Type: resolver.Backend})
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
	chilren, _, err := b.regst.GetChildren(rpath)
	if err != nil {
		return []string{}, fmt.Errorf("GetChildren服务地址出错 %s %w", rpath, err)
	}

	addrs = b.extractAddrs(chilren)
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
		return "", fmt.Errorf("service服务路径错误, %s", b.service)
	}

	//递归获取真实的路径
	path, ok, err := b.getPath(rpath, sp, 0)
	if err != nil {
		return "", err
	}

	if !ok {
		return "", fmt.Errorf("没有找到有效的服务路径, %s", b.service)
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

/*
// Watcher is the implementaion of grpc.naming.Watcher
type Watcher struct {
	client        registry.IRegistry
	isInitialized bool
	caches        map[string]bool
	plat          string
	path          string
	service       string
	sortPrefix    string
	closeCh       chan struct{}
	lastErr       error
	once          sync.Once
}

// Close do nothing
func (w *Watcher) Close() {
	w.once.Do(func() {
		close(w.closeCh)
	})
}

// Next 监控服务器地址变化,监控发生异常时移除所有服务,否则等待服务器地址变化
func (w *Watcher) Next() ([]*naming.Update, error) {
start:
	w.lastErr = nil
	if !w.isInitialized {
		path, err := w.initialize()

		if err != nil {
			return w.getUpdates([]string{}), nil
		}
		w.path = path

		resp, _, err := w.client.GetChildren(w.path)
		if err == nil {
			w.isInitialized = true
			addrs := w.extractAddrs(resp)
			return w.getUpdates(addrs), nil
		}
	}

	// generate etcd/zk Watcher
	watcherCh, err := w.client.WatchChildren(w.path)
	if err != nil {
		return nil, fmt.Errorf("rpc.client.未找到服务:%s(err:%v)", w.path, err)
	}
	var watcher r.ChildrenWatcher
	select {
	case watcher = <-watcherCh:
	case <-w.closeCh:
		return w.getUpdates([]string{}), w.lastErr
	}
	if err = watcher.GetError(); err != nil {
		return nil, err
	}
	chilren, _ := watcher.GetValue()
	if len(chilren) == 0 {
		w.isInitialized = false
		goto start
	}
	addrs := w.extractAddrs(chilren)
	return w.getUpdates(addrs), nil
}
func (w *Watcher) getUpdates(addrs []string) (updates []*naming.Update) {
	newCache := make(map[string]bool)
	for i := 0; i < len(addrs); i++ {
		newCache[addrs[i]] = true
		if _, ok := w.caches[addrs[i]]; !ok {
			updates = append(updates, &naming.Update{Op: naming.Add, Addr: addrs[i]})
		} else {
			w.caches[addrs[i]] = false
		}
	}
	for i, v := range w.caches {
		if v {
			updates = append(updates, &naming.Update{Op: naming.Delete, Addr: i})
		}
	}
	w.caches = newCache
	return
}
func (w *Watcher) extractAddrs(resp []string) []string {
	addrs := make([]string, 0, len(resp))
	for _, v := range resp {
		item := strings.SplitN(v, "_", 2)
		addrs = append(addrs, item[0])
	}
	if w.sortPrefix != "" {
		sort.Slice(addrs, func(i, j int) bool {
			return strings.HasPrefix(addrs[i], w.sortPrefix)
		})
	}
	return addrs
}

func (w *Watcher) initialize() (string, error) {

	//根据绝对路径查询服务
	rpath := registry.Join(w.plat, "services", "rpc", w.service, "providers")
	b, err := w.client.Exists(rpath)
	if err != nil {
		return "", fmt.Errorf("检查服务地址出错 %s %w", rpath, err)
	}
	if b {
		return rpath, nil
	}

	//查询模糊匹配的节点
	root := registry.Join(w.plat, "services", "rpc")
	items := registry.Split(w.service)

	list, err := w.findServicePath([]string{root}, items...)
	if err != nil {
		return "", fmt.Errorf("获取服务地址失败 %w", err)
	}

	//对查询到的路径进行排序
	sort.Strings(list)
	if len(list) == 0 {
		return "", fmt.Errorf("未找到服务提供程序:%s", rpath)
	}
	return list[0], nil
}

func (w *Watcher) findServicePath(roots []string, names ...string) ([]string, error) {

	//查找完成所有节点
	if len(names) == 0 {
		providers := make([]string, 0, len(roots))
		for _, r := range roots {
			providers = append(providers, registry.Join(r, "providers"))
		}
		return providers, nil
	}

	//查找匹配的节点
	matchPaths := make([]string, 0, 1)
	for _, root := range roots {
		paths, _, err := w.client.GetChildren(root)
		if err != nil {
			return nil, err
		}
		for _, p := range paths {
			if p == names[0] || strings.HasPrefix(p, ":") {
				matchPaths = append(matchPaths, registry.Join(root, p))
			}
		}
	}
	if len(matchPaths) == 0 {
		return nil, fmt.Errorf("未找到匹配的路径：%v %s", roots, names[0])
	}
	return w.findServicePath(matchPaths, names[1:]...)
}
*/
