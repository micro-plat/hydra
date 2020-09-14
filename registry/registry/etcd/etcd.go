// Package etcd provides an etcd service registry
package etcd

import (
	"crypto/tls"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/micro-plat/hydra/registry"
	"github.com/zhiyunliu/etcd/clientv3"
)

type etcdRegistry struct {
	client  *clientv3.Client
	options *registry.Options
	CloseCh chan struct{}
	// register and leases are grouped by domain
	sync.RWMutex
	leases sync.Map

	seqValue int32
}

type leases clientv3.LeaseID

// NewRegistry returns an initialized etcd registry
func NewRegistry(opts ...registry.Option) *etcdRegistry {

	e := &etcdRegistry{
		options: &registry.Options{},
	}
	e.CloseCh = make(chan struct{})
	configure(e, opts...)
	go e.leaseRemain()
	return e
}

func newClient(e *etcdRegistry) (*clientv3.Client, error) {
	config := clientv3.Config{
		Endpoints: []string{"127.0.0.1:2379"},
	}
	if e.options.Timeout == 0 {
		e.options.Timeout = 5 * time.Second
	}

	if e.options.TLSConfig != nil {
		tlsConfig := e.options.TLSConfig
		if tlsConfig == nil {
			tlsConfig = &tls.Config{
				InsecureSkipVerify: true,
			}
		}

		config.TLS = tlsConfig
	}
	if e.options.Auth != nil {
		config.Username = e.options.Auth.Username
		config.Password = e.options.Auth.Password
	}
	var cAddrs []string

	for _, address := range e.options.Addrs {
		if len(address) == 0 {
			continue
		}
		addr, port, err := net.SplitHostPort(address)
		if ae, ok := err.(*net.AddrError); ok && ae.Err == "missing port in address" {
			port = "2379"
			addr = address
			cAddrs = append(cAddrs, net.JoinHostPort(addr, port))
		} else if err == nil {
			cAddrs = append(cAddrs, net.JoinHostPort(addr, port))
		}
	}

	// if we got addrs then we'll update
	if len(cAddrs) > 0 {
		config.Endpoints = cAddrs
	}

	// check if the endpoints have https://
	if config.TLS != nil {
		for i, ep := range config.Endpoints {
			if !strings.HasPrefix(ep, "https://") {
				config.Endpoints[i] = "https://" + ep
			}
		}
	}

	cli, err := clientv3.New(config)
	if err != nil {
		return nil, err
	}

	return cli, nil
}

// configure will setup the registry with new options
func configure(e *etcdRegistry, opts ...registry.Option) error {
	for _, o := range opts {
		o(e.options)
	}
	// setup the client
	cli, err := newClient(e)
	if err != nil {
		return err
	}

	if e.client != nil {
		e.client.Close()
	}

	// setup new client
	e.client = cli

	return nil
}

func (e *etcdRegistry) Options() registry.Options {
	return *e.options
}

func (e *etcdRegistry) Watch(opts ...registry.WatchOption) (registry.Watcher, error) {
	cli, err := newClient(e)
	if err != nil {
		return nil, err
	}
	return newEtcdWatcher(cli, e.options.Timeout, opts...)
}

func (e *etcdRegistry) String() string {
	return "etcd"
}
