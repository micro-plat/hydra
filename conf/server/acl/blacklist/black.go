package blacklist

import (
	"fmt"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/registry"
)

//BlackList 黑名单配置
type BlackList struct {
	Disable bool     `json:"disable,omitempty" toml:"disable,omitempty"`
	IPS     []string `json:"black-list" valid:"required" toml:"black-list,omitempty"`
	ipm     *conf.PathMatch
}

//New 黑名单配置
func New(opts ...Option) *BlackList {
	f := &BlackList{IPS: make([]string, 0, 1)}
	for _, opt := range opts {
		opt(f)

	}
	f.ipm = conf.NewPathMatch(f.IPS...)
	return f
}

//IsDeny 验证当前请求是否在黑名单中
func (w *BlackList) IsDeny(ip string) bool {
	ok, _ := w.ipm.Match(ip)
	return ok
}

//GetConf 获取BlackList
func GetConf(cnf conf.IMainConf) (*BlackList, error) {
	ip := BlackList{}
	_, err := cnf.GetSubObject(registry.Join("acl", "black.list"), &ip)
	if err == conf.ErrNoSetting {
		return &BlackList{Disable: true}, nil
	}
	if err != nil && err != conf.ErrNoSetting {
		return nil, fmt.Errorf("black list配置有误:%v", err)
	}

	ip.ipm = conf.NewPathMatch(ip.IPS...)
	return &ip, nil
}
