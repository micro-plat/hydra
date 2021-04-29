package ctx

import (
	"net/url"
	"strings"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/conf/app"
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/hydra/services"
	"github.com/micro-plat/lib4go/encoding"
	"github.com/micro-plat/lib4go/types"
)

var _ context.IPath = &rpath{}

//rpath 处理请求的路径信息
type rpath struct {
	ctx      context.IInnerContext
	appConf  app.IAPPConf
	meta     conf.IMeta
	isLimit  bool
	fallback bool
	encoding string
	params   types.XMap
}

func NewRpath(ctx context.IInnerContext, appConf app.IAPPConf, meta conf.IMeta) *rpath {
	return &rpath{
		ctx:     ctx,
		appConf: appConf,
		meta:    meta,
	}
}

//GetMethod 获取服务请求方式
func (c *rpath) GetMethod() string {
	return c.ctx.GetMethod()
}

// Params 获取路径中的参数
func (c *rpath) Params() types.XMap {
	if c.params != nil {
		return c.params
	}
	c.params = c.ctx.GetParams()
	return c.params
}

//FormatService 通过GetProcessorConf 格式化服务名
func (c *rpath) FormatService(service string) string {
	p, err := c.appConf.GetProcessorConf()
	if err != nil {
		return service
	}
	return p.FormatService(service)
}

//GetService 获取服务名称
func (c *rpath) GetService() string {
	tp := c.appConf.GetServerConf().GetServerType()
	switch tp {
	case global.API, global.Web, global.WS, global.RPC, global.MQC, global.CRON:
		routerObj, err := services.GetRouter(tp).GetRouters()
		if err != nil {
			return ""
		}
		router, err := routerObj.Match(c.ctx.GetRouterPath(), c.ctx.GetMethod())
		if err != nil {
			return ""
		}
		return router.Service
	default:
		return c.ctx.GetRouterPath()
	}
}

//GetGroup 获取当前服务注册的group名
func (c *rpath) GetGroup() string {
	return services.Def.GetGroup(c.appConf.GetServerConf().GetServerType(), c.GetService())
}

//GetPageAndTag 获取服务对应的页面路径与tag标签(page:静态文件prefix+服务原始注册路径,tag：对象中的函数名)
func (c *rpath) GetPageAndTag() (page string, tag string, ok bool) {

	//获取服务注册的路径名，tag标签
	tp := c.appConf.GetServerConf().GetServerType()
	page, tag, ok = services.Def.GetRawPathAndTag(tp, c.GetService())
	if !ok {
		return "", "", false
	}

	//处理tag为空时，获取当前method
	if tag == "" {
		tag = c.ctx.GetMethod()
	}

	return page, tag, ok
}

func (c *rpath) GetEncoding() string {
	if c.encoding != "" {
		return c.encoding
	}

	if c.encoding = c.getEncoding(); c.encoding != "" {
		return c.encoding
	}

	//从请求header中获取
	charsetStr := strings.Join(c.ctx.GetHeaders()["Content-Type"], ",")
	if !strings.Contains(charsetStr, "charset=") {
		charsetStr = strings.Join(c.ctx.GetHeaders()["Accept-Charset"], ",")
	}
	switch {
	case strings.Contains(charsetStr, encoding.GB2312):
		c.encoding = encoding.GB2312
	case strings.Contains(charsetStr, encoding.GBK):
		c.encoding = encoding.GBK
	}
	c.encoding = types.GetString(c.encoding, encoding.UTF8)

	return c.encoding
}

//getEncoding 获取路由配置的编码
func (c *rpath) getEncoding() string {
	tp := c.appConf.GetServerConf().GetServerType()
	switch tp {
	case global.API, global.Web, global.WS, global.RPC, global.MQC, global.CRON:
		routerObj, err := services.GetRouter(tp).GetRouters()
		if err != nil {
			return ""
		}
		router, err := routerObj.Match(c.ctx.GetRouterPath(), c.ctx.GetMethod())
		if err != nil {
			return ""
		}
		return router.Encoding
	default:
		return ""
	}
}

//GetURL 获取请求路径
func (c *rpath) GetURL() *url.URL {
	return c.ctx.GetURL()
}

//GetRequestPath 获取请求路径
func (c *rpath) GetRequestPath() string {
	return c.ctx.GetURL().Path
}

//Limit 限流设置
func (c *rpath) Limit(isLimit bool, fallback bool) {
	c.isLimit = isLimit
	c.fallback = fallback
}

//IsLimited 是否已限流
func (c *rpath) IsLimited() bool {
	return c.isLimit
}

//AllowFallback 是否允许降级
func (c *rpath) AllowFallback() bool {
	return c.fallback
}
