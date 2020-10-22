package registry

import (
	"fmt"
	"testing"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/context/ctx"
	"github.com/micro-plat/hydra/registry/pub"
	"github.com/micro-plat/hydra/registry/watcher/wchild"
	"github.com/micro-plat/hydra/test/mocks"
	"github.com/micro-plat/lib4go/logger"
)

func TestChildWatcher_deleted(t *testing.T) {
	//构建配置对象
	confObj := mocks.NewConf()
	confObj.API(":8080")
	apiconf := confObj.GetAPIConf()
	c := apiconf.GetMainConf()
	r := c.GetRegistry()

	//发布节点到注册中心
	p := pub.New(c)
	p.Publish("192.168.5.115:9091", "192.168.5.115:9091", c.GetServerID(), apiconf.GetRouterConf().GetPath()...)
	p.Publish("192.168.5.116:9091", "192.168.5.116:9091", c.GetServerID(), apiconf.GetRouterConf().GetPath()...)

	//path := c.GetServerPubPath(c.GetClusterName())

	//子节点监控启动
	log := logger.GetSession(apiconf.GetMainConf().GetServerName(), ctx.NewUser(&mocks.TestContxt{}, conf.NewMeta()).GetRequestID())
	w := wchild.NewChildWatcher(r, "hydra/apiserver/api/test/", log)
	w.Start() //启动更新w.watchers

	fmt.Println("w:", w.Watchers)

	// //模拟注册中心节点不存在
	// err := r.Delete(path)
	// assert.Equal(t, nil, err, "删除节点")

	// //启动另外一个线程检查不存在的节点并进行删除操作
	// w.Start()

	// //保证能够执行删除操作
	// time.Sleep(time.Second * 2)

	// fmt.Println("w:", w.Watchers)

	//fmt.Println("children:", paths)

	//测试的最底层节点变动
}
