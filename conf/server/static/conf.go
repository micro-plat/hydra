package static

//处理配置目录或压缩包文件，创建到本地，并解压

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/mholt/archiver"
	"github.com/micro-plat/hydra/global"
)

//反归档为文件
func (s *Static) getFileOS() (IFS, error) {
	//检查路径是否配置
	if s.Path == "" {
		return nil, nil
	}
	fs, err := os.Stat(s.Path)
	if err != nil {
		return nil, fmt.Errorf("无法读取文件%s,%w", s.Path, err)
	}

	//配置为目录
	if fs.IsDir() {
		return newOSFS(s.Path), nil
	}
	return unarchive(s.Path)

}

func unarchive(path string) (IFS, error) {
	//非目录则按压缩包方式解压
	rootPath := filepath.Dir(os.Args[0])
	tmpDir, err := ioutil.TempDir(rootPath, TempDirName)
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
