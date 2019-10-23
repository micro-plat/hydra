package http

import (
	"fmt"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/servers/pkg/circuit"
)

//SetRouters 设置路由配置
func (s *WebServer) SetRouters(routers []*conf.Router) (err error) {
	s.engine.Handler, err = s.getHandler(routers)
	return
}

//SetJWT Server
func (s *WebServer) SetJWT(auth *conf.JWTAuth) error {
	s.conf.SetMetadata("jwt", auth)
	return nil
}

//SetResponse 设置response配置
func (s *WebServer) SetResponse(r *conf.Response) error {
	s.conf.SetMetadata("__response_conf_", r)
	return nil
}

//SetAjaxRequest 只允许ajax请求
func (s *WebServer) SetAjaxRequest(allow bool) error {
	s.conf.SetMetadata("ajax-request", allow)
	return nil
}

//SetTrace 显示跟踪信息
func (s *WebServer) SetTrace(b bool) {
	s.conf.SetMetadata("show-trace", b)
	return
}

//SetHosts 设置组件的host name
func (s *WebServer) SetHosts(hosts conf.Hosts) error {
	s.conf.SetMetadata("hosts", hosts)
	return nil
}

//SetStatic 设置静态文件路由
func (s *WebServer) SetStatic(static *conf.Static) error {
	s.conf.SetMetadata("static", static)
	return nil
}

//SetMetric 重置metric
func (s *WebServer) SetMetric(metric *conf.Metric) error {
	s.metric.Stop()
	if metric.Disable {
		return nil
	}
	if err := s.metric.Restart(metric.Host, metric.DataBase, metric.UserName, metric.Password, metric.Cron, s.Logger); err != nil {
		err = fmt.Errorf("metric设置有误:%v", err)
		return err
	}
	return nil
}

//SetHeader 设置http头
func (s *WebServer) SetHeader(headers conf.Headers) error {
	s.conf.SetMetadata("headers", headers)
	return nil
}

//StopMetric stop metric
func (s *WebServer) StopMetric() error {
	s.metric.Stop()
	return nil
}

//SetView 设置view参数
func (s *WebServer) SetView(view *conf.View) (err error) {
	s.conf.SetMetadata("view", view)
	if s.views, err = s.loadHTMLGlob(); err != nil {
		s.Logger.Debugf("%s未找到模板:%v", s.conf.Name, err)
		return err
	}
	return nil
}

//CloseCircuitBreaker 关闭熔断配置
func (s *WebServer) CloseCircuitBreaker() error {
	if c, ok := s.conf.GetMetadata("__circuit-breaker_").(*circuit.NamedCircuitBreakers); ok {
		c.Close()
	}
	return nil
}

//SetCircuitBreaker 设置熔断配置
func (s *WebServer) SetCircuitBreaker(c *conf.CircuitBreaker) error {
	s.conf.SetMetadata("__circuit-breaker_", circuit.NewNamedCircuitBreakers(c))
	return nil
}
