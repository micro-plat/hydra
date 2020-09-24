package rpc

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/micro-plat/hydra/components/rpcs/balancer"
	"github.com/micro-plat/hydra/components/rpcs/rpc/pb"
	"github.com/micro-plat/lib4go/logger"
	"github.com/micro-plat/lib4go/net"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

//IRequest RPC请求
type IRequest interface {
	Request(ctx context.Context, service string, form map[string]interface{}, opts ...RequestOption) (res *Response, err error)
}

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
	balancer          string
	localIP           string
	//resolver          balancer.ServiceResolver
	plat       string
	service    string
	sortPrefix string
	tls        []string
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

//WithTLS 设置TLS证书(pem,key)
func WithTLS(tls []string) ClientOption {
	return func(o *clientOption) {
		if len(tls) == 2 {
			o.tls = tls
		}
	}
}

type requestOption struct {
	name     string
	service  string
	headers  map[string]string
	failFast bool
	method   string
}

//RequestOption 客户端配置选项
type RequestOption func(*requestOption)

func newOption() *requestOption {
	return &requestOption{
		headers:  map[string]string{},
		failFast: true,
		method:   "GET",
	}
}

//WithHeader 请求头信息
func WithHeader(k string, v string) RequestOption {
	return func(o *requestOption) {
		o.headers[k] = v
	}
}

//WithHeaders 设置请求头
func WithHeaders(p map[string][]string) RequestOption {
	return func(o *requestOption) {
		for k, v := range p {
			o.headers[k] = v[0]
		}
	}
}

//WithHost 设置当前机器IP
func WithHost(s ...string) RequestOption {
	return func(o *requestOption) {
		if len(s) > 0 {
			o.headers["Host"] = strings.Join(s, ",")
		} else {
			o.headers["Host"] = net.GetLocalIPAddress()
		}
	}
}

//WithXRequestID 设置请求编号
func WithXRequestID(s string) RequestOption {
	return func(o *requestOption) {
		o.headers["X-Request-Id"] = s
	}
}

//WithDelay 设置请求延迟时长
func WithDelay(s time.Duration) RequestOption {
	return func(o *requestOption) {
		o.headers["X-Add-Delay"] = fmt.Sprint(s)
	}
}

//WithFailFast 快速失败
func WithFailFast(b bool) RequestOption {
	return func(o *requestOption) {
		o.failFast = b
	}
}

//WithMethod 设置请求方法
func WithMethod(m string) RequestOption {
	return func(o *requestOption) {
		o.method = strings.ToUpper(m)
	}
}

//WithContentType 设置请求类型
func WithContentType(m string) RequestOption {
	return func(o *requestOption) {
		o.headers["Content-Type"] = m
	}
}

//WithOperationName 设置请求延迟时长
func WithOperationName(name string) RequestOption {
	return func(o *requestOption) {
		o.name = name
	}
}

func (r *requestOption) getData(v interface{}) ([]byte, error) {
	buff, err := json.Marshal(&v)
	if err != nil {
		return nil, err
	}
	if len(buff) == 0 {
		return []byte("{}"), nil
	}
	return buff, nil
}

//WithRoundRobinBalancer 配置为轮询负载均衡器
func WithRoundRobinBalancer(plat, service string) ClientOption {
	return func(o *clientOption) {
		o.plat = plat
		o.service = service
		o.sortPrefix = ""
		//o.resolver = balancer.NewResolver(plat, service, "")
		//o.balancer = balancer.RoundRobin(service, o.resolver, o.log)
		o.balancer = balancer.RoundRobin
	}
}

//WithLocalFirstBalancer 配置为本地优先负载均衡器
func WithLocalFirstBalancer(plat, service string, local string) ClientOption {
	return func(o *clientOption) {
		o.plat = plat
		o.service = service
		o.sortPrefix = local
		//o.resolver = balancer.NewResolver(plat, service, local)
		//o.balancer = balancer.LocalFirst(service, local, o.resolver)
		o.balancer = balancer.LocalFirst
		o.localIP = local
	}
}

//WithBalancer 设置负载均衡器
func WithBalancer(plat, service string, lbname string) ClientOption {
	return func(o *clientOption) {
		o.plat = plat
		o.service = service
		o.balancer = lbname
	}
}

// //WithBalancer 设置负载均衡器
// func WithBalancer(service string, lb balancer.CustomerBalancer) ClientOption {
// 	return func(o *clientOption) {
// 		o.service = service
// 		o.balancer = lb
// 	}
// }

//NewClient 创建RPC客户端,地址是远程RPC服务器地址或注册中心地址
func NewClient(address string, opts ...ClientOption) (*Client, error) {
	client := &Client{address: address, clientOption: &clientOption{connectionTimeout: time.Second * 3}}
	for _, opt := range opts {
		opt(client.clientOption)
	}
	if client.log == nil {
		client.log = logger.GetSession("rpc.client", logger.CreateSession())
	}
	err := client.connect()
	if err != nil {
		err = fmt.Errorf("rpc.client连接到服务器失败:%s(%v)(err:%v)", address, client.connectionTimeout, err)
		return nil, err
	}
	time.Sleep(time.Second)
	return client, err
}

//Connect 连接到RPC服务器，如果当前无法连接系统会定时自动重连
//未使用压缩，由于传输数据默认限制为4M(已修改为20M)压缩后会影响系统并发能力
// grpc.WithDefaultCallOptions(grpc.UseCompressor(Snappy)),
// grpc.WithDecompressor(grpc.NewGZIPDecompressor()),
// grpc.WithCompressor(grpc.NewGZIPCompressor()),
func (c *Client) connect() (err error) {
	if c.IsConnect {
		return nil
	}
	if c.balancer == "" {
		c.balancer = balancer.RoundRobin
	}
	fmt.Println("client.connect.2", c.balancer, c.address)

	rb := balancer.NewResolverBuilder(c.address, c.plat, c.service, c.sortPrefix)

	ctx, _ := context.WithTimeout(context.Background(), c.connectionTimeout)
	c.conn, err = grpc.DialContext(ctx,
		c.address+"/mockrpc",
		grpc.WithInsecure(),
		grpc.WithBalancerName(c.balancer),
		grpc.WithResolvers(rb))

	if err != nil {
		fmt.Println("client.connect.3", err)
		return
	}
	c.client = pb.NewRPCClient(c.conn)
	return nil
}

//Request 发送Request请求
func (c *Client) Request(ctx context.Context, service string, form map[string]interface{}, opts ...RequestOption) (res *Response, err error) {
	//处理可选参数
	o := newOption()
	for _, opt := range opts {
		opt(o)
	}
	o.service = service
	response, err := c.clientRequest(ctx, o, form)
	fmt.Println("Client.Request.1")
	if err != nil {
		fmt.Println("Client.Request.2", err)
		return NewResponseByStatus(http.StatusInternalServerError, err), err
	}
	fmt.Println("Client.Request.3")
	return NewResponse(int(response.Status), response.GetHeader(), response.GetResult()), err
}

//Close 关闭RPC客户端连接
func (c *Client) Close() {
	c.isClose = true
	// if c.resolver != nil {
	// 	c.resolver.Close()
	// }
	if c.conn != nil {
		c.conn.Close()
	}
}
