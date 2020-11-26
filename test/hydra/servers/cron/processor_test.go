package cron

import (
	"fmt"
	"testing"
	"time"

	"github.com/micro-plat/hydra/global"

	"github.com/micro-plat/hydra/components"
	_ "github.com/micro-plat/hydra/components/caches/cache/redis"
	"github.com/micro-plat/hydra/conf/app"
	"github.com/micro-plat/hydra/conf/server/task"
	  "github.com/micro-plat/hydra/conf/vars/cache/cacheredis"
	 "github.com/micro-plat/hydra/conf/vars/queue/queueredis"
	confRedis "github.com/micro-plat/hydra/conf/vars/redis"
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/hydra/servers/cron"
	"github.com/micro-plat/hydra/services"
	"github.com/micro-plat/hydra/test/assert"
	"github.com/micro-plat/hydra/test/mocks"
)

func TestProcessor_Start(t *testing.T) {
	confObj := mocks.NewConfBy("cronserver_resserivece_testx", "testcronsdfx")
	global.Def.PlatName = "taosytest"
	confObj.CRON()
	confObj.Vars().Redis("redis", confRedis.New(nil, confRedis.WithAddrs("192.168.5.79:6379")))
	confObj.Vars().Cache().Redis("redis", cacheredis.New(cacheredis.WithConfigName("redis")))
	confObj.Vars().Queue().Redis("redis", queueredis.New(queueredis.WithConfigName("redis")))
	app.Cache.Save(confObj.GetCronConf())
	services.Def.CRON("/taosy/services1", func(ctx context.IContext) (r interface{}) {
		queueObj := components.Def.Queue().GetRegularQueue("redis")
		if err := queueObj.Push("services1:queue1", `1`); err != nil {
			ctx.Log().Errorf("发送queue1队列消息异常, err:%v", err)
		}
		return
	})

	services.Def.CRON("/taosy/services2", func(ctx context.IContext) (r interface{}) {
		queueObj := components.Def.Queue().GetRegularQueue("redis")
		if err := queueObj.Push("services2:queue2", `1`); err != nil {
			ctx.Log().Errorf("发送queue1队列消息异常, err:%v", err)
		}
		return
	})

	services.Def.CRON("/taosy/services3", func(ctx context.IContext) (r interface{}) {
		queueObj := components.Def.Queue().GetRegularQueue("redis")
		if err := queueObj.Push("services3:queue3", `1`); err != nil {
			ctx.Log().Errorf("发送queue1队列消息异常, err:%v", err)
		}
		return
	})

	services.Def.CRON("/taosy/services4", func(ctx context.IContext) (r interface{}) {
		queueObj := components.Def.Queue().GetRegularQueue("redis")
		if err := queueObj.Push("services4:queue4", `1`); err != nil {
			ctx.Log().Errorf("发送queue1队列消息异常, err:%v", err)
		}
		return
	})

	s := cron.NewProcessor()
	test1 := task.NewTask("@every 1s", "/taosy/services1")
	test2 := task.NewTask("@every 5s", "/taosy/services2")
	test3 := task.NewTask("@every 10s", "/taosy/services3")
	test4 := task.NewTask("@every 40s", "/taosy/services4")
	err := s.Add(test1, test2, test3, test4)
	assert.Equalf(t, true, err == nil, ",err")
	s.Resume()
	go s.Start()
	time.Sleep(51 * time.Second)
	s.Close()

	queueObj := components.Def.Queue().GetRegularQueue("redis")
	count1, _ := queueObj.Count("services1:queue1")
	count2, _ := queueObj.Count("services2:queue2")
	count3, _ := queueObj.Count("services3:queue3")
	count4, _ := queueObj.Count("services4:queue4")
	assert.Equal(t, int64(50), count1, fmt.Sprint("1s任务数量", count1))
	assert.Equal(t, int64(10), count2, fmt.Sprint("5s任务数量", count2))
	assert.Equal(t, int64(5), count3, fmt.Sprint("10s任务数量", count3))
	assert.Equal(t, int64(1), count4, fmt.Sprint("40s任务数量", count4))
	cacheObj := components.Def.Cache().GetRegularCache("redis")
	cacheObj.Delete("taosytest:services1:queue1")
	cacheObj.Delete("taosytest:services2:queue2")
	cacheObj.Delete("taosytest:services3:queue3")
	cacheObj.Delete("taosytest:services4:queue4")
}
