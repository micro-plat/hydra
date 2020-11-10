package whitelist

import (
	"fmt"

	"github.com/asaskevich/govalidator"
	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/registry"
)

const (
	//ParNodeName 白名单配置父节点名
	ParNodeName = "acl"
	//SubNodeName 白名单配置子节点名
	SubNodeName = "white.list"
)

//IPList ip列表
type IPList struct {
	Requests []string `json:"requests,omitempty" valid:"ascii,required" toml:"requests,omitempty"`
	IPS      []string `json:"ips,omitempty" valid:"ascii,required" toml:"ips,omitempty"`
	ipm      *conf.PathMatch
	rqm      *conf.PathMatch
}

//WhiteList 白名单配置
type WhiteList struct {
	Disable bool      `json:"disable,omitempty" toml:"disable,omitempty"`
	IPS     []*IPList `json:"whiteList,omitempty" toml:"whiteList,omitempty"`
}

//New 创建白名单规则服务
func New(opts ...Option) *WhiteList {
	f := &WhiteList{IPS: make([]*IPList, 0, 1)}
	for idx := range opts {
		opts[idx](f)
	}
	return f
}

//IsAllow 验证当前请求是否在白名单中
func (w *WhiteList) IsAllow(path string, ip string) bool {
	for _, cur := range w.IPS {
		if ok, _ := cur.rqm.Match(path); ok {
			ok, _ := cur.ipm.Match(ip, ".")
			return ok
		}
	}
	return true
}

//GetConf 获取WhiteList
func GetConf(cnf conf.IServerConf) (*WhiteList, error) {
	ip := WhiteList{}
	_, err := cnf.GetSubObject(registry.Join(ParNodeName, SubNodeName), &ip)
	if err == conf.ErrNoSetting {
		return &WhiteList{Disable: true}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("white list配置格式有误:%v", err)
	}

	for _, i := range ip.IPS {
		i.ipm = conf.NewPathMatch(i.IPS...)
		i.rqm = conf.NewPathMatch(i.Requests...)
		if b, err := govalidator.ValidateStruct(i); !b {
			return nil, fmt.Errorf("white list配置数据有误:%v", err)
		}

	}
	return &ip, nil
}
