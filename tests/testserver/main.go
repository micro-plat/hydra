package main

import (
	"github.com/micro-plat/hydra"
	"github.com/micro-plat/hydra/components"
	_ "github.com/micro-plat/hydra/components/caches/cache/redis"
	"github.com/micro-plat/hydra/conf/server/task"
	queueredis "github.com/micro-plat/hydra/conf/vars/queue/redis"
	confRedis "github.com/micro-plat/hydra/conf/vars/redis"
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/hydra/servers/cron"
)

var app = hydra.NewApp(
	hydra.WithServerTypes(cron.CRON),
	hydra.WithPlatName("taosytest"),
	hydra.WithSystemName("testserver"),
	hydra.WithClusterName("t"),
	hydra.WithRegistry("lm://."),
)

func init() {
	hydra.Conf.Vars().HTTP("http")
	hydra.Conf.Vars().Redis("redis", confRedis.New(nil, confRedis.WithAddrs("192.168.5.79:6379")))
	hydra.Conf.Vars().Queue().Redis("redis", queueredis.New(queueredis.WithConfigName("redis")))
	hydra.Conf.CRON().Task(task.NewTask("@every 1s", "/taosy/testserver"))
	app.CRON("/taosy/testserver", func(ctx context.IContext) (r interface{}) {
		ctx.Log().Info("cron-hour 接口服务测试")

		queueObj := components.Def.Queue().GetRegularQueue("redis")

		// b := rand.Intn(7) + 3
		for i := 0; i < 10; i++ {
			// go func() {
			httpObj, err := components.Def.HTTP().GetClient()
			if err != nil {
				ctx.Log().Error("获取http客户端异常", err)
				return
			}
			content, status, err := httpObj.Get("http://192.168.0.137:8070/taosy/testapi")
			if err != nil || status != 200 {
				ctx.Log().Errorf("获取http-Get请求异常,status:%d, err:%v", status, err)
				return
			}
			ctx.Log().Info("api-get请求结果:", content)

			content, status, err = httpObj.Post("http://192.168.0.137:8070/taosy/testapi", "")
			if err != nil || status != 200 {
				ctx.Log().Errorf("获取http-Post请求异常,status:%d, err:%v", status, err)
				return
			}
			ctx.Log().Info("api-post请求结果:", content)

			content, status, err = httpObj.Get("http://192.168.0.137:8071/README1.md", "")
			if err != nil || status != 200 {
				ctx.Log().Errorf("获取http-Post请求异常,status:%d, err:%v", status, err)
				return
			}
			ctx.Log().Info("api-static请求结果:", content)

			if err = queueObj.Push("queue1", `{"mqvc":"queue1-succ"}`); err != nil {
				ctx.Log().Errorf("发送queue1队列消息异常,status:%d, err:%v", status, err)
			}

			if err = queueObj.Push("queue2", `{"mqvc":"queue2-succ"}`); err != nil {
				ctx.Log().Errorf("发送queue2队列消息异常,status:%d, err:%v", status, err)
			}
			// }()
		}

		return
	})
}

func main() {
	app.Start()
}
