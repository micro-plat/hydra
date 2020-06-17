package jwt

import (
	"fmt"

	"github.com/asaskevich/govalidator"
	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/lib4go/utility"
)

//JWTAuth jwt配置信息
type JWTAuth struct {
	Name            string   `json:"name" valid:"ascii,required" toml:"name,omitempty"`
	ExpireAt        int64    `json:"expireAt" valid:"required" toml:"expireAt,omitzero"`
	Mode            string   `json:"mode" valid:"in(HS256|HS384|HS512|RS256|ES256|ES384|ES512|RS384|RS512|PS256|PS384|PS512),required" toml:"mode,omitempty"`
	Secret          string   `json:"secret" valid:"ascii,required" toml:"secret,omitempty"`
	Source          string   `json:"source,omitempty" valid:"in(header|cookie|HEADER|COOKIE|H)" toml:"source,omitempty"`
	Excludes        []string `json:"excludes,omitempty" toml:"exclude,omitempty"`
	Redirect        string   `json:"redirect,omitempty" valid:"ascii" toml:"redirect,omitempty"`
	Domain          string   `json:"domain,omitempty" valid:"ascii" toml:"domain,omitempty"`
	Disable         bool     `json:"disable,omitempty" toml:"disable,omitempty"`
	*conf.PathMatch `json:"-"`
}

//NewJWT 构建JWT配置参数发
func NewJWT(opts ...Option) *JWTAuth {
	jwt := &JWTAuth{
		Name:     "Authorization-Jwt",
		Mode:     "HS512",
		Secret:   utility.GetGUID(),
		ExpireAt: 86400,
		Source:   "COOKIE",
	}
	for _, opt := range opts {
		opt(jwt)
	}
	jwt.PathMatch = conf.NewPathMatch(jwt.Excludes...)
	return jwt
}

type ConfHandler func(cnf conf.IMainConf) *JWTAuth

func (h ConfHandler) Handle(cnf conf.IMainConf) interface{} {
	return h(cnf)
}

//GetConf 获取jwt
func GetConf(cnf conf.IMainConf) *JWTAuth {
	jwt := JWTAuth{}
	_, err := cnf.GetSubObject(registry.Join("auth", "fsa"), &jwt)
	if err == conf.ErrNoSetting {
		return &JWTAuth{Disable: true}
	}
	if err != nil && err != conf.ErrNoSetting {
		panic(fmt.Errorf("jwt配置有误:%v", err))
	}
	if b, err := govalidator.ValidateStruct(&jwt); !b {
		panic(fmt.Errorf("jwt配置有误:%v", err))
	}
	jwt.PathMatch = conf.NewPathMatch(jwt.Excludes...)

	return &jwt
}
