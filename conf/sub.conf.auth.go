package conf

import "strings"

//Authes 安全认证组
type Authes struct {
	JWT                  *JWTAuth         `json:"jwt,omitempty"`
	FixedScret           *FixedSecretAuth `json:"fixed-secret,omitempty"`
	RemotingServiceAuths ServiceAuths     `json:"remotings,omitempty"`
}

//JWTAuth jwt安全认证
type JWTAuth struct {
	Name       string   `json:"name" valid:"ascii,required"`
	ExpireAt   int64    `json:"expireAt" valid:"required"`
	Mode       string   `json:"mode" valid:"in(HS256|HS384|HS512|RS256|ES256|ES384|ES512|RS384|RS512|PS256|PS384|PS512),required"`
	Source     string   `json:"source,omitempty" valid:"in(header|cookie|HEADER|COOKIE|H)"`
	Secret     string   `json:"secret" valid:"ascii,required"`
	Exclude    []string `json:"exclude,omitempty"`
	FailedCode string   `json:"failed-code,omitempty"`
	Redirect   string   `json:"redirect,omitempty" valid:"ascii"`
	Domain     string   `json:"domain,omitempty" valid:"ascii"`
	Disable    bool     `json:"disable,omitempty"`
}

//NewAuthes  构建安全认证
func NewAuthes() *Authes {
	return &Authes{
		RemotingServiceAuths: make([]*ServiceAuth, 0, 2),
	}
}

//WithJWT 添加jwt验证
func (a *Authes) WithJWT(jwt *JWTAuth) *Authes {
	a.JWT = jwt
	return a
}

//NewJWT 构建JWT安全认证
func NewJWT(name string, mode string, secret string, expireAt int64, exclude ...string) *JWTAuth {
	return &JWTAuth{
		Name:     name,
		Mode:     mode,
		Secret:   secret,
		ExpireAt: expireAt,
		Exclude:  exclude,
	}
}

//WithHeaderStore 将jwt信息存储到header中
func (a *JWTAuth) WithHeaderStore() *JWTAuth {
	a.Source = "HEADER"
	return a
}

//WithCookieStore 将jwt信息存储到cookie中
func (a *JWTAuth) WithCookieStore(domain ...string) *JWTAuth {
	a.Source = "COOKIE"
	if len(domain) > 0 {
		a.Domain = domain[0]
	}
	return a
}

//WithFailedCode  设置jwt验证失败后的返回给客户端的错误码
func (a *JWTAuth) WithFailedCode(code string) *JWTAuth {
	a.FailedCode = code
	return a
}

//WithRedirect 设置jwt验证失败后的跳转地址
func (a *JWTAuth) WithRedirect(url string) *JWTAuth {
	a.Redirect = url
	return a
}

//WithDisable 禁用配置
func (a *JWTAuth) WithDisable() *JWTAuth {
	a.Disable = true
	return a
}

//WithEnable 启用配置
func (a *JWTAuth) WithEnable() *JWTAuth {
	a.Disable = false
	return a
}

var cacheExcludeSvcs = map[string]bool{}

//IsExcluded 是否是排除验证的服务
func (a *JWTAuth) IsExcluded(service string) bool {

	service = strings.ToLower(service)
	if v, e := cacheExcludeSvcs[service]; e && v {
		return true
	}

	sparties := strings.Split(service, "/")
	//排除指定请求
	for _, u := range a.Exclude {
		//完全匹配
		if strings.EqualFold(u, service) {
			cacheExcludeSvcs[service] = true
			return true
		}
		//分段模糊
		uparties := strings.Split(u, "/")
		//取较少的数组长度
		uc := len(uparties)
		sc := len(sparties)
		/*
			处理模式：
			1. /a/b/ *
			2. /a/ **
			3. /a/ * /d
		**/
		if uc != sc && !strings.HasSuffix(u, "**") {
			continue
		}
		if uc > sc {
			continue
		}
		isMatch := true
		for i := 0; i < uc; i++ {
			if uparties[i] == "**" {
				cacheExcludeSvcs[service] = true
				return true
			}
			if uparties[i] == "*" {
				for j := i + 1; j < uc; j++ {
					if uparties[j] != sparties[j] {
						isMatch = false
						break
					}
				}
				if !isMatch {
					break
				}
				cacheExcludeSvcs[service] = true
				return true
			}
			if uparties[i] != sparties[i] {
				break
			}
		}

	}
	return false
}
