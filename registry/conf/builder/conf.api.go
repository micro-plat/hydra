package builder

import (
	"github.com/micro-plat/hydra/registry/conf/server/api"
	"github.com/micro-plat/hydra/registry/conf/server/auth/fsa"
	"github.com/micro-plat/hydra/registry/conf/server/auth/jwt"
	"github.com/micro-plat/hydra/registry/conf/server/auth/ras"
	"github.com/micro-plat/hydra/registry/conf/server/header"
	"github.com/micro-plat/hydra/registry/conf/server/metric"
	"github.com/micro-plat/hydra/registry/conf/server/render"
	"github.com/micro-plat/hydra/registry/conf/server/static"
	"github.com/micro-plat/hydra/services"
)

type apiBuilder map[string]interface{}

//NewAPI 构建API生成器
func NewAPI(address string, opts ...api.Option) apiBuilder {
	b := make(map[string]interface{})
	b["main"] = api.New(address, opts...)
	routers, err := services.Registry.GetRouters("api")
	if err != nil {
		panic(err)
	}
	b["router"] = routers
	return b
}

//Jwt jwt配置
func (a apiBuilder) Jwt(opts ...jwt.Option) apiBuilder {
	a["jwt"] = jwt.NewJWT(opts...)
	return a
}

//Fsa fsa静态密钥错误
func (a apiBuilder) Fsa(secret string, opts ...fsa.FixedOption) apiBuilder {
	a["fsa"] = fsa.New(secret, opts...)
	return a
}

//Ras 远程认证服务配置
func (a apiBuilder) Ras(service string, opts ...ras.RemotingOption) apiBuilder {
	a["ras"] = ras.New(service, opts...)
	return a
}

//Header 头配置
func (a apiBuilder) Header(opts ...header.Option) apiBuilder {
	a["header"] = header.New(opts...)
	return a
}

//Header 头配置
func (a apiBuilder) Metric(host string, db string, cron string, opts ...metric.Option) apiBuilder {
	a["metric"] = metric.New(host, db, cron, opts...)
	return a
}

//Static 静态文件配置
func (a apiBuilder) Static(opts ...static.Option) apiBuilder {
	a["static"] = static.New(opts...)
	return a
}

//Render 响应渲染配置
func (a apiBuilder) Render(opts ...render.Option) apiBuilder {
	a["render"] = render.New(opts...)
	return a
}
