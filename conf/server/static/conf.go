package static

//处理配置目录或压缩包文件，创建到本地，并解压

import (
	"fmt"
	"os"
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
