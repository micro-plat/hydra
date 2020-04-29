package gin

import (
	"github.com/micro-plat/hydra/engines/context"
	"github.com/micro-plat/hydra/registry/conf"
	"github.com/micro-plat/hydra/registry/conf/server/auth/jwt"
	"github.com/micro-plat/hydra/registry/conf/server/header"
	"github.com/micro-plat/hydra/registry/conf/server/metric"
	"github.com/micro-plat/hydra/registry/conf/server/static"
)

var _ context.Iserver = &server{}

type server struct {
	conf.IMainConf
	path string
}

func newServer() *server {
	return &server{}
}

//GetHeaderConf 获取响应头配置
func (a *server) GetHeaderConf() header.Headers {
	return header.GetConf(a.IMainConf)
}

//GetJWTConf 获取jwt配置
func (a *server) GetJWTConf() *jwt.JWTAuth {
	return jwt.GetConf(a.IMainConf)
}

//GetMetricConf 获取metric配置
func (a *server) GetMetricConf() *metric.Metric {
	return metric.GetConf(a.IMainConf)
}

//GetStaticConf 获取静态文件配置
func (a *server) GetStaticConf() *static.Static {
	return static.GetConf(a.IMainConf)
}
func (a *server) SetServerType(path string) {
	if a.path != "" {
		panic("不能多次设置服务路径:", path)
	}
	a.path = path
	//加载IMainConf
}
