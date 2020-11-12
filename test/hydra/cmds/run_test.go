package cmds

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/micro-plat/hydra/compatible"

	"github.com/micro-plat/hydra"
	"github.com/micro-plat/hydra/hydra/servers/http"
	"github.com/micro-plat/hydra/test/assert"
)

const runRegistryAddr = "lm://."
const runOthRegistryAddr = "zk://192.168.0.101"

func Test_run_Normal(t *testing.T) {
	//正常的运行
	fileName := fmt.Sprint(time.Now().Nanosecond())
	file, _ := os.Create(fileName)
	orgStd := *os.Stdout
	*os.Stdout = *file

	defer os.Remove(fileName)
	//1. 正常的运行
	args := []string{"xxtest", "run", "-r", runRegistryAddr, "-c", "c"}
	var app = hydra.NewApp(
		hydra.WithServerTypes(http.API),
		hydra.WithPlatName("xxtest"),
		hydra.WithSystemName("apiserver"),
		hydra.WithClusterName("c"),
	)
	hydra.Conf.API(":19003")
	os.Args = args
	go app.Start()
	time.Sleep(time.Second * 2)

	//2. 关闭服务
	compatible.AppClose()

	time.Sleep(time.Second * 2)

	//还原std
	*os.Stdout = orgStd

	file.Close()
	time.Sleep(time.Second)
	bytes, err := ioutil.ReadFile(fileName)

	if err != nil {
		t.Error(err)
		return
	}
	line := string(bytes)

	result := strings.Contains(line, "启动成功")
	assert.Equal(t, true, result, "正常服务启动")

	time.Sleep(time.Second)
}

func Test_run_Normal_other_registry(t *testing.T) {
	//正常的运行
	fileName := fmt.Sprint(time.Now().Nanosecond())
	file, _ := os.Create(fileName)
	orgStd := *os.Stdout
	*os.Stdout = *file

	defer os.Remove(fileName)
	//1. 正常的运行
	args := []string{"xxtest", "run", "-r", runOthRegistryAddr, "-c", "c"}
	var app = hydra.NewApp(
		hydra.WithServerTypes(http.API),
		hydra.WithPlatName("xxtest"),
		hydra.WithSystemName("apiserver"),
		hydra.WithClusterName("c"),
	)
	hydra.Conf.API(":19003")
	os.Args = args
	go app.Start()
	time.Sleep(time.Second * 2)

	//2. 关闭服务
	compatible.AppClose()

	time.Sleep(time.Second * 2)

	//还原std
	*os.Stdout = orgStd

	file.Close()
	time.Sleep(time.Second)
	bytes, err := ioutil.ReadFile(fileName)

	if err != nil {
		t.Error(err)
		return
	}
	line := string(bytes)

	result := strings.Contains(line, "启动成功")
	assert.Equal(t, true, result, "正常服务启动")

	time.Sleep(time.Second)
}

func Test_run_Withtrace(t *testing.T) {
	//启用trace
	fileName := fmt.Sprint(time.Now().Nanosecond())
	file, _ := os.Create(fileName)
	orgStd := *os.Stdout
	*os.Stdout = *file

	defer os.Remove(fileName)
	//1. 正常的运行
	args := []string{"xxtest", "run", "-r", runRegistryAddr, "-c", "c", "-trace", "web"}
	var app = hydra.NewApp(
		hydra.WithServerTypes(http.API),
		hydra.WithPlatName("xxtest"),
		hydra.WithSystemName("apiserver"),
		hydra.WithClusterName("c"),
	)
	hydra.Conf.API(":19004")
	os.Args = args
	go app.Start()
	time.Sleep(time.Second * 2)

	//2. 关闭服务
	compatible.AppClose()

	time.Sleep(time.Second * 2)

	//还原std
	*os.Stdout = orgStd

	file.Close()
	time.Sleep(time.Second)
	bytes, err := ioutil.ReadFile(fileName)

	if err != nil {
		t.Error(err)
		return
	}
	line := string(bytes)

	result := strings.Contains(line, "启动成功:pprof.web")
	assert.Equal(t, true, result, "正常服务启动")

	_, err = os.Stat("trace.out")
	if os.IsNotExist(err) {
		t.Error("trace文件创建失败", err)
		return
	}
	os.Remove("trace.out")

	time.Sleep(time.Second)
}

func Test_run_error_registry(t *testing.T) {
	//错误的注册中心
	fileName := fmt.Sprint(time.Now().Nanosecond())
	file, _ := os.Create(fileName)
	orgStd := *os.Stdout
	*os.Stdout = *file

	//defer os.Remove(fileName)
	//1. 正常的运行
	args := []string{"xxtest", "run", "-r", "xx://xxx", "-c", "c"}
	var app = hydra.NewApp(
		hydra.WithServerTypes(http.API),
		hydra.WithPlatName("xxtest"),
		hydra.WithSystemName("apiserver"),
		hydra.WithClusterName("c"),
	)
	hydra.Conf.API(":19004")
	os.Args = args
	go app.Start()
	time.Sleep(time.Second * 2)

	//2. 关闭服务
	compatible.AppClose()

	time.Sleep(time.Second * 2)

	//还原std
	*os.Stdout = orgStd

	file.Close()
	time.Sleep(time.Second)
	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		t.Error(err)
		return
	}
	line := string(bytes)

	result := strings.Contains(line, "不支持的协议类型")
	assert.Equal(t, true, result, "错误的注册中心")

	time.Sleep(time.Second)
}
