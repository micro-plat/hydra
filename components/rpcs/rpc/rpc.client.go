package rpc

import (
	"fmt"
	"net/http"
	"time"

	"github.com/micro-plat/lib4go/types"

	"github.com/micro-plat/hydra/components/rpcs/balancer"
	"github.com/micro-plat/hydra/components/rpcs/rpc/pb"
	rpcconf "github.com/micro-plat/hydra/conf/vars/rpc"
	"github.com/micro-plat/lib4go/logger"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
)

//Client rpc client, 用于构建基础的RPC调用,并提供基于服务器的限流工具，轮询、本地优先等多种负载算法
type Client struct {
	address string //RPC Server Address 或 registry target
	conn    *grpc.ClientConn
	plat    string
	service string
	log     *logger.Logger
	*rpcconf.RPCConf
	client        pb.RPCClient
	hasRunChecker bool
	IsConnect     bool
	isClose       bool
}

//NewClient .
func NewClient(address, plat, service string, opts ...rpcconf.Option) (*Client, error) {
	conf := rpcconf.New(opts...)
	return NewClientByConf(address, plat, service, conf)
}

//NewClientByConf 创建RPC客户端,地址是远程RPC服务器地址或注册中心地址
func NewClientByConf(address, plat, service string, conf *rpcconf.RPCConf) (*Client, error) {
	client := &Client{address: address, plat: plat, service: service}
	client.RPCConf = conf

	if client.log == nil {
		client.log = logger.GetSession(types.GetStringByIndex([]string{client.Log}, 0, "rpc.client"), logger.CreateSession())
	}
	err := client.connect()
	if err != nil {
		err = fmt.Errorf("rpc.client连接到服务器失败:%s(%v)(err:%v)", address, client.ConntTimeout, err)
		return nil, err
	}
	time.Sleep(time.Second)
	return client, err
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
	if err != nil {
		return NewResponseByStatus(http.StatusInternalServerError, err), err
	}
	return NewResponse(int(response.Status), response.GetHeader(), response.GetResult()), err
}

//Close 关闭RPC客户端连接
func (c *Client) Close() {
	c.isClose = true
	if c.conn != nil {
		c.conn.Close()
	}
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

	balanc := balancer.RoundRobin
	if c.LocalFirst {
		balanc = balancer.LocalFirst
	}

	var rb resolver.Builder
	//兼容直接传服务器ip来进行访问
	if len(c.plat) > 0 {
		rb = balancer.NewResolverBuilder(c.address, c.plat, c.service, c.SortPrefix)
	}

	ctx, _ := context.WithTimeout(context.Background(), time.Duration(c.ConntTimeout)*time.Second)
	c.conn, err = grpc.DialContext(ctx,
		c.address+"/mockrpc",
		grpc.WithInsecure(),
		grpc.WithBalancerName(balanc),
		grpc.WithResolvers(rb))

	if err != nil {
		return
	}
	c.client = pb.NewRPCClient(c.conn)
	return nil
}
