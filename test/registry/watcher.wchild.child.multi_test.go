package registry

import (
	"testing"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/context/ctx"
	"github.com/micro-plat/hydra/registry/watcher/wchild"
	"github.com/micro-plat/hydra/test/assert"
	"github.com/micro-plat/hydra/test/mocks"
	"github.com/micro-plat/lib4go/logger"
)

func TestNewMultiChildWatcher(t *testing.T) {

	//构建配置对象
	confObj := mocks.NewConfBy("hydra_rgst_watcher_MultiChild", "rgtwatchMultiChildest")
	confObj.API(":8080")
	apiconf := confObj.GetAPIConf()
	c := apiconf.GetServerConf()
	log := logger.GetSession(apiconf.GetServerConf().GetServerName(), ctx.NewUser(&mocks.TestContxt{}, "", conf.NewMeta()).GetRequestID())

	w, _ := wchild.NewMultiChildWatcher(c.GetRegistry(), []string{"a", "b", "c"}, log)
	assert.Equal(t, 3, len(w.Watchers), "构建节点监控对象")
}

func TestMultiChildWatcher_Close(t *testing.T) {
	//构建配置对象
	confObj := mocks.NewConfBy("hydra_rgst_watcher_MultiChild1", "rgtwatchMultiChildest1")
	confObj.API(":8080")
	apiconf := confObj.GetAPIConf()
	c := apiconf.GetServerConf()
	log := logger.GetSession(apiconf.GetServerConf().GetServerName(), ctx.NewUser(&mocks.TestContxt{}, "", conf.NewMeta()).GetRequestID())

	w, _ := wchild.NewMultiChildWatcher(c.GetRegistry(), []string{"a", "b", "c"}, log)
	w.Close()
	for _, v := range w.Watchers {
		assert.Equal(t, true, v.Done, "节点监控对象关闭")
		_, ok := <-v.CloseChan
		assert.Equal(t, false, ok, "节点监控对象关闭")
	}
}
