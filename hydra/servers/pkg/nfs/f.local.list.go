package nfs

import (
	"path/filepath"
	"sort"
	"strings"
)

const (
	DIR = "dir"
)

type fileInfo struct {
	Path    string `json:"path"`
	Name    string `json:"name"`
	Type    string `json:"type"`
	ModTime string `json:"time,omitempty"`
	Size    int64  `json:"size"`
}
type fileList []*fileInfo

func (s fileList) Len() int {
	return len(s)
}

// 在比较的方法中，定义排序的规则
func (s fileList) Less(i, j int) bool {
	if s[i].Type == DIR {
		return true
	} else if s[j].Type == DIR {
		return false
	}
	return s[i].ModTime < s[i].ModTime
}

func (s fileList) Swap(i, j int) {
	temp := s[i]
	s[i] = s[j]
	s[j] = temp
}

//GetFileList 获以文件列表
func (l *local) GetFileList(q string) fileList {
	list := make(fileList, 0, 1)
	prefix := strings.Trim(q, "/")
	fps := l.GetFPs()
	for k, v := range fps {
		if prefix != "" && !strings.HasPrefix(k, prefix) {
			continue
		}
		path := strings.Trim(strings.Trim(k, prefix), "/")
		if strings.Contains(path, "/") {
			list = append(list, &fileInfo{
				Type: DIR,
				Path: k,
				Size: v.Size,
				Name: strings.Split(path, "/")[0]})
			continue
		}
		list = append(list, &fileInfo{
			Type:    strings.ToLower(strings.Trim(filepath.Ext(k), ".")),
			Path:    k,
			ModTime: v.ModTime.Format("2006/01/02 15:04:05"),
			Size:    v.Size,
			Name:    strings.Trim(path, filepath.Ext(k))})

	}
	sort.Sort(list)
	return list
}
