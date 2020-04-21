package cron

import (
	"strings"

	"github.com/micro-plat/hydra/servers"
	"github.com/micro-plat/lib4go/types"
)

//publish 将当前服务器的节点信息发布到注册中心
func (w *CronResponsiveServer) publish() (err error) {
	addr := w.server.GetAddress()
	serverName := strings.Split(addr, "://")[1]

	if err := w.pub.Publish(serverName, map[string]interface{}{
		"service":    addr,
		"cluster_id": w.currentConf.GetClusterID(),
	}); err != nil {
		return err
	}

	return
}
func (w *CronResponsiveServer) notify(isMaster bool, sharding int, total int) {
	servers.Tracef(w.Infof, "%s", types.DecodeString(isMaster, true, "master cron server", "slave cron server"))
	if isMaster {
		w.server.Resume()
	}
	w.server.Pause()
}

//unpublish 删除已发布的节点
func (w *CronResponsiveServer) unpublish() {
	w.pub.Clear()
}
