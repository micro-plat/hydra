package creator

import (
	"github.com/micro-plat/hydra/conf/server/api"
	"github.com/micro-plat/hydra/conf/server/auth/apikey"
	"github.com/micro-plat/hydra/conf/server/auth/basic"
	"github.com/micro-plat/hydra/conf/server/auth/jwt"
	"github.com/micro-plat/hydra/conf/server/auth/ras"
	"github.com/micro-plat/hydra/conf/server/header"
	"github.com/micro-plat/hydra/conf/server/metric"
	"github.com/micro-plat/hydra/conf/server/render"
	"github.com/micro-plat/hydra/conf/server/static"
	"github.com/micro-plat/hydra/services"
)

type httpBuilder = customerBuilder

//newHTTP 构建http生成器
func newHTTP(address string, opts ...api.Option) httpBuilder {
	b := make(map[string]interface{})
	b["main"] = api.New(address, opts...)
	return b
}

//Load 加载路由
func (b httpBuilder) Load() {
	routers, err := services.API.GetRouters()
	if err != nil {
		panic(err)
	}
	b["router"] = routers
	return
}

//Jwt jwt配置
func (b httpBuilder) Jwt(opts ...jwt.Option) httpBuilder {
	b["auth/jwt"] = jwt.NewJWT(opts...)
	return b
}

//Fsa fsa静态密钥错误
func (b httpBuilder) APIKEY(secret string, opts ...apikey.Option) httpBuilder {
	b["auth/apikey"] = apikey.New(secret, opts...)
	return b
}

//Fsa fsa静态密钥错误
func (b httpBuilder) Basic(opts ...basic.Option) httpBuilder {
	b["auth/basic"] = basic.NewBasic(opts...)
	return b
}

//Ras 远程认证服务配置
func (b httpBuilder) Ras(service string, opts ...ras.Option) httpBuilder {
	b["auth/ras"] = ras.New(service, opts...)
	return b
}

//Header 头配置
func (b httpBuilder) Header(opts ...header.Option) httpBuilder {
	b["header"] = header.New(opts...)
	return b
}

//Header 头配置
func (b httpBuilder) Metric(host string, db string, cron string, opts ...metric.Option) httpBuilder {
	b["metric"] = metric.New(host, db, cron, opts...)
	return b
}

//Static 静态文件配置
func (b httpBuilder) Static(opts ...static.Option) httpBuilder {
	b["static"] = static.New(opts...)
	return b
}

//Render 响应渲染配置
func (b httpBuilder) Render(opts ...render.Option) httpBuilder {
	b["render"] = render.NewRender(opts...)
	return b
}
