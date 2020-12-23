package update

import (
	"os"
	"path/filepath"

	logs "github.com/lib4dev/cli/logger"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/hydra/global/compatible"
	"github.com/micro-plat/lib4go/osext"
	"github.com/urfave/cli"
)

//构建
func doBuild(c *cli.Context) (err error) {

	path, err := osext.Executable()
	if err != nil {
		return err
	}
	p, err := filepath.Abs(global.Version)
	if err != nil {
		return err
	}

	if _, err := os.Stat(p); err == nil && !coverIfExists {
		logs.Log.Errorf("目录已存在%s%s", p, compatible.FAILED)
		return nil
	}

	if err := os.RemoveAll(p); err != nil {
		logs.Log.Errorf("无法移除目录%s%s", p, compatible.FAILED)
		return nil
	}

	if err := Archive(path, p, url); err != nil {
		return err
	}
	logs.Log.Infof("文件已生成到%s%s", p, compatible.SUCCESS)
	return nil
}
