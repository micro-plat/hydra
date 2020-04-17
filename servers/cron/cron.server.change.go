package cron

import (
	"fmt"

	"github.com/micro-plat/hydra/conf"
)

//SetMetric 重置metric
func (s *CronServer) SetMetric(metric *conf.Metric) error {
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

//SetTasks 设置定时任务
func (s *CronServer) SetTasks(tasks []*conf.Task) (err error) {
	s.Processor, err = s.getProcessor(tasks)
	return err
}

//ShowTrace 显示跟踪信息
func (s *CronServer) ShowTrace(b bool) {
	s.conf.SetMetadata("show-trace", b)
	return
}
