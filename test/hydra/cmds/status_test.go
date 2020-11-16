package cmds

import (
	"os"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/micro-plat/hydra"
	"github.com/micro-plat/hydra/hydra/servers/http"
	"github.com/micro-plat/hydra/test/assert"
)

func Test_status_Normal_running(t *testing.T) {
	resetServiceName(t.Name())
	execPrint(t)
	//正常
	defunc, fileCallback := injectStdOutFile()
	defer defunc()

	var app = hydra.NewApp(
		hydra.WithServerTypes(http.API),
		hydra.WithPlatName("xxtest"),
		hydra.WithSystemName("apiserver"),
		hydra.WithClusterName("c"),
	)
	hydra.Conf.API(":19020")

	//1. 删除服务
	os.Args = []string{"xxtest", "remove"}
	app.Start()
	time.Sleep(time.Second * 2)

	//2. 安装服务
	os.Args = []string{"xxtest", "install", "-r", runRegistryAddr, "-c", "c"}
	go app.Start()
	time.Sleep(time.Second * 2)

	//2. 启动服务
	os.Args = []string{"xxtest", "start"}
	go app.Start()

	time.Sleep(time.Second)

	//3. 检查服务状态
	os.Args = []string{"xxtest", "status"}
	go app.Start()

	time.Sleep(time.Second * 2)

	//4. 关闭服务状态
	os.Args = []string{"xxtest", "stop"}
	go app.Start()
	time.Sleep(time.Second * 2)

	//5. 清除服务
	os.Args = []string{"xxtest", "remove"}
	go app.Start()

	time.Sleep(time.Second * 2)
	bytes, err := fileCallback()
	if err != nil {
		t.Error(err)
		return
	}
	line := string(bytes)
	result := strings.Contains(line, "Starting") && strings.Contains(line, "running...")
	assert.Equal(t, true, result, "正常服务运行状态")

	time.Sleep(time.Second)
}

func Test_status_Not_installed(t *testing.T) {

	resetServiceName(t.Name())
	execPrint(t)

	//未安装服务
	defunc, fileCallback := injectStdOutFile()
	defer defunc()

	var app = hydra.NewApp(
		hydra.WithServerTypes(http.API),
		hydra.WithPlatName("xxtest"),
		hydra.WithSystemName("apiserver"),
		hydra.WithClusterName("c"),
	)
	hydra.Conf.API(":19010")

	//1. 清除服务(保证没有服务安装)
	os.Args = []string{"xxtest", "remove"}
	go app.Start()
	time.Sleep(time.Second * 2)

	//2. 启动服务
	os.Args = []string{"xxtest", "status"}
	go app.Start()
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
		assert.Equal(t, true, result, "未安装的服务状态")
	}
	if runtime.GOOS == "windows" {
		result := strings.Contains(line, "not exist as an installed service")
		assert.Equal(t, true, result, "未安装的服务状态")
	}

	time.Sleep(time.Second)
}
