package apikey

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/lib4go/security/md5"
	"github.com/micro-plat/lib4go/security/sha1"
	"github.com/micro-plat/lib4go/security/sha256"
)

//APIKeyAuth 创建固定密钥验证服务
type APIKeyAuth struct {
	Secret   string   `json:"secret" valid:"ascii,required" toml:"secret,omitempty"`
	Mode     string   `json:"mode" valid:"in(MD5|SHA1|SHA256),required" toml:"mode,omitempty"`
	Excludes []string `json:"excludes" valid:"required" toml:"excludes,omitempty"`
	Disable  bool     `json:"disable,omitempty" toml:"disable,omitempty"`
	*conf.PathMatch
}

//New 创建固定密钥验证服务
func New(secret string, opts ...Option) *APIKeyAuth {
	f := &APIKeyAuth{
		Secret: secret,
		Mode:   "MD5",
	}
	for _, opt := range opts {
		opt(f)
	}
	f.PathMatch = conf.NewPathMatch(f.Excludes...)
	return f
}

//Verify 验证签名是否通过
func (a *APIKeyAuth) Verify(raw string, secret string, sign string) error {
	var expect string
	switch strings.ToUpper(a.Mode) {
	case "MD5":
		expect = md5.Encrypt(raw + secret)
	case "SHA1":
		expect = sha1.Encrypt(raw + secret)
	case "SHA256":
		expect = sha256.Encrypt(raw + secret)
	default:
		return fmt.Errorf("不支持的签名验证方式:%v", a.Mode)
	}
	if strings.EqualFold(expect, sign) {
		return nil
	}
	return fmt.Errorf("签名错误:raw:%s,expect:%s,actual:%s", raw, expect, sign)

}

//GetConf 获取APIKeyAuth
func GetConf(cnf conf.IMainConf) *APIKeyAuth {
	fsa := APIKeyAuth{}
	_, err := cnf.GetSubObject(registry.Join("auth", "apikey"), &fsa)
	if err == conf.ErrNoSetting {
		return &APIKeyAuth{Disable: true, PathMatch: conf.NewPathMatch()}
	}
	if err != nil && err != conf.ErrNoSetting {
		panic(fmt.Errorf("apikey配置有误:%v", err))
	}
	if b, err := govalidator.ValidateStruct(&fsa); !b {
		panic(fmt.Errorf("apikey配置有误:%v", err))
	}
	fsa.PathMatch = conf.NewPathMatch(fsa.Excludes...)
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
