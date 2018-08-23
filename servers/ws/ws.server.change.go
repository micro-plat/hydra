package ws

import (
	"fmt"

	"github.com/asaskevich/govalidator"
	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/servers/pkg/circuit"
)

//SetRouters 设置路由配置
func (s *WSServer) SetRouters(routers []*conf.Router) (err error) {
	s.engine.Handler, err = s.getHandler(routers)
	return
}

//SetJWT Server
func (s *WSServer) SetJWT(auth *conf.Auth) error {
	s.conf.SetMetadata("jwt", auth)
	return nil
}

//SetHosts 设置组件的host name
func (s *WSServer) SetHosts(hosts conf.Hosts) error {
	for _, host := range hosts {
		if !govalidator.IsDNSName(host) {
			return fmt.Errorf("%s不是有效的dns名称", host)
		}
	}
	s.conf.SetMetadata("hosts", hosts)
	return nil
}

//SetStatic 设置静态文件路由
func (s *WSServer) SetStatic(static *conf.Static) error {
	s.conf.SetMetadata("static", static)
	return nil
}

//SetTrace 显示跟踪信息
func (s *WSServer) SetTrace(b bool) {
	s.conf.SetMetadata("show-trace", b)
	return
}

//SetMetric 重置metric
func (s *WSServer) SetMetric(metric *conf.Metric) error {
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

//StopMetric stop metric
func (s *WSServer) StopMetric() error {
	s.metric.Stop()
	return nil
}

//CloseCircuitBreaker 关闭熔断配置
func (s *WSServer) CloseCircuitBreaker() error {
	if c, ok := s.conf.GetMetadata("__circuit-breaker_").(*circuit.NamedCircuitBreakers); ok {
		c.Close()
	}
	return nil
}

//SetCircuitBreaker 设置熔断配置
func (s *WSServer) SetCircuitBreaker(c *conf.CircuitBreaker) error {
	s.conf.SetMetadata("__circuit-breaker_", circuit.NewNamedCircuitBreakers(c))
	return nil
}
