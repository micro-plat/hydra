package http

import (
	"fmt"

	"github.com/asaskevich/govalidator"
	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/servers/pkg/circuit"
)

//SetRouters 设置路由配置
func (s *ApiServer) SetRouters(routers []*conf.Router) (err error) {
	s.engine.Handler, err = s.getHandler(routers)
	return
}

//SetJWT Server
func (s *ApiServer) SetJWT(auth *conf.Auth) error {
	s.conf.SetMetadata("jwt", auth)
	return nil
}

//SetAjaxRequest 只允许ajax请求
func (s *ApiServer) SetAjaxRequest(allow bool) error {
	s.conf.SetMetadata("ajax-request", allow)
	return nil
}

//SetTrace 显示跟踪信息
func (s *ApiServer) SetTrace(b bool) {
	s.conf.SetMetadata("show-trace", b)
	return
}

//SetHosts 设置组件的host name
func (s *ApiServer) SetHosts(hosts conf.Hosts) error {
	for _, host := range hosts {
		if !govalidator.IsDNSName(host) {
			return fmt.Errorf("%s不是有效的dns名称", host)
		}
	}
	s.conf.SetMetadata("hosts", hosts)
	return nil
}

//SetStatic 设置静态文件路由
func (s *ApiServer) SetStatic(static *conf.Static) error {
	s.conf.SetMetadata("static", static)
	return nil
}

//SetMetric 重置metric
func (s *ApiServer) SetMetric(metric *conf.Metric) error {
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
func (s *ApiServer) SetHeader(headers conf.Headers) error {
	s.conf.SetMetadata("headers", headers)
	return nil
}

//StopMetric stop metric
func (s *ApiServer) StopMetric() error {
	s.metric.Stop()
	return nil
}

//CloseCircuitBreaker 关闭熔断配置
func (s *ApiServer) CloseCircuitBreaker() error {
	if c, ok := s.conf.GetMetadata("__circuit-breaker_").(*circuit.NamedCircuitBreakers); ok {
		c.Close()
	}
	return nil
}

//SetCircuitBreaker 设置熔断配置
func (s *ApiServer) SetCircuitBreaker(c *conf.CircuitBreaker) error {
	s.conf.SetMetadata("__circuit-breaker_", circuit.NewNamedCircuitBreakers(c))
	return nil
}
