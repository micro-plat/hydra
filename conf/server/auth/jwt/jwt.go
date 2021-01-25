package jwt

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/lib4go/errs"
	"github.com/micro-plat/lib4go/security/jwt"
	"github.com/micro-plat/lib4go/types"
	"github.com/micro-plat/lib4go/utility"
)

const (
	//JWTStatusTokenNotExsit jwt token 不存在
	JWTStatusTokenNotExsit = http.StatusUnauthorized
	//JWTStatusTokenExpired jwt token过期
	JWTStatusTokenExpired = http.StatusForbidden
	//JWTStatusTokenError  jwt token错误
	JWTStatusTokenError = http.StatusForbidden
	//JWTStatusConfError jwt配置错误
	JWTStatusConfError = http.StatusNotExtended
	//JWTStatusConfDataError jwt配置数据错误
	JWTStatusConfDataError = http.StatusInternalServerError
	//JWTStatusRedirect jwt跳转
	JWTStatusRedirect = http.StatusFound
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
	Name            string   `json:"name,omitempty" valid:"ascii,required" toml:"name,omitempty" label:"jwt名称"`
	ExpireAt        int64    `json:"expireAt,omitzero" valid:"required" toml:"expireAt,omitzero" label:"jwt过期时间"`
	Mode            string   `json:"mode,omitempty" valid:"in(HS256|HS384|HS512|RS256|ES256|ES384|ES512|RS384|RS512|PS256|PS384|PS512),required" toml:"mode,omitempty" label:"jwt认证方式"`
	Secret          string   `json:"secret,omitempty" valid:"ascii,required" toml:"secret,omitempty"  label:"jwt密钥"`
	Source          string   `json:"source,omitempty" valid:"in(header|cookie|HEADER|COOKIE|H)" toml:"source,omitempty" label:"jwt存储方式"`
	Excludes        []string `json:"excludes,omitempty" toml:"exclude,omitempty"`
	Domain          string   `json:"domain,omitempty" toml:"domain,omitempty"`
	AuthURL         string   `json:"authURL,omitempty" valid:"ascii" toml:"authURL,omitempty" label:"jwt认证跳转地址"`
	Disable         bool     `json:"disable,omitempty" toml:"disable,omitempty"`
	*conf.PathMatch `json:"-"`
}

//NewJWT 构建JWT配置参数
func NewJWT(opts ...Option) *JWTAuth {
	jwt := &JWTAuth{
		Name:     AuthorizationHeader,
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

//CheckJWT 检查jwt合法性
func (j *JWTAuth) CheckJWT(token string) (data interface{}, err error) {
	if token == "" {
		return nil, errs.NewError(JWTStatusTokenError, fmt.Errorf("未传入jwt.token(%s %s值为空)", j.Source, j.Name))
	}
	if !strings.HasPrefix(token, TokenBearerPrefix) {
		return nil, errs.NewError(JWTStatusTokenError, fmt.Errorf("jwt.token格式错误(%s)", token))
	}
	token = token[len(TokenBearerPrefix):]

	//2. 解密jwt判断是否有效，是否过期
	data, er := jwt.Decrypt(token, j.Secret)
	if er != nil {
		if strings.Contains(er.Error(), "Token is expired") {
			return nil, errs.NewError(JWTStatusTokenExpired, er)
		}
		return data, errs.NewError(JWTStatusTokenError, fmt.Errorf("jwt.token值(%s)有误 %w", token, er))
	}

	return data, nil
}

//GetJWTForRspns 获取jwt响应参数值
func (j *JWTAuth) GetJWTForRspns(token string, expired ...bool) (string, string, bool) {

	isExpired := types.GetBoolByIndex(expired, 0, false)
	token = TokenBearerPrefix + token
	switch strings.ToUpper(j.Source) {
	case SourceHeader, SourceHeaderShort: //"HEADER", "H":
		return AuthorizationHeader, token, isExpired == false
	default:
		expireVal := j.getExpireTime(isExpired)
		if j.Domain != "" {
			return "Set-Cookie", fmt.Sprintf("%s=%s;domain=%s;path=/;expires=%s;HttpOnly", j.Name, token, j.Domain, expireVal), true
		}
		return "Set-Cookie", fmt.Sprintf("%s=%s;path=/;expires=%s;HttpOnly", j.Name, token, expireVal), true
	}
}

//getExpireTime 获取jwt的超时时间
func (j *JWTAuth) getExpireTime(expired bool) string {
	expireTime := time.Now().Add(time.Hour * -24)
	if !expired {
		expireTime = time.Now().Add(time.Duration(time.Duration(j.ExpireAt)*time.Second - 8*60*60*time.Second))
	}
	return expireTime.Format("Mon, 02 Jan 2006 15:04:05 GMT")
}

//GetConf 获取jwt配置
func GetConf(cnf conf.IServerConf) (*JWTAuth, error) {
	jwt := JWTAuth{}
	_, err := cnf.GetSubObject(registry.Join(ParNodeName, SubNodeName), &jwt)
	if errors.Is(err, conf.ErrNoSetting) {
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
