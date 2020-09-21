package conf

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/micro-plat/hydra"
	"github.com/micro-plat/hydra/conf/server"
	"github.com/micro-plat/hydra/conf/server/cron"
	"github.com/micro-plat/hydra/conf/server/mqc"
	squeue "github.com/micro-plat/hydra/conf/server/queue"
	"github.com/micro-plat/hydra/conf/server/task"
	"github.com/micro-plat/hydra/conf/vars/cache"
	"github.com/micro-plat/hydra/conf/vars/db"
	"github.com/micro-plat/hydra/conf/vars/queue"
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/global"
	_ "github.com/micro-plat/hydra/registry/registry/localmemory"

	_ "github.com/micro-plat/hydra/hydra/servers/cron"
	_ "github.com/micro-plat/hydra/hydra/servers/http"
	_ "github.com/micro-plat/hydra/hydra/servers/mqc"
)

func StartServer() {
	hydra.Conf.Vars().DB("taosy_db", db.New("oracle", "connstring", db.WithConnect(10, 10, 10)))
	hydra.Conf.Vars().Cache("taosy_cache", cache.New("redis", []byte(`{"cache":"oracle-redis"}`)))
	hydra.Conf.Vars().Queue("taosy_queue", queue.New("redis", []byte(`{"queue":"oracle-redis"}`)))
	hydra.Conf.API(":8070")
	// hydra.Conf.RPC(":8071")
	queues := squeue.Queues{}
	queues.Append(squeue.NewQueue("queuename1", "/testmqc"))
	hydra.Conf.MQC("redis://192.168.0.111:6379", mqc.WithTrace(), mqc.WithTimeout(10)).Queue(queues.Queues...)
	tasks := task.Tasks{}
	tasks.Append(task.NewTask(fmt.Sprintf("@every %ds", 10), "/testcron"),
		task.NewTask(fmt.Sprintf("@every %ds", 10), "/testcron1"),
		task.NewTask(fmt.Sprintf("@every %ds", 10), "/testcron2"),
	)
	hydra.Conf.CRON(cron.WithEnable(), cron.WithTrace(), cron.WithTimeout(10), cron.WithSharding(1)).Task(tasks.Tasks...)
	app := hydra.NewApp(
		hydra.WithServerTypes(global.API),
		hydra.WithPlatName("taosytest"),
		hydra.WithSystemName("test-confcache"),
		hydra.WithClusterName("taosy"),
		//hydra.WithRegistry("lm://localhost"),
	)

	app.API("/testapi", func(ctx context.IContext) (r interface{}) {
		return nil
	})

	os.Args = []string{"taosytest", "run"}
	go app.Start()
	time.Sleep(time.Second * 2)
}

func TestGetCache(t *testing.T) {
	StartServer()
	serverConf, err := server.Cache.GetServerConf(global.API)
	if err != nil {
		t.Error("获取服务配置信息异常")
		return
	}

	router := serverConf.GetRouterConf()
	t.Errorf("api-router信息:%s", router.String())

	varConf, err := server.Cache.GetVarConf()
	if err != nil {
		t.Error("获取var配置信息异常")
		return
	}

	confInfo, err := varConf.GetConf("db", "taosy_db")
	if err != nil {
		t.Errorf("获取var配置信息异常1:%+v", err)
		return
	}

	t.Errorf("数据配置:%v \n", *confInfo)
	b, v, err := confInfo.GetJSON("connString")
	if err != nil {
		t.Errorf("获取var配置信息异常2:%+v", err)
		return
	}

	t.Errorf("db节点数据:%s \n", string(b))
	t.Errorf("db配置版本号:%d \n", v)
	return
}
