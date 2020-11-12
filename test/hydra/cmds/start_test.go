package cmds

import (
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/micro-plat/hydra"
	"github.com/micro-plat/hydra/hydra/servers/http"
	"github.com/micro-plat/hydra/test/assert"
)

func Test_start_Normal(t *testing.T) {
	resetServiceName(t.Name())
	execPrint(t)
	//正常的开启
	fileName := fmt.Sprint(time.Now().Nanosecond())
	file, _ := os.Create(fileName)
	orgStd := *os.Stdout
	*os.Stdout = *file
	defer os.Remove(fileName)

	//1. 安装服务
	args := []string{"xxtest", "install", "-r", runRegistryAddr, "-c", "c"}
	var app = hydra.NewApp(
		hydra.WithServerTypes(http.API),
		hydra.WithPlatName("xxtest"),
		hydra.WithSystemName("apiserver"),
		hydra.WithClusterName("c"),
	)
	hydra.Conf.API(":19010")
	os.Args = args
	go app.Start()
	time.Sleep(time.Second * 2)

	//2. 启动服务
	args = []string{"xxtest", "start"}
	os.Args = args
	go app.Start()

	time.Sleep(time.Second * 10)

	//3. 清除服务
	args = []string{"xxtest", "remove"}
	os.Args = args
	go app.Start()

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
	//启动 + 成功
	result := strings.Contains(line, "Starting") && strings.Contains(line, "OK")
	assert.Equal(t, true, result, "正常服务启动")

	time.Sleep(time.Second)
}

//sudo /usr/local/go/bin/go test github.com/micro-plat/hydra/test/hydra/cmds -run "^(Test_start_Not_installed)$"  -v -timeout=60s
func Test_start_Not_installed(t *testing.T) {
	resetServiceName(t.Name())
	//启动未安装的服务
	execPrint(t)

	fileName := fmt.Sprint(time.Now().Nanosecond())
	file, _ := os.Create(fileName)
	orgStd := *os.Stdout
	*os.Stdout = *file
	//defer os.Remove(fileName)

	//1. 启动服务

	var app = hydra.NewApp(
		hydra.WithServerTypes(http.API),
		hydra.WithPlatName("xxtest"),
		hydra.WithSystemName("apiserver"),
		hydra.WithClusterName("c"),
	)
	hydra.Conf.API(":19011")

	//2. 清除服务(保证没有服务安装)
	os.Args = []string{"xxtest", "remove"}
	go app.Start()
	time.Sleep(time.Second * 2)

	//2. 启动服务
	os.Args = []string{"xxtest", "start"}
	go app.Start()
	time.Sleep(time.Second * 2)

	//2. 删除服务 
	os.Args =  []string{"xxtest", "remove"}
	app.Start()
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

	if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
		//unbuntu/centos
		result := strings.Contains(line, "sudo") || strings.Contains(line, "Service is not installed")
		assert.Equal(t, true, result, "启动未安装的服务")
	}
	if runtime.GOOS == "windows" {
		result := strings.Contains(line, "not exist as an installed service")
		assert.Equal(t, true, result, "启动未安装的服务")
	}

	time.Sleep(time.Second)
}
