package container

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/micro-plat/hydra"
	"github.com/micro-plat/hydra/components/container"
	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/conf/server/api"
	"github.com/micro-plat/hydra/conf/server/apm"
	"github.com/micro-plat/hydra/context"
	xhttp "github.com/micro-plat/hydra/hydra/servers/http"
	_ "github.com/micro-plat/hydra/registry/registry/localmemory"
)

func StartServer() {

	//hydra.Conf.Vars().APM(varapm.New("skywalking", `{"server_address":"192.168.106.160:11800"}`))
	hydra.Conf.API(":8070", api.WithHeaderReadTimeout(30), api.WithTimeout(30, 30)).APM("skywalking", apm.WithDisable())
	app := hydra.NewApp(
		hydra.WithServerTypes(xhttp.API),
		hydra.WithPlatName("test"),
		hydra.WithSystemName("test-performance"),
		hydra.WithClusterName("t"),
		hydra.WithRegistry("lm://."),
	)

	app.API("/entrance", func(ctx context.IContext) (r interface{}) {
		fmt.Println("api:inbound")
		return nil
	})

	os.Args = []string{"test", "run"}
	go app.Start()
	time.Sleep(time.Second * 2)
}

func TestGetOrCreate(t *testing.T) {
	defer func() {
		if e := recover(); e != nil {
			err := e.(error)
			if !strings.Contains(err.Error(), "未找到var") {
				t.Errorf("var未配置用例检测不通过;err:%+v", err)
				return
			}
		}
	}()
	c := container.NewContainer()
	_, err := c.GetOrCreate("gocache", "cache", func(js *conf.RawConf) (interface{}, error) {
		return nil, nil
	})
	if err != nil && !strings.Contains(err.Error(), "未找到var") {
		t.Errorf("var主节点检测不通过:err:%+v", err)
		return
	}
	return
}

func TestGetOrCreate1(t *testing.T) {
	StartServer()
	c := container.NewContainer()
	_, err := c.GetOrCreate("gocache", "cache", func(js *conf.RawConf) (interface{}, error) {
		if string(js.GetRaw()) == "{}" {
			return nil, fmt.Errorf("cache节点不存在")
		}

		return nil, nil
	})
	if err != nil && !strings.Contains(err.Error(), "cache节点不存在") {
		t.Errorf("var主节点检测不通过:err:%+v", err)
		return
	}

	if err == nil {
		t.Errorf("var主节点检测不通过11")
		return
	}
	return
}

func TestGetOrCreate2(t *testing.T) {
	StartServer()
	c := container.NewContainer()
	_, err := c.GetOrCreate("gocache", "cache", func(js *conf.RawConf) (interface{}, error) {
		if string(js.GetRaw()) == "{}" {
			return nil, fmt.Errorf("cache节点不存在")
		}

		return nil, nil
	})
	if err != nil && !strings.Contains(err.Error(), "cache节点不存在") {
		t.Errorf("var主节点检测不通过:err:%+v", err)
		return
	}

	if err == nil {
		t.Errorf("var主节点检测不通过11")
		return
	}
	return
}
