package cron

import (
	"time"

	"github.com/micro-plat/hydra/conf/server/cron"
)

func (w *Responsive) watch() {

	//监控集群信息
START:
	cluster, err := w.conf.GetServerConf().GetCluster()
	if err != nil {
		w.log.Error("获取集群信息失败", err)
		select {
		case <-w.closeChan:
			return
		default:
		}
		time.Sleep(time.Second)
		goto START
	}
	watcher := cluster.Watch()
	notify := watcher.Notify()

	unavailableCount := 0
	//循环监控集群变化
LOOP:
	for {
		select {
		case <-w.closeChan:
			watcher.Close()
			break LOOP
		case <-notify:
			server, err := cron.GetConf(w.conf.GetServerConf())
			if err != nil {
				w.log.Errorf("加载cron配置失败：%w", err)
				continue
			}
			if !cluster.Current().IsAvailable() {
				unavailableCount++
				time.Sleep(500 * time.Millisecond)
				if unavailableCount >= 3 {
					w.log.Warn("cron-当前集群节点不可用")
				}
				continue
			}

			if server.Sharding == 0 || cluster.Current().IsMaster(server.Sharding) {
				ok, err := w.Server.Resume()
				if err != nil {
					w.log.Error("恢复服务器失败:", err)
					continue
				}
				if ok {
					unavailableCount = 0
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
				unavailableCount = 0
				w.update("run-mode", "slave")
				w.log.Debugf("this cron server is started as slave")
			}

		}
	}
}
