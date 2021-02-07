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
	"github.com/micro-plat/hydra/conf/server/render"
	"github.com/micro-plat/hydra/conf/server/static"
)

type httpBuilder struct {
	BaseBuilder
	tp string
}

//newHTTP 构建http生成器
func newHTTP(tp string, address string, opts ...api.Option) *httpBuilder {
	b := &httpBuilder{tp: tp, BaseBuilder: make(map[string]interface{})}
	b.BaseBuilder[ServerMainNodeName] = api.New(address, opts...)
	return b
}

//Load 加载路由
func (b *httpBuilder) Load() {
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
