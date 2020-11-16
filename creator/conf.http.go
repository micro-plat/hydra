package creator

import (
	"fmt"

	"github.com/micro-plat/hydra/conf/server/acl/blacklist"
	"github.com/micro-plat/hydra/conf/server/acl/limiter"
	"github.com/micro-plat/hydra/conf/server/acl/proxy"
	"github.com/micro-plat/hydra/conf/server/acl/whitelist"
	"github.com/micro-plat/hydra/conf/server/api"
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
	CustomerBuilder
	tp          string
	fnGetRouter func(string) *services.ORouter
}

//newHTTP 构建http生成器
func newHTTP(tp string, address string, f func(string) *services.ORouter, opts ...api.Option) *httpBuilder {
	b := &httpBuilder{tp: tp, fnGetRouter: f, CustomerBuilder: make(map[string]interface{})}
	b.CustomerBuilder[ServerMainNodeName] = api.New(address, opts...)
	return b
}

//Load 加载路由
func (b *httpBuilder) Load() {
	routers, err := b.fnGetRouter(b.tp).GetRouters()
	if err != nil {
		panic(err)
	}
	b.CustomerBuilder[router.TypeNodeName] = routers
	return
}

//Jwt jwt配置
func (b *httpBuilder) Jwt(opts ...jwt.Option) *httpBuilder {
	path := fmt.Sprintf("%s/%s", jwt.ParNodeName, jwt.SubNodeName)
	b.CustomerBuilder[path] = jwt.NewJWT(opts...)
	return b
}

//Fsa fsa静态密钥错误
func (b *httpBuilder) APIKEY(secret string, opts ...apikey.Option) *httpBuilder {
	path := fmt.Sprintf("%s/%s", apikey.ParNodeName, apikey.SubNodeName)
	b.CustomerBuilder[path] = apikey.New(secret, opts...)
	return b
}

//Fsa fsa静态密钥错误
func (b *httpBuilder) Basic(opts ...basic.Option) *httpBuilder {
	path := fmt.Sprintf("%s/%s", basic.ParNodeName, basic.SubNodeName)
	b.CustomerBuilder[path] = basic.NewBasic(opts...)
	return b
}

//WhiteList 设置白名单
func (b *httpBuilder) WhiteList(opts ...whitelist.Option) *httpBuilder {
	path := fmt.Sprintf("%s/%s", whitelist.ParNodeName, whitelist.SubNodeName)
	b.CustomerBuilder[path] = whitelist.New(opts...)
	return b
}

//BlackList 设置黑名单
func (b *httpBuilder) BlackList(opts ...blacklist.Option) *httpBuilder {
	path := fmt.Sprintf("%s/%s", blacklist.ParNodeName, blacklist.SubNodeName)
	b.CustomerBuilder[path] = blacklist.New(opts...)
	return b
}

//Ras 远程认证服务配置
func (b *httpBuilder) Ras(opts ...ras.Option) *httpBuilder {
	path := fmt.Sprintf("%s/%s", ras.ParNodeName, ras.SubNodeName)
	b.CustomerBuilder[path] = ras.NewRASAuth(opts...)
	return b
}

//Header 头配置
func (b *httpBuilder) Header(opts ...header.Option) *httpBuilder {
	b.CustomerBuilder[header.TypeNodeName] = header.New(opts...)
	return b
}

//Header 头配置
func (b *httpBuilder) Metric(host string, db string, cron string, opts ...metric.Option) *httpBuilder {
	b.CustomerBuilder[metric.TypeNodeName] = metric.New(host, db, cron, opts...)
	return b
}

//Static 静态文件配置
func (b *httpBuilder) Static(opts ...static.Option) *httpBuilder {
	b.CustomerBuilder[static.TypeNodeName] = static.New(opts...)
	return b
}

//Limit 服务器限流配置
func (b *httpBuilder) Limit(opts ...limiter.Option) *httpBuilder {
	path := fmt.Sprintf("%s/%s", limiter.ParNodeName, limiter.SubNodeName)
	b.CustomerBuilder[path] = limiter.New(opts...)
	return b
}

//Proxy 代理配置
func (b *httpBuilder) Proxy(script string) *httpBuilder {
	path := fmt.Sprintf("%s/%s", proxy.ParNodeName, proxy.SubNodeName)
	b.CustomerBuilder[path] = script
	return b
}

//Render 响应渲染配置
func (b *httpBuilder) Render(opts ...render.Option) *httpBuilder {
	b.CustomerBuilder[render.TypeNodeName] = render.NewRender(opts...)
	return b
}
