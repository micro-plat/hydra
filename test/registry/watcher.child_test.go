package registry

import (
	"testing"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/context/ctx"
	"github.com/micro-plat/hydra/registry/watcher"
	"github.com/micro-plat/hydra/test/assert"
	"github.com/micro-plat/hydra/test/mocks"
	"github.com/micro-plat/lib4go/logger"
)

func TestNewChildWatcher(t *testing.T) {

	confObj := mocks.NewConf()         //构建对象
	confObj.API(":8080")               //初始化参数
	serverConf := confObj.GetAPIConf() //获取配置

	log := logger.GetSession(serverConf.GetMainConf().GetServerName(), ctx.NewUser(&mocks.TestContxt{}, conf.NewMeta()).GetRequestID())

	//测试通过registryAddr构建ChildWatcher
	_, err := watcher.NewChildWatcher("lm://.", []string{}, log)
	assert.Equal(t, nil, err, "测试通过registryAddr构建ChildWatcher")

	//测试通过registry构建ChildWatcher
	r := serverConf.GetMainConf().GetRegistry()
	_, err = watcher.NewChildWatcherByRegistry(r, []string{}, log)
	assert.Equal(t, nil, err, "测试通过registry构建ChildWatcher")

}
