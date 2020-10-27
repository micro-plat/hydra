package gray

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/hydra/registry"
)

//Gray 灰度设置
type Gray struct {
	Disable   bool   `json:"disable,omitempty" toml:"disable,omitempty"`
	Filter    string `json:"filter" valid:"required" toml:"filter,omitempty"`
	UPCluster string `json:"upcluster" valid:"required" toml:"upcluster,omitempty"`
	//conf      conf.IMainConf
	cluster conf.ICluster
}

//New 灰度设置(该方法只用在注册中心安装时调用,如果要使用对象方法请通过GetConf获取对象)
func New(opts ...Option) *Gray {
	r := &Gray{
		Disable: false,
	}
	for _, f := range opts {
		f(r)
	}
	return r
}

//Allow 当前服务是否允许使用灰度
func (g *Gray) Allow() bool {
	if g.cluster == nil {
		return false
	}
	return g.cluster.GetType() == global.API || g.cluster.GetType() == global.RPC
}

//Check 检查当前是否需要转到上游服务器处理
func (g *Gray) Check(funcs map[string]interface{}, i interface{}) (bool, error) {
	r, err := conf.TmpltTranslate(g.Filter, g.Filter, funcs, i)
	if err != nil {
		return true, fmt.Errorf("%s 过滤器转换出错 %w", g.Filter, err)
	}
	return strings.EqualFold(r, "true"), nil
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

//GetConf 获取Gray
func GetConf(cnf conf.IMainConf) (*Gray, error) {
	gray := &Gray{}
	_, err := cnf.GetSubObject(registry.Join("acl", "gray"), gray)
	if err == conf.ErrNoSetting {
		return &Gray{Disable: true}, nil
	}

	if err != nil && err != conf.ErrNoSetting {
		return nil, fmt.Errorf("acl.gray配置有误:%v", err)
	}

	if b, err := govalidator.ValidateStruct(gray); !b {
		return nil, fmt.Errorf("acl.gray配置数据有误:%v %+v", err, gray)
	}

	if err := gray.checkServers(cnf); err != nil {
		return nil, fmt.Errorf("acl.gray服务检查错误:%v", err)
	}

	return gray, nil
}
