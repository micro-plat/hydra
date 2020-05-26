package fsa

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"

	"github.com/asaskevich/govalidator"
	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/lib4go/security/md5"
)

//FixedSecretAuth 创建固定密钥验证服务
type FixedSecretAuth struct {
	Secret  string   `json:"secret" valid:"ascii,required" toml:"secret,omitempty"`
	Mode    string   `json:"mode" valid:"in(MD5|SHA1|SHA256),required" toml:"mode,omitempty"`
	Include []string `json:"include" valid:"required" toml:"include,omitempty"`
	Disable bool     `json:"disable,omitempty" toml:"disable,omitempty"`
	*conf.Includes
}

//New 创建固定密钥验证服务
func New(secret string, opts ...FixedOption) *FixedSecretAuth {
	f := &FixedSecretAuth{
		Secret:  secret,
		Include: []string{"**"},
		Mode:    "MD5",
	}
	for _, opt := range opts {
		opt(f)
	}
	f.Includes = conf.NewInCludes(f.Include...)
	return f
}

//GetConf 获取FixedSecretAuth
func GetConf(cnf conf.IMainConf) *FixedSecretAuth {
	fsa := FixedSecretAuth{}
	_, err := cnf.GetSubObject("fsa", &fsa)
	if err == conf.ErrNoSetting {
		return &FixedSecretAuth{Disable: true, Includes: conf.NewInCludes()}
	}
	if err != nil && err != conf.ErrNoSetting {
		panic(fmt.Errorf("fixed-secret配置有误:%v", err))
	}
	if b, err := govalidator.ValidateStruct(&fsa); !b {
		panic(fmt.Errorf("fixed-secret配置有误:%v", err))
	}
	fsa.Includes = conf.NewInCludes(fsa.Include...)
	return &fsa
}

//CreateSecret 创建Secret
func CreateSecret() string {
	b := make([]byte, 48)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return md5.Encrypt(base64.URLEncoding.EncodeToString(b))
}
