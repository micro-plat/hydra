package gray

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/hydra/registry"
)

//Gray 灰度设置
type Gray struct {
	Disable   bool   `json:"disable,omitempty" toml:"disable,omitempty"`
	Filter    string `json:"filter" valid:"required" toml:"filter,omitempty"`
	UPCluster string `json:"upcluster" valid:"required" toml:"upcluster,omitempty"`
	conf      conf.IMainConf
	cluster   conf.ICluster
}

//New 灰度设置
func New(filter string, upcluster string) *Gray {
	return &Gray{
		Filter:    filter,
		UPCluster: upcluster,
	}
}

//Allow 当前服务
func (g *Gray) Allow() bool {
	return g.cluster.GetType() == global.API || g.cluster.GetType() == global.RPC
}

//Next 获取下一个可用的上游地址
func (g *Gray) Next() (u *url.URL, err error) {
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

func (g *Gray) checkServers(c conf.IMainConf) error {
	cluster, err := c.GetCluster(g.UPCluster)
	if err != nil {
		return err
	}
	g.cluster = cluster
	return nil
}

type ConfHandler func(cnf conf.IMainConf) *Gray

func (h ConfHandler) Handle(cnf conf.IMainConf) interface{} {
	return h(cnf)
}

//GetConf 获取BlackList
func GetConf(cnf conf.IMainConf) *Gray {
	gray := Gray{}
	_, err := cnf.GetSubObject(registry.Join("acl", "gray"), &gray)
	if err == conf.ErrNoSetting {
		return &Gray{Disable: true}
	}

	if err != nil && err != conf.ErrNoSetting {
		panic(fmt.Errorf("acl.gray配置有误:%v", err))
	}
	if err := gray.checkServers(cnf); err != nil {
		panic(err)
	}

	return &gray
}
