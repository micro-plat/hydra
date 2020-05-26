package jwt

import (
	"fmt"

	"github.com/asaskevich/govalidator"
	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/lib4go/utility"
)

//JWTAuth jwt配置信息
type JWTAuth struct {
	jwtOption
	*conf.Includes
}

//NewJWT 构建JWT配置参数发
func NewJWT(opts ...Option) *JWTAuth {
	jwt := &JWTAuth{
		jwtOption: jwtOption{
			Name:     "Authorization-Jwt",
			Mode:     "HS512",
			Secret:   utility.GetGUID(),
			ExpireAt: 86400,
			Source:   "COOKIE",
		},
	}
	for _, opt := range opts {
		opt(&jwt.jwtOption)
	}
	jwt.Includes = conf.NewInCludes(jwt.Exclude...)
	return jwt
}

//GetConf 获取jwt
func GetConf(cnf conf.IMainConf) *JWTAuth {
	jwt := JWTAuth{}
	_, err := cnf.GetSubObject("jwt", &jwt)
	if err == conf.ErrNoSetting {
		return &JWTAuth{jwtOption: jwtOption{Disable: true}}
	}
	if err != nil && err != conf.ErrNoSetting {
		panic(fmt.Errorf("jwt配置有误:%v", err))
	}
	if b, err := govalidator.ValidateStruct(&jwt); !b {
		panic(fmt.Errorf("jwt配置有误:%v", err))
	}
	jwt.Includes = conf.NewInCludes(jwt.Exclude...)

	return &jwt
}
