package static

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

//GetGzFile 获取gz 压缩包
func (s *Static) GetGzFile(rPath string) string {

	fi, ok := s.FileMap[rPath]
	if !ok {
		return ""
	}
	if !fi.HasGz {
		return ""
	}
	return fi.GzFile
}

//RereshData 刷新配置数据
func (s *Static) RereshData() {
	s.recursiveDir(strings.TrimPrefix(s.Dir, "./"))
}

//GetFileMap 获取静态文件中的gz文件字典
func (s *Static) GetFileMap() map[string]FileInfo {
	return s.FileMap
}

func (s *Static) recursiveDir(dir string) {
	children := []string{}
	const suffix = ".gz"
	list, err := ioutil.ReadDir(dir)
	if err != nil {
		return
	}
	var cur os.FileInfo
	for i := range list {
		cur = list[i]
		if !cur.IsDir() {
			if strings.HasSuffix(cur.Name(), suffix) {
				fpath := fmt.Sprintf("%s/%s", dir, strings.TrimSuffix(cur.Name(), suffix))
				s.FileMap[fpath] = FileInfo{
					HasGz:  true,
					GzFile: fpath + suffix,
				}
			}
			continue
		}
		children = append(children, fmt.Sprintf("%s/%s", dir, cur.Name()))
	}

	for i := range children {
		s.recursiveDir(children[i])
	}
	return
}
