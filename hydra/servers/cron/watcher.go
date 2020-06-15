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
				w.log.Errorf("加载cron服务失败：%w", err)
				continue
			}
			if cluster.Current().IsBefore(server.Sharding) {
				w.Server.Resume()
				continue
			}
			w.Server.Pause()
		}
	}
}
