package cmds

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/micro-plat/hydra"
	varredis "github.com/micro-plat/hydra/conf/vars/redis"
	"github.com/micro-plat/hydra/hydra/servers/http"
	"github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/hydra/test/assert"
	"github.com/micro-plat/lib4go/logger"
)

const registryAddr = "zk://192.168.0.101"

//@todo show 会无限循环写入
func xTestshowNow(t *testing.T) {
	fileName := fmt.Sprint(time.Now().Nanosecond())
	file, _ := os.Create(fileName)
	orgStd := *os.Stdout
	*os.Stdout = *file
	//show
	args := []string{"xxtest", "conf", "show", "-r", registryAddr, "-c", "c"}

	var app = hydra.NewApp(
		hydra.WithServerTypes(http.API),
		hydra.WithPlatName("xxtest"),
		hydra.WithSystemName("apiserver"),
		hydra.WithClusterName("c"),
	)
	os.Args = args
	app.Start()
	time.Sleep(time.Second)

	//还原std
	*os.Stdout = orgStd

	file.Close()
	time.Sleep(time.Second)
	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		orgStd.WriteString(err.Error())
	}
	lines := strings.Split(string(bytes), "\r")
	for _, row := range lines {
		assert.Equal(t, true, strings.Contains(row, "OK"), "正常参数的安装")
	}
}

func Test_installNow_Normal(t *testing.T) {
	fileName := fmt.Sprint(time.Now().Nanosecond())
	file, _ := os.Create(fileName)
	orgStd := *os.Stdout
	*os.Stdout = *file
	defer os.Remove(fileName)
	//正常的安装
	args := []string{"xxtest", "conf", "install", "-r", registryAddr, "-c", "c"}

	var app = hydra.NewApp(
		hydra.WithServerTypes(http.API),
		hydra.WithPlatName("xxtest"),
		hydra.WithSystemName("apiserver"),
		hydra.WithClusterName("c"),
	)
	os.Args = args
	app.Start()
	time.Sleep(time.Second)

	//还原std
	*os.Stdout = orgStd

	file.Close()
	time.Sleep(time.Second)
	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		orgStd.WriteString(err.Error())
	}
	lines := strings.Split(string(bytes), "\r")
	for _, row := range lines {
		assert.Equal(t, true, strings.Contains(row, "OK"), "正常参数的安装")
	}
}

func Test_installNow_NoFlag(t *testing.T) {
	fileName := fmt.Sprint(time.Now().Nanosecond())
	file, _ := os.Create(fileName)
	orgStd := *os.Stdout
	*os.Stdout = *file
	defer os.Remove(fileName)
	//正常的安装
	args := []string{"xxtest", "conf", "install", "-r", registryAddr, "-xc", "c"}

	var app = hydra.NewApp(
		hydra.WithServerTypes(http.API),
		hydra.WithPlatName("xxtest"),
		hydra.WithSystemName("apiserver"),
		hydra.WithClusterName("c"),
	)

	os.Args = args
	app.Start()
	time.Sleep(time.Second)

	//还原std
	*os.Stdout = orgStd

	file.Close()
	time.Sleep(time.Second)
	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		orgStd.WriteString(err.Error())
	}
	lines := strings.Split(string(bytes), "\r")
	for _, row := range lines {
		assert.Equal(t, true, strings.Contains(row, "Incorrect Usage"), "未定义的flag")
	}
}

func Test_installNow_Cover(t *testing.T) {
	fileName := fmt.Sprint(time.Now().Nanosecond())
	file, _ := os.Create(fileName)
	orgStd := *os.Stdout
	*os.Stdout = *file

	defer os.Remove(fileName)

	orgRedisAddr := "192.168.5.79:1000"
	newRedisAddr := "192.168.5.79:6379"

	//正常的安装
	args := []string{"xxtest", "conf", "install", "-r", registryAddr, "-c", "c"}
	var app = hydra.NewApp(
		hydra.WithServerTypes(http.API),
		hydra.WithPlatName("xxtest"),
		hydra.WithSystemName("apiserver"),
		hydra.WithClusterName("c"),
	)

	hydra.Conf.Vars().Redis("5.79", varredis.New([]string{orgRedisAddr}))

	os.Args = args
	app.Start()
	time.Sleep(time.Second)

	//执行覆盖
	args = []string{"xxtest", "conf", "install", "-r", registryAddr, "-c", "c", "-cover", "true"}
	app = hydra.NewApp(
		hydra.WithServerTypes(http.API),
		hydra.WithPlatName("xxtest"),
		hydra.WithSystemName("apiserver"),
		hydra.WithClusterName("c"),
	)

	hydra.Conf.Vars().Redis("5.79", varredis.New([]string{newRedisAddr}))

	os.Args = args
	app.Start()
	time.Sleep(time.Second)

	//还原std
	*os.Stdout = orgStd

	regist, err := registry.NewRegistry(registryAddr, logger.Nil())
	if err != nil {
		t.Error("registry.NewRegistry", err)
		return
	}

	data, _, err := regist.GetValue("/xxtest/var/redis/5.79")

	if err != nil {
		t.Error("regist.GetValue", err)
		return
	}

	if !strings.Contains(string(data), newRedisAddr) {
		t.Error("数据覆盖不成功", err)
		return
	}

	file.Close()
	time.Sleep(time.Second)
	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		orgStd.WriteString(err.Error())
	}
	lines := strings.Split(string(bytes), "\r")
	for _, row := range lines {
		if strings.Contains(row, "安装到配置中心") {
			assert.Equal(t, true, strings.Contains(row, "OK"), "带Cover参数")
		}
	}

}
