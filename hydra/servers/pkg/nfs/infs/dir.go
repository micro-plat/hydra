package infs

type DirInfo struct {
	ID       string     `json:"id"`
	Path     string     `json:"path,omitempty"`
	DPath    string     `json:"dpath,omitempty"`
	PID      string     `json:"pid,omitempty"`
	Name     string     `json:"name"`
	ModTime  string     `json:"time,omitempty"`
	Children []*DirInfo `json:"children,omitempty"`
	Size     int64      `json:"size"`
}

func (f *DirInfo) Copy() *DirInfo {
	n := *f
	return &n
}

type DirList []*DirInfo

func (s DirList) Len() int {
	return len(s)
}

// 在比较的方法中，定义排序的规则
func (s DirList) Less(i, j int) bool {
	return s[i].Path < s[j].Path
}

func (s DirList) Swap(i, j int) {
	temp := s[i]
	s[i] = s[j]
	s[j] = temp
}
func (s DirList) GetMultiLevel(p string) DirList {
	menus := make(DirList, 0, len(s))
	cache := make(map[string]*DirInfo)
	pids := make(map[string]*DirInfo)
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
