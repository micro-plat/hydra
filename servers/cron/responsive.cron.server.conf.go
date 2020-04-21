package cron

import (
	"github.com/micro-plat/hydra/registry/conf"
	"github.com/micro-plat/hydra/registry/conf/server/cron"
	"github.com/micro-plat/hydra/servers"
)

//Notify 服务器配置变更通知
func (w *CronResponsiveServer) Notify(c conf.IMainConf) error {
	w.comparer.Update(c)

	//配置未发生变化
	if w.comparer.IsChanged() {
		return nil
	}

	if w.comparer.IsValueChanged(cron.MainConfName...) ||
		w.comparer.IsSubConfChanged(cron.SubConfName...) {
		servers.Trace(w.Infof, "关键配置发生变化，准备重启服务器")
		return w.Restart(c)
	}
	servers.Trace(w.Infof, "配置发生变化，准备更新")
	return nil
}
