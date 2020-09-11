package gray

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/hydra/pkgs/lua"
	"github.com/micro-plat/hydra/registry"
)

//Gray 灰度设置
type Gray struct {
	Disable           bool   `json:"disable,omitempty" toml:"disable,omitempty"`
	Script            string `json:"script" valid:"required" toml:"script,omitempty"`
	GetUpStreamMethod string
	Go2UpStreamMethod string
	cluster           conf.ICluster
}

//New 灰度设置
func New(script string) *Gray {
	return &Gray{
		GetUpStreamMethod: "getUpStream",
		Go2UpStreamMethod: "go2UpStream",
		Script:            script,
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

//checkServers 检查服务器信息
func (g *Gray) checkServers(c conf.IMainConf) error {
	vm, err := lua.New(g.Script, lua.WithMainFuncMode())
	if err != nil {
		return err
	}
	defer vm.Shutdown()

	rts, err := vm.CallByMethod(g.GetUpStreamMethod)
	if err != nil {
		return err
	}
	if len(rts) < 1 {
		return fmt.Errorf("%s至少包含一个返回参数", g.GetUpStreamMethod)
	}

	cluster, err := c.GetCluster(rts[0])
	if err != nil {
		return fmt.Errorf("获取指定的集群信息失败：%w", err)
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
	raw, err := cnf.GetSubConf(registry.Join("acl", "gray"))
	if err == conf.ErrNoSetting {
		return &Gray{Disable: true}
	}
	if err != nil {
		panic(fmt.Errorf("脚本加载失败 %w", err))
	}
	if err != nil && err != conf.ErrNoSetting {
		panic(fmt.Errorf("acl.gray配置有误:%v", err))
	}
	gray := New(string(raw.GetRaw()))
	if err := gray.checkServers(cnf); err != nil {
		panic(err)
	}

	return gray
}
