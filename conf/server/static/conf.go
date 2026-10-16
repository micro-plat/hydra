package static

//处理配置目录或压缩包文件，创建到本地，并解压

import (
	"fmt"
	"os"
	"strings"
)

//反归档为文件
func (s *Static) getFileOS() (IFS, error) {
	//检查路径是否配置
	if s.Path == "" {
		return nil, nil
	}
	// HTTP(S) 远程资源由 middleware.downloadFromHTTP 处理，本地无需 os.Stat，
	// 否则 Windows 会把 URL 当本地路径，os.Stat 失败导致服务整体 510。
	if strings.HasPrefix(s.Path, "http://") || strings.HasPrefix(s.Path, "https://") {
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
