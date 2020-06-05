package whitelist

import (
	"fmt"

	"github.com/asaskevich/govalidator"
	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/registry"
)

//IPList ip列表
type IPList struct {
	Requests []string `json:"requests" valid:"ascii,required" toml:"requests,omitempty"`
	IPS      []string `json:"ips" valid:"required" toml:"ips,omitempty"`
	ipm      *conf.PathMatch
	rqm      *conf.PathMatch
}

//WhiteList 白名单配置
type WhiteList struct {
	Disable bool      `json:"disable,omitempty" toml:"disable,omitempty"`
	IPS     []*IPList `json:"white-list,omitempty" toml:"white-list,omitempty"`
}

//New 创建固定密钥验证服务
func New(ips ...*IPList) *WhiteList {
	f := &WhiteList{IPS: make([]*IPList, 0, 1)}
	for _, ip := range ips {
		ip.ipm = conf.NewPathMatch(ip.IPS...)
		ip.rqm = conf.NewPathMatch(ip.Requests...)
		f.IPS = append(f.IPS, ip)
	}
	return f
}

//IsAllow 验证当前请求是否在白名单中
func (w *WhiteList) IsAllow(path string, ip string) bool {
	for _, cur := range w.IPS {
		if ok, _ := cur.rqm.Match(path); ok {
			ok, _ := cur.ipm.Match(ip)
			return ok
		}
	}
	return true
}

//GetConf 获取WhiteList
func GetConf(cnf conf.IMainConf) *WhiteList {
	ip := WhiteList{}
	_, err := cnf.GetSubObject(registry.Join("acl", "white.list"), &ip)
	if err == conf.ErrNoSetting {
		return &WhiteList{Disable: true}
	}
	if err != nil && err != conf.ErrNoSetting {
		panic(fmt.Errorf("white list配置有误:%v", err))
	}

	for _, i := range ip.IPS {
		i.ipm = conf.NewPathMatch(i.IPS...)
		i.rqm = conf.NewPathMatch(i.Requests...)
		if b, err := govalidator.ValidateStruct(i); !b {
			panic(fmt.Errorf("white list配置有误:%v", err))
		}

	}
	return &ip
}
