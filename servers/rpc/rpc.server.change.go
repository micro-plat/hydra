package rpc

import (
	"fmt"

	"github.com/micro-plat/hydra/conf"
)

//SetRouters 设置路由配置
func (s *RpcServer) SetRouters(routers []*conf.Router) (err error) {
	if s.Processor, err = s.getProcessor(routers); err != nil {
		return
	}
	return
}

//SetJWT Server
func (s *RpcServer) SetJWT(auth *conf.Auth) error {
	s.conf.SetMetadata("jwt", auth)
	return nil
}

//SetHosts 设置组件的host name
func (s *RpcServer) SetHosts(hosts conf.Hosts) error {
	s.conf.SetMetadata("hosts", hosts)
	return nil
}

//SetTrace 显示跟踪信息
func (s *RpcServer) SetTrace(b bool) {
	s.conf.SetMetadata("show-trace", b)
	return
}

//SetMetric 重置metric
func (s *RpcServer) SetMetric(metric *conf.Metric) error {
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
func (s *RpcServer) StopMetric() error {
	s.metric.Stop()
	return nil
}

//SetHeader 设置http头
func (s *RpcServer) SetHeader(headers conf.Headers) error {
	s.conf.SetMetadata("headers", headers)
	return nil
}
