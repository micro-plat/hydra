package cmds

import (
	"fmt"
	"testing"

	"github.com/micro-plat/hydra/global"
)

func resetServiceName(name string) {
	global.AppName = "test12345678901234567890" + name
}

func execPrint(t *testing.T) {
	cmdVal := fmt.Sprintf(`unbuntu使用命令执行测试：sudo /usr/local/go/bin/go test github.com/micro-plat/hydra/test/hydra/cmds -run "^(%s)$"  -v -timeout=60s`, t.Name())
	fmt.Println(cmdVal)
}
