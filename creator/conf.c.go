package creator

import (
	"bytes"
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/micro-plat/hydra/conf/server/api"
	"github.com/micro-plat/hydra/conf/server/cron"
	"github.com/micro-plat/hydra/conf/server/mqc"
	"github.com/micro-plat/hydra/conf/server/rpc"
	"github.com/micro-plat/hydra/conf/server/static"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/hydra/hydra/servers"
	"github.com/micro-plat/hydra/services"
)

//IConf 配置注册管理
type IConf interface {

	//Var 参数配置
	Vars() vars

	//API api服务器配置
	API(address string, opts ...api.Option) *httpBuilder

	//GetAPI() 获取API服务器配置
	GetAPI() *httpBuilder

	//Web web服务器配置
	Web(address string, opts ...api.Option) *httpBuilder

	//GetWeb() 获取Web服务器配置
	GetWeb() *httpBuilder

	//WS ws服务器配置
	WS(address string, opts ...api.Option) *httpBuilder

	//GetWS() 获取ws服务器配置
	GetWS() *httpBuilder

	//RPC rpc服务器配置
	RPC(address string, opts ...rpc.Option) *rpcBuilder

	//GetRPC() 获取rpc服务器配置
	GetRPC() *rpcBuilder

	//Custom 自定义服务器配置
	Custom(tp string, s ...interface{}) CustomerBuilder

	//CRON 构建cron服务器配置
	CRON(opts ...cron.Option) *cronBuilder

	//GetCRON 获取CRON服务器配置
	GetCRON() *cronBuilder

	//CRON 构建mqc服务器配置
	MQC(addr string, opts ...mqc.Option) *mqcBuilder

	//GetMQC 获取MQC服务器配置
	GetMQC() *mqcBuilder

	//Pub 发布服务
	Pub(platName string, systemName string, clusterName string, registryAddr string, cover bool) error

	//Load 加载所有配置
	Load() error
}

//ServerMainNodeName 服务主节点名称
const ServerMainNodeName = "main"

//Conf 配置服务
var Conf = New()

//New 构建新的配置
func New() *conf {
	return &conf{
		data:         make(map[string]iCustomerBuilder),
		vars:         make(map[string]map[string]interface{}),
		routerLoader: services.GetRouter,
	}
}

//NewByLoader 设置路由加载器
func NewByLoader(routerLoader func(string) *services.ORouter) *conf {
	return &conf{
		data:         make(map[string]iCustomerBuilder),
		vars:         make(map[string]map[string]interface{}),
		routerLoader: routerLoader,
	}
}

type conf struct {
	data         map[string]iCustomerBuilder
	vars         map[string]map[string]interface{}
	routerLoader func(string) *services.ORouter
}

//Load 加载所有配置
func (c *conf) Load() error {

	types := servers.GetServerTypes()
	for _, t := range types {
		_, ok := c.data[t]
		if !ok {
			switch t {
			case global.API:
				c.data[global.API] = c.GetAPI()
			case global.Web:
				c.data[global.Web] = c.GetWeb()
			case global.WS:
				c.data[global.WS] = c.GetWS()
			case global.RPC:
				c.data[global.RPC] = c.GetRPC()
			case global.CRON:
				c.data[global.CRON] = c.GetCRON()
			case global.MQC:
				c.data[global.MQC] = c.GetMQC()
			default:
				c.data[t] = newCustomerBuilder()
			}
		}
		c.data[t].Load()

	}
	//添加其它服务器
	return nil
}

//API api服务器配置
func (c *conf) API(address string, opts ...api.Option) *httpBuilder {
	api := newHTTP(global.API, address, c.routerLoader, opts...)
	c.data[global.API] = api
	return api
}

//GetAPI 获取当前已配置的api服务器
func (c *conf) GetAPI() *httpBuilder {
	if api, ok := c.data[global.API]; ok {
		return api.(*httpBuilder)
	}

	return c.API(api.DefaultAPIAddress)
}

//Web web服务器配置
func (c *conf) Web(address string, opts ...api.Option) *httpBuilder {
	web := newHTTP(global.Web, address, c.routerLoader, opts...)
	web.Static(static.WithArchive(global.AppName))
	c.data[global.Web] = web
	return web
}

