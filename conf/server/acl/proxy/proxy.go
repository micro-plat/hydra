package proxy

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/lib4go/tgo"
)

const (
	//ParNodeName proxy配置父节点名
	ParNodeName = "acl"
	//SubNodeName proxy配置子节点名
	SubNodeName = "proxy"
)

//Proxy 代理设置
type Proxy struct {
	Disable bool          `json:"disable,omitempty" toml:"disable,omitempty"`
	Script  string        `json:"-"`
	cluster conf.ICluster `json:"-"`
	c       conf.IServerConf
	tengo   *tgo.VM
}

//New 代理设置(该方法只用在注册中心安装时调用,如果要使用对象方法请通过GetConf获取对象)
func New(opts ...Option) *Proxy {
	r := &Proxy{
		Disable: false,
	}
	for _, f := range opts {
		f(r)
	}
	return r
}

//Allow 当前服务是否允许使用代理
func (g *Proxy) Allow() bool {
	if g.cluster == nil {
		return false
	}
	return g.cluster.GetServerType() == global.API || g.cluster.GetServerType() == global.Web
}

//Check 检查当前是否需要转到上游服务器处理
func (g *Proxy) Check(funcs map[string]interface{}, i interface{}) (bool, error) {
	result, err := g.tengo.Run()
	if err != nil {
		return false, err
	}
	result.GetBool("")
	return false, nil
}

//Next 获取下一个可用的上游地址
func (g *Proxy) Next() (u *url.URL, err error) {
	if g.cluster == nil {
		return nil, errors.New("当前配置不可用")
	}
	node, ok := g.cluster.Next()
	if !ok {
		return nil, fmt.Errorf("无法获取到集群的下一个服务器")
	}
	path := fmt.Sprintf("http://%s:%s", node.GetHost(), node.GetPort())
	url, err := url.Parse(path)
	if err != nil {
		return nil, fmt.Errorf("集群的服务器地址不合法:%s %s", path, err)
	}
	return url, nil
}

func (g *Proxy) checkServers() (err error) {
	g.tengo, err = tgo.New(g.Script, tgo.WithModule(global.GetTGOModules()...))
	if err != nil {
		return err
	}
	result, err := g.tengo.Run()
	upstream := result.GetString("upstream")
	cluster, err := g.c.GetCluster(upstream)
	if err != nil {
		return err
	}
	g.cluster = cluster
	return err
}

//GetConf 获取Proxy
func GetConf(cnf conf.IServerConf) (*Proxy, error) {
	script, err := cnf.GetSubConf(registry.Join(ParNodeName, SubNodeName))
	if err == conf.ErrNoSetting {
		return &Proxy{Disable: true}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("acl.proxy配置有误:%v", err)
	}
	proxy := New(WithScript(string(script.GetRaw())))
	proxy.c = cnf
	if err := proxy.checkServers(); err != nil {
		return nil, fmt.Errorf("acl.proxy服务检查错误:%v", err)
	}
	return proxy, nil
}
