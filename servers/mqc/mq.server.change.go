package mqc

import (
	"fmt"

	"github.com/micro-plat/hydra/conf"
)

//SetMetric 重置metric
func (s *MqcServer) SetMetric(metric *conf.Metric) error {
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
func (s *MqcServer) StopMetric() error {
	s.metric.Stop()
	return nil
}

//SetQueues 设置监听队列
func (s *MqcServer) SetQueues(proto string, raw string, queues []*conf.Queue) (err error) {
	s.Processor, err = s.getProcessor(proto, raw, queues)
	if err != nil {
		err = fmt.Errorf("queue设置有误:%v", err)
		return err
	}
	return nil
}

//SetTrace 显示跟踪信息
func (s *MqcServer) SetTrace(b bool) {
	s.conf.SetMetadata("show-trace", b)
	return
}
