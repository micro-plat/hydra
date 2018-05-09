package rpc

import (
	"fmt"
	"time"

	"github.com/micro-plat/hydra/rpc/balancer"
	"github.com/micro-plat/hydra/servers/rpc/pb"
	"github.com/micro-plat/lib4go/jsons"
	"github.com/micro-plat/lib4go/logger"
	"github.com/micro-plat/lib4go/types"

	"errors"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

//Client rpc client, 用于构建基础的RPC调用,并提供基于服务器的限流工具，轮询、本地优先等多种负载算法
type Client struct {
	address string //RPC Server Address 或 registry target
	conn    *grpc.ClientConn
	*clientOption
	client        pb.RPCClient
	hasRunChecker bool
	IsConnect     bool
	isClose       bool
}

type clientOption struct {
	connectionTimeout time.Duration
	log               *logger.Logger
	balancer          balancer.CustomerBalancer
	resolver          balancer.ServiceResolver
	service           string
}

//ClientOption 客户端配置选项
type ClientOption func(*clientOption)

//WithConnectionTimeout 配置网络连接超时时长
func WithConnectionTimeout(t time.Duration) ClientOption {
	return func(o *clientOption) {
		o.connectionTimeout = t
	}
}

//WithLogger 配置日志记录器
func WithLogger(log *logger.Logger) ClientOption {
	return func(o *clientOption) {
		o.log = log
	}
}

//WithRoundRobinBalancer 配置为轮询负载均衡器
func WithRoundRobinBalancer(r balancer.ServiceResolver, service string, timeout time.Duration, limit map[string]int) ClientOption {
	return func(o *clientOption) {
		o.service = service
		o.resolver = r
		o.balancer = balancer.RoundRobin(service, r, limit, o.log)
	}
}

//WithLocalFirstBalancer 配置为本地优先负载均衡器
func WithLocalFirstBalancer(r balancer.ServiceResolver, service string, local string, limit map[string]int) ClientOption {
	return func(o *clientOption) {
		o.service = service
		o.resolver = r
		o.balancer = balancer.LocalFirst(service, local, r, limit)
	}
}

//WithBalancer 设置负载均衡器
func WithBalancer(service string, lb balancer.CustomerBalancer) ClientOption {
	return func(o *clientOption) {
		o.service = service
		o.balancer = lb
	}
}

//NewClient 创建RPC客户端,地址是远程RPC服务器地址或注册中心地址
func NewClient(address string, opts ...ClientOption) (*Client, error) {
	client := &Client{address: address, clientOption: &clientOption{connectionTimeout: time.Second * 3}}
	for _, opt := range opts {
		opt(client.clientOption)
	}
	if client.log == nil {
		client.log = logger.GetSession("rpc.client", logger.CreateSession())
	}
	//grpclog.SetLogger(client.log)
	err := client.connect()
	if err != nil {
		err = fmt.Errorf("rpc.client连接到服务器失败:%s(%v)(err:%v)", address, client.connectionTimeout, err)
		return nil, err
	}
	time.Sleep(time.Second)
	return client, err
}

//Connect 连接到RPC服务器，如果当前无法连接系统会定时自动重连
func (c *Client) connect() (err error) {
	if c.IsConnect {
		return nil
	}
	if c.balancer == nil {
		c.conn, err = grpc.Dial(c.address, grpc.WithInsecure(), grpc.WithTimeout(c.connectionTimeout))
	} else {
		ctx, _ := context.WithTimeout(context.Background(), c.connectionTimeout)
		c.conn, err = grpc.DialContext(ctx, c.address, grpc.WithInsecure(), grpc.WithBalancer(c.balancer))
	}
	if err != nil {
		return
	}
	c.client = pb.NewRPCClient(c.conn)
	return nil
}

//Request 发送Request请求
func (c *Client) Request(service string, method string, header map[string]string, form map[string]interface{}, failFast bool) (status int, result string, param map[string]string, err error) {
	h, err := jsons.Marshal(header)
	if err != nil {
		return
	}
	if len(h) == 0 {
		h = []byte("{}")
	}
	f, err := jsons.Marshal(form)
	if err != nil {
		return
	}
	if len(f) == 0 {
		h = []byte("{}")
	}
	response, err := c.client.Request(context.Background(),
		&pb.RequestContext{
			Method:  method,
			Service: service,
			Header:  string(h),
			Input:   string(f),
		},
		grpc.FailFast(failFast))
	if err != nil {
		status = 500
		return
	}

	if response.Header != "" {
		mh, err := jsons.Unmarshal([]byte(response.Header))
		if err != nil {
			return 400, "", nil, err
		}
		param, _ = types.ToStringMap(mh)
	}

	status = int(response.Status)
	result = response.GetResult()
	return
}

//UpdateLimiter 修改服务器限流规则
func (c *Client) UpdateLimiter(limit map[string]int) error {
	if c.balancer != nil {
		c.balancer.UpdateLimiter(limit)
		return nil
	}
	return errors.New("rpc.client.未指定balancer")
}

//Close 关闭RPC客户端连接
func (c *Client) Close() {
	c.isClose = true
	if c.resolver != nil {
		c.resolver.Close()
	}
	if c.conn != nil {
		c.conn.Close()
	}
}
