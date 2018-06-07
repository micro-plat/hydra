package cron

import (
	"fmt"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/servers"
)

//Notify 服务器配置变更通知
func (w *CronResponsiveServer) Notify(conf conf.IServerConf) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.restarted = false
	//检查是否需要重启服务器
	restart, err := w.NeedRestart(conf)
	if err != nil {
		return err
	}
	if restart { //服务器地址已变化，则重新启动新的server,并停止当前server
		servers.Trace(w.Infof, "关键配置发生变化，准备重启服务器")
		return w.Restart(conf)
	}

	servers.Trace(w.Infof, "配置发生变化，准备更新")
	if err = w.SetConf(false, conf); err != nil {
		return err
	}
	w.engine.UpdateVarConf(conf)
	w.currentConf = conf
	return nil
}

//NeedRestart 检查配置判断是否需要重启服务器
func (w *CronResponsiveServer) NeedRestart(cnf conf.IServerConf) (bool, error) {
	if cnf.ForceRestart() {
		return true, nil
	}
	comparer := conf.NewComparer(w.currentConf, cnf)
	if !comparer.IsChanged() {
		return false, nil
	}
	if comparer.IsVarChanged() {
		return true, nil
	}
	if comparer.IsValueChanged("status", "sharding") {
		return true, nil
	}
	ok, err := comparer.IsRequiredSubConfChanged("task")
	if ok {
		return ok, nil
	}
	if err != nil {
		return false, fmt.Errorf("task未配置或配置有误:%v", err)
	}
	if ok := comparer.IsSubConfChanged("app"); ok {
		return ok, nil
	}
	if ok := comparer.IsSubConfChanged("redis"); ok {
		return true, nil
	}
	return false, nil

}

//SetConf 设置配置参数
func (w *CronResponsiveServer) SetConf(restart bool, conf conf.IServerConf) (err error) {
	//设置分片数量
	w.shardingCount = conf.GetInt("sharding", 0)

	var ok bool
	//设置task
	if restart {
		if _, err := SetTasks(w.engine, w.server, conf, map[string]interface{}{
			"__get_sharding_index_": func() (int, int) {
				return w.shardingIndex, w.shardingCount
			},
		}); err != nil {
			return err
		}
	}
	if err != nil {
		return
	}

	//设置metric
	if ok, err = SetMetric(w.server, conf); err != nil {
		return err
	}
	servers.TraceIf(ok, w.Infof, w.Debugf, getEnableName(ok), "metric设置")
	return nil
}
func getEnableName(b bool) string {
	if b {
		return "启用"
	}
	return "未启用"
}
