package nfs

import (
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/micro-plat/hydra"
)

const (
	DIR = "dir"
)

type fileInfo struct {
	Path    string `json:"path"`
	DPath   string `json:"dpath"`
	Name    string `json:"name"`
	Type    string `json:"type"`
	ModTime string `json:"time,omitempty"`
	Size    int64  `json:"size"`
}

func (f *fileInfo) Copy() *fileInfo {
	n := *f
	return &n
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
	if s[i].ModTime == "" {
		return true
	}
	if s[j].ModTime == "" {
		return false
	}
	left, _ := time.Parse("2006/01/02 15:04:05", s[i].ModTime)
	right, _ := time.Parse("2006/01/02 15:04:05", s[j].ModTime)

	return !left.Before(right)
}

func (s fileList) Swap(i, j int) {
	temp := s[i]
	s[i] = s[j]
	s[j] = temp
}

//GetFileList 获取文件列表
func (l *local) GetFileList(path string, q string, all bool, index int, count int) fileList {
	list := l.getFileList(path, q, all)
	total := index + count
	if index >= list.Len() {
		return make(fileList, 0)
	}
	if total > list.Len() {
		return list[index:]
	}
	return list[index:total]
}

//GetFileList 获以文件列表
func (l *local) getFileList(path string, q string, all bool) fileList {
	list := make(fileList, 0, 1)
	fps := l.GetFPs()
	for k, v := range fps {

		hydra.G.Log().Debug("path:", filepath.Join(path, v.Name), k, path, v.Name)
		if all || !all && k == filepath.Join(path, v.Name) {

			if !strings.Contains(v.Name, q) {
				continue
			}

			list = append(list, &fileInfo{
				Type:    strings.ToLower(strings.Trim(filepath.Ext(k), ".")),
				Path:    k,
				ModTime: v.ModTime.Format("2006/01/02 15:04:05"),
				Size:    v.Size,
				DPath:   strings.ReplaceAll(k, "/", "|"),
				Name:    v.Name,
			})
		}
	}
	sort.Sort(list)
	return list
}
