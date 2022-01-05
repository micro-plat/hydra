package lnfs

import (
	"sort"
	"strings"

	"github.com/micro-plat/hydra/hydra/servers/pkg/nfs/infs"
	"github.com/micro-plat/lib4go/security/md5"
)

//GetDirList 获以文件列表
func (l *local) GetDirList(path string, deep int) infs.DirList {
	list := make(infs.DirList, 0, 1)
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
		list = append(list, &infs.DirInfo{
			ID:      md5.Encrypt(k)[8:16],
			Path:    k,
			ModTime: dir.ModTime.Format("2006/01/02 15:04:05"),
			Size:    dir.Size,
			DPath:   infs.UnMultiPath(k),
			PID:     getParent(k),
			Name:    dir.Name,
		})
	}
	sort.Sort(list)
	return list.GetMultiLevel(path)
}
func getParent(p string) string {
	if p == "" {
		return p
	}
	items := strings.Split(p, "/")
	return strings.Join(items[0:len(items)-1], "/")
}
