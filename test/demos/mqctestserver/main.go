package main

import (
	"github.com/micro-plat/hydra"
	"github.com/micro-plat/hydra/components"
	_ "github.com/micro-plat/hydra/components/caches/cache/gocache"
	_ "github.com/micro-plat/hydra/components/caches/cache/memcached"
	_ "github.com/micro-plat/hydra/components/caches/cache/redis"
	"github.com/micro-plat/lib4go/types"

	_ "github.com/micro-plat/hydra/components/queues/mq/lmq"
	_ "github.com/micro-plat/hydra/components/queues/mq/mqtt"
	_ "github.com/micro-plat/hydra/components/queues/mq/redis"
	_ "github.com/micro-plat/hydra/components/queues/mq/xmq"

	"github.com/micro-plat/hydra/conf/server/task"

	"github.com/micro-plat/hydra/conf/vars/queue/queueredis"
	confRedis "github.com/micro-plat/hydra/conf/vars/redis"
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/hydra/servers/cron"
)

var app = hydra.NewApp(
	hydra.WithServerTypes(cron.CRON),
	hydra.WithPlatName("hydratest"),
	hydra.WithSystemName("mqctestserver"),
	hydra.WithClusterName("t"),
	hydra.WithRegistry("zk://192.168.0.101"),
)

func init() {
	hydra.Conf.Vars().Redis("redis", "192.168.5.79:6379", confRedis.WithTimeout(10, 10, 10))
	hydra.Conf.Vars().Queue().Redis("redis", "", queueredis.WithConfigName("redis"))
	hydra.Conf.CRON().Task(task.NewTask("@every 1s", "/taosy/testserver"))
	app.CRON("/taosy/testserver", func(ctx context.IContext) (r interface{}) {
		ctx.Log().Info("mqc-test 消息丄生产，每秒1000个消息")
		queueObj := components.Def.Queue().GetRegularQueue("redis")
		for i := 0; i < 1000; i++ {
			queueName := "mqcserver:yltest" + types.GetString(i)
			if err := queueObj.Send(queueName, `{"mqvc":"queue1-succ"}`); err != nil {
				ctx.Log().Errorf("发送queue1队列消息异常,err:%v", err)
			}
		}
		return
	})
}

func main() {
	app.Start()
}
