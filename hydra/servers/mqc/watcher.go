package mqc

import (
	"time"
)

func (w *Responsive) watch() {
START:
	cluster, err := w.conf.GetServerConf().GetCluster()
	if err != nil {
		w.log.Errorf("无法获取到集群信息：%v", err)
		tk := time.After(time.Second)
		select {
		case <-tk:
			goto START
		case <-w.closeChan:
			return
		}
	}

	watcher := cluster.Watch()
	notify := watcher.Notify()
	unavailableCount := 0

LOOP:
	for {
		select {
		case <-w.closeChan:
			watcher.Close()
			break LOOP
		case <-notify:
			server, err := w.conf.GetMQCMainConf()
			if err != nil {
				w.log.Error("mqc主配置获取失败:", err)
				continue
			}
			if !cluster.Current().IsAvailable() {
				unavailableCount++
				time.Sleep(500 * time.Millisecond)
				if unavailableCount >= 3 {
					w.log.Warn("当前集群节点不可用")
				}
				continue
			}

			if server.Sharding == 0 || cluster.Current().IsMaster(server.Sharding) {
				ok, err := w.Server.Resume()
				if err != nil {
					w.log.Error("恢复mqc服务器失败:", err)
					continue
				}
				if ok {
					unavailableCount = 0
					w.update("run-mode", "master")
					w.log.Debugf("this mqc server is started as master")
				}
				continue
			}
			ok, err := w.Server.Pause()
			if err != nil {
				w.log.Error("暂停mqc服务器失败:", err)
				continue
			}
			if ok {
				unavailableCount = 0
				w.update("run-mode", "slave")
				w.log.Debugf("this mqc server is started as slave")
			}
		}
	}
}
