package basic

import (
	"fmt"
	"strconv"

	"github.com/asaskevich/govalidator"
	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/registry"
)

//BasicAuth http basic 认证配置
type BasicAuth struct {
	Excludes       []string          `json:"excludes,omitempty" toml:"exclude,omitempty"`
	Members        map[string]string `json:"members"`
	Disable        bool              `json:"disable,omitempty" toml:"disable,omitempty"`
	*conf.Includes `json:"-"`
	authorization  []*auth `json:"-"`
}

//NewBasic 构建http basic配置参数发
func NewBasic(opts ...Option) *BasicAuth {
	basic := &BasicAuth{
		Excludes: make([]string, 0, 1),
		Members:  make(map[string]string),
	}
	for _, opt := range opts {
		opt(basic)
	}
	basic.Includes = conf.NewInCludes(basic.Excludes...)
	basic.authorization = newAuthorization(basic.Members)
	return basic
}

//Verify 验证用户信息
func (b *BasicAuth) Verify(authValue string) (string, bool) {
	if authValue == "" {
		return "", false
	}
	for _, pair := range b.authorization {
		if pair.auth == authValue {
			return pair.userName, true
		}
	}
	return "", false
}

//GetRealm 获取认证域
func (b *BasicAuth) GetRealm(realm string) string {
	if realm == "" {
		realm = "Authorization Required"
	}
	return "Basic realm=" + strconv.Quote("Authorization Required")
}

//GetConf 获取basic
func GetConf(cnf conf.IMainConf) *BasicAuth {
	basic := BasicAuth{}
	_, err := cnf.GetSubObject(registry.Join("auth", "basic"), &basic)
	if err == conf.ErrNoSetting || len(basic.Members) == 0 {
		return &BasicAuth{Disable: true}
	}
	if err != nil && err != conf.ErrNoSetting {
		panic(fmt.Errorf("basic配置有误:%v", err))
	}
	if b, err := govalidator.ValidateStruct(&basic); !b {
		panic(fmt.Errorf("basic配置有误:%v", err))
	}
	basic.Includes = conf.NewInCludes(basic.Excludes...)
	basic.authorization = newAuthorization(basic.Members)
	return &basic
}
