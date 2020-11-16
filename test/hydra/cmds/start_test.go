package cmds

import (
	"fmt"
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

	defunc, fileCallback := injectStdOutFile()
	defer defunc()

	//1. 安装服务
	var app = hydra.NewApp(
		hydra.WithServerTypes(http.API),
		hydra.WithPlatName("xxtest"),
		hydra.WithSystemName("apiserver"),
		hydra.WithClusterName("c"),
	)
	hydra.Conf.API(":19010")
	fmt.Println("安装")
	os.Args = []string{"xxtest", "install", "-r", runRegistryAddr, "-c", "c"}
	go app.Start()
	time.Sleep(time.Second * 4)
	fmt.Println("启动")
	//2. 启动服务
	os.Args = []string{"xxtest", "start"}
	app.Start()

	time.Sleep(time.Second * 2)

	fmt.Println("停止")
	//3. 清除服务
	os.Args = []string{"xxtest", "stop"}
	app.Start()

	fmt.Println("删除")
	//3. 清除服务
	os.Args = []string{"xxtest", "remove"}
	app.Start()

	time.Sleep(time.Second * 2)
	bytes, err := fileCallback()
	if err != nil {
		t.Error(err)
		return
	}
	line := string(bytes)
	//启动 + 成功
	result := strings.Contains(line, "Starting") && strings.Contains(line, "OK")
	assert.Equal(t, true, result, "正常-服务start启动")

	time.Sleep(time.Second)
}

//sudo /usr/local/go/bin/go test github.com/micro-plat/hydra/test/hydra/cmds -run "^(Test_start_Not_installed)$"  -v -timeout=60s
func Test_start_Not_installed(t *testing.T) {
	resetServiceName(t.Name())
	//启动未安装的服务
	execPrint(t)

	defunc, fileCallback := injectStdOutFile()
	defer defunc()

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
	os.Args = []string{"xxtest", "remove"}
	app.Start()

	time.Sleep(time.Second * 2)
	bytes, err := fileCallback()

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