//GetWeb 获取当前已配置的web服务器
func (c *conf) GetWeb() *httpBuilder {
	if web, ok := c.data[global.Web]; ok {
		return web.(*httpBuilder)
	}
	return c.Web(api.DefaultWEBAddress)
}

//Web web服务器配置
func (c *conf) WS(address string, opts ...api.Option) *httpBuilder {
	ws := newHTTP(global.WS, address, c.routerLoader, opts...)
	ws.Static(static.WithArchive(global.AppName))
	c.data[global.WS] = ws
	return ws
}

//GetWeb 获取当前已配置的web服务器
func (c *conf) GetWS() *httpBuilder {
	if ws, ok := c.data[global.WS]; ok {
		return ws.(*httpBuilder)
	}
	return c.WS(api.DefaultWSAddress)
}

//RPC rpc服务器配置
func (c *conf) RPC(address string, opts ...rpc.Option) *rpcBuilder {
	rpc := newRPC(address, c.routerLoader, opts...)
	c.data[global.RPC] = rpc
	return rpc
}

//GetRPC 获取当前已配置的rpc服务器
func (c *conf) GetRPC() *rpcBuilder {
	if rpc, ok := c.data[global.RPC]; ok {
		return rpc.(*rpcBuilder)
	}
	return c.RPC(rpc.DefaultRPCAddress)
}

//CRON cron服务器配置
func (c *conf) CRON(opts ...cron.Option) *cronBuilder {
	cron := newCron(opts...)
	c.data[global.CRON] = cron
	return cron
}

//GetWeb 获取当前已配置的web服务器
func (c *conf) GetCRON() *cronBuilder {
	if cron, ok := c.data[global.CRON]; ok {
		return cron.(*cronBuilder)
	}
	return c.CRON()
}

//MQC mqc服务器配置
func (c *conf) MQC(addr string, opts ...mqc.Option) *mqcBuilder {
	mqc := newMQC(addr, opts...)
	c.data[global.MQC] = mqc
	return mqc
}

//GetMQC 获取当前已配置的mqc服务器
func (c *conf) GetMQC() *mqcBuilder {
	if mqc, ok := c.data[global.MQC]; ok {
		return mqc.(*mqcBuilder)
	}
	panic("未指定mqc服务器配置")
}

//Vars 平台变量配置
func (c *conf) Vars() vars {
	return c.vars
}

//Vars 平台变量配置
func (c *conf) GetVar(tp, name string) (val interface{}, ok bool) {
	tpv, ok := c.vars[tp]
	if !ok {
		return
	}
	val, ok = tpv[name]
	if !ok {
		return
	}
	return
}

//Custom 用户自定义配置服务
func (c *conf) Custom(tp string, s ...interface{}) CustomerBuilder {
	if _, ok := c.data[tp]; ok {
		panic(fmt.Sprintf("不能重复注册%s", tp))
	}
	customer := newCustomerBuilder(s...)
	c.data[tp] = customer
	return customer
}

//Encode 将当前配置序列化为toml格式
func (c *conf) Encode() (string, error) {
	var buffer bytes.Buffer
	encoder := toml.NewEncoder(&buffer)
	err := encoder.Encode(&c.data)
	return buffer.String(), err
}

//Encode2File 将当前配置内容保存到文件中
func (c *conf) Encode2File(path string, cover bool) error {
	if !cover {
		if _, err := os.Stat(path); err == nil || os.IsExist(err) {
			return fmt.Errorf("配置文件已存在 %s", path)
		}
	}
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return fmt.Errorf("无法打开文件:%s %w", path, err)
	}
	encoder := toml.NewEncoder(f)
	err = encoder.Encode(&c.data)
	if err != nil {
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return nil
}

//Decode 从配置文件中读取配置信息
func (c *conf) Decode(f string) error {
	_, err := toml.DecodeFile(f, &c.data)
	return err
}

//Decode 从配置文件中读取配置信息
func (c *conf) Decode1(f string) error {

	_, err := toml.DecodeFile(f, &c.data)
	return err
}
