package cron

import "github.com/micro-plat/hydra/conf/server/cron"

func (w *Responsive) watch() {
	cluster := w.conf.GetMainConf().GetCluster()
	watcher := cluster.Watch()
	notify := watcher.Notify()
LOOP:
	for {
		select {
		case <-w.closeChan:
			watcher.Close()
			break LOOP
		case <-notify:

			server, err := cron.GetConf(w.conf.GetMainConf())
			if err != nil {
				w.log.Errorf("加载cron配置失败：%w", err)
				continue
			}
			if !cluster.Current().IsAvailable() {
				w.log.Error("当前集群节点不可用")
				continue
			}

			if server.Sharding == 0 || cluster.Current().IsMaster(server.Sharding) {
				ok, err := w.Server.Resume()
				if err != nil {
					w.log.Error("恢复服务器失败:", err)
					continue
				}
				if ok {
					w.update("run-mode", "master")
					w.log.Debugf("this cron server is started as master")
				}
				continue
			}
			ok, err := w.Server.Pause()
			if err != nil {
				w.log.Error("暂停服务器失败:", err)
				continue
			}
			if ok {
				w.update("run-mode", "slave")
				w.log.Debugf("this cron server is started as slave")
			}
		}
	}
}
