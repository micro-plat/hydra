package basic

import (
	"fmt"
	"strconv"

	"github.com/asaskevich/govalidator"
	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/registry"
)

const (
	//ParNodeName auth-basic配置父节点名
	ParNodeName = "auth"
	//SubNodeName auth-basic配置子节点名
	SubNodeName = "basic"
)

//BasicAuth http basic 认证配置
type BasicAuth struct {
	//Excludes 排除路径列表
	Excludes        []string          `json:"excludes,omitempty" toml:"exclude,omitempty"`
	Members         map[string]string `json:"members,omitempty" toml:"members,omitempty"`
	Disable         bool              `json:"disable,omitempty" toml:"disable,omitempty"`
	*conf.PathMatch `json:"-"`
	authorization   []*auth `json:"-"`
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
	basic.PathMatch = conf.NewPathMatch(basic.Excludes...)
	basic.authorization = newAuthorization(basic.Members)
	return basic
}

//Verify 验证用户信息
func (b *BasicAuth) Verify(authValue string) (string, bool) {
	for _, pair := range b.authorization {
		if pair.auth == authValue {
			return pair.userName, true
		}
	}
	return "", false
}

//GetRealm 获取认证域
func (b *BasicAuth) GetRealm() string {
	return "Basic realm=" + strconv.Quote("Authorization Required")
}

//GetConf 获取basic
func GetConf(cnf conf.IServerConf) (*BasicAuth, error) {
	basic := BasicAuth{}
	_, err := cnf.GetSubObject(registry.Join(ParNodeName, SubNodeName), &basic)
	if err == conf.ErrNoSetting || len(basic.Members) == 0 {
		return &BasicAuth{Disable: true}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("basic配置格式有误:%v", err)
	}
	if b, err := govalidator.ValidateStruct(&basic); !b {
		return nil, fmt.Errorf("basic配置数据有误:%v", err)
	}
	basic.PathMatch = conf.NewPathMatch(basic.Excludes...)
	basic.authorization = newAuthorization(basic.Members)
	return &basic, nil
}
