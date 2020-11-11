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

const installRegistryAddr = "lm://."

func Test_install_Normal(t *testing.T) {
	//正常的安装
	fileName := fmt.Sprint(time.Now().Nanosecond())
	file, _ := os.Create(fileName)
	orgStd := *os.Stdout
	*os.Stdout = *file
	defer os.Remove(fileName)
	//正常的安装
	args := []string{"xxtest", "install", "-r", installRegistryAddr, "-c", "c"}

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

	fmt.Println("bytes:", string(bytes), err)

	lines := strings.Split(string(bytes), "\r")
	for _, row := range lines {
		if !strings.Contains(row, "Install") {
			continue
		}
		if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
			//unbuntu/centos
			result := strings.Contains(row, "sudo") || strings.Contains(row, "OK")
			assert.Equal(t, true, result, "正常参数的安装")
		}
		if runtime.GOOS == "windows" {
			result := strings.Contains(row, "OK")
			assert.Equal(t, true, result, "正常参数的安装")
		}
		return
	}
}

func Test_install_Less_param(t *testing.T) {
	//缺少参数的安装 -c
	fileName := fmt.Sprint(time.Now().Nanosecond())
	file, _ := os.Create(fileName)
	orgStd := *os.Stdout
	*os.Stdout = *file
	defer os.Remove(fileName)
	//缺少参数的安装 -c
	args := []string{"xxtest", "install", "-r", installRegistryAddr}

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

	fmt.Println("bytes:", string(bytes), err)

	lines := strings.Split(string(bytes), "\r")
	for _, row := range lines {
		if !strings.Contains(row, "Install") {
			continue
		}
		if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
			//unbuntu/centos
			result := strings.Contains(row, "sudo") || strings.Contains(row, "FAILED")
			assert.Equal(t, true, result, "缺少参数的安装 -c")
		}
		if runtime.GOOS == "windows" {
			result := strings.Contains(row, "FAILED")
			assert.Equal(t, true, result, "缺少参数的安装 -c")
		}
		return
	}
}

func Test_install_Cover(t *testing.T) {
	//覆盖安装 -c
	fileName := fmt.Sprint(time.Now().Nanosecond())
	file, _ := os.Create(fileName)
	orgStd := *os.Stdout
	*os.Stdout = *file
	defer os.Remove(fileName)
	//覆盖安装 -c
	args := []string{"xxtest", "install", "-r", installRegistryAddr, "-c", "c", "-cover", "true"}

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

	fmt.Println("bytes:", string(bytes), err)

	lines := strings.Split(string(bytes), "\r")
	for _, row := range lines {
		if !strings.Contains(row, "Install") {
			continue
		}
		if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
			//unbuntu/centos
			result := strings.Contains(row, "sudo") || strings.Contains(row, "OK")
			assert.Equal(t, true, result, "覆盖安装 -cover=true")
		}
		if runtime.GOOS == "windows" {
			result := strings.Contains(row, "OK")
			assert.Equal(t, true, result, "覆盖安装 -cover=true")
		}
		return
	}
}
