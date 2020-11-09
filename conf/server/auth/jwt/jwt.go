package jwt

import (
	"fmt"

	"github.com/asaskevich/govalidator"
	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/lib4go/utility"
)

const (
	//ParNodeName auth-jwt配置父节点名
	ParNodeName = "auth"
	//SubNodeName auth-jwt配置子节点名
	SubNodeName = "jwt"
)

const (
	//ModeHS256 加密模式HS256
	ModeHS256 = "HS256"
	//ModeHS384 加密模式HS384
	ModeHS384 = "HS384"
	//ModeHS512 加密模式HS512
	ModeHS512 = "HS512"
	//ModeRS256 加密模式RS256
	ModeRS256 = "RS256"
	//ModeES256 加密模式ES256
	ModeES256 = "ES256"
	//ModeES384 加密模式ES384
	ModeES384 = "ES384"
	//ModeES512 加密模式ES512
	ModeES512 = "ES512"
	//ModeRS384 加密模式RS384
	ModeRS384 = "RS384"
	//ModeRS512 加密模式RS512
	ModeRS512 = "RS512"
	//ModePS256 加密模式PS256
	ModePS256 = "PS256"
	//ModePS384 加密模式PS384
	ModePS384 = "PS384"
	//ModePS512 加密模式PS512
	ModePS512 = "PS512"
)

//JWTName 节点标识名
const JWTName = "Authorization-Jwt"

//JWTAuth jwt配置信息
type JWTAuth struct {
	Name     string   `json:"name,omitempty" valid:"ascii,required" toml:"name,omitempty"`
	ExpireAt int64    `json:"expireAt,omitempty" valid:"required" toml:"expireAt,omitzero"`
	Mode     string   `json:"mode,omitempty" valid:"in(HS256|HS384|HS512|RS256|ES256|ES384|ES512|RS384|RS512|PS256|PS384|PS512),required" toml:"mode,omitempty"`
	Secret   string   `json:"secret,omitempty" valid:"ascii,required" toml:"secret,omitempty"`
	Source   string   `json:"source,omitempty" valid:"in(header|cookie|HEADER|COOKIE|H)" toml:"source,omitempty"`
	Excludes []string `json:"excludes,omitempty" toml:"exclude,omitempty"`
	//Redirect        string   `json:"redirect,omitempty" valid:"ascii" toml:"redirect,omitempty"`
	AuthURL         string `json:"authURL,omitempty" valid:"ascii" toml:"authURL,omitempty"`
	Disable         bool   `json:"disable,omitempty" toml:"disable,omitempty"`
	*conf.PathMatch `json:"-"`
}

//NewJWT 构建JWT配置参数
func NewJWT(opts ...Option) *JWTAuth {
	jwt := &JWTAuth{
		Name:     JWTName,
		Mode:     ModeHS512,
		Secret:   utility.GetGUID(),
		ExpireAt: 86400,
		Source:   SourceCookie,
	}
	for _, opt := range opts {
		opt(jwt)
	}
	jwt.PathMatch = conf.NewPathMatch(jwt.Excludes...)
	return jwt
}

//GetConf 获取jwt配置
func GetConf(cnf conf.IServerConf) (*JWTAuth, error) {
	jwt := JWTAuth{}
	_, err := cnf.GetSubObject(registry.Join(ParNodeName, SubNodeName), &jwt)
	if err == conf.ErrNoSetting {
		return &JWTAuth{Disable: true}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("jwt配置格式有误:%v", err)
	}
	if b, err := govalidator.ValidateStruct(&jwt); !b {
		return nil, fmt.Errorf("jwt配置数据有误:%v", err)
	}
	jwt.PathMatch = conf.NewPathMatch(jwt.Excludes...)

	return &jwt, nil
}
