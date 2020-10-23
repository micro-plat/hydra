package creator

import (
	"github.com/micro-plat/hydra/conf/server/acl/blacklist"
	"github.com/micro-plat/hydra/conf/server/acl/whitelist"
	"github.com/micro-plat/hydra/conf/server/api"
	"github.com/micro-plat/hydra/conf/server/auth/apikey"
	"github.com/micro-plat/hydra/conf/server/auth/basic"
	"github.com/micro-plat/hydra/conf/server/auth/jwt"
	"github.com/micro-plat/hydra/conf/server/auth/ras"
	"github.com/micro-plat/hydra/conf/server/gray"
	"github.com/micro-plat/hydra/conf/server/header"
	"github.com/micro-plat/hydra/conf/server/limiter"
	"github.com/micro-plat/hydra/conf/server/metric"
	"github.com/micro-plat/hydra/conf/server/render"
	"github.com/micro-plat/hydra/conf/server/static"
	"github.com/micro-plat/hydra/services"
)

type httpBuilder struct {
	customerBuilder
	tp          string
	fnGetRouter func(string) *services.ORouter
}

//newHTTP 构建http生成器
func newHTTP(tp string, address string, f func(string) *services.ORouter, opts ...api.Option) *httpBuilder {
	b := &httpBuilder{tp: tp, fnGetRouter: f, customerBuilder: make(map[string]interface{})}
	b.customerBuilder["main"] = api.New(address, opts...)
	return b
}

//Load 加载路由
func (b *httpBuilder) Load() {
	routers, err := b.fnGetRouter(b.tp).GetRouters()
	if err != nil {
		panic(err)
	}
	b.customerBuilder["router"] = routers
	return
}

//Jwt jwt配置
func (b *httpBuilder) Jwt(opts ...jwt.Option) *httpBuilder {
	b.customerBuilder["auth/jwt"] = jwt.NewJWT(opts...)
	return b
}

//Fsa fsa静态密钥错误
func (b *httpBuilder) APIKEY(secret string, opts ...apikey.Option) *httpBuilder {
	b.customerBuilder["auth/apikey"] = apikey.New(secret, opts...)
	return b
}

//Fsa fsa静态密钥错误
func (b *httpBuilder) Basic(opts ...basic.Option) *httpBuilder {
	b.customerBuilder["auth/basic"] = basic.NewBasic(opts...)
	return b
}

//WhiteList 设置白名单
func (b *httpBuilder) WhiteList(opts ...whitelist.Option) *httpBuilder {
	b.customerBuilder["acl/white.list"] = whitelist.New(opts...)
	return b
}

//BlackList 设置黑名单
func (b *httpBuilder) BlackList(opts ...blacklist.Option) *httpBuilder {
	b.customerBuilder["acl/black.list"] = blacklist.New(opts...)
	return b
}

//Ras 远程认证服务配置
func (b *httpBuilder) Ras(opts ...ras.Option) *httpBuilder {
	b.customerBuilder["auth/ras"] = ras.NewRASAuth(opts...)
	return b
}

//Header 头配置
func (b *httpBuilder) Header(opts ...header.Option) *httpBuilder {
	b.customerBuilder["header"] = header.New(opts...)
	return b
}

//Header 头配置
func (b *httpBuilder) Metric(host string, db string, cron string, opts ...metric.Option) *httpBuilder {
	b.customerBuilder["metric"] = metric.New(host, db, cron, opts...)
	return b
}

//Static 静态文件配置
func (b *httpBuilder) Static(opts ...static.Option) *httpBuilder {
	b.customerBuilder["static"] = static.New(opts...)
	return b
}

//Limit 服务器限流配置
func (b *httpBuilder) Limit(opts ...limiter.Option) *httpBuilder {
	b.customerBuilder["acl/limit"] = limiter.New(opts...)
	return b
}

//Gray 灰度配置
func (b *httpBuilder) Gray(opts ...gray.Option) *httpBuilder {
	b.customerBuilder["acl/gray"] = gray.New(opts...)
	return b
}

//Render 响应渲染配置
func (b *httpBuilder) Render(opts ...render.Option) *httpBuilder {
	b.customerBuilder["render"] = render.NewRender(opts...)
	return b
}
