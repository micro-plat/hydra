package static

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/mholt/archiver"
	"github.com/micro-plat/hydra/global"
)

func unarchive(path string) (*osfs, error) {
	//非目录则按压缩包方式解压
	//rootPath := filepath.Dir(os.Args[0])
	tmpDir, err := ioutil.TempDir(".", TempDirName)
	if err != nil {
		return nil, fmt.Errorf("创建临时文件失败:%v", err)
	}
	err = archiver.Unarchive(path, tmpDir)
	if err != nil {
		return nil, fmt.Errorf("指定的文件%s解压失败:%v", path, err)
	}
	waitRemoveDir = append(waitRemoveDir, tmpDir)
	return newOSFS(tmpDir), nil
}

//文件删除处理
var waitRemoveDir = make([]string, 0, 1)

func init() {
	global.Def.AddCloser(func() error {
		for _, d := range waitRemoveDir {
			os.RemoveAll(d)
		}
		return nil
	})
}
