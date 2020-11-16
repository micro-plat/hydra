package proxy

import (
	"fmt"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/lib4go/tgo"
)

const (
	//ParNodeName proxy配置父节点名
	ParNodeName = "acl"

	//SubNodeName proxy配置子节点名
	SubNodeName = "proxy"

	//脚本中的上游集群参数名称
	upclusterName = "upcluster"
)

//Proxy 代理设置
type Proxy struct {
	Disable bool `json:""-`
	c       conf.IServerConf
	tengo   *tgo.VM
}

//New 代理设置(该方法只用在注册中心安装时调用,如果要使用对象方法请通过GetConf获取对象)
func New(opts ...Option) *Proxy {
	r := &Proxy{}
	for _, f := range opts {
		f(r)
	}
	return r
}

//Check 检查当前是否需要转到上游服务器处理
func (g *Proxy) Check() (*UpCluster, bool, error) {

	//执行脚本，检查当前请求是否需要转到上游服务器
	result, err := g.tengo.Run()
	if err != nil {
		return nil, false, err
	}

	//获取脚本执行结果
	upstream := result.GetString(upclusterName)
	if upstream == "" || upstream == g.c.GetClusterName() {
		return nil, false, nil
	}

	//保存到缓存，或从缓存获取上游集群信息
	_, cluster, err := clusters.SetIfAbsentCb(upstream, func(value ...interface{}) (interface{}, error) {
		up, err := g.c.GetCluster(upstream)
		if err != nil {
			return nil, err
		}
		return &UpCluster{c: up, name: upstream}, nil
	})
	if err != nil {
		return nil, false, err
	}
	return cluster.(*UpCluster), true, nil

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
	proxy := New()
	proxy.c = cnf
	proxy.tengo, err = tgo.New(string(script.GetRaw()))
	if err != nil {
		return nil, fmt.Errorf("acl.proxy脚本错误:%v", err)
	}
	return proxy, nil
}
