package lnfs

import (
	"sort"
	"strings"

	"github.com/micro-plat/lib4go/security/md5"
)

type dirInfo struct {
	ID       string     `json:"id"`
	Path     string     `json:"path,omitempty"`
	DPath    string     `json:"dpath,omitempty"`
	PID      string     `json:"pid,omitempty"`
	Name     string     `json:"name"`
	ModTime  string     `json:"time,omitempty"`
	Children []*dirInfo `json:"children,omitempty"`
	Size     int64      `json:"size"`
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
	return s[i].Path < s[j].Path
}

func (s dirList) Swap(i, j int) {
	temp := s[i]
	s[i] = s[j]
	s[j] = temp
}
func (s dirList) getMultiLevel(p string) dirList {
	menus := make(dirList, 0, len(s))
	cache := make(map[string]*dirInfo)
	pids := make(map[string]*dirInfo)
	for _, v := range s {

		//缓存每节菜单，用于后续快速查找父节节点
		if _, ok := cache[v.Path]; !ok {
			cache[v.Path] = v
		}

		//生成顶节菜单
		if _, ok := pids[v.Path]; !ok && v.PID == p {
			v.PID = ""
			pids[v.Path] = v
			menus = append(menus, v)
		}

		//非顶节查找，父节节点
		if m, ok := cache[v.PID]; ok {
			v.PID = m.ID
			m.Children = append(m.Children, v)
		}
	}
	return menus
}

//GetDirList 获以文件列表
func (l *local) GetDirList(path string, deep int) dirList {
	list := make(dirList, 0, 1)
	dirs := l.dirs.Items()

	for k, v := range dirs {

		if !strings.HasPrefix(k, path) || k == path {
			continue
		}
		icount := strings.Count(strings.Trim(strings.TrimLeft(k, path), "/"), "/")
		if icount > deep {
			continue
		}

		dir := v.(*eDirEntity)
		list = append(list, &dirInfo{
			ID:      md5.Encrypt(k)[8:16],
			Path:    k,
			ModTime: dir.ModTime.Format("2006/01/02 15:04:05"),
			Size:    dir.Size,
			DPath:   strings.ReplaceAll(k, "/", "|"),
			PID:     getParent(k),
			Name:    dir.Name,
		})
	}
	sort.Sort(list)
	return list.getMultiLevel(path)
}
func getParent(p string) string {
	if p == "" {
		return p
	}
	items := strings.Split(p, "/")
	return strings.Join(items[0:len(items)-1], "/")
}
