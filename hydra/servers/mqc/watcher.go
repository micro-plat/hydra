package mqc

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
			server := w.conf.GetMQCMainConf()
			if !cluster.Current().IsAvailable() {
				w.log.Error("当前集群节点不可用")
				continue
			}

			if server.Sharding == 0 || cluster.Current().IsMaster(server.Sharding) {
				ok, err := w.Server.Resume()
				if err != nil {
					w.log.Error("恢复mqc服务器失败:", err)
					continue
				}
				if ok {
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
				w.update("run-mode", "slave")
				w.log.Debugf("this mqc server is started as slave")
			}
		}
	}
}
