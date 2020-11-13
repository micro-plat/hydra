package cmds

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/micro-plat/hydra/global"
)

func resetServiceName(name string) {
	global.AppName = "test12345678901234567890" + name
}

func execPrint(t *testing.T) {
	cmdVal := fmt.Sprintf(`ubuntu使用命令执行测试：sudo /usr/local/go/bin/go test github.com/micro-plat/hydra/test/hydra/cmds -run "^(%s)$"  -v -timeout=60s`, t.Name())
	fmt.Println(cmdVal)
}

func injectStdOutFile() (func(), func() (string, error)) {
	fileName := fmt.Sprint("cmds_test" + time.Now().Nanosecond())
	file, _ := os.Create(fileName)
	orgStd := *os.Stdout
	*os.Stdout = *file
	return func() {
			os.Remove(fileName)
		},
		func() (string, error) {
			*os.Stdout = orgStd
			file.Close()
			bytes, err := ioutil.ReadFile(fileName)
			return string(bytes), err
		}
}
