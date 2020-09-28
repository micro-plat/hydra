package main

import (
	"fmt"

	"github.com/micro-plat/hydra"
	"github.com/micro-plat/hydra/conf/server/cron"
	"github.com/micro-plat/hydra/conf/server/task"
	"github.com/micro-plat/hydra/context"
	scron "github.com/micro-plat/hydra/hydra/servers/cron"
)

/*
	1.conf配置的status=stop无效;
	2.配置时,任务中的disable值不能设置;
	3.超时时间无作用;
*/

var app = hydra.NewApp(
	hydra.WithServerTypes(scron.CRON),
	hydra.WithPlatName("taosy-cron-test"),
	hydra.WithSystemName("test-cronn"),
	hydra.WithClusterName("taosycron"),
	hydra.WithRegistry("redis://192.168.5.79:6379"),
	// hydra.WithRegistry("redis://192.168.0.111:6379,192.168.0.112:6379,192.168.0.113:6379,192.168.0.114:6379,192.168.0.115:6379,192.168.0.116:6379"),
	// hydra.WithRegistry("zk://192.168.0.101:2181"),
)

func init() {
	tasks := task.Tasks{}
	tasks.Append(task.NewTask(fmt.Sprintf("@every %ds", 10), "/taosy/testcron"))
	// tasks.Append(task.NewTask("@once", "/taosy/testcron"))
	hydra.Conf.CRON(cron.WithEnable(), cron.WithTrace(), cron.WithTimeout(10), cron.WithDisable()).Task(tasks.Tasks...)
	app.CRON("/taosy/testcron", func(ctx context.IContext) (r interface{}) {
		ctx.Log().Info("CRON-CRON-CRON-CRON-CRON-CRON")

		return nil
	})
}

func main() {
	app.Start()
}
