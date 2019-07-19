package balancer

import (
	"time"

	"fmt"

	"github.com/micro-plat/hydra/registry"
	"google.golang.org/grpc/naming"
)

type ServiceResolver interface {
	// Resolve creates a Watcher for target.
	Resolve(target string) (naming.Watcher, error)
	Close()
}

//Resolver 服务解析器,用于解析不同的注册中心地址,创建注册中心watcher
type Resolver struct {
	timeout    time.Duration
	service    string
	sortPrefix string
	closeCh    []*Watcher
}

//NewResolver 返回服务解析器
func NewResolver(service string, timeout time.Duration, sortPrefix string) ServiceResolver {
	return &Resolver{timeout: timeout, service: service, sortPrefix: sortPrefix, closeCh: make([]*Watcher, 0, 2)}
}

// Resolve to resolve the service from zookeeper, target is the dial address of zookeeper
// target example: "zk://192.168.0.159:2181,192.168.0.154:2181"
func (v *Resolver) Resolve(target string) (naming.Watcher, error) {
	r, err := registry.NewRegistryWithAddress(target, nil)
	if err != nil {
		return nil, fmt.Errorf("rpc.client.resolver target err:%v", err)
	}
	rw := &Watcher{client: r, service: v.service, sortPrefix: v.sortPrefix, closeCh: make(chan struct{})}
	v.closeCh = append(v.closeCh, rw)
	return rw, nil
}

//Close 关闭所有watcher
func (v *Resolver) Close() {
	for _, c := range v.closeCh {
		c.Close()
	}
}
