package jwt

import (
	"strings"

	"github.com/micro-plat/lib4go/utility"
)

//JWTAuth jwt配置信息
type JWTAuth struct {
	*jwtOption
}

//NewJWT 构建JWT配置参数发
func NewJWT(opts ...Option) *JWTAuth {
	jwt := &JWTAuth{
		jwtOption: &jwtOption{
			Name:     "hsid",
			Mode:     "HS512",
			Secret:   utility.GetGUID(),
			ExpireAt: 86400,
			Source:   "COOKIE",
		},
	}
	for _, opt := range opts {
		opt(jwt.jwtOption)
	}
	return jwt
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
