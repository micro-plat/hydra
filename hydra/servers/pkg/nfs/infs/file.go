package infs

import "time"

type FileInfo struct {
	Path    string `json:"path"`
	DPath   string `json:"dpath"`
	Name    string `json:"name"`
	Type    string `json:"type"`
	ModTime string `json:"time,omitempty"`
	Size    int64  `json:"size"`
}

func (f *FileInfo) Copy() *FileInfo {
	n := *f
	return &n
}

type FileList []*FileInfo

func (s FileList) Len() int {
	return len(s)
}

// 在比较的方法中，定义排序的规则
func (s FileList) Less(i, j int) bool {
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

func (s FileList) Swap(i, j int) {
	temp := s[i]
	s[i] = s[j]
	s[j] = temp
}
