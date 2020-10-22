/*
根据请示指定限流规则，被限制的请求可以等待一段时间。当启用降级后，将调用对应的降级服务。
未指定降级服务，未提供降级服务时将调用默认的响应配置。如果未配置响应模板则默认返回服务不可用。
*/

package limiter

import (
	"fmt"

	"github.com/asaskevich/govalidator"
	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/lib4go/concurrent/cmap"
)

//Limiter 限流器
type Limiter struct {
	Rules    []*Rule         `json:"rules" toml:"rules,omitempty"`
	Disable  bool            `json:"disable,omitempty" toml:"disable,omitempty"`
	p        *conf.PathMatch `json:"-"`
	limiters cmap.ConcurrentMap
}

//New 构建Limit配置
func New(auth *Rule, opts ...Option) *Limiter {
	limiter := &Limiter{
		Rules:    []*Rule{auth},
		limiters: cmap.New(8),
		Disable:  false,
	}

	for _, f := range opts {
		f(limiter)
	}
	paths := make([]string, 0, len(limiter.Rules)+1)
	for _, v := range limiter.Rules {
		limiter.limiters.Set(v.Path, v)
		paths = append(paths, v.Path)
	}
	limiter.p = conf.NewPathMatch(paths...)
	fmt.Println("NewPathMatch:", paths)
	return limiter
}

//GetLimiter 获取限流器
func (l *Limiter) GetLimiter(path string) (bool, *Rule) {
	fmt.Println("GetLimiter:", path)

	ok, path := l.p.Match(path, "/")
	if !ok {
		return false, nil
	}
	fmt.Println(l.limiters.Items())
	rule, ok := l.limiters.Get(path)
	fmt.Println("l.limiters.Get:", rule, ok)
	if !ok {
		panic("从缓存中未找到limite组件")
	}
	return true, rule.(*Rule)
}

//GetConf 获取jwt
func GetConf(cnf conf.IMainConf) (*Limiter, error) {
	limiter := Limiter{}
	_, err := cnf.GetSubObject(registry.Join("acl", "limit"), &limiter)
	if err == conf.ErrNoSetting || len(limiter.Rules) == 0 {
		return &Limiter{Disable: true}, nil
	}
	if err != nil && err != conf.ErrNoSetting {
		return nil, fmt.Errorf("绑定limit配置有误:%v", err)
	}
	if b, err := govalidator.ValidateStruct(&limiter); !b {
		return nil, fmt.Errorf("limit配置数据有误:%v %+v", err, limiter)
	}
	return New(limiter.Rules[0], WithRuleList(limiter.Rules[1:]...)), nil
}
