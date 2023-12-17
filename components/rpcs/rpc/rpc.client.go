package rpc

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/micro-plat/hydra/pkgs"
	"github.com/micro-plat/lib4go/types"

	"github.com/micro-plat/hydra/components/rpcs/balancer"
	"github.com/micro-plat/hydra/components/rpcs/rpc/pb"
	rpcconf "github.com/micro-plat/hydra/conf/vars/rpc"
	"github.com/micro-plat/lib4go/logger"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Client rpc client, 用于构建基础的RPC调用,并提供基于服务器的限流工具，轮询、本地优先等多种负载算法
type Client struct {
	address string //RPC Server Address 或 registry target
	conn    *grpc.ClientConn
	plat    string
	service string
	log     *logger.Logger
	*rpcconf.RPCConf
	client          pb.RPCClient
	balancerBuilder *balancer.ResolverBuilder
	hasRunChecker   bool
	IsConnect       bool
	isClose         bool
}

// NewClient .
func NewClient(address, plat, service string, opts ...rpcconf.Option) (*Client, error) {
	conf := rpcconf.New(opts...)
	return NewClientByConf(address, plat, service, conf)
}

// NewClientByConf 创建RPC客户端,地址是远程RPC服务器地址或注册中心地址
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

// Request 发送Request请求
func (c *Client) Request(ctx context.Context, service string, form map[string]interface{}, opts ...RequestOption) (res *pkgs.Rspns, err error) {
	//处理可选参数
	buff, err := json.Marshal(form)
	if err != nil {
		return nil, err
	}
	if len(buff) == 0 {
		buff = []byte("{}")
	}

	return c.RequestByString(ctx, service, string(buff), opts...)
}

// RequestByString 发送Request请求
func (c *Client) RequestByString(ctx context.Context, service string, form string, opts ...RequestOption) (res *pkgs.Rspns, err error) {
	//处理可选参数
	o := newOption()
	for _, opt := range opts {
		opt(o)
	}

	for k := range o.headers {
		if strings.EqualFold("Content-Type", k) {
			o.headers[k] = "application/json;charset=utf-8"
			break
		}
	}

	o.service = service
	response, err := c.clientRequest(ctx, o, form)
	if err != nil {
		return pkgs.NewRspns(err), err
	}
	return pkgs.NewRspnsByHD(int(response.Status), response.GetHeader(), response.GetResult()), err
}

// Close 关闭RPC客户端连接
func (c *Client) Close() {
	c.isClose = true
	if c.conn != nil {
		c.conn.Close()
	}
	c.balancerBuilder.Close()
}

// Connect 连接到RPC服务器，如果当前无法连接系统会定时自动重连
// 未使用压缩，由于传输数据默认限制为4M(已修改为20M)压缩后会影响系统并发能力
// grpc.WithDefaultCallOptions(grpc.UseCompressor(Snappy)),
// grpc.WithDecompressor(grpc.NewGZIPDecompressor()),
// grpc.WithCompressor(grpc.NewGZIPCompressor()),
func (c *Client) connect() (err error) {
	if c.IsConnect {
		return nil
	}

	if c.Balancer == "" {
		c.Balancer = rpcconf.LocalFirst
	}

	c.balancerBuilder, err = balancer.NewResolverBuilder(c.address, c.plat, c.service, c.SortPrefix)
	if err != nil {
		return
	}

	ctx, _ := context.WithTimeout(context.Background(), time.Duration(c.ConntTimeout)*time.Second)
	c.conn, err = grpc.DialContext(ctx,
		c.address+"/rpcsrv",
		grpc.WithTransportCredentials(insecure.NewCredentials()), //2023/12/17
		grpc.WithDefaultServiceConfig(
			`{"loadBalancingPolicy":"round_robin"}`,
		),
		grpc.WithResolvers(c.balancerBuilder))

	if err != nil {
		return
	}
	c.client = pb.NewRPCClient(c.conn)
	return nil
}
