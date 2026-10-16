package jwt

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/conf/pkgs/security"
	"github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/lib4go/errs"
	"github.com/micro-plat/lib4go/security/jwt"
	"github.com/micro-plat/lib4go/types"
	"github.com/micro-plat/lib4go/utility"
)

const (
	JWTStatusTokenNotExsit = http.StatusUnauthorized
	JWTStatusTokenExpired  = http.StatusForbidden
	JWTStatusTokenError    = http.StatusForbidden
	JWTStatusConfError     = http.StatusNotExtended
	JWTStatusConfDataError = http.StatusInternalServerError
	JWTStatusRedirect      = http.StatusFound
)

const (
	ParNodeName = "auth"
	SubNodeName = "jwt"
)

const (
	ModeHS256 = "HS256"
	ModeHS384 = "HS384"
	ModeHS512 = "HS512"
	ModeRS256 = "RS256"
	ModeES256 = "ES256"
	ModeES384 = "ES384"
	ModeES512 = "ES512"
	ModeRS384 = "RS384"
	ModeRS512 = "RS512"
	ModePS256 = "PS256"
	ModePS384 = "PS384"
	ModePS512 = "PS512"
)

const JWTName = "Authorization-Jwt"

type JWTAuth struct {
	security.ConfEncrypt
	Name             string            `json:"name,omitempty" valid:"ascii,required" toml:"name,omitempty" label:"jwt name"`
	ExpireAt         int64             `json:"expireAt,omitzero" valid:"required" toml:"expireAt,omitzero" label:"jwt expire seconds"`
	Mode             string            `json:"mode,omitempty" valid:"in(HS256|HS384|HS512|RS256|ES256|ES384|ES512|RS384|RS512|PS256|PS384|PS512),required" toml:"mode,omitempty" label:"jwt mode"`
	Secret           string            `json:"secret,omitempty" valid:"ascii,required" toml:"secret,omitempty" label:"jwt secret"`
	Source           string            `json:"source,omitempty" valid:"in(header|cookie|HEADER|COOKIE|H)" toml:"source,omitempty" label:"jwt source"`
	Excludes         []string          `json:"excludes,omitempty" toml:"exclude,omitempty"`
	Domain           string            `json:"domain,omitempty" toml:"domain,omitempty"`
	AuthURL          string            `json:"authURL,omitempty" valid:"ascii" toml:"authURL,omitempty" label:"jwt auth redirect url"`
	AuthenticateURL  string            `json:"authenticateURL,omitempty" valid:"ascii" toml:"authenticateURL,omitempty" label:"jwt authenticate url"`
	AuthenticateURLs map[string]string `json:"authenticateURLs,omitempty" toml:"authenticateURLs,omitempty"`
	Disable          bool              `json:"disable,omitempty" toml:"disable,omitempty"`
	*conf.PathMatch  `json:"-"`
}

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

func (j *JWTAuth) CheckJWT(token string) (data interface{}, err error) {
	if token == "" {
		return nil, errs.NewError(JWTStatusTokenNotExsit, fmt.Errorf("jwt token is empty(%s %s)", j.Source, j.Name))
	}
	if !strings.HasPrefix(token, TokenBearerPrefix) {
		return nil, errs.NewError(JWTStatusTokenError, fmt.Errorf("jwt token format error(%s)", token))
	}
	token = token[len(TokenBearerPrefix):]

	data, er := jwt.Decrypt(token, j.Secret)
	if er != nil {
		if strings.Contains(er.Error(), "Token is expired") {
			return nil, errs.NewError(JWTStatusTokenExpired, er)
		}
		return data, errs.NewError(JWTStatusTokenError, fmt.Errorf("jwt token value(%s) error %w", token, er))
	}

	return data, nil
}

func (j *JWTAuth) GetAuthenticateURL(path string) string {
	if j.AuthenticateURLs != nil {
		if url := j.AuthenticateURLs[path]; url != "" {
			return url
		}
	}
	return j.AuthenticateURL
}

// HasPathAuthenticateURL 判断指定路径是否显式配置了路径级认证URL。
// 仅显式配置了 AuthenticateURLs[path] 的路径才在未登录时附加 return_to，
// 避免对未配置的路径（如 API）生成无意义的 return_to。
func (j *JWTAuth) HasPathAuthenticateURL(path string) bool {
	return j.AuthenticateURLs != nil && j.AuthenticateURLs[path] != ""
}

func (j *JWTAuth) SetAuthenticateURLForPath(path string, url string) {
	if j.AuthenticateURLs == nil {
		j.AuthenticateURLs = make(map[string]string)
	}
	j.AuthenticateURLs[path] = url
}

func (j *JWTAuth) AddExcludes(paths ...string) {
	j.Excludes = append(j.Excludes, paths...)
	j.PathMatch = conf.NewPathMatch(j.Excludes...)
}

func (j *JWTAuth) GetJWTForRspns(token string, expired ...bool) (string, string, bool) {
	isExpired := types.GetBoolByIndex(expired, 0, false)
	token = TokenBearerPrefix + token
	switch strings.ToUpper(j.Source) {
	case SourceHeader, SourceHeaderShort:
		return AuthorizationHeader, token, isExpired == false
	default:
		expireVal := j.getExpireTime(isExpired)
		if j.Domain != "" {
			return "Set-Cookie", fmt.Sprintf("%s=%s;domain=%s;path=/;expires=%s;HttpOnly", j.Name, token, j.Domain, expireVal), true
		}
		return "Set-Cookie", fmt.Sprintf("%s=%s;path=/;expires=%s;HttpOnly", j.Name, token, expireVal), true
	}
}

func (j *JWTAuth) getExpireTime(expired bool) string {
	expireTime := time.Now().Add(time.Hour * -24)
	if !expired {
		expireTime = time.Now().Add(time.Duration(time.Duration(j.ExpireAt)*time.Second - 8*60*60*time.Second))
	}
	return expireTime.Format("Mon, 02 Jan 2006 15:04:05 GMT")
}

func GetConf(cnf conf.IServerConf) (*JWTAuth, error) {
	jwt := JWTAuth{}
	_, err := cnf.GetSubObject(registry.Join(ParNodeName, SubNodeName), &jwt)
	if errors.Is(err, conf.ErrNoSetting) {
		return &JWTAuth{Disable: true}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("jwt conf format error:%v", err)
	}
	if b, err := govalidator.ValidateStruct(&jwt); !b {
		return nil, fmt.Errorf("jwt conf data error:%v", err)
	}
	jwt.PathMatch = conf.NewPathMatch(jwt.Excludes...)

	return &jwt, nil
}
