package nfs

import (
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type dirInfo struct {
	Path    string `json:"path"`
	DPath   string `json:"dpath"`
	Name    string `json:"name"`
	Type    string `json:"type"`
	ModTime string `json:"time,omitempty"`
	Size    int64  `json:"size"`
}

func (f *dirInfo) Copy() *dirInfo {
	n := *f
	return &n
}

type dirList []*dirInfo

func (s dirList) Len() int {
	return len(s)
}

// 在比较的方法中，定义排序的规则
func (s dirList) Less(i, j int) bool {
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

func (s dirList) Swap(i, j int) {
	temp := s[i]
	s[i] = s[j]
	s[j] = temp
}

//GetDirList 获以文件列表
func (l *local) GetDirList(path string, deep int) dirList {
	list := make(dirList, 0, 1)
	dirs := l.dirs.Items()
	for k, v := range dirs {
		if !strings.HasPrefix(k, path) {
			continue
		}
		icount := strings.Count(strings.Trim(strings.TrimLeft(k, path), "/"), "/")
		if icount >= deep {
			continue
		}

		dir := v.(*eDirEntity)
		list = append(list, &dirInfo{
			Type:    strings.ToLower(strings.Trim(filepath.Ext(k), ".")),
			Path:    k,
			ModTime: dir.ModTime.Format("2006/01/02 15:04:05"),
			Size:    dir.Size,
			DPath:   strings.ReplaceAll(k, "/", "|"),
			Name:    dir.Name,
		})

	}
	sort.Sort(list)
	return list
}
