package rpc

import (
	"fmt"
	"strings"
	"time"

	"github.com/micro-plat/hydra/rpc/balancer"
	"github.com/micro-plat/lib4go/concurrent/cmap"
	"github.com/micro-plat/lib4go/logger"
)

//Invoker RPC服务调用器，封装基于域及负载算法的RPC客户端
type Invoker struct {
	cache   cmap.ConcurrentMap
	address string
	opts    []ClientOption
	domain  string
	server  string
	lb      balancer.CustomerBalancer
	*invokerOption
}

type invokerOption struct {
	logger   *logger.Logger
	timerout time.Duration
	// balancerType int
	balancers map[string]BalancerMode
	servers   string
	// localPrefix  string
	tls map[string][]string
}

type BalancerMode struct {
	Mode  int
	Param string
}

const (
	//RoundRobin 轮询负载算法
	RoundRobin = iota + 1
	//LocalFirst 本地优先负载算法
	LocalFirst
)

//InvokerOption 客户端配置选项
type InvokerOption func(*invokerOption)

//WithInvokerLogger 设置日志记录器
func WithInvokerLogger(log *logger.Logger) InvokerOption {
	return func(o *invokerOption) {
		o.logger = log
	}
}

func WithBalancerMode(platName string, mode int, p string) InvokerOption {
	switch mode {
	case RoundRobin:
		return WithRoundRobin(platName)
	case LocalFirst:
		return WithLocalFirst(p, platName)
	default:
		return func(o *invokerOption) {
		}
	}
}

//WithRoundRobin 设置为轮询负载
func WithRoundRobin(platName ...string) InvokerOption {
	return func(o *invokerOption) {
		if len(platName) == 0 {
			o.balancers["*"] = BalancerMode{Mode: RoundRobin}
			return
		}
		for _, v := range platName {
			o.balancers[v] = BalancerMode{Mode: RoundRobin}
		}

	}
}

//WithLocalFirst 设置为本地优先负载
func WithLocalFirst(prefix string, platName ...string) InvokerOption {
	return func(o *invokerOption) {
		if prefix != "" {
			if len(platName) == 0 {
				o.balancers["*"] = BalancerMode{Mode: LocalFirst, Param: prefix}
				return
			}
			for _, v := range platName {
				o.balancers[v] = BalancerMode{Mode: LocalFirst, Param: prefix}
			}
		}
	}
}

//WithRPCTLS 设置TLS证书(pem,key)
func WithRPCTLS(platName string, tls []string) InvokerOption {
	return func(o *invokerOption) {
		if len(tls) == 2 {
			o.tls[platName] = tls
		}
	}
}

//NewInvoker 构建RPC服务调用器
//domain: 当前服务所在域
//server: 当前服务器名称
//addrss: 注册中心地址格式: zk://192.168.0.1166:2181或standalone://localhost
func NewInvoker(domain string, server string, address string, opts ...InvokerOption) (f *Invoker) {
	f = &Invoker{
		domain:  domain,
		server:  server,
		address: address,
		cache:   cmap.New(8),
		invokerOption: &invokerOption{
			balancers: map[string]BalancerMode{
				"*": BalancerMode{Mode: RoundRobin},
			},
			tls: make(map[string][]string),
		},
	}
	for _, opt := range opts {
		opt(f.invokerOption)
	}
	if f.invokerOption.logger == nil {
		f.invokerOption.logger = logger.GetSession("rpc.invoker", logger.CreateSession())
	}
	return
}

//RequestFailRetry 失败重试请求
func (r *Invoker) RequestFailRetry(service string, method string, header map[string]string, form map[string]interface{}, times int) (status int, result string, params map[string]string, err error) {
	for i := 0; i < times; i++ {
		status, result, params, err = r.Request(service, method, header, form, true)
		if err == nil || status < 500 {
			return
		}
	}
	return
}

//Request 使用RPC调用Request函数
func (r *Invoker) Request(service string, method string, header map[string]string, form map[string]interface{}, failFast bool) (status int, result string, params map[string]string, err error) {
	status = 500
	client, err := r.GetClient(service)
	if err != nil {
		return
	}
	_, rservice, d, s, _ := ResolvePath(service, r.domain, r.server)
	status, result, params, err = client.Request(rservice, method, header, form, failFast)
	if status != 200 || err != nil {
		if err != nil {
			err = fmt.Errorf("%s[@%s.%s]请求失败:%v(%d)", rservice, d, s, err, status)
		} else {
			err = fmt.Errorf("%s[@%s.%s]请求失败:%d)", rservice, d, s, status)
		}
	}
	return
}
func (r *Invoker) getBalancer(domain string) (int, string) {
	if b, ok := r.balancers[domain]; ok {
		return b.Mode, b.Param
	}
	if b, ok := r.balancers["*"]; ok {
		return b.Mode, b.Param
	}
	return RoundRobin, ""
}

//GetClient 获取RPC客户端
//addr 支持格式:
//order.request@merchant.hydra
//order.request,order.request@api.hydra
//order.request@api
func (r *Invoker) GetClient(addr string) (c *Client, err error) {
	isIP, rservice, domain, server, err := ResolvePath(addr, r.domain, r.server)
	if err != nil {
		return
	}
	//
	serviceKey := fmt.Sprintf("%s@%s.%s", rservice, server, domain)
	_, client, err := r.cache.SetIfAbsentCb(serviceKey, func(i ...interface{}) (interface{}, error) {
		plat := i[0].(string)
		server := i[1].(string)
		service := i[2].(string)
		isip := i[3].(bool)

		opts := make([]ClientOption, 0, 0)
		opts = append(opts, WithLogger(r.logger))

		//IP直接调用
		if isip {
			//设置安全证书
			switch len(r.tls[domain]) {
			case 2:
				opts = append(opts, WithTLS(r.tls[domain]))
			}
			return NewClient(server, opts...)
		}

		//非IP调用
		mode, p := r.getBalancer(domain)
		rs := balancer.NewResolver(plat, server, service, time.Second, p)
		servicePath := fmt.Sprintf(serviceRoot, strings.TrimPrefix(domain, "/"), server, service)
		//设置负载均衡算法
		switch mode {
		case RoundRobin:
			opts = append(opts, WithRoundRobinBalancer(rs, servicePath, time.Second, map[string]int{}))
		case LocalFirst:
			opts = append(opts, WithLocalFirstBalancer(rs, servicePath, p, map[string]int{}))
		default:
		}

		//设置安全证书
		switch len(r.tls[domain]) {
		case 2:
			opts = append(opts, WithTLS(r.tls[domain]))
		}
		return NewClient(r.address, opts...)
	}, domain, server, rservice, isIP)
	if err != nil {
		return
	}
	c = client.(*Client)
	return
}

//PreInit 预初始化服务器连接
func (r *Invoker) PreInit(services ...string) (err error) {
	for _, v := range services {
		_, err = r.GetClient(v)
		if err != nil {
			return
		}
	}
	return
}

//Close 关闭当前客户端与服务器的连接
func (r *Invoker) Close() {
	r.cache.RemoveIterCb(func(k string, v interface{}) bool {
		client := v.(*Client)
		client.Close()
		return true
	})
}
