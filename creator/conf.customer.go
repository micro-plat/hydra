package creator

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/micro-plat/hydra/conf/server/acl/blacklist"
	"github.com/micro-plat/hydra/conf/server/acl/limiter"
	"github.com/micro-plat/hydra/conf/server/acl/proxy"
	"github.com/micro-plat/hydra/conf/server/acl/whitelist"
	"github.com/micro-plat/hydra/conf/server/apm"
	"github.com/micro-plat/hydra/conf/server/auth/apikey"
	"github.com/micro-plat/hydra/conf/server/auth/basic"
	"github.com/micro-plat/hydra/conf/server/auth/jwt"
	"github.com/micro-plat/hydra/conf/server/auth/ras"
	"github.com/micro-plat/hydra/conf/server/header"
	"github.com/micro-plat/hydra/conf/server/metric"
	"github.com/micro-plat/hydra/conf/server/render"
	"github.com/micro-plat/hydra/conf/server/static"
)

type ISUB interface {
	Sub(name string, s ...interface{}) ISUB
}

type iCustomerBuilder interface {
	Load()
	ISUB
	Map() map[string]interface{}
}

var _ iCustomerBuilder = CustomerBuilder{}

type CustomerBuilder map[string]interface{}

//newHTTP 构建http生成器
func newCustomerBuilder(s ...interface{}) CustomerBuilder {
	b := make(map[string]interface{})
	if len(s) == 0 {
		b[ServerMainNodeName] = make(map[string]interface{})
		return b
	}
	b[ServerMainNodeName] = s[0]
	return b
}

//Sub 子配置
func (b CustomerBuilder) Sub(name string, s ...interface{}) ISUB {
	if len(s) == 0 {
		panic(fmt.Sprintf("配置：%s值不能为空", name))
	}
	tp := reflect.TypeOf(s[0])
	val := reflect.ValueOf(s[0])
	if tp.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	switch tp.Kind() {
	case reflect.String:
		b[name] = json.RawMessage([]byte(val.Interface().(string)))
	case reflect.Struct, reflect.Ptr, reflect.Map:
		b[name] = val.Interface()
	default:
		panic(fmt.Sprintf("配置：%s值类型不支持", name))
	}
	return b
}
func (b CustomerBuilder) Map() map[string]interface{} {
	return b
}
func (b CustomerBuilder) Load() {
}

//Jwt jwt配置
func (b *CustomerBuilder) Jwt(opts ...jwt.Option) *CustomerBuilder {
	path := fmt.Sprintf("%s/%s", jwt.ParNodeName, jwt.SubNodeName)
	(*b)[path] = jwt.NewJWT(opts...)
	return b
}

//Fsa fsa静态密钥错误
func (b *CustomerBuilder) APIKEY(secret string, opts ...apikey.Option) *CustomerBuilder {
	path := fmt.Sprintf("%s/%s", apikey.ParNodeName, apikey.SubNodeName)
	(*b)[path] = apikey.New(secret, opts...)
	return b
}

//Fsa fsa静态密钥错误
func (b *CustomerBuilder) Basic(opts ...basic.Option) *CustomerBuilder {
	path := fmt.Sprintf("%s/%s", basic.ParNodeName, basic.SubNodeName)
	(*b)[path] = basic.NewBasic(opts...)
	return b
}

//WhiteList 设置白名单
func (b *CustomerBuilder) WhiteList(opts ...whitelist.Option) *CustomerBuilder {
	path := fmt.Sprintf("%s/%s", whitelist.ParNodeName, whitelist.SubNodeName)
	(*b)[path] = whitelist.New(opts...)
	return b
}

//BlackList 设置黑名单
func (b *CustomerBuilder) BlackList(opts ...blacklist.Option) *CustomerBuilder {
	path := fmt.Sprintf("%s/%s", blacklist.ParNodeName, blacklist.SubNodeName)
	(*b)[path] = blacklist.New(opts...)
	return b
}

//Ras 远程认证服务配置
func (b *CustomerBuilder) Ras(opts ...ras.Option) *CustomerBuilder {
	path := fmt.Sprintf("%s/%s", ras.ParNodeName, ras.SubNodeName)
	(*b)[path] = ras.NewRASAuth(opts...)
	return b
}

//Header 头配置
func (b *CustomerBuilder) Header(opts ...header.Option) *CustomerBuilder {
	(*b)[header.TypeNodeName] = header.New(opts...)
	return b
}

//Header 头配置
func (b *CustomerBuilder) Metric(host string, db string, cron string, opts ...metric.Option) *CustomerBuilder {
	(*b)[metric.TypeNodeName] = metric.New(host, db, cron, opts...)
	return b
}

//Static 静态文件配置
func (b *CustomerBuilder) Static(opts ...static.Option) *CustomerBuilder {
	(*b)[static.TypeNodeName] = static.New(opts...)
	return b
}

//Limit 服务器限流配置
func (b *CustomerBuilder) Limit(opts ...limiter.Option) *CustomerBuilder {
	path := fmt.Sprintf("%s/%s", limiter.ParNodeName, limiter.SubNodeName)
	(*b)[path] = limiter.New(opts...)
	return b
}

//Proxy 代理配置
func (b *CustomerBuilder) Proxy(script string) *CustomerBuilder {
	path := fmt.Sprintf("%s/%s", proxy.ParNodeName, proxy.SubNodeName)
	(*b)[path] = script
	return b
}

//Render 响应渲染配置
func (b *CustomerBuilder) Render(script string) *CustomerBuilder {
	(*b)[render.TypeNodeName] = script
	return b
}

//APM 构建APM配置
func (b *CustomerBuilder) APM(address string) *CustomerBuilder {
	(*b)[apm.TypeNodeName] = apm.New(address)
	return b
}
