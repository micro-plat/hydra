package creator

import (
	"fmt"

	"github.com/micro-plat/hydra/conf/server/acl/blacklist"
	"github.com/micro-plat/hydra/conf/server/acl/limiter"
	"github.com/micro-plat/hydra/conf/server/acl/proxy"
	"github.com/micro-plat/hydra/conf/server/acl/whitelist"
	"github.com/micro-plat/hydra/conf/server/api"
	"github.com/micro-plat/hydra/conf/server/apm"
	"github.com/micro-plat/hydra/conf/server/auth/apikey"
	"github.com/micro-plat/hydra/conf/server/auth/basic"
	"github.com/micro-plat/hydra/conf/server/auth/jwt"
	"github.com/micro-plat/hydra/conf/server/auth/ras"
	"github.com/micro-plat/hydra/conf/server/header"
	"github.com/micro-plat/hydra/conf/server/metric"
	"github.com/micro-plat/hydra/conf/server/render"
	"github.com/micro-plat/hydra/conf/server/router"
	"github.com/micro-plat/hydra/conf/server/static"

	"github.com/micro-plat/hydra/services"
)

type httpBuilder struct {
	BaseBuilder
	tp          string
	fnGetRouter func(string) *services.ORouter
}

//newHTTP 构建http生成器
func newHTTP(tp string, address string, f func(string) *services.ORouter, opts ...api.Option) *httpBuilder {
	b := &httpBuilder{tp: tp, fnGetRouter: f, BaseBuilder: make(map[string]interface{})}
	b.BaseBuilder[ServerMainNodeName] = api.New(address, opts...)
	return b
}

//Load 加载路由
func (b *httpBuilder) Load() {
	routers, err := b.fnGetRouter(b.tp).GetRouters()
	if err != nil {
		panic(err)
	}
	b.BaseBuilder[router.TypeNodeName] = routers
	return
}

//Jwt jwt配置
func (b *httpBuilder) Jwt(opts ...jwt.Option) *httpBuilder {
	path := fmt.Sprintf("%s/%s", jwt.ParNodeName, jwt.SubNodeName)
	b.BaseBuilder[path] = jwt.NewJWT(opts...)
	return b
}

//Fsa fsa静态密钥错误
func (b *httpBuilder) APIKEY(secret string, opts ...apikey.Option) *httpBuilder {
	path := fmt.Sprintf("%s/%s", apikey.ParNodeName, apikey.SubNodeName)
	b.BaseBuilder[path] = apikey.New(secret, opts...)
	return b
}

//Fsa fsa静态密钥错误
func (b *httpBuilder) Basic(opts ...basic.Option) *httpBuilder {
	path := fmt.Sprintf("%s/%s", basic.ParNodeName, basic.SubNodeName)
	b.BaseBuilder[path] = basic.NewBasic(opts...)
	return b
}

//WhiteList 设置白名单
func (b *httpBuilder) WhiteList(opts ...whitelist.Option) *httpBuilder {
	path := fmt.Sprintf("%s/%s", whitelist.ParNodeName, whitelist.SubNodeName)
	b.BaseBuilder[path] = whitelist.New(opts...)
	return b
}

//BlackList 设置黑名单
func (b *httpBuilder) BlackList(opts ...blacklist.Option) *httpBuilder {
	path := fmt.Sprintf("%s/%s", blacklist.ParNodeName, blacklist.SubNodeName)
	b.BaseBuilder[path] = blacklist.New(opts...)
	return b
}

//Ras 远程认证服务配置
func (b *httpBuilder) Ras(opts ...ras.Option) *httpBuilder {
	path := fmt.Sprintf("%s/%s", ras.ParNodeName, ras.SubNodeName)
	b.BaseBuilder[path] = ras.NewRASAuth(opts...)
	return b
}

//Header 头配置
func (b *httpBuilder) Header(opts ...header.Option) *httpBuilder {
	b.BaseBuilder[header.TypeNodeName] = header.New(opts...)
	return b
}

//Header 头配置
func (b *httpBuilder) Metric(host string, db string, cron string, opts ...metric.Option) *httpBuilder {
	b.BaseBuilder[metric.TypeNodeName] = metric.New(host, db, cron, opts...)
	return b
}

//Static 静态文件配置
func (b *httpBuilder) Static(opts ...static.Option) *httpBuilder {
	b.BaseBuilder[static.TypeNodeName] = static.New(opts...)
	return b
}

//Limit 服务器限流配置
func (b *httpBuilder) Limit(opts ...limiter.Option) *httpBuilder {
	path := fmt.Sprintf("%s/%s", limiter.ParNodeName, limiter.SubNodeName)
	b.BaseBuilder[path] = limiter.New(opts...)
	return b
}

//Proxy 代理配置
func (b *httpBuilder) Proxy(script string) *httpBuilder {
	path := fmt.Sprintf("%s/%s", proxy.ParNodeName, proxy.SubNodeName)
	b.BaseBuilder[path] = script
	return b
}

//Render 响应渲染配置
func (b *httpBuilder) Render(script string) *httpBuilder {
	b.BaseBuilder[render.TypeNodeName] = script
	return b
}

//APM 构建APM配置
func (b *httpBuilder) APM(address string) *httpBuilder {
	b.BaseBuilder[apm.TypeNodeName] = apm.New(address)
	return b
}
