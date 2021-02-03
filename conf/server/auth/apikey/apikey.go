package apikey

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/pkgs"
	"github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/lib4go/security/md5"
	"github.com/micro-plat/lib4go/security/sha1"
	"github.com/micro-plat/lib4go/security/sha256"
)

const (
	//ParNodeName auth-apikey配置父节点名
	ParNodeName = "auth"
	//SubNodeName auth-apikey配置子节点名
	SubNodeName = "apikey"
)

const (
	//ModeMD5 md5加密模式
	ModeMD5 = "MD5"
	//ModeSHA1 SHA1加密模式
	ModeSHA1 = "SHA1"
	//ModeSHA256 SHA256加密模式
	ModeSHA256 = "SHA256"

	//ModeSRVC 服务service模式
	ModeSRVC = "SRVC"
)

//APIKeyAuth 创建固定密钥验证服务
type APIKeyAuth struct {
	Secret   string        `json:"secret,omitempty" valid:"ascii,required,stringlength(8|64)" toml:"secret,omitempty" label:"密钥验证服务secret"`
	Mode     string        `json:"mode,omitempty" valid:"in(MD5|SHA1|SHA256|SVS|SRVC),required" toml:"mode,omitempty" label:"密钥验证服务模式"`
	Excludes []string      `json:"excludes,omitempty" toml:"excludes,omitempty"` //排除不验证的路径
	Disable  bool          `json:"disable,omitempty" toml:"disable,omitempty"`
	Invoker  string        `json:"invoker,omitempty" toml:"invoker,omitempty"`
	invoker  *pkgs.Invoker `json:"-"`
	*conf.PathMatch
}

//New 创建固定密钥验证服务
//该对象支持的加密模式:MD5|SHA1|SHA256
func New(secret string, opts ...Option) *APIKeyAuth {
	f := &APIKeyAuth{
		Secret: secret,
		Mode:   ModeMD5,
	}
	for _, opt := range opts {
		opt(f)
	}
	f.PathMatch = conf.NewPathMatch(f.Excludes...)
	if f.Invoker != "" {
		f.invoker = pkgs.NewInvoker(f.Invoker)
		if f.invoker.Allow() {
			f.Mode = ModeSRVC
		}
	}
	return f
}

//Verify 验证签名是否通过
func (a *APIKeyAuth) Verify(raw string, sign string, invoke pkgs.FnInvoker) error {
	if a.Invoker != "" {
		//检查并执行本地服务调用
		if ok, rspns := a.invoker.CheckAndInvoke(invoke); ok {
			return rspns.GetError()
		}
	}
	//根据配置进行签名验证
	var expect string
	switch strings.ToUpper(a.Mode) {
	case ModeMD5:
		expect = md5.Encrypt(raw + a.Secret)
	case ModeSHA1:
		expect = sha1.Encrypt(raw + a.Secret)
	case ModeSHA256:
		expect = sha256.Encrypt(raw + a.Secret)
	default:
		return fmt.Errorf("不支持的签名验证方式:%v", a.Mode)
	}
	if strings.EqualFold(expect, sign) {
		return nil
	}
	return fmt.Errorf("签名错误:raw:%s,expect:%s,actual:%s", raw, expect, sign)

}

//GetConf 获取APIKeyAuth
func GetConf(cnf conf.IServerConf) (*APIKeyAuth, error) {
	f := APIKeyAuth{}
	_, err := cnf.GetSubObject(registry.Join(ParNodeName, SubNodeName), &f)
	if errors.Is(err, conf.ErrNoSetting) {
		return &APIKeyAuth{Disable: true, PathMatch: conf.NewPathMatch()}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("apikey配置格式有误:%v", err)
	}
	if f.Invoker != "" {
		f.invoker = pkgs.NewInvoker(f.Invoker)
		if f.invoker.Allow() {
			f.Mode = ModeSRVC
		}
	}
	if b, err := govalidator.ValidateStruct(&f); !b {
		return nil, fmt.Errorf("apikey配置数据有误:%v", err)
	}
	f.PathMatch = conf.NewPathMatch(f.Excludes...)

	return &f, nil
}

//CreateSecret 创建Secret
func CreateSecret() string {
	b := make([]byte, 48)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return md5.Encrypt(base64.URLEncoding.EncodeToString(b))
}
